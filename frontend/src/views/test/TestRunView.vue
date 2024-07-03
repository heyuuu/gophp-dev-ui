<template>
  <el-row @keyup.ctrl.enter="run">
    <el-col :span="8">
      <CodeEditor v-model="testContent" :height="400" />
    </el-col>
    <el-col :span="16" class="detail-card">
      <RunResultCard :info="runResult.info" :output="runResult.output" :expect="runResult.expect" />
    </el-col>
  </el-row>
</template>

<style scoped>
.main-card {
  padding: 10px 0 0 10px;
}
.detail-card {
  padding: 10px 10px 0 5px;
}
</style>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import CodeEditor from '@/components/CodeEditor.vue'
import { apiTestDetail, apiTestRunCustom } from '@/api/test'
import RunResultCard from '@/components/test/RunResultCard.vue'

// 从路由path获取参数
const props = defineProps<{ mode: string }>()
const mode = (props.mode || '') as string

// uri 参数
const route = useRoute()
const root = (route.query.root || '') as string
const path = (route.query.path || '') as string

// sections
const testContent = ref('')

// 初始化case数据
onMounted(async () => {
  const data = await apiTestDetail({
    mode: mode,
    root: root,
    path: path
  })

  testContent.value = data.content
  run()
})

// 执行
function run() {
  updateResult('执行中...', '', '')
  apiTestRunCustom({
    mode: mode,
    root: root,
    path: path,
    content: testContent.value
  }).then(
    (data) => {
      const result = data.result
      updateResult(result.info, result.output, result.expected)
    },
    () => {
      updateResult('调用 url 失败', '', '')
    }
  )
}

// 执行结果
const runResult = ref({
  info: '',
  output: '',
  expect: ''
})
function updateResult(info: string, output: string, expect: string) {
  runResult.value = { info, output, expect }
}
</script>
