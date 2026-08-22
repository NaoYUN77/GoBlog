import type { InjectionKey } from 'vue'

export const MOBILE_MENU_KEY: InjectionKey<() => void> = Symbol('open-mobile-menu')
