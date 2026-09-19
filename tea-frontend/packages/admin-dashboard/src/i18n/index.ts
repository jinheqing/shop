import { createI18n } from 'vue-i18n'
import en from './en.json'
import zh from './zh.json'

const saved = localStorage.getItem('locale') || ''
const browserLang = navigator.language.toLowerCase()
const defaultLocale = saved || (browserLang.startsWith('zh') ? 'zh' : 'en')

export const i18n = createI18n({
  legacy: false,
  locale: defaultLocale,
  fallbackLocale: 'en',
  messages: { en, zh },
})

export type AppLocale = 'zh' | 'en'

export function setLocale(locale: string) {
  i18n.global.locale.value = locale as AppLocale
  localStorage.setItem('locale', locale)
  // 同步 Element Plus 语言
  // 在组件中调用时需要自己处理 element-plus 的 locale
}
