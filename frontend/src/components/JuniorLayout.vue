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
        <router-link to="/child/junior/dashboard" class="nav-item">
          <span class="nav-icon">🏠</span>
          <span class="nav-label">Главная</span>
        </router-link>
        <router-link to="/child/junior/tasks" class="nav-item">
          <span class="nav-icon">📋</span>
          <span class="nav-label">Задания</span>
        </router-link>
        <router-link to="/child/junior/wishlist" class="nav-item">
          <span class="nav-icon">🎁</span>
          <span class="nav-label">Вишлист</span>
        </router-link>
      </nav>
      <button class="sidebar__logout" @click="logout">Выйти</button>
    </aside>

    <div class="main-wrap">
      <header class="mobile-header">
        <div class="mobile-header__left">
          <div class="mobile-header__logo">kid<span class="accent">TASK</span></div>
        </div>
        <div class="mobile-header__right" v-if="user">
          <span class="mobile-balance">⭐ {{ user.balance || 0 }}</span>
        </div>
      </header>

      <main class="main">
        <slot />
      </main>
    </div>

    <nav class="bottom-nav">
      <router-link to="/child/junior/dashboard" class="bottom-nav__item">
        <span class="bottom-nav__icon">🏠</span>
        <span class="bottom-nav__label">Главная</span>
      </router-link>
      <router-link to="/child/junior/tasks" class="bottom-nav__item">
        <span class="bottom-nav__icon">📋</span>
        <span class="bottom-nav__label">Задания</span>
      </router-link>
      <router-link to="/child/junior/wishlist" class="bottom-nav__item">
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
  name: 'JuniorLayout',
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
.layout { display: flex; min-height: 100vh; background: #fff7f0; }

.sidebar { width: 220px; background: #fff; border-right: 1px solid #fed7aa; display: flex; flex-direction: column; padding: 24px 0; flex-shrink: 0; position: sticky; top: 0; height: 100vh; }
.sidebar__logo { font-size: 20px; font-weight: 800; padding: 0 20px 20px; color: #1a1a1a; }
.accent { color: #ea580c; }
.sidebar__balance { display: flex; align-items: center; gap: 10px; padding: 14px 16px; background: #fff7f0; margin: 0 10px 16px; border-radius: 14px; border: 1px solid #fed7aa; }
.sidebar__avatar { width: 40px; height: 40px; border-radius: 50%; background: #ea580c; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 18px; font-weight: 700; flex-shrink: 0; overflow: hidden; }
.sidebar__avatar-img { width: 100%; height: 100%; object-fit: cover; }
.sidebar__name { font-size: 14px; font-weight: 700; color: #1a1a1a; }
.sidebar__coins { font-size: 16px; font-weight: 800; color: #ea580c; }
.sidebar__nav { flex: 1; display: flex; flex-direction: column; gap: 2px; padding: 0 10px; }
.nav-item { display: flex; align-items: center; gap: 10px; padding: 11px 12px; border-radius: 10px; color: #888; font-size: 15px; font-weight: 600; text-decoration: none; transition: all 0.15s; }
.nav-item:hover { background: #fff7f0; color: #ea580c; }
.nav-item.router-link-active { background: #fff0e6; color: #ea580c; }
.nav-icon { font-size: 20px; }
.sidebar__logout { margin: 16px 10px 0; padding: 11px 12px; border: none; border-radius: 10px; background: #fff5f5; color: #e53e3e; font-size: 14px; font-weight: 600; cursor: pointer; text-align: left; transition: background 0.15s; font-family: inherit; }
.sidebar__logout:hover { background: #ffe0e0; }

.mobile-header { display: none; }
.main-wrap { flex: 1; display: flex; flex-direction: column; min-width: 0; }
.main { flex: 1; padding: 28px 32px; max-width: 800px; width: 100%; }

.bottom-nav { display: none; }

@media (max-width: 768px) {
  .sidebar { display: none; }
  .mobile-header { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px; background: #fff; border-bottom: 2px solid #fed7aa; position: sticky; top: 0; z-index: 50; }
  .mobile-header__logo { font-size: 20px; font-weight: 800; color: #1a1a1a; }
  .mobile-balance { font-size: 18px; font-weight: 800; color: #ea580c; background: #fff7f0; padding: 5px 12px; border-radius: 20px; border: 1.5px solid #fed7aa; }
  .main { padding: 14px; padding-bottom: 88px; max-width: 100%; }
  .bottom-nav { display: flex; position: fixed; bottom: 0; left: 0; right: 0; background: #fff; border-top: 2px solid #fed7aa; padding: 8px 0 env(safe-area-inset-bottom, 8px); z-index: 100; }
  .bottom-nav__item { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 2px; padding: 6px 4px; text-decoration: none; color: #fed7aa; transition: color 0.15s; background: none; border: none; cursor: pointer; font-family: inherit; }
  .bottom-nav__item.router-link-active { color: #ea580c; }
  .bottom-nav__item--logout { color: #fca5a5; }
  .bottom-nav__icon { font-size: 22px; }
  .bottom-nav__label { font-size: 10px; font-weight: 700; }
}
</style>
