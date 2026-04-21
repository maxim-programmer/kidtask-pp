import axios from 'axios'

const api = axios.create({ baseURL: '/api' })

api.interceptors.request.use(cfg => {
  const token = localStorage.getItem('kt_token')
  if (token) cfg.headers.Authorization = `Bearer ${token}`
  return cfg
})

api.interceptors.response.use(
  r => r,
  err => {
    if (err.response?.status === 401) {
      localStorage.clear()
      window.location.href = '/login'
    }
    return Promise.reject(err)
  }
)

export function useApi() {
  const register = (data) => api.post('/auth/register', data)
  const login = (data) => api.post('/auth/login', data)
  const loginChild = (data) => api.post('/auth/child/login', data)
  const getMe = () => api.get('/me')

  const getChildren = () => api.get('/children')
  const createChild = (data) => api.post('/children', data)
  const deleteChild = (id) => api.delete(`/children/${id}`)

  const getTasks = (params) => api.get('/tasks', { params })
  const createTask = (data) => api.post('/tasks', data)
  const updateTask = (id, data) => api.patch(`/tasks/${id}`, data)
  const deleteTask = (id) => api.delete(`/tasks/${id}`)
  const submitTask = (id) => api.post(`/tasks/${id}/submit`)
  const approveTask = (id) => api.post(`/tasks/${id}/approve`)
  const rejectTask = (id, comment) => api.post(`/tasks/${id}/reject`, { comment })

  const getWishes = (childId, params) => api.get(`/children/${childId}/wishes`, { params })
  const createWish = (childId, data) => api.post(`/children/${childId}/wishes`, data)
  const updateWish = (childId, id, data) => api.patch(`/children/${childId}/wishes/${id}`, data)
  const deleteWish = (childId, id) => api.delete(`/children/${childId}/wishes/${id}`)
  const purchaseWish = (childId, id) => api.post(`/children/${childId}/wishes/${id}/purchase`)
  const deliverWish = (childId, id) => api.patch(`/children/${childId}/wishes/${id}/deliver`)

  const getStats = () => api.get('/stats')

  return {
    register, login, loginChild, getMe,
    getChildren, createChild, deleteChild,
    getTasks, createTask, updateTask, deleteTask, submitTask, approveTask, rejectTask,
    getWishes, createWish, updateWish, deleteWish, purchaseWish, deliverWish,
    getStats,
  }
}