<template>
  <div class="page">
    <header class="header">
      <div class="header__inner">
        <div class="logo">
          kid<span class="logo__accent">TASK</span>
        </div>
        <button class="btn btn--outline" disabled>Войти</button>
      </div>
    </header>

    <main>
      <section class="hero">
        <div class="hero__inner">
          <div class="hero__content">
            <h1 class="hero__title">
              kid<span class="hero__title-accent">TASK</span>
            </h1>
            <p class="hero__subtitle">
              Сервис геймификации домашних обязанностей и система финансового воспитания для детей
            </p>
            <p class="hero__description">
              Помогаем родителям мотивировать детей выполнять домашние обязанности, а детям — копить на свои желания через игровую валюту.
            </p>
            <button class="btn btn--primary btn--large" disabled>
              Присоединиться
            </button>
          </div>
          <div class="hero__visual">
            <div class="card-preview">
              <div class="card-preview__header">
                <span class="card-preview__name">Петя</span>
                <span class="card-preview__balance">⭐ 123</span>
              </div>
              <div class="card-preview__task">
                <span class="card-preview__task-title">Убрать комнату</span>
                <span class="card-preview__task-reward">+10 ⭐</span>
              </div>
              <div class="card-preview__task">
                <span class="card-preview__task-title">Помыть посуду</span>
                <span class="card-preview__task-reward">+5 ⭐</span>
              </div>
              <div class="card-preview__wish">
                <span class="card-preview__wish-title">🎮 Nintendo Switch</span>
                <div class="card-preview__progress-bar">
                  <div class="card-preview__progress-fill" style="width: 60%"></div>
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
              <div class="stat-card__value">
                <span v-if="loading" class="stat-card__loading">...</span>
                <span v-else>{{ stats.total }}</span>
              </div>
              <div class="stat-card__label">Пользователей</div>
            </div>
            <div class="stat-card">
              <div class="stat-card__value">
                <span v-if="loading" class="stat-card__loading">...</span>
                <span v-else>{{ stats.parents }}</span>
              </div>
              <div class="stat-card__label">Родителей</div>
            </div>
            <div class="stat-card">
              <div class="stat-card__value">
                <span v-if="loading" class="stat-card__loading">...</span>
                <span v-else>{{ stats.children }}</span>
              </div>
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

      <section class="cta">
        <div class="cta__inner">
          <h2 class="cta__title">Готовы начать?</h2>
          <p class="cta__text">Присоединяйтесь к KidTask и сделайте домашние обязанности интересными</p>
          <button class="btn btn--primary btn--large" disabled>
            Зарегистрироваться
          </button>
        </div>
      </section>
    </main>

    <footer class="footer">
      <div class="footer__inner">
        <span class="logo">kid<span class="logo__accent">TASK</span></span>
        <p class="footer__copy">© 2026 KidTask</p>
      </div>
    </footer>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'HomePage',
  data() {
    return {
      loading: true,
      stats: {
        total: 0,
        parents: 0,
        children: 0,
      }
    }
  },
  async mounted() {
    try {
      const res = await axios.get('/api/stats')
      this.stats = res.data
    } catch {
      this.stats = { total: 0, parents: 0, children: 0 }
    } finally {
      this.loading = false
    }
  }
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: #ffffff;
  border-bottom: 1px solid #e8e8e8;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header__inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 16px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a1a;
  letter-spacing: -0.5px;
}

.logo__accent {
  color: #4f7ef7;
}

.btn {
  padding: 10px 24px;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn--primary {
  background: #4f7ef7;
  color: #ffffff;
}

.btn--primary:not(:disabled):hover {
  background: #3a6be0;
}

.btn--outline {
  background: transparent;
  color: #4f7ef7;
  border: 2px solid #4f7ef7;
}

.btn--large {
  padding: 14px 36px;
  font-size: 17px;
  border-radius: 10px;
}

.hero {
  background: #ffffff;
  padding: 80px 0;
}

.hero__inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 60px;
  align-items: center;
}

.hero__title {
  font-size: 56px;
  font-weight: 800;
  color: #1a1a1a;
  margin-bottom: 16px;
  letter-spacing: -1px;
}

.hero__title-accent {
  color: #4f7ef7;
}

.hero__subtitle {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
  line-height: 1.4;
}

.hero__description {
  font-size: 16px;
  color: #666;
  line-height: 1.7;
  margin-bottom: 32px;
}

.hero__visual {
  display: flex;
  justify-content: center;
}

.card-preview {
  background: #ffffff;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.1);
  width: 320px;
  border: 1px solid #e8e8e8;
}

.card-preview__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-preview__name {
  font-size: 18px;
  font-weight: 700;
  color: #1a1a1a;
}

.card-preview__balance {
  font-size: 16px;
  font-weight: 600;
  color: #4f7ef7;
  background: #eef2ff;
  padding: 4px 12px;
  border-radius: 20px;
}

.card-preview__task {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: #f5f7ff;
  border-radius: 8px;
  margin-bottom: 8px;
}

.card-preview__task-title {
  font-size: 14px;
  color: #333;
}

.card-preview__task-reward {
  font-size: 13px;
  font-weight: 600;
  color: #4f7ef7;
}

.card-preview__wish {
  margin-top: 16px;
  padding: 12px;
  background: #fff8f0;
  border-radius: 8px;
  border: 1px solid #fde8c8;
}

.card-preview__wish-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  display: block;
  margin-bottom: 8px;
}

.card-preview__progress-bar {
  height: 6px;
  background: #e8e8e8;
  border-radius: 3px;
  margin-bottom: 6px;
}

.card-preview__progress-fill {
  height: 100%;
  background: #4f7ef7;
  border-radius: 3px;
}

.card-preview__wish-info {
  font-size: 12px;
  color: #888;
}

.stats {
  padding: 80px 0;
  background: #f5f7ff;
}

.stats__inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  text-align: center;
}

.stats__title {
  font-size: 36px;
  font-weight: 700;
  color: #1a1a1a;
  margin-bottom: 48px;
}

.stats__grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

.stat-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 40px 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}

.stat-card__value {
  font-size: 52px;
  font-weight: 800;
  color: #4f7ef7;
  margin-bottom: 8px;
}

.stat-card__loading {
  color: #aaa;
}

.stat-card__label {
  font-size: 16px;
  color: #666;
  font-weight: 500;
}

.features {
  padding: 80px 0;
  background: #ffffff;
}

.features__inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  text-align: center;
}

.features__title {
  font-size: 36px;
  font-weight: 700;
  color: #1a1a1a;
  margin-bottom: 48px;
}

.features__grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

.feature-card {
  padding: 36px 24px;
  background: #f5f7ff;
  border-radius: 16px;
  text-align: center;
}

.feature-card__icon {
  font-size: 40px;
  margin-bottom: 16px;
}

.feature-card__title {
  font-size: 20px;
  font-weight: 700;
  color: #1a1a1a;
  margin-bottom: 12px;
}

.feature-card__text {
  font-size: 15px;
  color: #666;
  line-height: 1.6;
}

.cta {
  padding: 80px 0;
  background: #4f7ef7;
}

.cta__inner {
  max-width: 600px;
  margin: 0 auto;
  padding: 0 24px;
  text-align: center;
}

.cta__title {
  font-size: 36px;
  font-weight: 700;
  color: #ffffff;
  margin-bottom: 16px;
}

.cta__text {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.85);
  margin-bottom: 32px;
  line-height: 1.6;
}

.cta .btn--primary {
  background: #ffffff;
  color: #4f7ef7;
}

.footer {
  background: #1a1a1a;
  padding: 32px 0;
  margin-top: auto;
}

.footer__inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.footer .logo {
  color: #ffffff;
}

.footer__copy {
  font-size: 14px;
  color: #888;
}

@media (max-width: 768px) {
  .hero__inner {
    grid-template-columns: 1fr;
    gap: 40px;
  }

  .hero__title {
    font-size: 40px;
  }

  .hero__visual {
    order: -1;
  }

  .stats__grid,
  .features__grid {
    grid-template-columns: 1fr;
  }
}
</style>