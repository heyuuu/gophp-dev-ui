import { apiGet, apiPost } from './base'

// api: GET /run/config
type ApiRunConfigParam = {
  mode: string
}
type ApiRunConfigResult = {
  types: string[]
}
export const apiRunConfig = async (params: ApiRunConfigParam) => {
  return await apiGet<ApiRunConfigResult>('run/config', params)
}

// api: POST /run/code

type ApiRunCodeParam = {
  mode: string
  code: string
}
type ApiRunCodeResult = {
  result: {
    type: string
    language: string
    content: string
  }[]
  error: string
}
export const apiRunCode = async (params: ApiRunCodeParam) => {
  return await apiPost<ApiRunCodeResult>('run/code', params)
}
