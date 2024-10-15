import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Blockchain",
  description: "Lightweight Go library for integration with ZLattice clients.",
  themeConfig: {
    logo: "/logo.svg",
    logoLink: "/",
    siteTitle: "ZLattice",
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "Guide", link: "/markdown-examples" },
      { text: "Developers", link: "/api-examples" },
      {
        text: "Resources",
        items: [
          { text: "Ledger", link: "/markdown-examples" },
          { text: "Contract", link: "/api-examples" },
          { text: "gRPC", link: "/api-examples" },
        ],
      },
    ],
    search: {
      provider: "local",
    },
    sidebar: [
      {
        text: "Examples",
        items: [
          { text: "Markdown Examples", link: "/markdown-examples" },
          { text: "Runtime API Examples", link: "/api-examples" },
        ],
      },
    ],

    socialLinks: [
      { icon: "github", link: "https://github.com/LatticeBCLab/go-lattice" },
    ],

    footer: {
      message: "Released under the MIT License.",
      copyright: "Copyright © 2024 ZKJG",
    },
  },
  head: [
    [
      "link",
      {
        rel: "stylesheet",
        href: "https://unpkg.com/tailwindcss@2.0.4/dist/tailwind.min.css",
      },
    ],
  ],
  vite: {
    build: {
      rollupOptions: {
        output: {
          globals: {
            self: "globalThis",
          },
        },
      },
    },
    define: {
      self: "globalThis",
    },
  },
});
