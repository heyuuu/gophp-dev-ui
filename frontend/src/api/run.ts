import { apiGet, apiPost } from './base'

// api: GET /run/config
type apiRunConfigResult = {
  types: string[]
}
export const apiRunConfig = async () => {
  return await apiGet<apiRunConfigResult>('run/config', {})
}

// api: POST /run/code

type apiRunCodeParam = {
  code: string
}
type apiRunCodeResult = {
  result: {
    type: string
    content: string
  }[]
}
export const apiRunCode = async (params: apiRunCodeParam) => {
  return await apiPost<apiRunCodeResult>('run/code', params)
}
