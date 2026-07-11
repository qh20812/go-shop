package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Product — struct ánh xạ 1 dòng trong bảng products
type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Price       int64     `json:"price"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"createdAt"`
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now()})
}

// ---------- Xử lý file ảnh ----------

var unsafeChars = regexp.MustCompile(`[^a-z0-9.]`)

// sanitizeFilename: "My Photo (1).PNG" → "my-photo--1-.png"
// filepath.Base chặn kiểu tấn công "../../etc/passwd"
func sanitizeFilename(name string) string {
	return unsafeChars.ReplaceAllString(strings.ToLower(filepath.Base(name)), "-")
}

// saveImage lưu ảnh vào uploads/, trả về "/uploads/xxx" (chuỗi rỗng nếu request không kèm ảnh)
func saveImage(c *gin.Context) (string, error) {
	file, err := c.FormFile("image")
	if err != nil {
		return "", nil // không gửi ảnh → không phải lỗi
	}
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return "", fmt.Errorf("chỉ cho phép upload file ảnh")
	}
	name := fmt.Sprintf("%d-%s", time.Now().UnixMilli(), sanitizeFilename(file.Filename))
	if err := c.SaveUploadedFile(file, filepath.Join("uploads", name)); err != nil {
		return "", err
	}
	return "/uploads/" + name, nil
}

// removeImage xoá file ảnh trên đĩa (bỏ qua nếu không tồn tại)
func removeImage(imagePath string) {
	if imagePath != "" {
		os.Remove("." + imagePath) // "/uploads/x.jpg" → "./uploads/x.jpg"
	}
}

// ---------- CRUD ----------

// GET /api/products — danh sách (mới nhất trước)
func listProducts(c *gin.Context) {
	rows, err := db.Query(`SELECT id, name, price, description, image, created_at
	                       FROM products ORDER BY created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer rows.Close()

	products := []Product{} // khởi tạo slice rỗng để JSON trả [] thay vì null
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Image, &p.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		products = append(products, p)
	}
	c.JSON(http.StatusOK, products)
}

// GET /api/products/:id — chi tiết
func getProduct(c *gin.Context) {
	var p Product
	err := db.QueryRow(`SELECT id, name, price, description, image, created_at
	                    FROM products WHERE id = $1`, c.Param("id")).
		Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Image, &p.CreatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Không tìm thấy sản phẩm"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

// POST /api/products — multipart/form-data: name, price, description, image (file)
func createProduct(c *gin.Context) {
	name := c.PostForm("name")
	price, err := strconv.ParseInt(c.PostForm("price"), 10, 64)
	if name == "" || err != nil || price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "name và price (số >= 0) là bắt buộc"})
		return
	}

	image, err := saveImage(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var p Product
	err = db.QueryRow(
		`INSERT INTO products (name, price, description, image)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, name, price, description, image, created_at`,
		name, price, c.PostForm("description"), image).
		Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Image, &p.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

// PUT /api/products/:id — cập nhật (ảnh mới là tuỳ chọn)
func updateProduct(c *gin.Context) {
	// 1. Lấy sản phẩm hiện tại
	var current Product
	err := db.QueryRow(`SELECT id, name, price, description, image
	                    FROM products WHERE id = $1`, c.Param("id")).
		Scan(&current.ID, &current.Name, &current.Price, &current.Description, &current.Image)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Không tìm thấy sản phẩm"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 2. Ghi đè các field được gửi lên
	if v := c.PostForm("name"); v != "" {
		current.Name = v
	}
	if v := c.PostForm("price"); v != "" {
		if price, err := strconv.ParseInt(v, 10, 64); err == nil && price >= 0 {
			current.Price = price
		}
	}
	if v, ok := c.GetPostForm("description"); ok {
		current.Description = v
	}

	// 3. Có ảnh mới → lưu ảnh mới, xoá ảnh cũ
	newImage, err := saveImage(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if newImage != "" {
		removeImage(current.Image)
		current.Image = newImage
	}

	// 4. Ghi xuống DB
	var p Product
	err = db.QueryRow(
		`UPDATE products SET name = $1, price = $2, description = $3, image = $4
		 WHERE id = $5
		 RETURNING id, name, price, description, image, created_at`,
		current.Name, current.Price, current.Description, current.Image, current.ID).
		Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Image, &p.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

// DELETE /api/products/:id — xoá sản phẩm + xoá luôn file ảnh
func deleteProduct(c *gin.Context) {
	var image string
	err := db.QueryRow(`DELETE FROM products WHERE id = $1 RETURNING image`,
		c.Param("id")).Scan(&image)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Không tìm thấy sản phẩm"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	removeImage(image)
	c.JSON(http.StatusOK, gin.H{"message": "Đã xoá", "id": c.Param("id")})
}