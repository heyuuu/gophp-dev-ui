<template>
    <el-row>
        <el-col :span="8">
            <el-form label-width="auto">
                <!-- sections -->
                <el-form-item v-for="(sec, index) in sections" :label="sec.name">
                    <template v-if="sec.type === 'input'">
                        <el-input v-model="sections[index].content"></el-input>
                    </template>
                    <template v-else-if="sec.type === 'text'">
                        <el-input v-model="sections[index].content" type="textarea"></el-input>
                    </template>
                    <template v-else>
                        <Editor :value="sec.content" :height='200' />
                    </template>
                </el-form-item>
                <!-- buttons -->
                <el-form-item>
                    <div style="display: flex; justify-content: flex-end; width: 100%;">
                        <el-button type="primary" @click="run">Run</el-button>
                    </div>
                </el-form-item>
            </el-form>
        </el-col>
        <el-col :span="16">222</el-col>
    </el-row>
</template>

<script setup lang="ts">
import { onMounted, ref, watch, type Ref } from "vue";
import { useRoute } from "vue-router";
import Editor from "@/components/Editor.vue"
import { ApiTestDetail } from '@/api/test';
import type { Sections } from '@/api/test';
import { pa, sr } from "element-plus/es/locales.mjs";

// uri 参数
const route = useRoute()
const src = (route.query.src || '') as string
const path = (route.query.path || '') as string

// sections
type SectionShowType = 'input' | 'text' | 'code'
type Section = {
    name: string,
    content: string,
    type: 'input' | 'text' | 'code',
}
const sections: Ref<Section[]> = ref([])
function sectionType(name: string): SectionShowType {
    switch (name) {
        case 'FILE':
            return 'code'
        default:
            return 'input'
    }
}

// 初始化case数据
onMounted(async () => {
    const rep = await ApiTestDetail({
        src: src,
        path: path
    })
    if (rep.code !== 0) {
        return
    }

    const sectionsData = rep.data.sections
    for (const name of Object.keys(sectionsData)) {
        sections.value.push({
            name: name,
            content: sectionsData[name],
            type: sectionType(name)
        })
    }

    console.log(sections.value)
})

// 执行
function run() {

}


</script>@/api/test@/api/test