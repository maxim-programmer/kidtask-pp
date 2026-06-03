<template>
  <ParentLayout>
    <h1 class="page-title">Мои дети</h1>

    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else>
      <div v-for="child in children" :key="child.child_id" class="child-card">
        <div class="child-card__avatar" :class="child.avatar_url ? 'avatar--photo' : avatarClass(child)" @click="openAvatarModal(child)">
          <img v-if="child.avatar_url" :src="child.avatar_url" class="avatar-img" />
          <span v-else>{{ child.name[0] }}</span>
          <div class="avatar-overlay">📷</div>
        </div>
        <div class="child-card__info">
          <div class="child-card__name">{{ child.name }}</div>
          <div class="child-card__meta">
            ⭐ {{ child.balance }}
            <span v-if="ageLabel(child.birthday)" class="age-badge">{{ ageLabel(child.birthday) }}</span>
            <span v-else class="age-badge age-badge--unknown">возраст не указан</span>
          </div>
        </div>
        <button class="icon-btn" @click="removeChild(child)" title="Удалить">🗑</button>
      </div>
      <button class="add-btn" @click="showModal = true">+ Добавить ребёнка</button>
    </div>

    <div v-if="showModal" class="modal-overlay" @click.self="closeAddModal">
      <div class="modal">
        <h3>Новый ребёнок</h3>
        <div class="field">
          <label>Имя</label>
          <input v-model="form.name" type="text" placeholder="Иван" />
        </div>
        <div class="field">
          <label>Логин</label>
          <input v-model="form.username" type="text" placeholder="ivan123" />
        </div>
        <div class="field">
          <label>Пароль</label>
          <input v-model="form.password" type="password" />
        </div>
        <div class="field">
          <label>День рождения</label>
          <input v-model="form.birthday" type="date" required />
        </div>
        <div v-if="error" class="error-msg">{{ error }}</div>
        <div class="modal-actions">
          <button class="btn-outline" @click="closeAddModal">Отмена</button>
          <button class="btn-primary" @click="addChild" :disabled="saving">
            {{ saving ? "Сохранение..." : "Добавить" }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="avatarModal.show" class="modal-overlay" @click.self="closeAvatarModal">
      <div class="modal modal--avatar">
        <h3>Фото {{ avatarModal.child && avatarModal.child.name }}</h3>
        <div class="avatar-preview-wrap">
          <div class="avatar-preview-circle" :class="avatarModal.preview || (avatarModal.child && avatarModal.child.avatar_url) ? 'avatar-preview-circle--img' : avatarClass(avatarModal.child)">
            <img v-if="avatarModal.preview || (avatarModal.child && avatarModal.child.avatar_url)"
              :src="avatarModal.preview || avatarModal.child.avatar_url"
              class="avatar-preview-img" />
            <span v-else class="avatar-preview-letter">{{ avatarModal.child && avatarModal.child.name && avatarModal.child.name[0] }}</span>
          </div>
          <button class="btn-pick" @click="$refs.fileInput.click()">
            {{ avatarModal.preview ? "🔄 Выбрать другое" : "📷 Выбрать фото" }}
          </button>
          <input ref="fileInput" type="file" accept="image/*" style="display:none" @change="onFileChange" />
        </div>
        <div v-if="avatarModal.error" class="error-msg">{{ avatarModal.error }}</div>
        <div class="modal-actions">
          <button class="btn-outline" @click="closeAvatarModal">Отмена</button>
          <button v-if="avatarModal.child && avatarModal.child.avatar_url" class="btn-danger" @click="removeAvatar" :disabled="avatarModal.saving">Удалить фото</button>
          <button class="btn-primary" @click="saveAvatar" :disabled="!avatarModal.preview || avatarModal.saving">
            {{ avatarModal.saving ? "Сохранение..." : "Сохранить" }}
          </button>
        </div>
      </div>
    </div>
  </ParentLayout>
</template>

<script>
import ParentLayout from "../../components/ParentLayout.vue"
import { useApi } from "../../composables/useApi"
import { ageLabel, calcAge } from "../../composables/useAge"
export default {
  name: "ParentChildren",
  components: { ParentLayout },
  data() {
    return {
      loading: true,
      children: [],
      showModal: false,
      saving: false,
      error: "",
      form: { name: "", username: "", password: "", birthday: "" },
      avatarModal: { show: false, child: null, preview: null, saving: false, error: "" }
    }
  },
  async mounted() { await this.load() },
  methods: {
    ageLabel,
    avatarClass(child) {
      if (!child) return "avatar--default"
      const age = calcAge(child.birthday)
      if (age === null) return "avatar--default"
      return age < 11 ? "avatar--junior" : "avatar--senior"
    },
    async load() {
      const { getChildren } = useApi()
      this.loading = true
      try { this.children = (await getChildren()).data.children || [] }
      finally { this.loading = false }
    },
    closeAddModal() {
      this.showModal = false
      this.error = ""
      this.form = { name: "", username: "", password: "", birthday: "" }
    },
    async addChild() {
      this.error = ""
      if (!this.form.name) { this.error = "Введите имя"; return }
      if (!this.form.username) { this.error = "Введите логин"; return }
      if (!this.form.password) { this.error = "Введите пароль"; return }
      if (!this.form.birthday) { this.error = "Укажите дату рождения"; return }
      this.saving = true
      const { createChild } = useApi()
      try {
        await createChild({
          name: this.form.name,
          username: this.form.username,
          password: this.form.password,
          birthday: new Date(this.form.birthday + "T00:00:00Z").toISOString()
        })
        this.closeAddModal()
        await this.load()
      } catch (e) { this.error = e.response?.data?.error?.message || "Ошибка" }
      finally { this.saving = false }
    },
    async removeChild(child) {
      if (!confirm("Удалить " + child.name + "?")) return
      const { deleteChild } = useApi()
      await deleteChild(child.child_id)
      await this.load()
    },
    openAvatarModal(child) {
      this.avatarModal = { show: true, child, preview: null, saving: false, error: "" }
    },
    closeAvatarModal() {
      this.avatarModal = { show: false, child: null, preview: null, saving: false, error: "" }
      if (this.$refs.fileInput) this.$refs.fileInput.value = ""
    },
    onFileChange(e) {
      const file = e.target.files[0]
      if (!file) return
      if (file.size > 5 * 1024 * 1024) { this.avatarModal.error = "Файл слишком большой (макс. 5 МБ)"; return }
      this.avatarModal.error = ""
      const reader = new FileReader()
      reader.onload = (ev) => {
        this.compressImage(ev.target.result, (compressed) => { this.avatarModal.preview = compressed })
      }
      reader.readAsDataURL(file)
    },
    compressImage(dataUrl, callback) {
      const img = new Image()
      img.onload = () => {
        const canvas = document.createElement("canvas")
        const MAX = 256
        let w = img.width, h = img.height
        if (w > h) { if (w > MAX) { h = Math.round(h * MAX / w); w = MAX } }
        else { if (h > MAX) { w = Math.round(w * MAX / h); h = MAX } }
        canvas.width = w
        canvas.height = h
        canvas.getContext("2d").drawImage(img, 0, 0, w, h)
        callback(canvas.toDataURL("image/jpeg", 0.85))
      }
      img.src = dataUrl
    },
    async saveAvatar() {
      if (!this.avatarModal.preview) return
      this.avatarModal.saving = true
      this.avatarModal.error = ""
      const { setChildAvatar } = useApi()
      try {
        await setChildAvatar(this.avatarModal.child.child_id, this.avatarModal.preview)
        await this.load()
        this.closeAvatarModal()
      } catch (e) { this.avatarModal.error = e.response?.data?.error?.message || "Ошибка" }
      finally { this.avatarModal.saving = false }
    },
    async removeAvatar() {
      this.avatarModal.saving = true
      this.avatarModal.error = ""
      const { setChildAvatar } = useApi()
      try {
        await setChildAvatar(this.avatarModal.child.child_id, "")
        await this.load()
        this.closeAvatarModal()
      } catch (e) { this.avatarModal.error = e.response?.data?.error?.message || "Ошибка" }
      finally { this.avatarModal.saving = false }
    }
  }
}
</script>

<style scoped>
.page-title { font-size: 24px; font-weight: 700; margin-bottom: 20px; }
.loading { text-align: center; padding: 60px; color: #888; }
.child-card { background: #fff; border-radius: 16px; padding: 16px; display: flex; align-items: center; gap: 14px; margin-bottom: 12px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); }
.child-card__avatar { width: 48px; height: 48px; border-radius: 50%; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 20px; font-weight: 700; flex-shrink: 0; cursor: pointer; position: relative; overflow: hidden; }
.child-card__avatar:hover .avatar-overlay { opacity: 1; }
.avatar--junior { background: #f97316; }
.avatar--senior { background: #4f7ef7; }
.avatar--default { background: #9ca3af; }
.avatar--photo { background: transparent; }
.avatar-img { width: 100%; height: 100%; object-fit: cover; }
.avatar-overlay { position: absolute; inset: 0; background: rgba(0,0,0,0.45); display: flex; align-items: center; justify-content: center; font-size: 18px; opacity: 0; transition: opacity 0.2s; }
.child-card__info { flex: 1; }
.child-card__name { font-size: 16px; font-weight: 700; }
.child-card__meta { font-size: 13px; color: #888; margin-top: 4px; display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.age-badge { background: #f0f4ff; color: #4f7ef7; border-radius: 8px; padding: 2px 8px; font-size: 12px; font-weight: 600; }
.age-badge--unknown { background: #f5f5f5; color: #aaa; }
.icon-btn { width: 34px; height: 34px; border-radius: 8px; border: none; background: #fee2e2; font-size: 16px; cursor: pointer; color: #dc2626; display: flex; align-items: center; justify-content: center; }
.add-btn { width: 100%; padding: 14px; background: #f0f4ff; border: 2px dashed #4f7ef7; border-radius: 16px; color: #4f7ef7; font-size: 16px; font-weight: 600; cursor: pointer; margin-top: 8px; font-family: inherit; }
.add-btn:hover { background: #e8eeff; }
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 300; }
.modal { background: #fff; border-radius: 16px; padding: 24px; width: 90%; max-width: 360px; }
.modal--avatar { max-width: 310px; }
.modal h3 { font-size: 20px; font-weight: 700; text-align: center; margin-bottom: 20px; }
.avatar-preview-wrap { display: flex; flex-direction: column; align-items: center; gap: 14px; margin-bottom: 20px; }
.avatar-preview-circle { width: 100px; height: 100px; border-radius: 50%; display: flex; align-items: center; justify-content: center; overflow: hidden; flex-shrink: 0; }
.avatar-preview-circle--img { background: #e8e8e8; }
.avatar--junior.avatar-preview-circle { background: #f97316; }
.avatar--senior.avatar-preview-circle { background: #4f7ef7; }
.avatar--default.avatar-preview-circle { background: #9ca3af; }
.avatar-preview-img { width: 100%; height: 100%; object-fit: cover; }
.avatar-preview-letter { font-size: 42px; font-weight: 800; color: #fff; }
.btn-pick { padding: 9px 18px; background: #f0f4ff; border: 1.5px solid #4f7ef7; border-radius: 10px; color: #4f7ef7; font-size: 14px; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-pick:hover { background: #e8eeff; }
.field { margin-bottom: 12px; }
.field label { display: block; font-size: 13px; color: #666; margin-bottom: 4px; font-weight: 500; }
.field input { width: 100%; padding: 10px 12px; border: 1.5px solid #ddd; border-radius: 8px; font-size: 15px; outline: none; box-sizing: border-box; font-family: inherit; }
.field input:focus { border-color: #4f7ef7; }
.error-msg { color: #e53e3e; font-size: 13px; margin-bottom: 10px; text-align: center; }
.modal-actions { display: flex; gap: 8px; margin-top: 4px; }
.btn-outline { flex: 1; padding: 11px; border: 1.5px solid #ddd; background: #fff; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; color: #666; font-family: inherit; }
.btn-primary { flex: 1; padding: 11px; background: #4f7ef7; color: #fff; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-primary:hover:not(:disabled) { background: #3a6be0; }
.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-danger { flex: 1; padding: 11px; background: #fee2e2; color: #dc2626; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-danger:hover:not(:disabled) { background: #fecaca; }
.btn-danger:disabled { opacity: 0.5; cursor: not-allowed; }
</style>