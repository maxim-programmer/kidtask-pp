<script setup lang="ts">
</script>

<style>
  @import '../assets/ch1114.css';
</style>

<template>
  <div class="form-wrapper">
    <form class="new-item-form">
      <label>Название</label>
      <input id="title-value" />
      <label>Описание</label>
      <input id="description-value" />
      <button @click="new_task">Добавить</button>
    </form>
  </div>
</template>
