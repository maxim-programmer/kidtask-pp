<template>
  <div class="page">
    <header class="header">
      <div class="header__inner">
        <div class="logo">kid<span class="logo__accent">TASK</span></div>
        <nav class="header__nav">
          <router-link to="/login" class="btn btn--outline">Войти</router-link>
          <router-link to="/register" class="btn btn--primary header__cta">Начать</router-link>
        </nav>
        <button class="burger" @click="mobileMenu = !mobileMenu" :class="{ open: mobileMenu }">
          <span></span><span></span><span></span>
        </button>
      </div>
      <div class="mobile-menu" :class="{ 'mobile-menu--open': mobileMenu }">
        <router-link to="/login" class="mobile-menu__link" @click="mobileMenu = false">Войти</router-link>
        <router-link to="/register" class="mobile-menu__link mobile-menu__link--accent" @click="mobileMenu = false">Зарегистрироваться</router-link>
        <router-link to="/admin/login" class="mobile-menu__link mobile-menu__link--muted" @click="mobileMenu = false">Для администраторов</router-link>
      </div>
    </header>

    <main>
      <section class="hero">
        <div class="hero__inner">
          <div class="hero__content">
            <div class="hero__badge">🏆 Финансовое воспитание</div>
            <h1 class="hero__title">kid<span class="hero__title-accent">TASK</span></h1>
            <p class="hero__subtitle">Превращаем домашние обязанности в игру и учим детей копить на мечту</p>
            <p class="hero__description">Родитель создаёт задания и назначает монеты. Ребёнок выполняет и копит на желаемое из вишлиста. Прозрачно, честно и весело.</p>
            <div class="hero__actions">
              <router-link to="/register" class="btn btn--primary btn--large">Присоединиться бесплатно</router-link>
              <router-link to="/login" class="btn btn--ghost btn--large">Уже есть аккаунт</router-link>
            </div>
          </div>
          <div class="hero__visual">
            <div class="card-preview">
              <div class="card-preview__header">
                <div class="card-preview__user">
                  <div class="card-preview__avatar">П</div>
                  <span class="card-preview__name">Петя</span>
                </div>
                <span class="card-preview__balance">⭐ 123</span>
              </div>
              <div class="card-preview__tasks">
                <div class="card-preview__task">
                  <span class="card-preview__task-title">🧹 Убрать комнату</span>
                  <span class="card-preview__task-reward">+10 ⭐</span>
                </div>
                <div class="card-preview__task">
                  <span class="card-preview__task-title">🍽 Помыть посуду</span>
                  <span class="card-preview__task-reward">+5 ⭐</span>
                </div>
              </div>
              <div class="card-preview__wish">
                <div class="card-preview__wish-header">
                  <span class="card-preview__wish-title">🎮 Nintendo Switch</span>
                  <span class="card-preview__wish-pct">62%</span>
                </div>
                <div class="card-preview__progress-bar">
                  <div class="card-preview__progress-fill" style="width: 62%"></div>
                </div>
                <span class="card-preview__wish-info">123 / 200 ⭐</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="stats">
        <div class="stats__inner">
          <h2 class="stats__title">Уже с нами</h2>
          <div class="stats__grid">
            <div class="stat-card">
              <div class="stat-card__value">{{ loading ? '...' : stats.total }}</div>
              <div class="stat-card__label">Пользователей</div>
            </div>
            <div class="stat-card">
              <div class="stat-card__value">{{ loading ? '...' : stats.parents }}</div>
              <div class="stat-card__label">Родителей</div>
            </div>
            <div class="stat-card">
              <div class="stat-card__value">{{ loading ? '...' : stats.children }}</div>
              <div class="stat-card__label">Детей</div>
            </div>
          </div>
        </div>
      </section>

      <section class="features">
        <div class="features__inner">
          <h2 class="features__title">Как это работает</h2>
          <div class="features__grid">
            <div class="feature-card">
              <div class="feature-card__icon">📋</div>
              <h3 class="feature-card__title">Задания</h3>
              <p class="feature-card__text">Родитель создаёт задания и назначает награду в монетах. Ребёнок выполняет и отмечает.</p>
            </div>
            <div class="feature-card">
              <div class="feature-card__icon">⭐</div>
              <h3 class="feature-card__title">Монеты</h3>
              <p class="feature-card__text">После подтверждения монеты начисляются на баланс ребёнка. Прозрачно и честно.</p>
            </div>
            <div class="feature-card">
              <div class="feature-card__icon">🎁</div>
              <h3 class="feature-card__title">Вишлист</h3>
              <p class="feature-card__text">Ребёнок добавляет желаемые вещи, родитель оценивает их в монетах. Копим вместе.</p>
            </div>
          </div>
        </div>
      </section>

      <section class="roles">
        <div class="roles__inner">
          <h2 class="roles__title">Для каждого своё</h2>
          <div class="roles__grid">
            <div class="role-card role-card--parent">
              <div class="role-card__icon">👨‍👩‍👧</div>
              <h3 class="role-card__title">Родителям</h3>
              <ul class="role-card__list">
                <li>Создавайте задания для детей</li>
                <li>Подтверждайте выполнение</li>
                <li>Оценивайте желания в вишлисте</li>
                <li>Следите за прогрессом</li>
              </ul>
              <router-link to="/register" class="role-card__btn">Начать →</router-link>
            </div>
            <div class="role-card role-card--child">
              <div class="role-card__icon">🧒</div>
              <h3 class="role-card__title">Детям</h3>
              <ul class="role-card__list">
                <li>Выполняйте задания и зарабатывайте</li>
                <li>Копи монеты на мечту</li>
                <li>Следи за прогрессом в вишлисте</li>
                <li>Общайся с родителями в чате</li>
              </ul>
              <router-link to="/login" class="role-card__btn role-card__btn--child">Войти →</router-link>
            </div>
          </div>
        </div>
      </section>

      <section class="cta">
        <div class="cta__inner">
          <h2 class="cta__title">Готовы начать?</h2>
          <p class="cta__text">Присоединяйтесь к KidTask и сделайте домашние обязанности интересными</p>
          <router-link to="/register" class="btn btn--white btn--large">Зарегистрироваться бесплатно</router-link>
        </div>
      </section>
    </main>

    <footer class="footer">
      <div class="footer__inner">
        <span class="logo logo--white">kid<span class="logo__accent">TASK</span></span>
        <p class="footer__copy">© 2026 KidTask</p>
        <router-link to="/admin/login" class="admin-link">Для администраторов</router-link>
      </div>
    </footer>
  </div>
</template>

<script>
import { useApi } from '../composables/useApi'
export default {
  name: 'HomePage',
  data() { return { loading: true, stats: { total: 0, parents: 0, children: 0 }, mobileMenu: false } },
  async mounted() {
    const { getStats } = useApi()
    try {
      const res = await getStats()
      this.stats = res.data
    } catch { this.stats = { total: 0, parents: 0, children: 0 } }
    finally { this.loading = false }
  }
}
</script>

<style scoped>
.page { min-height: 100vh; display: flex; flex-direction: column; }

.header { background: #fff; border-bottom: 1px solid #e8e8e8; position: sticky; top: 0; z-index: 100; }
.header__inner { max-width: 1200px; margin: 0 auto; padding: 14px 24px; display: flex; align-items: center; justify-content: space-between; }
.logo { font-size: 24px; font-weight: 800; color: #1a1a1a; }
.logo__accent { color: #4f7ef7; }
.logo--white { color: #fff; }

.header__nav { display: flex; gap: 10px; align-items: center; }
.btn { padding: 10px 20px; border-radius: 8px; font-size: 15px; font-weight: 600; border: none; transition: all 0.2s; display: inline-block; text-align: center; cursor: pointer; }
.btn--primary { background: #4f7ef7; color: #fff; }
.btn--primary:hover { background: #3a6be0; }
.btn--outline { background: transparent; color: #4f7ef7; border: 2px solid #4f7ef7; }
.btn--outline:hover { background: #f0f4ff; }
.btn--ghost { background: transparent; color: #555; border: 2px solid #ddd; }
.btn--ghost:hover { background: #f5f5f5; }
.btn--white { background: #fff; color: #4f7ef7; }
.btn--white:hover { background: #f0f4ff; }
.btn--large { padding: 14px 28px; font-size: 16px; border-radius: 10px; }

.burger { display: none; flex-direction: column; gap: 5px; background: none; border: none; cursor: pointer; padding: 4px; }
.burger span { display: block; width: 22px; height: 2px; background: #1a1a1a; border-radius: 2px; transition: all 0.25s; }
.burger.open span:nth-child(1) { transform: translateY(7px) rotate(45deg); }
.burger.open span:nth-child(2) { opacity: 0; }
.burger.open span:nth-child(3) { transform: translateY(-7px) rotate(-45deg); }

.mobile-menu { display: none; flex-direction: column; background: #fff; border-top: 1px solid #f0f0f0; overflow: hidden; max-height: 0; transition: max-height 0.3s ease; }
.mobile-menu--open { max-height: 300px; }
.mobile-menu__link { padding: 16px 24px; font-size: 16px; font-weight: 600; color: #333; border-bottom: 1px solid #f5f5f5; }
.mobile-menu__link--accent { color: #4f7ef7; }
.mobile-menu__link--muted { color: #aaa; font-size: 14px; }

.hero { background: linear-gradient(135deg, #f0f4ff 0%, #fff 60%); padding: 80px 0; }
.hero__inner { max-width: 1200px; margin: 0 auto; padding: 0 24px; display: grid; grid-template-columns: 1fr 1fr; gap: 60px; align-items: center; }
.hero__badge { display: inline-block; background: #eef2ff; color: #4f7ef7; padding: 6px 14px; border-radius: 20px; font-size: 13px; font-weight: 600; margin-bottom: 16px; }
.hero__title { font-size: 56px; font-weight: 900; color: #1a1a1a; margin-bottom: 14px; line-height: 1; }
.hero__title-accent { color: #4f7ef7; }
.hero__subtitle { font-size: 20px; font-weight: 600; color: #333; margin-bottom: 14px; line-height: 1.4; }
.hero__description { font-size: 16px; color: #666; line-height: 1.7; margin-bottom: 28px; }
.hero__actions { display: flex; gap: 12px; flex-wrap: wrap; }
.header__cta { display: inline-block; }

.card-preview { background: #fff; border-radius: 20px; padding: 22px; box-shadow: 0 8px 40px rgba(79,126,247,0.15); width: 100%; max-width: 320px; border: 1px solid #e8e8e8; }
.card-preview__header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.card-preview__user { display: flex; align-items: center; gap: 10px; }
.card-preview__avatar { width: 36px; height: 36px; border-radius: 50%; background: #f97316; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 16px; font-weight: 700; }
.card-preview__name { font-size: 16px; font-weight: 700; }
.card-preview__balance { font-size: 15px; font-weight: 700; color: #4f7ef7; background: #eef2ff; padding: 5px 12px; border-radius: 20px; }
.card-preview__tasks { margin-bottom: 14px; display: flex; flex-direction: column; gap: 8px; }
.card-preview__task { display: flex; justify-content: space-between; align-items: center; padding: 10px 12px; background: #f5f7ff; border-radius: 10px; }
.card-preview__task-title { font-size: 13px; color: #333; }
.card-preview__task-reward { font-size: 13px; font-weight: 700; color: #4f7ef7; }
.card-preview__wish { padding: 12px; background: #fff8f0; border-radius: 12px; border: 1px solid #fde8c8; }
.card-preview__wish-header { display: flex; justify-content: space-between; margin-bottom: 8px; }
.card-preview__wish-title { font-size: 13px; font-weight: 600; }
.card-preview__wish-pct { font-size: 12px; color: #f97316; font-weight: 700; }
.card-preview__progress-bar { height: 6px; background: #fde8c8; border-radius: 3px; margin-bottom: 6px; }
.card-preview__progress-fill { height: 100%; background: #f97316; border-radius: 3px; }
.card-preview__wish-info { font-size: 11px; color: #888; }

.stats { padding: 70px 0; background: #1a1a2e; }
.stats__inner { max-width: 1200px; margin: 0 auto; padding: 0 24px; text-align: center; }
.stats__title { font-size: 32px; font-weight: 700; color: #fff; margin-bottom: 40px; }
.stats__grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px; }
.stat-card { background: rgba(255,255,255,0.07); border-radius: 16px; padding: 32px 20px; border: 1px solid rgba(255,255,255,0.1); }
.stat-card__value { font-size: 48px; font-weight: 900; color: #4f7ef7; margin-bottom: 6px; }
.stat-card__label { font-size: 15px; color: rgba(255,255,255,0.6); font-weight: 500; }

.features { padding: 80px 0; background: #fff; }
.features__inner { max-width: 1200px; margin: 0 auto; padding: 0 24px; text-align: center; }
.features__title { font-size: 32px; font-weight: 700; margin-bottom: 40px; }
.features__grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px; }
.feature-card { padding: 32px 24px; background: #f5f7ff; border-radius: 20px; text-align: left; }
.feature-card__icon { font-size: 36px; margin-bottom: 14px; }
.feature-card__title { font-size: 18px; font-weight: 700; margin-bottom: 10px; }
.feature-card__text { font-size: 14px; color: #666; line-height: 1.6; }

.roles { padding: 80px 0; background: #f5f7ff; }
.roles__inner { max-width: 1200px; margin: 0 auto; padding: 0 24px; text-align: center; }
.roles__title { font-size: 32px; font-weight: 700; margin-bottom: 40px; }
.roles__grid { display: grid; grid-template-columns: 1fr 1fr; gap: 24px; }
.role-card { background: #fff; border-radius: 20px; padding: 32px; text-align: left; box-shadow: 0 2px 12px rgba(0,0,0,0.06); }
.role-card--child { background: linear-gradient(135deg, #1a1a2e, #16213e); }
.role-card__icon { font-size: 40px; margin-bottom: 14px; }
.role-card__title { font-size: 22px; font-weight: 700; margin-bottom: 16px; }
.role-card--child .role-card__title { color: #fff; }
.role-card__list { list-style: none; padding: 0; margin: 0 0 24px; display: flex; flex-direction: column; gap: 10px; }
.role-card__list li { font-size: 14px; color: #555; padding-left: 20px; position: relative; }
.role-card__list li::before { content: '✓'; position: absolute; left: 0; color: #4f7ef7; font-weight: 700; }
.role-card--child .role-card__list li { color: rgba(255,255,255,0.75); }
.role-card--child .role-card__list li::before { color: #f97316; }
.role-card__btn { display: inline-block; padding: 12px 22px; background: #4f7ef7; color: #fff; border-radius: 10px; font-size: 14px; font-weight: 700; transition: background 0.2s; }
.role-card__btn:hover { background: #3a6be0; }
.role-card__btn--child { background: #f97316; }
.role-card__btn--child:hover { background: #ea580c; }

.cta { padding: 80px 0; background: linear-gradient(135deg, #4f7ef7, #7c3aed); }
.cta__inner { max-width: 600px; margin: 0 auto; padding: 0 24px; text-align: center; }
.cta__title { font-size: 36px; font-weight: 800; color: #fff; margin-bottom: 14px; }
.cta__text { font-size: 17px; color: rgba(255,255,255,0.85); margin-bottom: 28px; line-height: 1.6; }

.footer { background: #0d0d1a; padding: 32px 0; margin-top: auto; }
.footer__inner { max-width: 1200px; margin: 0 auto; padding: 0 24px; display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 12px; }
.footer__copy { font-size: 14px; color: #555; }
.admin-link { font-size: 13px; color: #444; text-decoration: none; transition: color 0.15s; }
.admin-link:hover { color: #4f7ef7; }

@media (max-width: 900px) {
  .hero__inner { grid-template-columns: 1fr; gap: 40px; }
  .hero__title { font-size: 44px; }
  .hero__visual { display: flex; justify-content: center; }
  .features__grid { grid-template-columns: 1fr 1fr; }
  .roles__grid { grid-template-columns: 1fr; }
  .stats__grid { grid-template-columns: repeat(3, 1fr); }
}

@media (max-width: 640px) {
  .header__nav { display: none; }
  .burger { display: flex; }
  .mobile-menu { display: flex; }

  .hero { padding: 40px 0; }
  .hero__inner { padding: 0 16px; gap: 28px; }
  .hero__badge { font-size: 12px; }
  .hero__title { font-size: 36px; }
  .hero__subtitle { font-size: 17px; }
  .hero__description { font-size: 14px; margin-bottom: 20px; }
  .hero__actions { flex-direction: column; gap: 10px; }
  .hero__actions .btn { width: 100%; text-align: center; padding: 14px; }
  .card-preview { max-width: 100%; }

  .stats { padding: 40px 0; }
  .stats__inner { padding: 0 16px; }
  .stats__title { font-size: 24px; margin-bottom: 24px; }
  .stats__grid { grid-template-columns: 1fr; gap: 12px; }
  .stat-card { padding: 20px; }
  .stat-card__value { font-size: 36px; }

  .features { padding: 40px 0; }
  .features__inner { padding: 0 16px; }
  .features__title { font-size: 24px; margin-bottom: 24px; }
  .features__grid { grid-template-columns: 1fr; gap: 12px; }
  .feature-card { padding: 20px; }

  .roles { padding: 40px 0; }
  .roles__inner { padding: 0 16px; }
  .roles__title { font-size: 24px; margin-bottom: 24px; }
  .role-card { padding: 22px; }

  .cta { padding: 48px 0; }
  .cta__inner { padding: 0 16px; }
  .cta__title { font-size: 26px; }
  .cta__text { font-size: 15px; }
  .cta .btn--large { width: 100%; text-align: center; }

  .footer__inner { padding: 0 16px; flex-direction: column; align-items: center; gap: 8px; text-align: center; }
}
</style>
