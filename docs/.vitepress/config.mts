import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Yardl",
  description: "Yardl Documentation",
  head: [["link", { rel: "icon", href: "favicon.ico" }]],
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "Home", link: "/" },
      { text: "Python Documentation", link: "/python/introduction", activeMatch: "/python/" },
      { text: "C++ Documentation", link: "/cpp/introduction", activeMatch: "/cpp/" },
      { text: "Reference", link: "/reference/binary", activeMatch: "/reference/" },
    ],

    sidebar: {
      "/python/": [
        {
          text: "Getting Started (Python)",
          collapsed: false,
          items: [
            { text: "Introduction", link: "/python/introduction" },
            { text: "Quick Start", link: "/python/quickstart" },
          ],
        },
        {
          text: "Yardl Guide (Python)",
          collapsed: false,
          items: [
            { text: "Packages", link: "/python/packages" },
            { text: "The Yardl Language", link: "/python/language" },
          ],
        },
        { text: "Reference", link:"/reference/binary"}
      ],
      "/cpp/": [
        {
          text: "Getting Started (C++)",
          collapsed: false,
          items: [
            { text: "Introduction", link: "/cpp/introduction" },
            { text: "Installation", link: "/cpp/installation" },
            { text: "Quick Start", link: "/cpp/quickstart" },
          ],
        },
        {
          text: "Yardl Guide (C++)",
          collapsed: false,
          items: [
            { text: "Packages", link: "/cpp/packages" },
            { text: "The Yardl Language", link: "/cpp/language" },
            { text: "Performance Tips", link: "/cpp/performance" },
          ],
        },
        { text: "Reference", link:"/reference/binary"}
      ],
      "/reference/": [
        {
          text: "Reference",
          collapsed: false,
          items: [
            { text: "Binary Encoding Format", link: "/reference/binary" },
            { text: "NDJSON Encoding Format", link: "/reference/ndjson" },
            { text: "Protocol Schema JSON", link: "/reference/protocol-schema" },
          ],
        },
      ],
    },

    outline: {
      level: "deep"
    },

    socialLinks: [
      { icon: "github", link: "https://github.com/microsoft/yardl" },
    ],

    search: {
      provider: "local",
    },
  },
  base: "/yardl/",
  srcExclude: ["README.md"]
});
