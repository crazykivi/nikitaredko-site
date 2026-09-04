import { cpSync, existsSync, rmSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..')
const src = resolve(root, 'dist')
const dest = resolve(root, 'backend', 'dist')

if (!existsSync(src)) {
  console.error('dist/ not found. Run "npm run build" (or "make build") first.')
  process.exit(1)
}

rmSync(dest, { recursive: true, force: true })
cpSync(src, dest, { recursive: true })
console.log('Copied dist/ -> backend/dist/')
