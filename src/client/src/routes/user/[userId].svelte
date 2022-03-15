<script context="module">
    export { requiresAuthVerified as load } from '$lib/load'
</script>

<script>
    import { page } from '$app/stores'

    import Navbar from '$lib/components/Navbar.svelte'

    import { req } from '$lib/req'

    export let user

    let error = ''

    let loading = true
    let performingAction = false

    let username
    let decks
    let friendState
    let friendId

    const { userId } = $page.params

    req(`user/${userId}`, 'GET').then(user => {
        username = user.username
        decks = user.decks
        friendState = user.friendState
        friendId = user.friendId

        loading = false
    }).catch(({ message }) => {
        error = message
    })

    const sendFriendRequest = async () => {
        try {
            performingAction = true

            const { friendId: newFriendId } = await req(`user/${userId}/friend`, 'POST')

            friendId = newFriendId

            friendState = 'requested'
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }

    const acceptFriendRequest = async () => {
        try {
            performingAction = true

            await req(`friend/${friendId}/accept`, 'POST')

            friendState = 'friended'
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }

    const unfriendFriend = async () => {
        try {
            performingAction = true

            await req(`friend/${friendId}`, 'DELETE')

            friendState = 'unfriended'
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }
</script>

<svelte:head>
    <title>Studier &mdash; {username ?? 'User'}</title>
</svelte:head>

<div class="p-3">
    <Navbar {user} />

    {#if error}
        <div class="alert alert-danger">{error}</div>
    {/if}

    {#if loading}
        <h1>Loading...</h1>
    {:else}
        <div class="d-flex align-items-center justify-content-between mb-3">
            <h1>{username}</h1>
            {#if $page.params.userId !== user.userId}
                {#if friendState === 'unfriended'}
                    <button
                        class="btn btn-success"
                        disabled={performingAction}
                        on:click={sendFriendRequest}
                    >
                        Send Friend Request
                    </button>
                {:else if friendState === 'requested'}
                    <button
                        class="btn btn-warning"
                        disabled={performingAction}
                        on:click={unfriendFriend}
                    >
                        Cancel Friend Request
                    </button>
                {:else if friendState === 'requestedByFriend'}
                    <div>
                        <button
                            class="btn btn-success"
                            disabled={performingAction}
                            on:click={acceptFriendRequest}
                        >
                            Accept Friend Request
                        </button>
                        <button
                            class="btn btn-danger"
                            disabled={performingAction}
                            on:click={unfriendFriend}
                        >
                            Reject Friend Request
                        </button>
                    </div>
                {:else if friendState === 'friended'}
                    <button
                        class="btn btn-danger"
                        disabled={performingAction}
                        on:click={unfriendFriend}
                    >
                        Unfriend
                    </button>
                {/if}
            {/if}
        </div>
        <h5>Decks</h5>
        {#if decks.length}
            <ul class="list-group">
                {#each decks.sort((a, b) => a.name.localeCompare(b.name)) as deck}
                    <a
                        href="/deck/{deck.deckId}/study"
                        class="list-group-item list-group-item-action"
                    >{deck.name}</a>
                {/each}
            </ul>
        {:else}
            <div class="alert alert-info">
                This user does not have any decks.
            </div>
        {/if}
    {/if}
</div>
