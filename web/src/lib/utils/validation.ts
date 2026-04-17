export const USERNAME_REGEX = /^[a-zA-Z0-9_]+$/
export const NICKNAME_REGEX = /^[a-zA-Z0-9_가-힣]+$/

export function validateUsername(nickname: string): string {
    if (nickname.length < 2) return '아이디는 최소 4자리입니다.'
    if (!USERNAME_REGEX.test(nickname)) return '아이디에 특수문자는 사용할 수 없습니다.'
    return ''
}

export function validateNickname(nickname: string): string {
    if (nickname.length < 2) return '닉네임은 최소 2자리입니다.'
    if (!NICKNAME_REGEX.test(nickname)) return '닉네임에 특수문자는 사용할 수 없습니다.'
    return ''
}