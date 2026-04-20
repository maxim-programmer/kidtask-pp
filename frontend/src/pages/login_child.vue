<script setup lang="ts">
  import Header from '../components/parent/header.vue'
  import LoginItem from '../components/child-login-item.vue'

  const avstyle = "no-avatar"
  const avpath = "../public/white-square.png"
</script>

<style>
  @import '../assets/main.css';
</style>

<template>
  <Header :avatar_style="avstyle" :avatar_path="avpath" />
  <div class="greeting">
    <h2>Вход</h2>
  </div>
  <div class="page-content">
    <LoginItem />
  </div>
</template>
