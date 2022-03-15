<script context="module">
    export { requiresAuthVerified as load } from '$lib/load'
</script>

<script>
    import { goto } from '$app/navigation'

    import Modal from '$lib/components/Modal.svelte'
    import Navbar from '$lib/components/Navbar.svelte'

    import { req } from '$lib/req'

    export let user

    let error = ''

    let performingAction = false

    let changing = ''

    let newEmail = ''
    let confirmNewEmail = ''
    let password = ''

    let oldPassword = ''
    let newPassword = ''
    let confirmNewPassword = ''

    let deleteAccountConfirmed = false

    const cancelChangeEmail = () => {
        newEmail = ''
        confirmNewEmail = ''
        password = ''
        changing = ''
    }
    const changeEmail = async () => {
        try {
            performingAction = true

            await req('account/email', 'POST', { newEmail, password })

            goto('/login')
        } catch ({ message }) {
            error = message
            performingAction = false
            cancelChangeEmail()
        }
    }

    const cancelChangePassword = () => {
        newPassword = ''
        confirmNewPassword = ''
        changing = ''
    }
    const changePassword = async () => {
        try {
            performingAction = true

            await req('account/password', 'POST', { oldPassword, newPassword })

            goto('/login')
        } catch ({ message }) {
            error = message
            performingAction = false
            cancelChangePassword()
        }
    }

    const cancelDeleteAccount = () => {
        password = ''
        changing = ''
    }
    const deleteAccount = async () => {
        try {
            performingAction = true

            await req('account', 'DELETE', { password })

            for (const key of Object.keys(sessionStorage)) {
                if (key.startsWith('studier_deck_token_')) {
                    sessionStorage.removeItem(key)
                }
            }

            goto('/login')
        } catch ({ message }) {
            error = message
            performingAction = false
            cancelDeleteAccount()
        }
    }
</script>

<div class="p-3">
    <Navbar {user} />

    {#if error}
        <div class="alert alert-danger">
            {error}
        </div>
    {/if}

    <h1 class="mb-3">Settings</h1>
    <div class="fs-5 mb-3">
        Username: <b>{user.username}</b>
    </div>
    <div class="fs-5 mb-3">
        Email: <b>{user.email}</b>
        <button
            class="btn btn-primary btn-sm"
            on:click={() => changing = 'email'}
        >Change</button>
    </div>
    <div class="fs-5 mb-3">
        Password:
        <button
            class="btn btn-primary btn-sm"
            on:click={() => changing = 'password'}
        >Change</button>
    </div>

    <button
        class="btn btn-danger mb-3"
        on:click={() => changing = 'delete'}
    >Delete Account</button>

    <p class="mb-2">
        Studier is licensed under the <a href="https://github.com/gltchitm/studier/blob/master/LICENSE">
            MIT license</a>.
        View its source code <a href="https://github.com/gltchitm/studier">here</a>.
    </p>
    <p>
        Studier's logo is from <a href="https://github.com/twitter/twemoji">Twemoji</a>.
        It is used under the <a href="https://creativecommons.org/licenses/by/4.0/">CC BY 4.0 license.</a>
    </p>
</div>

{#if changing === 'email'}
    <Modal>
        <svelte:fragment slot="title">
            Change Email
        </svelte:fragment>
        <svelte:fragment slot="body">
            <div class="alert alert-danger">
                Once you change your email, you will be logged out and required to
                verify your new email. <b>Ensure you type the correct email</b>
                as you will be permanently locked out of your account if you make a mistake.
            </div>
            <label class="w-100 mb-2">
                New Email
                <input
                    type="email"
                    class="form-control"
                    disabled={performingAction}
                    bind:value={newEmail}
                />
            </label>
            <label class="w-100 mb-2">
                Confirm New Email
                <input
                    type="email"
                    class="form-control"
                    disabled={performingAction}
                    bind:value={confirmNewEmail}
                />
            </label>
            <label class="w-100">
                Password
                <input
                    type="password"
                    class="form-control"
                    disabled={performingAction}
                    bind:value={password}
                />
            </label>
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                class="btn btn-secondary"
                disabled={performingAction}
                on:click={cancelChangeEmail}
            >Cancel</button>
            <button
                class="btn btn-danger"
                disabled={performingAction || !newEmail.length || newEmail !== confirmNewEmail}
                on:click={changeEmail}
            >Change</button>
        </svelte:fragment>
    </Modal>
{:else if changing === 'password'}
    <Modal>
        <svelte:fragment slot="title">
            Change Password
        </svelte:fragment>
        <svelte:fragment slot="body">
            <div class="alert alert-danger">
                Once you change your password, you will be logged out and to required
                log back in. <b>Ensure you remember your new password</b>
                as you will be permanently locked out of your account if you forget it.
            </div>
            <label class="w-100 mb-2">
                Old Password
                <input
                    type="password"
                    class="form-control"
                    disabled={performingAction}
                    bind:value={oldPassword}
                />
            </label>
            <label class="w-100 mb-2">
                New Password
                <input
                    type="password"
                    class="form-control"
                    disabled={performingAction}
                    bind:value={newPassword}
                />
            </label>
            <label class="w-100">
                Confirm New Password
                <input
                    type="password"
                    class="form-control"
                    disabled={performingAction}
                    bind:value={confirmNewPassword}
                />
            </label>
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                class="btn btn-secondary"
                disabled={performingAction}
                on:click={cancelChangePassword}
            >Cancel</button>
            <button
                class="btn btn-danger"
                disabled={performingAction || !newPassword.length || newPassword !== confirmNewPassword}
                on:click={changePassword}
            >Change</button>
        </svelte:fragment>
    </Modal>
{:else if changing === 'delete'}
    <Modal>
        <svelte:fragment slot="title">
            Delete Account
        </svelte:fragment>
        <svelte:fragment slot="body">
            <div class="alert alert-danger">
                <b>WARNING</b>: This action is <u>irreversible</u>! There is no way to
                recover your account once it is deleted. Proceed with caution!
            </div>
            <label class="w-100 mb-2">
                Password
                <input
                    type="password"
                    class="form-control"
                    disabled={performingAction}
                    bind:value={password}
                />
            </label>
            <div class="form-check">
                <label class="form-check-label">
                    <input
                        type="checkbox"
                        class="form-check-input"
                        bind:value={deleteAccountConfirmed}
                    />
                    I understand that deleting my account is an irreversible action.
                </label>
            </div>
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button class="btn btn-secondary" on:click={cancelDeleteAccount}>Cancel</button>
            {#if deleteAccountConfirmed}
                <button class="btn btn-danger" on:click={deleteAccount}>Delete Account</button>
            {/if}
        </svelte:fragment>
    </Modal>
{/if}
