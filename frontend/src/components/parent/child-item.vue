<script setup lang="ts">
  defineProps<{
    username: string
  }>()
</script>

<style>
  @import '../../assets/parent.css';
</style>

<template>
  <div>
    <div class="child-avatar">
      <img src="{{ avatar_path }}">
    </div>
    <div>
      <<h1>{{ username }}</h1>
    </div>
  </div>
</template>
