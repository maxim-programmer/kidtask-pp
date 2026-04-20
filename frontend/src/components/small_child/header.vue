<script setup lang="ts">
  const props = defineProps<{
    avatar_style?: string,
    avatar_path?: string,
    age_group?: string
  }>()
</script>

<style>
  @import '../../assets/ch710.css';
</style>

<template>
  <div class="header">
    <div class="app-title">
      <router-link to="/">
        <span>kid</span>TASK
      </router-link>
    </div>
    <div class="avatar">
      <img :id="props.avatar_style" :src="props.avatar_path" />
    </div>
  </div>
</template>
