<template>
  <router-view />
</template>

<script>
export default { name: 'App' }
</script>

<style>
*, *::before, *::after {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  -webkit-tap-highlight-color: transparent;
}

html {
  height: -webkit-fill-available;
  font-size: 16px;
  scroll-behavior: smooth;
}

body {
  font-family: 'Segoe UI', -apple-system, BlinkMacSystemFont, 'Helvetica Neue', Arial, sans-serif;
  background: #f5f5f5;
  color: #1a1a1a;
  min-height: 100vh;
  min-height: -webkit-fill-available;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  overscroll-behavior-y: none;
}

a {
  text-decoration: none;
  color: inherit;
}

button {
  cursor: pointer;
  font-family: inherit;
  border: none;
  outline: none;
  -webkit-appearance: none;
}

input, textarea, select {
  font-family: inherit;
  outline: none;
  -webkit-appearance: none;
  border-radius: 0;
}

input[type="date"] {
  -webkit-appearance: textfield;
}

img {
  max-width: 100%;
  display: block;
}

::-webkit-scrollbar {
  width: 4px;
  height: 4px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.15);
  border-radius: 2px;
}

@media (max-width: 768px) {
  ::-webkit-scrollbar {
    display: none;
  }
  * {
    scrollbar-width: none;
  }
}

.router-link-active,
.router-link-exact-active {
  /* handled per-component */
}
</style>
