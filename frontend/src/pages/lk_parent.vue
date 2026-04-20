<script setup lang="ts">
  import Header from '../components/parent/header.vue'
  import { getGreeting } from '../functions/misc'
  import { getUserData, test_user } from '../functions/parent_actions'

  //const user = getUserData()
  const props = defineProps<{
    username?: string,
    avatar_style?: string,
    avatar_path?: string,
    greeting?: string
  }>()
  const user_greeting = getGreeting()
</script>

<style>
  @import '../assets/parent.css';
</style>

<template>
  <Header :avatar_style="test_user.avatar_style" :avatar_path="test_user.avatar_path" />

  <div class="greeting">
    <h2>{{ user_greeting }}, {{ test_user.username }} </h2>
  </div>

  <div class="tabs">
  </div>

  <div id="main-body"></div>

</template>


