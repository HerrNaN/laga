<script lang="ts">
    import { browserSupportsWebAuthn, startRegistration } from "@simplewebauthn/browser";
    import { navigate, p } from "../../router";
    import { refreshSession } from "./session";

    let pending = $state(false);
    let error = $state("");

    const signUp = async () => {
        pending = true;
        error = "";
        try {
            const begin = await fetch("/api/auth/register/begin", { method: "POST" });
            if (!begin.ok) throw new Error("Could not start registration");

            const credential = await startRegistration({ optionsJSON: await begin.json() });
            const finish = await fetch("/api/auth/register/finish", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(credential),
            });
            if (!finish.ok) throw new Error(await finish.text());

            await refreshSession();
            await navigate("/user");
        } catch (cause) {
            error = cause instanceof Error ? cause.message : "Registration failed";
        } finally {
            pending = false;
        }
    };
</script>

<section>
    <h1>Create account</h1>
    <p>Create an account using a passkey on this device.</p>
    {#if browserSupportsWebAuthn()}
        <button type="button" onclick={signUp} disabled={pending}>
            {pending ? "Creating account…" : "Create account with passkey"}
        </button>
    {:else}
        <p>This browser does not support passkeys.</p>
    {/if}
    {#if error}<p role="alert">{error}</p>{/if}
    <p>Already have an account? <a href={p("/sign-in")}>Sign in</a></p>
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
