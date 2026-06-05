<template>
  <div class="page">
    <div class="page__bg"></div>
    <header class="top-bar">
      <router-link to="/" class="logo">kid<span class="accent">TASK</span></router-link>
    </header>
    <div class="content">
      <div class="card">
        <div class="card__icon">👨‍👩‍👧</div>
        <h1 class="card__title">Регистрация</h1>
        <p class="card__sub">Создайте аккаунт родителя</p>

        <div class="field">
          <label>Ваше имя</label>
          <input v-model="form.name" type="text" placeholder="Иван Иванов" autocomplete="name" />
        </div>
        <div class="field">
          <label>Email</label>
          <input v-model="form.email" type="email" placeholder="parent@example.com" autocomplete="email" />
        </div>
        <div class="field">
          <label>Пароль</label>
          <div class="input-wrap">
            <input v-model="form.password" :type="showPass ? 'text' : 'password'" placeholder="Минимум 6 символов" autocomplete="new-password" />
            <button class="eye-btn" type="button" @click="showPass = !showPass">{{ showPass ? '🙈' : '👁' }}</button>
          </div>
        </div>
        <div class="field">
          <label>Повторите пароль</label>
          <div class="input-wrap">
            <input v-model="form.confirm" :type="showPass2 ? 'text' : 'password'" placeholder="Повторите пароль" autocomplete="new-password" @keyup.enter="register" />
            <button class="eye-btn" type="button" @click="showPass2 = !showPass2">{{ showPass2 ? '🙈' : '👁' }}</button>
          </div>
        </div>

        <div v-if="error" class="error-msg">{{ error }}</div>
        <div v-if="success" class="success-msg">Аккаунт создан! Перенаправляем...</div>

        <button class="btn-primary" @click="register" :disabled="loading">
          {{ loading ? 'Создаём аккаунт...' : 'Зарегистрироваться' }}
        </button>

        <p class="card__footer">Уже есть аккаунт? <router-link to="/login" class="link">Войти</router-link></p>
      </div>
    </div>
  </div>
</template>

<script>
import { useApi } from '../composables/useApi'
import { useAuth } from '../composables/useAuth'
export default {
  name: 'ParentRegister',
  data() {
    return { form: { name: '', email: '', password: '', confirm: '' }, loading: false, error: '', success: false, showPass: false, showPass2: false }
  },
  methods: {
    async register() {
      this.error = ''
      if (!this.form.name || !this.form.email || !this.form.password) { this.error = 'Заполните все поля'; return }
      if (this.form.password.length < 6) { this.error = 'Пароль должен быть не менее 6 символов'; return }
      if (this.form.password !== this.form.confirm) { this.error = 'Пароли не совпадают'; return }
      this.loading = true
      const { register } = useApi()
      const { saveAuth } = useAuth()
      try {
        const res = await register({ name: this.form.name, email: this.form.email, password: this.form.password })
        saveAuth({ token: res.data.token, role: 'parent', user: res.data.user })
        this.success = true
        setTimeout(() => this.$router.push('/parent/dashboard'), 800)
      } catch (e) {
        this.error = e.response?.data?.error?.message || 'Ошибка регистрации'
      } finally { this.loading = false }
    }
  }
}
</script>

<style scoped>
.page { min-height: 100vh; display: flex; flex-direction: column; align-items: center; background: #f0f4ff; position: relative; }
.page__bg { position: fixed; inset: 0; background: linear-gradient(135deg, #f0f4ff 0%, #fff 50%, #f5f0ff 100%); z-index: 0; }

.top-bar { width: 100%; padding: 16px 24px; display: flex; align-items: center; position: relative; z-index: 1; }
.logo { font-size: 22px; font-weight: 800; color: #1a1a1a; }
.accent { color: #4f7ef7; }

.content { flex: 1; display: flex; align-items: center; justify-content: center; width: 100%; padding: 16px; position: relative; z-index: 1; }

.card { background: #fff; border-radius: 24px; padding: 36px 32px; width: 100%; max-width: 420px; box-shadow: 0 8px 40px rgba(79,126,247,0.12); }
.card__icon { font-size: 40px; text-align: center; margin-bottom: 8px; }
.card__title { font-size: 26px; font-weight: 800; text-align: center; margin-bottom: 4px; color: #1a1a1a; }
.card__sub { text-align: center; color: #888; font-size: 14px; margin-bottom: 24px; }

.field { margin-bottom: 16px; }
.field label { display: block; font-size: 13px; font-weight: 600; color: #555; margin-bottom: 6px; }
.field input { width: 100%; padding: 13px 14px; border: 1.5px solid #e0e7ff; border-radius: 12px; font-size: 15px; outline: none; box-sizing: border-box; font-family: inherit; transition: border-color 0.15s; background: #fafbff; }
.field input:focus { border-color: #4f7ef7; background: #fff; }

.input-wrap { position: relative; }
.input-wrap input { padding-right: 44px; }
.eye-btn { position: absolute; right: 12px; top: 50%; transform: translateY(-50%); background: none; border: none; cursor: pointer; font-size: 16px; padding: 4px; }

.error-msg { color: #e53e3e; font-size: 13px; margin-bottom: 14px; text-align: center; background: #fff5f5; padding: 10px; border-radius: 8px; }
.success-msg { color: #155724; font-size: 13px; margin-bottom: 14px; text-align: center; background: #d4edda; padding: 10px; border-radius: 8px; }

.btn-primary { width: 100%; padding: 14px; background: #4f7ef7; color: #fff; border: none; border-radius: 12px; font-size: 16px; font-weight: 700; cursor: pointer; transition: background 0.15s; font-family: inherit; }
.btn-primary:hover:not(:disabled) { background: #3a6be0; }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }

.card__footer { text-align: center; font-size: 14px; color: #888; margin-top: 18px; }
.link { color: #4f7ef7; font-weight: 600; }

@media (max-width: 480px) {
  .card { padding: 24px 18px; border-radius: 20px; }
  .card__title { font-size: 22px; }
  .content { align-items: flex-start; padding-top: 12px; }
}
</style>
