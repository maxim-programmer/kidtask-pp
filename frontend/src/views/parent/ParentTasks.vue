<template>
  <ParentLayout>
    <div class="header-row">
      <h1 class="page-title">Список заданий</h1>
      <button class="btn-icon" @click="showFilter = !showFilter">≡ Фильтр</button>
    </div>

    <div v-if="showFilter" class="filter-bar">
      <select v-model="filterChild" class="filter-select">
        <option value="">Все дети</option>
        <option v-for="c in children" :key="c.child_id" :value="c.child_id">{{ c.name }}</option>
      </select>
      <select v-model="filterStatus" class="filter-select">
        <option value="">Все статусы</option>
        <option value="active">Активные</option>
        <option value="pending_review">На проверке</option>
        <option value="needs_rework">На доработке</option>
        <option value="completed">Выполненные</option>
      </select>
    </div>

    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else>
      <div v-for="t in filteredTasks" :key="t.task_id" :class="['task-card', `task-card--${statusClass(t.status)}`]">
        <div class="task-card__header">
          <span class="task-card__name">{{ t.title }} ({{ t.child_name }})</span>
          <span class="task-card__status">{{ statusLabel(t.status) }}</span>
        </div>
        <div class="task-card__desc" v-if="t.description">{{ t.description }}</div>
        <div class="task-card__footer">
          <span class="task-card__reward">⭐ {{ t.reward }}</span>
          <div class="task-card__actions">
            <button class="icon-btn" @click="openEdit(t)" v-if="t.status === 'active'">✏</button>
            <button class="icon-btn icon-btn--approve" @click="approve(t)" v-if="t.status === 'pending_review'" title="Одобрить">✓</button>
            <button class="icon-btn icon-btn--reject" @click="openReject(t)" v-if="t.status === 'pending_review'" title="Вернуть">✗</button>
            <button class="icon-btn icon-btn--delete" @click="remove(t)" title="Удалить">🗑</button>
          </div>
        </div>
      </div>

      <div v-if="filteredTasks.length === 0" class="empty">Нет заданий</div>

      <button class="add-btn" @click="showAddModal = true">+ Добавить задание</button>
    </div>

    <div v-if="showAddModal" class="modal-overlay" @click.self="showAddModal = false">
      <div class="modal">
        <h3>Новое задание</h3>
        <div class="field">
          <label>Название</label>
          <input v-model="addForm.title" type="text" placeholder="Убрать комнату" />
        </div>
        <div class="field">
          <label>Исполнитель</label>
          <select v-model="addForm.child_id" class="form-select">
            <option v-for="c in children" :key="c.child_id" :value="c.child_id">{{ c.name }}</option>
          </select>
        </div>
        <div class="field">
          <label>Награда (⭐)</label>
          <input v-model.number="addForm.reward" type="number" min="1" placeholder="10" />
        </div>
        <div class="field">
          <label>Описание</label>
          <textarea v-model="addForm.description" rows="2" class="form-textarea" placeholder="Дополнительно..."></textarea>
        </div>
        <div v-if="addError" class="error-msg">{{ addError }}</div>
        <button class="btn-primary" @click="createTask" :disabled="saving">
          {{ saving ? 'Добавление...' : 'Добавить' }}
        </button>
      </div>
    </div>

    <div v-if="editModal" class="modal-overlay" @click.self="editModal = false">
      <div class="modal">
        <h3>Редактировать задание</h3>
        <div class="field">
          <label>Название</label>
          <input v-model="editForm.title" type="text" />
        </div>
        <div class="field">
          <label>Награда (⭐)</label>
          <input v-model.number="editForm.reward" type="number" min="1" />
        </div>
        <div class="field">
          <label>Описание</label>
          <textarea v-model="editForm.description" rows="2" class="form-textarea"></textarea>
        </div>
        <div class="modal-actions">
          <button class="btn-outline" @click="editModal = false">Отмена</button>
          <button class="btn-primary" @click="saveEdit">Сохранить</button>
        </div>
      </div>
    </div>

    <div v-if="rejectModal" class="modal-overlay" @click.self="rejectModal = false">
      <div class="modal">
        <h3>Вернуть на доработку</h3>
        <textarea v-model="rejectComment" rows="3" class="form-textarea" placeholder="Комментарий..."></textarea>
        <div class="modal-actions">
          <button class="btn-outline" @click="rejectModal = false">Отмена</button>
          <button class="btn-primary" @click="confirmReject">Вернуть</button>
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
      loading: true, tasks: [], children: [],
      showFilter: false, filterChild: '', filterStatus: '',
      showAddModal: false, saving: false, addError: '',
      addForm: { title: '', child_id: '', reward: '', description: '' },
      editModal: false, editTarget: null, editForm: { title: '', reward: '', description: '' },
      rejectModal: false, rejectTarget: null, rejectComment: ''
    }
  },
  computed: {
    filteredTasks() {
      return this.tasks.filter(t => {
        if (this.filterChild && t.child_id !== this.filterChild) return false
        if (this.filterStatus && t.status !== this.filterStatus) return false
        return true
      })
    }
  },
  async mounted() {
    const { getTasks, getChildren } = useApi()
    try {
      const [tr, cr] = await Promise.all([getTasks(), getChildren()])
      this.tasks = tr.data.tasks || []
      this.children = cr.data.children || []
      if (this.children.length) this.addForm.child_id = this.children[0].child_id
    } finally { this.loading = false }
  },
  methods: {
    statusLabel(s) {
      return { active: 'Активное', pending_review: 'На проверке', needs_rework: 'На доработке', completed: 'Выполнено' }[s] || s
    },
    statusClass(s) {
      return { active: 'active', pending_review: 'pending', needs_rework: 'rework', completed: 'done' }[s] || ''
    },
    async createTask() {
      this.saving = true; this.addError = ''
      const { createTask } = useApi()
      try {
        await createTask({
          title: this.addForm.title, child_id: this.addForm.child_id,
          reward: this.addForm.reward, description: this.addForm.description || undefined
        })
        this.showAddModal = false
        this.addForm = { title: '', child_id: this.children[0]?.child_id || '', reward: '', description: '' }
        const { getTasks } = useApi()
        this.tasks = (await getTasks()).data.tasks || []
      } catch (e) { this.addError = e.response?.data?.error?.message || 'Ошибка' }
      finally { this.saving = false }
    },
    openEdit(t) { this.editTarget = t; this.editForm = { title: t.title, reward: t.reward, description: t.description || '' }; this.editModal = true },
    async saveEdit() {
      const { updateTask, getTasks } = useApi()
      await updateTask(this.editTarget.task_id, { title: this.editForm.title, reward: this.editForm.reward, description: this.editForm.description || undefined })
      this.editModal = false
      this.tasks = (await getTasks()).data.tasks || []
    },
    async approve(t) {
      const { approveTask, getTasks } = useApi()
      await approveTask(t.task_id)
      this.tasks = (await getTasks()).data.tasks || []
    },
    openReject(t) { this.rejectTarget = t; this.rejectComment = ''; this.rejectModal = true },
    async confirmReject() {
      const { rejectTask, getTasks } = useApi()
      await rejectTask(this.rejectTarget.task_id, this.rejectComment)
      this.rejectModal = false
      this.tasks = (await getTasks()).data.tasks || []
    },
    async remove(t) {
      if (!confirm('Удалить задание?')) return
      const { deleteTask, getTasks } = useApi()
      await deleteTask(t.task_id)
      this.tasks = (await getTasks()).data.tasks || []
    }
  }
}
</script>

<style scoped>
.header-row { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.page-title { font-size: 24px; font-weight: 700; }
.btn-icon { padding: 8px 14px; background: #f0f4ff; border: none; border-radius: 8px; color: #4f7ef7; font-weight: 600; font-size: 14px; cursor: pointer; }
.filter-bar { display: flex; gap: 10px; margin-bottom: 16px; }
.filter-select { flex: 1; padding: 8px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 14px; outline: none; background: #fff; }
.loading { text-align: center; padding: 60px; color: #888; }
.task-card { background: #fff; border-radius: 16px; padding: 16px; margin-bottom: 10px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); border-left: 4px solid #e8e8e8; }
.task-card--active { border-left-color: #4f7ef7; }
.task-card--pending { border-left-color: #f59e0b; }
.task-card--rework { border-left-color: #e53e3e; }
.task-card--done { border-left-color: #22c55e; }
.task-card__header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 4px; }
.task-card__name { font-size: 15px; font-weight: 600; }
.task-card__status { font-size: 12px; color: #888; background: #f5f5f5; padding: 2px 8px; border-radius: 10px; white-space: nowrap; margin-left: 8px; }
.task-card__desc { font-size: 13px; color: #888; margin-bottom: 8px; }
.task-card__footer { display: flex; align-items: center; justify-content: space-between; }
.task-card__reward { font-size: 14px; font-weight: 700; color: #4f7ef7; }
.task-card__actions { display: flex; gap: 6px; }
.icon-btn { width: 32px; height: 32px; border-radius: 8px; border: none; font-size: 15px; cursor: pointer; background: #f5f5f5; color: #666; }
.icon-btn--approve { background: #e6f9f0; color: #22c55e; }
.icon-btn--approve:hover { background: #22c55e; color: #fff; }
.icon-btn--reject { background: #fff5f5; color: #e53e3e; }
.icon-btn--reject:hover { background: #e53e3e; color: #fff; }
.icon-btn--delete { background: #f5f5f5; color: #999; }
.icon-btn--delete:hover { background: #e53e3e; color: #fff; }
.empty { text-align: center; color: #aaa; padding: 40px; }
.add-btn { width: 100%; padding: 14px; background: #f0f4ff; border: 2px dashed #4f7ef7; border-radius: 16px; color: #4f7ef7; font-size: 16px; font-weight: 600; cursor: pointer; margin-top: 8px; }
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 300; }
.modal { background: #fff; border-radius: 16px; padding: 24px; width: 90%; max-width: 380px; }
.modal h3 { font-size: 20px; font-weight: 700; margin-bottom: 16px; }
.field { margin-bottom: 12px; }
.field label { display: block; font-size: 13px; color: #666; margin-bottom: 4px; }
.field input, .form-select, .form-textarea { width: 100%; padding: 10px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 15px; outline: none; background: #fff; }
.field input:focus, .form-select:focus, .form-textarea:focus { border-color: #4f7ef7; }
.form-textarea { resize: vertical; }
.error-msg { color: #e53e3e; font-size: 13px; margin-bottom: 10px; }
.btn-primary { width: 100%; padding: 12px; background: #4f7ef7; color: #fff; border: none; border-radius: 8px; font-size: 16px; font-weight: 600; cursor: pointer; margin-top: 8px; }
.btn-primary:hover:not(:disabled) { background: #3a6be0; }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-outline { padding: 10px 20px; background: transparent; border: 2px solid #ddd; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; }
.modal-actions { display: flex; gap: 8px; margin-top: 16px; justify-content: flex-end; }
</style>