<script lang="ts">
    import { onMount } from "svelte";
    import WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
    import "@awesome.me/webawesome/dist/components/button/button.js";
    import "@awesome.me/webawesome/dist/components/icon/icon.js";
    import AddItem from "./AddItem.svelte";
    import ItemList from "./ItemList.svelte";
    import { listStore } from "./store";
    import { type Item } from "./list";
    import EditItem from "./EditItem.svelte";
    import { route, navigate } from "../../router";

    const { activeList, lists, setActiveList } = listStore;

    let editItemDialog = $state<WaDialog>();
    let itemToEdit = $state<Item>();
    const onClickEditItem = (item: Item) => {
        if (editItemDialog) editItemDialog.open = true;
        itemToEdit = item;
    };

    onMount(() => {
        const id = route.params.id;
        if (id && $lists.some((list) => list.id === id)) {
            setActiveList(id);
        } else {
            navigate("/lists", { replace: true });
        }
    });
</script>

<article>
    <header>
        <wa-button
            appearance="plain"
            aria-label="Back"
            onclick={() => navigate(-1)}
        >
            <wa-icon name="arrow-left"></wa-icon>
        </wa-button>
        <h2>{$activeList.name}</h2>
    </header>
    <div class="list-area">
        <ItemList items={$activeList.items} {onClickEditItem} />
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
        height: 2rem;
        z-index: 1;
        position: sticky;
        top: 0;
        display: flex;
        align-items: center;
        gap: var(--wa-space-s);
        padding: var(--wa-space-m);
        background-color: var(--wa-color-surface-raised);
        box-shadow: 0 var(--wa-shadow-offset-x-m) var(--wa-shadow-offset-y-m)
            var(--wa-color-shadow);
    }

    h2 {
        flex: 1;
        margin: 0;
        font-size: var(--wa-font-size-xl);
        font-weight: var(--wa-font-weight-normal);
        color: var(--wa-color-text-normal);
    }
</style>
