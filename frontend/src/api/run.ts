import { apiGet, apiPost } from './base'
import type { ApiResult } from './base'

type apiRunCodeParam = {
  code: string
}
type apiRunCodeResult = ApiResult<{
  result: {
    type: string
    content: string
  }[]
}>
export const apiRunCode = async (params: apiRunCodeParam) => {
  return await apiPost<apiRunCodeResult>('run/code', params)
}
