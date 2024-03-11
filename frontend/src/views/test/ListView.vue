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
            <el-button type="">重试当前测试</el-button>
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
            <el-link type="primary" target="_blank" :href="pageTestRun(src, scope.row.name)"
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
              {{ currTest.content }}
          </pre>
        </el-card>
        <div>
          <DiffEditor :original="currTest.output" :modified="currTest.expect" height="90vh" />
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
import { ref, type Ref, watch, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { pageTestRun } from '@/router/routes'
import { apiTestList, apiTestRun } from '@/api/test'
import type { RunStatus } from '@/api/test'
import { clipboardWriteText } from '@/utils/utils'
import { ElMessage } from 'element-plus'
import DiffEditor from '@/components/DiffEditor.vue'
import { off } from 'process'

// uri 参数
const route = useRoute()
const src = (route.query.src || '') as string
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
type Status = 'WAITING' | 'RUNNING' | 'NET_ERR' | RunStatus
const statusLevel: { [_: string]: number } = {
  WAITING: 1,
  RUNNING: 2,
  PASS: 3,
  SKIP: 4,
  NET_ERR: 5,
  FAIL: 6
}
function getStatusLevel(status: Status): number {
  return statusLevel[status] || 0
}

type Test = {
  idx: number
  name: string
  status: Status
  content: string
  code: string
  output: string
  expect: string
}

const tests: Ref<Test[]> = ref([])
const waitintIndexes: Ref<number[]> = ref([])
const statusCount = ref({
  WAITING: 0,
  RUNNING: 0,
  NET_ERR: 0,
  // RunStatus
  PASS: 0,
  BORK: 0,
  FAIL: 0,
  WARN: 0,
  LEAK: 0,
  XFAIL: 0,
  XLEAK: 0,
  SKIP: 0,
  SLOW: 0
})

// 初始化列表
onMounted(async () => {
  const rep = await apiTestList({
    src: src,
    path: path,
    offset: offset,
    limit: limit
  })
  if (rep.code !== 0) {
    return
  }

  tests.value = rep.data.list.map((name, idx) => {
    return {
      idx: idx,
      name: name,
      status: 'WAITING',
      content: '',
      code: '',
      output: '',
      expect: ''
    }
  })
  waitintIndexes.value = tests.value.map((test) => test.idx)
  statusCount.value.WAITING = waitintIndexes.value.length

  runNext() // 触发执行 loop
})

// 更新状态
function updateTest(test: Test, status: Status, content: string) {
  statusCount.value[test.status]--
  statusCount.value[status]++
  test.status = status
  test.content = content
  test.code = ''
}

// 执行指定case
function runTest(index: number) {
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
    src: src,
    path: test.name
  }).then(
    (res) => {
      if (res.code !== 0) {
        updateTest(test, 'FAIL', '执行失败: error=' + res.error)
        return
      }

      updateTest(test, res.data.status, res.data.info)
      test.code = res.data.code
      test.output = res.data.output
      test.expect = res.data.expect
      runNext()
    },
    () => {
      updateTest(test, 'NET_ERR', '调用 url 失败')
      runNext()
    }
  )
}

// 遍历待处理case
const runNext = () => {
  while (statusCount.value.RUNNING < 10 && waitintIndexes.value.length > 0) {
    const index = waitintIndexes.value.shift()
    console.log({ index })
    if (index !== undefined) {
      runTest(index)
    }
  }
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
const showTestClassName: (_: { row: Test }) => string = ({ row: test }) => {
  let classes = ''
  if (test.idx === activeIndex.value) {
    classes = 'active-row '
  }

  switch (test.status) {
    case 'PASS':
      return classes + 'success-row'
    case 'SKIP':
      return classes + 'warning-row'
    case 'FAIL':
      return classes + 'danger-row'
    default:
      return classes + 'info-row'
  }
}

// 当前选中案例index
const activeIndex = ref(-1)
function activeTest(test: Test) {
  console.log({ old: activeIndex.value, new: test.idx })
  activeIndex.value = test.idx
}
const currTest = computed(() => {
  const ct = tests.value[activeIndex.value]
  if (ct) {
    console.log({
      content: ct.content,
      output: ct.output,
      expect: ct.expect
    })
  }
  return tests.value[activeIndex.value]
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

  let offset = -1
  for (let i = 0; i < length; i++) {
    if (targets[i].idx === currIdx) {
      offset = i
      break
    }
  }

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

function filterTests() {
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
    if (
      searchContent.value &&
      !test.content.includes(searchContent.value) &&
      !test.code.includes(searchContent.value)
    ) {
      return false
    }
    return true
  })
  // 排序案例
  if (sortResult.value) {
    result.sort((t1, t2): number => {
      const statusDiff = getStatusLevel(t1.status) - getStatusLevel(t2.status)
      if (statusDiff !== 0) {
        return statusDiff
      }
      if (t1.content < t2.content) {
        return -1
      } else if (t1.content > t2.content) {
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
