import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import HomePage from './views/HomePage.vue'
import ParentRegister from './views/ParentRegister.vue'
import LoginPage from './views/LoginPage.vue'
import ParentDashboard from './views/parent/ParentDashboard.vue'
import ParentChildren from './views/parent/ParentChildren.vue'
import ParentTasks from './views/parent/ParentTasks.vue'
import ParentSettings from './views/parent/ParentSettings.vue'
import ChildDashboardJunior from './views/child/junior/ChildDashboard.vue'
import ChildTasksJunior from './views/child/junior/ChildTasks.vue'
import ChildWishlistJunior from './views/child/junior/ChildWishlist.vue'
import ChildDashboardSenior from './views/child/senior/ChildDashboard.vue'
import ChildTasksSenior from './views/child/senior/ChildTasks.vue'
import ChildWishlistSenior from './views/child/senior/ChildWishlist.vue'
import AdminLogin from './views/admin/AdminLogin.vue'
import AdminDashboard from './views/admin/AdminDashboard.vue'

const routes = [
  { path: '/', component: HomePage },
  { path: '/register', component: ParentRegister },
  { path: '/login', component: LoginPage },
  { path: '/parent/dashboard', component: ParentDashboard, meta: { requiresAuth: true, role: 'parent' } },
  { path: '/parent/children', component: ParentChildren, meta: { requiresAuth: true, role: 'parent' } },
  { path: '/parent/tasks', component: ParentTasks, meta: { requiresAuth: true, role: 'parent' } },
  { path: '/parent/settings', component: ParentSettings, meta: { requiresAuth: true, role: 'parent' } },
  { path: '/child/junior/dashboard', component: ChildDashboardJunior, meta: { requiresAuth: true, role: 'child', ageGroup: 'junior' } },
  { path: '/child/junior/tasks', component: ChildTasksJunior, meta: { requiresAuth: true, role: 'child', ageGroup: 'junior' } },
  { path: '/child/junior/wishlist', component: ChildWishlistJunior, meta: { requiresAuth: true, role: 'child', ageGroup: 'junior' } },
  { path: '/child/senior/dashboard', component: ChildDashboardSenior, meta: { requiresAuth: true, role: 'child', ageGroup: 'senior' } },
  { path: '/child/senior/tasks', component: ChildTasksSenior, meta: { requiresAuth: true, role: 'child', ageGroup: 'senior' } },
  { path: '/child/senior/wishlist', component: ChildWishlistSenior, meta: { requiresAuth: true, role: 'child', ageGroup: 'senior' } },
  { path: '/admin/login', component: AdminLogin },
  { path: '/admin', component: AdminDashboard, meta: { requiresAdmin: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  if (to.meta.requiresAdmin) {
    const saved = localStorage.getItem('kt_admin_secret')
    if (!saved) {
      next('/admin/login')
      return
    }
    next()
    return
  }

  const token = localStorage.getItem('kt_token')
  const role = localStorage.getItem('kt_role')
  const ageGroup = localStorage.getItem('kt_age_group')

  if (to.meta.requiresAuth && !token) {
    next('/login')
    return
  }
  if (to.meta.role && to.meta.role !== role) {
    if (role === 'parent') next('/parent/dashboard')
    else if (role === 'child') {
      const group = ageGroup || 'junior'
      next(`/child/${group}/dashboard`)
    } else next('/login')
    return
  }
  next()
})

createApp(App).use(router).mount('#app')