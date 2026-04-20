<script setup lang="ts">
  import Header from '../components/parent/header.vue'
  import RegisterItem from '../components/parent/register-item.vue'

  const avstyle = "no-avatar"
  const avpath = "../public/white-square.png"
</script>

<style>
  @import '../assets/main.css';
</style>

<template>
  <Header :avatar_style="avstyle" :avatar_path="avpath" />
  <div class="greeting">
    <h2>Регистрация</h2>
  </div>
  <div class="page-content">
    <RegisterItem />
  </div>
  <div class="page-link">
    <router-link to="/parent/login">Уже зарегистрированы? Войти</router-link>
  </div>
</template>
