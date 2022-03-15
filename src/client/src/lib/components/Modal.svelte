<script>
    import { fade, fly } from 'svelte/transition'
    import { quintOut } from 'svelte/easing'

    let modal

    const keydown = event => {
        if (event.key === 'Tab') {
            const tabbable = [...modal.querySelectorAll('*')].filter(element => element.tabIndex >= 0)

            const focused = tabbable.indexOf(document.activeElement)

            let nextIndex = (focused + (event.shiftKey ? -1 : 1)) % tabbable.length
            if (nextIndex === -1 && event.shiftKey) {
                nextIndex = tabbable.length - 1
            }

            tabbable[nextIndex].focus()

            event.preventDefault()
        }
    }
</script>

<svelte:window on:keydown={keydown} />

<div class="modal-backdrop show" transition:fade={{ duration: 150 }}></div>
<div
    class="modal d-block"
    bind:this={modal}
    in:fly={{ y: -50, duration: 300 }}
    out:fly={{ y: -50, duration: 300, easing: quintOut }}
>
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title">
                    <slot name="title" />
                </h5>
            </div>
            <div class="modal-body">
                <slot name="body" />
            </div>
            <div class="modal-footer">
                <slot name="footer" />
            </div>
        </div>
    </div>
</div>
