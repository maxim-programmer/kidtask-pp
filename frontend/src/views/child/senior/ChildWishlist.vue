<template>
  <SeniorLayout>
    <div class="page-header">
      <h1 class="page-title">Вишлист</h1>
      <button class="add-btn-header" @click="showModal = true">+ Добавить</button>
    </div>

    <div class="balance-bar">
      <span class="balance-label">Баланс:</span>
      <span class="balance-value">⭐ {{ balance }}</span>
    </div>

    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else>
      <div v-if="wishes.length === 0" class="empty">
        <p>Добавьте первую цель</p>
        <button class="empty-add-btn" @click="showModal = true">+ Добавить цель</button>
      </div>

      <div v-for="w in wishes" :key="w.wish_id" class="wish-card">
        <div class="wish-card__top">
          <div class="wish-info">
            <div class="wish-title">{{ w.title }}</div>
            <div class="wish-desc" v-if="w.description">{{ w.description }}</div>
          </div>
          <div class="wish-right">
            <div class="wish-price" v-if="w.price">⭐ {{ w.price }}</div>
            <div class="wish-price wish-price--unknown" v-else>⭐ ?</div>
            <span :class="['wish-status', `wish-status--${w.status}`]">{{ wishLabel(w.status) }}</span>
          </div>
        </div>

        <div class="steps">
          <div class="step-item step-item--done">
            <div class="step-dot step-dot--done"></div>
            <span>Создано</span>
          </div>
          <div class="step-line" :class="{ 'step-line--done': w.status === 'purchased' || w.status === 'delivered' }"></div>
          <div class="step-item" :class="{ 'step-item--done': w.status === 'purchased' || w.status === 'delivered' }">
            <div class="step-dot" :class="{ 'step-dot--done': w.status === 'purchased' || w.status === 'delivered' }"></div>
            <span>Куплено</span>
          </div>
          <div class="step-line" :class="{ 'step-line--done': w.status === 'delivered' }"></div>
          <div class="step-item" :class="{ 'step-item--done': w.status === 'delivered' }">
            <div class="step-dot" :class="{ 'step-dot--done': w.status === 'delivered' }"></div>
            <span>Доставлено</span>
          </div>
        </div>

        <div v-if="w.price" class="progress-wrap">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: calcProgress(w) + '%' }"></div>
          </div>
          <div class="progress-info">
            <span>{{ Math.min(balance, w.price) }} / {{ w.price }} ⭐</span>
            <span class="progress-pct">{{ calcProgress(w) }}%</span>
          </div>
        </div>

        <div class="wish-footer">
          <template v-if="w.status === 'awaiting_price'">
            <span class="footer-msg footer-msg--wait">Ожидает оценки от родителя</span>
          </template>
          <template v-else-if="w.status === 'available'">
            <button class="buy-btn" @click="purchase(w)"
              :disabled="balance < w.price || buying === w.wish_id">
              <span v-if="buying === w.wish_id">Покупаем...</span>
              <span v-else-if="balance >= w.price">🛒 Купить</span>
              <span v-else>Нужно ещё ⭐{{ w.price - balance }}</span>
            </button>
          </template>
          <template v-else-if="w.status === 'purchased'">
            <span class="footer-msg footer-msg--bought">🛒 Куплено, ожидайте доставку</span>
          </template>
          <template v-else-if="w.status === 'delivered'">
            <span class="footer-msg footer-msg--delivered">🎁 Доставлено!</span>
          </template>
          <button class="del-btn" @click="remove(w)" v-if="w.status !== 'purchased' && w.status !== 'delivered'">🗑</button>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <div class="modal__handle"></div>
        <h3>Новая цель</h3>
        <div class="field">
          <label>Название</label>
          <input v-model="form.title" type="text" placeholder="iPhone 16, PlayStation 5..." />
        </div>
        <div class="field">
          <label>Описание <span class="optional">(необязательно)</span></label>
          <textarea v-model="form.description" rows="3" placeholder="Зачем тебе это..."></textarea>
        </div>
        <div v-if="wishError" class="error-msg">{{ wishError }}</div>
        <div class="modal-actions">
          <button class="btn-cancel" @click="showModal = false; form = { title: '', description: '' }">Отмена</button>
          <button class="btn-submit" @click="addWish" :disabled="saving || !form.title.trim()">
            {{ saving ? 'Добавление...' : 'Добавить' }}
          </button>
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
    async remove(w) {
      if (!confirm('Удалить эту цель?')) return
      const { deleteWish } = useApi()
      await deleteWish(this.childId, w.wish_id)
      await this.load()
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
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 14px; }
.page-title { font-size: 22px; font-weight: 700; color: #f3f4f6; }
.add-btn-header { padding: 9px 18px; background: #6366f1; color: #fff; border: none; border-radius: 10px; font-size: 14px; font-weight: 600; cursor: pointer; font-family: inherit; }
.balance-bar { background: #1f2937; border: 1px solid #374151; border-radius: 12px; padding: 12px 16px; display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.balance-label { font-size: 14px; color: #9ca3af; }
.balance-value { font-size: 20px; font-weight: 800; color: #a5b4fc; }
.loading { text-align: center; padding: 60px; color: #6366f1; }
.empty { text-align: center; padding: 50px 20px; color: #6b7280; }
.empty p { font-size: 16px; margin-bottom: 16px; }
.empty-add-btn { padding: 12px 24px; background: #6366f1; color: #fff; border: none; border-radius: 10px; font-size: 14px; font-weight: 600; cursor: pointer; font-family: inherit; }

.wish-card { background: #1f2937; border-radius: 14px; padding: 16px; margin-bottom: 12px; border: 1px solid #374151; }
.wish-card__top { display: flex; justify-content: space-between; align-items: flex-start; gap: 10px; margin-bottom: 14px; }
.wish-info { flex: 1; min-width: 0; }
.wish-title { font-size: 15px; font-weight: 700; color: #f3f4f6; margin-bottom: 3px; }
.wish-desc { font-size: 13px; color: #9ca3af; }
.wish-right { display: flex; flex-direction: column; align-items: flex-end; gap: 6px; flex-shrink: 0; }
.wish-price { font-size: 18px; font-weight: 800; color: #a5b4fc; }
.wish-price--unknown { color: #6b7280; }
.wish-status { font-size: 11px; font-weight: 600; padding: 3px 8px; border-radius: 6px; }
.wish-status--awaiting_price { background: #78350f30; color: #fbbf24; }
.wish-status--available { background: #1e1b4b; color: #818cf8; }
.wish-status--purchased { background: #1e3a5f; color: #60a5fa; }
.wish-status--delivered { background: #14532d30; color: #4ade80; }

.steps { display: flex; align-items: center; margin-bottom: 12px; }
.step-item { display: flex; flex-direction: column; align-items: center; gap: 4px; }
.step-dot { width: 10px; height: 10px; border-radius: 50%; background: #374151; border: 2px solid #4b5563; }
.step-dot--done { background: #6366f1; border-color: #6366f1; }
.step-item span { font-size: 10px; color: #6b7280; white-space: nowrap; }
.step-item--done span { color: #a5b4fc; }
.step-line { flex: 1; height: 2px; background: #374151; margin: 0 4px; margin-bottom: 14px; }
.step-line--done { background: #6366f1; }

.progress-wrap { margin-bottom: 12px; }
.progress-bar { height: 8px; background: #374151; border-radius: 4px; margin-bottom: 6px; overflow: hidden; }
.progress-fill { height: 100%; background: linear-gradient(90deg, #6366f1, #8b5cf6); border-radius: 4px; transition: width 0.4s; }
.progress-info { display: flex; justify-content: space-between; font-size: 12px; color: #9ca3af; }
.progress-pct { color: #a5b4fc; font-weight: 700; }

.wish-footer { display: flex; align-items: center; justify-content: space-between; gap: 10px; }
.footer-msg { font-size: 13px; font-weight: 600; padding: 8px 12px; border-radius: 8px; flex: 1; text-align: center; }
.footer-msg--wait { background: #78350f20; color: #fbbf24; }
.footer-msg--bought { background: #1e3a5f; color: #60a5fa; }
.footer-msg--delivered { background: #14532d30; color: #4ade80; }
.buy-btn { flex: 1; padding: 11px 16px; background: #22c55e; color: #fff; border: none; border-radius: 10px; font-size: 14px; font-weight: 600; cursor: pointer; font-family: inherit; }
.buy-btn:disabled { background: #374151; color: #6b7280; cursor: not-allowed; }
.del-btn { width: 38px; height: 38px; background: transparent; border: 1px solid #374151; border-radius: 8px; color: #6b7280; font-size: 16px; cursor: pointer; display: flex; align-items: center; justify-content: center; flex-shrink: 0; font-family: inherit; }
.del-btn:hover { background: rgba(239,68,68,0.15); color: #ef4444; border-color: #ef4444; }

.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.65); display: flex; align-items: flex-end; justify-content: center; z-index: 300; }
.modal { background: #1f2937; border: 1px solid #374151; border-radius: 20px 20px 0 0; padding: 8px 20px 28px; width: 100%; max-width: 600px; max-height: 92vh; overflow-y: auto; }
.modal__handle { width: 40px; height: 4px; background: #374151; border-radius: 2px; margin: 0 auto 16px; }
.modal h3 { font-size: 18px; font-weight: 700; color: #e5e7eb; margin-bottom: 16px; }
.field { margin-bottom: 14px; }
.field label { display: block; font-size: 13px; color: #9ca3af; margin-bottom: 6px; font-weight: 500; }
.optional { color: #4b5563; font-size: 11px; font-weight: 400; }
.field input, .field textarea { width: 100%; padding: 12px 14px; background: #111827; border: 1px solid #374151; border-radius: 10px; font-size: 15px; color: #e5e7eb; outline: none; font-family: inherit; box-sizing: border-box; resize: vertical; }
.field input:focus, .field textarea:focus { border-color: #6366f1; }
.error-msg { color: #f87171; font-size: 13px; margin-bottom: 10px; }
.modal-actions { display: flex; gap: 8px; margin-top: 8px; }
.btn-cancel { flex: 1; padding: 12px; background: transparent; border: 1px solid #374151; border-radius: 10px; color: #9ca3af; font-size: 14px; cursor: pointer; font-family: inherit; }
.btn-submit { flex: 1; padding: 12px; background: #6366f1; color: #fff; border: none; border-radius: 10px; font-size: 14px; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-submit:disabled { opacity: 0.5; cursor: not-allowed; }
@media (min-width: 600px) {
  .modal-overlay { align-items: center; }
  .modal { border-radius: 14px; max-width: 400px; }
  .modal__handle { display: none; }
}
</style>
