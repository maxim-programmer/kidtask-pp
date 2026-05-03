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

const adminApi = axios.create({ baseURL: '/api/admin' })
adminApi.interceptors.request.use(cfg => {
  cfg.headers['X-Admin-Secret'] = localStorage.getItem('kt_admin_secret') || ''
  return cfg
})

export function useApi() {
  const register = (data) => api.post('/auth/register', data)
  const login = (data) => api.post('/auth/login', data)
  const loginChild = (data) => api.post('/auth/child/login', data)
  const getMe = () => api.get('/me')
  const updateMe = (data) => api.patch('/me', data)

  const getChildren = () => api.get('/children')
  const createChild = (data) => api.post('/children', data)
  const updateChild = (id, data) => api.patch(`/children/${id}`, data)
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

  const getSupportChat = (parentId) => api.get('/support/chat', { params: { parent_id: parentId } })
  const sendSupportMessage = (parentId, body) => api.post('/support/chat', { parent_id: parentId, body })
  const createComplaint = (parentId, subject, body) => api.post('/support/complaints', { parent_id: parentId, subject, body })

  return {
    register, login, loginChild, getMe, updateMe,
    getChildren, createChild, updateChild, deleteChild,
    getTasks, createTask, updateTask, deleteTask, submitTask, approveTask, rejectTask,
    getWishes, createWish, updateWish, deleteWish, purchaseWish, deliverWish,
    getStats,
    getSupportChat, sendSupportMessage, createComplaint,
  }
}

export function useAdminApi() {
  const getStats = () => adminApi.get('/stats')
  const getFamilies = () => adminApi.get('/families')
  const blockFamily = (id) => adminApi.post(`/families/${id}/block`)
  const unblockFamily = (id) => adminApi.post(`/families/${id}/unblock`)
  const deleteFamily = (id) => adminApi.delete(`/families/${id}`)

  const getChildren = () => adminApi.get('/children')
  const blockChild = (id) => adminApi.post(`/children/${id}/block`)
  const unblockChild = (id) => adminApi.post(`/children/${id}/unblock`)
  const adjustBalance = (id, delta, reason) => adminApi.post(`/children/${id}/balance`, { delta, reason })
  const getBalanceLogs = (id) => adminApi.get(`/children/${id}/logs`)

  const getComplaints = () => adminApi.get('/complaints')
  const resolveComplaint = (id) => adminApi.post(`/complaints/${id}/resolve`)

  const getChatParents = () => adminApi.get('/chat')
  const getChatMessages = (parentId) => adminApi.get(`/chat/${parentId}`)
  const sendChatMessage = (parentId, body) => adminApi.post(`/chat/${parentId}`, { body })

  const getWishes = (sort) => adminApi.get('/wishes', { params: { sort } })

  return {
    getStats, getFamilies, blockFamily, unblockFamily, deleteFamily,
    getChildren, blockChild, unblockChild, adjustBalance, getBalanceLogs,
    getComplaints, resolveComplaint,
    getChatParents, getChatMessages, sendChatMessage,
    getWishes,
  }
}