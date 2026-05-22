<script lang="ts">
    import type { Item } from "./item";
    import ItemRow from "./ItemRow.svelte";

    interface Props {
        items: Item[];
    }

    let { items: itemList }: Props = $props();

    const uncheckedItems = $derived(itemList.filter((item) => !item.checked));
    const checkedItems = $derived(itemList.filter((item) => item.checked));
</script>

<ul>
    {#each uncheckedItems as item (item.id)}
        <li>
            <ItemRow {item} />
        </li>
    {/each}
</ul>

{#if checkedItems.length > 0}
    <div class="checked-section">
        <h2>Checked</h2>
        <ul>
            {#each checkedItems as item (item.id)}
                <li>
                    <ItemRow {item} />
                </li>
            {/each}
        </ul>
    </div>
{/if}

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

    li {
        height: 3rem;
        margin-inline: var(--wa-space-2xs);
    }

    .checked-section {
        display: flex;
        flex-direction: column;
        gap: var(--wa-space-2xs);
    }

    h2 {
        margin: 0;
        padding: var(--wa-space-xs) var(--wa-space-m) 0;
        font-size: var(--wa-font-size-xs);
        text-transform: uppercase;
        letter-spacing: 0.05em;
        color: var(--wa-color-text-muted);
    }
</style>
