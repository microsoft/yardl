# Documentation Website

This site uses [VitePress](https://vitepress.dev/) to generate a static site for documentation. To run it locally, run

```bash
$ just start-docs-website
```

In most cases, the page in the browser should update when saving the markdown files. But changes to `.vitepress/config.mts` require a page reload.

Use Node.js 24, as in CI. To install the locked dependencies and check the site:

```bash
cd docs
npm ci
npm audit
npm run docs:build
```

Dependency upgrades must update both `package.json` and `package-lock.json`.
The lockfile omits registry-specific tarball URLs so installs can use the
configured npm registry or proxy while still verifying package integrity hashes.
CI builds the documentation on pull requests before publishing it from `main`.
The site currently uses VitePress 2 alpha because VitePress 1 depends on the
unsupported Vite 5 release line. Keep the prerelease pinned and validate future
upgrades with a full build and a check of navigation and local search.
