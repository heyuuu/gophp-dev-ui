import { computed, watch, type Ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

/**
 * 将指定参数记录到 route 的 hash 上，初始化时会应用值到对应 Ref 上
 * @param data 指定参数记录map. 键为存储key，值为对应 Ref
 */
export function useRouteHash(data: Record<string, Ref<any>>) {
  const route = useRoute()
  const router = useRouter()

  // 从路由获取参数
  const initData = decodeHash(route.hash)
  Object.entries(data).forEach(([key, ref]) => {
    ref.value = initData[key] || ref.value
  })

  // 监听参数并同步到 hash
  const watchData = computed(() => {
    return Object.fromEntries(Object.entries(data).map(([key, ref]) => [key, ref.value]))
  })
  watch(watchData, () => {
    const hash = encodeHash(watchData.value)
    router.push({ ...route, hash: hash })
  })
}

export function encodeHash(data: Record<string, any>): string {
  const text = JSON.stringify(data)

  console.log(`encode string[${text.length}]: ${text}`)
  const hash = base64Encode(text)
  console.log(`encoded string[${hash.length}]: ${hash}`)
  // todo encode
  return '#' + hash
}

export function decodeHash(hash: string): Record<string, any> {
  console.log(`decode hash: ${hash}`)

  let result: Object = {}

  if (hash.length > 1 && hash[0] === '#') {
    try {
      const text = base64Decode(hash.slice(1))
      const data = JSON.parse(text)
      if (data instanceof Object) {
        result = data
      }
    } catch {
      /* empty */
    }
  }

  console.log(`decoded hash: ${JSON.stringify(result)}`)
  return result
}

function base64Encode(raw: string): string {
  // 字符串转 codePoint 以支持多字节字符
  const bytes = new TextEncoder().encode(raw)
  const binString = Array.from(bytes, (byte) => String.fromCodePoint(byte)).join('')
  return btoa(binString)
}

function base64Decode(encoded: string): string {
  const binString = atob(encoded)
  const bytes = Uint8Array.from(binString, (m) => m.codePointAt(0) as number)
  return new TextDecoder().decode(bytes)
}
