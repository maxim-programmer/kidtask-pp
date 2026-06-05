<template>
  <SeniorLayout>
    <h1 class="page-title">Мои задания</h1>

    <div class="filter-bar">
      <button v-for="f in filters" :key="f.value"
        :class="['chip', { 'chip--active': active === f.value }]"
        @click="active = f.value">{{ f.label }}</button>
    </div>

    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else>
      <div v-if="filtered.length === 0" class="empty">Нет заданий</div>

      <div v-for="t in filtered" :key="t.task_id" :class="['task-card', `task-card--${t.status}`]">
        <div class="task-card__body">
          <div class="task-card__title">{{ t.title }}</div>
          <div class="task-card__desc" v-if="t.description">{{ t.description }}</div>
          <div class="task-card__rework" v-if="t.rejection_comment">💬 {{ t.rejection_comment }}</div>
          <div class="task-card__meta">
            <span class="reward">⭐ {{ t.reward }}</span>
            <span :class="['status-tag', `status-tag--${t.status}`]">{{ statusLabel(t.status) }}</span>
          </div>
        </div>
        <div class="task-card__action">
          <button
            v-if="t.status === 'active' || t.status === 'needs_rework'"
            class="done-btn"
            @click="submit(t)"
            :disabled="submitting === t.task_id">
            {{ submitting === t.task_id ? '...' : 'Сдать' }}
          </button>
          <span v-else-if="t.status === 'pending_review'" class="status-icon">⏳</span>
          <span v-else-if="t.status === 'completed'" class="status-icon status-icon--done">✓</span>
        </div>
      </div>
    </div>
  </SeniorLayout>
</template>

<script>
import SeniorLayout from '../../../components/SeniorLayout.vue'
import { useApi } from '../../../composables/useApi'
export default {
  name: 'ChildTasksSenior',
  components: { SeniorLayout },
  data() {
    return {
      loading: true, tasks: [], submitting: null, active: '',
      filters: [
        { value: '', label: 'Все' },
        { value: 'active', label: 'Активные' },
        { value: 'pending_review', label: 'На проверке' },
        { value: 'needs_rework', label: 'Доработка' },
        { value: 'completed', label: 'Выполненные' },
      ]
    }
  },
  computed: {
    filtered() { return this.active ? this.tasks.filter(t => t.status === this.active) : this.tasks }
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
      catch (e) { alert(e.response?.data?.error?.message || 'Ошибка') }
      finally { this.submitting = null }
    },
    statusLabel(s) {
      return { active: 'Активное', pending_review: 'На проверке', needs_rework: 'Доработка', completed: 'Выполнено' }[s] || s
    }
  }
}
</script>

<style scoped>
.page-title { font-size: 22px; font-weight: 700; color: #f3f4f6; margin-bottom: 16px; }
.filter-bar { display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 16px; overflow-x: auto; -webkit-overflow-scrolling: touch; padding-bottom: 4px; }
.chip { padding: 6px 13px; border-radius: 16px; border: 1px solid #374151; background: transparent; color: #9ca3af; font-size: 13px; font-weight: 500; cursor: pointer; font-family: inherit; white-space: nowrap; }
.chip--active { background: #6366f1; color: #fff; border-color: #6366f1; }
.loading { text-align: center; padding: 60px; color: #6366f1; }
.empty { text-align: center; color: #6b7280; font-size: 15px; padding: 40px; }
.task-card { display: flex; align-items: flex-start; gap: 12px; background: #1f2937; border-radius: 14px; padding: 16px; margin-bottom: 10px; border: 1px solid #374151; border-left: 4px solid #374151; }
.task-card--active { border-left-color: #6366f1; }
.task-card--pending_review { border-left-color: #f59e0b; }
.task-card--needs_rework { border-left-color: #ef4444; }
.task-card--completed { border-left-color: #22c55e; opacity: 0.7; }
.task-card__body { flex: 1; min-width: 0; }
.task-card__title { font-size: 15px; font-weight: 700; color: #f3f4f6; margin-bottom: 4px; }
.task-card__desc { font-size: 13px; color: #9ca3af; margin-bottom: 6px; }
.task-card__rework { font-size: 12px; color: #f87171; background: rgba(239,68,68,0.1); border-radius: 8px; padding: 6px 10px; margin-bottom: 6px; }
.task-card__meta { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.reward { font-size: 15px; font-weight: 700; color: #a5b4fc; }
.status-tag { font-size: 11px; font-weight: 600; padding: 3px 8px; border-radius: 8px; }
.status-tag--active { background: #1e1b4b; color: #818cf8; }
.status-tag--pending_review { background: #78350f30; color: #fbbf24; }
.status-tag--needs_rework { background: rgba(239,68,68,0.15); color: #f87171; }
.status-tag--completed { background: #14532d30; color: #4ade80; }
.task-card__action { flex-shrink: 0; display: flex; align-items: center; }
.done-btn { padding: 10px 16px; background: #6366f1; color: #fff; border: none; border-radius: 10px; font-size: 14px; font-weight: 600; cursor: pointer; font-family: inherit; }
.done-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.status-icon { font-size: 20px; color: #f59e0b; }
.status-icon--done { color: #4ade80; }
</style>
