<script context="module">
    export { requiresNoAuth as load } from '$lib/load'
</script>

<script>
    import { goto } from '$app/navigation'

    import { req } from '$lib/req'

    let state = 'enterEmail'

    let error = ''

    let submitting = false

    let email = ''
    let forgotPasswordToken = ''
    let ticket = ''

    let newPassword = ''
    let confirmPassword = ''

    const sendForgotPasswordToken = async () => {
        submitting = true
        error = ''

        try {
            await req('account/forgot/password', 'POST', { email })

            state = 'codeSent'
        } catch ({ message }) {
            error = message
        } finally {
            submitting = false
        }
    }

    const redeemForgotPasswordToken = async () => {
        submitting = true
        error = ''

        try {
            ticket = (await req('account/forgot/password/redeem', 'POST', {
                token: forgotPasswordToken })).ticket

            state = 'codeRedeemed'
        } catch ({ message }) {
            error = message
        } finally {
            submitting = false
        }
    }

    const changePassword = async () => {
        if (newPassword !== confirmPassword) {
            error = 'Passwords do not match.'
            return
        }

        submitting = true
        error = ''

        try {
            await req('account/password', 'POST', { ticket, newPassword })

            state = 'passwordChanged'
        } catch ({ message }) {
            error = message
        } finally {
            submitting = false
        }
    }
</script>

<style>
    .forgot-password-form {
        width: 270px;
    }
</style>

<svelte:head>
    <title>Studier &mdash; Forgot Password</title>
</svelte:head>

<div class="d-flex flex-column align-items-center justify-content-center vh-100 vw-100">
    <div class="forgot-password-form">
        <h1 class="h2 text-center">Forgot Password</h1>

        {#if state === 'enterEmail'}
            <p class="text-center mb-2">
                Enter your email to reset your password.
            </p>

            {#if error}
                <div class="alert alert-danger text-center mb-2">{error}</div>
            {/if}

            <label class="mb-3 w-100">
                Email
                <input
                    type="email"
                    class="form-control"
                    disabled={submitting}
                    bind:value={email}
                />
            </label>
            <button
                class="btn btn-primary w-100"
                disabled={submitting}
                on:click={sendForgotPasswordToken}
            >Continue</button>
        {:else if state === 'codeSent'}
            <p class="text-center mb-2">
                Enter the code that was emailed to you.
            </p>

            {#if error}
                <div class="alert alert-danger text-center mb-2">{error}</div>
            {/if}

            <label class="mb-3 w-100">
                Code
                <input
                    type="text"
                    class="form-control"
                    disabled={submitting}
                    bind:value={forgotPasswordToken}
                />
            </label>
            <button
                class="btn btn-primary w-100"
                disabled={submitting}
                on:click={redeemForgotPasswordToken}
            >Continue</button>
        {:else if state === 'codeRedeemed'}
            <p class="text-center mb-3">
                Choose a new password.
            </p>

            {#if error}
                <div class="alert alert-danger text-center mb-2">{error}</div>
            {/if}

            <label class="mb-2 w-100">
                Password
                <input type="password" class="form-control" bind:value={newPassword} />
            </label>
            <label class="mb-3 w-100">
                Confirm Password
                <input type="password" class="form-control" bind:value={confirmPassword} />
            </label>

            <button
                class="btn btn-primary w-100"
                disabled={submitting}
                on:click={changePassword}
            >Continue</button>
        {:else if state === 'passwordChanged'}
            <p class="text-center mb-3">
                Your password has been changed.
            </p>
            <button
                class="btn btn-primary w-100"
                disabled={submitting}
                on:click={() => goto('/login')}
            >Back to Log In</button>
        {/if}
    </div>
</div>
