import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { useAuth } from "../context/AuthContext"
import type { GithubComRhajizadaGazetteInternalServiceFeed as FeedType } from "../api/data-contracts"
import { Card, CardContent, CardFooter } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Star } from "lucide-react"

interface FeedProps {
  feed: FeedType
}

export function FeedPreview({ feed }: FeedProps) {
  const navigate = useNavigate()
  const { api, logout } = useAuth()
  const [subscribed, setSubscribed] = useState(feed.subscribed ?? false)
  const [loading, setLoading] = useState(false)

  const handleToggle = async (e: React.MouseEvent) => {
    e.stopPropagation()
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
      if (err.error === "Unauthorized") logout()
      else console.error("Subscription error:", err)
    } finally {
      setLoading(false)
    }
  }

  const handleCardClick = () => {
    if (feed.id) navigate(`/feeds/${feed.id}`)
  }

  return (
    <Card
      onClick={handleCardClick}
      className="group flex flex-col h-full relative overflow-hidden rounded-2xl shadow-lg hover:shadow-2xl transition-shadow duration-300 cursor-pointer"
    >
      {feed.image?.url && (
        <div className="relative h-48 w-full flex-shrink-0">
          <img
            src={feed.image.url}
            alt={feed.image.title ?? feed.title}
            className="absolute inset-0 h-full w-full object-cover"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-black via-transparent opacity-50" />
          <div className="absolute bottom-4 left-4">
            <h3 className="text-lg font-bold leading-tight text-white">{feed.title}</h3>
          </div>
        </div>
      )}

      <CardContent className="flex-1 pt-4 px-4 pb-2">
        {feed.description && (
          <p className="text-gray-700 line-clamp-3 mb-3">{feed.description}</p>
        )}

        <div className="flex flex-wrap gap-2 mb-2">
          {feed.categories?.slice(0, 3).map((cat, _) => (
            <Badge>
              {cat}
            </Badge>
          ))}
        </div>

        {feed.authors && feed.authors.length > 0 && (
          <p className="text-sm text-gray-500 mb-1">
            <strong>Author{feed.authors.length > 1 ? "s" : ""}:</strong>{" "}
            {feed.authors.map((a) => a.name || a.email).join(", ")}
          </p>
        )}

        {feed.language && (
          <p className="text-sm text-gray-500 mb-1">
            <strong>Language:</strong> {feed.language}
          </p>
        )}

        <p className="text-xs text-gray-400">
          Last updated: {feed.updated_parsed && new Date(feed.updated_parsed).toLocaleDateString()}
        </p>
      </CardContent>

      <CardFooter className="px-4 py-2 flex items-center justify-end space-x-4">
        <Button
          size="sm"
          onClick={handleToggle}
          disabled={loading}
          className="inline-flex items-center bg-white text-gray-800 px-3 py-1 rounded-full shadow hover:bg-gray-100 transition disabled:opacity-50"
        >
          <Star
            fill={subscribed ? "currentColor" : "none"}
            className={`w-5 h-5 mr-2 ${subscribed ? "text-yellow-400" : "text-gray-800"
              }`}
          />
          {subscribed ? "Unsubscribe" : "Subscribe"}
        </Button>
      </CardFooter>
    </Card>
  )
}

