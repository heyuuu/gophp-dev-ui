<template>
  <el-card class="show-card" body-style="padding: 0px 10px; position: relative">
    <span class="log-card-icon" @click="switchLogShowMode">+</span>
    <pre :class="'log-content-mode-' + logShowMode">{{ info }}</pre>
  </el-card>
  <div>
    <DiffEditor :original="output" :modified="expect" height="90vh" />
  </div>
</template>

<style scoped>
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
  position: absolute;
  top: 5px;
  right: 35px;

  width: 20px;
  height: 20px;
  text-align: center;
  font-size: 18px;
  line-height: 18px;
  font-weight: bolder;
  user-select: none;
}
</style>

<script setup lang="ts">
import { ref, computed } from 'vue'
import DiffEditor from '@/components/DiffEditor.vue'

// 传入属性
const props = defineProps({
  info: { type: String, default: '' },
  output: { type: String, default: '' },
  expect: { type: String, default: '' }
})

const info = computed(() => props.info.trim())
const output = computed(() => props.output)
const expect = computed(() => props.expect)

// 日志展示模式切换
const logShowMode = ref(0)
function switchLogShowMode() {
  logShowMode.value = (logShowMode.value + 1) % 3
}
</script>
