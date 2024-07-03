export const baseUri = import.meta.env.BASE_URL

function routeUri(path: string, query?: Record<string, string>): string {
  let uri = baseUri

  // path
  if (path.startsWith('/')) {
    path = path.slice(1)
  }
  uri += path

  // query
  if (query) {
    Object.entries(query).forEach(([key, value], index) => {
      uri += (index === 0 ? '?' : '&') + key + '=' + encodeURIComponent(value)
    })
  }

  return uri
}

// pages

export function pageTestList(mode: string, root: string, path: string): string {
  return routeUri(mode + '/test/list', { root, path })
}

export function pageTestRun(mode: string, root: string, path: string): string {
  return routeUri(mode + '/test/run', { root, path })
}
