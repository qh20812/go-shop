import { useEffect, useState } from 'react'
import axios from 'axios'
import './App.css'

const emptyForm = { name: '', price: '', description: '' }

function App() {
  const [products, setProducts] = useState([])
  const [form, setForm] = useState(emptyForm)
  const [imageFile, setImageFile] = useState(null)
  const [editingId, setEditingId] = useState(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    fetchProducts()
  }, [])

  async function fetchProducts() {
    const res = await axios.get('/api/products')
    setProducts(res.data)
  }

  function handleChange(e) {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setLoading(true)
    try {
      const data = new FormData()
      data.append('name', form.name)
      data.append('price', form.price)
      data.append('description', form.description)
      if (imageFile) data.append('image', imageFile)

      if (editingId) {
        await axios.put('/api/products/' + editingId, data)
      } else {
        await axios.post('/api/products', data)
      }
      resetForm()
      fetchProducts()
    } catch (err) {
      alert(err.response?.data?.message || 'Có lỗi xảy ra')
    } finally {
      setLoading(false)
    }
  }

  function startEdit(product) {
    setEditingId(product.id)
    setForm({
      name: product.name,
      price: product.price,
      description: product.description
    })
    setImageFile(null)
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  async function handleDelete(id) {
    if (!confirm('Xoá sản phẩm này?')) return
    await axios.delete('/api/products/' + id)
    fetchProducts()
  }

  function resetForm() {
    setForm(emptyForm)
    setImageFile(null)
    setEditingId(null)
    document.getElementById('imageInput').value = ''
  }

  return (
    <div className="container">
      <h1>🛒 Go Shop Quản lý sản phẩm</h1

      <form className="card" onSubmit={handleSubmit}>
        <h2>{editingId ? 'Sửa sản phẩm' : 'Thêm sản phẩm'}</h2>
        <input
          name="name"
          placeholder="Tên sản phẩm"
          value={form.name}
          onChange={handleChange}
          required
        />
        <input
          name="price"
          type="number"
          min="0"
          placeholder="Giá (VND)"
          value={form.price}
          onChange={handleChange}
          required
        />
        <textarea
          name="description"
          placeholder="Mô tả"
          value={form.description}
          onChange={handleChange}
        />
        <input
          id="imageInput"
          type="file"
          accept="image/*"
          onChange={(e) => setImageFile(e.target.files[0])}
        />
        <div className="actions">
          <button type="submit" disabled={loading}>
            {editingId ? 'Cập nhật' : 'Thêm mới'}
          </button>
          {editingId && (
            <button type="button" onClick={resetForm}>
              Huỷ
            </button>
          )}
        </div>
      </form>

      <div className="grid">
        {products.map((product) => (
          <div className="card product" key={product.id}>
            {product.image ? (
              <img src={product.image} alt={product.name} />
            ) : (
              <div className="no-image">Chưa có ảnh</div>
            )}
            <h3>{product.name}</h3>
            <p className="price">
              {Number(product.price).toLocaleString('vi-VN')} ₫
            </p>
            <p>{product.description}</p>
            <div className="actions">
              <button onClick={() => startEdit(product)}>Sửa</button>
              <button className="danger" onClick={() => handleDelete(product.id)}>
                Xoá
              </button>
            </div>
          </div>
        ))}
        {products.length === 0 && (
          <p>Chưa có sản phẩm nào — thêm sản phẩm đầu tiên ở form trên 👆</p>
        )}
      </div>
    </div>
  )
}

export default App
