<script lang="ts">
    import { browserSupportsWebAuthn, startAuthentication } from "@simplewebauthn/browser";
    import { navigate, p } from "../../router";
    import { refreshSession } from "./session";

    let pending = $state(false);
    let error = $state("");

    const signIn = async () => {
        pending = true;
        error = "";
        try {
            const begin = await fetch("/api/auth/login/begin", { method: "POST" });
            if (!begin.ok) throw new Error("Could not start sign-in");

            const credential = await startAuthentication({ optionsJSON: await begin.json() });
            const finish = await fetch("/api/auth/login/finish", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(credential),
            });
            if (!finish.ok) throw new Error(await finish.text());

            await refreshSession();
            await navigate("/user");
        } catch (cause) {
            error = cause instanceof Error ? cause.message : "Sign-in failed";
        } finally {
            pending = false;
        }
    };
</script>

<section>
    <h1>Sign in</h1>
    <p>Sign in to your account using your passkey.</p>
    {#if browserSupportsWebAuthn()}
        <button type="button" onclick={signIn} disabled={pending}>
            {pending ? "Signing in…" : "Sign in with passkey"}
        </button>
    {:else}
        <p>This browser does not support passkeys.</p>
    {/if}
    {#if error}<p role="alert">{error}</p>{/if}
    <p>New to Laga? <a href={p("/sign-up")}>Create an account</a></p>
    <p><a href={p("/lists")}>Back to lists</a></p>
</section>

<style>
    section {
        padding: var(--wa-space-l);
    }

    button {
        padding: var(--wa-space-m);
        min-height: 48px;
    }
</style>
