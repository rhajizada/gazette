import { FeedPreview } from "@/components/FeedPreview";
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
import { useEffect, useState } from "react";
import { toast } from "sonner";
import type { GithubComRhajizadaGazetteInternalServiceFeed as FeedModel } from "../api/data-contracts";
import { useAuth } from "../context/AuthContext";

const PAGE_SIZE = 15;
const labelMap: Record<"title" | "last_updated_at", string> = {
  title: "title",
  last_updated_at: "last updated",
};

type SortKey = "title" | "last_updated_at";

export function SubscribedFeeds() {
  const { api, logout } = useAuth();
  const [allFeeds, setAllFeeds] = useState<FeedModel[]>([]);
  const [loading, setLoading] = useState(true);

  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [sortKey, setSortKey] = useState<SortKey>("title");
  const [sortAsc, setSortAsc] = useState(true);

  useEffect(() => {
    (async () => {
      setLoading(true);
      try {
        const pageSize = 100;
        let offset = 0;
        let total = Infinity;
        const acc: FeedModel[] = [];
        while (offset < total) {
          const res = await api.feedsList(
            { subscribedOnly: true, limit: pageSize, offset },
            { secure: true, format: "json" },
          );
          const chunk = res.data.feeds ?? [];
          acc.push(...chunk);
          total = res.data.total_count ?? chunk.length;
          offset += chunk.length;
          if (!chunk.length) break;
        }
        setAllFeeds(acc);
      } catch (err: any) {
        if (err.error === "Unauthorized") logout();
        else {
          const message = await err.text();
          toast.error(message || "Failed to load subscribed feeds");
        }
      } finally {
        setLoading(false);
      }
    })();
  }, [api, logout]);

  const processed = allFeeds
    .filter(
      (f) =>
        f.title?.toLowerCase().includes(search.toLowerCase()) ||
        f.description?.toLowerCase().includes(search.toLowerCase()),
    )
    .sort((a, b) => {
      let cmp: number;
      if (sortKey === "title") {
        cmp = (a.title ?? "").localeCompare(b.title ?? "");
      } else {
        const at = a.last_updated_at
          ? new Date(a.last_updated_at).getTime()
          : 0;
        const bt = b.last_updated_at
          ? new Date(b.last_updated_at).getTime()
          : 0;
        cmp = at - bt;
      }
      return sortAsc ? cmp : -cmp;
    });

  const totalPages = Math.ceil(processed.length / PAGE_SIZE);
  const pageItems = processed.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE);

  // Build pagination array
  const pages: (number | "ellipsis")[] = [];
  if (totalPages > 0) {
    pages.push(1);
    if (page > 3) pages.push("ellipsis");
    const start = Math.max(2, page - 2);
    const end = Math.min(totalPages - 1, page + 2);
    for (let p = start; p <= end; p++) pages.push(p);
    if (page < totalPages - 2) pages.push("ellipsis");
    if (totalPages > 1) pages.push(totalPages);
  }

  return (
    <div className="flex flex-col space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <Input
          placeholder="Search feeds..."
          value={search}
          onChange={(e) => {
            setSearch(e.currentTarget.value);
            setPage(1);
          }}
          className="flex-1 min-w-0"
        />
        <Button
          variant="outline"
          size="sm"
          onClick={() => {
            if (sortKey === "title" && sortAsc) setSortAsc(false);
            else if (sortKey === "title" && !sortAsc) {
              setSortKey("last_updated_at");
              setSortAsc(true);
            } else if (sortKey === "last_updated_at" && sortAsc)
              setSortAsc(false);
            else {
              setSortKey("title");
              setSortAsc(true);
            }
            setPage(1);
          }}
        >
          sort by {labelMap[sortKey]} {sortAsc ? "↑" : "↓"}
        </Button>
      </div>

      {loading ? (
        <div className="flex justify-center py-10">
          <Spinner />
        </div>
      ) : pageItems.length === 0 ? (
        <p className="text-center text-gray-500">No subscribed feeds.</p>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {pageItems.map((f) => (
            <FeedPreview key={f.id} feed={f} />
          ))}
        </div>
      )}

      {totalPages > 1 && (
        <div className="flex justify-center pt-4">
          <Pagination>
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  aria-disabled={page === 1}
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
                    >
                      {p}
                    </PaginationLink>
                  </PaginationItem>
                ),
              )}
              <PaginationItem>
                <PaginationNext
                  onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                  aria-disabled={page === totalPages}
                />
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        </div>
      )}
    </div>
  );
}
