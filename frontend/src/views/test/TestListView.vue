<template>
  <el-row>
    <el-col :span="8" class="main-card">
      <!-- 搜索项 -->
      <el-form label-width="auto" size="small">
        <el-form-item label="Name">
          <el-input v-model="searchName"></el-input>
        </el-form-item>
        <el-form-item label="Content">
          <el-input v-model="searchContent"></el-input>
        </el-form-item>
        <el-form-item label="Switches">
          <el-checkbox v-model="hidePass" label="Hide Pass" />
          <el-checkbox v-model="hideSkip" label="Hide Skip" />
          <el-checkbox v-model="sortResult" label="Sort" />
          <el-checkbox v-model="lockList" label="Lock" />
        </el-form-item>
        <el-form-item>
          <div class="flex-right">
            <div class="stat-line">{{ statLine }}</div>
            <el-button @click="retryShowTests">重试当前测试</el-button>
          </div>
        </el-form-item>
      </el-form>
      <!-- 列表 -->
      <el-table
        :data="showTests"
        :row-class-name="showTestClassName"
        @row-click="activeTest"
        max-height="80vh"
        @keyup.prevent.up="moveActive(false)"
        @keyup.prevent.down="moveActive(true)"
      >
        <el-table-column label="#" prop="idx" width="50" />
        <el-table-column label="name" prop="name" />
        <el-table-column label="status" prop="status" width="100" />
        <el-table-column label="herf">
          <template #default="scope">
            <el-link type="primary" target="_blank" :href="pageTestRun(mode, root, scope.row.name)"
              >detail</el-link
            >&nbsp; <el-link type="primary" @click="runTest(scope.row.idx)">retry</el-link>&nbsp;
            <el-link type="primary" @click="copyTestName(scope.row.idx)">paste</el-link>&nbsp;
          </template>
        </el-table-column>
      </el-table>
    </el-col>
    <!-- 细节展示区 -->
    <el-col :span="16" class="detail-card">
      <template v-if="currTest">
        <el-card class="show-card" body-style="padding: 0px 10px">
          <span class="log-card-icon" @click="changeLogShowMode">+</span>
          <pre :class="'log-content-mode-' + logShowMode">
              {{ currTest.info }}
          </pre>
        </el-card>
        <div>
          <DiffEditor
            :original="currTest?.result?.output || ''"
            :modified="currTest?.result?.expected || ''"
            height="90vh"
          />
        </div>
      </template>
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
.stat-line {
  margin-right: 20px;
  font-size: 12px;
}
.show-card {
  color: var(--el-text-color-primary);
  background-color: white;
}
.log-content-mode-0 {
  overflow-y: auto;
  max-height: 10vh;
}
.log-content-mode-1 {
  overflow-y: auto;
  max-height: 30vh;
}
.log-content-mode-2 {
  overflow-y: auto;
  max-height: 100vh;
}
.log-card-icon {
  float: right;
  width: 20px;
  height: 20px;
  font-weight: bolder;
  font-size: 18px;
  user-select: none;
}
</style>

<script setup lang="ts">
import { ref, computed, onMounted, watch, type Ref } from 'vue'
import { useRoute } from 'vue-router'
import { pageTestRun } from '@/router/routes'
import { apiTestCaseList, apiTestRun, type TestResultStatus, type TestResult } from '@/api/test'
import { clipboardWriteText } from '@/utils/utils'
import { ElMessage } from 'element-plus'
import DiffEditor from '@/components/DiffEditor.vue'

// 常量
const BATCH_SIZE = 10 // 最大并发执行数

// 从路由path获取参数
const props = defineProps<{ mode: string }>()
const mode = computed(() => props.mode || '')

// 从路由query获取参数
const route = useRoute()
const root = (route.query.root || '') as string
const path = (route.query.path || '') as string
const offset = route.query.offset ? parseInt(route.query.offset as string) : 0
const limit = route.query.limit ? parseInt(route.query.limit as string) : 0

// 筛选项及开关
const searchName = ref('')
const searchContent = ref('')
const hidePass = ref(false)
const hideSkip = ref(false)
const sortResult = ref(false)
const lockList = ref(false)

// 列表原始数据
type Status = 'WAITING' | 'RUNNING' | 'NET_ERR' | TestResultStatus
const statusLevel: Partial<Record<Status, number>> = {
  WAITING: 1,
  RUNNING: 2,
  PASS: 3,
  SKIP: 4,
  NET_ERR: 5,
  FAIL: 6
}
function compareStatus(s1: Status, s2: Status): number {
  const l1 = statusLevel[s1] || 0
  const l2 = statusLevel[s2] || 0
  return l1 - l2
}

type Test = {
  idx: number
  name: string
  status: Status
  info: string
  result: TestResult | null
}

const tests: Ref<Test[]> = ref([])
const waitingIndexes: Ref<number[]> = ref([]) // 待执行的测试case索引列表
const statusCount: Ref<Record<Status, number>> = ref({
  WAITING: 0,
  RUNNING: 0,
  NET_ERR: 0,
  // RunStatus
  PASS: 0,
  FAIL: 0,
  BORK: 0,
  SKIP: 0
})

// 初始化列表
onMounted(async () => {
  const data = await apiTestCaseList({
    mode: mode.value,
    root: root,
    path: path,
    offset: offset,
    limit: limit
  })
  tests.value = data.list.map((name, idx) => {
    return {
      idx: idx,
      name: name,
      status: 'WAITING',
      info: '',
      result: null
    }
  })
  waitingIndexes.value = tests.value.map((test) => test.idx)
  statusCount.value.WAITING = waitingIndexes.value.length
})

// 更新状态
function updateTest(test: Test, status: Status, info: string) {
  statusCount.value[test.status]--
  statusCount.value[status]++
  test.status = status
  test.info = info
}

// 执行指定case
function runTest(index: number, onFinish: () => void = () => {}) {
  const test = tests.value[index]
  if (!test) {
    return
  }
  if (test.status === 'RUNNING') {
    return
  }
  updateTest(test, 'RUNNING', '')

  // api
  apiTestRun({
    mode: mode.value,
    root: root,
    path: test.name
  }).then(
    (data) => {
      const result = data.result
      updateTest(test, result.status, result.info)
      test.result = result

      onFinish()
    },
    (err) => {
      updateTest(test, 'NET_ERR', '调用 url 失败:' + err)
      test.result = null

      onFinish()
    }
  )
}

// 遍历待处理case
const runNext = () => {
  while (statusCount.value.RUNNING < BATCH_SIZE && waitingIndexes.value.length > 0) {
    const index = waitingIndexes.value.shift()
    console.log({ index })
    if (index !== undefined) {
      runTest(index, runNext)
    }
  }
}
watch(waitingIndexes, runNext)

// 重试当前展示的case列表
function retryShowTests() {
  const showTestIndexes = showTests.value.map((test) => test.idx)
  waitingIndexes.value.push(...showTestIndexes)
}

// 状态栏
const statLine = computed(() => {
  const total = tests.value.length
  const count = statusCount.value
  return `total: ${total}, running: ${count.RUNNING}, waiting: ${count.WAITING}, pass: ${count.PASS}, skip: ${count.SKIP}, fail: ${count.FAIL}, net_err: ${count.NET_ERR}`
})

// 复制case名
function copyTestName(index: number) {
  const test = tests.value[index]
  if (!test) {
    return
  }
  clipboardWriteText(test.name).then(
    (msg) => {
      ElMessage({
        message: msg,
        type: 'success'
      })
    },
    (msg) => {
      ElMessage({
        message: msg,
        type: 'warning'
      })
    }
  )
}

// 展示行的状态
function showTestClassName(data: { row: Test; rowIndex: number }): string {
  const test = data.row

  let classes = ''
  if (test.idx === activeIndex.value) {
    classes = 'active-row '
  }

  switch (test.status) {
    case 'PASS':
      return classes + 'success-row'
    case 'SKIP':
      return classes + 'warning-row'
    case 'WAITING':
    case 'RUNNING':
      return classes + 'info-row'
    case 'FAIL':
    default:
      return classes + 'danger-row'
  }
}

// 当前选中案例index
const activeIndex = ref(-1)
function activeTest(test: Test) {
  console.log({ old: activeIndex.value, new: test.idx })
  activeIndex.value = test.idx
}
const currTest = computed(() => {
  const test = tests.value[activeIndex.value] || null
  return test
})

// 向前/后切换选中案例
function moveActive(next: boolean) {
  console.log(next ? 'move next' : 'move prev')

  const currIdx = activeIndex.value
  const targets = showTests.value
  const length = targets.length
  if (length === 0) {
    return
  }

  // 找到原选中位置，未找到时为 -1
  let offset = targets.findIndex((test) => test.idx == currIdx)

  // 移动到新位置
  let newOffset = -1
  if (offset < 0) {
    // not found
    newOffset = next ? 0 : length - 1
  } else {
    newOffset = (next ? offset + 1 : offset + length - 1) % length
  }

  activeIndex.value = targets[newOffset].idx
}

// 列表筛选数据
let lockListIndexes: number[] = []
const showTests = computed(() => {
  console.log('showTests')
  // 判断是否锁列表
  if (lockList.value) {
    return lockListIndexes.map((idx) => tests.value[idx])
  } else {
    const resultTests = filterTests()
    lockListIndexes = resultTests.map((t) => t.idx)
    return resultTests
  }
})

function isMatchContent(test: Test, searchContent: string): boolean {
  if (test.info.includes(searchContent)) {
    return true
  }
  if (test.result?.code && test.result.code.includes(searchContent)) {
    return true
  }
  return false
}

function filterTests(): Test[] {
  // 过滤案例
  let result = tests.value.filter((test) => {
    if (test.idx === activeIndex.value) {
      return true
    }
    if (hidePass.value && test.status === 'PASS') {
      return false
    }
    if (hideSkip.value && test.status === 'SKIP') {
      return false
    }
    if (searchName.value && !test.name.includes(searchName.value)) {
      return false
    }
    if (searchContent.value && !isMatchContent(test, searchContent.value)) {
      return false
    }
    return true
  })
  // 排序案例
  if (sortResult.value) {
    result.sort((t1, t2): number => {
      const result = compareStatus(t1.status, t2.status)
      if (result) {
        return result
      }

      if (t1.info < t2.info) {
        return -1
      } else if (t1.info > t2.info) {
        return 1
      } else {
        return t1.idx - t2.idx
      }
    })
  }

  return result
}

// 日志展示
const logShowMode = ref(0)
function changeLogShowMode() {
  logShowMode.value = (logShowMode.value + 1) % 3
}
</script>
