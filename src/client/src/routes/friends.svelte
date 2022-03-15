<script context="module">
    export { requiresAuthVerified as load } from '$lib/load'
</script>

<script>
    import Modal from '$lib/components/Modal.svelte'
    import Navbar from '$lib/components/Navbar.svelte'

    import { req } from '$lib/req'

    export let user

    let error = ''

    let loading = true
    let performingAction = false

    let incomingFriendRequests = []
    let outgoingFriendRequests = []
    let acceptedFriends = []

    Promise.all([
        req('friend/list/incoming', 'GET').then(({ friendRequests }) =>{
            incomingFriendRequests = friendRequests
        }),
        req('friend/list/outgoing', 'GET').then(({ friendRequests }) =>{
            outgoingFriendRequests = friendRequests
        }),
        req('friend/list/accepted', 'GET').then(({ friends }) =>{
            acceptedFriends = friends
        })
    ]).then(() => loading = false).catch(({ message }) => error = message)

    let incomingFriendRequestsListGroup
    let outgoingFriendRequestsListGroup
    let friendsListGroup

    const acceptIncomingFriendRequest = async friend => {
        try {
            performingAction = true

            await req(`friend/${friend.friendId}/accept`, 'POST')

            user.friendRequests--

            incomingFriendRequests = incomingFriendRequests.filter(({ friendId }) =>
                friendId !== friend.friendId)

            acceptedFriends = [...acceptedFriends, {
                friendId: friend.friendId,
                friend: {
                    userId: friend.from.userId,
                    username: friend.from.username
                }
            }]
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }

    const rejectIncomingFriendRequest = async friend => {
        try {
            performingAction = true

            await req(`friend/${friend.friendId}`, 'DELETE')

            incomingFriendRequests = incomingFriendRequests.filter(
                ({ friendId }) => friendId !== friend.friendId)
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }

    const acceptAllIncomingFriendRequests = async () => {
        try {
            performingAction = true

            for (const friend of incomingFriendRequests) {
                await req(`friend/${friend.friendId}/accept`, 'POST')
                user.friendRequests--
                acceptedFriends = [...acceptedFriends, {
                    friendId: friend.friendId,
                    friend: {
                        userId: friend.from.userId,
                        username: friend.from.username
                    }
                }]
            }

            incomingFriendRequests = []
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }

    const rejectAllIncomingFriendRequests = async () => {
        try {
            performingAction = true

            for (const { friendId } of incomingFriendRequests) {
                await req(`friend/${friendId}`, 'DELETE')
            }

            incomingFriendRequests = []
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }

    const cancelOutgoingFriendRequest = async friend => {
        try {
            performingAction = true

            await req(`friend/${friend.friendId}`, 'DELETE')

            outgoingFriendRequests = outgoingFriendRequests.filter(
                ({ friendId }) => friendId !== friend.friendId)
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }

    const unfriendFriend = async friend => {
        try {
            performingAction = true

            await req(`friend/${friend.friendId}`, 'DELETE')

            acceptedFriends = acceptedFriends.filter(
                ({ friendId }) => friendId !== friend.friendId)
        } catch ({ message }) {
            error = message
        } finally {
            performingAction = false
        }
    }
</script>

<style>
    .list-group-item:last-child {
        border-bottom-width: 1px;
    }
</style>

<svelte:head>
    <title>Studier &mdash; Friends</title>
</svelte:head>

<div class="vh-100 d-flex">
    <div class="p-3 flex-grow-1 d-flex flex-column">
        <Navbar active="friends" {user} />

        {#if loading}
            <h1>Loading...</h1>
        {:else}
            <div class="d-flex gap-3 flex-grow-1 overflow-hidden">
                <div class="w-50 d-flex flex-column gap-2">
                    <div class="d-flex flex-column h-50">
                        <h4 class="d-flex align-items-center justify-content-between">
                            <div class="d-flex align-items-center">
                                Incoming Friend Requests
                                <span class="badge rounded-pill bg-primary fs-6 ms-2">
                                    {incomingFriendRequests.length}
                                </span>
                            </div>
                            <div>
                                <button
                                    class="btn btn-outline-success btn-sm"
                                    on:click={acceptAllIncomingFriendRequests}
                                >Accept All</button>
                                <button
                                    class="btn btn-outline-danger btn-sm"
                                    on:click={rejectAllIncomingFriendRequests}
                                >Reject All</button>
                            </div>
                        </h4>
                        <div class="card flex-grow-1 overflow-auto">
                            <ul
                                class="list-group list-group-flush h-100"
                                bind:this={incomingFriendRequestsListGroup}
                            >
                                {#each incomingFriendRequests
                                    .sort((a, b) => b.timestamp - a.timestamp) as incomingFriendRequest, i}
                                    {@const from = incomingFriendRequest.from}
                                    <li
                                        class="list-group-item rounded-0"
                                        class:border-0={incomingFriendRequest.length === i + 1 &&
                                            incomingFriendRequestsListGroup?.scrollHeight >
                                            incomingFriendRequestsListGroup?.clientHeight}
                                    >
                                        <div class="d-flex align-items-center justify-content-between">
                                            <a
                                                href="/user/{from.userId}"
                                                class="text-break"
                                            >{from.username}</a>
                                            <div>
                                                <button
                                                    class="btn btn-success btn-sm"
                                                    on:click={() => acceptIncomingFriendRequest(
                                                        incomingFriendRequest)}
                                                >Accept</button>
                                                <button
                                                    class="btn btn-danger btn-sm"
                                                    on:click={() => rejectIncomingFriendRequest(
                                                        incomingFriendRequest)}
                                                >Reject</button>
                                            </div>
                                        </div>
                                    </li>
                                {/each}
                            </ul>
                        </div>
                    </div>
                    <div class="d-flex flex-column h-50">
                        <h4 class="d-flex align-items-center">
                            Outgoing Friend Requests
                            <span class="badge rounded-pill bg-primary fs-6 ms-2">
                                {outgoingFriendRequests.length}
                            </span>
                        </h4>
                        <div class="card flex-grow-1 overflow-auto">
                            <ul
                                class="list-group list-group-flush h-100"
                                bind:this={outgoingFriendRequestsListGroup}
                            >
                                {#each outgoingFriendRequests
                                    .sort((a, b) => b.timestamp - a.timestamp) as outgoingFriendRequest, i}
                                    {@const to = outgoingFriendRequest.to}
                                    <li
                                        class="list-group-item rounded-0"
                                        class:border-0={outgoingFriendRequest.length === i + 1 &&
                                            outgoingFriendRequestsListGroup?.scrollHeight >
                                            outgoingFriendRequestsListGroup?.clientHeight}
                                    >
                                        <div class="d-flex align-items-center justify-content-between">
                                            <a
                                                href="/user/{to.userId}"
                                                class="text-break"
                                            >{to.username}</a>
                                            <div>
                                                <button
                                                    class="btn btn-warning btn-sm"
                                                    on:click={() => cancelOutgoingFriendRequest(
                                                        outgoingFriendRequest)}
                                                >Cancel</button>
                                            </div>
                                        </div>
                                    </li>
                                {/each}
                            </ul>
                        </div>
                    </div>
                </div>
                <div class="w-50 d-flex flex-column">
                    <h4>Friends</h4>
                    <div class="card flex-grow-1 overflow-auto">
                        <ul class="list-group list-group-flush h-100" bind:this={friendsListGroup}>
                            {#each acceptedFriends.sort((a, b) =>
                                a.friend.username.localeCompare(b.friend.username)) as acceptedFriend, i}
                                <li
                                    class="list-group-item rounded-0"
                                    class:border-0={acceptedFriends.length === i + 1 &&
                                        friendsListGroup?.scrollHeight >
                                        friendsListGroup?.clientHeight}
                                >
                                    <div class="d-flex align-items-center justify-content-between">
                                        <a
                                            href="/user/{acceptedFriend.friend.userId}"
                                            class="text-break"
                                        >{acceptedFriend.friend.username}</a>
                                        <button
                                            class="btn btn-danger btn-sm"
                                            on:click={() => unfriendFriend(acceptedFriend)}
                                        >Unfriend</button>
                                    </div>
                                </li>
                            {/each}
                        </ul>
                    </div>
                </div>
            </div>
        {/if}
    </div>
</div>

{#if error}
    <Modal>
        <svelte:fragment slot="title">
            Error
        </svelte:fragment>
        <svelte:fragment slot="body">
            {error}
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button class="btn btn-primary" on:click={() => error = ''}>OK</button>
        </svelte:fragment>
    </Modal>
{/if}
