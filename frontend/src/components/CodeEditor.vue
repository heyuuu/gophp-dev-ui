<template>
  <div class="editor" ref="editorDom" :style="{ height: height, width: '100%' }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import type { Ref } from 'vue'
import * as monaco from 'monaco-editor'
// import * as monaco from 'monaco-editor/esm/vs/editor/editor.api'
// import 'monaco-editor/esm/vs/basic-languages/php/php.contribution'

type Editor = monaco.editor.IStandaloneCodeEditor
type EditorOptions = monaco.editor.IStandaloneEditorConstructionOptions

const text = defineModel<string>()
const props = defineProps<{
  language?: EditorOptions['language']

  height?: number | string
  readonly?: EditorOptions['readOnly']
  lineNumbers?: EditorOptions['lineNumbers']
}>()

// 内容对象及更新

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
  editor = monaco.editor.create(editorDom.value)

  // 监听修改文本事件
  editor.onDidChangeModelContent(() => {
    text.value = editor?.getValue()
  })

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
  text: text.value || '',
  language: props.language
}))
function updateModel() {
  if (!editor) {
    return
  }
  if (model.value.text === editor.getValue()) {
    return
  }
  editor.setModel(monaco.editor.createModel(model.value.text, model.value.language))
}
watch(model, updateModel)
</script>
