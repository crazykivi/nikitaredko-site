import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { execSync } from 'child_process'

function getGitInfo() {
  try {
    const commitHash = execSync('git rev-parse HEAD').toString().trim()
    const commitShort = execSync('git rev-parse --short HEAD').toString().trim()
 
    let currentTag = ''
    try {
      currentTag = execSync('git describe --tags --exact-match 2>/dev/null').toString().trim()
    } catch {
      currentTag = ''
    }
    
    return {
      commitHash,
      commitShort,
      gitTag: currentTag,
      isRelease: !!currentTag
    }
  } catch (error) {
    console.warn('Failed to get git info:', error)
    return {
      commitHash: 'unknown',
      commitShort: 'unknown',
      gitTag: '',
      isRelease: false
    }
  }
}

const gitInfo = getGitInfo()

export default defineConfig({
  plugins: [vue()],
  define: {
    'import.meta.env.VITE_GIT_COMMIT': JSON.stringify(gitInfo.commitHash),
    'import.meta.env.VITE_GIT_COMMIT_SHORT': JSON.stringify(gitInfo.commitShort),
    'import.meta.env.VITE_GIT_TAG': JSON.stringify(gitInfo.gitTag),
    'import.meta.env.VITE_GIT_IS_RELEASE': JSON.stringify(gitInfo.isRelease),
    'import.meta.env.VITE_REPO_URL': JSON.stringify('https://github.com/crazykivi/nikitaredko-site'),
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