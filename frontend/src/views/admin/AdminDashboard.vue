<template>
  <div class="admin-layout">
    <aside class="sidebar">
      <div class="sidebar__logo">kid<span class="accent">TASK</span> <span class="admin-badge">Admin</span></div>
      <div class="sidebar__nav">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          class="nav-btn"
          :class="{ active: activeTab === tab.id }"
          @click="activeTab = tab.id"
        >
          <span class="nav-icon">{{ tab.icon }}</span>
          {{ tab.label }}
        </button>
      </div>
      <button class="logout-btn" @click="logout">Выйти</button>
    </aside>

    <main class="admin-main">

      <div v-if="activeTab === 'stats'" class="tab-content">
        <h2 class="tab-title">Аналитика и мониторинг</h2>
        <div v-if="stats" class="stats-grid">
          <div class="stat-card">
            <div class="stat-value">{{ stats.total_families }}</div>
            <div class="stat-label">Семей</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ stats.total_children }}</div>
            <div class="stat-label">Детей</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ stats.completed_tasks }}</div>
            <div class="stat-label">Выполнено задач</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ stats.avg_tasks_done.toFixed(1) }}</div>
            <div class="stat-label">Задач на ребёнка</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ stats.total_wishes }}</div>
            <div class="stat-label">Желаний</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ stats.total_balance }}</div>
            <div class="stat-label">Монет</div>
          </div>
          <div class="stat-card stat-card--perf">
            <div class="stat-value">{{ responseTime }}мс</div>
            <div class="stat-label">Отклик API /health</div>
          </div>
        </div>
        <div v-else class="loading">Загрузка...</div>
      </div>

      <div v-if="activeTab === 'families'" class="tab-content">
        <h2 class="tab-title">Управление семьями</h2>
        <div v-if="loadingFamilies" class="loading">Загрузка...</div>
        <table v-else class="data-table">
          <thead>
            <tr>
              <th>Родитель</th>
              <th>Email</th>
              <th>Детей</th>
              <th>Задач выполнено</th>
              <th>Статус</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="f in families" :key="f.parent_id">
              <td>{{ f.parent_name }}</td>
              <td>{{ f.parent_email }}</td>
              <td>{{ f.child_count }}</td>
              <td>{{ f.tasks_done }}</td>
              <td><span class="badge" :class="f.is_blocked ? 'badge--blocked' : 'badge--active'">{{ f.is_blocked ? 'Заблокирован' : 'Активен' }}</span></td>
              <td class="actions">
                <button v-if="!f.is_blocked" class="btn-sm btn-warn" @click="blockFamily(f)">Блок</button>
                <button v-else class="btn-sm btn-ok" @click="unblockFamily(f)">Разблок</button>
                <button class="btn-sm btn-danger" @click="deleteFamily(f)">Удалить</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="activeTab === 'children'" class="tab-content">
        <h2 class="tab-title">Управление детьми</h2>
        <div v-if="loadingChildren" class="loading">Загрузка...</div>
        <table v-else class="data-table">
          <thead>
            <tr>
              <th>Имя</th>
              <th>Никнейм</th>
              <th>Родитель</th>
              <th>Баланс</th>
              <th>Статус</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="c in children" :key="c.child_id">
              <td>{{ c.name }}</td>
              <td>@{{ c.username }}</td>
              <td>{{ c.parent_name }}</td>
              <td>⭐ {{ c.balance }}</td>
              <td><span class="badge" :class="c.is_blocked ? 'badge--blocked' : 'badge--active'">{{ c.is_blocked ? 'Заблокирован' : 'Активен' }}</span></td>
              <td class="actions">
                <button v-if="!c.is_blocked" class="btn-sm btn-warn" @click="blockChild(c)">Блок</button>
                <button v-else class="btn-sm btn-ok" @click="unblockChild(c)">Разблок</button>
                <button class="btn-sm btn-primary" @click="openBalanceModal(c)">Баланс</button>
                <button class="btn-sm btn-secondary" @click="openLogsModal(c)">Логи</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="activeTab === 'complaints'" class="tab-content">
        <h2 class="tab-title">Жалобы пользователей</h2>
        <div v-if="loadingComplaints" class="loading">Загрузка...</div>
        <div v-else>
          <div v-if="!complaints.length" class="empty">Жалоб нет</div>
          <div v-for="c in complaints" :key="c.complaint_id" class="complaint-card">
            <div class="complaint-header">
              <span class="complaint-author">{{ c.parent_name }}</span>
              <span class="complaint-date">{{ formatDate(c.created_at) }}</span>
              <span class="badge" :class="c.status === 'open' ? 'badge--open' : 'badge--resolved'">{{ c.status === 'open' ? 'Открыта' : 'Решена' }}</span>
            </div>
            <div class="complaint-subject">{{ c.subject }}</div>
            <div class="complaint-body">{{ c.body }}</div>
            <button v-if="c.status === 'open'" class="btn-sm btn-ok" @click="resolveComplaint(c)">Отметить решённой</button>
          </div>
        </div>
      </div>

      <div v-if="activeTab === 'chat'" class="tab-content chat-layout">
        <div class="chat-sidebar">
          <h3 class="chat-sidebar-title">Пользователи</h3>
          <div
            v-for="p in chatParents"
            :key="p.parent_id"
            class="chat-user"
            :class="{ active: selectedChatParent && selectedChatParent.parent_id === p.parent_id }"
            @click="selectChatParent(p)"
          >
            <div class="chat-user-name">{{ p.name }}</div>
            <div class="chat-user-email">{{ p.email }}</div>
          </div>
          <div v-if="!chatParents.length" class="empty">Нет пользователей</div>
        </div>
        <div class="chat-window" v-if="selectedChatParent">
          <div class="chat-messages" ref="chatMessages">
            <div v-if="!chatMessages.length" class="empty">Нет сообщений</div>
            <div
              v-for="msg in chatMessages"
              :key="msg.message_id"
              class="chat-bubble"
              :class="msg.from_admin ? 'chat-bubble--admin' : 'chat-bubble--user'"
            >
              <div class="chat-bubble-body">{{ msg.body }}</div>
              <div class="chat-bubble-time">{{ formatDate(msg.created_at) }}</div>
            </div>
          </div>
          <div class="chat-input-row">
            <input v-model="chatInput" @keyup.enter="sendAdminMessage" placeholder="Написать сообщение..." class="chat-input" />
            <button class="btn-primary" @click="sendAdminMessage">Отправить</button>
          </div>
        </div>
        <div v-else class="chat-window chat-window--empty">
          <div class="empty">Выберите пользователя</div>
        </div>
      </div>

      <div v-if="activeTab === 'wishes'" class="tab-content">
        <h2 class="tab-title">Вишлист (все желания)</h2>
        <div class="wish-filters">
          <button class="filter-btn" :class="{ active: wishSort === 'newest' }" @click="setWishSort('newest')">Новые</button>
          <button class="filter-btn" :class="{ active: wishSort === 'oldest' }" @click="setWishSort('oldest')">Старые</button>
          <button class="filter-btn" :class="{ active: wishSort === 'price_desc' }" @click="setWishSort('price_desc')">Дорогие</button>
          <button class="filter-btn" :class="{ active: wishSort === 'price_asc' }" @click="setWishSort('price_asc')">Дешёвые</button>
        </div>
        <div v-if="loadingWishes" class="loading">Загрузка...</div>
        <table v-else class="data-table">
          <thead>
            <tr>
              <th>Название</th>
              <th>Ребёнок</th>
              <th>Семья</th>
              <th>Цена</th>
              <th>Статус</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="w in wishes" :key="w.wish_id">
              <td>{{ w.title }}</td>
              <td>{{ w.child_name }}</td>
              <td>{{ w.parent_name }}</td>
              <td>{{ w.price != null ? '⭐ ' + w.price : '—' }}</td>
              <td><span class="badge" :class="wishStatusClass(w.status)">{{ wishStatusLabel(w.status) }}</span></td>
            </tr>
          </tbody>
        </table>
      </div>

    </main>

    <div v-if="balanceModal.show" class="modal-overlay" @click.self="balanceModal.show = false">
      <div class="modal">
        <h3>Корректировка баланса</h3>
        <div class="modal-info">{{ balanceModal.child?.name }} · @{{ balanceModal.child?.username }}</div>
        <div class="field">
          <label>Изменение (положительное или отрицательное)</label>
          <input v-model.number="balanceModal.delta" type="number" placeholder="например: 50 или -10" />
        </div>
        <div class="field">
          <label>Причина</label>
          <input v-model="balanceModal.reason" type="text" placeholder="Технический сбой, ручная корректировка..." />
        </div>
        <div v-if="balanceModal.error" class="error-msg">{{ balanceModal.error }}</div>
        <div class="modal-actions">
          <button class="btn-secondary" @click="balanceModal.show = false">Отмена</button>
          <button class="btn-primary" @click="submitBalance" :disabled="balanceModal.saving">
            {{ balanceModal.saving ? 'Сохранение...' : 'Применить' }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="logsModal.show" class="modal-overlay" @click.self="logsModal.show = false">
      <div class="modal modal--wide">
        <h3>Лог баланса — {{ logsModal.child?.name }}</h3>
        <div v-if="logsModal.loading" class="loading">Загрузка...</div>
        <div v-else>
          <div v-if="!logsModal.logs.length" class="empty">Нет записей</div>
          <table v-else class="data-table">
            <thead>
              <tr><th>Дата</th><th>Изменение</th><th>Причина</th></tr>
            </thead>
            <tbody>
              <tr v-for="l in logsModal.logs" :key="l.log_id">
                <td>{{ formatDate(l.created_at) }}</td>
                <td :class="l.delta > 0 ? 'pos' : 'neg'">{{ l.delta > 0 ? '+' : '' }}{{ l.delta }}</td>
                <td>{{ l.reason }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="modal-actions">
          <button class="btn-secondary" @click="logsModal.show = false">Закрыть</button>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
import { useAdminApi } from '../../composables/useApi'

export default {
  name: 'AdminDashboard',
  data() {
    return {
      activeTab: 'stats',
      tabs: [
        { id: 'stats', icon: '📊', label: 'Аналитика' },
        { id: 'families', icon: '👨‍👩‍👧', label: 'Семьи' },
        { id: 'children', icon: '👶', label: 'Дети' },
        { id: 'complaints', icon: '📋', label: 'Жалобы' },
        { id: 'chat', icon: '💬', label: 'Чат поддержки' },
        { id: 'wishes', icon: '🎁', label: 'Вишлист' },
      ],
      stats: null,
      responseTime: null,

      families: [],
      loadingFamilies: false,

      children: [],
      loadingChildren: false,

      complaints: [],
      loadingComplaints: false,

      chatParents: [],
      selectedChatParent: null,
      chatMessages: [],
      chatInput: '',

      wishes: [],
      loadingWishes: false,
      wishSort: 'newest',

      balanceModal: { show: false, child: null, delta: 0, reason: '', saving: false, error: '' },
      logsModal: { show: false, child: null, logs: [], loading: false },
    }
  },
  watch: {
    activeTab(tab) {
      if (tab === 'stats') this.loadStats()
      else if (tab === 'families') this.loadFamilies()
      else if (tab === 'children') this.loadChildren()
      else if (tab === 'complaints') this.loadComplaints()
      else if (tab === 'chat') this.loadChatParents()
      else if (tab === 'wishes') { this.wishSort = 'newest'; this.loadWishes() }
    }
  },
  async mounted() {
    await this.loadStats()
  },
  methods: {
    async loadStats() {
      const { getStats } = useAdminApi()
      try {
        const t0 = Date.now()
        await fetch('/health')
        this.responseTime = Date.now() - t0
        const res = await getStats()
        this.stats = res.data
      } catch {}
    },
    async loadFamilies() {
      const { getFamilies } = useAdminApi()
      this.loadingFamilies = true
      try {
        const res = await getFamilies()
        this.families = res.data.families || []
      } finally { this.loadingFamilies = false }
    },
    async blockFamily(f) {
      if (!confirm(`Заблокировать ${f.parent_name}?`)) return
      const { blockFamily } = useAdminApi()
      await blockFamily(f.parent_id)
      f.is_blocked = true
    },
    async unblockFamily(f) {
      const { unblockFamily } = useAdminApi()
      await unblockFamily(f.parent_id)
      f.is_blocked = false
    },
    async deleteFamily(f) {
      if (!confirm(`Удалить семью ${f.parent_name} и все данные? Это необратимо.`)) return
      const { deleteFamily } = useAdminApi()
      await deleteFamily(f.parent_id)
      this.families = this.families.filter(x => x.parent_id !== f.parent_id)
    },

    async loadChildren() {
      const { getChildren } = useAdminApi()
      this.loadingChildren = true
      try {
        const res = await getChildren()
        this.children = res.data.children || []
      } finally { this.loadingChildren = false }
    },
    async blockChild(c) {
      if (!confirm(`Заблокировать ${c.name}?`)) return
      const { blockChild } = useAdminApi()
      await blockChild(c.child_id)
      c.is_blocked = true
    },
    async unblockChild(c) {
      const { unblockChild } = useAdminApi()
      await unblockChild(c.child_id)
      c.is_blocked = false
    },
    openBalanceModal(c) {
      this.balanceModal = { show: true, child: c, delta: 0, reason: '', saving: false, error: '' }
    },
    async submitBalance() {
      const m = this.balanceModal
      if (!m.delta) { m.error = 'Укажите ненулевое значение'; return }
      if (!m.reason) { m.error = 'Укажите причину'; return }
      m.saving = true; m.error = ''
      try {
        const { adjustBalance } = useAdminApi()
        const res = await adjustBalance(m.child.child_id, m.delta, m.reason)
        const idx = this.children.findIndex(x => x.child_id === m.child.child_id)
        if (idx !== -1) this.children[idx].balance = res.data.balance
        m.show = false
      } catch (e) {
        m.error = e.response?.data?.error?.message || 'Ошибка'
      } finally { m.saving = false }
    },
    async openLogsModal(c) {
      this.logsModal = { show: true, child: c, logs: [], loading: true }
      try {
        const { getBalanceLogs } = useAdminApi()
        const res = await getBalanceLogs(c.child_id)
        this.logsModal.logs = res.data.logs || []
      } finally { this.logsModal.loading = false }
    },

    async loadComplaints() {
      const { getComplaints } = useAdminApi()
      this.loadingComplaints = true
      try {
        const res = await getComplaints()
        this.complaints = res.data.complaints || []
      } finally { this.loadingComplaints = false }
    },
    async resolveComplaint(c) {
      const { resolveComplaint } = useAdminApi()
      await resolveComplaint(c.complaint_id)
      c.status = 'resolved'
    },

    async loadChatParents() {
      const { getChatParents } = useAdminApi()
      try {
        const res = await getChatParents()
        this.chatParents = res.data.parents || []
      } catch {}
    },
    async selectChatParent(p) {
      this.selectedChatParent = p
      this.chatMessages = []
      const { getChatMessages } = useAdminApi()
      try {
        const res = await getChatMessages(p.parent_id)
        this.chatMessages = res.data.messages || []
        p.unread = 0
        this.$nextTick(() => this.scrollChat())
      } catch {}
    },
    async sendAdminMessage() {
      if (!this.chatInput.trim()) return
      const { sendChatMessage } = useAdminApi()
      try {
        const res = await sendChatMessage(this.selectedChatParent.parent_id, this.chatInput.trim())
        this.chatMessages.push(res.data.message)
        this.chatInput = ''
        this.$nextTick(() => this.scrollChat())
      } catch {}
    },
    scrollChat() {
      const el = this.$refs.chatMessages
      if (el) el.scrollTop = el.scrollHeight
    },

    async loadWishes() {
      const { getWishes } = useAdminApi()
      this.loadingWishes = true
      try {
        const res = await getWishes(this.wishSort)
        this.wishes = res.data.wishes || []
      } finally { this.loadingWishes = false }
    },
    async setWishSort(sort) {
      if (this.wishSort === sort) return
      this.wishSort = sort
      await this.loadWishes()
    },

    logout() {
      localStorage.removeItem('kt_admin_secret')
      this.$router.push('/admin/login')
    },

    formatDate(d) {
      return new Date(d).toLocaleString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' })
    },
    wishStatusLabel(s) {
      return { awaiting_price: 'Без цены', available: 'Доступно', purchased: 'Куплено', delivered: 'Доставлено' }[s] || s
    },
    wishStatusClass(s) {
      return { awaiting_price: 'badge--warn', available: 'badge--active', purchased: 'badge--ok', delivered: 'badge--resolved' }[s] || ''
    },
  }
}
</script>

<style scoped>
.admin-layout { display: flex; min-height: 100vh; background: #f5f7ff; }

.sidebar { width: 220px; background: #1e2235; display: flex; flex-direction: column; padding: 24px 0; flex-shrink: 0; }
.sidebar__logo { color: #fff; font-size: 18px; font-weight: 800; padding: 0 20px 24px; }
.accent { color: #4f7ef7; }
.admin-badge { font-size: 10px; font-weight: 600; background: #4f7ef7; color: #fff; border-radius: 6px; padding: 2px 6px; vertical-align: middle; margin-left: 6px; }
.sidebar__nav { display: flex; flex-direction: column; gap: 2px; padding: 0 10px; flex: 1; }
.logout-btn { margin: 16px 10px 0; padding: 11px 12px; border: none; border-radius: 10px; background: rgba(239,68,68,0.15); color: #f87171; font-size: 14px; font-weight: 600; cursor: pointer; text-align: left; transition: background 0.15s; }
.logout-btn:hover { background: rgba(239,68,68,0.28); }
.nav-btn { display: flex; align-items: center; gap: 10px; width: 100%; padding: 11px 12px; border: none; border-radius: 10px; background: transparent; color: #aab; font-size: 14px; font-weight: 500; cursor: pointer; text-align: left; transition: background 0.15s, color 0.15s; }
.nav-btn:hover { background: rgba(255,255,255,0.06); color: #fff; }
.nav-btn.active { background: #4f7ef7; color: #fff; }
.nav-icon { font-size: 16px; }

.admin-main { flex: 1; padding: 32px; overflow-y: auto; }
.tab-title { font-size: 22px; font-weight: 700; margin-bottom: 20px; }
.tab-content { }

.stats-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 16px; }
.stat-card { background: #fff; border-radius: 16px; padding: 20px; text-align: center; box-shadow: 0 2px 8px rgba(0,0,0,0.06); }
.stat-card--perf { border: 2px solid #4f7ef7; }
.stat-value { font-size: 32px; font-weight: 800; color: #1a1a1a; }
.stat-label { font-size: 13px; color: #888; margin-top: 4px; }

.data-table { width: 100%; border-collapse: collapse; background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.06); }
.data-table th { background: #f5f7ff; padding: 12px 14px; text-align: left; font-size: 13px; color: #666; font-weight: 600; }
.data-table td { padding: 12px 14px; font-size: 14px; border-top: 1px solid #f0f0f0; }
.data-table tr:hover td { background: #fafbff; }

.actions { display: flex; gap: 6px; flex-wrap: wrap; }
.btn-sm { padding: 5px 10px; border: none; border-radius: 7px; font-size: 12px; font-weight: 600; cursor: pointer; }
.btn-sm.btn-warn { background: #fff3cd; color: #856404; }
.btn-sm.btn-ok { background: #d4edda; color: #155724; }
.btn-sm.btn-danger { background: #ffe0e0; color: #c53030; }
.btn-sm.btn-primary { background: #4f7ef7; color: #fff; }
.btn-sm.btn-secondary { background: #eef; color: #4f7ef7; }

.badge { display: inline-block; padding: 3px 8px; border-radius: 6px; font-size: 12px; font-weight: 600; }
.badge--active { background: #d4edda; color: #155724; }
.badge--blocked { background: #ffe0e0; color: #c53030; }
.badge--open { background: #fff3cd; color: #856404; }
.badge--resolved { background: #d4edda; color: #155724; }
.badge--warn { background: #fff3cd; color: #856404; }
.badge--ok { background: #cce5ff; color: #004085; }

.complaint-card { background: #fff; border-radius: 12px; padding: 16px; margin-bottom: 12px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); }
.complaint-header { display: flex; align-items: center; gap: 10px; margin-bottom: 8px; }
.complaint-author { font-weight: 700; font-size: 14px; }
.complaint-date { font-size: 12px; color: #888; flex: 1; }
.complaint-subject { font-size: 15px; font-weight: 600; margin-bottom: 6px; }
.complaint-body { font-size: 14px; color: #555; margin-bottom: 12px; }

.chat-layout { display: flex; gap: 0; height: calc(100vh - 96px); }
.chat-sidebar { width: 220px; background: #fff; border-radius: 12px 0 0 12px; overflow-y: auto; border-right: 1px solid #eee; flex-shrink: 0; }
.chat-sidebar-title { padding: 16px; font-size: 14px; font-weight: 700; color: #555; border-bottom: 1px solid #eee; }
.chat-user { padding: 12px 16px; cursor: pointer; position: relative; border-bottom: 1px solid #f5f5f5; }
.chat-user:hover { background: #f5f7ff; }
.chat-user.active { background: #eef2ff; }
.chat-user-name { font-size: 14px; font-weight: 600; }
.chat-user-email { font-size: 12px; color: #888; }
.unread-badge { position: absolute; top: 10px; right: 12px; background: #e53e3e; color: #fff; border-radius: 10px; font-size: 11px; font-weight: 700; padding: 2px 6px; }
.chat-window { flex: 1; background: #fff; border-radius: 0 12px 12px 0; display: flex; flex-direction: column; }
.chat-window--empty { align-items: center; justify-content: center; }
.chat-messages { flex: 1; overflow-y: auto; padding: 16px; display: flex; flex-direction: column; gap: 10px; }
.chat-bubble { max-width: 70%; padding: 10px 14px; border-radius: 14px; }
.chat-bubble--user { background: #f0f0f0; align-self: flex-start; border-bottom-left-radius: 4px; }
.chat-bubble--admin { background: #4f7ef7; color: #fff; align-self: flex-end; border-bottom-right-radius: 4px; }
.chat-bubble-body { font-size: 14px; }
.chat-bubble-time { font-size: 11px; margin-top: 4px; opacity: 0.6; text-align: right; }
.chat-input-row { display: flex; gap: 10px; padding: 12px 16px; border-top: 1px solid #eee; }
.chat-input { flex: 1; padding: 10px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 14px; outline: none; }
.chat-input:focus { border-color: #4f7ef7; }

.wish-filters { display: flex; gap: 8px; margin-bottom: 16px; }
.filter-btn { padding: 8px 16px; border: 1.5px solid #ddd; border-radius: 8px; background: #fff; font-size: 13px; font-weight: 600; color: #555; cursor: pointer; }
.filter-btn.active, .filter-btn:hover { border-color: #4f7ef7; color: #4f7ef7; background: #eef2ff; }

.loading { text-align: center; padding: 60px; color: #888; }
.empty { text-align: center; padding: 40px; color: #bbb; font-size: 14px; }

.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 300; }
.modal { background: #fff; border-radius: 16px; padding: 24px; width: 90%; max-width: 380px; }
.modal--wide { max-width: 600px; }
.modal h3 { font-size: 18px; font-weight: 700; margin-bottom: 6px; }
.modal-info { font-size: 13px; color: #888; margin-bottom: 16px; }
.field { margin-bottom: 14px; }
.field label { display: block; font-size: 13px; color: #666; margin-bottom: 4px; }
.field input { width: 100%; padding: 10px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 15px; outline: none; box-sizing: border-box; }
.field input:focus { border-color: #4f7ef7; }
.error-msg { color: #e53e3e; font-size: 13px; margin-bottom: 10px; }
.modal-actions { display: flex; gap: 10px; margin-top: 12px; }
.btn-secondary { flex: 1; padding: 10px; background: #f5f5f5; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; color: #555; }
.btn-primary { flex: 1; padding: 10px; background: #4f7ef7; color: #fff; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; }
.btn-primary:hover:not(:disabled) { background: #3a6be0; }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
.pos { color: #155724; font-weight: 700; }
.neg { color: #c53030; font-weight: 700; }
</style>