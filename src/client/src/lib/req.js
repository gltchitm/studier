const rawReq = (url, method, body = {}) => {
    const request = { method }

    if (method !== 'GET') {
        request.headers = { 'Content-Type': 'application/json' }
        request.body = JSON.stringify(body)
    }

    return fetch(new URL(url, new URL('/api/', window.location.href)).href, request)
}

const req = async (url, method, body) => {
    const response = await rawReq(url, method, body)

    const json = await response.json().catch(() => ({ error: response.statusText }))
    if (!response.ok) {
        throw new Error(json.error)
    }

    return json
}

export { rawReq, req }
