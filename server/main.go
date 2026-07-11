package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Đọc file .env nếu có (trên VPS, systemd sẽ nạp env qua EnvironmentFile)
	godotenv.Load()

	if err := connectDB(); err != nil {
		log.Fatal("❌ Không kết nối được PostgreSQL: ", err)
	}
	log.Println("✅ Đã kết nối PostgreSQL")

	// Tự tạo thư mục uploads nếu chưa có
	os.MkdirAll("uploads", 0755)

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20 // giới hạn upload 8MB

	// Serve ảnh đã upload: /uploads/xxx.jpg
	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		api.GET("/health", healthCheck)
		api.GET("/products", listProducts)
		api.GET("/products/:id", getProduct)
		api.POST("/products", createProduct)
		api.PUT("/products/:id", updateProduct)
		api.DELETE("/products/:id", deleteProduct)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("🚀 API chạy tại http://localhost:" + port)
	r.Run(":" + port)
}