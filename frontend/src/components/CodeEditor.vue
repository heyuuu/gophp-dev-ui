<template>
  <div class="editor" ref="editorDom" :style="{ height: height, width: '100%' }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import type { Ref } from 'vue'
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api'

const model = defineModel<string>()
const props = defineProps<{
  height: number | string
}>()

const height = computed(() => {
  const h = props.height
  if (typeof h == 'number') {
    return h + 'px'
  } else if (typeof h == 'string' && h !== '') {
    return h
  } else {
    return '200px'
  }
})

// 编辑器绑定 dom
const editorDom = ref(null) as Ref<unknown> as Ref<HTMLElement>
let editor: monaco.editor.IStandaloneCodeEditor | null = null
onMounted(() => {
  // 初始化编辑器
  editor = monaco.editor.create(editorDom.value, {
    language: 'php',
    theme: 'vs-dark',
    automaticLayout: true
  })

  // 监听修改文本事件
  editor.onDidChangeModelContent((e) => {
    model.value = editor?.getValue()
  })

  updateModel()
})

// 更新文本
function updateModel() {
  if (!editor) {
    return
  }

  const newValue = model.value || ''
  if (editor.getValue() !== newValue) {
    editor.setValue(newValue)
  }
}
watch(model, updateModel)
</script>
