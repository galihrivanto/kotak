import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import yaml from 'js-yaml'
import fs from 'fs'

// https://vite.dev/config/
export default defineConfig(({ command, mode }) => {
  const config = yaml.load(fs.readFileSync('../config.yaml', 'utf8'))

  return {
    plugins: [react()],
    
    define: {
      'import.meta.env.VITE_API_HOST': JSON.stringify(config.http_server.api_host),
      'import.meta.env.VITE_API_BASE': JSON.stringify(config.http_server.api_base),
    }
  }
})