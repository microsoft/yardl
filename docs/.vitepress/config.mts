import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Yardl",
  description: "Yardl Documentation",
  head: [["link", { rel: "icon", href: "/favicon.ico" }]],
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
        _render(src, env, md) {
          switch (env.relativePath.split("/")[0]) {
            case "cpp":
              var searchDiscriminator = "C++";
              break;
            case "python":
              var searchDiscriminator = "Python";
              break;
            default:
              var searchDiscriminator = "";
          }

          if (searchDiscriminator != "") {
            // HACK: Insert a fake title of {searchDiscriminator}
            // and demote all other headings to be under it.
            // The objective is to make the search results show
            // the top-level directory. e.g.
            // Python > The Yardl Language > Maps
            // c++ > The Yardl Language > Maps

            src = src
              .replace(/^#+ .+$/gm, `#$&`)
              .replace(/^## .+$/m, `# ${searchDiscriminator}\n$&`);
          }

          return md.render(src, env);
        },

        miniSearch: {
          searchOptions: {
            boostDocument(documentId, term, storedFields:any) {
              // filter out results that match the fake title we inserted above.
              if (storedFields.titles.length == 0) {
                switch (storedFields.title) {
                  case "C++":
                  case "Python":
                    return 0;
                }
              }
              return 1;
            },

            filter(result) {
              console.log("filtering!");
              return true;
            },
          },
        },
      },
    },
  },
  base: "/yardl/",
  srcExclude: ["README.md"],
});
