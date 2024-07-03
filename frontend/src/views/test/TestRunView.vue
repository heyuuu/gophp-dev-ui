<template>
  <el-row @keyup.ctrl.enter="run">
    <el-col :span="8">
      <el-form label-width="auto" class="main-card">
        <!-- sections -->
        <el-form-item v-for="(sec, index) in sections" :key="index" :label="sec.type">
          <template v-if="sectionShowType(sec.type) === 'input'">
            <el-input v-model="sec.text"></el-input>
          </template>
          <template v-else-if="sectionShowType(sec.type) === 'text'">
            <el-input v-model="sec.text" type="textarea"></el-input>
          </template>
          <template v-else>
            <CodeEditor v-model="sec.text" :height="400" />
          </template>
        </el-form-item>
        <!-- buttons -->
        <el-form-item>
          <div class="flex-right">
            <el-button type="primary" @click="run">Run</el-button>
          </div>
        </el-form-item>
      </el-form>
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
import { onMounted, ref, type Ref } from 'vue'
import { useRoute } from 'vue-router'
import CodeEditor from '@/components/CodeEditor.vue'
import { apiTestDetail, apiTestRunCustom } from '@/api/test'
import type { SectionType, Section } from '@/models/test'
import { sectionMapToList, sectionListToMap } from '@/models/test'
import RunResultCard from '@/components/test/RunResultCard.vue'

// 从路由path获取参数
const props = defineProps<{ mode: string }>()
const mode = (props.mode || '') as string

// uri 参数
const route = useRoute()
const root = (route.query.root || '') as string
const path = (route.query.path || '') as string

// sections
type SectionShowType = 'input' | 'text' | 'code'
const sections: Ref<Section[]> = ref([])
function sectionShowType(type: SectionType): SectionShowType {
  switch (type) {
    case 'FILE':
    case 'SKIPIF':
    case 'CLEAN':
      return 'code'
    case 'EXPECT':
    case 'EXPECTF':
    case 'EXPECTREGEX':
      return 'text'
    default:
      return 'input'
  }
}

// 初始化case数据
onMounted(async () => {
  const data = await apiTestDetail({
    mode: mode,
    root: root,
    path: path
  })

  sections.value = sectionMapToList(data.sections)
  console.log(sections.value)
  run()
})

// 执行
function run() {
  updateResult('执行中...', '', '')
  apiTestRunCustom({
    mode: mode,
    root: root,
    path: path,
    sections: sectionListToMap(sections.value)
  }).then(
    (data) => {
      updateResult(data.info, data.output, data.expect)
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
