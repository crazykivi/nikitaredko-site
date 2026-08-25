import { createApp } from 'vue'
import { createHead } from '@unhead/vue/client'
import './styles/style.css'
import './styles/markdown.css'
import App from './App.vue'
import router from './router'
const isCapacitor = import.meta.env.VITE_CAPACITOR_BUILD === 'true'

if (!import.meta.env.DEV && !isCapacitor) {
  import('virtual:pwa-register').then(({ registerSW }) => {
    registerSW({
      immediate: true,
      onOfflineReady() {
        console.info('[PWA] Приложение готово к офлайн-работе')
      },
      onRegisterError(error) {
        console.error('[PWA] Ошибка регистрации service worker:', error)
      }
    })
  })
}

const app = createApp(App)
const head = createHead()

app.use(head)
app.use(router)

app.mount('#app')