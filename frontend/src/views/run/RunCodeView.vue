<template>
  <el-row>
    <!-- 代码栏 -->
    <el-col :span="6" class="code-pane">
      <!-- 开关/操作 -->
      <div class="code-options">
        <el-switch v-model="autoRefresh" active-text="定时刷新" />
        <el-button type="primary" @click="runCode">Run</el-button>
      </div>

      <!-- 代码编辑器 -->
      <CodeEditor
        v-model="code"
        height="80vh"
        class="code-editor"
        :class="error ? 'code-editor-error' : ''"
      />

      <!-- 错误信息 -->
      <div>
        <span style="color: red">{{ error }}</span>
      </div>
    </el-col>

    <!-- 结果栏 -->
    <el-col :span="18" class="result-pane">
      <!-- 结果展示相关选项 -->
      <div class="result-options" v-if="allTypes">
        <!-- 对比模式开关 -->
        <el-switch v-model="openDiffMode" active-text="对比模式" />
        <!-- 执行结果开关 -->
        <template v-if="!openDiffMode">
          <el-button-group>
            <el-button
              v-for="typ in allTypes"
              :key="typ"
              :type="isShowType(typ) ? 'primary' : 'info'"
              @click="switchShowType(typ)"
              >{{ typ }}</el-button
            >
          </el-button-group>
        </template>
        <!-- 对比模式选项 -->
        <template v-if="openDiffMode">
          <p>左侧</p>
          <el-select v-model="diffTypeLeft" style="width: 120px">
            <el-option key="@" label="@" value="@" />
            <el-option v-for="typ in allTypes" :key="typ" :label="typ" :value="typ" />
          </el-select>
          <p>右侧</p>
          <el-select v-model="diffTypeRight" style="width: 120px">
            <el-option key="@" label="@" value="@" />
            <el-option v-for="typ in allTypes" :key="typ" :label="typ" :value="typ" />
          </el-select>
        </template>
      </div>

      <!-- 分栏结果栏 -->
      <div class="result-container" v-if="!openDiffMode">
        <template v-for="typ in allTypes" :key="typ">
          <el-card v-if="isShowType(typ)" class="result-card" :header="typ">
            <pre>{{ getTypeContent(typ) }}</pre>
          </el-card>
        </template>
      </div>

      <!-- 对比结果栏 -->
      <DiffEditor
        v-if="openDiffMode"
        :original="getTypeContent(diffTypeLeft)"
        :modified="getTypeContent(diffTypeRight)"
        height="100vh"
      />
    </el-col>
  </el-row>
</template>

<style scoped>
.code-options {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0 10px 10px 10px;
  color: var(--vt-c-white-soft);
  --el-text-color-primary: var(--vt-c-white-soft);
}
.code-pane {
  padding: 10px 5px;
}
.result-options {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0 10px 10px 10px;
  color: var(--vt-c-white-soft);
  --el-text-color-primary: var(--vt-c-white-soft);
}
.result-pane {
  padding: 10px 5px;
}

.code-editor {
  border: 1px solid white;
}
.code-editor-error {
  border: 1px solid red;
}

.result-container {
  display: flex;
}
.result-card {
  flex: 1;
  min-width: 0;
  margin: 0 10px;

  padding-top: 0;
  --el-card-padding: 10px;
  height: 100vh;
  background-color: #1e1e1e;
  color: #d4d4d4;
  font-size: 14px;
  overflow: auto;
}
</style>

<script setup lang="ts">
import { onMounted, ref, watch, type Ref } from 'vue'
import { ElMessage } from 'element-plus'

import { apiRunConfig, apiRunCode } from '@/api/run'
import CodeEditor from '@/components/CodeEditor.vue'
import DiffEditor from '@/components/DiffEditor.vue'
import { useRouteHash } from '@/utils/hash'

// 常量
const TypeSrc = '@'

// 从路由获取参数
const code = ref('<?php\n')
const autoRefresh = ref(false) // 自动刷新开关
const openDiffMode = ref(false) // 对比模式开关
const showTypes: Ref<string[]> = ref([]) // 展示类型列表
const diffTypeLeft = ref(TypeSrc) // 对比类型选项-左侧
const diffTypeRight = ref(TypeSrc) // 对比类型选项-右侧

// 参数同步到路由上
useRouteHash({ code, autoRefresh, openDiffMode, showTypes, diffTypeLeft, diffTypeRight })

// 初始化配置
const allTypes: Ref<string[]> = ref([]) // 全量的类型列表
onMounted(async () => {
  apiRunConfig().then((res) => {
    allTypes.value = res.types
    if (showTypes.value.length === 0) {
      showTypes.value = allTypes.value
    }
  })
})

// 执行结果开关
function isShowType(typ: string): boolean {
  return showTypes.value.includes(typ)
}
function switchShowType(typ: string) {
  if (isShowType(typ)) {
    showTypes.value = showTypes.value.filter((v) => v !== typ)
  } else {
    showTypes.value = [...showTypes.value, typ]
  }
}

// 执行结果相关
type ResultItem = {
  type: string
  content: string
}

const result: Ref<ResultItem[]> = ref([])
const error = ref('')

function updateResult(res: ResultItem[], err: string) {
  if (err) {
    error.value = err
    return
  }

  result.value = res
  error.value = err
}
function showMessage(type: 'success' | 'error', message: string) {
  ElMessage({ message, type, grouping: true, offset: 8 })
}
function getTypeContent(typ: string): string {
  if (typ === TypeSrc) {
    return code.value
  }

  const item = result.value.find((item) => item.type === typ)
  return item ? item.content : ''
}

// 调用 Api 执行代码
let runIndex = 0
function runCode() {
  if (code.value === '') {
    updateResult([], '')
    return
  }

  runIndex++
  const currIndex = runIndex

  apiRunCode({ code: code.value }).then(
    (res) => {
      if (currIndex !== runIndex) {
        return
      }
      showMessage('success', '请求成功')
      updateResult(res.result, '')
    },
    (err) => {
      if (currIndex !== runIndex) {
        return
      }
      updateResult([], `请求失败: ${err}`)
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
