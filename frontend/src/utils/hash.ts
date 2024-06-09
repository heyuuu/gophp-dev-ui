import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

export function useRouteHash(options: {
  init?: (v: Record<string, any>) => void
  calc: () => Record<string, any>
}) {
  const route = useRoute()
  const router = useRouter()

  // 从路由获取参数
  if (options.init) {
    const initData = decodeHash(route.hash)
    options.init(initData)
  }

  // 监听参数并同步到 hash
  const data = computed(options.calc)
  watch(data, () => {
    const hash = encodeHash(data.value)
    router.push({ hash: hash })
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
