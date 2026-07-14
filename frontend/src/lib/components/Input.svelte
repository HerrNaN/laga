<script lang="ts">
    import WaInput from "@awesome.me/webawesome/dist/components/input/input.js";
    import type { WaInputProps } from "@awesome.me/webawesome/dist/custom-elements-jsx.d.ts";

    type Props = {
        value: string;
        onsave: (value: string) => void;
    } & WaInputProps;

    const { value, onsave, ...waProps }: Props = $props();

    let editing = $state(false);
    let input = $state<WaInput>();

    const commitChange = () => {
        editing = false;
        if (input && input.value !== null) {
            onsave(input.value);
        }
    };
</script>

<div>
    {#if editing}
        <wa-input
            bind:this={input}
            {...waProps}
            onchange={() => commitChange()}
            onblur={() => commitChange()}
            {value}
        ></wa-input>
    {:else}
        {value}
        <wa-button
            appearance="plain"
            aria-label="Back"
            onclick={() => {
                editing = true;
                setTimeout(() => {
                    if (input) input.focus();
                });
            }}
        >
            <wa-icon name="pen"></wa-icon>
        </wa-button>
    {/if}
</div>

<style>
    div {
        display: flex;
        align-items: center;

        gap: var(--wa-space-s);
    }
</style>
