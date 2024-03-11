<template>
  <div class="editor" ref="editorDom" :style="{ height: props.height, width: '100%' }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import type { Ref } from 'vue'
import * as monaco from 'monaco-editor'

const props = defineProps({
  original: { type: String, required: true },
  modified: { type: String, required: true },
  height: { type: String, default: '200px' }
})

const model = computed(() => ({
  original: props.original,
  modified: props.modified
}))

// 编辑器绑定 dom
const editorDom = ref(null) as Ref<unknown> as Ref<HTMLElement>
let editor: monaco.editor.IStandaloneDiffEditor | null = null
onMounted(() => {
  // console.log({ dom: editorDom.value })
  // console.log({
  //   original: props.original,
  //   modified: props.modified,
  //   height: props.height
  // })
  editor = monaco.editor.createDiffEditor(editorDom.value, {
    theme: 'vs-dark',
    automaticLayout: true
  })
  updateModel()
})

// 更新文本
function updateModel() {
  // console.log('updateModel')
  if (editor) {
    editor.setModel({
      original: monaco.editor.createModel(model.value.original, 'text'),
      modified: monaco.editor.createModel(model.value.modified, 'text')
    })
  }
}
watch(model, updateModel)
</script>
