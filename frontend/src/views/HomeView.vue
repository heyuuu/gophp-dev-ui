<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <main class="main">
    <el-button type="primary" size="default" @click="router.push('run')">运行代码</el-button>
    <el-button type="primary" size="default" @click="router.push('test')">运行Case</el-button>
  </main>
</template>

<style scoped>
.main {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-content: center;
  flex-wrap: wrap;
}
</style>
