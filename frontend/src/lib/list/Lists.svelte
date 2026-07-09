<script lang="ts">
    import { listStore } from "./store";
    import { p } from "../../router";
    import "@awesome.me/webawesome/dist/components/icon/icon.js";

    const { lists } = listStore;
</script>

<article>
    <header>
        <h2>Lists</h2>
    </header>
    <ul>
        {#each $lists as list (list.id)}
            <li>
                <a href={p("/lists/:id", { params: { id: list.id } })}>
                    <wa-icon name="list-check"></wa-icon>
                    <span class="name">{list.name}</span>
                    <span class="count">{list.itemCount}</span>
                    <wa-icon name="chevron-right"></wa-icon>
                </a>
            </li>
        {/each}
    </ul>
</article>

<style>
    article {
        height: 100%;
        display: flex;
        flex-direction: column;
        background-color: var(--wa-color-surface-default);
        overflow-y: hidden;
    }

    header {
        height: 2rem;
        z-index: 1;
        position: sticky;
        top: 0;
        padding: var(--wa-space-m) var(--wa-space-l);
        background-color: var(--wa-color-surface-raised);
        box-shadow: 0 var(--wa-shadow-offset-x-m) var(--wa-shadow-offset-y-m)
            var(--wa-color-shadow);
    }

    h2 {
        margin: 0;
        font-size: var(--wa-font-size-xl);
        font-weight: var(--wa-font-weight-normal);
        color: var(--wa-color-text-normal);
    }

    ul {
        flex: 1;
        list-style: none;
        margin: 0;
        padding: var(--wa-space-2xs);
        overflow-y: auto;
        display: flex;
        flex-direction: column;
        gap: var(--wa-space-2xs);
    }

    a {
        display: flex;
        align-items: center;
        gap: var(--wa-space-s);
        padding: var(--wa-space-m);
        text-decoration: none;
        color: var(--wa-color-text-normal);
        background-color: var(--wa-color-surface-lowered);
        border-radius: var(--wa-border-width-l);
    }

    .name {
        flex: 1;
    }

    .count {
        color: var(--wa-color-neutral-on-quiet);
        font-size: var(--wa-font-size-s);
    }
</style>
