import { axiosInstance, type ApiResult } from "./base";



type ApiTestsPathListParam = {
    src: String,
}
type ApiTestsPathListResult = ApiResult<{
    count: number,
    list: string[],
}>
export const ApiTestsPathList = async function(params: ApiTestsPathListParam) {
    return await axiosInstance.get<ApiTestsPathListResult>('tests/path_list', {params})
}



type ApiTestsParam = {
    src: String,
}
type ApiTestsResult = ApiResult<{
    count: number,
    list: string[],
}>
export const ApiTestsList = async function(params:ApiTestsParam) {
    return await axiosInstance.get<ApiTestsResult>("tests/list", {params})
}