<template>
  <el-row class="main" @keyup.ctrl.enter="runCode">
    <!-- 代码栏 -->
    <el-col :span="6" class="code-pane">
      <!-- 自动更新按钮 -->
      <el-switch v-model="autoRefresh" active-text="定时自动更新" class="run-switch" />
      <el-switch v-model="openDiffMode" active-text="对比模式" class="run-switch" />

      <!-- 代码编辑器 -->
      <CodeEditor v-model="code" height="80vh" />

      <!-- 主动运行按钮 -->
      <el-button type="primary" size="default" @click="runCode">Run</el-button>

      <!-- 错误信息 -->
      <div>
        <span style="color: red">{{ error }}</span>
      </div>
    </el-col>
    <!-- 结果栏 -->
    <el-col :span="18" class="result-pane">
      <ResultPane v-if="!openDiffMode" :items="result" />
      <DiffEditor
        v-if="openDiffMode"
        :original="runResult"
        :modified="runRawResult"
        height="100vh"
      />
    </el-col>
  </el-row>
</template>

<style scoped>
.main {
  background-color: var(--color-background-soft);
}

.code-pane {
  padding: 10px 5px;
}
.result-pane {
  padding: 10px 5px;
}

.run-switch {
  --el-text-color-primary: var(--vt-c-white-soft);
  margin-bottom: 10px;
  padding-right: 20px;
}
</style>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import type { Ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import { apiRunCode } from '@/api/run'
import CodeEditor from '@/components/CodeEditor.vue'
import DiffEditor from '@/components/DiffEditor.vue'
import ResultPane from './ResultPane.vue'
import type { Item as ResultItem } from './ResultPane.vue'

// 从路由获取参数
const route = useRoute()
const code = ref((route.query.code as string) || '')
const autoRefresh = ref(!!route.query.refresh) // 自动刷新开关
const openDiffMode = ref(!!route.query.diff) // 对比模式开关

// 更新路由
const query = computed(() => {
  const qs: Record<string, string> = {}
  if (code.value !== '') {
    qs['code'] = code.value
  }
  if (autoRefresh.value) {
    qs['refresh'] = '1'
  }
  if (openDiffMode.value) {
    qs['diff'] = '1'
  }
  return qs
})

const router = useRouter()
watch(query, () => {
  router.push({ query: query.value })
})

// 执行结果
const result: Ref<ResultItem[]> = ref([])
const error = ref('')
const runResult = ref('')
const runRawResult = ref('')

function updateResult(res: ResultItem[], err: string) {
  result.value = res
  error.value = err

  runResult.value = ''
  runRawResult.value = ''
  res.forEach((item) => {
    if (item.type === 'Run') {
      runResult.value = item.content
    } else if (item.type === 'Run-Raw') {
      runRawResult.value = item.content
    }
  })

  // 触发弹窗
  if (res.length > 0 && err == '') {
    showMessage('success', '请求成功')
  } else if (err) {
    showMessage('error', err)
  }
}
function showMessage(type: 'success' | 'error', message: string) {
  ElMessage({ message, type, grouping: true, offset: 8 })
}

// 调用 Api 执行代码
let runIndex = 0
function runCode() {
  if (code.value == '') {
    updateResult([], '')
    return
  }

  runIndex++
  const currIndex = runIndex

  apiRunCode({
    code: code.value
  }).then(
    (res) => {
      if (currIndex !== runIndex) {
        return
      }
      if (res.code !== 0) {
        updateResult([], '请求失败: ' + res.error)
        return
      }

      updateResult(res.data.result, '')
    },
    () => {
      if (currIndex !== runIndex) {
        return
      }
      updateResult([], '请求失败，确认服务是否可用')
    }
  )
}
watch(code, runCode)
onMounted(runCode)

// 自动刷新
let autoRefreshId = 0
function updateAutoRefresh() {
  // 强制先清理定时任务
  if (autoRefreshId !== 0) {
    clearInterval(autoRefreshId)
    autoRefreshId = 0
  }
  // 开启定时任务
  if (autoRefresh.value) {
    autoRefreshId = setInterval(runCode, 2000)
  }
}
watch(autoRefresh, updateAutoRefresh)
onMounted(updateAutoRefresh)
</script>
