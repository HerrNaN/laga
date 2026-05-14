<script lang="ts">
    import { items } from "../stores";
    import "@awesome.me/webawesome/dist/components/input/input.js";
    import "@awesome.me/webawesome/dist/components/button/button.js";

    let inputValue = $state("");

    function handleAdd() {
        const text = inputValue.trim();
        if (text) {
            items.addItem(text);
            inputValue = "";
        }
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Enter") {
            handleAdd();
        }
    }
</script>

<section>
    <wa-input
        placeholder="e.g. 2 kg potatoes"
        value={inputValue}
        oninput={(e: Event) => {
            inputValue = (e.target as HTMLInputElement).value;
        }}
        onkeydown={handleKeydown}
    ></wa-input>
    <wa-button variant="brand" onclick={handleAdd}> Add </wa-button>
</section>

<style>
    section {
        display: flex;
        gap: var(--wa-space-m);
        padding: var(--wa-space-m);
    }

    wa-input {
        flex: 1;
    }

    wa-button {
        flex-shrink: 0;
    }
</style>
