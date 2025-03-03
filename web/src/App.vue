<script setup lang="ts">
import { RouterView } from 'vue-router'
import { NConfigProvider,NDialogProvider } from 'naive-ui'

</script>

<template>
  <NConfigProvider>
    <NDialogProvider>
      <RouterView />
    </NDialogProvider>
  </NConfigProvider>
</template>

<style scoped>

</style>
