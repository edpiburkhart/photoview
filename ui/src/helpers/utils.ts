/* eslint-disable @typescript-eslint/no-explicit-any */

export interface DebouncedFn<F extends (...args: any[]) => any> {
  (...args: Parameters<F>): void
  cancel(): void
}

export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number,
  triggerRising?: boolean
): DebouncedFn<T> {
  let timeout: number | undefined = undefined

  const debounced = (...args: Parameters<T>) => {
    if (timeout) {
      clearTimeout(timeout)
      timeout = undefined
    } else if (triggerRising) {
      func(...args)
    }

    timeout = window.setTimeout(() => {
      timeout = undefined
      func(...args)
    }, wait)
  }

  debounced.cancel = () => {
    clearTimeout(timeout)
    timeout = undefined
  }

  return debounced
}

export function isNil(value: any): value is undefined | null {
  return value === undefined || value === null
}
