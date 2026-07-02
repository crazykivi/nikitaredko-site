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
    __COMMIT_HASH__: JSON.stringify(gitInfo.commitHash),
    __COMMIT_SHORT__: JSON.stringify(gitInfo.commitShort),
    __GIT_TAG__: JSON.stringify(gitInfo.gitTag),
    __IS_RELEASE__: gitInfo.isRelease,
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