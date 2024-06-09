<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import MainHander from '@/components/layout/MainHeader.vue'
</script>

<template>
  <el-container>
    <el-header style="padding: 0">
      <MainHander />
    </el-header>
    <el-main style="padding: 10px">
      <RouterView />
    </el-main>
  </el-container>
</template>
