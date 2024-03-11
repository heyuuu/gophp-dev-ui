<template>
  <div class="editor" ref="editorDom" :style="{ height: props.height + 'px', width: '100%' }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Ref } from 'vue'
import * as monaco from 'monaco-editor'

const props = defineProps({
  value: { type: String, default: '' },
  height: { type: Number, default: 200 }
})

// 编辑器绑定 dom
const editorDom = ref(null) as Ref<unknown> as Ref<HTMLElement>
onMounted(() => {
  console.log({ dom: editorDom.value })
  console.log({
    value: props.value,
    height: props.height
  })
  const editor = monaco.editor.create(editorDom.value, {
    value: props.value,
    language: 'php',
    theme: 'vs-dark',
    automaticLayout: true
  })
})
</script>
