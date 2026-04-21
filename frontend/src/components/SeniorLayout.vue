<template>
  <div class="layout">
    <header class="header">
      <div class="header__inner">
        <div class="logo">kid<span class="accent">TASK</span></div>
        <div class="header__right">
          <div class="balance">⭐ {{ balance }}</div>
          <div class="avatar" @click="menuOpen = !menuOpen">
            <span>{{ userInitial }}</span>
          </div>
          <div class="dropdown" v-if="menuOpen">
            <button class="dropdown__item" @click="go('/child/senior/wishlist')">Мой вишлист</button>
            <button class="dropdown__item" @click="go('/child/senior/tasks')">Мои задания</button>
            <button class="dropdown__item dropdown__item--sep" @click="logout">Выйти</button>
          </div>
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
import { useApi } from '../composables/useApi'
export default {
  name: 'SeniorLayout',
  data() { return { menuOpen: false, balance: 0 } },
  computed: {
    userInitial() { return useAuth().user.value?.name?.[0]?.toUpperCase() || 'R' }
  },
  async mounted() {
    document.addEventListener('click', this.handleOutside)
    try {
      const { getMe } = useApi()
      const r = await getMe()
      this.balance = r.data.user?.balance ?? 0
    } catch {}
  },
  beforeUnmount() { document.removeEventListener('click', this.handleOutside) },
  methods: {
    go(path) { this.menuOpen = false; this.$router.push(path) },
    logout() { useAuth().logout(); this.$router.push('/') },
    handleOutside(e) {
      if (!this.$el.querySelector('.header__right')?.contains(e.target)) this.menuOpen = false
    }
  }
}
</script>

<style scoped>
.layout { min-height: 100vh; background: #111827; display: flex; flex-direction: column; color: #e5e7eb; }
.header { background: #1f2937; border-bottom: 1px solid #374151; position: sticky; top: 0; z-index: 100; }
.header__inner { max-width: 700px; margin: 0 auto; padding: 12px 20px; display: flex; align-items: center; justify-content: space-between; }
.logo { font-size: 22px; font-weight: 800; color: #f3f4f6; }
.accent { color: #6366f1; }
.header__right { display: flex; align-items: center; gap: 12px; position: relative; }
.balance { background: #374151; border-radius: 20px; padding: 6px 14px; font-weight: 700; color: #a5b4fc; font-size: 15px; }
.avatar { width: 36px; height: 36px; border-radius: 50%; background: #6366f1; color: #fff; display: flex; align-items: center; justify-content: center; font-weight: 700; font-size: 16px; cursor: pointer; }
.dropdown { position: absolute; top: calc(100% + 8px); right: 0; background: #1f2937; border: 1px solid #374151; border-radius: 12px; box-shadow: 0 4px 24px rgba(0,0,0,0.4); min-width: 180px; overflow: hidden; z-index: 200; }
.dropdown__item { display: block; width: 100%; padding: 12px 16px; border: none; background: transparent; text-align: left; font-size: 15px; color: #e5e7eb; cursor: pointer; transition: background 0.15s; }
.dropdown__item:hover { background: #374151; }
.dropdown__item--sep { border-top: 1px solid #374151; color: #f87171; }
.main { flex: 1; max-width: 700px; width: 100%; margin: 0 auto; padding: 24px 20px; }
</style>