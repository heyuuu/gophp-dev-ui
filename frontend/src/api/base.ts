import axios from 'axios'
import type { Method } from 'axios'
import { ElMessage } from 'element-plus'

export const axiosInstance = axios.create({
  baseURL: '/api'
})

export type ApiResult<T = any> = {
  code: number
  error: string
  data: T
}

const showApiFailMessage = (errMsg: string) => {
  ElMessage({
    type: 'error',
    message: '请求失败: ' + errMsg,
    grouping: true,
    offset: 8
  })
}

export async function apiRequest<T>(
  method: Method,
  url: string,
  params: any,
  data?: any
): Promise<T> {
  const rep = await axiosInstance
    .request<ApiResult<T>>({
      method: method,
      url: url,
      params: params,
      data: data
    })
    .catch((err) => {
      showApiFailMessage(`请求错误: ${err}`)
      throw err
    })

  const result = rep.data
  if (result.code !== 0) {
    showApiFailMessage(result.error)
    throw new Error(result.error)
  }

  return result.data
}

export async function apiGet<T>(url: string, params: any) {
  return apiRequest<T>('GET', url, params)
}

export async function apiPost<T>(url: string, data: any) {
  return apiRequest<T>('POST', url, {}, data)
}
