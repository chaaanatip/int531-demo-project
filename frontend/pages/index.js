// frontend/pages/index.js
import { useEffect, useState } from 'react'
import { fetchUsersClient } from '../lib/api'

export default function Home() {
    const [users, setUsers] = useState(null)
    const [err, setErr] = useState(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        let mounted = true
        fetchUsersClient()
            .then((data) => {
                if (mounted) setUsers(data)
            })
            .catch((e) => {
                if (mounted) setErr(e.message)
            })
            .finally(() => {
                if (mounted) setLoading(false)
            })
        return () => {
            mounted = false
        }
    }, [])

    return (
        <main style={{ padding: 24, fontFamily: 'system-ui, sans-serif' }}>
            <h1>Production-ready Next.js App via Kong Gateway</h1>

            <section style={{ marginTop: 20 }}>
                <h2>fetch users from backend</h2>
                {loading && <div>Loading…</div>}
                {err && <div style={{ color: 'crimson' }}>Error: {err}</div>}
                {users && (
                    <ul>
                        {users.map((u) => (
                            <li key={u.id}>
                                <strong>{u.id}</strong> <small>({u.name})</small>
                            </li>
                        ))}
                    </ul>
                )}
            </section>
        </main>
    )
}
