<script setup lang="ts">
  import { check_login_credentials } from '../functions/child_actions'
</script>

<style>
  @import '../assets/main.css';
</style>

<template>
  <div class="form-wrapper">
    <form class="login-form">
      <label>Почта</label>
      <input id="email-value" />
      <label>Пароль</label>
      <input id="pass-value" />
      <a>Забыли пароль?</a>
      <!--<router-link></router-link>-->
      <button @click="check_login_credentials">Вход</button>
    </form>
  </div>
</template>
