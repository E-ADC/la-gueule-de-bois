import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // API Go (spec : backend écoute sur :8080 en dev)
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      // Photos de soirées servies en statique par nginx en prod (voir spec) ;
      // même prestataire en dev pour tester l'upload sans reconfigurer les URLs.
      '/uploads': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
