import dayjs from 'dayjs'

export function formatDateTime(value: string | null): string {
  if (!value) return '—'
  const d = dayjs(value)
  return d.isValid() ? d.format('YYYY-MM-DD HH:mm') : '—'
}

export function formatRelative(value: string | null): string {
  if (!value) return '—'
  const d = dayjs(value)
  if (!d.isValid()) return '—'
  const diffMinutes = dayjs().diff(d, 'minute')
  if (diffMinutes < 1) return '刚刚'
  if (diffMinutes < 60) return `${diffMinutes} 分钟前`
  const diffHours = dayjs().diff(d, 'hour')
  if (diffHours < 24) return `${diffHours} 小时前`
  const diffDays = dayjs().diff(d, 'day')
  if (diffDays < 30) return `${diffDays} 天前`
  return d.format('YYYY-MM-DD')
}

// 中文按字符计，英文按单词计，用于估算阅读时长
export function estimateReadingStats(markdown: string): { chars: number; minutes: number } {
  const chars = markdown.replace(/\s/g, '').length
  const minutes = Math.max(1, Math.round(chars / 400))
  return { chars, minutes }
}

export function slugify(title: string): string {
  return title
    .trim()
    .toLowerCase()
    .replace(/[^\p{Letter}\p{Number}]+/gu, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, 100)
}
