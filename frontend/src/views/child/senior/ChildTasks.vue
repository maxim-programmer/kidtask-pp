<template>
  <SeniorLayout>
    <div class="page-header">
      <h1 class="page-title">Мои задания</h1>
      <button class="filter-btn" @click="showFilter = !showFilter">≡ Фильтр</button>
    </div>

    <div v-if="showFilter" class="filter-bar">
      <button v-for="f in filters" :key="f.value"
        :class="['chip', { 'chip--active': activeFilter === f.value }]"
        @click="activeFilter = f.value">{{ f.label }}</button>
    </div>

    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else>
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
          <button class="icon-del" @click="remove(t)" v-if="false">🗑</button>
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
      loading: true, tasks: [], submitting: null,
      showFilter: false, activeFilter: '',
      filters: [
        { value: '', label: 'Все' },
        { value: 'active', label: 'Активные' },
        { value: 'pending_review', label: 'На проверке' },
        { value: 'needs_rework', label: 'На доработке' },
        { value: 'completed', label: 'Завершённые' },
      ]
    }
  },
  computed: {
    filteredTasks() {
      if (!this.activeFilter) return this.tasks
      return this.tasks.filter(t => t.status === this.activeFilter)
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
      catch (e) { alert(e.response?.data?.error?.message || 'Ошибка') }
      finally { this.submitting = null }
    },
    async remove(task) {
      if (!confirm('Удалить?')) return
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
.page-title { font-size: 22px; font-weight: 800; color: #e5e7eb; }
.filter-btn { padding: 8px 14px; background: #1f2937; border: 1px solid #374151; border-radius: 8px; color: #9ca3af; font-size: 13px; font-weight: 600; cursor: pointer; }
.filter-bar { display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 14px; }
.chip { padding: 5px 12px; border-radius: 16px; border: 1px solid #374151; background: transparent; color: #9ca3af; font-size: 13px; font-weight: 500; cursor: pointer; transition: all 0.15s; }
.chip--active { background: #6366f1; color: #fff; border-color: #6366f1; }
.loading { text-align: center; padding: 60px; color: #6366f1; }
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
.done-btn { padding: 8px 18px; background: #22c55e; color: #fff; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; }
.done-btn:hover:not(:disabled) { background: #16a34a; }
.done-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.badge { padding: 6px 12px; border-radius: 8px; font-size: 13px; font-weight: 600; }
.badge--pending { background: #78350f30; color: #fbbf24; }
.badge--done { background: #14532d30; color: #4ade80; }
</style>