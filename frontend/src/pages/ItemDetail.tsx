import { useEffect, useState } from "react"
import { useParams, Navigate } from "react-router-dom"
import { Spinner } from "@/components/ui/spinner"
import { useAuth } from "../context/AuthContext"
import { Navbar } from "@/components/Navbar"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"
import { Heart } from "lucide-react"
import type { GithubComRhajizadaGazetteInternalServiceItem as ItemModel } from "../api/data-contracts"

export default function ItemDetail() {
  const { itemID } = useParams<{ itemID: string }>()
  const { api, logout } = useAuth()

  const [item, setItem] = useState<ItemModel | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [notFound, setNotFound] = useState(false)

  // Like state
  const [liked, setLiked] = useState<boolean>(false)
  const [likeLoading, setLikeLoading] = useState<boolean>(false)

  useEffect(() => {
    if (!itemID) return
    setLoading(true)
    api
      .itemsDetail(itemID, { secure: true, format: "json" })
      .then((res) => {
        setItem(res.data)
        setLiked(res.data.liked ?? false)
      })
      .catch((err) => {
        if (err.status === 404 || err.error === "Not Found") setNotFound(true)
        else if (err.error === "Unauthorized") logout()
        else {
          console.error(err)
          setError("Failed to load item")
        }
      })
      .finally(() => setLoading(false))
  }, [api, itemID, logout])

  const toggleLike = async () => {
    if (likeLoading || !item) return
    setLikeLoading(true)
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
      setLikeLoading(false)
    }
  }

  if (loading) {
    return (
      <>
        <Navbar />
        <div className="flex items-center justify-center min-h-screen">
          <Spinner />
        </div>
      </>
    )
  }

  if (notFound) return <Navigate to="*" replace />
  if (error || !item) return (
    <>
      <Navbar />
      <div className="p-6 text-red-600">{error || "Item not available"}</div>
    </>
  )

  return (
    <>
      <Navbar />
      {/* Main layout grid */}
      <div className="mx-auto max-w-5xl px-6 py-8 grid grid-cols-1 lg:grid-cols-[3fr_1fr] gap-8">
        {/* Main content */}
        <div>
          <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
            {item.title}
          </h1>

          {item.authors && item.authors.length > 0 && (
            <div className="flex flex-wrap gap-2 mt-4">
              {item.authors.map((a, i) => (
                <Badge key={i}>{a.name || a.email}</Badge>
              ))}
            </div>
          )}

          <p className="leading-7 [&:not(:first-child)]:mt-6 text-sm text-gray-500">
            Published: {item.published_parsed && new Date(item.published_parsed).toLocaleString()}
          </p>
          <p className="leading-7 [&:not(:first-child)]:mt-6 text-sm text-gray-500">
            Updated: {item.updated_parsed && new Date(item.updated_parsed).toLocaleString()}
          </p>

          {item.image?.url && (
            <img
              src={item.image.url}
              alt={item.image.title ?? item.title}
              className="float-left mr-6 mb-6 w-1/3 rounded-lg"
            />
          )}

          {item.description && (
            <div
              className="leading-7 [&:not(:first-child)]:mt-6"
              dangerouslySetInnerHTML={{ __html: item.description }}
            />
          )}

          {item.content && (
            <div
              className="leading-7 [&:not(:first-child)]:mt-6"
              dangerouslySetInnerHTML={{ __html: item.content }}
            />
          )}

          {item.categories && item.categories.length > 0 && (
            <ul className="my-6 ml-6 list-disc [&>li]:mt-2">
              <h2 className="mt-8 scroll-m-20 text-2xl font-semibold tracking-tight">Categories</h2>
              <li>{item.categories.join(", ")}</li>
            </ul>
          )}

          {item.link && (
            <blockquote className="mt-6 border-l-2 pl-6 italic">
              <Button asChild variant="link">
                <a href={item.link} target="_blank" rel="noopener noreferrer">
                  Read Original Source
                </a>
              </Button>
            </blockquote>
          )}

          {/* Like button at bottom */}
          <div className="mt-8">
            <button
              onClick={toggleLike}
              disabled={likeLoading}
              aria-label={liked ? "Unlike item" : "Like item"}
              className={`p-2 rounded-full transition-colors ${liked ? "text-red-500 hover:bg-red-100" : "text-gray-500 hover:bg-gray-100"
                } ${likeLoading ? "opacity-50 cursor-not-allowed" : ""}`}
            >
              <Heart fill={liked ? "currentColor" : "none"} strokeWidth={2} className="w-6 h-6" />
            </button>
          </div>
        </div>

        {item.enclosures && item.enclosures.length > 0 && (
          <aside>
            <h2 className="mt-0 scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight">
              Attachments
            </h2>
            <div className="mt-4 flex flex-col gap-4">
              {item.enclosures.map((enc: any, i: number) => {
                if (enc.type?.startsWith("audio/")) {
                  return <audio key={i} controls src={enc.url} className="w-full" />
                }
                if (enc.type?.startsWith("video/")) {
                  return <video key={i} controls src={enc.url} className="w-full" />
                }
                return (
                  <a
                    key={i}
                    href={enc.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="font-medium text-primary underline underline-offset-4"
                  >
                    Download ({enc.type || 'file'})
                  </a>
                )
              })}
            </div>
          </aside>
        )}
      </div>
    </>
  )
}

