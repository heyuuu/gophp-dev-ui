import { request } from "./base";
import type { ApiResult } from "./base";


type ApiTestPathListParam = {
    src: string,
}
type ApiTestPathListResult = ApiResult<{
    count: number,
    list: string[],
}>
export const ApiTestPathList = async function(params: ApiTestPathListParam) {
    return await request<ApiTestPathListResult>('GET', 'test/path_list', params)
}


type ApiTestListParam = {
    src: string,
    path?: string,
    offset?: number,
    limit?: number,
}
type ApiTestListResult = ApiResult<{
    list: string[],
    count: number,
    total: number,
}>
export const ApiTestList = async function(params:ApiTestListParam) {
    return await request<ApiTestListResult>('GET', "test/list", params)
}


type ApiTestDetailParam = {
    src: string,
    path: string,
}
type ApiTestDetailResult = ApiResult<{
    src: string,
    path: string,
    sections: object,
}>
export const ApiTestDetail = async function(params:ApiTestListParam) {
    return await request<ApiTestListResult>('GET', "test/detail", params)
}

type Sections = {[K:string]: string}

type ApiTestRunParam = {
    src: string,
    path: string,
}
type ApiTestRunResult = ApiResult<{
    src: string,
    path: string,
    sections: Sections,
}>
export const ApiTestRun = async function(params: ApiTestRunParam) {
    return await request<ApiTestRunResult>('POST', "test/detail", params)
}

type ApiTestRunCustomParam = {
    src: string,
    path?: string,
    sections: Sections,
}
type ApiTestRunCustomResult = ApiTestRunResult
export const ApiTestRunCustom = async function(params: ApiTestRunCustomParam) {
    return await request<ApiTestRunCustomResult>('GET', "test/detail", {
        src: params.src,
        path: params.path
    }, {
        sections: params.sections
    })
}