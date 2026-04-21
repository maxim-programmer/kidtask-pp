import { ref, computed } from 'vue'

const token = ref(localStorage.getItem('kt_token'))
const role = ref(localStorage.getItem('kt_role'))
const user = ref(JSON.parse(localStorage.getItem('kt_user') || 'null'))
const ageGroup = ref(localStorage.getItem('kt_age_group'))

export function useAuth() {
  const isLoggedIn = computed(() => !!token.value)

  function saveAuth(data) {
    token.value = data.token
    role.value = data.role

    const u = data.user || {}
    // Нормализуем child_id — бэкенд отдаёт child_id в объекте child
    if (data.role === 'child' && !u.child_id && u.child_id !== 0) {
      u.child_id = u.child_id || u.id || u.user_id
    }

    user.value = u
    localStorage.setItem('kt_token', data.token)
    localStorage.setItem('kt_role', data.role)
    localStorage.setItem('kt_user', JSON.stringify(u))

    const ag = u.age_group
    if (ag) {
      ageGroup.value = ag
      localStorage.setItem('kt_age_group', ag)
    }
  }

  function logout() {
    token.value = null
    role.value = null
    user.value = null
    ageGroup.value = null
    localStorage.removeItem('kt_token')
    localStorage.removeItem('kt_role')
    localStorage.removeItem('kt_user')
    localStorage.removeItem('kt_age_group')
  }

  return { token, role, user, ageGroup, isLoggedIn, saveAuth, logout }
}