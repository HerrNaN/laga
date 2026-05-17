<script lang="ts">
    import WaInput from "@awesome.me/webawesome/dist/components/input/input.js";
    import { items } from "../stores";

    let input = $state<WaInput>();

    function handleAdd() {
        if (!input) return;

        const text = input.value?.trim();
        if (text) {
            items.addItem(text);
            input.value = "";
        }
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Enter") {
            handleAdd();
        }
    }
</script>

<div class="input-wrapper">
    <div
        class="backdrop"
        role="presentation"
        onclick={() => input?.blur()}
    ></div>
    <section>
        <wa-input
            placeholder="e.g. 2 kg potatoes"
            bind:this={input}
            onkeydown={handleKeydown}
        ></wa-input>
    </section>
</div>

<style>
    .backdrop {
        position: fixed;
        inset: 0;
        background-color: rgba(0, 0, 0, 0.4);
        pointer-events: none;
        opacity: 0;
        transition: opacity var(--wa-transition-fast) ease;
    }

    .input-wrapper:focus-within .backdrop {
        opacity: 1;
        pointer-events: auto;
    }

    section {
        position: relative;
        padding: var(--wa-space-m);
        box-shadow: 0 -0.25rem 0.75rem var(--wa-color-shadow);
        background-color: var(--wa-color-surface-default);
    }
</style>
