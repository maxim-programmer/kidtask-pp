<template>
  <div class="page">
    <div class="page__bg"></div>
    <header class="top-bar">
      <router-link to="/" class="logo">kid<span class="accent">TASK</span></router-link>
    </header>
    <div class="content">
      <div class="card">
        <div class="card__icon">{{ isChild ? '🧒' : '👨‍👩‍👧' }}</div>
        <h1 class="card__title">Войти</h1>

        <div class="role-tabs">
          <button :class="['role-tab', { active: !isChild }]" @click="isChild = false">Родитель</button>
          <button :class="['role-tab', { active: isChild }]" @click="isChild = true">Ребёнок</button>
        </div>

        <div v-if="!isChild">
          <div class="field">
            <label>Email</label>
            <input v-model="email" type="email" placeholder="parent@example.com" @keyup.enter="login" autocomplete="email" />
          </div>
          <div class="field">
            <label>Пароль</label>
            <div class="input-wrap">
              <input v-model="password" :type="showPass ? 'text' : 'password'" placeholder="Пароль" @keyup.enter="login" autocomplete="current-password" />
              <button class="eye-btn" type="button" @click="showPass = !showPass">{{ showPass ? '🙈' : '👁' }}</button>
            </div>
          </div>
        </div>

        <div v-else>
          <div class="field">
            <label>Логин</label>
            <input v-model="username" type="text" placeholder="ivan123" @keyup.enter="login" autocomplete="username" />
          </div>
          <div class="field">
            <label>Пароль</label>
            <div class="input-wrap">
              <input v-model="password" :type="showPass ? 'text' : 'password'" placeholder="Пароль" @keyup.enter="login" autocomplete="current-password" />
              <button class="eye-btn" type="button" @click="showPass = !showPass">{{ showPass ? '🙈' : '👁' }}</button>
            </div>
          </div>
        </div>

        <div v-if="error" class="error-msg" :class="{ 'error-msg--blocked': isBlocked }">{{ error }}</div>

        <button class="btn-primary" @click="login" :disabled="loading">
          {{ loading ? 'Входим...' : 'Войти' }}
        </button>

        <p class="card__footer" v-if="!isChild">
          Нет аккаунта? <router-link to="/register" class="link">Зарегистрироваться</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { useAuth } from '../composables/useAuth'
import { useApi } from '../composables/useApi'
export default {
  name: 'LoginPage',
  data() {
    return { email: '', password: '', username: '', isChild: false, loading: false, error: '', showPass: false, isBlocked: false }
  },
  methods: {
    async login() {
      this.error = ''
      this.isBlocked = false
      if (this.isChild) {
        if (!this.username || !this.password) { this.error = 'Заполните все поля'; return }
      } else {
        if (!this.email || !this.password) { this.error = 'Заполните все поля'; return }
      }
      this.loading = true
      const { login, loginChild } = useApi()
      const { saveAuth } = useAuth()
      try {
        if (this.isChild) {
          const res = await loginChild({ username: this.username, password: this.password })
          const child = res.data.child
          saveAuth({ token: res.data.token, role: 'child', user: { ...child, user_id: child.child_id } })
          const group = child.age_group || 'junior'
          this.$router.push(`/child/${group}/dashboard`)
        } else {
          const res = await login({ email: this.email, password: this.password })
          saveAuth({ token: res.data.token, role: 'parent', user: res.data.user })
          this.$router.push('/parent/dashboard')
        }
      } catch (e) {
        const code = e.response?.data?.error?.code
        if (code === 'USER_BLOCKED') {
          this.isBlocked = true
          this.error = 'Вас заблокировал Администратор'
        } else {
          this.error = e.response?.data?.error?.message || 'Неверные данные'
        }
      } finally { this.loading = false }
    }
  }
}
</script>

<style scoped>
.page { min-height: 100vh; display: flex; flex-direction: column; align-items: center; background: #f0f4ff; position: relative; overflow: hidden; }
.page__bg { position: fixed; inset: 0; background: linear-gradient(135deg, #f0f4ff 0%, #fff 50%, #f5f0ff 100%); z-index: 0; }

.top-bar { width: 100%; padding: 16px 24px; display: flex; align-items: center; position: relative; z-index: 1; }
.logo { font-size: 22px; font-weight: 800; color: #1a1a1a; }
.accent { color: #4f7ef7; }

.content { flex: 1; display: flex; align-items: center; justify-content: center; width: 100%; padding: 16px; position: relative; z-index: 1; }

.card { background: #fff; border-radius: 24px; padding: 36px 32px; width: 100%; max-width: 420px; box-shadow: 0 8px 40px rgba(79,126,247,0.12); }
.card__icon { font-size: 40px; text-align: center; margin-bottom: 8px; }
.card__title { font-size: 26px; font-weight: 800; text-align: center; margin-bottom: 20px; color: #1a1a1a; }

.role-tabs { display: flex; background: #f0f4ff; border-radius: 12px; padding: 4px; margin-bottom: 24px; }
.role-tab { flex: 1; padding: 10px; border: none; background: transparent; color: #666; font-size: 14px; font-weight: 600; cursor: pointer; border-radius: 8px; transition: all 0.15s; font-family: inherit; }
.role-tab.active { background: #fff; color: #4f7ef7; box-shadow: 0 2px 8px rgba(0,0,0,0.08); }

.field { margin-bottom: 16px; }
.field label { display: block; font-size: 13px; font-weight: 600; color: #555; margin-bottom: 6px; }
.field input { width: 100%; padding: 13px 14px; border: 1.5px solid #e0e7ff; border-radius: 12px; font-size: 15px; outline: none; box-sizing: border-box; font-family: inherit; transition: border-color 0.15s; background: #fafbff; }
.field input:focus { border-color: #4f7ef7; background: #fff; }

.input-wrap { position: relative; }
.input-wrap input { padding-right: 44px; }
.eye-btn { position: absolute; right: 12px; top: 50%; transform: translateY(-50%); background: none; border: none; cursor: pointer; font-size: 16px; padding: 4px; }

.error-msg { color: #e53e3e; font-size: 13px; margin-bottom: 14px; text-align: center; background: #fff5f5; padding: 10px; border-radius: 8px; }
.error-msg--blocked { color: #c53030; font-size: 14px; font-weight: 700; background: #ffe0e0; border: 1.5px solid #fc8181; padding: 12px; }

.btn-primary { width: 100%; padding: 14px; background: #4f7ef7; color: #fff; border: none; border-radius: 12px; font-size: 16px; font-weight: 700; cursor: pointer; transition: background 0.15s; font-family: inherit; }
.btn-primary:hover:not(:disabled) { background: #3a6be0; }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }

.card__footer { text-align: center; font-size: 14px; color: #888; margin-top: 18px; }
.link { color: #4f7ef7; font-weight: 600; }

@media (max-width: 480px) {
  .card { padding: 28px 20px; border-radius: 20px; }
  .card__title { font-size: 22px; }
  .content { align-items: flex-start; padding-top: 20px; }
}
</style>