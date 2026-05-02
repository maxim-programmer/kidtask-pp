<template>
  <ParentLayout>
    <h1 class="page-title">Настройки</h1>

    <section class="section">
      <h2 class="section-title">Мои данные</h2>
      <div class="card">
        <div class="profile-row">
          <div class="avatar">{{ user?.name?.[0] }}</div>
          <div>
            <div class="name">{{ user?.name }}</div>
            <div class="email">{{ user?.email }}</div>
          </div>
          <button class="edit-btn" @click="openParentModal">Изменить</button>
        </div>
      </div>
    </section>

    <section class="section">
      <h2 class="section-title">Данные детей</h2>
      <div v-if="loadingChildren" class="loading">Загрузка...</div>
      <div v-else>
        <div v-for="child in children" :key="child.child_id" class="card child-row">
          <div class="child-avatar">{{ child.name[0] }}</div>
          <div class="child-info">
            <div class="child-name">{{ child.name }}</div>
            <div class="child-meta">@{{ child.username }} · {{ child.age_group === 'junior' ? '7–10 лет' : '11–14 лет' }}</div>
          </div>
          <button class="edit-btn" @click="openChildModal(child)">Изменить</button>
        </div>
        <div v-if="!children.length" class="empty">Нет добавленных детей</div>
      </div>
    </section>

    <button class="logout-btn" @click="logout">Выйти</button>

    <div v-if="showParentModal" class="modal-overlay" @click.self="closeParentModal">
      <div class="modal">
        <h3>Мои данные</h3>
        <div class="field">
          <label>Имя</label>
          <input v-model="parentForm.name" type="text" placeholder="Ваше имя" />
        </div>
        <div class="field">
          <label>Email</label>
          <input v-model="parentForm.email" type="email" placeholder="email@example.com" />
        </div>
        <div class="field">
          <label>Новый пароль <span class="hint">(оставьте пустым, чтобы не менять)</span></label>
          <input v-model="parentForm.password" type="password" placeholder="Минимум 6 символов" />
        </div>
        <div v-if="parentError" class="error-msg">{{ parentError }}</div>
        <div class="modal-actions">
          <button class="btn-secondary" @click="closeParentModal">Отмена</button>
          <button class="btn-primary" @click="saveParent" :disabled="parentSaving">
            {{ parentSaving ? 'Сохранение...' : 'Сохранить' }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="showChildModal" class="modal-overlay" @click.self="closeChildModal">
      <div class="modal">
        <h3>Данные ребёнка</h3>
        <div class="field">
          <label>Имя</label>
          <input v-model="childForm.name" type="text" placeholder="Имя ребёнка" />
        </div>
        <div class="field">
          <label>Никнейм</label>
          <input v-model="childForm.username" type="text" placeholder="Логин для входа" />
        </div>
        <div class="field">
          <label>Новый пароль <span class="hint">(оставьте пустым, чтобы не менять)</span></label>
          <input v-model="childForm.password" type="password" placeholder="Минимум 4 символа" />
        </div>
        <div class="field">
          <label>Дата рождения</label>
          <input v-model="childForm.birthday" type="date" />
        </div>
        <div v-if="childError" class="error-msg">{{ childError }}</div>
        <div class="modal-actions">
          <button class="btn-secondary" @click="closeChildModal">Отмена</button>
          <button class="btn-primary" @click="saveChild" :disabled="childSaving">
            {{ childSaving ? 'Сохранение...' : 'Сохранить' }}
          </button>
        </div>
      </div>
    </div>
  </ParentLayout>
</template>

<script>
import ParentLayout from '../../components/ParentLayout.vue'
import { useAuth } from '../../composables/useAuth'
import { useApi } from '../../composables/useApi'

export default {
  name: 'ParentSettings',
  components: { ParentLayout },
  data() {
    return {
      children: [],
      loadingChildren: true,

      showParentModal: false,
      parentForm: { name: '', email: '', password: '' },
      parentSaving: false,
      parentError: '',

      showChildModal: false,
      editingChild: null,
      childForm: { name: '', username: '', password: '', birthday: '' },
      childSaving: false,
      childError: '',
    }
  },
  computed: {
    user() {
      return useAuth().user.value
    }
  },
  async mounted() {
    await this.loadChildren()
  },
  methods: {
    async loadChildren() {
      this.loadingChildren = true
      try {
        const { getChildren } = useApi()
        const res = await getChildren()
        this.children = res.data.children || []
      } finally {
        this.loadingChildren = false
      }
    },

    openParentModal() {
      const u = this.user || {}
      this.parentForm = { name: u.name || '', email: u.email || '', password: '' }
      this.parentError = ''
      this.showParentModal = true
    },
    closeParentModal() {
      this.showParentModal = false
    },
    async saveParent() {
      this.parentError = ''
      const payload = {}
      const u = this.user || {}
      if (this.parentForm.name !== u.name) payload.name = this.parentForm.name
      if (this.parentForm.email !== u.email) payload.email = this.parentForm.email
      if (this.parentForm.password) payload.password = this.parentForm.password

      if (!Object.keys(payload).length) {
        this.showParentModal = false
        return
      }

      this.parentSaving = true
      try {
        const { updateMe } = useApi()
        const res = await updateMe(payload)
        const { saveAuth } = useAuth()
        const token = localStorage.getItem('kt_token')
        saveAuth({ token, role: 'parent', user: res.data.user })
        this.showParentModal = false
      } catch (e) {
        this.parentError = e.response?.data?.error?.message || 'Ошибка сохранения'
      } finally {
        this.parentSaving = false
      }
    },

    openChildModal(child) {
      this.editingChild = child
      this.childForm = {
        name: child.name,
        username: child.username,
        password: '',
        birthday: child.birthday ? child.birthday.substring(0, 10) : '',
      }
      this.childError = ''
      this.showChildModal = true
    },
    closeChildModal() {
      this.showChildModal = false
      this.editingChild = null
    },
    async saveChild() {
      this.childError = ''
      const c = this.editingChild
      const payload = {}
      if (this.childForm.name !== c.name) payload.name = this.childForm.name
      if (this.childForm.username !== c.username) payload.username = this.childForm.username
      if (this.childForm.password) payload.password = this.childForm.password

      const newBirthday = this.childForm.birthday
        ? new Date(this.childForm.birthday + 'T00:00:00Z').toISOString()
        : null
      const oldBirthday = c.birthday ? new Date(c.birthday).toISOString() : null
      if (newBirthday !== oldBirthday) {
        payload.birthday = newBirthday
      }

      if (!Object.keys(payload).length) {
        this.showChildModal = false
        return
      }

      this.childSaving = true
      try {
        const { updateChild } = useApi()
        const res = await updateChild(c.child_id, payload)
        const idx = this.children.findIndex(x => x.child_id === c.child_id)
        if (idx !== -1) this.children[idx] = res.data.child
        this.showChildModal = false
      } catch (e) {
        this.childError = e.response?.data?.error?.message || 'Ошибка сохранения'
      } finally {
        this.childSaving = false
      }
    },

    logout() {
      useAuth().logout()
      this.$router.push('/')
    }
  }
}
</script>

<style scoped>
.page-title { font-size: 24px; font-weight: 700; margin-bottom: 20px; }

.section { margin-bottom: 24px; }
.section-title { font-size: 16px; font-weight: 600; color: #555; margin-bottom: 10px; }

.card { background: #fff; border-radius: 16px; padding: 16px 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); margin-bottom: 10px; }

.profile-row { display: flex; align-items: center; gap: 16px; }
.avatar { width: 52px; height: 52px; border-radius: 50%; background: #4f7ef7; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 24px; font-weight: 700; flex-shrink: 0; }
.name { font-size: 17px; font-weight: 700; }
.email { font-size: 13px; color: #888; margin-top: 2px; }

.child-row { display: flex; align-items: center; gap: 14px; }
.child-avatar { width: 44px; height: 44px; border-radius: 50%; background: #4f7ef7; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 20px; font-weight: 700; flex-shrink: 0; }
.child-info { flex: 1; }
.child-name { font-size: 16px; font-weight: 700; }
.child-meta { font-size: 13px; color: #888; margin-top: 2px; }

.edit-btn { margin-left: auto; padding: 6px 14px; border: 1.5px solid #4f7ef7; border-radius: 8px; background: #fff; color: #4f7ef7; font-size: 13px; font-weight: 600; cursor: pointer; white-space: nowrap; }
.edit-btn:hover { background: #f0f4ff; }

.loading { text-align: center; padding: 40px; color: #888; }
.empty { text-align: center; color: #bbb; padding: 20px; font-size: 14px; }

.logout-btn { width: 100%; padding: 14px; background: #fff5f5; border: 2px solid #e53e3e; border-radius: 12px; color: #e53e3e; font-size: 16px; font-weight: 600; cursor: pointer; margin-top: 8px; }
.logout-btn:hover { background: #e53e3e; color: #fff; }

.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 300; }
.modal { background: #fff; border-radius: 16px; padding: 24px; width: 90%; max-width: 360px; }
.modal h3 { font-size: 20px; font-weight: 700; text-align: center; margin-bottom: 20px; }

.field { margin-bottom: 14px; }
.field label { display: block; font-size: 13px; color: #666; margin-bottom: 4px; }
.hint { font-weight: 400; color: #aaa; }
.field input { width: 100%; padding: 10px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 15px; outline: none; box-sizing: border-box; }
.field input:focus { border-color: #4f7ef7; }

.error-msg { color: #e53e3e; font-size: 13px; margin-bottom: 10px; }

.modal-actions { display: flex; gap: 10px; margin-top: 8px; }
.btn-secondary { flex: 1; padding: 12px; background: #f5f5f5; border: none; border-radius: 8px; font-size: 15px; font-weight: 600; cursor: pointer; color: #555; }
.btn-secondary:hover { background: #ebebeb; }
.btn-primary { flex: 1; padding: 12px; background: #4f7ef7; color: #fff; border: none; border-radius: 8px; font-size: 15px; font-weight: 600; cursor: pointer; }
.btn-primary:hover:not(:disabled) { background: #3a6be0; }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
</style>