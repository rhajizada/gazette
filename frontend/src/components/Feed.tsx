import { useState } from "react"
import type {
  GithubComRhajizadaGazetteInternalServiceFeed as FeedType,
} from "../api/data-contracts"
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { useAuth } from "../context/AuthContext"
import { Link } from "react-router-dom";


interface FeedProps {
  feed: FeedType
}

export function Feed({ feed }: FeedProps) {
  const { api, logout } = useAuth()
  const [subscribed, setSubscribed] = useState(feed.subscribed ?? false)
  const [loading, setLoading] = useState(false)

  const handleToggle = async () => {
    if (!feed.id) return
    setLoading(true)
    try {
      if (subscribed) {
        await api.feedsSubscribeDelete(feed.id, { format: "json" })
        setSubscribed(false)
      } else {
        await api.feedsSubscribeUpdate(feed.id, { format: "json" })
        setSubscribed(true)
      }
    } catch (err: any) {
      if (err.error === "Unauthorized") {
        logout()
      } else {
        console.error("Subscription error:", err)
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card className="w-full max-w-lg mx-auto">
      {feed.image?.url && (
        <img
          src={feed.image.url}
          alt={feed.image.title ?? feed.title}
          className="w-full h-48 object-cover rounded-t-lg"
        />
      )}


      <CardHeader>
        <CardTitle>
          <Link to={`/feed/${feed.id}`}>{feed.title}</Link>
        </CardTitle>
        <CardDescription>{feed.description}</CardDescription>
      </CardHeader>


      <CardContent className="space-y-2 px-6">
        {feed.authors && feed.authors.length > 0 && (
          <p className="text-sm text-muted-foreground">
            <strong>Author{feed.authors.length > 1 ? "s" : ""}:</strong>{" "}
            {feed.authors
              .map((a) => (a.name ? a.name : a.email))
              .filter(Boolean)
              .join(", ")}
          </p>
        )}

        {feed.categories && feed.categories.length > 0 && (
          <p className="text-sm text-muted-foreground">
            <strong>Categories:</strong> {feed.categories.join(", ")}
          </p>
        )}

        {feed.language && (
          <p className="text-sm text-muted-foreground">
            <strong>Language:</strong> {feed.language}
          </p>
        )}




        {feed.links?.length! > 0 && (
          <p className="text-sm text-muted-foreground">
            <strong>Source{feed.links!.length > 1 ? "s" : ""}:</strong>{" "}
            {feed.links!.map((url, i) => (
              <span key={url}>
                <a
                  href={url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="underline hover:text-primary"
                >
                  {url}
                </a>
                {i < feed.links!.length - 1 ? ", " : null}
              </span>
            ))}
          </p>
        )}

        {feed.updated_parsed && (
          <p className="text-xs text-muted-foreground">
            <strong>Last updated:</strong>{" "}
            {new Date(feed.updated_parsed).toLocaleString()}
          </p>
        )}
      </CardContent>

      <CardFooter className="justify-end">
        <Button
          variant={subscribed ? "outline" : "default"}
          onClick={handleToggle}
          disabled={loading}
        >
          {subscribed ? "Unsubscribe" : "Subscribe"}
        </Button>
      </CardFooter>
    </Card>
  )
}
