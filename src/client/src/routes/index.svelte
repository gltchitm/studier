<script context="module">
    export { requiresAnyAuth as load } from '$lib/load'
</script>

<script>
    import { goto } from '$app/navigation'

    import { logout } from '$lib/logout'

    export let user
</script>

<svelte:head>
    <title>Studier &mdash; Welcome</title>
</svelte:head>

<div class="d-flex flex-column align-items-center justify-content-center vh-100 vw-100">
    <h1 class="display-1">Studier</h1>
    <h3 class="fw-light">Open-source online studying platform</h3>

    {#if user}
        <span class="mb-2">Logged in as {user.username}</span>
        <div class="d-flex gap-2">
            <button class="btn btn-primary" on:click={() => goto('/home')}>Home</button>
            <button class="btn btn-danger" on:click={logout}>Log Out</button>
        </div>
    {:else}
        <div class="d-flex gap-2">
            <button class="btn btn-primary" on:click={() => goto('/login')}>Log In</button>
            <button class="btn btn-secondary" on:click={() => goto('/signup')}>Sign Up</button>
        </div>
    {/if}
</div>
