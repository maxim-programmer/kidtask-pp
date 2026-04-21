<template>
  <ParentLayout>
    <h1 class="page-title">Настройки</h1>
    <div class="card">
      <div class="profile-row">
        <div class="avatar">{{ user?.name?.[0] }}</div>
        <div>
          <div class="name">{{ user?.name }}</div>
          <div class="email">{{ user?.email }}</div>
        </div>
      </div>
    </div>
    <button class="logout-btn" @click="logout">Выйти</button>
  </ParentLayout>
</template>

<script>
import ParentLayout from '../../components/ParentLayout.vue'
import { useAuth } from '../../composables/useAuth'
export default {
  name: 'ParentSettings',
  components: { ParentLayout },
  computed: { user() { return useAuth().user.value } },
  methods: {
    logout() { useAuth().logout(); this.$router.push('/') }
  }
}
</script>

<style scoped>
.page-title { font-size: 24px; font-weight: 700; margin-bottom: 20px; }
.card { background: #fff; border-radius: 16px; padding: 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); margin-bottom: 16px; }
.profile-row { display: flex; align-items: center; gap: 16px; }
.avatar { width: 52px; height: 52px; border-radius: 50%; background: #4f7ef7; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 24px; font-weight: 700; }
.name { font-size: 18px; font-weight: 700; }
.email { font-size: 14px; color: #888; margin-top: 2px; }
.logout-btn { width: 100%; padding: 14px; background: #fff5f5; border: 2px solid #e53e3e; border-radius: 12px; color: #e53e3e; font-size: 16px; font-weight: 600; cursor: pointer; }
.logout-btn:hover { background: #e53e3e; color: #fff; }
</style>