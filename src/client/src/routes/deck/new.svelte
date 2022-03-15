<script context="module">
    export { requiresAuthVerified as load } from '$lib/load'
</script>

<script>
    import { goto } from '$app/navigation'

    import FlashcardEditor from '$lib/components/FlashcardEditor.svelte'
    import Navbar from '$lib/components/Navbar.svelte'

    import { req } from '$lib/req'

    export let user

    let nextFlashcardId = 0

    const emptyFlashcard = () => ({
        term: '',
        definition: '',
        flashcardId: nextFlashcardId++
    })

    let submitting = false

    let name = ''
    let description = ''
    let access = 'everyone'
    let password = ''
    let flashcards = [emptyFlashcard(), emptyFlashcard(), emptyFlashcard()]

    let error = ''

    const newFlashcard = () => {
        if (submitting) {
            return
        }

        flashcards = [...flashcards, emptyFlashcard()]
    }

    const createDeck = async () => {
        error = ''
        submitting = true

        try {
            const { deckId } = await req('deck/new', 'POST', {
                name,
                description,
                access,
                password,
                flashcards: flashcards.map(flashcard => ({
                    term: flashcard.term,
                    definition: flashcard.definition
                }))
            })

            goto(`/deck/${deckId}/study`)
        } catch ({ message }) {
            error = message
        } finally {
            submitting = false
        }
    }
</script>

<svelte:head>
    <title>Studier &mdash; New Deck</title>
</svelte:head>

<div class="p-3">
    <Navbar showCreateDeck={false} {user} />

    <div class="d-flex align-items-center justify-content-between">
        <h2>New Deck</h2>
        <button
            class="btn btn-primary"
            disabled={submitting}
            on:click={createDeck}
        >Create</button>
    </div>

    {#if error}
        <div class="alert alert-danger w-100 my-3">
            {error}
        </div>
    {/if}

    <div class="d-flex gap-2 mb-4">
        <div class="w-50">
            <label class="w-100 mb-2">
                Name
                <input
                    type="text"
                    class="form-control"
                    disabled={submitting}
                    bind:value={name}
                />
            </label>
            <label class="w-100">
                Description
                <input
                    type="text"
                    class="form-control"
                    disabled={submitting}
                    bind:value={description}
                />
            </label>
        </div>
        <div class="w-50">
            <label class="w-100 mb-2">
                Access
                <select
                    class="form-select"
                    disabled={submitting}
                    bind:value={access}
                    on:change={() => password = ''}
                >
                    <option value="everyone">Everyone</option>
                    <option value="friends">My friends</option>
                    <option value="password">People with the password</option>
                    <option value="me">Only me</option>
                </select>
            </label>

            {#if access === 'password'}
                <label class="w-100">
                    Password
                    <input
                        type="password"
                        class="form-control"
                        disabled={submitting}
                        bind:value={password}
                    />
                </label>
            {/if}
        </div>
    </div>

    {#each flashcards as flashcard, index (flashcard.flashcardId)}
        <FlashcardEditor
            number={index + 1}
            {flashcard}
            disabled={submitting}
            on:termChanged={({ detail }) => flashcard.term = detail}
            on:definitionChanged={({ detail }) => flashcard.definition = detail}
            on:delete={() => flashcards = flashcards.filter(
                ({ flashcardId }) => flashcard.flashcardId !== flashcardId)}
        />
    {/each}

    <FlashcardEditor
        number={flashcards.length + 1}
        flashcard={{ term: '', definition: '' }}
        phantom
        on:newFlashcard={newFlashcard}
    />
</div>
