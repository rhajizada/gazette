import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardFooter } from "@/components/ui/card"
import { Star } from "lucide-react"
import { useState, useEffect } from "react"
import { useNavigate } from "react-router-dom"
import type { GithubComRhajizadaGazetteInternalServiceFeed as FeedType } from "../api/data-contracts"
import { useAuth } from "../context/AuthContext"

const MIN_WIDTH = 100
const MIN_HEIGHT = 100

interface FeedProps {
  feed: FeedType
}

export function FeedPreview({ feed }: FeedProps) {
  const navigate = useNavigate()
  const { api, logout } = useAuth()
  const [subscribed, setSubscribed] = useState(feed.subscribed ?? false)
  const [loading, setLoading] = useState(false)
  const [showImage, setShowImage] = useState(false)

  useEffect(() => {
    if (!feed?.image?.url) return
    const img = new Image()
    img.src = feed.image.url
    img.onload = () => {
      if (img.naturalWidth > MIN_WIDTH && img.naturalHeight > MIN_HEIGHT) {
        setShowImage(true)
      }
    }
  }, [feed?.image?.url])


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
      className="group flex flex-col h-full overflow-hidden rounded-2xl shadow-lg hover:shadow-2xl transition-shadow duration-300 cursor-pointer"
    >
      <CardContent className="flex-1 relative pt-4 px-4 pb-2">
        <>
          {showImage && (
            <img
              src={feed.image.url}
              alt={feed.image.title ?? feed.title}
              className="
            float-left
            mr-6
            mb-6
            rounded-lg
            max-w-[33.333%]
            h-auto
          "
            />
          )}
        </>

        <h3 className="text-lg font-bold text-gray-900 hover:text-primary transition-colors">
          {feed.title}
        </h3>

        {feed.description && (
          <p className="text-gray-700 leading-relaxed mt-2">
            {feed.description}
          </p>
        )}

        <div className="flex flex-wrap gap-2 mt-4 clear-left">
          {feed.categories?.slice(0, 3).map((cat, _) => (
            <Badge>
              {cat}
            </Badge>
          ))}
        </div>
        {feed.authors && feed.authors.length > 0 && (
          <p className="text-sm text-gray-500 mt-2">
            <strong>Author{feed.authors.length > 1 ? "s" : ""}:</strong> {feed.authors.map(a => a.name || a.email).join(", ")}
          </p>
        )}
        {feed.language && (
          <p className="text-sm text-gray-500 mt-1">
            <strong>Language:</strong> {feed.language}
          </p>
        )}
        <p className="text-sm text-gray-500 mt-1">
          <strong>Last updated:</strong> {feed.updated_parsed && new Date(feed.updated_parsed).toLocaleDateString()}
        </p>
      </CardContent>

      <CardFooter className="px-4 py-2 flex items-center justify-end">
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

