import { FeedPreview } from "@/components/FeedPreview";
import { Footer } from "@/components/Footer";
import { Navbar } from "@/components/Navbar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { Spinner } from "@/components/ui/spinner";
import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";
import type { GithubComRhajizadaGazetteInternalServiceFeed as FeedModel } from "../api/data-contracts";
import { useAuth } from "../context/AuthContext";

export default function Feeds() {
  const { api, logout } = useAuth();
  const PAGE_SIZE = 15;

  const [feeds, setFeeds] = useState<FeedModel[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);

  const [newUrl, setNewUrl] = useState("");
  const [creating, setCreating] = useState(false);

  const [search, setSearch] = useState("");
  const [sortKey, setSortKey] = useState<"title" | "last_updated_at">("title");
  const [sortAsc, setSortAsc] = useState(true);

  // Fetch paginated feeds
  const fetchFeeds = useCallback(() => {
    setLoading(true);
    api
      .feedsList(
        { subscribedOnly: false, limit: PAGE_SIZE, offset: (page - 1) * PAGE_SIZE },
        { secure: true, format: "json" }
      )
      .then((res) => {
        setFeeds(res.data.feeds ?? []);
        setTotal(res.data.total_count ?? 0);
      })
      .catch((err) => {
        if ((err as any).error === "Unauthorized") {
          logout();
        } else {
          console.error(err);
          setError("Failed to load feeds");
        }
      })
      .finally(() => setLoading(false));
  }, [api, logout, page]);

  useEffect(() => {
    fetchFeeds();
  }, [fetchFeeds]);

  const handleCreate = async () => {
    if (!newUrl.trim()) return;
    setCreating(true);
    try {
      await api.feedsCreate({ feed_url: newUrl.trim() }, { secure: true, format: "json" });
      setNewUrl("");
      setPage(1);
      fetchFeeds();
      toast.success("Feed imported");
    } catch (raw: unknown) {
      console.error("Failed to create feed:", raw);
      const e = raw as any;
      let description = await e.text?.().catch(() => "");
      if (!description) description = e.error?.toString() || e.message || "Unknown error";
      toast.error("Failed to create feed", { description });
    } finally {
      setCreating(false);
    }
  };

  // Filter & sort
  const filtered = feeds
    .filter(f => f.title?.toLowerCase().includes(search.toLowerCase()))
    .sort((a, b) => {
      const aVal = (a[sortKey] || "").toString();
      const bVal = (b[sortKey] || "").toString();
      if (aVal < bVal) return sortAsc ? -1 : 1;
      if (aVal > bVal) return sortAsc ? 1 : -1;
      return 0;
    });

  const totalPages = Math.ceil(total / PAGE_SIZE);
  const pages: (number | "ellipsis")[] = [];
  if (totalPages <= 7) {
    for (let i = 1; i <= totalPages; i++) pages.push(i);
  } else if (page <= 4) pages.push(1, 2, 3, 4, 5, "ellipsis", totalPages);
  else if (page > totalPages - 4)
    pages.push(1, "ellipsis", totalPages - 4, totalPages - 3, totalPages - 2, totalPages - 1, totalPages);
  else pages.push(1, "ellipsis", page - 1, page, page + 1, "ellipsis", totalPages);

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <main className="flex-grow bg-gray-50 p-6">
        <div className="mb-4 flex justify-between items-center">
          {error && (
            <div className="mb-4 rounded border border-red-300 bg-red-50 px-4 py-2 text-red-800">
              {error}
            </div>
          )}
          <div className="flex items-center gap-2">
            <Input
              placeholder="Search feeds..."
              value={search}
              onChange={e => setSearch(e.currentTarget.value)}
              className="w-48"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                setSortKey(k => (k === "title" ? "last_updated_at" : "title"));
                setSortAsc(a => !a);
              }}
            >
              Sort by {sortKey === "title" ? "title" : "last updated"} {sortAsc ? '↑' : '↓'}
            </Button>
          </div>
          <div className="flex items-center gap-2">
            <Input
              placeholder="Enter feed URL"
              value={newUrl}
              onChange={e => setNewUrl(e.currentTarget.value)}
              className="w-64"
            />
            <Button onClick={handleCreate} disabled={creating || !newUrl.trim()}>
              {creating ? <Spinner size="sm" /> : "Import"}
            </Button>
          </div>
        </div>

        {loading ? (
          <div className="flex justify-center py-10"><Spinner /></div>
        ) : filtered.length === 0 ? (
          <p className="text-center">No feeds found.</p>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {filtered.map(f => (
              <FeedPreview key={f.id} feed={f} />
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
                {pages.map((p, idx) =>
                  p === "ellipsis" ? (
                    <PaginationItem key={idx}>
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
      </main>
      <Footer />
    </div>
  );
}

