<script>
    import { goto } from '$app/navigation'
    import { page } from '$app/stores'

    import { rawReq } from '$lib/req'

    const { deckId, destination } = $page.params

    let loading = true

    let error = ''

    let name = ''

    const token = sessionStorage.getItem('studier_deck_token_' + deckId)
    rawReq(`deck/${deckId}${token ? '?token=' + token : ''}`, 'GET').then(async response => {
        const json = await response.json().catch(() => ({ error: response.statusText }))

        if (response.ok) {
            goto(`/deck/${deckId}/${destination}`)
        } else if (response.status === 403 && json.unlockable) {
            goto(`/deck/${deckId}/${destination}/unlock`)
        } else if (response.status === 403) {
            name = json.name
            loading = false
        } else {
            error = json.error
        }
    }).catch(({ message }) => {
        error = message
    })
</script>

<svelte:head>
    <title>Studier &mdash; Access Denied</title>
</svelte:head>

{#if error}
    <div class="alert alert-danger">{error}</div>
{/if}

{#if loading}
    <h1>Loading...</h1>
{:else}
    <h1>Access Denied</h1>

    <p class="text-break">You do not have access to the deck "{name}".</p>
{/if}
