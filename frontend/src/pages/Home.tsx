import { useEffect, useState } from "react"
import { useAuth } from "../context/AuthContext"
import type {
  GithubComRhajizadaGazetteInternalServiceFeed,
} from "../api/data-contracts"
import { Spinner } from "@/components/ui/spinner"

export default function Home() {
  const { api, logout } = useAuth()

  const [feeds, setFeeds] = useState<GithubComRhajizadaGazetteInternalServiceFeed[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setLoading(true)
    api
      .feedsList({ subscribedOnly: false, limit: 20, offset: 0 })
      .then(response => {
        // on success, response.data.feeds is your array
        setFeeds(response.data.feeds ?? [])
      })
      .catch(err => {
        // if your token expired server-side, you might get a 401
        if (err.error === "Unauthorized") {
          logout()
        } else {
          setError("Failed to load feeds")
          console.error(err)
        }
      })
      .finally(() => setLoading(false))
  }, [api, logout])

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <Spinner />
      </div>
    )
  }

  if (error) {
    return (
      <div className="p-4 text-red-600">
        {error}
      </div>
    )
  }

  return (
    <div className="p-6 space-y-4">
      {feeds.length === 0 ? (
        <p>No feeds yet—go add some!</p>
      ) : (
        feeds.map(feed => (
          <div
            key={feed.id}
            className="p-4 border rounded hover:shadow cursor-pointer"
            onClick={() => {
              /* e.g. navigate to /feeds/${feed.id} */
            }}
          >
            <h3 className="font-semibold">{feed.title}</h3>
            <p className="text-sm text-muted-foreground">
              {feed.description}
            </p>
          </div>
        ))
      )}
    </div>
  )
}

