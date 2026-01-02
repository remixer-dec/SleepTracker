import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import App from './App.vue'
import messages from './i18n/messages.js'
import './assets/styles.css'

function detectLanguage() {
  const browserLang = navigator.language || navigator.userLanguage
  const lang = browserLang.split('-')[0]
  if (['en', 'es', 'ru'].includes(lang)) {
    return lang
  }
  return 'en'
}

const i18n = createI18n({
  legacy: false,
  locale: detectLanguage(),
  fallbackLocale: 'en',
  messages
})

const pinia = createPinia()

const app = createApp(App)
app.use(pinia)
app.use(i18n)
app.mount('#app')
