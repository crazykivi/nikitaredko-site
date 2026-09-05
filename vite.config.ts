import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { VitePWA } from 'vite-plugin-pwa'
const isCapacitorBuild = process.env.CAPACITOR_BUILD === 'true'
export default defineConfig({
  plugins: [
    vue(),
    ...(isCapacitorBuild ? [] : [VitePWA({
      registerType: 'autoUpdate',
      injectRegister: null,

      includeAssets: [
        'favicon.svg',
        'robots.txt',
        'icons/*.png'
      ],

      manifest: {
        name: 'Nikita Redko — Блог',
        short_name: 'nikitaredko',
        description:
          'Блог Никиты Редко: статьи о разработке, проектах, мыслях вслух и всём подряд.',
        lang: 'ru',
        id: '/',
        start_url: '/',
        scope: '/',
        display: 'standalone',
        background_color: '#09090b',
        theme_color: '#09090b',
        icons: [
          {
            src: '/icons/icon-192.png',
            sizes: '192x192',
            type: 'image/png'
          },
          {
            src: '/icons/icon-512.png',
            sizes: '512x512',
            type: 'image/png'
          },
          {
            src: '/icons/icon-512-maskable.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'maskable'
          }
        ]
      },

      workbox: {
        globPatterns: ['**/*.{js,css,html,svg,webmanifest}'],
        navigateFallback: '/index.html',
        navigateFallbackDenylist: [
          /\/api\//,
          /^\/sitemap\.xml$/,
          /^\/robots\.txt$/
        ],

        cleanupOutdatedCaches: true,
        clientsClaim: true,
        skipWaiting: true,
        maximumFileSizeToCacheInBytes: 5 * 1024 * 1024,
        inlineWorkboxRuntime: true,

        runtimeCaching: [
          {
            urlPattern: /\/api\/(collections|articles\/structured|articles\/feed|about|uses)/,
            method: 'GET',
            handler: 'StaleWhileRevalidate',
            options: {
              cacheName: 'api-lists',
              cacheableResponse: {
                statuses: [0, 200]
              },
              expiration: {
                maxEntries: 100,
                maxAgeSeconds: 60 * 60 * 24 * 7
              }
            }
          },
          {
            urlPattern: /\/api\/articles\/(?!search)[^/?#]+/,
            method: 'GET',
            handler: 'NetworkFirst',
            options: {
              cacheName: 'api-article',
              networkTimeoutSeconds: 4,
              cacheableResponse: {
                statuses: [0, 200]
              },
              expiration: {
                maxEntries: 200,
                maxAgeSeconds: 60 * 60 * 24 * 30
              }
            }
          },
          {
            urlPattern: /\/api\/attachments\.redirect/,
            method: 'GET',
            handler: 'CacheFirst',
            options: {
              cacheName: 'attachments',
              cacheableResponse: {
                statuses: [0, 200]
              },
              expiration: {
                maxEntries: 100,
                maxAgeSeconds: 60 * 60 * 24 * 30
              }
            }
          },
          {
            urlPattern: /https:\/\/fonts\.googleapis\.com/,
            method: 'GET',
            handler: 'CacheFirst',
            options: {
              cacheName: 'google-fonts-cache',
              cacheableResponse: {
                statuses: [0, 200]
              },
              expiration: {
                maxEntries: 20,
                maxAgeSeconds: 60 * 60 * 24 * 365
              }
            }
          },
          {
            urlPattern: /https:\/\/fonts\.gstatic\.com/,
            method: 'GET',
            handler: 'CacheFirst',
            options: {
              cacheName: 'gstatic-fonts-cache',
              cacheableResponse: {
                statuses: [0, 200]
              },
              expiration: {
                maxEntries: 20,
                maxAgeSeconds: 60 * 60 * 24 * 365
              }
            }
          }
        ]
      },

      devOptions: {
        enabled: false
      }
    })]),
  ],
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
