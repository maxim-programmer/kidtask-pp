<template>
  <JuniorLayout>
    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else>
      <div class="view-tabs">
        <button :class="['tab', { 'tab--active': view === 'progress' }]" @click="view = 'progress'">📊</button>
        <button :class="['tab', { 'tab--active': view === 'tasks' }]" @click="view = 'tasks'">📋</button>
        <button :class="['tab', { 'tab--active': view === 'wishlist' }]" @click="view = 'wishlist'">🎁</button>
        <button :class="['tab', { 'tab--active': view === 'history' }]" @click="loadHistory">📜</button>
        <button :class="['tab', { 'tab--active': view === 'chat' }]" @click="loadChat">💬</button>
      </div>
      <div class="view-labels">
        <span>Прогресс</span><span>Задания</span><span>Вишлист</span><span>История</span><span>Чат</span>
      </div>

      <div v-if="view === 'progress'">
        <div class="balance-card">
          <div class="balance-card__label">Мой баланс</div>
          <div class="balance-card__value">⭐ {{ user?.balance || 0 }}</div>
          <div class="balance-card__sub">монет накоплено</div>
        </div>
        <div class="progress-row">
          <div class="progress-item">
            <div class="circle-wrap">
              <svg viewBox="0 0 36 36"><circle cx="18" cy="18" r="15" fill="none" stroke="#fed7aa" stroke-width="3"/><circle cx="18" cy="18" r="15" fill="none" stroke="#ea580c" stroke-width="3" :stroke-dasharray="`${wishPct} 100`" stroke-linecap="round" transform="rotate(-90 18 18)"/></svg>
              <span class="circle-text">{{ wishDone }}/{{ wishTotal }}</span>
            </div>
            <span class="progress-label">Целей</span>
          </div>
          <div class="progress-item">
            <div class="circle-wrap">
              <svg viewBox="0 0 36 36"><circle cx="18" cy="18" r="15" fill="none" stroke="#fed7aa" stroke-width="3"/><circle cx="18" cy="18" r="15" fill="none" stroke="#ea580c" stroke-width="3" :stroke-dasharray="`${taskPct} 100`" stroke-linecap="round" transform="rotate(-90 18 18)"/></svg>
              <span class="circle-text">{{ taskDone }}/{{ taskTotal }}</span>
            </div>
            <span class="progress-label">Заданий</span>
          </div>
        </div>
        <div class="quick-actions">
          <button class="quick-btn" @click="view = 'tasks'">📋 Мои задания</button>
          <button class="quick-btn" @click="view = 'wishlist'">🎁 Вишлист</button>
        </div>
      </div>

      <div v-if="view === 'tasks'">
        <div class="filter-bar">
          <button v-for="f in taskFilters" :key="f.value" :class="['chip', { 'chip--active': taskFilter === f.value }]" @click="taskFilter = f.value">{{ f.label }}</button>
        </div>
        <div v-if="filteredTasks.length === 0" class="empty">Нет заданий 🎉</div>
        <div v-for="t in filteredTasks" :key="t.task_id" :class="['task-card', `task-card--${t.status}`]">
          <div class="task-card__top">
            <div>
              <div class="task-card__title">{{ t.title }}</div>
              <div class="task-card__desc" v-if="t.description">{{ t.description }}</div>
              <div class="task-card__rework" v-if="t.rejection_comment">💬 {{ t.rejection_comment }}</div>
            </div>
          </div>
          <div class="task-card__footer">
            <span class="reward">⭐ {{ t.reward }}</span>
            <button v-if="t.status === 'active' || t.status === 'needs_rework'" class="done-btn" @click="submit(t)" :disabled="submitting === t.task_id">{{ submitting === t.task_id ? '...' : '✓ Сделал!' }}</button>
            <span v-else-if="t.status === 'pending_review'" class="badge badge--pending">На проверке ⏳</span>
            <span v-else-if="t.status === 'completed'" class="badge badge--done">Выполнено ✓</span>
          </div>
        </div>
      </div>

      <div v-if="view === 'wishlist'">
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
              <span class="step step--done">⭐ Создано</span>
              <span class="arrow">→</span>
              <span :class="['step', { 'step--done': w.status === 'purchased' || w.status === 'delivered' }]">🛒 Куплено</span>
              <span class="arrow">→</span>
              <span :class="['step', { 'step--done': w.status === 'delivered' }]">🎁 Доставлено</span>
            </div>
          </div>
          <div v-if="w.price" class="progress-bar">
            <div class="progress-fill" :style="{ width: wishProgress(w) + '%' }"></div>
          </div>
          <div class="wish-actions">
            <button v-if="w.status === 'available'" class="buy-btn" @click="purchase(w)" :disabled="(user?.balance || 0) < w.price || buying === w.wish_id">
              {{ buying === w.wish_id ? '...' : (user?.balance || 0) >= w.price ? '🛒 Купить!' : `Нужно ещё ⭐ ${w.price - (user?.balance || 0)}` }}
            </button>
            <span v-else-if="w.status === 'purchased'" class="badge badge--purchased">🛒 Куплено</span>
            <span v-else-if="w.status === 'delivered'" class="badge badge--delivered">🎁 Доставлено!</span>
            <span v-else-if="w.status === 'awaiting_price'" class="badge badge--wait">Ждём цену ⏳</span>
          </div>
        </div>
        <button class="add-btn" @click="showModal = true">+ Добавить цель</button>
      </div>

      <div v-if="view === 'history'">
        <div v-if="historyLoading" class="loading">Загрузка...</div>
        <div v-else-if="balanceLogs.length === 0" class="empty">История пуста</div>
        <div v-for="log in balanceLogs" :key="log.log_id" :class="['log-item', log.delta > 0 ? 'log-item--plus' : 'log-item--minus']">
          <div class="log-icon">{{ log.delta > 0 ? '⭐' : '🛒' }}</div>
          <div class="log-info">
            <div class="log-reason">{{ log.reason }}</div>
            <div class="log-date">{{ formatDate(log.created_at) }}</div>
          </div>
          <div class="log-delta">{{ log.delta > 0 ? '+' : '' }}{{ log.delta }} ⭐</div>
        </div>
      </div>

      <div v-if="view === 'chat'" class="chat-wrap">
        <div class="chat-messages" ref="chatMessages">
          <div v-if="chatLoading" class="chat-empty">Загрузка...</div>
          <div v-else-if="chatMsgs.length === 0" class="chat-empty">Напиши что-нибудь родителю 💬</div>
          <div v-for="m in chatMsgs" :key="m.message_id" :class="['chat-bubble', m.from_child ? 'chat-bubble--me' : 'chat-bubble--parent']">
            <div class="chat-bubble__name">{{ m.from_child ? 'Я' : 'Родитель' }}</div>
            <div class="chat-bubble__body">{{ m.body }}</div>
            <div class="chat-bubble__time">{{ formatDate(m.created_at) }}</div>
          </div>
        </div>
        <div class="chat-input-row">
          <input v-model="chatInput" type="text" placeholder="Написать сообщение..." class="chat-input" @keyup.enter="sendChat" />
          <button class="chat-send-btn" @click="sendChat" :disabled="!chatInput.trim()">➤</button>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <div class="modal__handle"></div>
        <h3>Новая цель 🎯</h3>
        <div class="field"><label>Название</label><input v-model="form.title" type="text" placeholder="Nintendo Switch" /></div>
        <div class="field"><label>Описание (необязательно)</label><input v-model="form.description" type="text" placeholder="Объясни почему хочешь это" /></div>
        <div v-if="wishError" class="error-msg">{{ wishError }}</div>
        <button class="btn-primary" @click="addWish" :disabled="saving">{{ saving ? 'Добавление...' : 'Хочу это! 🎁' }}</button>
      </div>
    </div>
  </JuniorLayout>
</template>

<script>
import JuniorLayout from '../../../components/JuniorLayout.vue'
import { useApi } from '../../../composables/useApi'
import { useAuth } from '../../../composables/useAuth'
export default {
  name: 'ChildDashboardJunior',
  components: { JuniorLayout },
  data() {
    return {
      loading: true, tasks: [], wishes: [], user: null,
      view: 'progress', taskFilter: '',
      taskFilters: [{ value: '', label: 'Все' }, { value: 'active', label: 'Активные' }, { value: 'pending_review', label: 'На проверке' }, { value: 'needs_rework', label: 'Доработка' }, { value: 'completed', label: 'Готово' }],
      submitting: null, buying: null,
      showModal: false, saving: false, wishError: '',
      form: { title: '', description: '' },
      balanceLogs: [], historyLoading: false,
      chatMsgs: [], chatLoading: false, chatInput: '',
    }
  },
  computed: {
    childId() { const u = useAuth().user.value; return u?.child_id || u?.user_id || u?.id },
    filteredTasks() { return this.taskFilter ? this.tasks.filter(t => t.status === this.taskFilter) : this.tasks },
    wishDone() { return this.wishes.filter(w => w.status === 'delivered').length },
    wishTotal() { return this.wishes.length },
    taskDone() { return this.tasks.filter(t => t.status === 'completed').length },
    taskTotal() { return this.tasks.length },
    wishPct() { return this.wishTotal ? Math.round((this.wishDone / this.wishTotal) * 100) : 0 },
    taskPct() { return this.taskTotal ? Math.round((this.taskDone / this.taskTotal) * 100) : 0 },
  },
  async mounted() { await this.load() },
  methods: {
    async load() {
      const { getMe, getTasks, getWishes } = useApi()
      const authUser = useAuth().user.value
      this.loading = true
      try {
        const [me, tr, wr] = await Promise.all([getMe(), getTasks(), getWishes(authUser.child_id || authUser.user_id || authUser.id)])
        this.user = me.data.user; this.tasks = tr.data.tasks || []; this.wishes = wr.data.wishes || []
      } catch {
        const tr = await getTasks().catch(() => ({ data: { tasks: [] } }))
        this.tasks = tr.data.tasks || []; this.user = authUser
      } finally { this.loading = false }
    },
    wishProgress(w) {
      if (!w.price) return 0
      if (w.status === 'purchased' || w.status === 'delivered') return 100
      return Math.min(100, Math.round(((this.user?.balance || 0) / w.price) * 100))
    },
    async submit(task) { this.submitting = task.task_id; const { submitTask } = useApi(); try { await submitTask(task.task_id); await this.load() } catch (e) { alert(e.response?.data?.error?.message || 'Попробуй снова') } finally { this.submitting = null } },
    async purchase(wish) {
      this.buying = wish.wish_id
      const { purchaseWish, getMe } = useApi()
      const { updateBalance } = useAuth()
      try {
        await purchaseWish(this.childId, wish.wish_id)
        const me = await getMe()
        updateBalance(me.data.user.balance)
        await this.load()
      } catch (e) { alert(e.response?.data?.error?.message || 'Ошибка') }
      finally { this.buying = null }
    },
    async loadHistory() {
      this.view = 'history'
      if (this.balanceLogs.length) return
      this.historyLoading = true
      const { getMyBalanceLogs } = useApi()
      try { const res = await getMyBalanceLogs(); this.balanceLogs = res.data.logs || [] } finally { this.historyLoading = false }
    },
    async loadChat() {
      this.view = 'chat'; this.chatLoading = true
      const { getMyChat } = useApi()
      try { const res = await getMyChat(); this.chatMsgs = res.data.messages || []; this.$nextTick(() => this.scrollChat()) } finally { this.chatLoading = false }
    },
    async sendChat() {
      if (!this.chatInput.trim()) return
      const { sendMyChat } = useApi()
      const res = await sendMyChat(this.chatInput.trim())
      this.chatMsgs.push(res.data.message); this.chatInput = ''
      this.$nextTick(() => this.scrollChat())
    },
    scrollChat() { const el = this.$refs.chatMessages; if (el) el.scrollTop = el.scrollHeight },
    formatDate(d) { return new Date(d).toLocaleString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' }) },
    async addWish() {
      this.saving = true; this.wishError = ''
      const { createWish } = useApi()
      try { await createWish(this.childId, { title: this.form.title, description: this.form.description || undefined }); this.showModal = false; this.form = { title: '', description: '' }; await this.load() }
      catch (e) { this.wishError = e.response?.data?.error?.message || 'Ошибка' }
      finally { this.saving = false }
    }
  }
}
</script>

<style scoped>
.loading { text-align: center; padding: 60px; color: #ea580c; font-size: 18px; }

.view-tabs { display: flex; background: #fde8d0; border-radius: 16px; padding: 4px; gap: 4px; margin-bottom: 4px; }
.tab { flex: 1; padding: 12px 4px; border: none; background: transparent; color: #c2410c; font-size: 20px; cursor: pointer; border-radius: 12px; transition: all 0.15s; font-family: inherit; }
.tab--active { background: #fff; box-shadow: 0 2px 8px rgba(234,88,12,0.15); }
.view-labels { display: flex; margin-bottom: 16px; }
.view-labels span { flex: 1; text-align: center; font-size: 10px; color: #c2410c; font-weight: 600; }

.balance-card { background: linear-gradient(135deg, #ea580c, #f97316); border-radius: 20px; padding: 22px 24px; margin-bottom: 16px; text-align: center; box-shadow: 0 6px 20px rgba(234,88,12,0.3); }
.balance-card__label { font-size: 14px; color: rgba(255,255,255,0.8); font-weight: 600; margin-bottom: 4px; }
.balance-card__value { font-size: 44px; font-weight: 900; color: #fff; line-height: 1; margin-bottom: 4px; }
.balance-card__sub { font-size: 12px; color: rgba(255,255,255,0.7); }

.progress-row { display: flex; gap: 12px; margin-bottom: 16px; }
.progress-item { background: #fff; border-radius: 16px; padding: 18px; flex: 1; display: flex; flex-direction: column; align-items: center; gap: 8px; box-shadow: 0 2px 8px rgba(234,88,12,0.1); }
.circle-wrap { position: relative; width: 72px; height: 72px; }
.circle-wrap svg { width: 100%; height: 100%; }
.circle-text { position: absolute; top: 50%; left: 50%; transform: translate(-50%,-50%); font-size: 13px; font-weight: 800; color: #ea580c; }
.progress-label { font-size: 14px; font-weight: 700; color: #555; }

.quick-actions { display: flex; gap: 12px; }
.quick-btn { flex: 1; padding: 14px; background: #fff; border: 2px solid #fed7aa; border-radius: 16px; color: #ea580c; font-size: 15px; font-weight: 700; cursor: pointer; font-family: inherit; transition: all 0.15s; }
.quick-btn:hover { background: #fff7f0; border-color: #ea580c; }

.filter-bar { display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 12px; overflow-x: auto; -webkit-overflow-scrolling: touch; padding-bottom: 4px; }
.chip { padding: 6px 13px; border-radius: 20px; border: 2px solid #fed7aa; background: #fff; color: #ea580c; font-size: 13px; font-weight: 600; cursor: pointer; transition: all 0.15s; font-family: inherit; white-space: nowrap; }
.chip--active { background: #ea580c; color: #fff; border-color: #ea580c; }

.empty { text-align: center; color: #aaa; font-size: 18px; padding: 40px; }

.task-card { background: #fff; border-radius: 20px; padding: 16px; margin-bottom: 12px; box-shadow: 0 2px 8px rgba(234,88,12,0.08); border-left: 5px solid #fed7aa; }
.task-card--active { border-left-color: #ea580c; }
.task-card--pending_review { border-left-color: #f59e0b; }
.task-card--needs_rework { border-left-color: #e53e3e; }
.task-card--completed { border-left-color: #22c55e; opacity: 0.85; }
.task-card__top { margin-bottom: 10px; }
.task-card__title { font-size: 17px; font-weight: 800; color: #1a1a1a; margin-bottom: 4px; }
.task-card__desc { font-size: 14px; color: #888; }
.task-card__rework { font-size: 13px; color: #e53e3e; margin-top: 6px; background: #fff5f5; border-radius: 8px; padding: 6px 10px; }
.task-card__footer { display: flex; align-items: center; justify-content: space-between; }
.reward { font-size: 20px; font-weight: 800; color: #ea580c; }
.done-btn { padding: 10px 20px; background: #22c55e; color: #fff; border: none; border-radius: 12px; font-size: 15px; font-weight: 700; cursor: pointer; font-family: inherit; }
.done-btn:disabled { opacity: 0.6; cursor: not-allowed; }
.badge { padding: 8px 14px; border-radius: 12px; font-size: 13px; font-weight: 600; }
.badge--pending { background: #fef3c7; color: #92400e; }
.badge--done { background: #dcfce7; color: #166534; }

.wish-card { background: #fff; border-radius: 20px; padding: 16px; margin-bottom: 14px; box-shadow: 0 2px 8px rgba(234,88,12,0.08); }
.wish-card__top { display: flex; gap: 12px; margin-bottom: 12px; align-items: flex-start; }
.wish-card__emoji { font-size: 32px; }
.wish-card__title { font-size: 17px; font-weight: 800; }
.wish-card__desc { font-size: 13px; color: #888; margin-top: 2px; }
.wish-price-row { display: flex; align-items: center; gap: 10px; margin-bottom: 10px; flex-wrap: wrap; }
.wish-price { font-size: 20px; font-weight: 900; color: #ea580c; }
.wish-status-badges { display: flex; align-items: center; gap: 5px; flex-wrap: wrap; }
.step { font-size: 11px; color: #aaa; font-weight: 600; }
.step--done { color: #ea580c; }
.arrow { color: #fed7aa; font-size: 12px; }
.progress-bar { height: 10px; background: #fed7aa; border-radius: 5px; margin-bottom: 12px; overflow: hidden; }
.progress-fill { height: 100%; background: #ea580c; border-radius: 5px; transition: width 0.4s; }
.wish-actions { display: flex; gap: 8px; }
.buy-btn { flex: 1; padding: 12px; background: #22c55e; color: #fff; border: none; border-radius: 14px; font-size: 15px; font-weight: 700; cursor: pointer; font-family: inherit; }
.buy-btn:disabled { background: #fed7aa; color: #ea580c; cursor: not-allowed; }
.badge--purchased { background: #fef3c7; color: #92400e; padding: 8px 14px; border-radius: 12px; font-size: 13px; font-weight: 700; }
.badge--delivered { background: #dcfce7; color: #166534; padding: 8px 14px; border-radius: 12px; font-size: 13px; font-weight: 700; }
.badge--wait { background: #f5f5f5; color: #888; padding: 8px 14px; border-radius: 12px; font-size: 13px; font-weight: 700; }

.add-btn { width: 100%; padding: 16px; background: #fff5eb; border: 3px dashed #ea580c; border-radius: 20px; color: #ea580c; font-size: 17px; font-weight: 700; cursor: pointer; margin-top: 8px; font-family: inherit; }

.log-item { display: flex; align-items: center; gap: 12px; background: #fff; border-radius: 16px; padding: 14px 16px; margin-bottom: 10px; box-shadow: 0 2px 8px rgba(234,88,12,0.07); }
.log-item--plus { border-left: 4px solid #22c55e; }
.log-item--minus { border-left: 4px solid #ea580c; }
.log-icon { font-size: 22px; flex-shrink: 0; }
.log-info { flex: 1; min-width: 0; }
.log-reason { font-size: 13px; font-weight: 600; color: #1a1a1a; }
.log-date { font-size: 11px; color: #aaa; margin-top: 2px; }
.log-delta { font-size: 17px; font-weight: 800; flex-shrink: 0; }
.log-item--plus .log-delta { color: #22c55e; }
.log-item--minus .log-delta { color: #ea580c; }

.chat-wrap { display: flex; flex-direction: column; height: calc(100vh - 260px); min-height: 300px; background: #fff; border-radius: 20px; box-shadow: 0 2px 8px rgba(234,88,12,0.08); overflow: hidden; }
.chat-messages { flex: 1; overflow-y: auto; padding: 16px; display: flex; flex-direction: column; gap: 10px; }
.chat-empty { text-align: center; color: #bbb; font-size: 15px; margin: auto; }
.chat-bubble { max-width: 80%; display: flex; flex-direction: column; gap: 2px; }
.chat-bubble--me { align-self: flex-end; align-items: flex-end; }
.chat-bubble--parent { align-self: flex-start; align-items: flex-start; }
.chat-bubble__name { font-size: 11px; color: #aaa; font-weight: 600; margin-bottom: 2px; }
.chat-bubble__body { padding: 10px 14px; border-radius: 18px; font-size: 15px; line-height: 1.4; word-break: break-word; }
.chat-bubble--me .chat-bubble__body { background: #ea580c; color: #fff; border-bottom-right-radius: 4px; }
.chat-bubble--parent .chat-bubble__body { background: #fff5eb; color: #1a1a1a; border-bottom-left-radius: 4px; }
.chat-bubble__time { font-size: 10px; color: #bbb; }
.chat-input-row { display: flex; gap: 8px; padding: 12px 16px; border-top: 2px solid #fed7aa; }
.chat-input { flex: 1; padding: 11px 14px; border: 2px solid #fed7aa; border-radius: 24px; font-size: 15px; outline: none; font-family: inherit; }
.chat-input:focus { border-color: #ea580c; }
.chat-send-btn { width: 44px; height: 44px; border-radius: 50%; background: #ea580c; color: #fff; border: none; font-size: 18px; cursor: pointer; flex-shrink: 0; font-family: inherit; }
.chat-send-btn:disabled { background: #fed7aa; cursor: not-allowed; }

.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: flex-end; justify-content: center; z-index: 300; }
.modal { background: #fff; border-radius: 24px 24px 0 0; padding: 8px 20px 28px; width: 100%; max-width: 600px; max-height: 92vh; overflow-y: auto; }
.modal__handle { width: 40px; height: 4px; background: #fed7aa; border-radius: 2px; margin: 0 auto 16px; }
.modal h3 { font-size: 22px; font-weight: 800; color: #ea580c; margin-bottom: 16px; text-align: center; }
.field { margin-bottom: 14px; }
.field label { display: block; font-size: 14px; font-weight: 600; color: #666; margin-bottom: 6px; }
.field input { width: 100%; padding: 13px 14px; border: 2px solid #fed7aa; border-radius: 14px; font-size: 16px; outline: none; font-family: inherit; box-sizing: border-box; }
.field input:focus { border-color: #ea580c; }
.error-msg { color: #e53e3e; font-size: 14px; margin-bottom: 10px; text-align: center; }
.btn-primary { width: 100%; padding: 15px; background: #ea580c; color: #fff; border: none; border-radius: 16px; font-size: 18px; font-weight: 800; cursor: pointer; font-family: inherit; }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }

@media (min-width: 600px) {
  .modal-overlay { align-items: center; }
  .modal { border-radius: 20px; max-width: 380px; }
  .modal__handle { display: none; }
  .chat-wrap { height: 460px; }
}
</style>
