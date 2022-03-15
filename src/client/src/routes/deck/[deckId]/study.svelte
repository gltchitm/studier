<script>
    import { fly } from 'svelte/transition'

    import Modal from '$lib/components/Modal.svelte'

    import { rawReq, req } from '$lib/req'

    import { goto } from '$app/navigation'
    import { page } from '$app/stores'

    const { deckId } = $page.params

    let error = ''

    let loadingDeck = true
    let performingAction = false

    let name = ''
    let description = ''
    let flashcards = []
    let editable = false
    let author = null
    let pinned = false

    let flashcardIndex = 0
    $: jumpToFlashcardNumber = flashcardIndex + 1

    let flipped = false

    let animatingFlashcardOut = false
    let inX = 0
    let outX = 0

    let confirmDeleteFlashcardId = ''

    let editingFlashcard = null

    let showNewFlashcardModal = false
    let newFlashcardTerm = ''
    let newFlashcardDefinition = ''

    let showJumpToFlashcardModal = false

    const token = sessionStorage.getItem('studier_deck_token_' + deckId)
    rawReq(`deck/${deckId}${token ? '?token=' + token : ''}`, 'GET').then(async response => {
        const json = await response.json().catch(() => ({ error: response.statusText }))

        if (response.ok) {
            name = json.name
            description = json.description
            flashcards = json.flashcards
            editable = json.editable
            author = json.author
            pinned = json.pinned

            loadingDeck = false
        } else if (response.status === 403) {
            goto(`/deck/${deckId}/study/${json.unlockable ? 'unlock' : 'locked'}`,
                { replaceState: true })
        } else {
            error = json.error
        }
    }).catch(({ message }) => {
        error = message
    })

    const flip = () => {
        if (loadingDeck) {
            return
        }

        flipped = !flipped
    }

    const back = () => {
        inX = 50
        outX = -50
        flashcardIndex--
        animatingFlashcardOut = true
    }

    const forward = () => {
        inX = -50
        outX = 50
        flashcardIndex++
        animatingFlashcardOut = true
    }

    const pinDeck = async () => {
        try {
            error = ''
            performingAction = true

            await req(`deck/${deckId}/pin`, 'POST')

            pinned = true
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }

    const unpinDeck = async () => {
        try {
            error = ''
            performingAction = true

            await req(`deck/${deckId}/pin`, 'DELETE')

            pinned = false
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }

    const deleteFlashcard = async () => {
        try {
            error = ''
            performingAction = true

            await req('flashcard/' + confirmDeleteFlashcardId, 'DELETE')

            if (flashcardIndex && flashcards[flashcardIndex].flashcardId === confirmDeleteFlashcardId) {
                flashcardIndex -= 1
            }

            flashcards = flashcards.filter(flashcard => flashcard.flashcardId !== confirmDeleteFlashcardId)
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
            confirmDeleteFlashcardId = ''
        }
    }

    const editFlashcard = async () => {
        try {
            error = ''
            performingAction = true

            await req('flashcard/' + editingFlashcard.flashcardId, 'PUT', {
                term: editingFlashcard.term,
                definition: editingFlashcard.definition
            })

            flashcards = flashcards.map(flashcard => {
                if (flashcard.flashcardId === editingFlashcard.flashcardId) {
                    flashcard.term = editingFlashcard.term
                    flashcard.definition = editingFlashcard.definition
                }

                return flashcard
            })
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
            editingFlashcard = null
        }
    }

    const closeNewFlashcardModal = () => {
        showNewFlashcardModal = false
        newFlashcardTerm = ''
        newFlashcardDefinition = ''
    }

    const newFlashcard = async () => {
        try {
            error = ''
            performingAction = true

            const { flashcardId } = await req(`deck/${deckId}/flashcard`, 'POST', {
                term: newFlashcardTerm,
                definition: newFlashcardDefinition
            })

            flashcards = [...flashcards, {
                flashcardId,
                term: newFlashcardTerm,
                definition: newFlashcardDefinition
            }]
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
            closeNewFlashcardModal()
        }
    }

    const moveFlashcardUp = async flashcardId => {
        try {
            error = ''
            performingAction = true

            await req(`flashcard/${flashcardId}/move`, 'POST', {
                direction: 'up'
            })

            const previousIndex = flashcards.findIndex(flashcard => flashcard.flashcardId === flashcardId)

            flashcards = flashcards.map((flashcard, index) => {
                if (index === previousIndex - 1) {
                    return flashcards[previousIndex]
                } else if (index === previousIndex) {
                    return flashcards[previousIndex - 1]
                }

                return flashcard
            })
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }

    const moveFlashcardDown = async flashcardId => {
        try {
            error = ''
            performingAction = true

            await req(`flashcard/${flashcardId}/move`, 'POST', {
                direction: 'down'
            })

            const previousIndex = flashcards.findIndex(flashcard => flashcard.flashcardId === flashcardId)

            flashcards = flashcards.map((flashcard, index) => {
                if (index === previousIndex + 1) {
                    return flashcards[previousIndex]
                } else if (index === previousIndex) {
                    return flashcards[previousIndex + 1]
                }

                return flashcard
            })
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }

    const cancelJumpToFlashcard = () => {
        jumpToFlashcardNumber = flashcardIndex + 1
        showJumpToFlashcardModal = false
    }
    const jumpToFlashcard = () => {
        const newIndex = jumpToFlashcardNumber - 1

        if (newIndex > flashcardIndex) {
            forward()
        } else if (newIndex < flashcardIndex) {
            back()
        }

        flashcardIndex = newIndex
        showJumpToFlashcardModal = false
    }
</script>

<style>
    .flashcard-container {
        width: 500px;
        height: 300px;
    }
    .flashcard {
        cursor: pointer;
        transform-style: preserve-3d;
        transition: transform 0.5s;
    }
    .flashcard.flipped {
        transform: rotateX(180deg);
    }
    .flashcard-front, .flashcard-back {
        position: absolute;
        backface-visibility: hidden;
        transform-style: preserve-3d;
    }
    .flashcard-back {
        transform: rotateX(180deg);
    }
    .current-flashcard {
        position: absolute;
        left: 50%;
        transform: translateX(-50%);
        cursor: pointer;
    }
    .flashcard-info {
        position: relative;
    }
    .jump-to-flashcard-number {
        width: 70px;
    }
</style>

<svelte:head>
    <title>Studier &mdash; Study</title>
</svelte:head>

<div class="d-flex align-items-center justify-content-center flex-column w-100">
    {#if error}
        <div class="alert alert-danger w-100">
            {error}
        </div>
    {/if}

    <h2 class="mb-{editable || loadingDeck ? 2 : 0} text-center text-break">
        {#if name}
            {name}
        {:else}
            Loading...
        {/if}
    </h2>

    {#if !editable && !loadingDeck}
        <small class="text-muted mb-2">
            Author: <a href="/user/{author.userId}">{author.username}</a>
        </small>
    {/if}

    {#if description}
        <p class="text-break text-center mb-2">{description}</p>
    {/if}

    {#if editable}
        {#if pinned}
            <button
                class="btn btn-warning btn-sm mb-3"
                disabled={performingAction}
                on:click={() => unpinDeck()}
            >
                Unpin
            </button>
        {:else}
            <button
                class="btn btn-success btn-sm mb-3"
                disabled={performingAction}
                on:click={() => pinDeck()}
            >
                Pin
            </button>
        {/if}
    {/if}

    <div class="flashcard-container d-flex flex-column gap-2 mb-5">
        {#if !animatingFlashcardOut}
            <div
                class="flashcard-animation-container h-100 w-100"
                in:fly|local={{ x: inX, duration: 125 }}
                out:fly|local={{ x: outX, duration: 125 }}
                on:outroend={() => animatingFlashcardOut = false}
            >
                <div
                    class="flashcard w-100 h-100 placeholder-glow"
                    class:flipped={flipped}
                    on:click={flip}
                >
                    <div
                        class="flashcard-front h-100 w-100 card d-flex justify-content-center p-3"
                        class:placeholder={loadingDeck}
                    >
                        {#if !loadingDeck}
                            <h3 class="text-break text-center">
                                {flashcards[flashcardIndex].term}
                            </h3>
                        {/if}
                    </div>
                    <div
                        class="flashcard-back h-100 w-100 card d-flex justify-content-center p-3"
                        class:placeholder={loadingDeck}
                    >
                        {#if !loadingDeck}
                            <h2 class="text-break text-center">
                                {flashcards[flashcardIndex].definition}
                            </h2>
                        {/if}
                    </div>
                </div>
            </div>
        {/if}

        <div class="d-flex justify-content-between w-100 flashcard-info">
            <div class="fw-light">
                {#if !loadingDeck}
                    {#if !flipped}
                        Term
                    {:else}
                        Definition
                    {/if}
                {/if}
            </div>
            <div class="current-flashcard fw-light" on:click={() => showJumpToFlashcardModal = true}>
                {#if !loadingDeck}
                    {flashcardIndex + 1}/{flashcards.length}
                {/if}
            </div>
            <div>
                <button
                    class="btn btn-primary btn-sm"
                    disabled={loadingDeck || flashcardIndex === 0}
                    on:click={back}
                >Back</button>
                <button
                    class="btn btn-primary btn-sm"
                    disabled={loadingDeck || flashcardIndex === flashcards.length - 1}
                    on:click={forward}
                >Forward</button>
            </div>
        </div>
    </div>

    <div class="d-flex align-items-center justify-content-between w-100">
        <h3>Flashcards</h3>
        {#if editable}
            <button
                class="btn btn-primary btn-sm"
                disabled={performingAction}
                on:click={() => showNewFlashcardModal = true}
            >New Flashcard</button>
        {/if}
    </div>

    {#if loadingDeck}
        <div class="alert alert-dark w-100">Loading...</div>
    {:else}
        <div class="table-responsive w-100">
            <table class="table table-striped table-bordered table-sm align-middle">
                <thead>
                    <tr>
                        <th>Term</th>
                        <th>Definition</th>
                        {#if editable}
                            <th>Actions</th>
                        {/if}
                    </tr>
                </thead>
                <tbody>
                    {#each flashcards as flashcard, index}
                        <tr>
                            <td class="text-break">{flashcard.term}</td>
                            <td class="text-break">{flashcard.definition}</td>
                            {#if editable}
                                <td>
                                    <div class="d-flex flex-nowrap gap-1">
                                        <button
                                            class="btn btn-primary btn-sm"
                                            disabled={performingAction}
                                            on:click={() => editingFlashcard = { ...flashcard }}
                                        >Edit</button>
                                        <button
                                            class="btn btn-danger btn-sm"
                                            disabled={performingAction}
                                            on:click={() => confirmDeleteFlashcardId = flashcard.flashcardId}
                                        >Delete</button>
                                        <button
                                            class="btn btn-info btn-sm"
                                            disabled={performingAction || !index}
                                            on:click={() => moveFlashcardUp(flashcard.flashcardId)}
                                        >&uarr;</button>
                                        <button
                                            class="btn btn-info btn-sm"
                                            disabled={performingAction || index === flashcards.length - 1}
                                            on:click={() => moveFlashcardDown(flashcard.flashcardId)}
                                        >&darr;</button>
                                    </div>
                                </td>
                            {/if}
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>

{#if confirmDeleteFlashcardId}
    <Modal>
        <svelte:fragment slot="title">
            Delete Flashcard
        </svelte:fragment>
        <svelte:fragment slot="body">
            Are you sure you want to delete this flashcard?
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                class="btn btn-secondary"
                disabled={performingAction}
                on:click={() => confirmDeleteFlashcardId = ''}
            >No</button>
            <button
                class="btn btn-primary"
                disabled={performingAction}
                on:click={deleteFlashcard}
            >Yes</button>
        </svelte:fragment>
    </Modal>
{/if}

{#if editingFlashcard}
    <Modal>
        <svelte:fragment slot="title">
            Edit Flashcard
        </svelte:fragment>
        <svelte:fragment slot="body">
            <label class="w-100 mb-2">
                Term
                <input type="text" class="form-control" bind:value={editingFlashcard.term} />
            </label>
            <label class="w-100">
                Definition
                <input type="text" class="form-control" bind:value={editingFlashcard.definition} />
            </label>
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                class="btn btn-secondary"
                disabled={performingAction}
                on:click={() => editingFlashcard = null}
            >Cancel</button>
            <button
                class="btn btn-primary"
                disabled={performingAction}
                on:click={editFlashcard}
            >Save</button>
        </svelte:fragment>
    </Modal>
{/if}

{#if showNewFlashcardModal}
    <Modal>
        <svelte:fragment slot="title">
            New Flashcard
        </svelte:fragment>
        <svelte:fragment slot="body">
            <label class="w-100 mb-2">
                Term
                <input type="text" class="form-control" bind:value={newFlashcardTerm} />
            </label>
            <label class="w-100">
                Definition
                <input type="text" class="form-control" bind:value={newFlashcardDefinition} />
            </label>
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                class="btn btn-secondary"
                disabled={performingAction}
                on:click={closeNewFlashcardModal}
            >Cancel</button>
            <button
                class="btn btn-primary"
                disabled={performingAction}
                on:click={newFlashcard}
            >Create</button>
        </svelte:fragment>
    </Modal>
{/if}

{#if showJumpToFlashcardModal}
    <Modal>
        <svelte:fragment slot="title">
            Jump to Flashcard
        </svelte:fragment>
        <svelte:fragment slot="body">
            <div class="d-flex align-items-center">
                <input
                    type="number"
                    class="form-control me-2 jump-to-flashcard-number"
                    min={1}
                    max={flashcards.length}
                    bind:value={jumpToFlashcardNumber}
                /> / {flashcards.length}
            </div>
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                class="btn btn-secondary"
                disabled={performingAction}
                on:click={cancelJumpToFlashcard}
            >Cancel</button>
            <button
                class="btn btn-primary"
                disabled={performingAction || !jumpToFlashcardNumber ||
                    jumpToFlashcardNumber > flashcards.length || jumpToFlashcardNumber < 0}
                on:click={jumpToFlashcard}
            >Go</button>
        </svelte:fragment>
    </Modal>
{/if}
