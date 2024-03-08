import axios from "axios";

export const axiosInstance = axios.create({
    baseURL: "/api"
})

export type ApiResult<T=any> = {
    code: number,
    error: string,
    data: T,
}