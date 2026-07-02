import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  define: {
    __COMMIT_HASH__: JSON.stringify(process.env.VITE_GIT_COMMIT || 'unknown'),
    __COMMIT_SHORT__: JSON.stringify(process.env.VITE_GIT_COMMIT_SHORT || 'unknown'),
    __GIT_TAG__: JSON.stringify(process.env.VITE_GIT_TAG || ''),
    __IS_RELEASE__: process.env.VITE_GIT_IS_RELEASE === 'true',
    __REPO_URL__: JSON.stringify('https://github.com/crazykivi/nikitaredko-site'),
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})