<script setup lang="ts">
  import { add_task } from '../../functions/parent-actions'
</script>

<style>
  @import '../../assets/parent.css';
</style>

<template>
  <div class="form-wrapper">
    <form class="new-item-form">
      <label>Название</label>
      <input id="title-value" />
      <label>Описание</label>
      <input id="description-value" />
      <label>Цена</label>
      <input id="price-value" />
      <button @click="add_task">Добавить</button>
    </form>
  </div>
</template>
