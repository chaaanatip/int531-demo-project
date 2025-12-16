const API_BASE = process.env.NEXT_PUBLIC_API_BASE || 'http://10.13.104.89:8000'

export async function fetchUsersClient() {
  const url = `${API_BASE}/api/users`
  const res = await fetch(url)
  if (!res.ok) {
    const text = await res.text()
    throw new Error(`fetchUsersClient failed: ${res.status} ${text}`)
  }
  return res.json()
}
