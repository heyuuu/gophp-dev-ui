import { apiPost } from './base'

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
