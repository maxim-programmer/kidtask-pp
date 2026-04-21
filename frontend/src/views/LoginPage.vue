<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="logo">kid<span class="accent">TASK</span></div>

      <div class="role-select">
        <p class="role-label">Кто вы?</p>
        <div class="role-btns">
          <button :class="['role-btn', { 'role-btn--active': mode === 'child' }]" @click="mode = 'child'">Ребёнок</button>
          <button :class="['role-btn', { 'role-btn--active': mode === 'parent' }]" @click="mode = 'parent'">Родитель</button>
        </div>
      </div>

      <h2 class="form-title">Вход</h2>

      <form @submit.prevent="login">
        <div class="field" v-if="mode === 'child'">
          <label>Логин</label>
          <input v-model="form.username" type="text" :placeholder="mode === 'child' ? 'ivan123' : ''" />
        </div>
        <div class="field" v-if="mode === 'parent'">
          <label>Почта</label>
          <input v-model="form.email" type="email" placeholder="example@mail.ru" />
        </div>
        <div class="field">
          <label>Пароль</label>
          <input v-model="form.password" type="password" />
        </div>
        <p class="forgot"><a href="#">Забыли пароль? Восстановить</a></p>
        <div v-if="error" class="error-msg">{{ error }}</div>
        <button type="submit" class="btn-primary" :disabled="loading">
          {{ loading ? 'Загрузка...' : 'Войти' }}
        </button>
        <p class="link-row" v-if="mode === 'parent'">Нет аккаунта? <router-link to="/register">Зарегистрироваться</router-link></p>
      </form>

      <div class="about">
        <p class="about-title">О нас</p>
        <p class="about-text">Сервис геймификации домашних обязанностей и система финансового воспитания для детей. Задача → начисление валюты → накопление → реальная покупка.</p>
      </div>
    </div>
  </div>
</template>

<script>
import { useApi } from '../composables/useApi'
import { useAuth } from '../composables/useAuth'
export default {
  name: 'LoginPage',
  data() { return { mode: 'parent', loading: false, error: '', form: { email: '', username: '', password: '' } } },
  methods: {
    async login() {
      this.loading = true; this.error = ''
      const { login, loginChild, getMe } = useApi()
      const { saveAuth } = useAuth()
      try {
        let res
        if (this.mode === 'parent') {
          res = await login({ email: this.form.email, password: this.form.password })
          saveAuth({ token: res.data.token, role: 'parent', user: res.data.user })
          this.$router.push('/parent/dashboard')
        } else {
          res = await loginChild({ username: this.form.username, password: this.form.password })
          const user = res.data.child || res.data.user
          saveAuth({ token: res.data.token, role: 'child', user })
          const group = user.age_group || 'junior'
          this.$router.push(`/child/${group}/dashboard`)
        }
      } catch (e) {
        this.error = e.response?.data?.error?.message || 'Неверные данные'
      } finally { this.loading = false }
    }
  }
}
</script>

<style scoped>
.auth-page { min-height: 100vh; display: flex; align-items: flex-start; justify-content: center; background: #f5f7ff; padding: 40px 16px; }
.auth-card { background: #fff; border-radius: 16px; padding: 32px; width: 100%; max-width: 360px; box-shadow: 0 4px 24px rgba(0,0,0,0.08); }
.logo { font-size: 28px; font-weight: 800; color: #1a1a1a; text-align: center; margin-bottom: 24px; }
.accent { color: #4f7ef7; }
.role-label { text-align: center; font-size: 16px; color: #333; margin-bottom: 12px; font-weight: 500; }
.role-btns { display: flex; gap: 8px; justify-content: center; margin-bottom: 20px; }
.role-btn { padding: 8px 24px; border-radius: 8px; border: 2px solid #4f7ef7; color: #4f7ef7; font-weight: 600; font-size: 14px; background: transparent; cursor: pointer; transition: all 0.2s; }
.role-btn--active { background: #4f7ef7; color: #fff; }
.form-title { font-size: 22px; font-weight: 700; text-align: center; margin-bottom: 20px; }
.field { margin-bottom: 14px; }
.field label { display: block; font-size: 13px; color: #666; margin-bottom: 4px; }
.field input { width: 100%; padding: 10px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 15px; outline: none; transition: border-color 0.2s; }
.field input:focus { border-color: #4f7ef7; }
.forgot { font-size: 13px; color: #4f7ef7; text-align: right; margin-bottom: 12px; }
.forgot a { color: #4f7ef7; }
.error-msg { color: #e53e3e; font-size: 13px; margin-bottom: 10px; text-align: center; }
.btn-primary { width: 100%; padding: 12px; background: #4f7ef7; color: #fff; border: none; border-radius: 8px; font-size: 16px; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-primary:hover:not(:disabled) { background: #3a6be0; }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
.link-row { text-align: center; font-size: 13px; color: #666; margin-top: 12px; }
.link-row a { color: #4f7ef7; font-weight: 600; }
.about { margin-top: 24px; border-top: 1px solid #eee; padding-top: 16px; }
.about-title { font-size: 14px; font-weight: 700; margin-bottom: 6px; }
.about-text { font-size: 13px; color: #666; line-height: 1.6; }
</style>