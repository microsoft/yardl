import { defineConfig } from "vitepress";

const base = "/yardl"
// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Yardl",
  description: "Yardl Documentation",
  head: [["link", { rel: "icon", href: `${base}/favicon.ico` }]],
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "Home", link: "/" },
      {
        text: "Python Documentation",
        link: "/python/introduction",
        activeMatch: "/python/",
      },
      {
        text: "C++ Documentation",
        link: "/cpp/introduction",
        activeMatch: "/cpp/",
      },
      {
        text: "Matlab Documentation",
        link: "/matlab/introduction",
        activeMatch: "/matlab/",
      },
      {
        text: "Reference",
        link: "/reference/binary",
        activeMatch: "/reference/",
      },
    ],

    sidebar: {
      "/python/": [
        {
          text: "Getting Started (Python)",
          collapsed: false,
          items: [
            { text: "Introduction", link: "/python/introduction" },
            { text: "Installation", link: "/python/installation" },
            { text: "Quick Start", link: "/python/quickstart" },
          ],
        },
        {
          text: "Yardl Guide (Python)",
          collapsed: false,
          items: [
            { text: "Packages", link: "/python/packages" },
            { text: "The Yardl Language", link: "/python/language" },
            // { text: "Schema Evolution", link: "/python/evolution" },
          ],
        },
        { text: "Reference", link: "/reference/binary" },
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
            { text: "Configurable Arrays", link: "/cpp/arrays"},
            { text: "Schema Evolution", link: "/cpp/evolution" },
          ],
        },
        { text: "Reference", link: "/reference/binary" },
      ],
      "/matlab/": [
        {
          text: "Getting Started (Matlab)",
          collapsed: false,
          items: [
            { text: "Introduction", link: "/matlab/introduction" },
            { text: "Installation", link: "/matlab/installation" },
            { text: "Quick Start", link: "/matlab/quickstart" },
          ],
        },
        {
          text: "Yardl Guide (Matlab)",
          collapsed: false,
          items: [
            { text: "Packages", link: "/matlab/packages" },
            { text: "The Yardl Language", link: "/matlab/language" },
            // { text: "Schema Evolution", link: "/matlab/evolution" },
          ],
        },
        { text: "Reference", link: "/reference/binary" },
      ],
      "/reference/": [
        {
          text: "Reference",
          collapsed: false,
          items: [
            { text: "Binary Encoding Format", link: "/reference/binary" },
            { text: "NDJSON Encoding Format", link: "/reference/ndjson" },
            {
              text: "Protocol Schema JSON",
              link: "/reference/protocol-schema",
            },
          ],
        },
      ],
    },

    outline: {
      level: "deep",
    },

    socialLinks: [
      { icon: "github", link: "https://github.com/microsoft/yardl" },
    ],

    search: {
      provider: "local",
      options: {
        miniSearch: {
          options: {
            extractField(document, fieldName) {
              const fieldValue = document[fieldName];
              if (fieldName == "titles") {
                // Several documents have the same title in the Python and C++
                // documentation, which makes is hard to know which language a
                // search result is for. So we augment the "titles" field with
                // either C++ or Python if the document is under one of those paths.

                var documentId: string = document["id"];
                if (documentId.startsWith("/yardl/cpp")) {
                  // Include "C++"" in the search preview "path"
                  return ["C++"].concat(fieldValue);
                }

                if (documentId.startsWith("/yardl/python")) {
                  return ["Python"].concat(fieldValue);
                }

                if (documentId.startsWith("/yardl/matlab")) {
                  return ["Matlab"].concat(fieldValue);
                }
              }

              return fieldValue;
            },
          },
        },
      },
    },
  },
  base: base,
  srcExclude: ["README.md"],
});
