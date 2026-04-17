export const decodeJwtPayload = (token: string) => {
    const base64 = token.split('.')[1]
        .replace(/-/g, '+')
        .replace(/_/g, '/')
    
    const jsonStr = decodeURIComponent(
        atob(base64)
            .split('')
            .map(c => '%' + c.charCodeAt(0).toString(16).padStart(2, '0'))
            .join('')
    )
    return JSON.parse(jsonStr)
}