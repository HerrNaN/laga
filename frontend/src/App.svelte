<script lang="ts">
    import List from "./lib/list/List.svelte";
    import { pwaHack } from "./lib/list/utils/pwa-hack";
</script>

<main {@attach pwaHack("100dvh")}>
    <List />
</main>

<style>
    :global(body) {
        --add-item-height: calc(var(--wa-space-m) * 2 + 3rem);
        margin: 0;
    }

    main {
        height: 100dvh;
        display: flex;
        flex-direction: column;
        background-color: var(--wa-color-surface-default);
    }
</style>
