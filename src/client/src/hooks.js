export const handle = async ({ event, resolve }) => resolve(event, { ssr: false })

export const getSession = event => ({ user: event.locals.user })
