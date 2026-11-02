const TOKEN_KEY = 'wtv_admin_token'
const PROFILE_KEY = 'wtv_admin_profile'

export function getToken() {
  return localStorage.getItem(TOKEN_KEY) || ''
}

export function getProfile() {
  try {
    return JSON.parse(localStorage.getItem(PROFILE_KEY) || 'null')
  } catch {
    return null
  }
}

export function setSession(token, profile) {
  localStorage.setItem(TOKEN_KEY, token)
  localStorage.setItem(PROFILE_KEY, JSON.stringify(profile))
}

export function clearSession() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(PROFILE_KEY)
}

export function hasPerm(profile, perm) {
  return Array.isArray(profile?.permissions) && profile.permissions.includes(perm)
}
