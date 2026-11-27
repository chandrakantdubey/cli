import { readFileSync, readdirSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { defineConfig, type Plugin } from 'vite'

const indexId = '\0virtual:cli-silly-index'
const courseRoot = new URL('./course/', import.meta.url)
const referenceRoot = new URL('./src/data/modules/', import.meta.url)

function courseIndex(): Plugin {
  return {
    name: 'cli-silly-course-index',
    resolveId(id) {
      if (id === 'virtual:cli-silly-index') return indexId
    },
    load(id) {
      if (id !== indexId) return
      const read = (url: URL) => {
        this.addWatchFile(fileURLToPath(url))
        return JSON.parse(readFileSync(url, 'utf8'))
      }
      const catalog = read(new URL('catalog.json', courseRoot)) as {
        workshops: { id: string; title: string; lessons: string[] }[]
      }
      const workshops = catalog.workshops.map((group) => ({
        ...group,
        lessons: group.lessons.map((lessonId) => {
          const { id, title, reviewOf, referenceId } = read(new URL(`workshops/${lessonId}/lesson.json`, courseRoot))
          if (id !== lessonId) throw new Error(`Unexpected manifest ID for ${lessonId}`)
          return { id, title, reviewOf, referenceId }
        }),
      }))
      const promoted = new Map<string, string>()
      for (const lesson of workshops.flatMap((group) => group.lessons)) {
        if (!lesson.referenceId) continue
        if (promoted.has(lesson.referenceId)) throw new Error(`Duplicate reference ID: ${lesson.referenceId}`)
        promoted.set(lesson.referenceId, lesson.id)
      }
      const topics = readdirSync(referenceRoot)
        .filter((name) => name.endsWith('.json'))
        .sort((a, b) => parseInt(a) - parseInt(b))
        .map((name) => {
          const { id, title, content } = read(new URL(name, referenceRoot))
          const lessonId = promoted.get(id)
          if (!lessonId && !content) throw new Error(`Missing reading for reference ${id}`)
          promoted.delete(id)
          return { id, title, file: name, lessonId }
        })
      if (promoted.size > 0) throw new Error(`Unknown reference IDs: ${[...promoted.keys()].join(', ')}`)
      return `export const workshops = ${JSON.stringify(workshops)};\nexport const topics = ${JSON.stringify(topics)};`
    },
    handleHotUpdate({ file, server }) {
      const name = file.replaceAll('\\', '/')
      if (!name.endsWith('/course/catalog.json') && !name.endsWith('/lesson.json') &&
          !(name.endsWith('.json') && name.includes('/src/data/modules/'))) return
      const index = server.moduleGraph.getModuleById(indexId)
      if (index) {
        server.moduleGraph.invalidateModule(index)
        server.ws.send({ type: 'full-reload' })
        return []
      }
    },
  }
}

export default defineConfig({
  plugins: [react(), tailwindcss(), courseIndex()],
})
