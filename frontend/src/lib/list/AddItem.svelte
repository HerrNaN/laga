<script lang="ts">
    import WaInput from "@awesome.me/webawesome/dist/components/input/input.js";
    import "@awesome.me/webawesome/dist/components/drawer/drawer.js";
    import { items } from "./store";

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

    const id = $props.id();
</script>

<div class="input-wrapper">
    <wa-drawer
        {id}
        placement="bottom"
        light-dismiss
        without-header
        onwa-show={() => setTimeout(() => input?.focus())}
    >
        <wa-input
            bind:this={input}
            enterkeyhint="go"
            placeholder="t.ex. 2 kg potatis"
            bind:this={input}
            onkeydown={handleKeydown}
        ></wa-input>
    </wa-drawer>
    <wa-button size="l" variant="brand" data-drawer="open {id}">
        <wa-icon name="plus"></wa-icon>
    </wa-button>
</div>

<style>
    .input-wrapper {
        height: 100%;
        position: sticky;
        bottom: 0;
        display: flex;
        justify-content: flex-end;
        align-items: flex-end;
    }

    wa-button {
        margin: 1rem;
    }

    wa-drawer {
        --size: min-content;
        --show-duration: 0;

        :global(&[open] + wa-button) {
            display: none;
        }
    }

    wa-input {
        width: 100%;
    }
</style>
