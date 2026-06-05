<template>
  <div class="layout">
    <aside class="sidebar">
      <div class="sidebar__logo">kid<span class="accent">TASK</span></div>
      <nav class="sidebar__nav">
        <router-link to="/parent/dashboard" class="nav-item">
          <span class="nav-icon">🏠</span>
          <span class="nav-label">Главная</span>
        </router-link>
        <router-link to="/parent/children" class="nav-item">
          <span class="nav-icon">👶</span>
          <span class="nav-label">Мои дети</span>
        </router-link>
        <router-link to="/parent/tasks" class="nav-item">
          <span class="nav-icon">📋</span>
          <span class="nav-label">Задания</span>
        </router-link>
        <router-link to="/parent/settings" class="nav-item">
          <span class="nav-icon">⚙️</span>
          <span class="nav-label">Настройки</span>
        </router-link>
      </nav>
      <button class="sidebar__logout" @click="logout">Выйти</button>
    </aside>

    <div class="main-wrap">
      <header class="mobile-header">
        <div class="mobile-header__logo">kid<span class="accent">TASK</span></div>
        <div class="mobile-header__right">
          <span class="mobile-header__role">Родитель</span>
        </div>
      </header>

      <main class="main">
        <slot />
      </main>
    </div>

    <nav class="bottom-nav">
      <router-link to="/parent/dashboard" class="bottom-nav__item">
        <span class="bottom-nav__icon">🏠</span>
        <span class="bottom-nav__label">Главная</span>
      </router-link>
      <router-link to="/parent/children" class="bottom-nav__item">
        <span class="bottom-nav__icon">👶</span>
        <span class="bottom-nav__label">Дети</span>
      </router-link>
      <router-link to="/parent/tasks" class="bottom-nav__item">
        <span class="bottom-nav__icon">📋</span>
        <span class="bottom-nav__label">Задания</span>
      </router-link>
      <router-link to="/parent/settings" class="bottom-nav__item">
        <span class="bottom-nav__icon">⚙️</span>
        <span class="bottom-nav__label">Профиль</span>
      </router-link>
    </nav>
  </div>
</template>

<script>
import { useAuth } from '../composables/useAuth'
export default {
  name: 'ParentLayout',
  methods: {
    logout() {
      useAuth().logout()
      this.$router.push('/')
    }
  }
}
</script>

<style scoped>
.layout { display: flex; min-height: 100vh; background: #f5f7ff; }

.sidebar { width: 220px; background: #fff; border-right: 1px solid #e8e8e8; display: flex; flex-direction: column; padding: 24px 0; flex-shrink: 0; position: sticky; top: 0; height: 100vh; }
.sidebar__logo { font-size: 20px; font-weight: 800; padding: 0 20px 24px; color: #1a1a1a; }
.accent { color: #4f7ef7; }
.sidebar__nav { flex: 1; display: flex; flex-direction: column; gap: 2px; padding: 0 10px; }
.nav-item { display: flex; align-items: center; gap: 10px; padding: 11px 12px; border-radius: 10px; color: #666; font-size: 14px; font-weight: 500; text-decoration: none; transition: all 0.15s; }
.nav-item:hover { background: #f5f7ff; color: #333; }
.nav-item.router-link-active { background: #eef2ff; color: #4f7ef7; font-weight: 600; }
.nav-icon { font-size: 18px; }
.sidebar__logout { margin: 16px 10px 0; padding: 11px 12px; border: none; border-radius: 10px; background: #fff5f5; color: #e53e3e; font-size: 14px; font-weight: 600; cursor: pointer; text-align: left; transition: background 0.15s; font-family: inherit; }
.sidebar__logout:hover { background: #ffe0e0; }

.mobile-header { display: none; }
.main-wrap { flex: 1; display: flex; flex-direction: column; min-width: 0; }
.main { flex: 1; padding: 28px 32px; overflow-y: auto; max-width: 900px; width: 100%; }

.bottom-nav { display: none; }

@media (max-width: 768px) {
  .sidebar { display: none; }
  .mobile-header { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px; background: #fff; border-bottom: 1px solid #e8e8e8; position: sticky; top: 0; z-index: 50; }
  .mobile-header__logo { font-size: 20px; font-weight: 800; color: #1a1a1a; }
  .mobile-header__role { font-size: 12px; font-weight: 600; background: #eef2ff; color: #4f7ef7; padding: 4px 10px; border-radius: 12px; }
  .main { padding: 16px; padding-bottom: 80px; max-width: 100%; }
  .bottom-nav { display: flex; position: fixed; bottom: 0; left: 0; right: 0; background: #fff; border-top: 1px solid #e8e8e8; padding: 8px 0 env(safe-area-inset-bottom, 8px); z-index: 100; }
  .bottom-nav__item { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 2px; padding: 6px 4px; text-decoration: none; color: #aaa; transition: color 0.15s; }
  .bottom-nav__item.router-link-active { color: #4f7ef7; }
  .bottom-nav__icon { font-size: 22px; }
  .bottom-nav__label { font-size: 10px; font-weight: 600; }
}
</style>
