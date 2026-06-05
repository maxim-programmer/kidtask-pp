<template>
  <div class="page">
    <div class="page__bg"></div>
    <header class="top-bar">
      <div class="logo" @click="$router.push('/')">kid<span class="accent">TASK</span> <span class="badge">Admin</span></div>
    </header>
    <div class="content">
      <div class="card">
        <div class="card__icon">🔐</div>
        <h1 class="title">Вход для администратора</h1>
        <p class="card__sub">Введите секретный ключ доступа</p>
        <div class="field">
          <label>Секретный ключ</label>
          <div class="input-wrap">
            <input
              v-model="secret"
              :type="showPass ? 'text' : 'password'"
              placeholder="Введите admin secret"
              @keyup.enter="login"
              autofocus
            />
            <button class="eye-btn" type="button" @click="showPass = !showPass">{{ showPass ? '🙈' : '👁' }}</button>
          </div>
        </div>
        <div v-if="error" class="error-msg">{{ error }}</div>
        <button class="btn" @click="login" :disabled="loading">{{ loading ? 'Проверка...' : 'Войти' }}</button>
        <router-link to="/" class="back-link">← На главную</router-link>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AdminLogin',
  data() {
    return { secret: '', error: '', loading: false, showPass: false }
  },
  methods: {
    async login() {
      this.error = ''
      if (!this.secret) { this.error = 'Введите секретный ключ'; return }
      if (/[^\x00-\x7F]/.test(this.secret)) { this.error = 'Секретный ключ должен содержать только латинские символы'; return }
      this.loading = true
      try {
        const res = await fetch('/api/admin/stats', { headers: { 'X-Admin-Secret': this.secret } })
        if (res.status === 401 || res.status === 403 || !res.ok) { this.error = 'Неверный секретный ключ'; return }
        localStorage.setItem('kt_admin_secret', this.secret)
        this.$router.push('/admin')
      } catch { this.error = 'Ошибка соединения с сервером' }
      finally { this.loading = false }
    }
  }
}
</script>

<style scoped>
.page { min-height: 100vh; display: flex; flex-direction: column; align-items: center; position: relative; overflow: hidden; }
.page__bg { position: fixed; inset: 0; background: linear-gradient(135deg, #1e2235 0%, #0d1117 100%); z-index: 0; }

.top-bar { width: 100%; padding: 16px 24px; display: flex; align-items: center; position: relative; z-index: 1; }
.logo { font-size: 20px; font-weight: 800; color: #fff; cursor: pointer; user-select: none; transition: opacity 0.15s; }
.logo:hover { opacity: 0.75; }
.accent { color: #4f7ef7; }
.badge { font-size: 10px; font-weight: 700; background: #4f7ef7; color: #fff; border-radius: 6px; padding: 2px 8px; vertical-align: middle; margin-left: 6px; }

.content { flex: 1; display: flex; align-items: center; justify-content: center; width: 100%; padding: 16px; position: relative; z-index: 1; }

.card { background: #1e2235; border: 1px solid rgba(255,255,255,0.08); border-radius: 24px; padding: 36px 32px; width: 100%; max-width: 400px; box-shadow: 0 20px 60px rgba(0,0,0,0.4); }
.card__icon { font-size: 44px; text-align: center; margin-bottom: 8px; }
.title { font-size: 22px; font-weight: 800; color: #fff; text-align: center; margin-bottom: 6px; }
.card__sub { text-align: center; color: #6b7280; font-size: 13px; margin-bottom: 24px; }

.field { margin-bottom: 18px; }
.field label { display: block; font-size: 13px; font-weight: 600; color: #9ca3af; margin-bottom: 8px; }
.field input { width: 100%; padding: 13px 44px 13px 14px; background: #111827; border: 1.5px solid #374151; border-radius: 12px; color: #f3f4f6; font-size: 15px; outline: none; box-sizing: border-box; font-family: inherit; transition: border-color 0.15s; }
.field input:focus { border-color: #4f7ef7; }

.input-wrap { position: relative; }
.eye-btn { position: absolute; right: 12px; top: 50%; transform: translateY(-50%); background: none; border: none; cursor: pointer; font-size: 16px; padding: 4px; }

.error-msg { color: #f87171; font-size: 13px; margin-bottom: 14px; text-align: center; background: rgba(239,68,68,0.1); padding: 10px; border-radius: 8px; }

.btn { width: 100%; padding: 14px; background: #4f7ef7; color: #fff; border: none; border-radius: 12px; font-size: 16px; font-weight: 700; cursor: pointer; transition: background 0.15s; font-family: inherit; }
.btn:hover:not(:disabled) { background: #3a6be0; }
.btn:disabled { opacity: 0.6; cursor: not-allowed; }

.back-link { display: block; text-align: center; margin-top: 16px; font-size: 13px; color: #4b5563; text-decoration: none; transition: color 0.15s; }
.back-link:hover { color: #9ca3af; }

@media (max-width: 480px) {
  .card { padding: 28px 20px; border-radius: 20px; }
  .content { align-items: flex-start; padding-top: 20px; }
}
</style>
