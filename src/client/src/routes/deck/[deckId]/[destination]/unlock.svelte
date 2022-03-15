<script>
    import { goto } from '$app/navigation'
    import { page } from '$app/stores'

    import { rawReq, req } from '$lib/req'

    const { deckId, destination } = $page.params

    let loadingDeck = true
    let performingAction = false

    let error = ''

    let name = ''

    let password = ''

    const token = sessionStorage.getItem('studier_deck_token_' + deckId)
    rawReq(`deck/${deckId}${token ? '?token=' + token : ''}`, 'GET').then(async response => {
        const json = await response.json().catch(() => ({ error: response.statusText }))

        if (response.ok) {
            goto(`/deck/${deckId}/${destination}`)
        } else if (response.status === 403 && !json.unlockable) {
            goto(`/deck/${deckId}/${destination}/locked`)
        } else if (response.status === 403) {
            name = json.name
            loadingDeck = false
        } else {
            error = json.error
        }
    }).catch(({ message }) => {
        error = message
    })

    const unlock = async () => {
        try {
            error = ''
            performingAction = true

            const { token } = await req(`deck/${deckId}/unlock`, 'POST', { password })

            sessionStorage.setItem(`studier_deck_token_${deckId}`, token)

            goto(`/deck/${deckId}/${destination}`)
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }
</script>

<svelte:head>
    <title>Studier &mdash; Unlock Deck</title>
</svelte:head>

{#if error}
    <div class="alert alert-danger">{error}</div>
{/if}

{#if loadingDeck}
    <h1>Loading...</h1>
{:else}
    <h1>Enter Password</h1>

    <p class="text-break">Enter the password to unlock the deck "{name}".</p>

    <label class="mb-2">
        Password
        <input
            type="password"
            class="form-control"
            disabled={performingAction}
            bind:value={password}
        />
    </label>

    <button
        class="btn btn-primary d-block"
        disabled={performingAction}
        on:click={unlock}
    >Unlock</button>
{/if}
