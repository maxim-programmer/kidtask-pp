<template>
  <JuniorLayout>
    <h1 class="page-title">Мой вишлист 🎁</h1>

    <div class="balance-banner">
      <span class="balance-label">Мой баланс</span>
      <span class="balance-value">⭐ {{ balance }}</span>
    </div>

    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else>
      <div v-if="wishes.length === 0" class="empty">
        <div class="empty__icon">🌟</div>
        <p>Добавь свою первую цель!</p>
      </div>

      <div v-for="w in wishes" :key="w.wish_id" :class="['wish-card', `wish-card--${w.status}`]">
        <div class="wish-card__header">
          <div class="wish-card__icon">{{ wishIcon(w.status) }}</div>
          <div class="wish-card__info">
            <div class="wish-card__title">{{ w.title }}</div>
            <div class="wish-card__desc" v-if="w.description">{{ w.description }}</div>
          </div>
          <span :class="['wish-badge', `wish-badge--${w.status}`]">{{ wishLabel(w.status) }}</span>
        </div>

        <div class="steps">
          <span class="step step--done">⭐ Создано</span>
          <span class="step-arrow">→</span>
          <span :class="['step', { 'step--done': w.status === 'purchased' || w.status === 'delivered' }]">🛒 Куплено</span>
          <span class="step-arrow">→</span>
          <span :class="['step', { 'step--done': w.status === 'delivered' }]">🎁 Доставлено</span>
        </div>

        <div v-if="w.price" class="progress-wrap">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: calcProgress(w) + '%' }"></div>
          </div>
          <div class="progress-labels">
            <span>⭐ {{ balance }}</span>
            <span class="progress-pct">{{ calcProgress(w) }}%</span>
            <span>⭐ {{ w.price }}</span>
          </div>
        </div>

        <div class="wish-card__footer">
          <div v-if="w.status === 'awaiting_price'" class="no-price-msg">
            Ждём оценку от родителя ⏳
          </div>
          <div v-else-if="w.status === 'available'" class="buy-row">
            <span class="wish-price-big">⭐ {{ w.price }}</span>
            <button class="buy-btn" @click="purchase(w)"
              :disabled="balance < w.price || buying === w.wish_id">
              {{ buying === w.wish_id ? '...' : balance >= w.price ? '🛒 Купить!' : `Ещё ⭐${w.price - balance}` }}
            </button>
          </div>
          <div v-else-if="w.status === 'purchased'" class="status-msg status-msg--bought">
            🛒 Уже куплено, скоро доставят!
          </div>
          <div v-else-if="w.status === 'delivered'" class="status-msg status-msg--delivered">
            🎉 Получено! Ура!
          </div>
        </div>
      </div>

      <button class="add-btn" @click="showModal = true">
        <span class="add-btn__icon">+</span> Добавить цель
      </button>
    </div>

    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <div class="modal__handle"></div>
        <div class="modal__icon">🎯</div>
        <h3>Новая цель!</h3>
        <div class="field">
          <label>Что ты хочешь?</label>
          <input v-model="form.title" type="text" placeholder="Nintendo Switch, велосипед..." />
        </div>
        <div class="field">
          <label>Почему хочешь это?</label>
          <textarea v-model="form.description" rows="2" placeholder="Объясни родителям..."></textarea>
        </div>
        <div v-if="wishError" class="error-msg">{{ wishError }}</div>
        <button class="modal-submit-btn" @click="addWish" :disabled="saving || !form.title.trim()">
          {{ saving ? 'Добавляем...' : '🎁 Хочу это!' }}
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
      loading: true, wishes: [],
      showModal: false, saving: false, wishError: '',
      form: { title: '', description: '' },
      buying: null,
    }
  },
  computed: {
    childId() { const u = useAuth().user.value; return u?.child_id || u?.user_id || u?.id },
    balance() { return useAuth().user.value?.balance || 0 }
  },
  async mounted() { await this.load() },
  methods: {
    async load() {
      const { getWishes } = useApi()
      this.loading = true
      try { this.wishes = (await getWishes(this.childId)).data.wishes || [] }
      finally { this.loading = false }
    },
    calcProgress(w) {
      if (!w.price) return 0
      if (w.status === 'purchased' || w.status === 'delivered') return 100
      return Math.min(100, Math.round((this.balance / w.price) * 100))
    },
    wishIcon(s) {
      return { awaiting_price: '⏳', available: '🌟', purchased: '🛒', delivered: '🎁' }[s] || '🎁'
    },
    wishLabel(s) {
      return { awaiting_price: 'Без цены', available: 'Доступно', purchased: 'Куплено', delivered: 'Доставлено' }[s] || s
    },
    async purchase(w) {
      this.buying = w.wish_id
      const { purchaseWish, getMe } = useApi()
      const { updateBalance } = useAuth()
      try {
        await purchaseWish(this.childId, w.wish_id)
        const me = await getMe()
        updateBalance(me.data.user.balance)
        await this.load()
      }
      catch (e) { alert(e.response?.data?.error?.message || 'Ошибка') }
      finally { this.buying = null }
    },
    async addWish() {
      if (!this.form.title.trim()) return
      this.saving = true; this.wishError = ''
      const { createWish } = useApi()
      try {
        await createWish(this.childId, { title: this.form.title.trim(), description: this.form.description || undefined })
        this.showModal = false; this.form = { title: '', description: '' }
        await this.load()
      } catch (e) { this.wishError = e.response?.data?.error?.message || 'Ошибка' }
      finally { this.saving = false }
    }
  }
}
</script>

<style scoped>
.page-title { font-size: 24px; font-weight: 800; color: #1a1a1a; margin-bottom: 12px; }
.balance-banner { background: linear-gradient(135deg, #ea580c, #f97316); border-radius: 14px; padding: 14px 20px; display: flex; align-items: center; justify-content: space-between; margin-bottom: 18px; }
.balance-label { font-size: 14px; font-weight: 600; color: rgba(255,255,255,0.8); }
.balance-value { font-size: 24px; font-weight: 900; color: #fff; }
.loading { text-align: center; padding: 60px; color: #ea580c; font-size: 18px; }
.empty { text-align: center; padding: 50px 20px; }
.empty__icon { font-size: 52px; margin-bottom: 10px; }
.empty p { font-size: 16px; color: #aaa; }

.wish-card { background: #fff; border-radius: 20px; padding: 16px; margin-bottom: 14px; box-shadow: 0 2px 10px rgba(234,88,12,0.09); }
.wish-card--available { border: 2px solid #fed7aa; }
.wish-card--delivered { border: 2px solid #22c55e; }
.wish-card__header { display: flex; align-items: flex-start; gap: 10px; margin-bottom: 12px; }
.wish-card__icon { font-size: 30px; flex-shrink: 0; }
.wish-card__info { flex: 1; min-width: 0; }
.wish-card__title { font-size: 17px; font-weight: 800; color: #1a1a1a; }
.wish-card__desc { font-size: 13px; color: #888; margin-top: 2px; }
.wish-badge { font-size: 11px; font-weight: 700; padding: 3px 8px; border-radius: 8px; flex-shrink: 0; }
.wish-badge--awaiting_price { background: #fef3c7; color: #92400e; }
.wish-badge--available { background: #fff0e6; color: #ea580c; }
.wish-badge--purchased { background: #dbeafe; color: #1e40af; }
.wish-badge--delivered { background: #dcfce7; color: #166534; }

.steps { display: flex; align-items: center; gap: 6px; margin-bottom: 12px; flex-wrap: wrap; }
.step { font-size: 12px; color: #d1d5db; font-weight: 700; }
.step--done { color: #ea580c; }
.step-arrow { color: #e5e7eb; font-size: 13px; }

.progress-wrap { margin-bottom: 12px; }
.progress-bar { height: 12px; background: #fed7aa; border-radius: 6px; overflow: hidden; margin-bottom: 6px; }
.progress-fill { height: 100%; background: linear-gradient(90deg, #ea580c, #f97316); border-radius: 6px; transition: width 0.4s; }
.progress-labels { display: flex; justify-content: space-between; font-size: 12px; color: #888; }
.progress-pct { font-weight: 700; color: #ea580c; }

.wish-card__footer { }
.no-price-msg { color: #92400e; font-size: 14px; font-weight: 600; background: #fef9c3; border-radius: 10px; padding: 10px 14px; }
.buy-row { display: flex; align-items: center; justify-content: space-between; gap: 10px; }
.wish-price-big { font-size: 22px; font-weight: 900; color: #ea580c; }
.buy-btn { padding: 11px 22px; background: #22c55e; color: #fff; border: none; border-radius: 14px; font-size: 16px; font-weight: 800; cursor: pointer; font-family: inherit; transition: transform 0.1s; }
.buy-btn:active:not(:disabled) { transform: scale(0.95); }
.buy-btn:disabled { background: #fed7aa; color: #ea580c; cursor: not-allowed; font-size: 13px; }
.status-msg { font-size: 15px; font-weight: 700; border-radius: 10px; padding: 10px 14px; text-align: center; }
.status-msg--bought { background: #dbeafe; color: #1e40af; }
.status-msg--delivered { background: #dcfce7; color: #166534; }

.add-btn { width: 100%; padding: 18px; background: #fff5eb; border: 3px dashed #ea580c; border-radius: 20px; color: #ea580c; font-size: 17px; font-weight: 700; cursor: pointer; margin-top: 6px; font-family: inherit; display: flex; align-items: center; justify-content: center; gap: 8px; }
.add-btn__icon { font-size: 22px; font-weight: 900; }

.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: flex-end; justify-content: center; z-index: 300; }
.modal { background: #fff; border-radius: 24px 24px 0 0; padding: 8px 20px 32px; width: 100%; max-width: 600px; max-height: 92vh; overflow-y: auto; }
.modal__handle { width: 40px; height: 4px; background: #fed7aa; border-radius: 2px; margin: 0 auto 10px; }
.modal__icon { font-size: 40px; text-align: center; margin-bottom: 6px; }
.modal h3 { font-size: 22px; font-weight: 800; color: #ea580c; margin-bottom: 18px; text-align: center; }
.field { margin-bottom: 14px; }
.field label { display: block; font-size: 14px; font-weight: 700; color: #555; margin-bottom: 6px; }
.field input, .field textarea { width: 100%; padding: 13px 14px; border: 2px solid #fed7aa; border-radius: 14px; font-size: 16px; outline: none; font-family: inherit; box-sizing: border-box; resize: vertical; }
.field input:focus, .field textarea:focus { border-color: #ea580c; }
.error-msg { color: #e53e3e; font-size: 13px; margin-bottom: 10px; text-align: center; }
.modal-submit-btn { width: 100%; padding: 16px; background: #ea580c; color: #fff; border: none; border-radius: 16px; font-size: 18px; font-weight: 800; cursor: pointer; font-family: inherit; }
.modal-submit-btn:disabled { opacity: 0.5; cursor: not-allowed; }
@media (min-width: 600px) {
  .modal-overlay { align-items: center; }
  .modal { border-radius: 20px; max-width: 380px; }
  .modal__handle { display: none; }
}
</style>
