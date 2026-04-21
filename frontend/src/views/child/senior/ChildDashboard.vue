<template>
  <SeniorLayout>
    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else>
      <div class="view-tabs">
        <button :class="['tab', { 'tab--active': view === 'progress' }]" @click="view = 'progress'">Прогресс</button>
        <button :class="['tab', { 'tab--active': view === 'tasks' }]" @click="view = 'tasks'">Задания</button>
        <button :class="['tab', { 'tab--active': view === 'wishlist' }]" @click="view = 'wishlist'">Вишлист</button>
      </div>

      <!-- ПРОГРЕСС -->
      <div v-if="view === 'progress'">
        <div class="balance-card">
          <div class="balance-label">Мой баланс</div>
          <div class="balance-value">⭐ {{ user?.balance || 0 }}</div>
        </div>
        <div class="progress-row">
          <div class="progress-item">
            <div class="circle-wrap">
              <svg viewBox="0 0 36 36">
                <circle cx="18" cy="18" r="15" fill="none" stroke="#374151" stroke-width="3"/>
                <circle cx="18" cy="18" r="15" fill="none" stroke="#6366f1" stroke-width="3"
                  :stroke-dasharray="`${wishPct} 100`" stroke-linecap="round" transform="rotate(-90 18 18)"/>
              </svg>
              <span class="circle-text">{{ wishDone }}/{{ wishTotal }}</span>
            </div>
            <span class="progress-label">Целей</span>
          </div>
          <div class="progress-item">
            <div class="circle-wrap">
              <svg viewBox="0 0 36 36">
                <circle cx="18" cy="18" r="15" fill="none" stroke="#374151" stroke-width="3"/>
                <circle cx="18" cy="18" r="15" fill="none" stroke="#6366f1" stroke-width="3"
                  :stroke-dasharray="`${taskPct} 100`" stroke-linecap="round" transform="rotate(-90 18 18)"/>
              </svg>
              <span class="circle-text">{{ taskDone }}/{{ taskTotal }}</span>
            </div>
            <span class="progress-label">Заданий</span>
          </div>
        </div>
      </div>

      <!-- ЗАДАНИЯ -->
      <div v-if="view === 'tasks'">
        <div class="filter-bar">
          <button v-for="f in taskFilters" :key="f.value"
            :class="['chip', { 'chip--active': taskFilter === f.value }]"
            @click="taskFilter = f.value">{{ f.label }}</button>
        </div>
        <div v-if="filteredTasks.length === 0" class="empty">Нет заданий</div>
        <div v-for="t in filteredTasks" :key="t.task_id" :class="['task-card', `task-card--${t.status}`]">
          <div class="task-card__header">
            <div class="task-card__left">
              <span class="task-title">{{ t.title }}</span>
              <span class="task-desc" v-if="t.description">{{ t.description }}</span>
              <span class="task-rework" v-if="t.rejection_comment">💬 {{ t.rejection_comment }}</span>
            </div>
            <button class="audio-btn" @click="speak(t.title)">🔊</button>
          </div>
          <div class="task-card__footer">
            <span class="reward">⭐ {{ t.reward }}</span>
            <button
              v-if="t.status === 'active' || t.status === 'needs_rework'"
              class="done-btn"
              @click="submit(t)"
              :disabled="submitting === t.task_id">
              {{ submitting === t.task_id ? '...' : '✓ Сдать' }}
            </button>
            <span v-else-if="t.status === 'pending_review'" class="badge badge--pending">На проверке</span>
            <span v-else-if="t.status === 'completed'" class="badge badge--done">Выполнено</span>
          </div>
        </div>
      </div>

      <!-- ВИШЛИСТ -->
      <div v-if="view === 'wishlist'">
        <div v-if="wishes.length === 0" class="empty">Добавь свою первую цель</div>
        <div v-for="w in wishes" :key="w.wish_id" class="wish-card">
          <div class="wish-card__top">
            <div class="wish-info">
              <div class="wish-title">{{ w.title }}</div>
              <div class="wish-desc" v-if="w.description">{{ w.description }}</div>
            </div>
            <div class="wish-price-col">
              <span class="wish-price">⭐ {{ w.price ?? '?' }}</span>
            </div>
          </div>
          <div class="steps">
            <span class="step step--done">⭐ Создано</span>
            <span class="step-arrow">→</span>
            <span :class="['step', { 'step--done': w.status === 'purchased' || w.status === 'delivered' }]">🛒 Куплено</span>
            <span class="step-arrow">→</span>
            <span :class="['step', { 'step--done': w.status === 'delivered' }]">🎁 Доставлено</span>
          </div>
          <div v-if="w.price" class="progress-bar">
            <div class="progress-fill" :style="{ width: Math.min(100, ((user?.balance || 0) / w.price) * 100) + '%' }"></div>
          </div>
          <div class="wish-actions">
            <button
              v-if="w.status === 'available'"
              class="buy-btn"
              @click="purchase(w)"
              :disabled="(user?.balance || 0) < w.price || buying === w.wish_id">
              {{ (user?.balance || 0) >= (w.price || Infinity) ? '🛒 Купить' : `Нужно ещё ⭐${w.price - (user?.balance || 0)}` }}
            </button>
            <span v-else-if="w.status === 'purchased'" class="badge badge--bought">🛒 Куплено</span>
            <span v-else-if="w.status === 'delivered'" class="badge badge--done">🎁 Доставлено</span>
            <span v-else-if="w.status === 'awaiting_price'" class="badge badge--wait">Ожидает оценки</span>
            <button class="del-btn" @click="removeWish(w)">🗑</button>
          </div>
        </div>
        <button class="add-btn" @click="showModal = true">+ Добавить цель</button>
      </div>
    </div>

    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <h3>Новая цель</h3>
        <div class="field">
          <label>Название</label>
          <input v-model="form.title" type="text" placeholder="iPhone 16" />
        </div>
        <div class="field">
          <label>Описание</label>
          <textarea v-model="form.description" rows="2" class="form-textarea" placeholder="Зачем мне это нужно..."></textarea>
        </div>
        <div v-if="wishError" class="error-msg">{{ wishError }}</div>
        <div class="modal-actions">
          <button class="btn-cancel" @click="showModal = false">Отмена</button>
          <button class="btn-primary" @click="addWish" :disabled="saving">{{ saving ? '...' : 'Добавить' }}</button>
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
  name: 'ChildDashboardSenior',
  components: { SeniorLayout },
  data() {
    return {
      loading: true, tasks: [], wishes: [], user: null,
      view: 'progress',
      taskFilter: '',
      taskFilters: [
        { value: '', label: 'Все' },
        { value: 'active', label: 'Активные' },
        { value: 'pending_review', label: 'На проверке' },
        { value: 'needs_rework', label: 'На доработке' },
        { value: 'completed', label: 'Завершённые' },
      ],
      submitting: null, buying: null,
      showModal: false, saving: false, wishError: '',
      form: { title: '', description: '' },
    }
  },
  computed: {
    childId() {
      const u = useAuth().user.value
      return u?.child_id || u?.user_id || u?.id
    },
    filteredTasks() {
      if (!this.taskFilter) return this.tasks
      return this.tasks.filter(t => t.status === this.taskFilter)
    },
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
        const childId = authUser?.child_id || authUser?.user_id || authUser?.id
        const [me, tr, wr] = await Promise.all([
          getMe(),
          getTasks(),
          getWishes(childId).catch(() => ({ data: { wishes: [] } }))
        ])
        this.user = me.data.user
        this.tasks = tr.data.tasks || []
        this.wishes = wr.data.wishes || []
      } catch {
        this.user = authUser
        const tr = await getTasks().catch(() => ({ data: { tasks: [] } }))
        this.tasks = tr.data.tasks || []
      } finally { this.loading = false }
    },
    async submit(task) {
      if (task.status !== 'active' && task.status !== 'needs_rework') return
      this.submitting = task.task_id
      const { submitTask } = useApi()
      try { await submitTask(task.task_id); await this.load() }
      catch (e) { alert(e.response?.data?.error?.message || 'Ошибка') }
      finally { this.submitting = null }
    },
    async purchase(wish) {
      this.buying = wish.wish_id
      const { purchaseWish } = useApi()
      try { await purchaseWish(this.childId, wish.wish_id); await this.load() }
      catch (e) { alert(e.response?.data?.error?.message || 'Ошибка') }
      finally { this.buying = null }
    },
    async removeWish(wish) {
      if (!confirm('Удалить цель?')) return
      const { deleteWish } = useApi()
      await deleteWish(this.childId, wish.wish_id)
      await this.load()
    },
    async addWish() {
      this.saving = true; this.wishError = ''
      const { createWish } = useApi()
      try {
        await createWish(this.childId, { title: this.form.title, description: this.form.description || undefined })
        this.showModal = false
        this.form = { title: '', description: '' }
        await this.load()
      } catch (e) { this.wishError = e.response?.data?.error?.message || 'Ошибка' }
      finally { this.saving = false }
    },
    speak(text) {
      if (!window.speechSynthesis) return
      const u = new SpeechSynthesisUtterance(text)
      u.lang = 'ru-RU'
      window.speechSynthesis.speak(u)
    }
  }
}
</script>

<style scoped>
.loading { text-align: center; padding: 60px; color: #6366f1; }

.view-tabs { display: flex; gap: 4px; background: #1f2937; border-radius: 12px; padding: 4px; margin-bottom: 20px; }
.tab { flex: 1; padding: 10px; border: none; background: transparent; color: #9ca3af; font-size: 14px; font-weight: 600; cursor: pointer; border-radius: 8px; transition: all 0.15s; font-family: inherit; }
.tab--active { background: #374151; color: #e5e7eb; }

.balance-card { background: linear-gradient(135deg, #4f46e5, #7c3aed); border-radius: 16px; padding: 20px; margin-bottom: 16px; display: flex; align-items: center; justify-content: space-between; }
.balance-label { font-size: 14px; color: rgba(255,255,255,0.75); font-weight: 500; }
.balance-value { font-size: 30px; font-weight: 900; color: #fff; }

.progress-row { display: flex; gap: 12px; }
.progress-item { background: #1f2937; border-radius: 14px; padding: 20px; flex: 1; display: flex; flex-direction: column; align-items: center; gap: 8px; border: 1px solid #374151; }
.circle-wrap { position: relative; width: 72px; height: 72px; }
.circle-wrap svg { width: 100%; height: 100%; }
.circle-text { position: absolute; top: 50%; left: 50%; transform: translate(-50%,-50%); font-size: 12px; font-weight: 700; color: #a5b4fc; }
.progress-label { font-size: 14px; color: #9ca3af; font-weight: 500; }

.filter-bar { display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 12px; }
.chip { padding: 5px 12px; border-radius: 16px; border: 1px solid #374151; background: transparent; color: #9ca3af; font-size: 13px; font-weight: 500; cursor: pointer; transition: all 0.15s; font-family: inherit; }
.chip--active { background: #6366f1; color: #fff; border-color: #6366f1; }

.empty { text-align: center; color: #6b7280; font-size: 15px; padding: 40px; }

.task-card { background: #1f2937; border-radius: 14px; padding: 14px 16px; margin-bottom: 10px; border: 1px solid #374151; border-left: 4px solid #374151; }
.task-card--active { border-left-color: #6366f1; }
.task-card--pending_review { border-left-color: #f59e0b; }
.task-card--needs_rework { border-left-color: #ef4444; }
.task-card--completed { border-left-color: #22c55e; opacity: 0.7; }
.task-card__header { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 10px; }
.task-card__left { flex: 1; display: flex; flex-direction: column; gap: 3px; }
.task-title { font-size: 15px; font-weight: 700; color: #f3f4f6; }
.task-desc { font-size: 13px; color: #6b7280; }
.task-rework { font-size: 12px; color: #f87171; }
.audio-btn { background: #374151; border: none; border-radius: 8px; width: 32px; height: 32px; font-size: 15px; cursor: pointer; flex-shrink: 0; margin-left: 8px; }
.task-card__footer { display: flex; align-items: center; justify-content: space-between; }
.reward { font-size: 16px; font-weight: 700; color: #a5b4fc; }
.done-btn { padding: 8px 18px; background: #22c55e; color: #fff; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; font-family: inherit; }
.done-btn:hover:not(:disabled) { background: #16a34a; }
.done-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.badge { padding: 6px 12px; border-radius: 8px; font-size: 13px; font-weight: 600; }
.badge--pending { background: #78350f30; color: #fbbf24; }
.badge--done { background: #14532d30; color: #4ade80; }

.wish-card { background: #1f2937; border-radius: 14px; padding: 16px; margin-bottom: 12px; border: 1px solid #374151; }
.wish-card__top { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 12px; }
.wish-info { flex: 1; }
.wish-title { font-size: 16px; font-weight: 700; color: #f3f4f6; margin-bottom: 4px; }
.wish-desc { font-size: 13px; color: #6b7280; }
.wish-price-col { margin-left: 12px; }
.wish-price { font-size: 20px; font-weight: 800; color: #a5b4fc; }
.steps { display: flex; align-items: center; gap: 6px; margin-bottom: 10px; flex-wrap: wrap; }
.step { font-size: 12px; color: #4b5563; font-weight: 600; }
.step--done { color: #a5b4fc; }
.step-arrow { color: #374151; }
.progress-bar { height: 6px; background: #374151; border-radius: 3px; margin-bottom: 12px; overflow: hidden; }
.progress-fill { height: 100%; background: #6366f1; border-radius: 3px; transition: width 0.4s; }
.wish-actions { display: flex; gap: 8px; align-items: center; }
.buy-btn { flex: 1; padding: 10px 16px; background: #22c55e; color: #fff; border: none; border-radius: 10px; font-size: 14px; font-weight: 600; cursor: pointer; font-family: inherit; }
.buy-btn:hover:not(:disabled) { background: #16a34a; }
.buy-btn:disabled { background: #374151; color: #6b7280; cursor: not-allowed; }
.badge--bought { background: #78350f30; color: #fbbf24; }
.badge--wait { background: #1e1b4b; color: #818cf8; }
.del-btn { background: transparent; border: 1px solid #374151; border-radius: 8px; padding: 6px 10px; color: #6b7280; font-size: 14px; cursor: pointer; font-family: inherit; }
.del-btn:hover { background: #ef444420; color: #ef4444; border-color: #ef4444; }

.add-btn { width: 100%; padding: 13px; background: transparent; border: 1px dashed #4b5563; border-radius: 12px; color: #6b7280; font-size: 15px; font-weight: 600; cursor: pointer; margin-top: 8px; font-family: inherit; transition: all 0.15s; }
.add-btn:hover { border-color: #6366f1; color: #a5b4fc; }

.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.6); display: flex; align-items: center; justify-content: center; z-index: 300; }
.modal { background: #1f2937; border: 1px solid #374151; border-radius: 14px; padding: 24px; width: 90%; max-width: 360px; }
.modal h3 { font-size: 18px; font-weight: 700; color: #e5e7eb; margin-bottom: 16px; }
.field { margin-bottom: 14px; }
.field label { display: block; font-size: 13px; color: #9ca3af; margin-bottom: 6px; font-weight: 500; }
.field input, .form-textarea { width: 100%; padding: 10px 12px; background: #111827; border: 1px solid #374151; border-radius: 10px; font-size: 15px; color: #e5e7eb; outline: none; font-family: inherit; }
.field input:focus, .form-textarea:focus { border-color: #6366f1; }
.form-textarea { resize: vertical; }
.error-msg { color: #f87171; font-size: 13px; margin-bottom: 10px; }
.modal-actions { display: flex; gap: 8px; margin-top: 16px; justify-content: flex-end; }
.btn-cancel { padding: 10px 18px; background: transparent; border: 1px solid #374151; border-radius: 8px; color: #9ca3af; font-size: 14px; cursor: pointer; font-family: inherit; }
.btn-primary { padding: 10px 20px; background: #6366f1; color: #fff; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-primary:hover:not(:disabled) { background: #4f46e5; }
.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }
</style>