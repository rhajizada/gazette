import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { useAuth } from "../context/AuthContext"
import type { GithubComRhajizadaGazetteInternalServiceItem as ItemModel } from "../api/data-contracts"
import { Card, CardContent, CardFooter } from "@/components/ui/card"
import { Heart } from "lucide-react"
import { toast } from "sonner"

interface ItemPreviewProps {
  item: ItemModel
}

export function ItemPreview({ item }: ItemPreviewProps) {
  const { api } = useAuth()
  const navigate = useNavigate()
  const [liked, setLiked] = useState<boolean>(item.liked ?? false)
  const [loading, setLoading] = useState<boolean>(false)

  const toggleLike = async () => {
    if (loading) return
    setLoading(true)
    try {
      if (liked) {
        await api.itemsLikeDelete(item.id!, { format: "json" })
        setLiked(false)
      } else {
        await api.itemsLikeCreate(item.id!, { format: "json" })
        setLiked(true)
      }
    } catch (err) {
      console.error("Failed to toggle like:", err)
      toast.error("Could not update like status")
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card className="group overflow-hidden rounded-2xl shadow hover:shadow-lg transition-shadow duration-300 cursor-pointer">
      <CardContent className="flex break-words">
        {/* Thumbnail on left */}
        {item.image?.url && (
          <img
            src={item.image.url}
            alt={item.image.title ?? item.title}
            className="w-24 h-24 object-cover rounded-lg mr-4 flex-shrink-0"
            onClick={() => navigate(`/items/${item.id}`)}
          />
        )}
        <div onClick={() => navigate(`/items/${item.id}`)} className="flex-1 break-words">
          <h3 className="text-lg font-semibold text-gray-900 hover:text-primary transition-colors">
            {item.title}
          </h3>
          {item.description && (
            <p className="text-gray-700 line-clamp-3 mt-1 break-words">
              {item.description}
            </p>
          )}
          {item.authors && item.authors.length > 0 && (
            <div className="flex flex-wrap gap-1 mt-2">
              {item.authors.map((a, i) => (
                <span
                  key={i}
                  className="px-2 py-0.5 bg-gray-100 text-xs rounded-full"
                >
                  {a.name || a.email}
                </span>
              ))}
            </div>
          )}
        </div>
      </CardContent>

      <CardFooter className="px-4 py-2 flex items-center justify-between">
        <button
          onClick={toggleLike}
          disabled={loading}
          aria-label={liked ? "Unlike item" : "Like item"}
          className={`p-2 rounded-full transition-colors ${liked ? "text-red-500 hover:bg-red-100" : "text-gray-500 hover:bg-gray-100"
            } ${loading ? "opacity-50 cursor-not-allowed" : ""}`}
        >
          <Heart
            fill={liked ? "currentColor" : "none"}
            strokeWidth={2}
            className="w-5 h-5"
          />
        </button>
        {item.published_parsed && (
          <span className="text-xs text-gray-500">
            {new Date(item.published_parsed).toLocaleDateString()}
          </span>
        )}
      </CardFooter>
    </Card>
  )
}

