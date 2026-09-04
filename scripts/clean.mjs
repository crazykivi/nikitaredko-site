import { rmSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..')

const paths = [
  resolve(root, 'dist'),
  resolve(root, 'backend', 'dist'),
  resolve(root, 'backend', 'server'),
  resolve(root, 'backend', 'server.exe')
]

for (const p of paths) {
  rmSync(p, { recursive: true, force: true })
  console.log(`Removed ${p}`)
}
