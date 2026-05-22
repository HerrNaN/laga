<script lang="ts">
    import { items } from "../stores";
    import type { Item } from "./item";
    import SwipeAction from "./SwipeAction.svelte";
    import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
    import "@awesome.me/webawesome/dist/components/icon/icon.js";

    type Props = {
        item: Item;
    };

    let { item }: Props = $props();

    const handleToggle = () => {
        items.toggleItem(item.id);
    };
</script>

<li>
    <SwipeAction ontrigger={() => items.deleteItem(item.id)}>
        {#snippet underlay()}
            <div class="underlay">
                <wa-icon name="trash"></wa-icon>
            </div>
        {/snippet}

        <div class="content">
            <wa-checkbox checked={item.checked} onchange={handleToggle}
                >{item.text}</wa-checkbox
            >
        </div>
    </SwipeAction>
</li>

<style>
    li {
        height: 3.25rem;
        display: flex;
        align-items: center;
        justify-items: stretch;
        margin-inline: var(--wa-space-2xs);
    }

    .content {
        position: relative;
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        padding-inline: var(--wa-space-s);
        background-color: var(--wa-color-surface-default);
        border-color: var(--wa-color-neutral-border-normal);
        border-radius: var(--wa-border-width-s);
        border-style: var(--wa-border-style);
        border-width: var(--wa-border-width-s);

        &:has(wa-checkbox:state(checked)) {
            background-color: var(--wa-color-neutral-fill-quiet);

            & wa-checkbox::part(label) {
                text-decoration: line-through;
                color: var(--wa-color-neutral-on-quiet);
            }
        }
    }

    .underlay {
        position: absolute;
        inset: 0;
        display: flex;
        align-items: center;
        padding: var(--wa-space-s) var(--wa-space-m);
        justify-content: end;
        background-color: var(--wa-color-danger-fill-loud);
        color: var(--wa-color-brand-on-loud);
        border-radius: var(--wa-border-width-s);
        border-style: var(--wa-border-style);
        border-width: var(--wa-border-width-s);
    }

    wa-checkbox {
        width: 100%;
        -webkit-tap-highlight-color: transparent;
    }
</style>
