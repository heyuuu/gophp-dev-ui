<template>
  <div class="editor" ref="editorDom" :style="{ height: props.height + 'px', width: '100%' }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import type { Ref } from 'vue'
import * as monaco from 'monaco-editor'

const model = defineModel<string>()

const props = defineProps({
  value: { type: String, default: '' },
  height: { type: Number, default: 200 }
})

// 编辑器绑定 dom
const editorDom = ref(null) as Ref<unknown> as Ref<HTMLElement>
let editor: monaco.editor.IStandaloneCodeEditor | null = null
onMounted(() => {
  // console.log({ dom: editorDom.value })
  // console.log({
  //   value: props.value,
  //   height: props.height
  // })

  // 初始化编辑器
  editor = monaco.editor.create(editorDom.value, {
    value: props.value,
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
