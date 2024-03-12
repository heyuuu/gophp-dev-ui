<template>
  <el-row>
    <el-col :span="24">
      <el-button-group style="margin-bottom: 10px">
        <el-button
          v-for="typ in types"
          :key="typ"
          :type="isShowType(typ) ? 'primary' : 'info'"
          @click="swithShowType(typ)"
          >{{ typ }}</el-button
        >
      </el-button-group>
    </el-col>
  </el-row>
  <el-row>
    <el-col
      v-for="(item, index) in showItems"
      :key="item.type"
      :span="spanWidth(index)"
      class="show-col"
    >
      <el-card class="show-card" :header="item.type">
        <pre>{{ item.content }}</pre>
      </el-card>
    </el-col>
  </el-row>
</template>

<style scoped>
.show-col {
  padding: 0px 5px;
}

.show-card {
  --el-card-padding: 10px;
  height: 100vh;
  background-color: #1e1e1e;
  color: #d4d4d4;
  font-size: 14px;
  overflow: auto;
}
</style>

<script setup lang="ts">
import { ref, computed } from 'vue'

export type Item = {
  type: string
  content: string
}
const props = defineProps<{
  items: Item[]
}>()

// 展示类型
const __all__ = '__all__'
const types = computed(() => props.items.map((item) => item.type))
const isShow = ref({ __all__: true } as { [K: string]: boolean })
function isShowType(typ: string): boolean {
  return isShow.value[__all__] || isShow.value[typ]
}
function swithShowType(typ: string) {
  // 若未初始化，则初始化 isShow
  if (isShow.value[__all__]) {
    isShow.value = {}
    props.items.forEach((item) => (isShow.value[item.type] = true))
  }
  // 切换 type 是否展示
  isShow.value[typ] = !isShow.value[typ]
}

// 展示内容
const showItems = computed(() => props.items.filter((item) => isShowType(item.type)))

// 分配宽度
function spanWidth(index: number): number {
  const total = 24
  const count = showItems.value.length
  // 平均宽度
  const width = Math.floor(total / count)
  // 不能平均时优先分配到左侧
  if (total - width * count > index) {
    return width + 1
  } else {
    return width
  }
}
</script>
