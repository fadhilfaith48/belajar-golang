import { useEffect, useState } from 'react'
import './App.css'

function App() {
  const [tasks, setTasks] = useState([])
  const [stats, setStats] = useState({})
  const [filter, setFilter] = useState('')
  const [search, setSearch] = useState('')
  const [page, setPage] = useState(1)
  const [info, setInfo] = useState({ total: 0, total_pages: 1 })
  const [loading, setLoading] = useState(true)
  const [pesan, setPesan] = useState('')
  const [form, setForm] = useState({ judul: '', kategori: '', prioritas: 'sedang', deadline: '' })
  const [prioritas, setPrioritas] = useState('')
  const [sort, setSort] = useState('id_asc')

  const limit = 10

  // Ambil daftar task + statistik dari API
  const ambilData = async () => {
    setLoading(true)
    try {
      // URLSearchParams otomatis meng-encode nilai,
      // jadi spasi otomatis jadi %20 dan tidak merusak URL.
      const params = new URLSearchParams()
      if (filter) params.set('status', filter)
      if (search.trim()) params.set('search', search.trim())
      if (prioritas) params.set('prioritas', prioritas)
      if (sort) {
        const [sb, so] = sort.split('_')
        if (sb) params.set('sort_by', sb)
        if (so) params.set('order', so)
      }
      params.set('page', page)
      params.set('limit', limit)

      const res = await fetch(`/api/tasks?${params}`)
      const hasil = await res.json()
      setTasks(hasil.data)
      setInfo({ total: hasil.total, total_pages: hasil.total_pages })

      const sres = await fetch('/api/tasks/stats')
      setStats(await sres.json())
    } catch (err) {
      setPesan('Gagal memuat data: ' + err.message)
    } finally {
      setLoading(false)
    }
  }

  // Debounce: tunggu 300ms setelah berhenti mengetik baru request.
  // Tanpa ini, tiap ketikan huruf akan memicu satu request ke server.
  useEffect(() => {
    const timer = setTimeout(ambilData, 300)
    return () => clearTimeout(timer)
  }, [filter, search, page])

  // Ganti filter/search selalu kembali ke halaman 1,
  // supaya tidak terjebak di halaman 5 padahal hasilnya cuma 1 halaman.
  const gantiFilter = (kode) => {
    setFilter(kode)
    setPage(1)
  }

  const gantiSearch = (nilai) => {
    setSearch(nilai)
    setPage(1)
  }
  const gantiPrioritas = (v) => { setPrioritas(v); setPage(1) }
  const gantiSort = (v) => { setSort(v); setPage(1) }

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
      setForm({ judul: '', kategori: '', prioritas: 'sedang', deadline: '' })
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
        <input
          type="date"
          value={form.deadline}
          onChange={(e) => setForm({ ...form, deadline: e.target.value })}
        />
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
            onClick={() => gantiFilter(f.kode)}
          >
            {f.label}
          </button>
        ))}
      </div>

      <div className="cari">
        <input
          placeholder="Cari judul tugas..."
          value={search}
          onChange={(e) => gantiSearch(e.target.value)}
        />
        <span>{info.total} tugas ditemukan</span>
      </div>

      <div className="filter">
        <span>Filter prioritas: </span>
        {[
          { kode: '', label: 'Semua' },
          { kode: 'rendah', label: 'Rendah' },
          { kode: 'sedang', label: 'Sedang' },
          { kode: 'tinggi', label: 'Tinggi' },
        ].map((p) => (
          <button key={p.kode} className={prioritas === p.kode ? 'aktif' : ''} onClick={() => gantiPrioritas(p.kode)}>{p.label}</button>
        ))}
        <select value={sort} onChange={(e) => gantiSort(e.target.value)} style={{ marginLeft: 8 }}>
          <option value="id_asc">ID A-Z</option>
          <option value="id_desc">ID Z-A</option>
          <option value="judul_asc">Judul A-Z</option>
          <option value="judul_desc">Judul Z-A</option>
          <option value="deadline_asc">Deadline Terdekat</option>
          <option value="deadline_desc">Deadline Terjauh</option>
          <option value="prioritas_desc">Prioritas Tinggi→Rendah</option>
          <option value="prioritas_asc">Prioritas Rendah→Tinggi</option>
        </select>
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
              <th>Deadline</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            {tasks.length === 0 ? (
              <tr>
                <td colSpan="7">Belum ada tugas.</td>
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
                  <td>{t.deadline || '-'}</td>
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

      <div className="pagination">
        <button
          disabled={page <= 1}
          onClick={() => setPage((p) => p - 1)}
        >
          ‹ Sebelumnya
        </button>
        <span>
          Halaman {page} dari {info.total_pages}
        </span>
        <button
          disabled={page >= info.total_pages}
          onClick={() => setPage((p) => p + 1)}
        >
          Berikutnya ›
        </button>
      </div>
    </div>
  )
}

export default App
