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
            placeholder="t.ex. 2 kg potatis"
            bind:this={input}
            onkeydown={handleKeydown}
        ></wa-input>
    </section>
</div>

<style>
    .input-wrapper {
        position: fixed;
        bottom: 0;
        left: 0;
        right: 0;
        z-index: 10;
    }

    .backdrop {
        position: fixed;
        inset: 0;
        background-color: var(--wa-color-overlay-modal);
        pointer-events: none;
        opacity: 0;
        transition: opacity var(--wa-transition-fast) ease;
    }

    .input-wrapper:focus-within .backdrop {
        opacity: 1;
        pointer-events: auto;
    }

    wa-input {
        width: 100%;
    }

    section {
        position: relative;
        height: var(--add-item-height);
        padding: 0 var(--wa-space-m);
        display: flex;
        align-items: center;
        background-color: var(--wa-color-surface-raised);
    }
</style>
