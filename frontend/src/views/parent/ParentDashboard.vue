<template>
  <ParentLayout>
    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else>
      <h1 class="greeting">Добрый день, {{ user?.name }}!</h1>

      <div v-if="!children.length" class="no-child">
        <p>Добавьте ребёнка, чтобы начать</p>
        <button class="btn-primary" @click="$router.push('/parent/children')">Мои дети</button>
      </div>

      <div v-else>
        <div class="child-tabs">
          <button
            v-for="c in children" :key="c.child_id"
            :class="['tab', { 'tab--active': selectedChild?.child_id === c.child_id }]"
            @click="selectChild(c)"
          >{{ c.name }}</button>
        </div>

        <div v-if="selectedChild">
          <div class="balance-card">
            <div>
              <div class="balance-label">Баланс {{ selectedChild.name }}</div>
              <div class="balance-meta">{{ selectedChild.age_group === 'junior' ? '7–10 лет' : '11–14 лет' }}</div>
            </div>
            <div class="balance-value">⭐ {{ selectedChild.balance }}</div>
          </div>

          <div class="progress-row">
            <div class="progress-item">
              <div class="progress-circle">
                <svg viewBox="0 0 36 36">
                  <circle cx="18" cy="18" r="15" fill="none" stroke="#e8e8e8" stroke-width="3"/>
                  <circle cx="18" cy="18" r="15" fill="none" stroke="#4f7ef7" stroke-width="3"
                    :stroke-dasharray="`${wishProgress} 100`" stroke-linecap="round" transform="rotate(-90 18 18)"/>
                </svg>
                <span class="progress-text">{{ progress.wishes_done }}/{{ progress.wishes_total }}</span>
              </div>
              <span class="progress-label">Целей</span>
            </div>
            <div class="progress-item">
              <div class="progress-circle">
                <svg viewBox="0 0 36 36">
                  <circle cx="18" cy="18" r="15" fill="none" stroke="#e8e8e8" stroke-width="3"/>
                  <circle cx="18" cy="18" r="15" fill="none" stroke="#4f7ef7" stroke-width="3"
                    :stroke-dasharray="`${taskProgress} 100`" stroke-linecap="round" transform="rotate(-90 18 18)"/>
                </svg>
                <span class="progress-text">{{ progress.tasks_done }}/{{ progress.tasks_total }}</span>
              </div>
              <span class="progress-label">Заданий</span>
            </div>
          </div>

          <div class="view-tabs">
            <button :class="['vtab', { 'vtab--active': view === 'overview' }]" @click="view = 'overview'">Обзор</button>
            <button :class="['vtab', { 'vtab--active': view === 'tasks' }]" @click="view = 'tasks'">Задания</button>
            <button :class="['vtab', { 'vtab--active': view === 'wishlist' }]" @click="view = 'wishlist'">Вишлист</button>
          </div>

          <!-- ОБЗОР -->
          <div v-if="view === 'overview'">
            <div class="section">
              <h3 class="section-title">Задания на проверке</h3>
              <div v-if="pendingTasks.length === 0" class="empty">Нет заданий на проверке</div>
              <div v-for="t in pendingTasks" :key="t.task_id" class="task-card">
                <div class="task-card__info">
                  <div class="task-card__title">{{ t.title }}</div>
                  <div class="task-card__comment" v-if="t.description">{{ t.description }}</div>
                </div>
                <div class="task-card__actions">
                  <span class="reward-badge">⭐ {{ t.reward }}</span>
                  <button class="icon-btn icon-btn--approve" @click="approve(t)" title="Одобрить">✓</button>
                  <button class="icon-btn icon-btn--reject" @click="openReject(t)" title="Вернуть">✗</button>
                  <button class="icon-btn icon-btn--delete" @click="deleteTask(t)" title="Удалить">🗑</button>
                </div>
              </div>
            </div>

            <div class="section">
              <h3 class="section-title">Цели без цены</h3>
              <div v-if="awaitingWishes.length === 0" class="empty">Все цели оценены</div>
              <div v-for="w in awaitingWishes" :key="w.wish_id" class="wish-card">
                <div class="wish-card__info">
                  <div class="wish-card__title">{{ w.title }}</div>
                  <div class="wish-card__desc" v-if="w.description">{{ w.description }}</div>
                </div>
                <div class="wish-card__price-row">
                  <input type="number" v-model.number="priceInputs[w.wish_id]" placeholder="⭐ Стоимость" class="price-input" min="1" />
                  <button class="btn-small" @click="setPrice(w)">Установить</button>
                </div>
              </div>
            </div>
          </div>

          <!-- ЗАДАНИЯ -->
          <div v-if="view === 'tasks'">
            <div class="filter-bar">
              <button
                v-for="f in taskFilters" :key="f.value"
                :class="['chip', { 'chip--active': taskFilter === f.value }]"
                @click="taskFilter = f.value">{{ f.label }}</button>
            </div>
            <div class="add-task-row">
              <button class="btn-add" @click="showAddTask = true">+ Новое задание</button>
            </div>
            <div v-if="filteredTasks.length === 0" class="empty">Нет заданий</div>
            <div v-for="t in filteredTasks" :key="t.task_id" :class="['task-card', `task-card--${t.status}`]">
              <div class="task-card__info">
                <div class="task-card__title">{{ t.title }}</div>
                <div class="task-card__comment" v-if="t.description">{{ t.description }}</div>
                <div class="task-card__rework" v-if="t.rejection_comment">💬 {{ t.rejection_comment }}</div>
                <div class="task-card__status-label">{{ statusLabel(t.status) }}</div>
              </div>
              <div class="task-card__actions">
                <span class="reward-badge">⭐ {{ t.reward }}</span>
                <button class="icon-btn icon-btn--approve" v-if="t.status === 'pending_review'" @click="approve(t)" title="Одобрить">✓</button>
                <button class="icon-btn icon-btn--reject" v-if="t.status === 'pending_review'" @click="openReject(t)" title="Вернуть">✗</button>
                <button class="icon-btn icon-btn--edit" v-if="t.status === 'active'" @click="openEdit(t)" title="Редактировать">✏</button>
                <button class="icon-btn icon-btn--delete" @click="deleteTask(t)" title="Удалить">🗑</button>
              </div>
            </div>
          </div>

          <!-- ВИШЛИСТ -->
          <div v-if="view === 'wishlist'">
            <div class="filter-bar">
              <button
                v-for="f in wishFilters" :key="f.value"
                :class="['chip', { 'chip--active': wishFilter === f.value }]"
                @click="wishFilter = f.value">{{ f.label }}</button>
            </div>
            <div v-if="filteredWishes.length === 0" class="empty">Нет целей</div>
            <div v-for="w in filteredWishes" :key="w.wish_id" class="wish-card wish-card--full">
              <div class="wish-card__top">
                <div class="wish-card__info">
                  <div class="wish-card__title">{{ w.title }}</div>
                  <div class="wish-card__desc" v-if="w.description">{{ w.description }}</div>
                </div>
                <div class="wish-card__status-col">
                  <span :class="['wish-status', `wish-status--${w.status}`]">{{ wishStatusLabel(w.status) }}</span>
                </div>
              </div>

              <div class="wish-steps">
                <span :class="['step', 'step--done']">⭐ Создано</span>
                <span class="arrow">→</span>
                <span :class="['step', { 'step--done': w.status === 'purchased' || w.status === 'delivered' }]">🛒 Куплено</span>
                <span class="arrow">→</span>
                <span :class="['step', { 'step--done': w.status === 'delivered' }]">🎁 Доставлено</span>
              </div>

              <div v-if="w.price" class="wish-progress-bar">
                <div class="wish-progress-fill" :style="{ width: Math.min(100, (selectedChild.balance / w.price) * 100) + '%' }"></div>
              </div>

              <div class="wish-card__actions">
                <div v-if="w.status === 'awaiting_price'" class="price-set-row">
                  <input type="number" v-model.number="priceInputs[w.wish_id]" placeholder="⭐ Стоимость" class="price-input" min="1" />
                  <button class="btn-small" @click="setPrice(w)">Установить</button>
                </div>
                <div v-else class="wish-price-display">⭐ {{ w.price ?? '?' }}</div>

                <div class="wish-btns">
                  <button
                    v-if="w.status === 'purchased'"
                    class="btn-deliver"
                    @click="deliver(w)"
                    :disabled="delivering === w.wish_id">
                    {{ delivering === w.wish_id ? '...' : '🎁 Доставлено' }}
                  </button>
                  <button class="icon-btn icon-btn--delete" @click="deleteWish(w)" title="Удалить">🗑</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Модалки -->
    <div v-if="rejectModal" class="modal-overlay" @click.self="rejectModal = false">
      <div class="modal">
        <h3>Вернуть на доработку</h3>
        <textarea v-model="rejectComment" placeholder="Комментарий для ребёнка..." rows="3" class="modal-textarea"></textarea>
        <div class="modal-actions">
          <button class="btn-outline" @click="rejectModal = false">Отмена</button>
          <button class="btn-primary-sm" @click="confirmReject">Вернуть</button>
        </div>
      </div>
    </div>

    <div v-if="showAddTask" class="modal-overlay" @click.self="showAddTask = false">
      <div class="modal">
        <h3>Новое задание</h3>
        <div class="field">
          <label>Название</label>
          <input v-model="addForm.title" type="text" placeholder="Убрать комнату" />
        </div>
        <div class="field">
          <label>Награда (⭐)</label>
          <input v-model.number="addForm.reward" type="number" min="1" placeholder="10" />
        </div>
        <div class="field">
          <label>Описание</label>
          <textarea v-model="addForm.description" rows="2" class="modal-textarea" placeholder="Необязательно..."></textarea>
        </div>
        <div v-if="addError" class="error-msg">{{ addError }}</div>
        <div class="modal-actions">
          <button class="btn-outline" @click="showAddTask = false">Отмена</button>
          <button class="btn-primary-sm" @click="createTask" :disabled="saving">
            {{ saving ? '...' : 'Добавить' }}
          </button>
        </div>
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
          <textarea v-model="editForm.description" rows="2" class="modal-textarea"></textarea>
        </div>
        <div class="modal-actions">
          <button class="btn-outline" @click="editModal = false">Отмена</button>
          <button class="btn-primary-sm" @click="saveEdit">Сохранить</button>
        </div>
      </div>
    </div>
  </ParentLayout>
</template>

<script>
import ParentLayout from '../../components/ParentLayout.vue'
import { useApi } from '../../composables/useApi'
import { useAuth } from '../../composables/useAuth'
export default {
  name: 'ParentDashboard',
  components: { ParentLayout },
  data() {
    return {
      loading: true, children: [], selectedChild: null,
      tasks: [], wishes: [],
      progress: { wishes_done: 0, wishes_total: 0, tasks_done: 0, tasks_total: 0 },
      priceInputs: {},
      view: 'overview',
      taskFilter: '', wishFilter: '',
      taskFilters: [
        { value: '', label: 'Все' },
        { value: 'active', label: 'Активные' },
        { value: 'pending_review', label: 'На проверке' },
        { value: 'needs_rework', label: 'На доработке' },
        { value: 'completed', label: 'Выполненные' },
      ],
      wishFilters: [
        { value: '', label: 'Все' },
        { value: 'awaiting_price', label: 'Без цены' },
        { value: 'available', label: 'Доступные' },
        { value: 'purchased', label: 'Куплено' },
        { value: 'delivered', label: 'Доставлено' },
      ],
      rejectModal: false, rejectComment: '', rejectTarget: null,
      showAddTask: false, saving: false, addError: '',
      addForm: { title: '', reward: '', description: '' },
      editModal: false, editTarget: null, editForm: { title: '', reward: '', description: '' },
      delivering: null,
    }
  },
  computed: {
    user() { return useAuth().user.value },
    pendingTasks() { return this.tasks.filter(t => t.status === 'pending_review') },
    awaitingWishes() { return this.wishes.filter(w => w.status === 'awaiting_price') },
    filteredTasks() {
      if (!this.taskFilter) return this.tasks
      return this.tasks.filter(t => t.status === this.taskFilter)
    },
    filteredWishes() {
      if (!this.wishFilter) return this.wishes
      return this.wishes.filter(w => w.status === this.wishFilter)
    },
    wishProgress() {
      if (!this.progress.wishes_total) return 0
      return Math.round((this.progress.wishes_done / this.progress.wishes_total) * 100)
    },
    taskProgress() {
      if (!this.progress.tasks_total) return 0
      return Math.round((this.progress.tasks_done / this.progress.tasks_total) * 100)
    }
  },
  async mounted() {
    const { getChildren } = useApi()
    try {
      const res = await getChildren()
      this.children = res.data.children || []
      if (this.children.length) await this.selectChild(this.children[0])
    } finally { this.loading = false }
  },
  methods: {
    async selectChild(child) {
      this.selectedChild = child
      const { getWishes, getTasks } = useApi()
      const [w, t] = await Promise.all([
        getWishes(child.child_id),
        getTasks({ child_id: child.child_id })
      ])
      this.wishes = w.data.wishes || []
      this.tasks = t.data.tasks || []
      this.progress = {
        wishes_total: this.wishes.length,
        wishes_done: this.wishes.filter(x => x.status === 'delivered').length,
        tasks_total: this.tasks.length,
        tasks_done: this.tasks.filter(x => x.status === 'completed').length
      }
    },
    async reload() { await this.selectChild(this.selectedChild) },
    statusLabel(s) {
      return { active: 'Активное', pending_review: 'На проверке', needs_rework: 'На доработке', completed: 'Выполнено' }[s] || s
    },
    wishStatusLabel(s) {
      return { awaiting_price: 'Без цены', available: 'Доступно', purchased: 'Куплено', delivered: 'Доставлено' }[s] || s
    },
    async approve(task) {
      const { approveTask } = useApi()
      await approveTask(task.task_id)
      await this.reload()
    },
    openReject(task) { this.rejectTarget = task; this.rejectComment = ''; this.rejectModal = true },
    async confirmReject() {
      const { rejectTask } = useApi()
      await rejectTask(this.rejectTarget.task_id, this.rejectComment)
      this.rejectModal = false
      await this.reload()
    },
    async deleteTask(task) {
      if (!confirm('Удалить задание?')) return
      const { deleteTask } = useApi()
      await deleteTask(task.task_id)
      await this.reload()
    },
    openEdit(t) {
      this.editTarget = t
      this.editForm = { title: t.title, reward: t.reward, description: t.description || '' }
      this.editModal = true
    },
    async saveEdit() {
      const { updateTask } = useApi()
      await updateTask(this.editTarget.task_id, {
        title: this.editForm.title,
        reward: this.editForm.reward,
        description: this.editForm.description || undefined
      })
      this.editModal = false
      await this.reload()
    },
    async createTask() {
      if (!this.addForm.title || !this.addForm.reward) { this.addError = 'Заполните название и награду'; return }
      this.saving = true; this.addError = ''
      const { createTask } = useApi()
      try {
        await createTask({
          title: this.addForm.title,
          child_id: this.selectedChild.child_id,
          reward: this.addForm.reward,
          description: this.addForm.description || undefined
        })
        this.showAddTask = false
        this.addForm = { title: '', reward: '', description: '' }
        await this.reload()
      } catch (e) { this.addError = e.response?.data?.error?.message || 'Ошибка' }
      finally { this.saving = false }
    },
    async setPrice(wish) {
      const price = this.priceInputs[wish.wish_id]
      if (!price || price <= 0) return
      const { updateWish } = useApi()
      await updateWish(this.selectedChild.child_id, wish.wish_id, { price })
      await this.reload()
    },
    async deliver(wish) {
      this.delivering = wish.wish_id
      const { deliverWish } = useApi()
      try { await deliverWish(this.selectedChild.child_id, wish.wish_id); await this.reload() }
      catch (e) { alert(e.response?.data?.error?.message || 'Ошибка') }
      finally { this.delivering = null }
    },
    async deleteWish(wish) {
      if (!confirm('Удалить цель?')) return
      const { deleteWish } = useApi()
      await deleteWish(this.selectedChild.child_id, wish.wish_id)
      await this.reload()
    }
  }
}
</script>

<style scoped>
.loading { text-align: center; padding: 60px; color: #888; }
.greeting { font-size: 24px; font-weight: 700; margin-bottom: 20px; }
.no-child { text-align: center; padding: 60px 20px; color: #888; }
.no-child p { margin-bottom: 16px; font-size: 16px; }
.btn-primary { padding: 12px 24px; background: #4f7ef7; color: #fff; border: none; border-radius: 8px; font-size: 15px; font-weight: 600; cursor: pointer; }

.child-tabs { display: flex; gap: 8px; margin-bottom: 16px; flex-wrap: wrap; }
.tab { padding: 8px 18px; border-radius: 20px; border: 2px solid #4f7ef7; background: transparent; color: #4f7ef7; font-weight: 600; font-size: 14px; cursor: pointer; transition: all 0.2s; }
.tab--active { background: #4f7ef7; color: #fff; }

.balance-card { background: #fff; border-radius: 16px; padding: 18px 20px; margin-bottom: 14px; display: flex; align-items: center; justify-content: space-between; box-shadow: 0 2px 8px rgba(0,0,0,0.06); }
.balance-label { font-size: 15px; color: #555; font-weight: 600; }
.balance-meta { font-size: 12px; color: #aaa; margin-top: 2px; }
.balance-value { font-size: 28px; font-weight: 800; color: #4f7ef7; }

.progress-row { display: flex; gap: 12px; margin-bottom: 16px; }
.progress-item { background: #fff; border-radius: 14px; padding: 16px; flex: 1; display: flex; flex-direction: column; align-items: center; gap: 6px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); }
.progress-circle { position: relative; width: 52px; height: 52px; }
.progress-circle svg { width: 100%; height: 100%; }
.progress-text { position: absolute; top: 50%; left: 50%; transform: translate(-50%,-50%); font-size: 10px; font-weight: 700; color: #333; }
.progress-label { font-size: 12px; color: #666; }

.view-tabs { display: flex; background: #eef1ff; border-radius: 12px; padding: 4px; gap: 4px; margin-bottom: 16px; }
.vtab { flex: 1; padding: 9px; border: none; background: transparent; color: #666; font-size: 14px; font-weight: 600; cursor: pointer; border-radius: 8px; transition: all 0.15s; font-family: inherit; }
.vtab--active { background: #fff; color: #4f7ef7; box-shadow: 0 2px 8px rgba(0,0,0,0.08); }

.filter-bar { display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 12px; }
.chip { padding: 5px 12px; border-radius: 16px; border: 1.5px solid #d0d8ff; background: #fff; color: #4f7ef7; font-size: 13px; font-weight: 500; cursor: pointer; transition: all 0.15s; }
.chip--active { background: #4f7ef7; color: #fff; border-color: #4f7ef7; }

.add-task-row { margin-bottom: 12px; }
.btn-add { padding: 9px 18px; background: #f0f4ff; border: 2px dashed #4f7ef7; border-radius: 10px; color: #4f7ef7; font-size: 14px; font-weight: 600; cursor: pointer; }

.section { background: #fff; border-radius: 16px; padding: 18px; margin-bottom: 14px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); }
.section-title { font-size: 15px; font-weight: 700; margin-bottom: 12px; color: #333; }
.empty { color: #bbb; font-size: 14px; text-align: center; padding: 20px; }

.task-card { display: flex; align-items: flex-start; justify-content: space-between; padding: 12px 14px; background: #f8f9ff; border-radius: 10px; margin-bottom: 8px; border-left: 4px solid #e0e0e0; gap: 8px; }
.task-card--active { border-left-color: #4f7ef7; }
.task-card--pending_review { border-left-color: #f59e0b; }
.task-card--needs_rework { border-left-color: #e53e3e; }
.task-card--completed { border-left-color: #22c55e; opacity: 0.75; }
.task-card__info { flex: 1; }
.task-card__title { font-size: 14px; font-weight: 600; color: #1a1a1a; }
.task-card__comment { font-size: 12px; color: #888; margin-top: 2px; }
.task-card__rework { font-size: 12px; color: #e53e3e; margin-top: 4px; }
.task-card__status-label { font-size: 11px; color: #aaa; margin-top: 4px; }
.task-card__actions { display: flex; align-items: center; gap: 5px; flex-shrink: 0; flex-wrap: wrap; justify-content: flex-end; }
.reward-badge { font-size: 13px; font-weight: 700; color: #4f7ef7; background: #eef2ff; padding: 3px 8px; border-radius: 10px; white-space: nowrap; }
.icon-btn { width: 30px; height: 30px; border-radius: 7px; border: none; font-size: 14px; cursor: pointer; display: flex; align-items: center; justify-content: center; transition: all 0.15s; }
.icon-btn--approve { background: #e6f9f0; color: #22c55e; }
.icon-btn--approve:hover { background: #22c55e; color: #fff; }
.icon-btn--reject { background: #fff5f5; color: #e53e3e; }
.icon-btn--reject:hover { background: #e53e3e; color: #fff; }
.icon-btn--edit { background: #f0f4ff; color: #4f7ef7; }
.icon-btn--edit:hover { background: #4f7ef7; color: #fff; }
.icon-btn--delete { background: #f5f5f5; color: #bbb; }
.icon-btn--delete:hover { background: #e53e3e; color: #fff; }

.wish-card { padding: 14px; background: #f8f9ff; border-radius: 12px; margin-bottom: 10px; }
.wish-card--full { border: 1px solid #e8eeff; }
.wish-card__top { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 10px; }
.wish-card__info { flex: 1; }
.wish-card__title { font-size: 15px; font-weight: 600; margin-bottom: 2px; }
.wish-card__desc { font-size: 13px; color: #888; }
.wish-card__status-col { margin-left: 10px; }
.wish-status { padding: 4px 10px; border-radius: 10px; font-size: 12px; font-weight: 600; }
.wish-status--awaiting_price { background: #fef9e7; color: #b45309; }
.wish-status--available { background: #eef2ff; color: #4f7ef7; }
.wish-status--purchased { background: #fef3c7; color: #92400e; }
.wish-status--delivered { background: #dcfce7; color: #166534; }
.wish-steps { display: flex; align-items: center; gap: 6px; flex-wrap: wrap; margin-bottom: 10px; }
.step { font-size: 12px; color: #ccc; font-weight: 600; }
.step--done { color: #4f7ef7; }
.arrow { color: #ddd; font-size: 12px; }
.wish-progress-bar { height: 6px; background: #e0e0e0; border-radius: 3px; margin-bottom: 10px; overflow: hidden; }
.wish-progress-fill { height: 100%; background: #4f7ef7; border-radius: 3px; transition: width 0.3s; }
.wish-card__actions { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.price-set-row { display: flex; gap: 6px; flex: 1; }
.price-input { flex: 1; padding: 7px 10px; border: 1px solid #ddd; border-radius: 8px; font-size: 13px; outline: none; min-width: 0; }
.price-input:focus { border-color: #4f7ef7; }
.wish-price-display { font-size: 18px; font-weight: 700; color: #4f7ef7; }
.btn-small { padding: 7px 12px; background: #4f7ef7; color: #fff; border: none; border-radius: 8px; font-size: 13px; font-weight: 600; cursor: pointer; white-space: nowrap; }
.wish-btns { display: flex; gap: 6px; align-items: center; }
.btn-deliver { padding: 8px 14px; background: #22c55e; color: #fff; border: none; border-radius: 8px; font-size: 13px; font-weight: 600; cursor: pointer; white-space: nowrap; }
.btn-deliver:hover:not(:disabled) { background: #16a34a; }
.btn-deliver:disabled { opacity: 0.6; cursor: not-allowed; }

.wish-card__price-row { display: flex; gap: 8px; margin-top: 8px; }

.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.45); display: flex; align-items: center; justify-content: center; z-index: 300; }
.modal { background: #fff; border-radius: 16px; padding: 24px; width: 90%; max-width: 380px; }
.modal h3 { font-size: 18px; font-weight: 700; margin-bottom: 16px; }
.field { margin-bottom: 12px; }
.field label { display: block; font-size: 13px; color: #666; margin-bottom: 4px; font-weight: 500; }
.field input { width: 100%; padding: 10px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 15px; outline: none; font-family: inherit; }
.field input:focus { border-color: #4f7ef7; }
.modal-textarea { width: 100%; border: 1px solid #ddd; border-radius: 8px; padding: 10px; font-size: 14px; outline: none; resize: vertical; font-family: inherit; }
.modal-textarea:focus { border-color: #4f7ef7; }
.modal-actions { display: flex; gap: 8px; margin-top: 16px; justify-content: flex-end; }
.btn-outline { padding: 10px 18px; background: transparent; border: 2px solid #ddd; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-primary-sm { padding: 10px 20px; background: #4f7ef7; color: #fff; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-primary-sm:hover:not(:disabled) { background: #3a6be0; }
.btn-primary-sm:disabled { opacity: 0.6; cursor: not-allowed; }
.error-msg { color: #e53e3e; font-size: 13px; margin-bottom: 8px; }
</style>