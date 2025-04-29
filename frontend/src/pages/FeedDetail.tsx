import { useEffect, useState, useCallback } from "react"
import { useParams, Navigate } from "react-router-dom"
import { Spinner } from "@/components/ui/spinner"
import { ItemPreview } from "@/components/ItemPreview"
import { useAuth } from "../context/AuthContext"
import { Navbar } from "@/components/Navbar"
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationPrevious,
  PaginationEllipsis,
  PaginationLink,
  PaginationNext,
} from "@/components/ui/pagination"
import { Card, CardContent } from "@/components/ui/card"
import { Star } from "lucide-react"

import type {
  GithubComRhajizadaGazetteInternalServiceFeed as FeedType,
  GithubComRhajizadaGazetteInternalServiceItem as ItemModel,
} from "../api/data-contracts"

export default function FeedDetail() {
  const { feedID } = useParams<{ feedID: string }>()
  const { api, logout } = useAuth()
  const PAGE_SIZE = 24

  const [feed, setFeed] = useState<FeedType | null>(null)
  const [items, setItems] = useState<ItemModel[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const [itemsLoading, setItemsLoading] = useState(false)
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
        if (err.status === 404 || err.error === "Not Found") setNotFound(true)
        else if (err.error === "Unauthorized") logout()
        else console.error(err)
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

  useEffect(() => { fetchItems() }, [fetchItems])

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
  if (notFound || !feed) return <Navigate to="*" replace />

  // Pagination UI
  const totalPages = Math.ceil(total / PAGE_SIZE)
  const pages: (number | "ellipsis")[] = []
  if (totalPages <= 7) {
    for (let i = 1; i <= totalPages; i++) pages.push(i)
  } else if (page <= 4) pages.push(1, 2, 3, 4, 5, "ellipsis", totalPages)
  else if (page > totalPages - 4) pages.push(1, "ellipsis", totalPages - 4, totalPages - 3, totalPages - 2, totalPages - 1, totalPages)
  else pages.push(1, "ellipsis", page - 1, page, page + 1, "ellipsis", totalPages)

  return (
    <>
      <Navbar />
      <Card className="relative overflow-hidden rounded-2xl shadow-lg mx-auto max-w-4xl mb-8">
        {feed.image?.url && (
          <div className="relative h-64 w-full">
            <img
              src={feed.image.url}
              alt={feed.title}
              className="absolute inset-0 w-full h-full object-cover"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-black via-transparent opacity-60" />
          </div>
        )}
        <CardContent className="absolute bottom-4 left-4 text-white">
          <h1 className="text-3xl font-bold truncate">{feed.title}</h1>
          <p className="mt-2 line-clamp-2 text-sm">{feed.description}</p>

          <button
            onClick={toggleSub}
            disabled={subLoading}
            className="mt-4 inline-flex items-center bg-white text-gray-800 px-3 py-1 rounded-full shadow hover:bg-gray-100 transition disabled:opacity-50"
          >
            <Star
              fill={subscribed ? "currentColor" : "none"}
              className={`w-5 h-5 mr-2 ${subscribed ? "text-yellow-400" : "text-gray-800"}`}
            />
            {subscribed ? "Unsubscribe" : "Subscribe"}
          </button>

        </CardContent>
      </Card>

      <div className="mx-auto max-w-5xl px-6">
        {itemsLoading ? <Loader /> : (
          <div className="grid grid-cols-1 sm:grid-cols-1 md:grid-cols-2 gap-6">
            {items.map(item => <ItemPreview key={item.id} item={item} />)}
          </div>
        )}

        {totalPages > 1 && (
          <div className="mt-8 flex justify-center">
            <Pagination>
              <PaginationContent>
                <PaginationItem>
                  <PaginationPrevious onClick={() => setPage(page - 1)} aria-disabled={page === 1} />
                </PaginationItem>
                {pages.map((p, i) => p === "ellipsis"
                  ? <PaginationItem key={i}><PaginationEllipsis /></PaginationItem>
                  : <PaginationItem key={p}><PaginationLink onClick={() => setPage(p as number)} isActive={p === page}>{p}</PaginationLink></PaginationItem>
                )}
                <PaginationItem>
                  <PaginationNext onClick={() => setPage(page + 1)} aria-disabled={page === totalPages} />
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          </div>
        )}
      </div>
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


