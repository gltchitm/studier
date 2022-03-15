<script>
    import { goto } from '$app/navigation'

    import { logout } from '$lib/logout'

    export let user

    export let active = null
    export let showCreateDeck = true

    let dropdownContainer
    let showDropdown = false

    const mousedown = ({ target }) => {
        if (showDropdown && !dropdownContainer.contains(target)) {
            showDropdown = false
        }
    }

    const keydown = event => {
        if (showDropdown && event.key === 'Tab') {
            const tabbable = [...dropdownContainer.querySelectorAll('*')].filter(
                element => element.tabIndex >= 0)

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
<style>
    .dropdown-menu {
        right: 12px;
        border-radius: 8px;
    }
    .pointer {
        cursor: pointer;
    }
    .log-out:hover {
        background: var(--bs-danger);
    }
</style>

<svelte:window on:mousedown={mousedown} on:keydown={keydown} />

<nav class="navbar navbar-expand-lg navbar-dark bg-dark rounded-3 mb-3">
    <div class="container-fluid d-flex align-items-center justify-content-center justify-content-md-between">
        <a class="navbar-brand" href="/">Studier</a>
        <ul class="navbar-nav">
            <li class="nav-item">
                <a href="/home" class="nav-link" class:active={active === 'home'}>Home</a>
            </li>
            <li class="nav-item">
                <a href="/friends" class="nav-link" class:active={active === 'friends'}>
                    Friends
                    <span class="badge rounded-pill bg-danger">{user.friendRequests}</span>
                </a>
            </li>
        </ul>
        <div>
            <button
                class="btn btn-outline-light btn-sm"
                class:invisible={!showCreateDeck}
                on:click={() => goto('/deck/new')}
            >
                Create Deck
            </button>
        </div>
        <div bind:this={dropdownContainer}>
            <button
                class="btn btn-dark btn-sm dropdown-toggle"
                on:click={() => showDropdown = !showDropdown}
            >
                {user.username}
            </button>
            <ul
                class="dropdown-menu dropdown-menu-dark p-2 d-flex flex-column gap-1"
                class:d-flex={showDropdown}
            >
                <li
                    class="dropdown-item rounded pointer"
                    tabindex={0}
                    on:click={() => goto('/user/' + user.userId)}
                >My Profile</li>
                <li
                    class="dropdown-item rounded pointer"
                    tabindex={0}
                    on:click={() => goto('/settings')}
                >Settings</li>
                <li
                    class="dropdown-item rounded pointer log-out"
                    tabindex={0}
                    on:click={logout}
                >Log Out</li>
            </ul>
        </div>
    </div>
</nav>
