import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

function buildCsp(isDev) {
  const directives = [
    "default-src 'self'",
    "base-uri 'self'",
    "frame-ancestors 'none'",
    "form-action 'self'",
    "object-src 'none'",
    "img-src 'self' data: blob:",
    "font-src 'self' data:",
    "style-src 'self' 'unsafe-inline'",
    isDev
      ? "script-src 'self' 'unsafe-inline'"
      : "script-src 'self'",
    isDev
      ? "connect-src 'self' ws: wss: http://localhost:8081"
      : "connect-src 'self'",
  ]

  return directives.join('; ')
}

export default defineConfig(({ mode }) => {
  const isDev = mode === 'development'
  const cspHeader = buildCsp(isDev)

  return {
    plugins: [react()],
    envDir: path.resolve(__dirname, '..'),
    server: {
      port: 5173,
      headers: {
        'Content-Security-Policy': cspHeader,
      },
      proxy: {
        '/api': {
          target: 'http://localhost:8081',
          changeOrigin: true,
        }
      }
    },
    preview: {
      headers: {
        'Content-Security-Policy': buildCsp(false),
      },
    },
  }
})
