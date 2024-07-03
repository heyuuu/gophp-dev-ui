<template>
  <!-- 搜索框 -->
  <el-row style="margin-top: 20px">
    <el-col :span="12" :offset="6">
      <el-input v-model="root" class="input">
        <template #prepend>root</template>
        <template #append>
          <el-button :icon="Search" @click="update"></el-button>
        </template>
      </el-input>
    </el-col>
  </el-row>
  <!-- 结果列表 -->
  <el-row style="margin-top: 20px">
    <el-col :span="18" :offset="3">
      <el-table :data="tableData" height="800">
        <el-table-column type="index" />
        <el-table-column prop="name" label="Name" />
        <el-table-column prop="name" label="Link">
          <template #default="scope">
            <el-link type="primary" target="_blank" :href="pageTestList(mode, root, scope.row.name)"
              >列表页</el-link
            >
          </template>
        </el-table-column>
      </el-table>
    </el-col>
  </el-row>
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted, type Ref } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { apiTestConfig, apiTestPathList } from '@/api/test'
import { pageTestList } from '@/router/routes'

// 从路由path获取参数
const props = defineProps<{ mode: string }>()
const mode = props.mode || ''

//
const root = ref('')
const dirList: Ref<string[]> = ref([])
const tableData = computed(() => dirList.value.map((name) => ({ name })))

// 初始化配置
onMounted(async () => {
  const config = await apiTestConfig({ mode: mode })
  root.value = config.defaultTestRoot
})

// 更新列表
async function update() {
  if (root.value === '') {
    dirList.value = []
    return
  }

  const data = await apiTestPathList({
    mode: mode,
    root: root.value
  })
  dirList.value = data.list
  console.log({ data })
}
watch(root, update)
</script>
