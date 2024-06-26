import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// 默认的 sanitizeFileName 实现 (https://github.com/rollup/rollup/blob/master/src/utils/sanitizeFileName.ts)
// https://datatracker.ietf.org/doc/html/rfc2396
// eslint-disable-next-line no-control-regex
const INVALID_CHAR_REGEX = /[\u0000-\u001F"#$&*+,:;<=>?[\]^`{|}\u007F]/g
const DRIVE_LETTER_REGEX = /^[a-z]:/i

export function sanitizeFileName(name: string): string {
  const match = DRIVE_LETTER_REGEX.exec(name)
  const driveLetter = match ? match[0] : ''

  // A `:` is only allowed as part of a windows drive letter (ex: C:\foo)
  // Otherwise, avoid them because they can refer to NTFS alternate data streams.
  return driveLetter + name.slice(driveLetter.length).replace(INVALID_CHAR_REGEX, '_')
}

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true
      }
    }
  },
  base: '/static/',
  build: {
    outDir: '../static/',
    rollupOptions: {
      output: {
        sanitizeFileName(name: string): string {
          let newName = sanitizeFileName(name)
          if (newName[0] === '_') {
            newName = '0' + newName
          }

          console.log({ name, newName })
          return newName
        }
      }
    }
  }
})
