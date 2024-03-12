<template>
  <el-row class="main">
    <!-- 代码栏 -->
    <el-col :span="6" class="code-pane">
      <!-- 自动更新按钮 -->
      <el-switch v-model="autoRefresh" active-text="定时自动更新" class="run-switch" />
      <el-switch v-model="openDiffMode" active-text="对比模式" class="run-switch" />

      <!-- 代码编辑器 -->
      <CodeEditor v-model="code" height="80vh" />

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
import { onMounted, ref, watch } from 'vue'
import type { Ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { apiRunCode } from '@/api/run'
import CodeEditor from '@/components/CodeEditor.vue'
import DiffEditor from '@/components/DiffEditor.vue'
import ResultPane from './ResultPane.vue'
import type { Item as ResultItem } from './ResultPane.vue'

const route = useRoute()
const router = useRouter()
const code = ref((route.query.code as string) || '')
watch(code, () => {
  router.push({
    query: { code: code.value }
  })
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

// 对比模式
const openDiffMode = ref(false)

// 自动刷新
const autoRefresh = ref(false)
</script>
