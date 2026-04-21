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
            <button class="dropdown__item" @click="go('/child/junior/wishlist')">❤ Мой вишлист</button>
            <button class="dropdown__item" @click="go('/child/junior/tasks')">📋 Мои задания</button>
            <button class="dropdown__item" @click="logout">🚪 Выйти</button>
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
  name: 'JuniorLayout',
  data() { return { menuOpen: false, balance: 0 } },
  computed: {
    userInitial() { return useAuth().user.value?.name?.[0]?.toUpperCase() || 'R' }
  },
  async mounted() {
    document.addEventListener('click', this.handleOutside)
    await this.loadBalance()
  },
  beforeUnmount() { document.removeEventListener('click', this.handleOutside) },
  methods: {
    async loadBalance() {
      const { getMe } = useApi()
      try { const r = await getMe(); this.balance = r.data.user?.balance ?? 0 } catch {}
    },
    go(path) { this.menuOpen = false; this.$router.push(path) },
    logout() { useAuth().logout(); this.$router.push('/') },
    handleOutside(e) {
      if (!this.$el.querySelector('.header__right')?.contains(e.target)) this.menuOpen = false
    }
  }
}
</script>

<style scoped>
.layout { min-height: 100vh; background: #fff5eb; display: flex; flex-direction: column; }
.header { background: #ffedd5; border-bottom: 2px solid #fed7aa; position: sticky; top: 0; z-index: 100; }
.header__inner { max-width: 600px; margin: 0 auto; padding: 12px 20px; display: flex; align-items: center; justify-content: space-between; }
.logo { font-size: 22px; font-weight: 800; color: #ea580c; }
.accent { color: #c2410c; }
.header__right { display: flex; align-items: center; gap: 12px; position: relative; }
.balance { background: #fff; border-radius: 20px; padding: 6px 14px; font-weight: 700; color: #ea580c; font-size: 16px; box-shadow: 0 2px 8px rgba(234,88,12,0.15); }
.avatar { width: 36px; height: 36px; border-radius: 50%; background: #ea580c; color: #fff; display: flex; align-items: center; justify-content: center; font-weight: 700; font-size: 16px; cursor: pointer; }
.dropdown { position: absolute; top: calc(100% + 8px); right: 0; background: #fff; border-radius: 16px; box-shadow: 0 4px 24px rgba(234,88,12,0.18); min-width: 180px; overflow: hidden; z-index: 200; }
.dropdown__item { display: block; width: 100%; padding: 14px 16px; border: none; background: transparent; text-align: left; font-size: 16px; color: #333; cursor: pointer; font-weight: 600; }
.dropdown__item:hover { background: #fff5eb; }
.main { flex: 1; max-width: 600px; width: 100%; margin: 0 auto; padding: 20px 16px; }
</style>