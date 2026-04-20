<script setup lang="ts">
  import { edit_wish, remove_wish } from '../../functions/parent_actions'

  defineProps<{
    title: string
    description: string
    price: number | string
  }>()
</script>

<style>
  @import '../../assets/parent.css';
</style>

<template>
  <div class="wish-card">
    <h1>{{ title }}</h1>
    <p>{{ description }}</p>
    <span>{{ price }}</span>
    <div class="card-action-buttons">
      <button @click="edit_wish">done</button>
      <button @click="remove_wish">remove</button>
    </div>
  </div>
</template>
