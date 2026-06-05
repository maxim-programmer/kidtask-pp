<template>
  <ParentLayout>
    <div class="page-header">
      <h1 class="page-title">Задания</h1>
      <button class="btn-add" @click="showAddTask = true">+ Новое</button>
    </div>

    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else>
      <div v-if="!children.length" class="empty-state">
        <div class="empty-state__icon">👶</div>
        <p>Сначала добавьте ребёнка</p>
        <router-link to="/parent/children" class="btn-primary">Мои дети</router-link>
      </div>

      <div v-else>
        <div class="child-tabs">
          <button v-for="c in children" :key="c.child_id"
            :class="['tab', { 'tab--active': selectedChildId === c.child_id }]"
            @click="selectChild(c.child_id)">
            <span v-if="c.avatar_url" class="tab-avatar-wrap"><img :src="c.avatar_url" class="tab-avatar" /></span>
            {{ c.name }}
          </button>
        </div>

        <div class="filter-bar">
          <button v-for="f in taskFilters" :key="f.value"
            :class="['chip', { 'chip--active': taskFilter === f.value }]"
            @click="taskFilter = f.value">{{ f.label }}</button>
        </div>

        <div v-if="filteredTasks.length === 0" class="empty">Нет заданий</div>

        <div v-for="t in filteredTasks" :key="t.task_id" :class="['task-card', `task-card--${t.status}`]">
          <div class="task-card__body">
            <div class="task-card__top">
              <span class="task-status-dot"></span>
              <span class="task-title">{{ t.title }}</span>
            </div>
            <div class="task-desc" v-if="t.description">{{ t.description }}</div>
            <div class="task-rework" v-if="t.rejection_comment">💬 {{ t.rejection_comment }}</div>
            <div class="task-meta">
              <span class="reward">⭐ {{ t.reward }}</span>
              <span class="status-label">{{ statusLabel(t.status) }}</span>
              <span class="child-name-tag" v-if="t.child_name">{{ t.child_name }}</span>
            </div>
          </div>
          <div class="task-card__actions">
            <button v-if="t.status === 'pending_review'" class="act-btn act-btn--approve" @click="approve(t)">✓</button>
            <button v-if="t.status === 'pending_review'" class="act-btn act-btn--reject" @click="openReject(t)">✗</button>
            <button v-if="t.status === 'active'" class="act-btn act-btn--edit" @click="openEdit(t)">✏</button>
            <button class="act-btn act-btn--delete" @click="deleteTask(t)">🗑</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showAddTask" class="modal-overlay" @click.self="showAddTask = false">
      <div class="modal">
        <div class="modal__handle"></div>
        <h3>Новое задание</h3>
        <div class="field">
          <label>Для кого</label>
          <select v-model="addForm.child_id">
            <option value="" disabled>Выберите ребёнка</option>
            <option v-for="c in children" :key="c.child_id" :value="c.child_id">{{ c.name }}</option>
          </select>
        </div>
        <div class="field"><label>Название</label><input v-model="addForm.title" type="text" placeholder="Убрать комнату" /></div>
        <div class="field"><label>Награда (⭐)</label><input v-model.number="addForm.reward" type="number" min="1" placeholder="10" /></div>
        <div class="field"><label>Описание</label><textarea v-model="addForm.description" rows="2" placeholder="Необязательно..."></textarea></div>
        <div v-if="addError" class="error-msg">{{ addError }}</div>
        <div class="modal-actions">
          <button class="btn-outline" @click="showAddTask = false">Отмена</button>
          <button class="btn-primary" @click="createTask" :disabled="saving">{{ saving ? '...' : 'Добавить' }}</button>
        </div>
      </div>
    </div>

    <div v-if="rejectModal" class="modal-overlay" @click.self="rejectModal = false">
      <div class="modal">
        <div class="modal__handle"></div>
        <h3>Вернуть на доработку</h3>
        <div class="field"><label>Комментарий</label><textarea v-model="rejectComment" rows="3" placeholder="Объясните что нужно исправить..."></textarea></div>
        <div class="modal-actions">
          <button class="btn-outline" @click="rejectModal = false">Отмена</button>
          <button class="btn-primary" @click="confirmReject">Вернуть</button>
        </div>
      </div>
    </div>

    <div v-if="editModal" class="modal-overlay" @click.self="editModal = false">
      <div class="modal">
        <div class="modal__handle"></div>
        <h3>Редактировать задание</h3>
        <div class="field"><label>Название</label><input v-model="editForm.title" type="text" /></div>
        <div class="field"><label>Награда (⭐)</label><input v-model.number="editForm.reward" type="number" min="1" /></div>
        <div class="field"><label>Описание</label><textarea v-model="editForm.description" rows="2"></textarea></div>
        <div class="modal-actions">
          <button class="btn-outline" @click="editModal = false">Отмена</button>
          <button class="btn-primary" @click="saveEdit">Сохранить</button>
        </div>
      </div>
    </div>
  </ParentLayout>
</template>

<script>
import ParentLayout from '../../components/ParentLayout.vue'
import { useApi } from '../../composables/useApi'
export default {
  name: 'ParentTasks',
  components: { ParentLayout },
  data() {
    return {
      loading: true, children: [], tasks: [],
      selectedChildId: '',
      taskFilter: '',
      taskFilters: [
        { value: '', label: 'Все' },
        { value: 'active', label: 'Активные' },
        { value: 'pending_review', label: 'На проверке' },
        { value: 'needs_rework', label: 'Доработка' },
        { value: 'completed', label: 'Выполненные' },
      ],
      showAddTask: false, saving: false, addError: '',
      addForm: { child_id: '', title: '', reward: '', description: '' },
      rejectModal: false, rejectComment: '', rejectTarget: null,
      editModal: false, editTarget: null, editForm: { title: '', reward: '', description: '' },
    }
  },
  computed: {
    filteredTasks() {
      let t = this.tasks
      if (this.selectedChildId) t = t.filter(x => x.child_id === this.selectedChildId)
      if (this.taskFilter) t = t.filter(x => x.status === this.taskFilter)
      return t
    }
  },
  async mounted() { await this.load() },
  methods: {
    async load() {
      this.loading = true
      const { getChildren, getTasks } = useApi()
      try {
        const [cr, tr] = await Promise.all([getChildren(), getTasks()])
        this.children = cr.data.children || []
        this.tasks = tr.data.tasks || []
        if (this.children.length && !this.selectedChildId) this.selectedChildId = this.children[0].child_id
      } finally { this.loading = false }
    },
    selectChild(id) { this.selectedChildId = id; this.taskFilter = '' },
    statusLabel(s) { return { active: 'Активное', pending_review: 'На проверке', needs_rework: 'Доработка', completed: 'Выполнено' }[s] || s },
    async approve(task) { const { approveTask } = useApi(); await approveTask(task.task_id); await this.load() },
    openReject(task) { this.rejectTarget = task; this.rejectComment = ''; this.rejectModal = true },
    async confirmReject() { const { rejectTask } = useApi(); await rejectTask(this.rejectTarget.task_id, this.rejectComment); this.rejectModal = false; await this.load() },
    async deleteTask(task) { if (!confirm('Удалить задание?')) return; const { deleteTask } = useApi(); await deleteTask(task.task_id); await this.load() },
    openEdit(t) { this.editTarget = t; this.editForm = { title: t.title, reward: t.reward, description: t.description || '' }; this.editModal = true },
    async saveEdit() { const { updateTask } = useApi(); await updateTask(this.editTarget.task_id, { title: this.editForm.title, reward: this.editForm.reward, description: this.editForm.description || undefined }); this.editModal = false; await this.load() },
    async createTask() {
      if (!this.addForm.child_id || !this.addForm.title || !this.addForm.reward) { this.addError = 'Заполните обязательные поля'; return }
      this.saving = true; this.addError = ''
      const { createTask } = useApi()
      try { await createTask({ title: this.addForm.title, child_id: this.addForm.child_id, reward: this.addForm.reward, description: this.addForm.description || undefined }); this.showAddTask = false; this.addForm = { child_id: this.addForm.child_id, title: '', reward: '', description: '' }; await this.load() }
      catch (e) { this.addError = e.response?.data?.error?.message || 'Ошибка' }
      finally { this.saving = false }
    }
  }
}
</script>

<style scoped>
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.page-title { font-size: 24px; font-weight: 700; }
.btn-add { padding: 10px 18px; background: #4f7ef7; color: #fff; border: none; border-radius: 10px; font-size: 15px; font-weight: 600; cursor: pointer; font-family: inherit; }
.loading { text-align: center; padding: 60px; color: #888; }
.empty-state { text-align: center; padding: 60px 20px; }
.empty-state__icon { font-size: 48px; margin-bottom: 12px; }
.empty-state p { font-size: 16px; color: #888; margin-bottom: 16px; }
.btn-primary { display: inline-block; padding: 12px 24px; background: #4f7ef7; color: #fff; border: none; border-radius: 10px; font-size: 15px; font-weight: 600; cursor: pointer; font-family: inherit; }
.child-tabs { display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 14px; }
.tab { padding: 7px 16px; border-radius: 20px; border: 2px solid #4f7ef7; background: transparent; color: #4f7ef7; font-weight: 600; font-size: 13px; cursor: pointer; transition: all 0.15s; font-family: inherit; display: flex; align-items: center; gap: 6px; }
.tab--active { background: #4f7ef7; color: #fff; }
.tab-avatar-wrap { width: 18px; height: 18px; border-radius: 50%; overflow: hidden; }
.tab-avatar { width: 100%; height: 100%; object-fit: cover; }
.filter-bar { display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 14px; }
.chip { padding: 5px 11px; border-radius: 16px; border: 1.5px solid #d0d8ff; background: #fff; color: #4f7ef7; font-size: 12px; font-weight: 500; cursor: pointer; font-family: inherit; }
.chip--active { background: #4f7ef7; color: #fff; border-color: #4f7ef7; }
.empty { text-align: center; color: #bbb; font-size: 15px; padding: 40px; }
.task-card { background: #fff; border-radius: 14px; padding: 14px 16px; margin-bottom: 10px; display: flex; align-items: flex-start; justify-content: space-between; gap: 10px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); border-left: 4px solid #e0e0e0; }
.task-card--active { border-left-color: #4f7ef7; }
.task-card--pending_review { border-left-color: #f59e0b; }
.task-card--needs_rework { border-left-color: #e53e3e; }
.task-card--completed { border-left-color: #22c55e; opacity: 0.75; }
.task-card__body { flex: 1; min-width: 0; }
.task-card__top { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.task-status-dot { width: 8px; height: 8px; border-radius: 50%; background: currentColor; flex-shrink: 0; }
.task-card--active .task-status-dot { color: #4f7ef7; }
.task-card--pending_review .task-status-dot { color: #f59e0b; }
.task-card--needs_rework .task-status-dot { color: #e53e3e; }
.task-card--completed .task-status-dot { color: #22c55e; }
.task-title { font-size: 15px; font-weight: 700; color: #1a1a1a; }
.task-desc { font-size: 13px; color: #888; margin-bottom: 6px; }
.task-rework { font-size: 12px; color: #e53e3e; margin-bottom: 6px; }
.task-meta { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.reward { font-size: 14px; font-weight: 700; color: #f59e0b; }
.status-label { font-size: 11px; color: #aaa; }
.child-name-tag { font-size: 12px; background: #f0f4ff; color: #4f7ef7; padding: 2px 8px; border-radius: 8px; font-weight: 600; }
.task-card__actions { display: flex; gap: 5px; flex-direction: column; }
.act-btn { width: 32px; height: 32px; border-radius: 8px; border: none; font-size: 14px; cursor: pointer; display: flex; align-items: center; justify-content: center; font-family: inherit; }
.act-btn--approve { background: #dcfce7; color: #16a34a; }
.act-btn--reject { background: #fee2e2; color: #dc2626; }
.act-btn--edit { background: #e0f2fe; color: #0284c7; }
.act-btn--delete { background: #f5f5f5; color: #888; }
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: flex-end; justify-content: center; z-index: 300; }
.modal { background: #fff; border-radius: 20px 20px 0 0; padding: 8px 20px 28px; width: 100%; max-width: 600px; max-height: 92vh; overflow-y: auto; }
.modal__handle { width: 40px; height: 4px; background: #ddd; border-radius: 2px; margin: 0 auto 16px; }
.modal h3 { font-size: 20px; font-weight: 700; text-align: center; margin-bottom: 18px; }
.field { margin-bottom: 14px; }
.field label { display: block; font-size: 13px; color: #666; margin-bottom: 6px; }
.field input, .field textarea, .field select { width: 100%; padding: 12px 14px; border: 1.5px solid #e0e7ff; border-radius: 12px; font-size: 15px; outline: none; box-sizing: border-box; font-family: inherit; resize: vertical; background: #fafbff; }
.field input:focus, .field textarea:focus, .field select:focus { border-color: #4f7ef7; }
.error-msg { color: #e53e3e; font-size: 13px; margin-bottom: 10px; }
.modal-actions { display: flex; gap: 10px; margin-top: 8px; }
.btn-outline { flex: 1; padding: 13px; border: 1.5px solid #ddd; background: #fff; border-radius: 12px; font-size: 15px; font-weight: 600; cursor: pointer; color: #666; font-family: inherit; }
.btn-primary { flex: 1; padding: 13px; background: #4f7ef7; color: #fff; border: none; border-radius: 12px; font-size: 15px; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
@media (min-width: 600px) {
  .modal-overlay { align-items: center; }
  .modal { border-radius: 20px; max-width: 420px; }
  .modal__handle { display: none; }
  .task-card__actions { flex-direction: row; }
}
</style>
