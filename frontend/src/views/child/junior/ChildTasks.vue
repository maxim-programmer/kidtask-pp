<template>
  <JuniorLayout>
    <h1 class="page-title">Мои задания 📋</h1>

    <div class="filter-bar">
      <button v-for="f in filters" :key="f.value"
        :class="['chip', { 'chip--active': active === f.value }]"
        @click="active = f.value">{{ f.label }}</button>
    </div>

    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else>
      <div v-if="filtered.length === 0" class="empty">
        <div class="empty__icon">✨</div>
        <p>{{ active ? 'Нет заданий в этой группе' : 'Пока нет заданий' }}</p>
      </div>

      <div v-for="t in filtered" :key="t.task_id" :class="['task-card', `task-card--${t.status}`]">
        <div class="task-card__icon">{{ statusIcon(t.status) }}</div>
        <div class="task-card__body">
          <div class="task-card__title">{{ t.title }}</div>
          <div class="task-card__desc" v-if="t.description">{{ t.description }}</div>
          <div class="task-card__rework" v-if="t.rejection_comment">
            💬 {{ t.rejection_comment }}
          </div>
          <div class="task-card__meta">
            <span class="reward">⭐ {{ t.reward }}</span>
            <span class="status-tag" :class="`status-tag--${t.status}`">{{ statusLabel(t.status) }}</span>
          </div>
        </div>
        <div class="task-card__action">
          <button
            v-if="t.status === 'active' || t.status === 'needs_rework'"
            class="done-btn"
            @click="submit(t)"
            :disabled="submitting === t.task_id">
            {{ submitting === t.task_id ? '⏳' : '✓' }}
          </button>
          <span v-else-if="t.status === 'pending_review'" class="check-icon">⏳</span>
          <span v-else-if="t.status === 'completed'" class="check-icon check-icon--done">✓</span>
        </div>
      </div>
    </div>
  </JuniorLayout>
</template>

<script>
import JuniorLayout from '../../../components/JuniorLayout.vue'
import { useApi } from '../../../composables/useApi'
export default {
  name: 'ChildTasksJunior',
  components: { JuniorLayout },
  data() {
    return {
      loading: true,
      tasks: [],
      submitting: null,
      active: '',
      filters: [
        { value: '', label: 'Все' },
        { value: 'active', label: 'Активные' },
        { value: 'pending_review', label: 'На проверке' },
        { value: 'needs_rework', label: 'Доработка' },
        { value: 'completed', label: 'Готово' },
      ]
    }
  },
  computed: {
    filtered() {
      return this.active ? this.tasks.filter(t => t.status === this.active) : this.tasks
    }
  },
  async mounted() { await this.load() },
  methods: {
    async load() {
      const { getTasks } = useApi()
      this.loading = true
      try { this.tasks = (await getTasks()).data.tasks || [] }
      finally { this.loading = false }
    },
    async submit(task) {
      this.submitting = task.task_id
      const { submitTask } = useApi()
      try { await submitTask(task.task_id); await this.load() }
      catch (e) { alert(e.response?.data?.error?.message || 'Попробуй снова') }
      finally { this.submitting = null }
    },
    statusLabel(s) {
      return { active: 'Активное', pending_review: 'На проверке', needs_rework: 'Доработка', completed: 'Выполнено' }[s] || s
    },
    statusIcon(s) {
      return { active: '📋', pending_review: '⏳', needs_rework: '🔄', completed: '✅' }[s] || '📋'
    }
  }
}
</script>

<style scoped>
.page-title { font-size: 24px; font-weight: 800; color: #1a1a1a; margin-bottom: 16px; }

.filter-bar { display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 16px; overflow-x: auto; -webkit-overflow-scrolling: touch; padding-bottom: 4px; }
.chip { padding: 7px 14px; border-radius: 20px; border: 2px solid #fed7aa; background: #fff; color: #ea580c; font-size: 13px; font-weight: 600; cursor: pointer; transition: all 0.15s; font-family: inherit; white-space: nowrap; }
.chip--active { background: #ea580c; color: #fff; border-color: #ea580c; }

.loading { text-align: center; padding: 60px; color: #ea580c; font-size: 18px; }

.empty { text-align: center; padding: 50px 20px; }
.empty__icon { font-size: 48px; margin-bottom: 10px; }
.empty p { font-size: 16px; color: #aaa; }

.task-card { display: flex; align-items: flex-start; gap: 12px; background: #fff; border-radius: 20px; padding: 16px; margin-bottom: 12px; box-shadow: 0 2px 8px rgba(234,88,12,0.08); border-left: 5px solid #fed7aa; }
.task-card--active { border-left-color: #ea580c; }
.task-card--pending_review { border-left-color: #f59e0b; }
.task-card--needs_rework { border-left-color: #e53e3e; background: #fff8f8; }
.task-card--completed { border-left-color: #22c55e; opacity: 0.8; }

.task-card__icon { font-size: 26px; flex-shrink: 0; margin-top: 2px; }
.task-card__body { flex: 1; min-width: 0; }
.task-card__title { font-size: 16px; font-weight: 800; color: #1a1a1a; margin-bottom: 4px; }
.task-card__desc { font-size: 13px; color: #888; margin-bottom: 6px; line-height: 1.4; }
.task-card__rework { font-size: 13px; color: #e53e3e; background: #fff5f5; border-radius: 8px; padding: 6px 10px; margin-bottom: 8px; line-height: 1.4; }
.task-card__meta { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.reward { font-size: 18px; font-weight: 800; color: #ea580c; }
.status-tag { font-size: 11px; font-weight: 700; padding: 3px 8px; border-radius: 8px; }
.status-tag--active { background: #fff0e6; color: #ea580c; }
.status-tag--pending_review { background: #fef3c7; color: #92400e; }
.status-tag--needs_rework { background: #ffe4e6; color: #be123c; }
.status-tag--completed { background: #dcfce7; color: #166534; }

.task-card__action { flex-shrink: 0; display: flex; align-items: center; }
.done-btn { width: 44px; height: 44px; background: #22c55e; border: none; border-radius: 50%; color: #fff; font-size: 20px; font-weight: 700; cursor: pointer; display: flex; align-items: center; justify-content: center; transition: transform 0.1s; font-family: inherit; }
.done-btn:active:not(:disabled) { transform: scale(0.92); }
.done-btn:disabled { background: #a3e635; cursor: not-allowed; }
.check-icon { font-size: 24px; }
.check-icon--done { color: #22c55e; }
</style>
