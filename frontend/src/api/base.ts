import axios from 'axios';
import type { Method, AxiosResponse } from 'axios';

export const axiosInstance = axios.create({
    baseURL: "/api"
})

export type ApiResult<T=any> = {
    code: number,
    error: string,
    data: T,
}

export async function request<T, D = any>(method: Method, url: string, params: any, data?: D) {
    const rep = await axiosInstance.request<T>({
        method: method,
        url: url,
        params: params,
        data: data,
    })
    return rep.data
}