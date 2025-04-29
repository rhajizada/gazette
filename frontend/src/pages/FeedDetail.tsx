
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
  const [itemTotal, setItemTotal] = useState(0)
  const [itemPage, setItemPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const [itemsLoading, setItemsLoading] = useState(false)
  const [notFound, setNotFound] = useState(false)

  // Fetch feed details
  useEffect(() => {
    if (!feedID) return
    setLoading(true)
    api
      .feedsDetail(feedID, { secure: true, format: "json" })
      .then((res) => setFeed(res.data))
      .catch((err) => {
        if (err.status === 404 || err.error === "Not Found") {
          setNotFound(true)
        } else if (err.error === "Unauthorized") {
          logout()
        } else {
          console.error(err)
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
        { limit: PAGE_SIZE, offset: (itemPage - 1) * PAGE_SIZE },
        { secure: true, format: "json" }
      )
      .then((res) => {
        setItems(res.data.items ?? [])
        setItemTotal(res.data.total_count ?? 0)
      })
      .catch((err) => {
        if (err.error === "Unauthorized") logout()
        else console.error(err)
      })
      .finally(() => setItemsLoading(false))
  }, [api, feedID, itemPage, logout])

  useEffect(() => { fetchItems() }, [fetchItems])

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <Spinner />
      </div>
    )
  }

  if (notFound || !feed) {
    return <Navigate to="*" replace />
  }

  // Pagination logic
  const totalPages = Math.ceil(itemTotal / PAGE_SIZE)
  const pages: (number | "ellipsis")[] = []
  if (totalPages <= 7) {
    for (let i = 1; i <= totalPages; i++) pages.push(i)
  } else if (itemPage <= 4) {
    pages.push(1, 2, 3, 4, 5, "ellipsis", totalPages)
  } else if (itemPage > totalPages - 4) {
    pages.push(
      1,
      "ellipsis",
      totalPages - 4,
      totalPages - 3,
      totalPages - 2,
      totalPages - 1,
      totalPages
    )
  } else {
    pages.push(
      1,
      "ellipsis",
      itemPage - 1,
      itemPage,
      itemPage + 1,
      "ellipsis",
      totalPages
    )
  }

  return (
    <>
      <Navbar />
      <div className="p-6 space-y-6">
        <header className="space-y-1">
          <h1 className="text-2xl font-semibold">{feed.title}</h1>
          {feed.updated_parsed && (
            <p className="text-sm text-gray-500">
              Last updated: {new Date(feed.updated_parsed).toLocaleString()}
            </p>
          )}
        </header>

        <section>
          {itemsLoading ? (
            <Spinner />
          ) : items.length === 0 ? (
            <p className="text-center text-gray-500">No items in this feed.</p>
          ) : (
            // Updated to four columns
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6">
              {items.map((item) => (
                <ItemPreview key={item.id} item={item} />
              ))}
            </div>
          )}
        </section>

        {/* Pagination controls */}
        {totalPages > 1 && (
          <div className="flex justify-center">
            <Pagination>
              <PaginationContent>
                <PaginationItem>
                  <PaginationPrevious
                    onClick={() => setItemPage(itemPage - 1)}
                    aria-disabled={itemPage === 1}
                    tabIndex={itemPage === 1 ? -1 : undefined}
                    className={itemPage === 1 ? "pointer-events-none opacity-50" : undefined}
                  />
                </PaginationItem>
                {pages.map((p, i) =>
                  p === "ellipsis" ? (
                    <PaginationItem key={`e${i}`}>
                      <PaginationEllipsis />
                    </PaginationItem>
                  ) : (
                    <PaginationItem key={p}>
                      <PaginationLink
                        onClick={() => setItemPage(p as number)}
                        isActive={p === itemPage}
                      >
                        {p}
                      </PaginationLink>
                    </PaginationItem>
                  )
                )}
                <PaginationItem>
                  <PaginationNext
                    onClick={() => setItemPage(itemPage + 1)}
                    aria-disabled={itemPage === totalPages}
                    tabIndex={itemPage === totalPages ? -1 : undefined}
                    className={itemPage === totalPages ? "pointer-events-none opacity-50" : undefined}
                  />
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          </div>
        )}
      </div>
    </>
  )
}

