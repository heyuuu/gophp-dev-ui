<template>
    <el-row>
        <el-col :span="8" class="main-card">
            <!-- 搜索项 -->
            <el-form label-width="auto" size="small">
                <el-form-item label="Name">
                    <el-input v-model="searchName"></el-input>
                </el-form-item>
                <el-form-item label="Content">
                    <el-input v-model="searchContent"></el-input>
                </el-form-item>
                <el-form-item label="Switches">
                    <el-checkbox v-model="hidePass" label="Hide Pass" />
                    <el-checkbox v-model="hideSkip" label="Hide Skip" />
                    <el-checkbox v-model="sortResult" label="Sort" />
                    <el-checkbox v-model="lockList" label="Lock" />
                </el-form-item> 
                <el-form-item>
                    <div class="flex-right">
                        <div class="stat-line">{{ statLine }}</div>
                        <el-button type="">重试当前测试</el-button>
                    </div>
                </el-form-item> 
            </el-form label="111">
            <!-- 列表 -->
            <el-table :data="showTests">
                <el-table-column label="#" type="index" width="50"/>
                <el-table-column label="name"   prop="name" />
                <el-table-column label="status" prop="status" width="100" />
                <el-table-column label="herf">
                    <template #default="scope">
                        <el-link type="primary" target="_blank" :href='pageTestRun(src, scope.row.name)'>detail</el-link>&nbsp;
                        <el-link type="primary" @click="runTest(scope.row.idx)">retry</el-link>&nbsp;
                        <el-link type="primary" @click="copyTestName(scope.row.idx)">paste</el-link>&nbsp;
                    </template>
                </el-table-column>
            </el-table>
        </el-col>
        <!-- 细节展示区 -->
        <el-col :span="16" class="detail-card">
            222
        </el-col>
    </el-row>
</template>

<style scoped>
.main-card {
    padding: 10px 0 0 10px;
}
.detail-card {
    padding: 10px 10px 0 5px;
}
.stat-line {
    margin-right: 20px;
    font-size: 12px;
}

</style>

<script setup lang="ts">
import { ref, type Ref, watch, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { pageTestRun } from '@/router/routes'
import { ApiTestList, ApiTestRun } from '@/api/test';
import type { RunStatus } from '@/api/test';
import { clipboardWriteText } from "@/utils/utils";
import { ElMessage } from "element-plus";

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
type Status = 'WAIT' | 'RUNNING' | 'NET_ERR' | RunStatus
type Test = {
    idx: number
    name: string
    status: Status
    content: string
    code: string
}

const tests: Ref<Test[]> = ref([])
const waitintIndexes: Ref<number[]> = ref([])
const statusCount = ref({
    'WAIT': 0,
    'RUNNING': 0,
    'NET_ERR': 0,
    // RunStatus
    'PASS':0,
    'BORK':0,
    'FAIL':0,
    'WARN':0,
    'LEAK':0,
    'XFAIL':0,
    'XLEAK':0,
    'SKIP':0,
    'SLOW':0,
})

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
    statusCount.value.WAIT = waitintIndexes.value.length

    runNext() // 触发执行 loop
})

// 更新状态
function updateTest(test: Test, status: Status, content: string) {
    statusCount.value[test.status]--
    statusCount.value[status]++
    test.status = status
    test.content = content
    test.code = ""
}

// 执行指定case
function runTest(index: number) {
    const test = tests.value[index]
    if (!test) {
        return
    }
    if (test.status === 'RUNNING') {
        return
    }
    updateTest(test, 'RUNNING', '')

    // api
    ApiTestRun({
        src: src,
        path: test.name,
    }).then(
        res => {
            if (res.code !== 0) {
                updateTest(test, 'FAIL',  "执行失败: error=" + res.error)
                return
            }

            updateTest(test, res.data.status, res.data.info)
            // test.code = res.data.code
            // test.output = res.data.output
            // test.expect = res.data.expect
            runNext()
        },
        () => {
            updateTest(test, 'NET_ERR', '调用 url 失败')
            runNext()
        }
    )
}

// 遍历待处理case
const runNext = () => {
    if (statusCount.value.WAIT < 10 && waitintIndexes.value.length > 0) {
        const index = waitintIndexes.value.shift()
        if (index) {
            runTest(index)
        }
    }
}

// 状态栏
const statLine = computed(() => {
    return "total: " + tests.value.length +
        ", wait: " + statusCount.value.WAIT +
        ", pass: " + statusCount.value.PASS +
        ", skip: " + statusCount.value.SKIP +
        ", fail: " + statusCount.value.FAIL +
        ", net_err: " + statusCount.value.NET_ERR
})

// 复制case名
function copyTestName(index: number) {
    const test = tests.value[index]
    if (!test) {
        return
    }
    clipboardWriteText(test.name).then(
        (msg) => { 
            ElMessage({
                message: msg,
                type: 'success',
            })
        },
        (msg) => { 
            ElMessage({
                message: msg,
                type: 'warning',
            })
         },
    )
}

// 列表筛选数据
let lockListIndexes: number[] = []
const showTests = computed(() => {
    return tests.value
})

// todo

</script>