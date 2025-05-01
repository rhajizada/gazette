import { Footer } from "@/components/Footer"
import { ItemPreview } from "@/components/ItemPreview"
import { Navbar } from "@/components/Navbar"
import { Button } from "@/components/ui/button"
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination"
import { Spinner } from "@/components/ui/spinner"
import { Star } from "lucide-react"
import { useCallback, useEffect, useState } from "react"
import { Navigate, useParams } from "react-router-dom"
import { toast } from "sonner"
import type {
  GithubComRhajizadaGazetteInternalServiceFeed as FeedType,
  GithubComRhajizadaGazetteInternalServiceItem as ItemModel,
} from "../api/data-contracts"
import { useAuth } from "../context/AuthContext"

export default function FeedDetail() {
  const { feedID } = useParams<{ feedID: string }>()
  const { api, logout } = useAuth()
  const PAGE_SIZE = 10

  const [feed, setFeed] = useState<FeedType | null>(null)
  const [items, setItems] = useState<ItemModel[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const [itemsLoading, setItemsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [notFound, setNotFound] = useState(false)
  const [subscribed, setSubscribed] = useState(false)
  const [subLoading, setSubLoading] = useState(false)

  // Fetch feed details
  useEffect(() => {
    if (!feedID) return
    setLoading(true)
    api
      .feedsDetail(feedID, { secure: true, format: "json" })
      .then((res) => {
        setFeed(res.data)
        setSubscribed(res.data.subscribed ?? false)
      })
      .catch((err) => {
        if (err.status === 404 || err.status === 400) {
          setNotFound(true)
        } else {
          toast.error("Failed to load feed")
          setError("Failed to load feed")
          if (err.error === "Unauthorized") logout()
        }
      })
      .finally(() => setLoading(false))
  }, [api, feedID, logout])

  // Fetch paginated items
  const fetchItems = useCallback(() => {
    if (!feedID) return
    setItemsLoading(true)
    api
      .feedsItemsList(
        feedID,
        { limit: PAGE_SIZE, offset: (page - 1) * PAGE_SIZE },
        { secure: true, format: "json" }
      )
      .then((res) => {
        setItems(res.data.items ?? [])
        setTotal(res.data.total_count ?? 0)
      })
      .catch((err) => {
        if (err.error === "Unauthorized") logout()
        else console.error(err)
      })
      .finally(() => setItemsLoading(false))
  }, [api, feedID, page, logout])

  useEffect(() => {
    fetchItems()
  }, [fetchItems])

  const toggleSub = async () => {
    if (!feed || subLoading) return
    setSubLoading(true)
    try {
      if (subscribed) {
        await api.feedsSubscribeDelete(feed.id!, { format: "json" })
        setSubscribed(false)
      } else {
        await api.feedsSubscribeUpdate(feed.id!, { format: "json" })
        setSubscribed(true)
      }
    } catch (e) {
      console.error(e)
    } finally {
      setSubLoading(false)
    }
  }

  if (loading) return <Loader />
  if (notFound || !feed) return <Navigate to="/*" />
  if (error || !feed) return (
    <>
      <Navbar />
      <div className="h-2 ju text-red-600">{error || "Feed not available"}</div>
    </>
  )

  const totalPages = Math.ceil(total / PAGE_SIZE)
  const pages: (number | "ellipsis")[] = []
  if (totalPages <= 7) {
    for (let i = 1; i <= totalPages; i++) pages.push(i)
  } else if (page <= 4) pages.push(1, 2, 3, 4, 5, "ellipsis", totalPages)
  else if (page > totalPages - 4)
    pages.push(
      1,
      "ellipsis",
      totalPages - 4,
      totalPages - 3,
      totalPages - 2,
      totalPages - 1,
      totalPages
    )
  else pages.push(1, "ellipsis", page - 1, page, page + 1, "ellipsis", totalPages)

  return (
    <>
      <Navbar />

      <div className="mx-auto max-w-4xl px-6 py-8">
        <h1 className="text-4xl font-bold text-black">{feed.title}</h1>
        {feed.description && (
          <p className="text-lg text-gray-700 mt-2">{feed.description}</p>
        )}
        {feed.updated_parsed && (
          <p className="text-sm text-gray-500 mt-1">Last updated: {new Date(feed.updated_parsed).toLocaleDateString()}</p>
        )}
        <br />
        <Button
          size="sm"
          onClick={toggleSub}
          disabled={subLoading}
          className="inline-flex items-center bg-white text-gray-800 px-3 py-1 shadow hover:bg-gray-100 transition disabled:opacity-50"
        >
          <Star
            fill={subscribed ? "currentColor" : "none"}
            className={`w-5 h-5 mr-2 ${subscribed ? "text-yellow-400" : "text-gray-800"
              }`}
          />
          {subscribed ? "Unsubscribe" : "Subscribe"}
        </Button>
      </div>
      <div className="mx-auto max-w-5xl px-6">
        {itemsLoading ? (
          <Loader />
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {items.map((item) => (
              <ItemPreview key={item.id} item={item} />
            ))}
          </div>
        )}

        {totalPages > 1 && (
          <div className="mt-8 flex justify-center">
            <Pagination>
              <PaginationContent>
                <PaginationItem>
                  <PaginationPrevious onClick={() => setPage(page - 1)} aria-disabled={page === 1} />
                </PaginationItem>
                {pages.map((p, i) =>
                  p === "ellipsis" ? (
                    <PaginationItem key={i}>
                      <PaginationEllipsis />
                    </PaginationItem>
                  ) : (
                    <PaginationItem key={p}>
                      <PaginationLink onClick={() => setPage(p as number)} isActive={p === page}>
                        {p}
                      </PaginationLink>
                    </PaginationItem>
                  )
                )}
                <PaginationItem>
                  <PaginationNext onClick={() => setPage(page + 1)} aria-disabled={page === totalPages} />
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          </div>
        )}
      </div>
      <Footer />
    </>
  )
}

function Loader() {
  return (
    <div className="flex items-center justify-center min-h-[300px]">
      <Spinner />
    </div>
  )
}

