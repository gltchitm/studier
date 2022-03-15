<script>
    import { createEventDispatcher } from 'svelte'

    import { scale } from 'svelte/transition'
    import { quintOut } from 'svelte/easing'

    export let number
    export let flashcard
    export let disabled = false
    export let phantom = false

    const dispatch = createEventDispatcher()

    const newFlashcard = () => {
        if (phantom) {
            dispatch('newFlashcard')
        }
    }
</script>

<style>
    .phantom-card {
        cursor: pointer;
    }
</style>

<div
    class="card mb-3"
    class:opacity-50={phantom}
    class:phantom-card={phantom}
    class:user-select-none={phantom}
    on:click={newFlashcard}
    out:scale|local={{ duration: phantom ? 0 : 250, easing: quintOut }}
>
    <div
        class="card-header d-flex align-items-center justify-content-between"
        class:pe-none={phantom}
    >
        Flashcard {number}
        <button
            class="btn btn-danger btn-sm"
            {disabled}
            on:click={() => dispatch('delete')}
        >Delete</button>
    </div>
    <div
        class="card-body d-flex align-items-center justify-content-between gap-3"
        class:pe-none={phantom}
    >
        <label class="w-50">
            Term
            <input
                type="text"
                class="form-control"
                tabindex={phantom ? -1 : 0}
                {disabled}
                bind:value={flashcard.term}
                on:input={() => dispatch('termChanged', flashcard.term)}
            />
        </label>
        <label class="w-50">
            Definition
            <input
                type="text"
                class="form-control"
                tabindex={phantom ? -1 : 0}
                {disabled}
                bind:value={flashcard.definition}
                on:input={() => dispatch('definitionChanged', flashcard.definition)}
            />
        </label>
    </div>
</div>
