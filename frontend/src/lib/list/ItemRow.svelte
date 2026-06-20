<script lang="ts">
    import { items } from "./store";
    import type { Item } from "./list";
    import SwipeAction from "./SwipeAction.svelte";
    import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
    import "@awesome.me/webawesome/dist/components/icon/icon.js";

    type Props = {
        item: Item;
    };

    let { item }: Props = $props();

    const checked = $derived(item.checkedAt !== undefined);

    let isRemoving = $state(false);

    const handleToggle = () => {
        items.toggleItem(item.id);
    };

    const handleTransitionEnd = (e: TransitionEvent) => {
        if (e.propertyName === "opacity" && isRemoving) {
            items.deleteItem(item.id);
        }
    };

    const id = $props.id();
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

        <div
            class="content"
            onclick={handleToggle}
            role="checkbox"
            aria-checked={checked}
            tabindex="0"
        >
            <wa-checkbox {id} defaultChecked={checked}></wa-checkbox>
            <label for={id}>{item.text}</label>
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
        gap: var(--wa-space-s);
        align-items: center;
        padding-inline: var(--wa-space-2xs);
        background-color: var(--wa-color-surface-default);
        border-radius: var(--wa-border-width-l);

        &:has(wa-checkbox:state(checked)) {
            label {
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
        -webkit-tap-highlight-color: transparent;
        padding: 0.675rem;
        border-radius: var(--wa-border-width-l);
        background-color: var(--wa-color-neutral-fill-quiet);

        &:state(checked) {
            background-color: var(--wa-color-success-fill-quiet);
            &::part(control) {
                color: var(--wa-color-success-on-loud);
                background-color: var(--wa-color-success-fill-loud);
                border-color: var(--wa-color-success-fill-loud);
            }
        }

        &::part(control) {
            margin-inline-end: 0;
        }
    }
</style>
