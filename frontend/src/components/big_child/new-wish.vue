<script setup lang="ts">
  import { add_wish } from '../../functions/child_actions'

  function handleSubmit() {
    const title = (document.getElementById('title-value') as HTMLInputElement).value
    const description = (document.getElementById('description-value') as HTMLInputElement).value
    add_wish(title, description)
  }
</script>

<style>
  @import '../../assets/ch1114.css';
</style>

<template>
  <div class="form-wrapper">
    <form class="new-item-form">
      <label>Название</label>
      <input id="title-value" />
      <label>Описание</label>
      <input id="description-value" />
      <button @click="handleSubmit">Добавить</button>
    </form>
  </div>
</template>