import type { CapacitorConfig } from '@capacitor/cli'

const config: CapacitorConfig = {
  appId: 'ru.nikitaredko.app',
  appName: 'Nikita Redko',
  webDir: 'dist',
  server: {
    androidScheme: 'https',
  },
}

export default config