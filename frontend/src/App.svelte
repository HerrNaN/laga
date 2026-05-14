<script lang="ts">
    import { onMount } from "svelte";
    import { items } from "./stores";
    import { getAllItems } from "./db";
    import AddItem from "./lib/AddItem.svelte";
    import ItemList from "./lib/ItemList.svelte";

    onMount(async () => {
        const dbItems = await getAllItems();
        items.hydrate(dbItems);
    });
</script>

<main>
    <header>
        <h1>laga</h1>
    </header>
    <AddItem />
    <ItemList items={$items} />
</main>

<style>
    main {
        max-width: 100%;
        min-height: 100dvh;
        background-color: var(--wa-color-surface-default);
    }

    header {
        padding: var(--wa-space-m);
        border-bottom: 1px solid var(--wa-color-surface-border);
    }

    h1 {
        margin: 0;
        font-size: var(--wa-font-size-xl);
        color: var(--wa-color-text-normal);
    }
</style>
