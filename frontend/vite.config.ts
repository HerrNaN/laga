import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import wasm from "vite-plugin-wasm";

// https://vite.dev/config/
export default defineConfig({
  build: {
    outDir: "../backend/internal/app/web/dist",
    emptyOutDir: true,
  },
  plugins: [svelte(), wasm()],
});
