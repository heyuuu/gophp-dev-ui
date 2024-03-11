export function clipboardWriteText(text: string): Promise<string> {
  if (navigator.clipboard) {
    const cb = navigator.clipboard
    return cb.writeText(text).then(
      () => Promise.resolve('copy success'),
      () => Promise.reject('copy fail')
    )
  }
  return Promise.reject('copy unsupported')
}
