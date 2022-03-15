<script context="module">
    export { requiresAuthNotVerified as load } from '$lib/load'
</script>

<script>
    import { goto } from '$app/navigation'

    import { req } from '$lib/req'

    export let user

    let error = ''

    let showResent = false
    let submitting = false

    let verificationCode = ''

    const verify = async () => {
        error = ''
        showResent = false
        submitting = true

        try {
            await req('account/verify', 'POST', { verificationCode })

            goto('/home')
        } catch ({ message }) {
            error = message
        } finally {
            submitting = false
        }
    }

    const resend = async () => {
        error = ''
        showResent = false
        submitting = true
        try {

            await req('account/verify/resend', 'POST')

            showResent = true
        } catch ({ message }) {
            error = message
        } finally {
            submitting = false
        }
    }

    const logout = async () => {
        error = ''
        showResent = false
        submitting = true

        try {
            await req('auth/logout', 'DELETE')

            goto('/login')
        } catch ({ message }) {
            error = message
        } finally {
            submitting = false
        }
    }
</script>

<style>
    .verify-form {
        width: 270px;
    }
</style>

<svelte:head>
    <title>Studier &mdash; Verify Account</title>
</svelte:head>

<div class="d-flex flex-column align-items-center justify-content-center vh-100 vw-100">
    <div class="verify-form">
        <h1 class="text-center">Verify Account</h1>
        <p class="text-center">
            Enter the code emailed to you to verify your account.
        </p>

        <div class="alert alert-primary text-center mb-3 text-break">
            Logged in as <b class="break-all">{user.username}</b>.
        </div>

        {#if error}
            <div class="alert alert-danger text-center mb-3">{error}</div>
        {/if}

        {#if showResent}
            <div class="alert alert-warning text-center mb-3">
                Verification code resent.
            </div>
        {/if}

        <label class="mb-3 w-100">
            Verification Code
            <input
                type="text"
                class="form-control"
                disabled={submitting}
                bind:value={verificationCode}
            />
        </label>

        <button
            class="btn btn-primary w-100 mb-3"
            disabled={submitting}
            on:click={verify}
        >Verify</button>

        <div class="d-flex gap-3">
            <button
                class="btn btn-warning w-50"
                disabled={submitting}
                on:click={resend}
            >
                Resend
            </button>
            <button
                class="btn btn-danger w-50"
                disabled={submitting}
                on:click={logout}
            >Log Out</button>
        </div>
    </div>
</div>
