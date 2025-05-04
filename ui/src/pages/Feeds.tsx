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
import { PlusCircle } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import type { GithubComRhajizadaGazetteInternalServiceFeed as FeedModel } from "../api/data-contracts";
import { useAuth } from "../context/AuthContext";

export default function Feeds() {
  const { api, logout } = useAuth();
  const PAGE_SIZE = 15;

  const [allFeeds, setAllFeeds] = useState<FeedModel[]>([]);
  const [loadingAll, setLoadingAll] = useState(true);

  const [page, setPage] = useState(1);
  const [newUrl, setNewUrl] = useState("");
  const [creating, setCreating] = useState(false);
  const [search, setSearch] = useState("");

  const [sortKey, setSortKey] = useState<"title" | "updated_parsed">("title");
  const [sortAsc, setSortAsc] = useState(true);

  const labelMap: Record<"title" | "updated_parsed", string> = {
    title: "title",
    updated_parsed: "updated",
  };

  useEffect(() => {
    (async () => {
      setLoadingAll(true);
      try {
        const pageSize = 100;
        let offset = 0;
        let total = Infinity;
        const acc: FeedModel[] = [];
        while (offset < total) {
          const res = await api.feedsList(
            { subscribedOnly: false, limit: pageSize, offset },
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
          toast.error(message || "failed to load feeds");
        }
      } finally {
        setLoadingAll(false);
      }
    })();
  }, [api, logout]);

  const handleCreate = async () => {
    if (!newUrl.trim()) return;
    setCreating(true);
    try {
      await api.feedsCreate(
        { feed_url: newUrl.trim() },
        { secure: true, format: "json" },
      );
      setPage(1);
      setSearch("");
      setSortKey("title");
      setSortAsc(true);
      const reload = async () => {
        const pageSize = 100;
        let offset = 0;
        let total = Infinity;
        const acc: FeedModel[] = [];
        while (offset < total) {
          const res = await api.feedsList(
            { subscribedOnly: false, limit: pageSize, offset },
            { secure: true, format: "json" },
          );
          const chunk = res.data.feeds ?? [];
          acc.push(...chunk);
          total = res.data.total_count ?? chunk.length;
          offset += chunk.length;
          if (!chunk.length) break;
        }
        setAllFeeds(acc);
      };
      await reload();
      setNewUrl("");
      toast.success("feed imported");
    } catch (err: any) {
      if (err.error === "Unauthorized") logout();
      else {
        const message = await err.text();
        toast.error(message || "failed to create feed");
      }
    } finally {
      setCreating(false);
    }
  };

  const processed = allFeeds
    .filter((f) => f.title?.toLowerCase().includes(search.toLowerCase()))
    .sort((a, b) => {
      let cmp: number;
      if (sortKey === "title") {
        cmp = a.title!.localeCompare(b.title!);
      } else {
        const at = a.updated_parsed ? new Date(a.updated_parsed).getTime() : 0;
        const bt = b.updated_parsed ? new Date(b.updated_parsed).getTime() : 0;
        cmp = at - bt;
      }
      return sortAsc ? cmp : -cmp;
    });

  const totalPages = Math.ceil(processed.length / PAGE_SIZE);
  const pageItems = processed.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE);

  const pages: (number | "ellipsis")[] = [];
  if (totalPages <= 7) {
    for (let i = 1; i <= totalPages; i++) pages.push(i);
  } else if (page <= 4) {
    pages.push(1, 2, 3, 4, 5, "ellipsis", totalPages);
  } else if (page > totalPages - 4) {
    pages.push(
      1,
      "ellipsis",
      totalPages - 4,
      totalPages - 3,
      totalPages - 2,
      totalPages - 1,
      totalPages,
    );
  } else {
    pages.push(1, "ellipsis", page - 1, page, page + 1, "ellipsis", totalPages);
  }

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />

      <main className="flex-grow bg-gray-50 p-6">
        <div className="container mx-auto px-4">
          {loadingAll ? (
            <div className="flex justify-center py-10">
              <Spinner />
            </div>
          ) : (
            <>
              {/* Controls */}
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6">
                {/* Import */}
                <div className="flex items-center gap-2">
                  <Input
                    placeholder="Feed URL"
                    value={newUrl}
                    onChange={(e) => setNewUrl(e.currentTarget.value)}
                    className="flex-1 min-w-0"
                  />
                  <Button
                    size="sm"
                    onClick={handleCreate}
                    disabled={creating || !newUrl.trim()}
                    className="flex-shrink-0"
                  >
                    {creating ? (
                      <Spinner size="sm" />
                    ) : (
                      <PlusCircle className="mr-1" />
                    )}
                    Import
                  </Button>
                </div>

                {/* Search & Sort */}
                <div className="flex items-center gap-2">
                  <Input
                    placeholder="Search..."
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
                      if (sortKey === "title" && sortAsc) {
                        setSortAsc(false);
                      } else if (sortKey === "title" && !sortAsc) {
                        setSortKey("updated_parsed");
                        setSortAsc(true);
                      } else if (sortKey === "updated_parsed" && sortAsc) {
                        setSortAsc(false);
                      } else {
                        setSortKey("title");
                        setSortAsc(true);
                      }
                      setPage(1);
                    }}
                    className="flex-shrink-0"
                  >
                    sort by {labelMap[sortKey]} {sortAsc ? "↑" : "↓"}
                  </Button>
                </div>
              </div>

              {/* Feed grid or empty */}
              {pageItems.length === 0 ? (
                <p className="text-center">No feeds found.</p>
              ) : (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                  {pageItems.map((f) => (
                    <FeedPreview key={f.id} feed={f} />
                  ))}
                </div>
              )}

              {/* Pagination */}
              {totalPages > 1 && (
                <div className="mt-8 flex justify-center">
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
                          <PaginationItem key={i}>
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
                          onClick={() =>
                            setPage((p) => Math.min(totalPages, p + 1))
                          }
                          aria-disabled={page === totalPages}
                        />
                      </PaginationItem>
                    </PaginationContent>
                  </Pagination>
                </div>
              )}
            </>
          )}
        </div>
      </main>

      <div className="mt-auto">
        <Footer />
      </div>
    </div>
  );
}
