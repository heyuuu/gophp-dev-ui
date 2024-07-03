import { apiGet, apiPost } from './base'
import type { SectionMap } from '@/models/test'

// api: GET /test/config
type ApiTestConfigParam = {
  mode: string
}
type ApiTestConfigResult = {
  defaultTestRoot: string
}
export const apiTestConfig = async (params: ApiTestConfigParam) => {
  return await apiGet<ApiTestConfigResult>('test/config', params)
}

// api: GET /test/path_list
type ApiTestPathListParam = {
  mode: string
  root: string
}
type ApiTestPathListResult = {
  count: number
  list: string[]
}
export const apiTestPathList = async (params: ApiTestPathListParam) => {
  return await apiGet<ApiTestPathListResult>('test/path_list', params)
}

// api: GET /test/case_list
type ApiTestCaseListParam = {
  mode: string
  root: string
  path?: string
  offset?: number
  limit?: number
}
type ApiTestCaseListResult = {
  list: string[]
  count: number
  total: number
}
export const apiTestCaseList = async (params: ApiTestCaseListParam) => {
  return await apiGet<ApiTestCaseListResult>('test/case_list', params)
}

// api: GET /test/detail
type ApiTestDetailParam = {
  mode: string
  root: string
  path: string
}
type ApiTestDetailResult = {
  root: string
  path: string
  sections: SectionMap
}
export const apiTestDetail = async (params: ApiTestDetailParam) => {
  return await apiGet<ApiTestDetailResult>('test/detail', params)
}

// api: POST /test/run
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
  mode: string
  root: string
  path: string
}
type ApiTestRunResult = {
  fileName: string
  filePath: string

  code: string
  expect: string

  status: RunStatus
  output: string
  info: string
  useTime: number
}
export const apiTestRun = async (params: ApiTestRunParam) => {
  return await apiPost<ApiTestRunResult>('test/run', params)
}

// api: POST /test/run_custom
type ApiTestRunCustomParam = {
  mode: string
  root: string
  path?: string
  sections: SectionMap
}
type ApiTestRunCustomResult = ApiTestRunResult
export const apiTestRunCustom = async (params: ApiTestRunCustomParam) => {
  return await apiPost<ApiTestRunCustomResult>('test/run_custom', params)
}
