<script setup lang="ts">
  const props = defineProps<{
    avatar_style?: string,
    avatar_path?: string
  }>()
</script>

<style>
  @import '../../assets/parent.css';
</style>

<template>
  <div class="header">
    <router-link to="/">
      <div class="app-title">
        <span>kid</span>TASK
      </div>
    </router-link>
    <div class="avatar">
      <img :id="props.avatar_style" :src="props.avatar_path" />
    </div>
  </div>
</template>
