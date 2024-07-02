<template>
  <div class="editor" ref="editorDom" :style="{ height: height, width: '100%' }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import type { Ref } from 'vue'
import * as monaco from 'monaco-editor'
// import * as monaco from 'monaco-editor/esm/vs/editor/editor.api'
// import 'monaco-editor/esm/vs/basic-languages/php/php.contribution'

type Editor = monaco.editor.IStandaloneDiffEditor
type EditorOptions = monaco.editor.IStandaloneEditorConstructionOptions

const props = defineProps<{
  original: string
  modified: string
  originalLanguage?: EditorOptions['language']
  modifiedLanguage?: EditorOptions['language']

  height?: number | string
  readonly?: EditorOptions['readOnly']
  lineNumbers?: EditorOptions['lineNumbers']
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
let editor: Editor | null = null
onMounted(() => {
  // 初始化编辑器
  editor = monaco.editor.createDiffEditor(editorDom.value)

  // 初始化
  updateOptions()
  updateModel()
})

// 更新配置
const options = computed(() => ({
  theme: 'vs-dark',
  automaticLayout: true,
  readOnly: props.readonly,
  lineNumbers: props.lineNumbers
}))
function updateOptions() {
  editor?.updateOptions(options.value)
}
watch(options, updateOptions)

// 更新文本及语言
const model = computed(() => ({
  original: props.original,
  modified: props.modified,
  originalLanguage: props.originalLanguage,
  modifiedLanguage: props.modifiedLanguage
}))
function updateModel() {
  if (!editor) {
    return
  }
  editor.setModel({
    original: monaco.editor.createModel(model.value.original, model.value.originalLanguage),
    modified: monaco.editor.createModel(model.value.modified, model.value.modifiedLanguage)
  })
}
watch(model, updateModel)
</script>
