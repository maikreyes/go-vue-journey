import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

export default defineConfig(({ mode }) => {
  const envDir = path.resolve(process.cwd(), '..');
  const env = loadEnv(mode, envDir, 'VITE_');

  return {
    envDir: envDir,
    plugins: [vue(), tailwindcss()],
    
   
    define: {
      __API_URL__: JSON.stringify(env.VITE_API_URL),
      __TOP_URL__: JSON.stringify(env.VITE_API_URL + env.VITE_TOP_STOKS_ENDPOINT),
      __LINKEDIN_URL__: JSON.stringify(env.VITE_LINKEDIN_URL),
      __GITHUB_URL__: JSON.stringify(env.VITE_GITHUB_URL)
    }
  }
})