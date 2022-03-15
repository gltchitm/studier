<script context="module">
    import { requiresAuthVerified } from '$lib/load'

    export const load = ({ params }) => {
        if (!['study', 'writing', 'info'].includes(params.destination)) {
            return { status: 404 }
        }

        return requiresAuthVerified()
    }
</script>

<script>
    import Navbar from '$lib/components/Navbar.svelte'

    export let user
</script>

<div class="p-3">
    <Navbar {user} />

    <slot />
</div>
