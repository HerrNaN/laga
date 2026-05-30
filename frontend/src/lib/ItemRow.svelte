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

    let isRemoving = $state(false);

    const handleToggle = () => {
        items.toggleItem(item.id);
    };

    const handleTransitionEnd = (e: TransitionEvent) => {
        if (e.propertyName === "opacity" && isRemoving) {
            items.deleteItem(item.id);
        }
    };
</script>

<div
    class="container"
    class:removing={isRemoving}
    ontransitionend={handleTransitionEnd}
>
    <SwipeAction ontrigger={() => (isRemoving = true)} disabled={isRemoving}>
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
</div>

<style>
    .container {
        height: 100%;
        transition: opacity var(--wa-transition-normal) ease;

        &.removing {
            opacity: 0;
            pointer-events: none;
        }
    }

    .content {
        position: relative;
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        padding-inline: var(--wa-space-s);
        background-color: var(--wa-color-surface-default);
        /*border-color: var(--wa-color-neutral-border-normal);*/
        border-radius: var(--wa-border-width-l);
        /*border-style: var(--wa-border-style);
        border-width: var(--wa-border-width-s);*/

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
        border-radius: var(--wa-border-width-l);
        border-style: var(--wa-border-style);
        border-width: var(--wa-border-width-s);
    }

    wa-checkbox {
        width: 100%;
        -webkit-tap-highlight-color: transparent;
    }
</style>
