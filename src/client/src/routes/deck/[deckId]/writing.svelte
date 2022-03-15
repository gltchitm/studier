<script>
    import { goto } from '$app/navigation'
    import { page } from '$app/stores'

    import { rawReq, req } from '$lib/req'

    const { deckId } = $page.params

    let error = ''

    let loadingDeck = true
    let updatingPinned = false

    let name = ''
    let description = ''
    let flashcards = []
    let editable = false
    let author = null
    let pinned = false

    const writingConfig = {
        answerWith: 'definition',
        allowOverride: true,
        caseSensitivity: false,
        shuffle: false
    }

    let started = false
    let remainingFlashcards = []
    let correct = 0
    let incorrect = 0

    let answer = ''

    let result = ''

    $: question = remainingFlashcards?.[0]?.[writingConfig.answerWith === 'term' ? 'definition' : 'term']
    $: correctAnswer = remainingFlashcards?.[0]?.[writingConfig.answerWith]

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
            goto(`/deck/${deckId}/writing/${json.unlockable ? 'unlock' : 'locked'}`,
                { replaceState: true })
        } else {
            error = json.error
        }
    }).catch(({ message }) => {
        error = message
    })

    const pinDeck = async () => {
        try {
            error = ''
            updatingPinned = true

            await req(`deck/${deckId}/pin`, 'POST')

            pinned = true
        } catch ({ message }) {
            error = message
        } finally {
            updatingPinned = false
        }
    }

    const unpinDeck = async () => {
        try {
            error = ''
            updatingPinned = true

            await req(`deck/${deckId}/pin`, 'DELETE')

            pinned = false
        } catch ({ message }) {
            error = message
        } finally {
            updatingPinned = false
        }
    }

    const startWriting = () => {
        if (writingConfig.shuffle) {
            remainingFlashcards = []

            let indexes = Array(flashcards.length).fill().map((_, i) => i)

            while (indexes.length) {
                const index = indexes[Math.floor(Math.random() * indexes.length)]
                remainingFlashcards = [...remainingFlashcards, flashcards[index]]
                indexes.splice(indexes.indexOf(index), 1)
            }
        } else {
            remainingFlashcards = [...flashcards]
        }

        answer = ''
        result = ''
        correct = 0
        incorrect = 0
        started = true
    }

    const submit = () => {
        if (answer === correctAnswer || (!writingConfig.caseSensitivity &&
            answer.toLowerCase() === correctAnswer.toLowerCase())) {
            correct++
            result = 'correct'
        } else {
            incorrect++
            result = 'incorrect'
        }
    }

    const next = () => {
        answer = ''
        result = ''
        remainingFlashcards = remainingFlashcards.slice(1)
    }

    const override = () => {
        incorrect--
        correct++
        next()
    }
</script>

<svelte:head>
    <title>Studier &mdash; Writing</title>
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
                disabled={updatingPinned}
                on:click={() => unpinDeck()}
            >
                Unpin
            </button>
        {:else}
            <button
                class="btn btn-success btn-sm mb-3"
                disabled={updatingPinned}
                on:click={() => pinDeck()}
            >
                Pin
            </button>
        {/if}
    {/if}

    {#if started}
        <div class="card mb-3">
            <div class="card-body d-flex gap-3">
                <div class="d-flex align-items-center flex-column">
                    <h1 class="text-success mb-0">{correct}</h1>
                    <small>correct</small>
                </div>
                <div class="d-flex align-items-center flex-column">
                    <h1 class="text-danger mb-0">{incorrect}</h1>
                    <small>incorrect</small>
                </div>
            </div>
        </div>
    {/if}

    {#if !loadingDeck}
        {#if !started}
            <div class="mt-3">
                <h2 class="text-center">Writing</h2>

                <div class="d-flex gap-3">
                    <div class="mb-3 d-flex flex-column" style="width: 175px">
                        <label>
                            Answer With
                            <select class="form-select" bind:value={writingConfig.answerWith}>
                                <option value="term">Term</option>
                                <option value="definition">Definition</option>
                            </select>
                        </label>
                    </div>

                    <div class="mb-3 d-flex flex-column" style="width: 175px">
                        <label>
                            Allow Override
                            <select class="form-select" bind:value={writingConfig.allowOverride}>
                                <option value={true}>On</option>
                                <option value={false}>Off</option>
                            </select>
                        </label>
                    </div>
                </div>

                <div class="d-flex gap-3">
                    <div class="mb-3 d-flex flex-column" style="width: 175px">
                        <label>
                            Case Sensitivity
                            <select class="form-select" bind:value={writingConfig.caseSensitivity}>
                                <option value={true}>On</option>
                                <option value={false}>Off</option>
                            </select>
                        </label>
                    </div>

                    <div class="mb-3 d-flex flex-column" style="width: 175px">
                        <label>
                            Shuffle
                            <select class="form-select" bind:value={writingConfig.shuffle}>
                                <option value={true}>On</option>
                                <option value={false}>Off</option>
                            </select>
                        </label>
                    </div>
                </div>

                <button
                    class="btn btn-primary w-100"
                    on:click={startWriting}
                >Start Writing</button>
            </div>
        {:else if !remainingFlashcards.length}
            <h2 class="my-2">
                {correct}/{flashcards.length} = {Math.trunc(correct / flashcards.length * 100)}%
            </h2>
            <button class="btn btn-primary mt-2" on:click={() => started = false}>Finish</button>
        {:else}
            <div class="card" style="width: 650px;">
                <div class="card-header d-flex align-items-center justify-content-between">
                    {flashcards.length - remainingFlashcards.length + 1}/{flashcards.length}
                    <div class="d-flex gap-2">
                        <button class="btn btn-danger btn-sm" on:click={() => started = false}>
                            Abandon
                        </button>
                    </div>
                </div>
                <div class="card-body">
                    {#if result === 'correct'}
                        <h2 class="text-success">Correct</h2>

                        <div>
                            <small>Question:</small>
                            <h5>{question}</h5>
                        </div>
                        <div>
                            <small>Your answer:</small>
                            <h5 class="text-success">{answer}</h5>
                        </div>
                        <div>
                            <small>Correct answer:</small>
                            <h5 class="text-success">{correctAnswer}</h5>
                        </div>
                        <div class="d-flex align-items-center justify-content-end w-100">
                            <button class="btn btn-success" on:click={next}>Next</button>
                        </div>
                    {:else if result === 'incorrect'}
                        <h2 class="text-danger">Incorrect</h2>

                        <div>
                            <small>Question:</small>
                            <h5>{question}</h5>
                        </div>
                        <div>
                            <small>Your answer:</small>
                            <h5 class="text-danger">{answer}</h5>
                        </div>
                        <div>
                            <small>Correct answer:</small>
                            <h5 class="text-success">{correctAnswer}</h5>
                        </div>
                        <div class="d-flex align-items-center justify-content-end w-100">
                            {#if writingConfig.allowOverride}
                                <button
                                    class="btn btn-warning me-2"
                                    on:click={override}
                                >Override</button>
                            {/if}
                            <button class="btn btn-danger" on:click={next}>Next</button>
                        </div>
                    {:else}
                        <h4 class="mb-2">{question}</h4>
                        <label class="w-100 mb-2">
                            {#if writingConfig.answerWith === 'term'}
                                Term
                            {:else}
                                Definition
                            {/if}
                            <input
                                type="text"
                                class="form-control"
                                maxlength={flashcards.sort((a, b) => b[writingConfig.answerWith].length -
                                    a[writingConfig.answerWith].length)[0][writingConfig.answerWith].length}
                                bind:value={answer}
                            />
                        </label>
                        <div class="d-flex align-items-center justify-content-end w-100">
                            <button
                                class="btn btn-primary"
                                disabled={!answer.length}
                                on:click={submit}
                            >Submit</button>
                        </div>
                    {/if}
                </div>
            </div>
        {/if}
    {/if}
</div>
