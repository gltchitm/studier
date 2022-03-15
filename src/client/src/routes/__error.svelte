<script>
    import { dev } from '$app/env'
    import { page } from '$app/stores'
</script>

<svelte:head>
    <title>Studier &mdash; Error {$page.status}</title>
</svelte:head>

<div class="p-3">
    <div class="text-center container">
        <div class="card">
            <div class="card-header">
                <h3>Error {$page.status}</h3>
            </div>
            <div class="card-body">
                {$page.error.message}
                {#if dev && $page.error.stack}
                    <pre>{$page.error.stack}</pre>
                {/if}
            </div>
        </div>
    </div>
</div>
