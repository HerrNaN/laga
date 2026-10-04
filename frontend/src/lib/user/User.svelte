<script lang="ts">
    import { p } from "../../router";
    import { session } from "./session";
</script>

<section>
    <h1>Account</h1>
    {#if $session}
        <p>Signed in as {$session.name}.</p>
        <p>Shopping lists are still stored on this device for now.</p>
    {:else}
        <p>You are not signed in.</p>
        <p><a href={p("/sign-in")}>Sign in</a></p>
        <p><a href={p("/sign-up")}>Create an account</a></p>
    {/if}
    <p><a href={p("/lists")}>Back to lists</a></p>
</section>

<style>
    section {
        padding: var(--wa-space-l);
    }
</style>
