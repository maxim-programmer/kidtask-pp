<template>
  <div class="page">
    <div class="card">
      <div class="logo">kid<span class="accent">TASK</span> <span class="badge">Admin</span></div>
      <h1 class="title">Вход для администратора</h1>
      <div class="field">
        <label>Секретный ключ</label>
        <input
          v-model="secret"
          type="password"
          placeholder="Введите admin secret"
          @keyup.enter="login"
          autofocus
        />
      </div>
      <div v-if="error" class="error-msg">{{ error }}</div>
      <button class="btn" @click="login">Войти</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AdminLogin',
  data() {
    return { secret: '', error: '' }
  },
  methods: {
    async login() {
      this.error = ''
      if (!this.secret) {
        this.error = 'Введите секретный ключ'
        return
      }
      try {
        const res = await fetch('/api/admin/stats', {
          headers: { 'X-Admin-Secret': this.secret }
        })
        if (res.status === 401) {
          this.error = 'Неверный секретный ключ'
          return
        }
        localStorage.setItem('kt_admin_secret', this.secret)
        this.$router.push('/admin')
      } catch {
        this.error = 'Ошибка соединения с сервером'
      }
    }
  }
}
</script>

<style scoped>
.page { min-height: 100vh; background: #f5f7ff; display: flex; align-items: center; justify-content: center; }
.card { background: #fff; border-radius: 20px; padding: 40px 36px; width: 90%; max-width: 380px; box-shadow: 0 4px 24px rgba(0,0,0,0.10); }
.logo { font-size: 22px; font-weight: 800; color: #1a1a1a; text-align: center; margin-bottom: 8px; }
.accent { color: #4f7ef7; }
.badge { font-size: 11px; font-weight: 700; background: #4f7ef7; color: #fff; border-radius: 6px; padding: 2px 8px; vertical-align: middle; margin-left: 6px; }
.title { font-size: 18px; font-weight: 700; text-align: center; color: #444; margin-bottom: 28px; }
.field { margin-bottom: 16px; }
.field label { display: block; font-size: 13px; color: #666; margin-bottom: 6px; font-weight: 500; }
.field input { width: 100%; padding: 12px 14px; border: 1.5px solid #ddd; border-radius: 10px; font-size: 15px; outline: none; box-sizing: border-box; transition: border-color 0.15s; }
.field input:focus { border-color: #4f7ef7; }
.error-msg { color: #e53e3e; font-size: 13px; margin-bottom: 12px; text-align: center; }
.btn { width: 100%; padding: 14px; background: #4f7ef7; color: #fff; border: none; border-radius: 10px; font-size: 16px; font-weight: 700; cursor: pointer; transition: background 0.15s; }
.btn:hover { background: #3a6be0; }
</style>