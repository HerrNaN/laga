<script lang="ts">
    import { items } from "../stores";
    import type { Item } from "./item";
    import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";

    interface Props {
        items: Item[];
    }

    let { items: itemList }: Props = $props();

    function handleToggle(id: string) {
        items.toggleItem(id);
    }
</script>

<ul>
    {#each itemList as item (item.id)}
        <li>
            <wa-checkbox
                checked={item.checked}
                onchange={() => handleToggle(item.id)}
            ></wa-checkbox>
            <span>{item.text}</span>
        </li>
    {/each}
</ul>

<style>
    ul {
        list-style: none;
        margin: 0;
        padding: 0;
    }

    li {
        display: flex;
        flex-direction: column;
        gap: var(--wa-space-2xs);
        padding-block: var(--wa-space-2xs);
    }
</style>
