<script lang="ts">
    import { p } from "../../router";
    import { session } from "./session";

    const logout = () =>
      fetch("/api/auth/logout", { method: "POST" }).then(() => location.reload())

</script>

<section>
    <h1>Profile</h1>
    {#if $session}
        <p>You are signed in as: <span class="username">{$session.name}</span>.</p>
        <wa-button appearance="outlined" onclick={logout}>
            <wa-icon slot="start" name="arrow-right-from-bracket"></wa-icon>
            Log out
        </wa-button>
    {:else}
        <p>You are not signed in.</p>
        <p><a href={p("/sign-in")}>Sign in</a></p>
        <p><a href={p("/sign-up")}>Create an account</a></p>
    {/if}
</section>

<style>
    section {
        padding: var(--wa-space-l);
    }
    .username {
        font-weight: var(--wa-font-weight-bold);
    }
</style>
