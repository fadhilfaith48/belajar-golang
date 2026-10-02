import { useEffect, useState } from 'react'
import './App.css'

function App() {
  const [tasks, setTasks] = useState([])
  const [stats, setStats] = useState({})
  const [filter, setFilter] = useState('')
  const [loading, setLoading] = useState(true)
  const [pesan, setPesan] = useState('')
  const [form, setForm] = useState({ judul: '', kategori: '', prioritas: 'sedang' })

  // Ambil daftar task + statistik dari API
  const ambilData = async () => {
    setLoading(true)
    try {
      const q = filter ? `?status=${filter}` : ''
      const res = await fetch(`/api/tasks${q}`)
      const data = await res.json()
      setTasks(data)

      const sres = await fetch('/api/tasks/stats')
      setStats(await sres.json())
    } catch (err) {
      setPesan('Gagal memuat data: ' + err.message)
    } finally {
      setLoading(false)
    }
  }

  // Jalankan ulang setiap filter berubah
  useEffect(() => {
    ambilData()
  }, [filter])

  const tambahTask = async (e) => {
    e.preventDefault()
    if (!form.judul.trim()) {
      setPesan('Judul tidak boleh kosong')
      return
    }
    setPesan('')
    try {
      const res = await fetch('/api/tasks', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form),
      })
      if (!res.ok) {
        setPesan('Gagal menambah: ' + (await res.text()))
        return
      }
      setForm({ judul: '', kategori: '', prioritas: 'sedang' })
      ambilData()
    } catch (err) {
      setPesan('Error: ' + err.message)
    }
  }

  const hapusTask = async (id) => {
    await fetch(`/api/tasks/${id}`, { method: 'DELETE' })
    ambilData()
  }

  const ubahStatus = async (t, status) => {
    await fetch(`/api/tasks/${t.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ...t, status }),
    })
    ambilData()
  }

  return (
    <div className="container">
      <h1>Manajemen Tugas</h1>

      <div className="stats">
        <span className="stat">Todo: {stats.todo ?? 0}</span>
        <span className="stat">Doing: {stats.doing ?? 0}</span>
        <span className="stat">Done: {stats.done ?? 0}</span>
        <span className="stat">Total request: {stats.total_request ?? 0}</span>
      </div>

      <form className="form" onSubmit={tambahTask}>
        <input
          placeholder="Judul tugas"
          value={form.judul}
          onChange={(e) => setForm({ ...form, judul: e.target.value })}
        />
        <input
          placeholder="Kategori (mis. belajar)"
          value={form.kategori}
          onChange={(e) => setForm({ ...form, kategori: e.target.value })}
        />
        <select
          value={form.prioritas}
          onChange={(e) => setForm({ ...form, prioritas: e.target.value })}
        >
          <option value="rendah">Prioritas: Rendah</option>
          <option value="sedang">Prioritas: Sedang</option>
          <option value="tinggi">Prioritas: Tinggi</option>
        </select>
        <button type="submit">Tambah</button>
      </form>

      {pesan && <p className="pesan">{pesan}</p>}

      <div className="filter">
        <span>Filter status: </span>
        {[
          { kode: '', label: 'Semua' },
          { kode: 'todo', label: 'Todo' },
          { kode: 'doing', label: 'Doing' },
          { kode: 'done', label: 'Done' },
        ].map((f) => (
          <button
            key={f.kode}
            className={filter === f.kode ? 'aktif' : ''}
            onClick={() => setFilter(f.kode)}
          >
            {f.label}
          </button>
        ))}
      </div>

      {loading ? (
        <p>Memuat...</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Judul</th>
              <th>Kategori</th>
              <th>Prioritas</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            {tasks.length === 0 ? (
              <tr>
                <td colSpan="6">Belum ada tugas.</td>
              </tr>
            ) : (
              tasks.map((t) => (
                <tr key={t.id}>
                  <td>{t.id}</td>
                  <td>{t.judul}</td>
                  <td>{t.kategori}</td>
                  <td>{t.prioritas}</td>
                  <td>
                    <select value={t.status} onChange={(e) => ubahStatus(t, e.target.value)}>
                      <option value="todo">Todo</option>
                      <option value="doing">Doing</option>
                      <option value="done">Done</option>
                    </select>
                  </td>
                  <td>
                    <button className="hapus" onClick={() => hapusTask(t.id)}>
                      Hapus
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      )}
    </div>
  )
}

export default App
