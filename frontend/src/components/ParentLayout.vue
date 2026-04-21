<template>
  <div class="layout">
    <header class="header">
      <div class="header__inner">
        <router-link to="/parent/dashboard" class="logo">kid<span class="accent">TASK</span></router-link>
        <div class="header__right" ref="menuRef">
          <div class="avatar" @click.stop="menuOpen = !menuOpen">
            <span>{{ userInitial }}</span>
          </div>
          <transition name="fade">
            <div class="dropdown" v-if="menuOpen">
              <button class="dropdown__item" @click="go('/parent/dashboard')">Главная</button>
              <button class="dropdown__item" @click="go('/parent/children')">Мои дети</button>
              <button class="dropdown__item" @click="go('/parent/settings')">Настройки</button>
              <button class="dropdown__item dropdown__item--danger" @click="logout">Выйти</button>
            </div>
          </transition>
        </div>
      </div>
    </header>
    <main class="main">
      <slot />
    </main>
  </div>
</template>

<script>
import { useAuth } from '../composables/useAuth'
export default {
  name: 'ParentLayout',
  data() { return { menuOpen: false } },
  computed: {
    userInitial() {
      const { user } = useAuth()
      return user.value?.name?.[0]?.toUpperCase() || 'P'
    }
  },
  mounted() {
    document.addEventListener('click', this.handleOutside)
  },
  beforeUnmount() {
    document.removeEventListener('click', this.handleOutside)
  },
  methods: {
    handleOutside(e) {
      if (this.$refs.menuRef && !this.$refs.menuRef.contains(e.target)) {
        this.menuOpen = false
      }
    },
    go(path) { this.menuOpen = false; this.$router.push(path) },
    logout() {
      const { logout } = useAuth()
      logout()
      this.$router.push('/')
    }
  }
}
</script>

<style scoped>
.layout { min-height: 100vh; background: #f5f7ff; display: flex; flex-direction: column; }
.header { background: #fff; border-bottom: 1px solid #e8e8e8; position: sticky; top: 0; z-index: 100; }
.header__inner { max-width: 900px; margin: 0 auto; padding: 12px 20px; display: flex; align-items: center; justify-content: space-between; }
.logo { font-size: 22px; font-weight: 800; color: #1a1a1a; text-decoration: none; }
.accent { color: #4f7ef7; }
.header__right { display: flex; align-items: center; gap: 12px; position: relative; }
.avatar { width: 38px; height: 38px; border-radius: 50%; background: #4f7ef7; color: #fff; display: flex; align-items: center; justify-content: center; font-weight: 700; font-size: 16px; cursor: pointer; user-select: none; transition: opacity 0.15s; }
.avatar:hover { opacity: 0.85; }
.dropdown { position: absolute; top: calc(100% + 10px); right: 0; background: #fff; border-radius: 14px; box-shadow: 0 8px 32px rgba(0,0,0,0.14); min-width: 200px; overflow: hidden; z-index: 200; border: 1px solid #eee; }
.dropdown__item { display: flex; align-items: center; gap: 8px; width: 100%; padding: 13px 18px; border: none; background: transparent; text-align: left; font-size: 15px; color: #333; cursor: pointer; transition: background 0.15s; font-family: inherit; }
.dropdown__item:hover { background: #f5f7ff; }
.dropdown__item--danger { color: #e53e3e; border-top: 1px solid #f0f0f0; margin-top: 4px; }
.dropdown__item--danger:hover { background: #fff5f5; }
.main { flex: 1; max-width: 900px; width: 100%; margin: 0 auto; padding: 24px 20px; }
.fade-enter-active, .fade-leave-active { transition: opacity 0.15s, transform 0.15s; }
.fade-enter-from, .fade-leave-to { opacity: 0; transform: translateY(-6px); }
</style>