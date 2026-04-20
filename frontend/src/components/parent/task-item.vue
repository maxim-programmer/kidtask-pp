<script setup lang="ts">
  import { mark_task_done, mark_task_fix, remove_task } from '../../functions/parent_actions'

  defineProps<{
    title: string
    description: string
  }>()
</script>

<style>
  @import '../../assets/parent.css';
</style>

<template>
  <divc class="task-card">
    <h1>{{ title }}</h1>
    <p>{{ description }}</p>
    <div class="card-action-buttons">
      <button @click="mark_task_done">done</button>
      <button @click="mark_task_fix">fix</button>
      <button @click="remove_task">remove</button>
    </div>
  </div>
</template>
