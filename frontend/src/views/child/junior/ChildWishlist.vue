<template>
  <JuniorLayout>
    <div class="page-header">
      <h1 class="page-title">Мой вишлист</h1>
      <button class="add-icon-btn" @click="showModal = true">＋</button>
    </div>

    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else>
      <div v-if="wishes.length === 0" class="empty">Добавь свою первую цель! 🎯</div>
      <div v-for="w in wishes" :key="w.wish_id" class="wish-card">
        <div class="wish-card__top">
          <div class="wish-card__emoji">🎁</div>
          <div class="wish-card__info">
            <div class="wish-card__title">{{ w.title }}</div>
            <div class="wish-card__desc" v-if="w.description">{{ w.description }}</div>
          </div>
        </div>

        <div class="wish-price-row">
          <span class="wish-price">⭐ {{ w.price ?? '?' }}</span>
          <div class="wish-status-badges">
            <span class="step" :class="{ 'step--done': true }">⭐ Создано</span>
            <span class="arrow">→</span>
            <span class="step" :class="{ 'step--done': w.status === 'purchased' || w.status === 'delivered' }">🛒 Куплено</span>
            <span class="arrow">→</span>
            <span class="step" :class="{ 'step--done': w.status === 'delivered' }">🎁 Доставлено</span>
          </div>
        </div>

        <div v-if="w.price" class="progress-bar">
          <div class="progress-fill" :style="{ width: Math.min(100, ((balance) / w.price) * 100) + '%' }"></div>
        </div>

        <div class="wish-actions">
          <button
            v-if="w.status === 'available'"
            class="buy-btn"
            @click="purchase(w)"
            :disabled="balance < w.price || buying === w.wish_id">
            {{ buying === w.wish_id ? '...' : balance >= w.price ? '🛒 Купить!' : `Нужно ещё ⭐ ${w.price - balance}` }}
          </button>
          <span v-else-if="w.status === 'purchased'" class="badge badge--purchased">🛒 Куплено</span>
          <span v-else-if="w.status === 'delivered'" class="badge badge--delivered">🎁 Доставлено!</span>
          <span v-else-if="w.status === 'awaiting_price'" class="badge badge--wait">Ждём цену ⏳</span>
        </div>
      </div>

      <button class="add-btn" @click="showModal = true">+ Добавить цель</button>
    </div>

    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <h3>Новая цель 🎯</h3>
        <div class="field">
          <label>Название</label>
          <input v-model="form.title" type="text" placeholder="Nintendo Switch" />
        </div>
        <div class="field">
          <label>Описание (необязательно)</label>
          <input v-model="form.description" type="text" placeholder="Объясни почему хочешь это" />
        </div>
        <div v-if="error" class="error-msg">{{ error }}</div>
        <button class="btn-primary" @click="addWish" :disabled="saving">
          {{ saving ? 'Добавление...' : 'Хочу это! 🎁' }}
        </button>
      </div>
    </div>
  </JuniorLayout>
</template>

<script>
import JuniorLayout from '../../../components/JuniorLayout.vue'
import { useApi } from '../../../composables/useApi'
import { useAuth } from '../../../composables/useAuth'
export default {
  name: 'ChildWishlistJunior',
  components: { JuniorLayout },
  data() {
    return {
      loading: true, wishes: [], balance: 0, buying: null,
      showModal: false, saving: false, error: '',
      form: { title: '', description: '' }
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
        this.showModal = false
        this.form = { title: '', description: '' }
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
    }
  }
}
</script>

<style scoped>
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.page-title { font-size: 26px; font-weight: 800; color: #ea580c; }
.add-icon-btn { width: 40px; height: 40px; border-radius: 50%; background: #ea580c; color: #fff; border: none; font-size: 24px; cursor: pointer; display: flex; align-items: center; justify-content: center; }
.loading { text-align: center; padding: 60px; color: #ea580c; font-size: 18px; }
.empty { text-align: center; color: #aaa; font-size: 18px; padding: 40px; }
.wish-card { background: #fff; border-radius: 20px; padding: 16px; margin-bottom: 14px; box-shadow: 0 2px 8px rgba(234,88,12,0.08); }
.wish-card__top { display: flex; gap: 12px; margin-bottom: 12px; align-items: flex-start; }
.wish-card__emoji { font-size: 32px; }
.wish-card__title { font-size: 18px; font-weight: 800; }
.wish-card__desc { font-size: 13px; color: #888; margin-top: 2px; }
.wish-price-row { display: flex; align-items: center; gap: 12px; margin-bottom: 10px; flex-wrap: wrap; }
.wish-price { font-size: 22px; font-weight: 900; color: #ea580c; }
.wish-status-badges { display: flex; align-items: center; gap: 6px; flex-wrap: wrap; }
.step { font-size: 12px; color: #aaa; font-weight: 600; }
.step--done { color: #ea580c; }
.arrow { color: #fed7aa; font-size: 14px; }
.progress-bar { height: 12px; background: #fed7aa; border-radius: 6px; margin-bottom: 12px; overflow: hidden; }
.progress-fill { height: 100%; background: #ea580c; border-radius: 6px; transition: width 0.4s; }
.wish-actions { display: flex; gap: 8px; }
.buy-btn { flex: 1; padding: 12px; background: #22c55e; color: #fff; border: none; border-radius: 14px; font-size: 16px; font-weight: 700; cursor: pointer; transition: background 0.2s; }
.buy-btn:hover:not(:disabled) { background: #16a34a; }
.buy-btn:disabled { background: #fed7aa; color: #ea580c; cursor: not-allowed; }
.badge { padding: 8px 16px; border-radius: 12px; font-size: 14px; font-weight: 700; }
.badge--purchased { background: #fef3c7; color: #92400e; }
.badge--delivered { background: #dcfce7; color: #166534; }
.badge--wait { background: #f5f5f5; color: #888; }
.add-btn { width: 100%; padding: 16px; background: #fff5eb; border: 3px dashed #ea580c; border-radius: 20px; color: #ea580c; font-size: 18px; font-weight: 700; cursor: pointer; margin-top: 8px; }
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center; z-index: 300; }
.modal { background: #fff; border-radius: 20px; padding: 24px; width: 90%; max-width: 360px; }
.modal h3 { font-size: 22px; font-weight: 800; color: #ea580c; margin-bottom: 16px; text-align: center; }
.field { margin-bottom: 14px; }
.field label { display: block; font-size: 14px; font-weight: 600; color: #666; margin-bottom: 6px; }
.field input { width: 100%; padding: 12px; border: 2px solid #fed7aa; border-radius: 12px; font-size: 16px; outline: none; }
.field input:focus { border-color: #ea580c; }
.error-msg { color: #e53e3e; font-size: 14px; margin-bottom: 10px; text-align: center; }
.btn-primary { width: 100%; padding: 14px; background: #ea580c; color: #fff; border: none; border-radius: 14px; font-size: 18px; font-weight: 800; cursor: pointer; }
.btn-primary:hover:not(:disabled) { background: #c2410c; }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
</style>