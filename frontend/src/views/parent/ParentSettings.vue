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

    <section class="section">
      <div class="section-header">
        <h2 class="section-title">Чат поддержки</h2>
        <button class="toggle-chat-btn" @click="chatOpen = !chatOpen">
          {{ chatOpen ? 'Свернуть' : 'Открыть' }}
        </button>
      </div>
      <div v-if="chatOpen" class="chat-box">
        <div class="chat-messages" ref="chatBox">
          <div v-if="loadingChat" class="loading">Загрузка...</div>
          <div v-else-if="!chatMessages.length" class="empty">Нет сообщений. Напишите нам!</div>
          <div
            v-for="msg in chatMessages"
            :key="msg.message_id"
            class="chat-bubble"
            :class="msg.from_admin ? 'chat-bubble--admin' : 'chat-bubble--user'"
          >
            <div class="bubble-label">{{ msg.from_admin ? 'Поддержка' : 'Вы' }}</div>
            <div class="bubble-body">{{ msg.body }}</div>
            <div class="bubble-time">{{ formatDate(msg.created_at) }}</div>
          </div>
        </div>
        <div class="chat-input-row">
          <input
            v-model="chatInput"
            @keyup.enter="sendMessage"
            class="chat-input"
            placeholder="Написать сообщение..."
          />
          <button class="btn-send" @click="sendMessage" :disabled="!chatInput.trim()">➤</button>
        </div>
      </div>
    </section>

    <section class="section">
      <div class="section-header">
        <h2 class="section-title">Сообщить об ошибке</h2>
        <button class="toggle-chat-btn" @click="complaintOpen = !complaintOpen">
          {{ complaintOpen ? 'Свернуть' : 'Открыть' }}
        </button>
      </div>
      <div v-if="complaintOpen" class="card">
        <div class="field">
          <label>Тема</label>
          <input v-model="complaintForm.subject" type="text" placeholder="Кратко опишите проблему" />
        </div>
        <div class="field">
          <label>Описание</label>
          <textarea v-model="complaintForm.body" placeholder="Подробно опишите ошибку..." rows="4"></textarea>
        </div>
        <div v-if="complaintSuccess" class="success-msg">Жалоба отправлена. Мы разберёмся!</div>
        <div v-if="complaintError" class="error-msg">{{ complaintError }}</div>
        <button class="submit-btn" @click="sendComplaint" :disabled="complaintSending">
          {{ complaintSending ? 'Отправка...' : 'Отправить жалобу' }}
        </button>
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

      chatOpen: false,
      chatMessages: [],
      loadingChat: false,
      chatInput: '',

      complaintOpen: false,
      complaintForm: { subject: '', body: '' },
      complaintSending: false,
      complaintError: '',
      complaintSuccess: false,
    }
  },
  computed: {
    user() {
      return useAuth().user.value
    },
    parentId() {
      return this.user?.parent_id || null
    }
  },
  async mounted() {
    await this.loadChildren()
  },
  watch: {
    async chatOpen(val) {
      if (val && !this.chatMessages.length) {
        await this.loadChat()
      }
    }
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

    async loadChat() {
      if (!this.parentId) return
      this.loadingChat = true
      try {
        const { getSupportChat } = useApi()
        const res = await getSupportChat(this.parentId)
        this.chatMessages = res.data.messages || []
        this.$nextTick(() => this.scrollChat())
      } finally {
        this.loadingChat = false
      }
    },
    async sendMessage() {
      if (!this.chatInput.trim() || !this.parentId) return
      const { sendSupportMessage } = useApi()
      try {
        const res = await sendSupportMessage(this.parentId, this.chatInput.trim())
        this.chatMessages.push(res.data.message)
        this.chatInput = ''
        this.$nextTick(() => this.scrollChat())
      } catch {}
    },
    scrollChat() {
      const el = this.$refs.chatBox
      if (el) el.scrollTop = el.scrollHeight
    },

    async sendComplaint() {
      this.complaintError = ''
      this.complaintSuccess = false
      if (!this.complaintForm.subject.trim() || !this.complaintForm.body.trim()) {
        this.complaintError = 'Заполните тему и описание'
        return
      }
      if (!this.parentId) return
      this.complaintSending = true
      try {
        const { createComplaint } = useApi()
        await createComplaint(this.parentId, this.complaintForm.subject, this.complaintForm.body)
        this.complaintSuccess = true
        this.complaintForm = { subject: '', body: '' }
      } catch (e) {
        this.complaintError = e.response?.data?.error?.message || 'Ошибка отправки'
      } finally {
        this.complaintSending = false
      }
    },

    formatDate(d) {
      return new Date(d).toLocaleString('ru-RU', { day: '2-digit', month: '2-digit', hour: '2-digit', minute: '2-digit' })
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
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px; }
.section-title { font-size: 16px; font-weight: 600; color: #555; margin-bottom: 0; }
.toggle-chat-btn { padding: 5px 12px; border: 1.5px solid #4f7ef7; border-radius: 8px; background: #fff; color: #4f7ef7; font-size: 13px; font-weight: 600; cursor: pointer; }
.toggle-chat-btn:hover { background: #eef2ff; }

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

.chat-box { background: #fff; border-radius: 16px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); overflow: hidden; }
.chat-messages { height: 280px; overflow-y: auto; padding: 16px; display: flex; flex-direction: column; gap: 10px; background: #fafbff; }
.chat-bubble { max-width: 75%; padding: 10px 14px; border-radius: 14px; }
.chat-bubble--user { background: #eef2ff; align-self: flex-end; border-bottom-right-radius: 4px; }
.chat-bubble--admin { background: #fff; border: 1px solid #eee; align-self: flex-start; border-bottom-left-radius: 4px; }
.bubble-label { font-size: 11px; font-weight: 700; color: #4f7ef7; margin-bottom: 3px; }
.chat-bubble--user .bubble-label { color: #888; text-align: right; }
.bubble-body { font-size: 14px; color: #1a1a1a; }
.bubble-time { font-size: 11px; color: #aaa; margin-top: 4px; text-align: right; }
.chat-input-row { display: flex; gap: 8px; padding: 10px 12px; border-top: 1px solid #eee; background: #fff; }
.chat-input { flex: 1; padding: 10px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 14px; outline: none; }
.chat-input:focus { border-color: #4f7ef7; }
.btn-send { width: 42px; height: 42px; border-radius: 10px; background: #4f7ef7; color: #fff; border: none; font-size: 18px; cursor: pointer; flex-shrink: 0; }
.btn-send:disabled { opacity: 0.4; cursor: not-allowed; }

.field { margin-bottom: 14px; }
.field label { display: block; font-size: 13px; color: #666; margin-bottom: 4px; }
.field input, .field textarea { width: 100%; padding: 10px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 15px; outline: none; box-sizing: border-box; resize: vertical; }
.field input:focus, .field textarea:focus { border-color: #4f7ef7; }
.submit-btn { width: 100%; padding: 12px; background: #4f7ef7; color: #fff; border: none; border-radius: 10px; font-size: 15px; font-weight: 600; cursor: pointer; margin-top: 4px; }
.submit-btn:hover:not(:disabled) { background: #3a6be0; }
.submit-btn:disabled { opacity: 0.6; cursor: not-allowed; }

.success-msg { color: #155724; background: #d4edda; border-radius: 8px; padding: 10px 12px; font-size: 14px; margin-bottom: 10px; }
.error-msg { color: #e53e3e; font-size: 13px; margin-bottom: 10px; }

.loading { text-align: center; padding: 40px; color: #888; }
.empty { text-align: center; color: #bbb; padding: 20px; font-size: 14px; }

.logout-btn { width: 100%; padding: 14px; background: #fff5f5; border: 2px solid #e53e3e; border-radius: 12px; color: #e53e3e; font-size: 16px; font-weight: 600; cursor: pointer; margin-top: 8px; }
.logout-btn:hover { background: #e53e3e; color: #fff; }

.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 300; }
.modal { background: #fff; border-radius: 16px; padding: 24px; width: 90%; max-width: 360px; }
.modal h3 { font-size: 20px; font-weight: 700; text-align: center; margin-bottom: 20px; }
.hint { font-weight: 400; color: #aaa; }
.modal-actions { display: flex; gap: 10px; margin-top: 8px; }
.btn-secondary { flex: 1; padding: 12px; background: #f5f5f5; border: none; border-radius: 8px; font-size: 15px; font-weight: 600; cursor: pointer; color: #555; }
.btn-secondary:hover { background: #ebebeb; }
.btn-primary { flex: 1; padding: 12px; background: #4f7ef7; color: #fff; border: none; border-radius: 8px; font-size: 15px; font-weight: 600; cursor: pointer; }
.btn-primary:hover:not(:disabled) { background: #3a6be0; }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
</style>