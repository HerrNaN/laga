<script lang="ts">
    import "@awesome.me/webawesome/dist/components/icon/icon.js";
    import "@awesome.me/webawesome/dist/components/button/button.js";
</script>

<nav>
    <ul>
        <li>
            <wa-button disabled href="/#" appearance="plain">
                <wa-icon slot="start" name="kitchen-set"></wa-icon>
                Recepies (Soon)
            </wa-button>
        </li>
        <li class="active">
            <wa-button href="/#" appearance="plain">
                <wa-icon slot="start" name="list-check"></wa-icon>
                Lists
            </wa-button>
        </li>
    </ul>
</nav>

<style>
    nav {
        height: 4rem;
        width: 100%;
        z-index: 1;

        box-shadow: 0 var(--wa-shadow-offset-x-m) var(--wa-shadow-offset-y-m)
            var(--wa-color-shadow);
    }

    ul {
        height: 100%;
        width: 100%;
        margin: 0;
        padding: 0;
        display: flex;
        align-items: center;
    }

    li {
        height: 100%;
        list-style: none;
        flex: 1;
        display: flex;
        justify-content: center;
        align-items: center;

        &.active {
            wa-icon,
            wa-button::part(label) {
                color: var(--wa-color-brand-on-quiet);
            }
            wa-button::part(base) {
                background-color: var(--wa-color-brand-fill-quiet);
            }
        }
    }
</style>
