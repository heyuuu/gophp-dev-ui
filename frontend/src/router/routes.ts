
export const BaseUri = import.meta.env.BASE_URL

export function routeUri(uri: string): string {
    if (uri.startsWith('/')) {
        uri = uri.slice(1)
    }
    return BaseUri + uri
}

// pages 

export function pageTestList(src: string, path: string): string {
    return routeUri('test/list?src=' + src + '&path=' + path)
}

export function pageTestRun(src: string, path: string): string {
    return routeUri('test/run?src=' + src + '&path=' + path)
}