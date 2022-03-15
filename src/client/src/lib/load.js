import { rawReq } from './req'

const requiresAuthVerified = async () => {
    const response = await rawReq('account', 'GET')

    if (response.ok) {
        const user = await response.json()

        if (!user.verified) {
            return { status: 307, redirect: '/verify' }
        }

        return { props: { user } }
    } else if (response.status === 403) {
        return { status: 307, redirect: '/login' }
    }

    throw new Error((await response.json()).error)
}

const requiresAuthNotVerified = async () => {
    const response = await rawReq('account', 'GET')

    if (response.ok) {
        const user = await response.json()

        if (user.verified) {
            return { status: 307, redirect: '/' }
        }

        return { props: { user } }
    } else if (response.status === 403) {
        return { status: 307, redirect: '/login' }
    }

    throw new Error((await response.json()).error)
}

const requiresNoAuth = async () => {
    const response = await rawReq('account', 'GET')

    if (response.ok) {
        return { status: 307, redirect: '/' }
    } else if (response.status !== 403) {
        throw new Error((await response.json()).error)
    }

    return {}
}

const requiresAnyAuth = async () => {
    const response = await rawReq('account', 'GET')

    if (response.ok) {
        return { props: { user: await response.json() } }
    } else if (response.status === 403) {
        return { props: { user: null } }
    }

    throw new Error((await response.json()).error)
}

export { requiresAuthVerified, requiresAuthNotVerified, requiresNoAuth, requiresAnyAuth }
