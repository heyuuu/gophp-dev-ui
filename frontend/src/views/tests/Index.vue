<template>
  <!-- 搜索框 -->
  <el-row style="margin-top: 20px;">
    <el-col :span="12" :offset="6">
      <el-input v-model="src" class="input">
        <template #prepend>src</template>
        <template #append>
          <el-button :icon="Search"></el-button>
        </template>
      </el-input>
    </el-col>
  </el-row>
  <!-- 结果列表 -->
  <el-row style="margin-top: 20px;">
    <el-col :span="18" :offset="3">
      <el-table :data="tableData" height="800">
        <el-table-column type="index"></el-table-column>
        <el-table-column prop="name" label="Name"></el-table-column>
        <el-table-column prop="name" label="Link">
          <template #default="scope">
            <el-link type="primary" target="_blank" :href='"/tests/list?src=" + src + "&path=" + scope.row.name'>列表页</el-link>
          </template>
        </el-table-column>
      </el-table>†
    </el-col>
  </el-row>
</template>

<script setup lang="ts">
import { ref, watch, computed, type Ref } from "vue";
import { Search } from "@element-plus/icons-vue";
import { ApiTestsPathList } from "../../api/tests"

const src = ref("/Users/heyu/Code/src/php-7.4.33")
const dirList: Ref<String[]> = ref([])

const tableData = computed(() => {
  return dirList.value.map((dir) => {
    return {
      name: dir,
    }
  })
})

// 更新列表
async function update() {
  if (src.value == "") {
    dirList.value = []
    return
  }

  const rep = await ApiTestsPathList({
    src: src.value,
  })
  dirList.value = rep.data.data.list
  console.log({rep})
}

watch(src, update)
update()

</script>
