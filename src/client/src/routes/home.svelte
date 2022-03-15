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

    let loading = true
    let performingAction = false

    let allDecks
    let pinnedDecks

    let confirmDeleteDeck = null

    let search = ''

    Promise.all([
        req('deck/list/all', 'GET').then(({ decks }) => allDecks = decks),
        req('deck/list/pinned', 'GET').then(({ decks }) => pinnedDecks = decks)
    ]).then(() => loading = false).catch(({ message }) => error = message)

    const pinDeck = async deck => {
        try {
            error = ''
            performingAction = true

            await req(`deck/${deck.deckId}/pin`, 'POST')

            deck.pinned = true
            pinnedDecks = [deck, ...pinnedDecks]
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }

    const unpinDeck = async deck => {
        try {
            error = ''
            performingAction = true

            await req(`deck/${deck.deckId}/pin`, 'DELETE')

            deck.pinned = false
            pinnedDecks = pinnedDecks.filter(({ deckId }) => deckId !== deck.deckId)
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }

    const deleteDeck = async () => {
        try {
            error = ''
            performingAction = true

            await req(`deck/${confirmDeleteDeck.deckId}`, 'DELETE')

            allDecks = allDecks.filter(deck => confirmDeleteDeck.deckId !== deck.deckId)

            if (pinnedDecks.includes(confirmDeleteDeck)) {
                pinnedDecks = pinnedDecks.filter(deck => confirmDeleteDeck.deckId !== deck.deckId)
            }

            confirmDeleteDeck = null
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }
</script>
<style>
    .search {
        width: 225px;
    }
</style>

<svelte:head>
    <title>Studier &mdash; Home</title>
</svelte:head>

<div class="p-3">
    <Navbar active="home" {user} />

    {#if error}
        <div class="alert alert-danger w-100">{error}</div>
    {/if}

    {#if loading}
        <h1>Loading...</h1>
    {:else}
        <div class="mb-3">
            <h4>Pinned Decks</h4>
            <div class="d-flex gap-3 overflow-auto">
                {#if !pinnedDecks.length}
                    <div class="alert alert-info w-100">
                        You don't have any pinned decks.
                    </div>
                {/if}
                {#each pinnedDecks as deck}
                    <div class="card flex-shrink-0" style="width: 325px; height: 180px;">
                        <div class="card-body">
                            <h3 class="card-title h5 text-truncate m-0 mb-1">{deck.name}</h3>
                            <p class="text-truncate p-0 m-0 mb-2">{deck.description}</p>
                            <h5 class="card-subtitle h6 text-muted mb-1">
                                {deck.flashcards} flashcards
                            </h5>
                        </div>
                        <div class="card-footer">
                            <button
                                class="btn btn-primary w-100"
                                on:click={() => goto(`/deck/${deck.deckId}/study`)}
                            >Study</button>
                        </div>
                    </div>
                {/each}
            </div>
        </div>

        <div>
            <h4 class="d-flex justify-content-between">
                My Decks
                <input
                    type="text"
                    class="form-control search"
                    placeholder="Search&hellip;"
                    bind:value={search}
                />
            </h4>
            <table class="table table-striped table-bordered table-sm align-middle overflow-auto mb-0">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Flashcards</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {#each allDecks.filter(({ name }) => name.includes(search))
                        .sort((a, b) => a.name.localeCompare(b.name)) as deck}
                        <tr>
                            <td>{deck.name}</td>
                            <td>{deck.flashcards}</td>
                            <td>
                                <button
                                    class="btn btn-primary btn-sm"
                                    on:click={() => goto(`/deck/${deck.deckId}/study`)}
                                >Study</button>
                                <button
                                    class="btn btn-danger btn-sm"
                                    on:click={() => confirmDeleteDeck = deck}
                                >Delete</button>
                                {#if deck.pinned}
                                    <button
                                        class="btn btn-warning btn-sm"
                                        disabled={performingAction}
                                        on:click={() => unpinDeck(deck)}
                                    >Unpin</button>
                                {:else}
                                    <button
                                        class="btn btn-success btn-sm"
                                        disabled={performingAction}
                                        on:click={() => pinDeck(deck)}
                                    >Pin</button>
                                {/if}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>

            {#if !allDecks.length}
                <div class="text-center border border-top-0 p-2">
                    You have no decks. <a href="/deck/new">Create one?</a>
                </div>
            {/if}
        </div>
    {/if}
</div>


{#if confirmDeleteDeck}
    <Modal>
        <svelte:fragment slot="title">
            Delete Deck
        </svelte:fragment>
        <svelte:fragment slot="body">
            Are you sure you want to delete the deck "{confirmDeleteDeck.name}"?
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                class="btn btn-secondary"
                disabled={performingAction}
                on:click={() => confirmDeleteDeck = null}
            >Cancel</button>
            <button
                class="btn btn-danger"
                disabled={performingAction}
                on:click={deleteDeck}
             >Delete</button>
        </svelte:fragment>
    </Modal>
{/if}
