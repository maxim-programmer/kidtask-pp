<template>
  <ParentLayout>
    <h1 class="page-title">Мои дети</h1>

    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else>
      <div v-for="child in children" :key="child.child_id" class="child-card">
        <div class="child-card__avatar">{{ child.name[0] }}</div>
        <div class="child-card__info">
          <div class="child-card__name">{{ child.name }}</div>
          <div class="child-card__meta">⭐ {{ child.balance }} · {{ child.age_group === 'junior' ? '7-10 лет' : '11-14 лет' }}</div>
        </div>
        <button class="icon-btn" @click="removeChild(child)" title="Удалить">⋮</button>
      </div>

      <button class="add-btn" @click="showModal = true">+ Добавить ребёнка</button>
    </div>

    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <h3>Новый ребёнок</h3>
        <div class="avatar-preview">
          <div class="avatar-circle">👤</div>
          <span class="avatar-hint">Изменить изображение</span>
        </div>
        <div class="field">
          <label>Имя</label>
          <input v-model="form.name" type="text" placeholder="Иван" />
        </div>
        <div class="field">
          <label>Логин</label>
          <input v-model="form.username" type="text" placeholder="ivan123" />
        </div>
        <div class="field">
          <label>Пароль</label>
          <input v-model="form.password" type="password" />
        </div>
        <div class="field">
          <label>День рождения</label>
          <input v-model="form.birthday" type="date" />
        </div>
        <div v-if="error" class="error-msg">{{ error }}</div>
        <button class="btn-primary" @click="addChild" :disabled="saving">
          {{ saving ? 'Сохранение...' : 'Добавить' }}
        </button>
      </div>
    </div>
  </ParentLayout>
</template>

<script>
import ParentLayout from '../../components/ParentLayout.vue'
import { useApi } from '../../composables/useApi'
export default {
  name: 'ParentChildren',
  components: { ParentLayout },
  data() {
    return {
      loading: true, children: [], showModal: false, saving: false, error: '',
      form: { name: '', username: '', password: '', birthday: '' }
    }
  },
  async mounted() {
    await this.load()
  },
  methods: {
    async load() {
      const { getChildren } = useApi()
      this.loading = true
      try {
        const res = await getChildren()
        this.children = res.data.children || []
      } finally { this.loading = false }
    },
    async addChild() {
      this.saving = true; this.error = ''
      const { createChild } = useApi()
      try {
        const payload = {
          name: this.form.name,
          username: this.form.username,
          password: this.form.password,
        }
        if (this.form.birthday) {
          payload.birthday = new Date(this.form.birthday + 'T00:00:00Z').toISOString()
        }
        await createChild(payload)
        this.showModal = false
        this.form = { name: '', username: '', password: '', birthday: '' }
        await this.load()
      } catch (e) {
        this.error = e.response?.data?.error?.message || 'Ошибка'
      } finally { this.saving = false }
    },
    async removeChild(child) {
      if (!confirm(`Удалить ${child.name}?`)) return
      const { deleteChild } = useApi()
      await deleteChild(child.child_id)
      await this.load()
    }
  }
}
</script>

<style scoped>
.page-title { font-size: 24px; font-weight: 700; margin-bottom: 20px; }
.loading { text-align: center; padding: 60px; color: #888; }
.child-card { background: #fff; border-radius: 16px; padding: 16px; display: flex; align-items: center; gap: 14px; margin-bottom: 12px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); }
.child-card__avatar { width: 44px; height: 44px; border-radius: 50%; background: #4f7ef7; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 20px; font-weight: 700; }
.child-card__info { flex: 1; }
.child-card__name { font-size: 16px; font-weight: 700; }
.child-card__meta { font-size: 13px; color: #888; margin-top: 2px; }
.icon-btn { width: 32px; height: 32px; border-radius: 8px; border: none; background: #f5f5f5; font-size: 18px; cursor: pointer; color: #666; }
.add-btn { width: 100%; padding: 14px; background: #f0f4ff; border: 2px dashed #4f7ef7; border-radius: 16px; color: #4f7ef7; font-size: 16px; font-weight: 600; cursor: pointer; transition: background 0.2s; margin-top: 8px; }
.add-btn:hover { background: #e8eeff; }
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 300; }
.modal { background: #fff; border-radius: 16px; padding: 24px; width: 90%; max-width: 360px; }
.modal h3 { font-size: 20px; font-weight: 700; text-align: center; margin-bottom: 16px; }
.avatar-preview { display: flex; flex-direction: column; align-items: center; gap: 6px; margin-bottom: 16px; }
.avatar-circle { width: 64px; height: 64px; border-radius: 50%; background: #e8e8e8; display: flex; align-items: center; justify-content: center; font-size: 28px; }
.avatar-hint { font-size: 12px; color: #4f7ef7; cursor: pointer; }
.field { margin-bottom: 12px; }
.field label { display: block; font-size: 13px; color: #666; margin-bottom: 4px; }
.field input { width: 100%; padding: 10px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 15px; outline: none; }
.field input:focus { border-color: #4f7ef7; }
.error-msg { color: #e53e3e; font-size: 13px; margin-bottom: 10px; }
.btn-primary { width: 100%; padding: 12px; background: #4f7ef7; color: #fff; border: none; border-radius: 8px; font-size: 16px; font-weight: 600; cursor: pointer; margin-top: 8px; }
.btn-primary:hover:not(:disabled) { background: #3a6be0; }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
</style>