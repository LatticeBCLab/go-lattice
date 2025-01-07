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
      {
        text: "Guide",
        items: [
          { text: "Introduction", link: "/guide/introduction" },
          { text: "Getting Started", link: "/guide/getting-started" },
        ],
      },
      { text: "Developers", link: "/api-examples" },
      {
        text: "Resources",
        items: [
          { text: "Ledger", link: "/resources/ledger" },
          { text: "Contract", link: "/resources/contract" },
          { text: "gRPC", link: "/api-examples" },
        ],
      },
    ],
    search: {
      provider: "local",
    },
    /* sidebar: [
      {
        text: "Guide",
        items: [
          { text: "Introduction", link: "/guide/introduction" },
          { text: "Getting Started", link: "/guide/getting-started" },
        ],
        collapsed: true,
      },
      {
        text: "Developers",
      },
      {
        text: "Resources",
        items: [
          { text: "Ledger", link: "/markdown-examples" },
          { text: "Contract", link: "/api-examples" },
          { text: "gRPC", link: "/api-examples" },
        ],
        collapsed: true,
      },
    ], */

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
