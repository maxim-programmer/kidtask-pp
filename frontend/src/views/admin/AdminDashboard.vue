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
          @click="switchTab(tab.id)"
        >
          <span class="nav-icon">{{ tab.icon }}</span>
          {{ tab.label }}
        </button>
      </div>
      <button class="logout-btn" @click="logout">Выйти</button>
    </aside>

    <div class="main-wrap">
      <header class="mobile-header">
        <div class="mobile-header__logo">kid<span class="accent">TASK</span> <span class="admin-badge">Admin</span></div>
        <button class="mobile-menu-btn" @click="mobileTabsOpen = !mobileTabsOpen">☰</button>
      </header>

      <div class="mobile-tabs-dropdown" :class="{ open: mobileTabsOpen }">
        <button
          v-for="tab in tabs" :key="tab.id"
          class="mobile-tab-item"
          :class="{ active: activeTab === tab.id }"
          @click="switchTab(tab.id); mobileTabsOpen = false"
        >
          <span>{{ tab.icon }}</span> {{ tab.label }}
        </button>
        <button class="mobile-tab-item mobile-tab-item--logout" @click="logout">👋 Выйти</button>
      </div>

      <main class="admin-main">

        <div v-if="activeTab === 'stats'" class="tab-content">
          <h2 class="tab-title">📊 Аналитика и мониторинг</h2>
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
              <div class="stat-value">{{ stats.avg_tasks_done != null ? Number(stats.avg_tasks_done).toFixed(1) : '0.0' }}</div>
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
              <div class="stat-value">{{ responseTime != null ? responseTime + 'мс' : '—' }}</div>
              <div class="stat-label">Отклик API /health</div>
            </div>
          </div>
          <div v-else class="loading">Загрузка...</div>
        </div>

        <div v-if="activeTab === 'families'" class="tab-content">
          <h2 class="tab-title">👨‍👩‍👧 Управление семьями</h2>
          <div v-if="loadingFamilies" class="loading">Загрузка...</div>
          <div v-else>
            <div v-if="families.length === 0" class="empty">Нет данных</div>
            <div v-for="f in families" :key="f.parent_id" class="mobile-card">
              <div class="mobile-card__header">
                <div class="mobile-card__title">{{ f.parent_name }}</div>
                <span class="badge" :class="f.is_blocked ? 'badge--blocked' : 'badge--active'">{{ f.is_blocked ? 'Заблокирован' : 'Активен' }}</span>
              </div>
              <div class="mobile-card__meta">{{ f.parent_email }}</div>
              <div class="mobile-card__row">
                <span class="mobile-card__info">👶 {{ f.child_count }} детей · ✅ {{ f.tasks_done }} задач</span>
              </div>
              <div class="mobile-card__actions">
                <button v-if="!f.is_blocked" class="btn-sm btn-warn" @click="blockFamily(f)">Блокировать</button>
                <button v-else class="btn-sm btn-ok" @click="unblockFamily(f)">Разблокировать</button>
                <button class="btn-sm btn-danger" @click="deleteFamily(f)">Удалить</button>
              </div>
            </div>
          </div>
        </div>

        <div v-if="activeTab === 'children'" class="tab-content">
          <h2 class="tab-title">👶 Управление детьми</h2>
          <div v-if="loadingChildren" class="loading">Загрузка...</div>
          <div v-else>
            <div v-if="children.length === 0" class="empty">Нет данных</div>
            <div v-for="c in children" :key="c.child_id" class="mobile-card">
              <div class="mobile-card__header">
                <div class="mobile-card__title">{{ c.name }} <span class="username">@{{ c.username }}</span></div>
                <span class="badge" :class="c.is_blocked ? 'badge--blocked' : 'badge--active'">{{ c.is_blocked ? 'Заблокирован' : 'Активен' }}</span>
              </div>
              <div class="mobile-card__meta">{{ c.parent_name }} · <span v-if="ageLabel(c.birthday)">{{ ageLabel(c.birthday) }}</span><span v-else>возраст не указан</span></div>
              <div class="mobile-card__coins">⭐ {{ c.balance }}</div>
              <div class="mobile-card__actions">
                <button v-if="!c.is_blocked" class="btn-sm btn-warn" @click="blockChild(c)">Блок</button>
                <button v-else class="btn-sm btn-ok" @click="unblockChild(c)">Разблок</button>
                <button class="btn-sm btn-primary" @click="openBalanceModal(c)">Баланс</button>
                <button class="btn-sm btn-secondary" @click="openLogsModal(c)">Логи</button>
              </div>
            </div>
          </div>
        </div>

        <div v-if="activeTab === 'complaints'" class="tab-content">
          <h2 class="tab-title">📋 Жалобы пользователей</h2>
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
            <div class="chat-window__header">
              <button class="chat-back-btn" @click="selectedChatParent = null">← Назад</button>
              <span class="chat-window__name">{{ selectedChatParent.name }}</span>
            </div>
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
              <button class="btn-primary" @click="sendAdminMessage">➤</button>
            </div>
          </div>
          <div v-else class="chat-window chat-window--empty">
            <div class="empty">Выберите пользователя слева</div>
          </div>
        </div>

        <div v-if="activeTab === 'wishes'" class="tab-content">
          <h2 class="tab-title">🎁 Вишлист (все желания)</h2>
          <div class="wish-filters">
            <button class="filter-btn" :class="{ active: wishSort === 'newest' }" @click="setWishSort('newest')">Новые</button>
            <button class="filter-btn" :class="{ active: wishSort === 'oldest' }" @click="setWishSort('oldest')">Старые</button>
            <button class="filter-btn" :class="{ active: wishSort === 'price_desc' }" @click="setWishSort('price_desc')">Дорогие</button>
            <button class="filter-btn" :class="{ active: wishSort === 'price_asc' }" @click="setWishSort('price_asc')">Дешёвые</button>
          </div>
          <div v-if="loadingWishes" class="loading">Загрузка...</div>
          <div v-else>
            <div v-if="wishes.length === 0" class="empty">Нет данных</div>
            <div v-for="w in wishes" :key="w.wish_id" class="mobile-card">
              <div class="mobile-card__header">
                <div class="mobile-card__title">{{ w.title }}</div>
                <span class="badge" :class="wishStatusClass(w.status)">{{ wishStatusLabel(w.status) }}</span>
              </div>
              <div class="mobile-card__meta">{{ w.child_name }} · {{ w.parent_name }}</div>
              <div class="mobile-card__coins" v-if="w.price != null">⭐ {{ w.price }}</div>
              <div class="mobile-card__coins" v-else style="color:#aaa">Цена не установлена</div>
            </div>
          </div>
        </div>

      </main>
    </div>

    <nav class="bottom-nav">
      <button v-for="tab in tabs.slice(0,5)" :key="tab.id"
        class="bottom-nav__item"
        :class="{ active: activeTab === tab.id }"
        @click="switchTab(tab.id)">
        <span class="bottom-nav__icon">{{ tab.icon }}</span>
        <span class="bottom-nav__label">{{ tab.shortLabel || tab.label }}</span>
      </button>
    </nav>

    <div v-if="balanceModal.show" class="modal-overlay" @click.self="balanceModal.show = false">
      <div class="modal">
        <h3>Корректировка баланса</h3>
        <div class="modal-info">{{ balanceModal.child?.name }} · @{{ balanceModal.child?.username }}</div>
        <div class="field">
          <label>Изменение</label>
          <input v-model.number="balanceModal.delta" type="number" placeholder="например: 50 или -10" />
        </div>
        <div class="field">
          <label>Причина</label>
          <input v-model="balanceModal.reason" type="text" placeholder="Техническая корректировка..." />
        </div>
        <div v-if="balanceModal.error" class="error-msg">{{ balanceModal.error }}</div>
        <div class="modal-actions">
          <button class="btn-secondary" @click="balanceModal.show = false">Отмена</button>
          <button class="btn-primary-modal" @click="submitBalance" :disabled="balanceModal.saving">
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
          <div v-for="l in logsModal.logs" :key="l.log_id" class="log-row">
            <div class="log-row__date">{{ formatDate(l.created_at) }}</div>
            <div class="log-row__reason">{{ l.reason }}</div>
            <div :class="['log-row__delta', l.delta > 0 ? 'pos' : 'neg']">{{ l.delta > 0 ? '+' : '' }}{{ l.delta }}</div>
          </div>
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
import { ageLabel, calcAge } from '../../composables/useAge'

export default {
  name: 'AdminDashboard',
  data() {
    return {
      activeTab: 'stats',
      mobileTabsOpen: false,
      tabs: [
        { id: 'stats', icon: '📊', label: 'Аналитика', shortLabel: 'Стат' },
        { id: 'families', icon: '👨‍👩‍👧', label: 'Семьи', shortLabel: 'Семьи' },
        { id: 'children', icon: '👶', label: 'Дети', shortLabel: 'Дети' },
        { id: 'complaints', icon: '📋', label: 'Жалобы', shortLabel: 'Жалобы' },
        { id: 'chat', icon: '💬', label: 'Чат поддержки', shortLabel: 'Чат' },
        { id: 'wishes', icon: '🎁', label: 'Вишлист', shortLabel: 'Вишлист' },
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
  async mounted() {
    if (!localStorage.getItem('kt_admin_secret')) {
      this.$router.replace('/admin/login')
      return
    }
    await Promise.all([this.loadStats(), this.loadFamilies(), this.loadChildren(), this.loadComplaints(), this.loadChatParents(), this.loadWishes()])
  },
  methods: {
    switchTab(tab) {
      this.activeTab = tab
      if (tab === 'stats') this.loadStats()
      else if (tab === 'families') this.loadFamilies()
      else if (tab === 'children') this.loadChildren()
      else if (tab === 'complaints') this.loadComplaints()
      else if (tab === 'chat') this.loadChatParents()
      else if (tab === 'wishes') this.loadWishes()
    },
    async loadStats() {
      const { getStats } = useAdminApi()
      try {
        const t0 = Date.now()
        const res = await getStats()
        this.responseTime = Date.now() - t0
        this.stats = res.data
      } catch (e) {
        this.stats = { total_families: 0, total_children: 0, completed_tasks: 0, avg_tasks_done: 0, total_wishes: 0, total_balance: 0 }
        this.responseTime = null
      }
    },
    async loadFamilies() {
      const { getFamilies } = useAdminApi()
      this.loadingFamilies = true
      try { const res = await getFamilies(); this.families = res.data.families || [] }
      catch { this.families = [] }
      finally { this.loadingFamilies = false }
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
      if (!confirm(`Удалить семью ${f.parent_name}? Это необратимо.`)) return
      const { deleteFamily } = useAdminApi()
      await deleteFamily(f.parent_id)
      this.families = this.families.filter(x => x.parent_id !== f.parent_id)
    },
    async loadChildren() {
      const { getChildren } = useAdminApi()
      this.loadingChildren = true
      try { const res = await getChildren(); this.children = res.data.children || [] }
      catch { this.children = [] }
      finally { this.loadingChildren = false }
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
      } catch (e) { m.error = e.response?.data?.error?.message || 'Ошибка' }
      finally { m.saving = false }
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
      try { const res = await getComplaints(); this.complaints = res.data.complaints || [] }
      catch { this.complaints = [] }
      finally { this.loadingComplaints = false }
    },
    async resolveComplaint(c) {
      const { resolveComplaint } = useAdminApi()
      await resolveComplaint(c.complaint_id)
      c.status = 'resolved'
    },
    async loadChatParents() {
      const { getChatParents } = useAdminApi()
      try { const res = await getChatParents(); this.chatParents = res.data.parents || [] }
      catch {}
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
      try { const res = await getWishes(this.wishSort); this.wishes = res.data.wishes || [] }
      catch { this.wishes = [] }
      finally { this.loadingWishes = false }
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
    ageLabel,
    calcAge,
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

.sidebar { width: 220px; background: #1e2235; display: flex; flex-direction: column; padding: 24px 0; flex-shrink: 0; position: sticky; top: 0; height: 100vh; }
.sidebar__logo { color: #fff; font-size: 17px; font-weight: 800; padding: 0 20px 24px; }
.accent { color: #4f7ef7; }
.admin-badge { font-size: 10px; font-weight: 600; background: #4f7ef7; color: #fff; border-radius: 6px; padding: 2px 6px; vertical-align: middle; margin-left: 6px; }
.sidebar__nav { display: flex; flex-direction: column; gap: 2px; padding: 0 10px; flex: 1; }
.logout-btn { margin: 16px 10px 0; padding: 11px 12px; border: none; border-radius: 10px; background: rgba(239,68,68,0.15); color: #f87171; font-size: 14px; font-weight: 600; cursor: pointer; text-align: left; transition: background 0.15s; font-family: inherit; }
.logout-btn:hover { background: rgba(239,68,68,0.28); }
.nav-btn { display: flex; align-items: center; gap: 10px; width: 100%; padding: 11px 12px; border: none; border-radius: 10px; background: transparent; color: #aab; font-size: 14px; font-weight: 500; cursor: pointer; text-align: left; transition: background 0.15s, color 0.15s; font-family: inherit; }
.nav-btn:hover { background: rgba(255,255,255,0.06); color: #fff; }
.nav-btn.active { background: #4f7ef7; color: #fff; }
.nav-icon { font-size: 16px; }

.mobile-header { display: none; }
.mobile-tabs-dropdown { display: none; }

.main-wrap { flex: 1; display: flex; flex-direction: column; min-width: 0; }
.admin-main { flex: 1; padding: 28px 32px; overflow-y: auto; }
.tab-title { font-size: 20px; font-weight: 700; margin-bottom: 18px; }

.stats-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(150px, 1fr)); gap: 14px; }
.stat-card { background: #fff; border-radius: 16px; padding: 18px; text-align: center; box-shadow: 0 2px 8px rgba(0,0,0,0.06); }
.stat-card--perf { border: 2px solid #4f7ef7; }
.stat-value { font-size: 30px; font-weight: 800; color: #1a1a1a; }
.stat-label { font-size: 12px; color: #888; margin-top: 4px; }

.mobile-card { background: #fff; border-radius: 14px; padding: 16px; margin-bottom: 12px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); }
.mobile-card__header { display: flex; align-items: center; justify-content: space-between; gap: 8px; margin-bottom: 6px; }
.mobile-card__title { font-size: 15px; font-weight: 700; color: #1a1a1a; }
.mobile-card__meta { font-size: 13px; color: #888; margin-bottom: 4px; }
.mobile-card__coins { font-size: 16px; font-weight: 700; color: #f59e0b; margin-bottom: 10px; }
.mobile-card__row { margin-bottom: 8px; }
.mobile-card__info { font-size: 13px; color: #666; }
.mobile-card__actions { display: flex; gap: 8px; flex-wrap: wrap; }
.username { color: #888; font-weight: 400; font-size: 13px; }

.badge { display: inline-block; padding: 3px 8px; border-radius: 6px; font-size: 11px; font-weight: 600; flex-shrink: 0; }
.badge--active { background: #d4edda; color: #155724; }
.badge--blocked { background: #ffe0e0; color: #c53030; }
.badge--open { background: #fff3cd; color: #856404; }
.badge--resolved { background: #d4edda; color: #155724; }
.badge--warn { background: #fff3cd; color: #856404; }
.badge--ok { background: #cce5ff; color: #004085; }

.btn-sm { padding: 7px 12px; border: none; border-radius: 8px; font-size: 12px; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-sm.btn-warn { background: #fff3cd; color: #856404; }
.btn-sm.btn-ok { background: #d4edda; color: #155724; }
.btn-sm.btn-danger { background: #ffe0e0; color: #c53030; }
.btn-sm.btn-primary { background: #4f7ef7; color: #fff; }
.btn-sm.btn-secondary { background: #eef; color: #4f7ef7; }

.complaint-card { background: #fff; border-radius: 12px; padding: 16px; margin-bottom: 12px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); }
.complaint-header { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; flex-wrap: wrap; }
.complaint-author { font-weight: 700; font-size: 14px; }
.complaint-date { font-size: 12px; color: #888; flex: 1; }
.complaint-subject { font-size: 15px; font-weight: 600; margin-bottom: 6px; }
.complaint-body { font-size: 14px; color: #555; margin-bottom: 12px; }

.chat-layout { display: flex; height: calc(100vh - 100px); gap: 0; }
.chat-sidebar { width: 220px; background: #fff; border-radius: 12px 0 0 12px; overflow-y: auto; border-right: 1px solid #eee; flex-shrink: 0; }
.chat-sidebar-title { padding: 14px 16px; font-size: 13px; font-weight: 700; color: #555; border-bottom: 1px solid #eee; }
.chat-user { padding: 12px 16px; cursor: pointer; border-bottom: 1px solid #f5f5f5; }
.chat-user:hover { background: #f5f7ff; }
.chat-user.active { background: #eef2ff; }
.chat-user-name { font-size: 14px; font-weight: 600; }
.chat-user-email { font-size: 12px; color: #888; }
.chat-window { flex: 1; background: #fff; border-radius: 0 12px 12px 0; display: flex; flex-direction: column; }
.chat-window--empty { align-items: center; justify-content: center; }
.chat-window__header { display: none; }
.chat-messages { flex: 1; overflow-y: auto; padding: 16px; display: flex; flex-direction: column; gap: 10px; }
.chat-bubble { max-width: 70%; padding: 10px 14px; border-radius: 14px; }
.chat-bubble--user { background: #f0f0f0; align-self: flex-start; border-bottom-left-radius: 4px; }
.chat-bubble--admin { background: #4f7ef7; color: #fff; align-self: flex-end; border-bottom-right-radius: 4px; }
.chat-bubble-body { font-size: 14px; }
.chat-bubble-time { font-size: 11px; margin-top: 4px; opacity: 0.6; text-align: right; }
.chat-input-row { display: flex; gap: 8px; padding: 12px 16px; border-top: 1px solid #eee; }
.chat-input { flex: 1; padding: 10px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 14px; outline: none; font-family: inherit; }
.chat-input:focus { border-color: #4f7ef7; }

.wish-filters { display: flex; gap: 8px; margin-bottom: 14px; flex-wrap: wrap; }
.filter-btn { padding: 7px 14px; border: 1.5px solid #ddd; border-radius: 8px; background: #fff; font-size: 13px; font-weight: 600; color: #555; cursor: pointer; font-family: inherit; }
.filter-btn.active { border-color: #4f7ef7; color: #4f7ef7; background: #eef2ff; }

.loading { text-align: center; padding: 60px; color: #888; }
.empty { text-align: center; padding: 40px; color: #bbb; font-size: 14px; }

.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 300; padding: 16px; }
.modal { background: #fff; border-radius: 16px; padding: 24px; width: 100%; max-width: 380px; max-height: 90vh; overflow-y: auto; }
.modal--wide { max-width: 560px; }
.modal h3 { font-size: 18px; font-weight: 700; margin-bottom: 6px; }
.modal-info { font-size: 13px; color: #888; margin-bottom: 16px; }
.field { margin-bottom: 14px; }
.field label { display: block; font-size: 13px; color: #666; margin-bottom: 4px; }
.field input { width: 100%; padding: 10px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 15px; outline: none; box-sizing: border-box; font-family: inherit; }
.field input:focus { border-color: #4f7ef7; }
.error-msg { color: #e53e3e; font-size: 13px; margin-bottom: 10px; }
.modal-actions { display: flex; gap: 10px; margin-top: 12px; }
.btn-secondary { flex: 1; padding: 10px; background: #f5f5f5; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; color: #555; font-family: inherit; }
.btn-primary-modal { flex: 1; padding: 10px; background: #4f7ef7; color: #fff; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-primary-modal:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-primary { padding: 10px 16px; background: #4f7ef7; color: #fff; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; font-family: inherit; }
.pos { color: #155724; font-weight: 700; }
.neg { color: #c53030; font-weight: 700; }

.log-row { display: flex; align-items: center; gap: 8px; padding: 10px 0; border-bottom: 1px solid #f5f5f5; }
.log-row__date { font-size: 12px; color: #aaa; width: 100px; flex-shrink: 0; }
.log-row__reason { flex: 1; font-size: 13px; color: #555; }
.log-row__delta { font-size: 14px; font-weight: 700; flex-shrink: 0; }

.bottom-nav { display: none; }

@media (max-width: 900px) {
  .sidebar { display: none; }
  .mobile-header { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px; background: #1e2235; position: sticky; top: 0; z-index: 50; }
  .mobile-header__logo { color: #fff; font-size: 17px; font-weight: 800; }
  .mobile-menu-btn { background: none; border: none; color: #fff; font-size: 22px; cursor: pointer; padding: 4px 8px; }
  .mobile-tabs-dropdown { display: flex; flex-direction: column; background: #1e2235; border-bottom: 1px solid rgba(255,255,255,0.1); overflow: hidden; max-height: 0; transition: max-height 0.3s ease; position: sticky; top: 52px; z-index: 40; }
  .mobile-tabs-dropdown.open { max-height: 500px; }
  .mobile-tab-item { display: flex; align-items: center; gap: 10px; padding: 13px 20px; border: none; background: none; color: #aab; font-size: 15px; font-weight: 500; cursor: pointer; text-align: left; font-family: inherit; border-bottom: 1px solid rgba(255,255,255,0.05); }
  .mobile-tab-item.active { color: #4f7ef7; background: rgba(79,126,247,0.1); }
  .mobile-tab-item--logout { color: #f87171; }
  .admin-main { padding: 16px; padding-bottom: 80px; }
  .chat-layout { height: auto; flex-direction: column; }
  .chat-sidebar { width: 100%; border-radius: 12px; margin-bottom: 12px; border-right: none; max-height: 220px; }
  .chat-window { border-radius: 12px; height: 400px; }
  .chat-window__header { display: flex; align-items: center; gap: 12px; padding: 12px 16px; border-bottom: 1px solid #eee; }
  .chat-back-btn { background: none; border: none; color: #4f7ef7; font-size: 14px; font-weight: 600; cursor: pointer; padding: 0; font-family: inherit; }
  .chat-window__name { font-size: 14px; font-weight: 700; }
  .bottom-nav { display: flex; position: fixed; bottom: 0; left: 0; right: 0; background: #1e2235; border-top: 1px solid rgba(255,255,255,0.1); padding: 6px 0 env(safe-area-inset-bottom, 6px); z-index: 100; }
  .bottom-nav__item { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 2px; padding: 6px 2px; background: none; border: none; cursor: pointer; color: #555; font-family: inherit; transition: color 0.15s; }
  .bottom-nav__item.active { color: #4f7ef7; }
  .bottom-nav__icon { font-size: 20px; }
  .bottom-nav__label { font-size: 9px; font-weight: 600; color: inherit; }
}
</style>
