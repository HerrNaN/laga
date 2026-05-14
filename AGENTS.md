# Frontend
- **Prefer arrow syntx functions:**
  - DO: `const f = () => { ... }`
  - DONT: `function f() { ... }`
- **Utilize scoped styling:** Svelte scopes styles to the component files by default. Utilize this and avoid unnecessary classes that scope the component
  - DO: `section { ... }; input { ... }`
  - DONT: `.component__base { ... }; .component__input { ... }`
- **Locally import webawesome components when used:** In every file where a webawesome component is used (`<wa-...`) that component MUST be imported in that file.
