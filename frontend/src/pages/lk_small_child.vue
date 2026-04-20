<script setup lang="ts">
  import Header from '../components/small_child/header.vue'
  import { getChildUserData, test_child_user2 } from '../functions/child_actions'

  //const user = getChildUserData()
  const props = defineProps<{
    username?: string,
    avatar_style?: string,
    avatar_path?: string
  }>()
</script>

<style scoped>
  @import '../assets/ch710.css';
</style>

<template>
  <Header :avatar_style="test_child_user2.avatar_style" :avatar_path="test_child_user2.avatar_path"/>

  <div class="greeting">
    <h2>Привет, {{ test_child_user2.username }} </h2>
  </div>

  <div class="tabs">
  </div>

  <div id="main-body"></div>
</template>

