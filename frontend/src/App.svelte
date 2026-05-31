<script lang="ts">
    import { onMount } from "svelte";
    import { items } from "./stores";
    import { getAllItems } from "./db";
    import AddItem from "./lib/AddItem.svelte";
    import ItemList from "./lib/ItemList.svelte";

    onMount(async () => {
        const dbItems = await getAllItems();
        items.hydrate(dbItems);
        setTimeout(() => (document.body.style.height = "100dvh"), 10);
    });
</script>

<main>
    <header>
        <h1>Inköpslista</h1>
    </header>
    <div class="list-area">
        <ItemList items={$items} />
    </div>
    <AddItem />
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

    .list-area {
        flex: 1;
        overflow-y: auto;
        padding-bottom: var(--add-item-height);
        background-color: var(--wa-color-surface-lowered);
    }

    header {
        z-index: 1;
        position: sticky;
        top: 0;
        padding: var(--wa-space-m);
        background-color: var(--wa-color-surface-raised);
        box-shadow: 0 var(--wa-shadow-offset-x-m) var(--wa-shadow-offset-y-m)
            var(--wa-color-shadow);
    }

    h1 {
        margin: 0;
        font-size: var(--wa-font-size-xl);
        color: var(--wa-color-text-normal);
    }
</style>
