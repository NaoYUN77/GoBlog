import { defineStore } from 'pinia'

export type ToastKind = 'success' | 'error' | 'info'

export interface ToastItem {
  id: number
  kind: ToastKind
  message: string
}

let nextId = 1

export const useToastStore = defineStore('toast', {
  state: () => ({
    items: [] as ToastItem[],
  }),
  actions: {
    push(message: string, kind: ToastKind = 'info', duration = 3600) {
      const id = nextId++
      this.items.push({ id, kind, message })
      window.setTimeout(() => this.dismiss(id), duration)
      return id
    },
    success(message: string) {
      return this.push(message, 'success')
    },
    error(message: string) {
      return this.push(message, 'error', 5000)
    },
    info(message: string) {
      return this.push(message, 'info')
    },
    dismiss(id: number) {
      this.items = this.items.filter((item) => item.id !== id)
    },
  },
})
