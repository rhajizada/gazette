import { useEffect, useState, useCallback } from "react"
import { useAuth } from "../context/AuthContext"
import type { GithubComRhajizadaGazetteInternalServiceFeed as FeedModel } from "../api/data-contracts"
import { Spinner } from "@/components/ui/spinner"
import { FeedPreview } from "@/components/FeedPreview"
import { Navbar } from "@/components/Navbar"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination"

export default function Feeds() {
  const { api, logout } = useAuth()
  const PAGE_SIZE = 15

  const [feeds, setFeeds] = useState<FeedModel[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [page, setPage] = useState(1)

  const [newUrl, setNewUrl] = useState("")
  const [creating, setCreating] = useState(false)

  const fetchFeeds = useCallback(() => {
    setLoading(true)
    api
      .feedsList(
        { subscribedOnly: false, limit: PAGE_SIZE, offset: (page - 1) * PAGE_SIZE },
        { format: "json" }
      )
      .then((res) => {
        setFeeds(res.data.feeds ?? [])
        setTotal(res.data.total_count ?? 0)
      })
      .catch((err) => {
        if ((err as any).error === "Unauthorized") {
          logout()
        } else {
          console.error(err)
          setError("Failed to load feeds")
        }
      })
      .finally(() => setLoading(false))
  }, [api, logout, page])

  useEffect(() => {
    fetchFeeds()
  }, [fetchFeeds])

  const handleCreate = async () => {
    if (!newUrl) return
    setCreating(true)
    try {
      await api.feedsCreate({ feed_url: newUrl }, { format: "json" })
      setNewUrl("")
      setPage(1)
    } catch (raw: unknown) {
      console.error("Failed to create feed:", raw)
      const e = raw as any
      let description = await e.text?.().catch(() => "")
      if (!description) description = e.error?.toString() || e.message || "Unknown error"
      toast.error("Failed to create feed", { description })
    } finally {
      setCreating(false)
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

  if (error) {
    return (
      <>
        <Navbar />
        <div className="p-4 text-red-600">{error}</div>
      </>
    )
  }

  const totalPages = Math.ceil(total / PAGE_SIZE)
  // build page links
  const pages: (number | "ellipsis")[] = []
  if (totalPages <= 7) {
    for (let i = 1; i <= totalPages; i++) pages.push(i)
  } else if (page <= 4) {
    pages.push(1, 2, 3, 4, 5, "ellipsis", totalPages)
  } else if (page > totalPages - 4) {
    pages.push(1, "ellipsis", totalPages - 4, totalPages - 3, totalPages - 2, totalPages - 1, totalPages)
  } else {
    pages.push(1, "ellipsis", page - 1, page, page + 1, "ellipsis", totalPages)
  }

  return (
    <>
      <Navbar />
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="mb-6 flex space-x-2">
          <Input
            value={newUrl}
            onChange={(e) => setNewUrl(e.currentTarget.value)}
            placeholder="Enter feed URL"
            className="flex-1"
          />
          <Button onClick={handleCreate} disabled={creating || !newUrl}>
            {creating ? "Adding..." : "Import"}
          </Button>
        </div>

        {feeds.length === 0 ? (
          <p className="text-center">No feeds yet—go add some!</p>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {feeds.map((f) => (
              <FeedPreview key={f.id} feed={f} />
            ))}
          </div>
        )}

        {totalPages > 1 && (
          <div className="mt-8 flex justify-center">
            <Pagination>
              <PaginationContent>
                <PaginationItem>
                  <PaginationPrevious
                    onClick={() => setPage(page - 1)}
                    aria-disabled={page === 1}
                    tabIndex={page === 1 ? -1 : undefined}
                    className={page === 1 ? "pointer-events-none opacity-50" : undefined}
                    href="#"
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
                        onClick={() => setPage(p as number)}
                        isActive={p === page}
                        href="#"
                      >
                        {p}
                      </PaginationLink>
                    </PaginationItem>
                  )
                )}
                <PaginationItem>
                  <PaginationNext
                    onClick={() => setPage(page + 1)}
                    aria-disabled={page === totalPages}
                    tabIndex={page === totalPages ? -1 : undefined}
                    className={page === totalPages ? "pointer-events-none opacity-50" : undefined}
                    href="#"
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

