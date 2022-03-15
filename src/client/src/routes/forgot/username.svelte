<script context="module">
    export { requiresNoAuth as load } from '$lib/load'
</script>

<script>
    import { goto } from '$app/navigation'

    import { req } from '$lib/req'

    let error = ''

    let submitting = false
    let success = false

    let email = ''

    const forgotUsername = async () => {
        success = false
        submitting = true

        try {
            await req('account/forgot/username', 'POST', { email })

            success = true
        } catch ({ message }) {
            error = message
        } finally {
            submitting = false
        }
    }
</script>

<style>
    .forgot-username-form {
        width: 270px;
    }
</style>

<svelte:head>
    <title>Studier &mdash; Forgot Username</title>
</svelte:head>

<div class="d-flex flex-column align-items-center justify-content-center vh-100 vw-100">
    <div class="forgot-username-form">
        <h1 class="h2 text-center">Forgot Username</h1>

        {#if success}
            <p class="text-center mb-3">
                Your account's username has been emailed to you.
            </p>
            <button
                class="btn btn-primary w-100"
                on:click={() => goto('/login')}
            >Back to Log In</button>
        {:else}
            <p class="text-center mb-2">
                Enter your email to have your username sent to you.
            </p>

            {#if error}
                <div class="alert alert-danger text-center mb-2">{error}</div>
            {/if}

            <label class="mb-3 w-100">
                Email
                <input type="email" class="form-control" bind:value={email} />
            </label>

            <button
                class="btn btn-primary w-100"
                disabled={submitting}
                on:click={forgotUsername}
            >Submit</button>
        {/if}
    </div>
</div>
