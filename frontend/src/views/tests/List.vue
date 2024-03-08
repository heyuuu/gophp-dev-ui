<template>
    <el-row>
        <el-col :span="8">
            <!-- 搜索项 -->
            <el-form label-width="auto" style="padding: 10px 20px;" size="small">
                <el-form-item label="Name">
                    <el-input v-mode="searchName"></el-input>
                </el-form-item>
                <el-form-item label="Content">
                    <el-input v-mode="searchContent"></el-input>
                </el-form-item>
                <el-form-item label="Switches">
                    <el-checkbox v-model="hidePass" label="Hide Pass" />
                    <el-checkbox v-model="hideSkip" label="Hide Skip" />
                    <el-checkbox v-model="sortResult" label="Sort" />
                    <el-checkbox v-model="lockList" label="Lock" />
                </el-form-item> 
                <el-form-item>
                    <el-button type="link">重试当前测试</el-button>
                </el-form-item> 
            </el-form label="111">
            <!-- 列表 -->
            <el-table :data="showTests">
                <el-table-column label="#" type="index" width="50"/>
                <el-table-column label="name"   prop="name" />
                <el-table-column label="status" prop="status" width="100" />
                <el-table-column label="herf">
                    <el-link type="primary" href="" target="_blank">detail</el-link>&nbsp;
                    <el-link type="primary" href="" target="_blank">retry</el-link>&nbsp;
                    <el-link type="primary" href="" target="_blank">paste</el-link>
                </el-table-column>
            </el-table>
        </el-col>
        <!-- 细节展示区 -->
        <el-col :span="16">
            222
        </el-col>
    </el-row>
</template>

<script setup lang="ts">
import { ApiTestList } from "@/api/tests";
import { id } from "element-plus/es/locales.mjs";
import { ref, type Ref, watch, computed, onMounted } from "vue";
import { useRoute } from "vue-router";

// uri 参数
const route = useRoute()
const src = (route.query.src || '') as string
const path = (route.query.path || '') as string
const offset = route.query.offset ? parseInt(route.query.offset as string) : 0
const limit = route.query.limit ? parseInt(route.query.limit as string) : 0

// 筛选项及开关
const searchName = ref("")
const searchContent = ref("")
const hidePass = ref(false)
const hideSkip = ref(false)
const sortResult = ref(false)
const lockList = ref(false)

// 列表原始数据
type Status = 'WAIT' | 'RUNNING' | 'PASS' | 'SKIP' | 'FAIL' | 'NET'
type Test = {
    idx: number
    name: string
    status: Status
    content: string
    code: string
}

const tests: Ref<Test[]> = ref([])
const waitintIndexes: Ref<number[]> = ref([])

// 初始化列表
onMounted(async () => {
    const rep = await ApiTestList({
        src: src,
        path: path,
    })
    if (rep.code !== 0) {
        return
    }

    tests.value = rep.data.list.map((name, idx) => {
        return {
            idx: idx,
            name: name,
            status: 'WAIT',
            content: '',
            code: '',
        }
    })
    waitintIndexes.value = tests.value.map(test => test.idx)
})

// 执行指定case
function runTest(index: number) {
    const test = tests.value[index]
    if (!test) {
        return
    }
    if (test.status === 'RUNNING') {
        return
    }
    test.status = "RUNNING"
    test.content = ""

    // todo
}

// 遍历待处理case
const runNext = () => {
    // todo
}

// 列表筛选数据
let lockListIndexes: number[] = []
const showTests = computed(() => {
    return tests.value
})
// todo

</script>