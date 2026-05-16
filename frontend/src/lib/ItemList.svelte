<script lang="ts">
    import type { Item } from "./item";
    import ItemRow from "./ItemRow.svelte";

    interface Props {
        items: Item[];
    }

    let { items: itemList }: Props = $props();
</script>

<ul>
    {#each itemList as item (item.id)}
        <ItemRow {item} />
    {/each}
</ul>

<style>
    ul {
        list-style: none;
        margin: 0;
        padding: 0;
        display: flex;
        flex-direction: column;
        gap: var(--wa-space-2xs);
        padding-block: var(--wa-space-2xs);
    }
</style>
