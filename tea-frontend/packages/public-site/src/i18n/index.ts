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
