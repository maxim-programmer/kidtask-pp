<template>
  <JuniorLayout>
    <div class="page-header">
      <h1 class="page-title">Мои задания</h1>
      <button class="filter-btn" @click="showFilter = !showFilter">≡ Фильтр</button>
    </div>

    <div v-if="showFilter" class="filter-bar">
      <button v-for="f in filters" :key="f.value"
        :class="['filter-chip', { 'filter-chip--active': activeFilter === f.value }]"
        @click="activeFilter = f.value">{{ f.label }}</button>
    </div>

    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else>
      <div v-if="filteredTasks.length === 0" class="empty">Нет заданий 🎉</div>
      <div v-for="t in filteredTasks" :key="t.task_id" :class="['task-card', `task-card--${t.status}`]">
        <div class="task-card__top">
          <div>
            <div class="task-card__title">{{ t.title }}</div>
            <div class="task-card__desc" v-if="t.description">{{ t.description }}</div>
            <div class="task-card__rework" v-if="t.rejection_comment">
              💬 {{ t.rejection_comment }}
            </div>
          </div>
          <button class="audio-btn" @click="speak(t.title)" title="Прочитать">🔊</button>
        </div>
        <div class="task-card__footer">
          <span class="reward">⭐ {{ t.reward }}</span>
          <div class="task-actions">
            <button
              v-if="t.status === 'active' || t.status === 'needs_rework'"
              class="done-btn"
              @click="submit(t)"
              :disabled="submitting === t.task_id">
              {{ submitting === t.task_id ? '...' : '✓ Сделал!' }}
            </button>
            <span v-else-if="t.status === 'pending_review'" class="status-badge status-badge--pending">На проверке ⏳</span>
            <span v-else-if="t.status === 'completed'" class="status-badge status-badge--done">Выполнено ✓</span>
          </div>
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
      loading: true, tasks: [], submitting: null,
      showFilter: false, activeFilter: '',
      filters: [
        { value: '', label: 'Все' },
        { value: 'active', label: 'Активные' },
        { value: 'pending_review', label: 'На проверке' },
        { value: 'needs_rework', label: 'На доработке' },
        { value: 'completed', label: 'Выполненные' },
      ]
    }
  },
  computed: {
    filteredTasks() {
      if (!this.activeFilter) return this.tasks
      return this.tasks.filter(t => t.status === this.activeFilter)
    }
  },
  async mounted() {
    await this.load()
  },
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
      catch (e) { alert('Ошибка: ' + (e.response?.data?.error?.message || 'Попробуй снова')) }
      finally { this.submitting = null }
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
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.page-title { font-size: 26px; font-weight: 800; color: #ea580c; }
.filter-btn { padding: 8px 14px; background: #fff; border: 2px solid #ea580c; border-radius: 10px; color: #ea580c; font-weight: 700; font-size: 14px; cursor: pointer; }
.filter-bar { display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 14px; }
.filter-chip { padding: 6px 14px; border-radius: 20px; border: 2px solid #fed7aa; background: #fff; color: #ea580c; font-size: 13px; font-weight: 600; cursor: pointer; transition: all 0.15s; }
.filter-chip--active { background: #ea580c; color: #fff; border-color: #ea580c; }
.loading { text-align: center; padding: 60px; color: #ea580c; font-size: 18px; }
.empty { text-align: center; color: #aaa; font-size: 18px; padding: 40px; }
.task-card { background: #fff; border-radius: 20px; padding: 16px; margin-bottom: 12px; box-shadow: 0 2px 8px rgba(234,88,12,0.08); border-left: 5px solid #fed7aa; }
.task-card--active { border-left-color: #ea580c; }
.task-card--pending_review { border-left-color: #f59e0b; }
.task-card--needs_rework { border-left-color: #e53e3e; }
.task-card--completed { border-left-color: #22c55e; }
.task-card__top { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 10px; }
.task-card__title { font-size: 18px; font-weight: 800; color: #1a1a1a; margin-bottom: 4px; }
.task-card__desc { font-size: 14px; color: #888; }
.task-card__rework { font-size: 13px; color: #e53e3e; margin-top: 6px; background: #fff5f5; border-radius: 8px; padding: 6px 10px; }
.audio-btn { background: #fff5eb; border: none; border-radius: 10px; width: 36px; height: 36px; font-size: 18px; cursor: pointer; flex-shrink: 0; }
.task-card__footer { display: flex; align-items: center; justify-content: space-between; }
.reward { font-size: 20px; font-weight: 800; color: #ea580c; }
.task-actions { display: flex; gap: 8px; }
.done-btn { padding: 10px 20px; background: #22c55e; color: #fff; border: none; border-radius: 12px; font-size: 16px; font-weight: 700; cursor: pointer; transition: background 0.2s; }
.done-btn:hover:not(:disabled) { background: #16a34a; }
.done-btn:disabled { opacity: 0.6; cursor: not-allowed; }
.status-badge { padding: 8px 14px; border-radius: 12px; font-size: 14px; font-weight: 600; }
.status-badge--pending { background: #fef3c7; color: #92400e; }
.status-badge--done { background: #dcfce7; color: #166534; }
</style>