<template>
  <SeniorLayout>
    <div class="page-header">
      <div class="page-tabs">
        <button :class="['ptab', { 'ptab--active': view === 'wishlist' }]" @click="view = 'wishlist'">
          ❤ Мой вишлист
        </button>
        <button :class="['ptab', { 'ptab--active': view === 'new' }]" @click="view = 'new'">
          + Новая цель
        </button>
      </div>
    </div>

    <div v-if="loading" class="loading">Загрузка...</div>

    <div v-else-if="view === 'wishlist'">
      <div v-if="wishes.length === 0" class="empty">Добавь свою первую цель</div>
      <div v-for="w in wishes" :key="w.wish_id" class="wish-card">
        <div class="wish-card__top">
          <div class="wish-info">
            <div class="wish-title">{{ w.title }}</div>
            <div class="wish-desc" v-if="w.description">{{ w.description }}</div>
          </div>
          <div class="wish-price-col">
            <span class="wish-price">⭐ {{ w.price ?? '?' }}</span>
            <button class="edit-btn" @click="openEdit(w)" v-if="w.status === 'awaiting_price' || w.status === 'available'">✏</button>
          </div>
        </div>

        <div class="steps">
          <span :class="['step', 'step--done']">⭐ Создано</span>
          <span class="step-arrow">→</span>
          <span :class="['step', { 'step--done': w.status === 'purchased' || w.status === 'delivered' }]">🛒 Куплено</span>
          <span class="step-arrow">→</span>
          <span :class="['step', { 'step--done': w.status === 'delivered' }]">🎁 Доставлено</span>
        </div>

        <div v-if="w.price" class="progress-bar">
          <div class="progress-fill" :style="{ width: Math.min(100, (balance / w.price) * 100) + '%' }"></div>
        </div>

        <div class="wish-actions">
          <button
            v-if="w.status === 'available'"
            class="buy-btn"
            @click="purchase(w)"
            :disabled="balance < w.price || buying === w.wish_id">
            {{ balance >= (w.price || Infinity) ? '🛒 Купить' : `Нужно ещё ⭐${w.price - balance}` }}
          </button>
          <span v-else-if="w.status === 'purchased'" class="badge badge--bought">🛒 Куплено</span>
          <span v-else-if="w.status === 'delivered'" class="badge badge--done">🎁 Доставлено</span>
          <span v-else-if="w.status === 'awaiting_price'" class="badge badge--wait">Ожидает оценки</span>
          <button class="del-btn" @click="remove(w)" v-if="w.status !== 'purchased' && w.status !== 'delivered'">🗑</button>
        </div>
      </div>
    </div>

    <div v-else-if="view === 'new'">
      <div class="new-form">
        <h2 class="new-title">Новая цель</h2>
        <div class="field">
          <label>Название</label>
          <input v-model="form.title" type="text" placeholder="iPhone 16" />
        </div>
        <div class="field">
          <label>Описание</label>
          <textarea v-model="form.description" rows="3" placeholder="Зачем мне это нужно..." class="form-textarea"></textarea>
        </div>
        <div v-if="error" class="error-msg">{{ error }}</div>
        <button class="btn-primary" @click="addWish" :disabled="saving">
          {{ saving ? 'Добавление...' : 'Добавить' }}
        </button>
      </div>
    </div>

    <div v-if="editModal" class="modal-overlay" @click.self="editModal = false">
      <div class="modal">
        <h3>Редактировать</h3>
        <div class="field">
          <label>Название</label>
          <input v-model="editForm.title" type="text" />
        </div>
        <div class="field">
          <label>Описание</label>
          <textarea v-model="editForm.description" rows="2" class="form-textarea"></textarea>
        </div>
        <div class="modal-actions">
          <button class="btn-cancel" @click="editModal = false">Отмена</button>
          <button class="btn-primary-sm" @click="saveEdit">Сохранить</button>
        </div>
      </div>
    </div>
  </SeniorLayout>
</template>

<script>
import SeniorLayout from '../../../components/SeniorLayout.vue'
import { useApi } from '../../../composables/useApi'
import { useAuth } from '../../../composables/useAuth'
export default {
  name: 'ChildWishlistSenior',
  components: { SeniorLayout },
  data() {
    return {
      loading: true, wishes: [], balance: 0, buying: null,
      view: 'wishlist', saving: false, error: '',
      form: { title: '', description: '' },
      editModal: false, editTarget: null, editForm: { title: '', description: '' }
    }
  },
  computed: {
    childId() {
      const u = useAuth().user.value
      return u?.child_id || u?.user_id || u?.id
    }
  },
  async mounted() { await this.load() },
  methods: {
    async load() {
      const { getWishes, getMe } = useApi()
      this.loading = true
      try {
        const [wr, me] = await Promise.all([getWishes(this.childId), getMe()])
        this.wishes = wr.data.wishes || []
        this.balance = me.data.user?.balance || 0
      } finally { this.loading = false }
    },
    async addWish() {
      this.saving = true; this.error = ''
      const { createWish } = useApi()
      try {
        await createWish(this.childId, { title: this.form.title, description: this.form.description || undefined })
        this.form = { title: '', description: '' }
        this.view = 'wishlist'
        await this.load()
      } catch (e) { this.error = e.response?.data?.error?.message || 'Ошибка' }
      finally { this.saving = false }
    },
    async purchase(wish) {
      this.buying = wish.wish_id
      const { purchaseWish } = useApi()
      try { await purchaseWish(this.childId, wish.wish_id); await this.load() }
      catch (e) { alert(e.response?.data?.error?.message || 'Ошибка') }
      finally { this.buying = null }
    },
    async remove(wish) {
      if (!confirm('Удалить цель?')) return
      const { deleteWish } = useApi()
      await deleteWish(this.childId, wish.wish_id)
      await this.load()
    },
    openEdit(w) { this.editTarget = w; this.editForm = { title: w.title, description: w.description || '' }; this.editModal = true },
    async saveEdit() {
      const { updateWish } = useApi()
      await updateWish(this.childId, this.editTarget.wish_id, { title: this.editForm.title, description: this.editForm.description || undefined })
      this.editModal = false
      await this.load()
    }
  }
}
</script>

<style scoped>
.page-header { margin-bottom: 20px; }
.page-tabs { display: flex; gap: 4px; background: #1f2937; border-radius: 12px; padding: 4px; }
.ptab { flex: 1; padding: 10px; border: none; background: transparent; color: #9ca3af; font-size: 14px; font-weight: 600; cursor: pointer; border-radius: 8px; transition: all 0.15s; }
.ptab--active { background: #374151; color: #e5e7eb; }
.loading { text-align: center; padding: 60px; color: #6366f1; }
.empty { text-align: center; color: #6b7280; font-size: 15px; padding: 40px; }
.wish-card { background: #1f2937; border-radius: 14px; padding: 16px; margin-bottom: 12px; border: 1px solid #374151; }
.wish-card__top { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 12px; }
.wish-info { flex: 1; }
.wish-title { font-size: 16px; font-weight: 700; color: #f3f4f6; margin-bottom: 4px; }
.wish-desc { font-size: 13px; color: #6b7280; }
.wish-price-col { display: flex; flex-direction: column; align-items: flex-end; gap: 6px; }
.wish-price { font-size: 20px; font-weight: 800; color: #a5b4fc; }
.edit-btn { background: #374151; border: none; border-radius: 6px; padding: 4px 8px; color: #9ca3af; font-size: 13px; cursor: pointer; }
.steps { display: flex; align-items: center; gap: 6px; margin-bottom: 10px; flex-wrap: wrap; }
.step { font-size: 12px; color: #4b5563; font-weight: 600; }
.step--done { color: #a5b4fc; }
.step-arrow { color: #374151; }
.progress-bar { height: 6px; background: #374151; border-radius: 3px; margin-bottom: 12px; overflow: hidden; }
.progress-fill { height: 100%; background: #6366f1; border-radius: 3px; transition: width 0.4s; }
.wish-actions { display: flex; gap: 8px; align-items: center; }
.buy-btn { flex: 1; padding: 10px 16px; background: #22c55e; color: #fff; border: none; border-radius: 10px; font-size: 14px; font-weight: 600; cursor: pointer; }
.buy-btn:hover:not(:disabled) { background: #16a34a; }
.buy-btn:disabled { background: #374151; color: #6b7280; cursor: not-allowed; }
.badge { padding: 6px 12px; border-radius: 8px; font-size: 13px; font-weight: 600; }
.badge--bought { background: #78350f30; color: #fbbf24; }
.badge--done { background: #14532d30; color: #4ade80; }
.badge--wait { background: #1e1b4b; color: #818cf8; }
.del-btn { background: transparent; border: 1px solid #374151; border-radius: 8px; padding: 6px 10px; color: #6b7280; font-size: 14px; cursor: pointer; }
.del-btn:hover { background: #ef444420; color: #ef4444; border-color: #ef4444; }
.new-form { background: #1f2937; border-radius: 14px; padding: 24px; border: 1px solid #374151; }
.new-title { font-size: 20px; font-weight: 700; color: #e5e7eb; margin-bottom: 20px; }
.field { margin-bottom: 14px; }
.field label { display: block; font-size: 13px; color: #9ca3af; margin-bottom: 6px; font-weight: 500; }
.field input, .form-textarea { width: 100%; padding: 10px 12px; background: #111827; border: 1px solid #374151; border-radius: 10px; font-size: 15px; color: #e5e7eb; outline: none; }
.field input:focus, .form-textarea:focus { border-color: #6366f1; }
.form-textarea { resize: vertical; }
.error-msg { color: #f87171; font-size: 13px; margin-bottom: 10px; }
.btn-primary { width: 100%; padding: 12px; background: #6366f1; color: #fff; border: none; border-radius: 10px; font-size: 16px; font-weight: 600; cursor: pointer; }
.btn-primary:hover:not(:disabled) { background: #4f46e5; }
.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.6); display: flex; align-items: center; justify-content: center; z-index: 300; }
.modal { background: #1f2937; border: 1px solid #374151; border-radius: 14px; padding: 24px; width: 90%; max-width: 360px; }
.modal h3 { font-size: 18px; font-weight: 700; color: #e5e7eb; margin-bottom: 16px; }
.modal-actions { display: flex; gap: 8px; margin-top: 16px; justify-content: flex-end; }
.btn-cancel { padding: 10px 18px; background: transparent; border: 1px solid #374151; border-radius: 8px; color: #9ca3af; font-size: 14px; cursor: pointer; }
.btn-primary-sm { padding: 10px 18px; background: #6366f1; color: #fff; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; }
</style>