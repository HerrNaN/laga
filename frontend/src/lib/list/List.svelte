<script lang="ts">
    import AddItem from "./AddItem.svelte";
    import ItemList from "./ItemList.svelte";
    import { items } from "./store";
</script>

<article>
    <header>
        <h1>Inköpslista</h1>
    </header>
    <div class="list-area">
        <ItemList items={$items} />
    </div>
    <AddItem />
</article>

<style>
    article {
        flex: 1;
        display: flex;
        flex-direction: column;
        background-color: var(--wa-color-surface-default);
        overflow-y: hidden;
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
