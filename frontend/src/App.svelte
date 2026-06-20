<script lang="ts">
    import List from "./lib/list/List.svelte";
    import { pwaHack } from "./lib/utils/pwa-hack";
    import NavBar from "./lib/NavBar.svelte";
</script>

<main {@attach pwaHack("100dvh")}>
    <List />
    <NavBar />
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
    }
</style>
