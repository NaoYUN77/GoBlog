import dayjs from 'dayjs'

export function formatDate(iso: string): string {
  if (!iso) return ''
  return dayjs(iso).format('YYYY年M月D日')
}

export function formatViewCount(count: number): string {
  if (count >= 10000) {
    return `${(count / 10000).toFixed(1)}万`
  }
  return String(count)
}
