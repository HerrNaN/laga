/** @type {import("@sveltejs/vite-plugin-svelte").SvelteConfig} */
export default {
  compilerOptions: {
    warningFilter: (warning) => !ignoreWarning[warning.code]?.(warning),
  },
};

const ignoreWarning = {
  a11y_no_static_element_interactions: (w) => w.message.startsWith("`<wa-"),
  a11y_click_events_have_key_events: () => true,
};
