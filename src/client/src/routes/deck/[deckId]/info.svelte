<script>
    import { goto } from '$app/navigation'
    import { page } from '$app/stores'

    import Modal from '$lib/components/Modal.svelte'

    import { rawReq, req } from '$lib/req'

    const { deckId } = $page.params

    let error = ''

    let loadingDeck = true
    let performingAction = false

    let name = ''
    let description = ''
    let access = ''
    let flashcards = []
    let editable = false
    let author = null

    let editing = ''

    let newName = ''
    let newDescription = ''
    let newAccess = ''
    let newPassword = ''

    let confirmingDelete = false

    const token = sessionStorage.getItem('studier_deck_token_' + deckId)
    rawReq(`deck/${deckId}${token ? '?token=' + token : ''}`, 'GET').then(async response => {
        const json = await response.json().catch(() => ({ error: response.statusText }))

        if (response.ok) {
            name = json.name
            description = json.description
            access = json.access
            flashcards = json.flashcards
            editable = json.editable
            author = json.author

            newName = name
            newDescription = description
            newAccess = access

            loadingDeck = false
        } else if (response.status === 403) {
            goto(`/deck/${deckId}/writing/${json.unlockable ? 'unlock' : 'locked'}`,
                { replaceState: true })
        } else {
            error = json.error
        }
    }).catch(({ message }) => {
        error = message
    })

    const closeEditing = () => {
        editing = ''
        newName = name
        newDescription = description
        newAccess = access
        newPassword = ''
    }

    const editName = async () => {
        try {
            error = ''
            performingAction = true

            await req(`deck/${deckId}/name`, 'PUT', { name: newName })

            name = newName
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
            closeEditing()
        }
    }

    const editDescription = async () => {
        try {
            error = ''
            performingAction = true

            await req(`deck/${deckId}/description`, 'PUT', { description: newDescription })

            description = newDescription
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
            closeEditing()
        }
    }

    const editAccess = async () => {
        try {
            error = ''
            performingAction = true

            await req(`deck/${deckId}/access`, 'PUT', { access: newAccess, password: newPassword })

            access = newAccess
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
            closeEditing()
        }
    }

    const deleteDeck = async () => {
        try {
            error = ''
            performingAction = true

            await req(`deck/${deckId}`, 'DELETE')

            goto('/home')
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
            confirmingDelete = false
        }
    }
</script>

<svelte:head>
    <title>Studier &mdash; Deck Info</title>
</svelte:head>

{#if error}
    <div class="alert alert-danger">{error}</div>
{/if}

{#if loadingDeck}
    <h1>Loading...</h1>
{:else}
    <div class="mb-2">
        <small>Name</small>
        <div class="d-flex align-items-center gap-2">
            <h3 class="text-break">{name}</h3>
            {#if editable}
                <button
                    class="btn btn-primary btn-sm"
                    on:click={() => editing = 'name'}
                >Edit</button>
            {/if}
        </div>
    </div>

    <div class="mb-2">
        <small>Description</small>
        <div class="d-flex align-items-center gap-2">
            <h3 class="text-break">{description}</h3>
            {#if editable}
                <button
                    class="btn btn-primary btn-sm"
                    on:click={() => editing = 'description'}
                >Edit</button>
            {/if}
        </div>
    </div>

    {#if !editable && !loadingDeck}
        <div class="mb-2">
            <small>Author</small>
            <div class="d-flex align-items-center gap-2">
                <h3 class="text-break"><a href="/user/{author.userId}">{author.username}</a></h3>
            </div>
        </div>
    {/if}

    <div class="mb-2">
        <small>Access</small>
        <div class="d-flex align-items-center gap-2">
            <h3 class="text-break">{access[0].toUpperCase() + access.slice(1)}</h3>
            {#if editable}
                <button
                    class="btn btn-primary btn-sm"
                    on:click={() => editing = 'access'}
                >Edit</button>
            {/if}
        </div>
    </div>

    <div class="mb-2">
        <small>Flashcards</small>
        <h3>{flashcards.length}</h3>
    </div>

    {#if editable}
        <button
            class="btn btn-danger mt-2"
            on:click={() => confirmingDelete = true}
        >Delete Deck</button>
    {/if}
{/if}

{#if editing === 'name'}
    <Modal>
        <svelte:fragment slot="title">
            Edit Name
        </svelte:fragment>
        <svelte:fragment slot="body">
            <label class="w-100">
                Name
                <input
                    type="text"
                    class="form-control"
                    disabled={performingAction}
                    bind:value={newName}
                />
            </label>
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                class="btn btn-secondary"
                disabled={performingAction}
                on:click={() => closeEditing()}
            >Cancel</button>
            <button
                class="btn btn-primary"
                disabled={performingAction}
                on:click={editName}
            >Save</button>
        </svelte:fragment>
    </Modal>
{:else if editing === 'description'}
    <Modal>
        <svelte:fragment slot="title">
            Edit Description
        </svelte:fragment>
        <svelte:fragment slot="body">
            <label class="w-100">
                Description
                <input
                    type="text"
                    class="form-control"
                    disabled={performingAction}
                    bind:value={newDescription}
                />
            </label>
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                class="btn btn-secondary"
                disabled={performingAction}
                on:click={() => closeEditing()}
            >Cancel</button>
            <button
                class="btn btn-primary"
                disabled={performingAction}
                on:click={editDescription}
            >Save</button>
        </svelte:fragment>
    </Modal>
{:else if editing === 'access'}
    <Modal>
        <svelte:fragment slot="title">
            Access
        </svelte:fragment>
        <svelte:fragment slot="body">
            <label class="w-100 mb-2">
                Access
                <select
                    class="form-select"
                    disabled={performingAction}
                    bind:value={newAccess}
                >
                    <option value="everyone">Everyone</option>
                    <option value="friends">My friends</option>
                    <option value="password">People with the password</option>
                    <option value="me">Only me</option>
                </select>
            </label>
            {#if newAccess === 'password'}
                <label class="w-100">
                    New Password
                    <input type="password" class="form-control" bind:value={newPassword} />
                </label>
            {/if}
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                class="btn btn-secondary"
                disabled={performingAction}
                on:click={() => closeEditing()}
            >Cancel</button>
            <button
                class="btn btn-primary"
                disabled={performingAction}
                on:click={editAccess}
            >Save</button>
        </svelte:fragment>
    </Modal>
{/if}

{#if confirmingDelete}
    <Modal>
        <svelte:fragment slot="title">
            Delete Deck
        </svelte:fragment>
        <svelte:fragment slot="body">
            Are you sure you want to delete the deck "{name}"?
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                class="btn btn-secondary"
                disabled={performingAction}
                on:click={() => confirmingDelete = false}
            >Cancel</button>
            <button
                class="btn btn-danger"
                disabled={performingAction}
                on:click={deleteDeck}
            >Delete</button>
        </svelte:fragment>
    </Modal>
{/if}
