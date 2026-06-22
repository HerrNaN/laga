<script lang="ts">
    import WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
    import AddItem from "./AddItem.svelte";
    import ItemList from "./ItemList.svelte";
    import { items } from "./store";
    import { type Item } from "./list";
    import EditItem from "./EditItem.svelte";

    let editItemDialog = $state<WaDialog>();
    let itemToEdit = $state<Item>();
    const onClickEditItem = (item: Item) => {
        if (editItemDialog) editItemDialog.open = true;
        itemToEdit = item;
    };
</script>

<article>
    <header>
        <h1>Inköpslista</h1>
    </header>
    <div class="list-area">
        <ItemList items={$items} {onClickEditItem} />
        <AddItem />
    </div>
</article>

<wa-dialog bind:this={editItemDialog} label="Edit Item" light-dismiss>
    {#if itemToEdit}
        <EditItem item={itemToEdit} />
    {/if}
    <wa-button slot="footer" variant="brand" data-dialog="close">
        Close
    </wa-button>
</wa-dialog>

<style>
    article {
        height: 100%;
        display: flex;
        flex-direction: column;
        background-color: var(--wa-color-surface-default);
        overflow-y: hidden;
    }

    .list-area {
        flex: 1;
        display: flex;
        flex-direction: column;
        overflow-y: auto;
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
