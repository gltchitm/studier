<script context="module">
    export { requiresNoAuth as load } from '$lib/load'
</script>

<script>
    import { goto } from '$app/navigation'

    import { req } from '$lib/req'

    let error = ''

    let submitting = false

    let username = ''
    let password = ''

    const login = async () => {
        error = ''
        submitting = true

        try {
            const { verified } = await req('auth/login', 'POST', {
                username,
                password
            })

            if (verified) {
                goto('/home')
            } else {
                goto('/verify')
            }
        } catch ({ message }) {
            error = message
        } finally {
            submitting = false
        }
    }
</script>

<style>
    .login-form {
        width: 270px;
    }
</style>

<svelte:head>
    <title>Studier &mdash; Log In</title>
</svelte:head>

<div class="d-flex flex-column align-items-center justify-content-center vh-100 vw-100">
    <div class="login-form">
        <h1 class="text-center">Log In</h1>

        {#if error}
            <div class="alert alert-danger text-center text-break">
                {error}
            </div>
        {/if}

        <div class="mb-2">
            <label class="w-100">
                Username
                <input
                    type="text"
                    class="form-control"
                    disabled={submitting}
                    bind:value={username}
                />
            </label>
            <a href="/forgot/username">Forgot your username?</a>
        </div>
        <div class="mb-3">
            <label class="w-100">
                Password
                <input
                    type="password"
                    class="form-control"
                    disabled={submitting}
                    bind:value={password}
                />
            </label>
            <a href="/forgot/password">Forgot your password?</a>
        </div>
        <button
            class="btn btn-primary w-100 mb-2"
            disabled={submitting}
            on:click={login}
        >Log In</button>
        <small class="d-block text-center">
            <a href="/">Home</a> &mdash; <a href="/signup">Sign Up</a>
        </small>
    </div>
</div>
