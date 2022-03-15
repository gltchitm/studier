import { goto } from '$app/navigation'

import { req } from './req'

const logout = async () => {
    await req('auth/logout', 'DELETE', '')
    for (const key of Object.keys(sessionStorage)) {
        if (key.startsWith('studier_deck_token_')) {
            sessionStorage.removeItem(key)
        }
    }
    goto('/login')
}

export { logout }
