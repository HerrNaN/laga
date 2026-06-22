<script lang="ts">
    import "@awesome.me/webawesome/dist/components/tab-group/tab-group.js";
    import "@awesome.me/webawesome/dist/components/tab-panel/tab-panel.js";
    import "@awesome.me/webawesome/dist/components/tab/tab.js";
    import List from "./lib/list/List.svelte";
</script>

<main>
    <wa-tab-group placement="bottom">
        <wa-tab panel="recipes" disabled>
            <wa-icon name="kitchen-set"></wa-icon>
            Recipes (soon!)
        </wa-tab>
        <wa-tab panel="lists" active>
            <wa-icon name="list-check"></wa-icon>
            Lists
        </wa-tab>
        <wa-tab-panel name="recipes"></wa-tab-panel>
        <wa-tab-panel name="lists">
            <List />
        </wa-tab-panel>
    </wa-tab-group>
</main>

<style>
    :global(body) {
        --add-item-height: calc(var(--wa-space-m) * 2 + 3rem);
        margin: 0;
    }

    main {
        height: 100dvh;
        max-height: 100svh;
        display: flex;
        flex-direction: column;
    }

    wa-tab {
        width: 100%;
        height: 100%;

        &::part(base) {
            width: 100%;
            display: flex;
            justify-content: center;
        }
    }

    wa-tab-panel::part(base) {
        padding-block: 0;
    }

    wa-tab-group,
    wa-tab-group::part(base),
    wa-tab-group::part(body),
    wa-tab-panel,
    wa-tab-panel::part(base) {
        height: 100%;
        overflow-y: hidden;
    }
</style>
