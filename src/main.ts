import { createApp } from 'vue'
import { createHead } from '@unhead/vue/client' 
import './styles/style.css'
import './styles/markdown.css'
import App from './App.vue'
import router from './router'

const app = createApp(App)
const head = createHead()

app.use(head)
app.use(router)
app.mount('#app')