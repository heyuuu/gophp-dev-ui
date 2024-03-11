import { apiGet, apiPost } from './base'
import type { ApiResult } from './base'

type ApiTestPathListParam = {
  src: string
}
type ApiTestPathListResult = ApiResult<{
  count: number
  list: string[]
}>
export const apiTestPathList = async function (params: ApiTestPathListParam) {
  return await apiGet<ApiTestPathListResult>('test/path_list', params)
}

type ApiTestListParam = {
  src: string
  path?: string
  offset?: number
  limit?: number
}
type ApiTestListResult = ApiResult<{
  list: string[]
  count: number
  total: number
}>
export const apiTestList = async function (params: ApiTestListParam) {
  return await apiGet<ApiTestListResult>('test/list', params)
}

type ApiTestDetailParam = {
  src: string
  path: string
}
type ApiTestDetailResult = ApiResult<{
  src: string
  path: string
  sections: Sections
}>
export const apiTestDetail = async function (params: ApiTestDetailParam) {
  return await apiGet<ApiTestDetailResult>('test/detail', params)
}

export type Sections = { [K: string]: string }
export type RunStatus =
  | 'PASS'
  | 'BORK'
  | 'FAIL'
  | 'WARN'
  | 'LEAK'
  | 'XFAIL'
  | 'XLEAK'
  | 'SKIP'
  | 'SLOW'

type ApiTestRunParam = {
  src: string
  path: string
}
type ApiTestRunResult = ApiResult<{
  fileName: string
  filePath: string

  code: string
  expect: string

  status: RunStatus
  output: string
  info: string
  useTime: number
}>
export const apiTestRun = async function (params: ApiTestRunParam) {
  return await apiPost<ApiTestRunResult>('test/run', params)
}

type ApiTestRunCustomParam = {
  src: string
  path?: string
  sections: Sections
}
type ApiTestRunCustomResult = ApiTestRunResult
export const apiTestRunCustom = async function (params: ApiTestRunCustomParam) {
  return await apiPost<ApiTestRunCustomResult>('test/run_custom', params)
}
