<script context="module">
    export { requiresAuthVerified as load } from '$lib/load'
</script>

<script>
    import Navbar from '$lib/components/Navbar.svelte'

    import { goto } from '$app/navigation'
    import { page } from '$app/stores'

    export let user

    $: active = $page.url.pathname.slice($page.url.pathname.lastIndexOf('/') + 1)
</script>

<style>
    .pills {
        width: 200px;
    }
</style>

<div class="p-3">
    <Navbar {user} />

    <div class="d-flex align-items-start gap-3">
        <div class="nav flex-column nav-pills pills">
            <button
                class="nav-link"
                class:active={active === 'study'}
                on:click={() => goto('study')}
            >Study</button>
            <button
                class="nav-link"
                class:active={active === 'writing'}
                on:click={() => goto('writing')}
            >Writing</button>
            <button
                class="nav-link"
                class:active={active === 'info'}
                on:click={() => goto('info')}
            >Deck Info</button>
        </div>
        <div class="tab-content w-100">
            <slot />
        </div>
    </div>
</div>
