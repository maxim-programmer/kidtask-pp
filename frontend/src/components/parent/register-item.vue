<script setup lang="ts">
  import check_reg_credentials from '../../functions/parent_actions'
</script>

<style>
  @import '../../assets/parent.css';
</style>

<template>
  <div class="form-wrapper">
    <form class="login-form">
      <label>Имя</label>
      <input id="username-value-reg" />
      <label>Почта</label>
      <input id="email-value-reg" />
      <label>Пароль</label>
      <input id="pass-value-reg" />
      <!--<router-link></router-link>-->
      <button @click="check_reg_credentials">Регистрация</button>
    </form>
  </div>
</template>
