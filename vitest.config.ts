import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'happy-dom',
    globals: true,
    coverage: {
      provider: 'v8',
      reporter: ['text', 'lcov'],
      all: true,
      include: ['src/**/*.{ts,vue}'],
      exclude: [
        'src/**/__tests__/**',
        'src/**/*.d.ts',
        'src/types',
        'src/utils/mcSounds.ts',
        'src/utils/digSound.ts',
        'src/utils/mcTransition.ts',
        'src/components/BlockBreakOverlay.vue',
        'src/components/MinecraftDepthBackground.vue',
      ],
    },
  },
})