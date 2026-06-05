<template>
  <div class="layout">
    <aside class="sidebar">
      <div class="sidebar__logo">kid<span class="accent">TASK</span></div>
      <div class="sidebar__balance" v-if="user">
        <div class="sidebar__avatar">
          <img v-if="user.avatar_url" :src="user.avatar_url" class="sidebar__avatar-img" />
          <span v-else>{{ (user.name || '?')[0] }}</span>
        </div>
        <div>
          <div class="sidebar__name">{{ user.name }}</div>
          <div class="sidebar__coins">⭐ {{ user.balance || 0 }}</div>
        </div>
      </div>
      <nav class="sidebar__nav">
        <router-link to="/child/senior/dashboard" class="nav-item">
          <span class="nav-icon">🏠</span>
          <span class="nav-label">Главная</span>
        </router-link>
        <router-link to="/child/senior/tasks" class="nav-item">
          <span class="nav-icon">📋</span>
          <span class="nav-label">Задания</span>
        </router-link>
        <router-link to="/child/senior/wishlist" class="nav-item">
          <span class="nav-icon">🎁</span>
          <span class="nav-label">Вишлист</span>
        </router-link>
      </nav>
      <button class="sidebar__logout" @click="logout">Выйти</button>
    </aside>

    <div class="main-wrap">
      <header class="mobile-header">
        <div class="mobile-header__logo">kid<span class="accent">TASK</span></div>
        <div class="mobile-header__right" v-if="user">
          <span class="mobile-balance">⭐ {{ user.balance || 0 }}</span>
        </div>
      </header>

      <main class="main">
        <slot />
      </main>
    </div>

    <nav class="bottom-nav">
      <router-link to="/child/senior/dashboard" class="bottom-nav__item">
        <span class="bottom-nav__icon">🏠</span>
        <span class="bottom-nav__label">Главная</span>
      </router-link>
      <router-link to="/child/senior/tasks" class="bottom-nav__item">
        <span class="bottom-nav__icon">📋</span>
        <span class="bottom-nav__label">Задания</span>
      </router-link>
      <router-link to="/child/senior/wishlist" class="bottom-nav__item">
        <span class="bottom-nav__icon">🎁</span>
        <span class="bottom-nav__label">Вишлист</span>
      </router-link>
      <button class="bottom-nav__item bottom-nav__item--logout" @click="logout">
        <span class="bottom-nav__icon">👋</span>
        <span class="bottom-nav__label">Выйти</span>
      </button>
    </nav>
  </div>
</template>

<script>
import { useAuth } from '../composables/useAuth'
export default {
  name: 'SeniorLayout',
  computed: {
    user() { return useAuth().user.value }
  },
  methods: {
    logout() {
      useAuth().logout()
      this.$router.push('/')
    }
  }
}
</script>

<style scoped>
.layout { display: flex; min-height: 100vh; background: #111827; }

.sidebar { width: 220px; background: #1f2937; border-right: 1px solid #374151; display: flex; flex-direction: column; padding: 24px 0; flex-shrink: 0; position: sticky; top: 0; height: 100vh; }
.sidebar__logo { font-size: 20px; font-weight: 800; padding: 0 20px 20px; color: #f3f4f6; }
.accent { color: #6366f1; }
.sidebar__balance { display: flex; align-items: center; gap: 10px; padding: 14px 16px; background: #111827; margin: 0 10px 16px; border-radius: 14px; border: 1px solid #374151; }
.sidebar__avatar { width: 40px; height: 40px; border-radius: 50%; background: #6366f1; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 18px; font-weight: 700; flex-shrink: 0; overflow: hidden; }
.sidebar__avatar-img { width: 100%; height: 100%; object-fit: cover; }
.sidebar__name { font-size: 14px; font-weight: 700; color: #f3f4f6; }
.sidebar__coins { font-size: 16px; font-weight: 800; color: #a5b4fc; }
.sidebar__nav { flex: 1; display: flex; flex-direction: column; gap: 2px; padding: 0 10px; }
.nav-item { display: flex; align-items: center; gap: 10px; padding: 11px 12px; border-radius: 10px; color: #6b7280; font-size: 15px; font-weight: 600; text-decoration: none; transition: all 0.15s; }
.nav-item:hover { background: #374151; color: #d1d5db; }
.nav-item.router-link-active { background: #312e81; color: #a5b4fc; }
.nav-icon { font-size: 20px; }
.sidebar__logout { margin: 16px 10px 0; padding: 11px 12px; border: none; border-radius: 10px; background: rgba(239,68,68,0.1); color: #f87171; font-size: 14px; font-weight: 600; cursor: pointer; text-align: left; transition: background 0.15s; font-family: inherit; }
.sidebar__logout:hover { background: rgba(239,68,68,0.2); }

.mobile-header { display: none; }
.main-wrap { flex: 1; display: flex; flex-direction: column; min-width: 0; }
.main { flex: 1; padding: 28px 32px; max-width: 800px; width: 100%; }

.bottom-nav { display: none; }

@media (max-width: 768px) {
  .sidebar { display: none; }
  .mobile-header { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px; background: #1f2937; border-bottom: 1px solid #374151; position: sticky; top: 0; z-index: 50; }
  .mobile-header__logo { font-size: 20px; font-weight: 800; color: #f3f4f6; }
  .mobile-balance { font-size: 17px; font-weight: 800; color: #a5b4fc; background: #312e81; padding: 5px 12px; border-radius: 20px; }
  .main { padding: 14px; padding-bottom: 88px; max-width: 100%; }
  .bottom-nav { display: flex; position: fixed; bottom: 0; left: 0; right: 0; background: #1f2937; border-top: 1px solid #374151; padding: 8px 0 env(safe-area-inset-bottom, 8px); z-index: 100; }
  .bottom-nav__item { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 2px; padding: 6px 4px; text-decoration: none; color: #4b5563; transition: color 0.15s; background: none; border: none; cursor: pointer; font-family: inherit; }
  .bottom-nav__item.router-link-active { color: #a5b4fc; }
  .bottom-nav__item--logout { color: #6b7280; }
  .bottom-nav__icon { font-size: 22px; }
  .bottom-nav__label { font-size: 10px; font-weight: 600; color: inherit; }
}
</style>
