<template>
  <div id="app">
    <HomePage />
  </div>
</template>

<script>
import HomePage from './views/HomePage.vue'

export default {
  components: { HomePage }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background-color: #f5f5f5;
  color: #1a1a1a;
}
</style>