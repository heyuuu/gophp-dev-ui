<template>
  <el-menu
    :default-active="activeIndex"
    class="el-menu-demo"
    mode="horizontal"
    @select="handleSelect"
  >
    <template v-for="item in menuItems">
      <!-- 折叠菜单项 -->
      <template v-if="item.children">
        <el-sub-menu :index="item.path" :key="item.path">
          <template #title>{{ item.name }}</template>
          <el-menu-item
            v-for="subItem in item.children"
            :index="subItem.path"
            :key="subItem.path"
            >{{ subItem.name }}</el-menu-item
          >
        </el-sub-menu>
      </template>
      <!-- 单独菜单项 -->
      <template v-else>
        <el-menu-item :index="item.path" :key="item.path">{{ item.name }}</el-menu-item>
      </template>
    </template>
  </el-menu>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { menuItems } from '@/router'

const router = useRouter()
const route = useRoute()

const activeIndex = ref(route.path)
const handleSelect = (key: string) => {
  console.log({ route: key })
  router.push(key)
}
</script>
