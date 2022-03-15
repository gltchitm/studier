<script context="module">
    export { requiresNoAuth as load } from '$lib/load'
</script>

<script>
    import { goto } from '$app/navigation'

    import { req } from '$lib/req'

    let error = ''

    let submitting = false

    let email = ''
    let username = ''
    let password = ''
    let confirmPassword = ''

    const signup = async () => {
        if (password !== confirmPassword) {
            error = 'Passwords do not match.'
            return
        }

        error = ''
        submitting = true

        try {
            await req('account/signup', 'POST', {
                email,
                username,
                password
            })

            goto('/login')
        } catch ({ message }) {
            error = message
        } finally {
            submitting = false
        }
    }
</script>

<style>
    .signup-form {
        width: 270px;
    }
</style>

<svelte:head>
    <title>Studier &mdash; Sign Up</title>
</svelte:head>

<div class="d-flex flex-column align-items-center justify-content-center vh-100 vw-100">
    <div class="signup-form">
        <h1 class="text-center">Sign Up</h1>

        {#if error}
            <div class="alert alert-danger text-center text-break">
                {error}
            </div>
        {/if}

        <label class="w-100 mb-2">
            Email
            <input
                type="email"
                class="form-control"
                disabled={submitting}
                bind:value={email}
            />
        </label>
        <label class="w-100 mb-2">
            Username
            <input
                type="text"
                class="form-control"
                disabled={submitting}
                bind:value={username}
            />
        </label>
        <label class="w-100 mb-3">
            Password
            <input
                type="password"
                class="form-control"
                disabled={submitting}
                bind:value={password}
            />
        </label>
        <label class="w-100 mb-3">
            Confirm Password
            <input
                type="password"
                class="form-control"
                disabled={submitting}
                bind:value={confirmPassword}
            />
        </label>
        <button
            class="btn btn-primary w-100 mb-2"
            disabled={submitting}
            on:click={signup}
        >Sign Up</button>
        <small class="d-block text-center">
            <a href="/">Home</a> &mdash; <a href="/login">Log In</a>
        </small>
    </div>
</div>
