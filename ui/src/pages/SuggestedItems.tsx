import { Footer } from "@/components/Footer";
import { ItemPreview } from "@/components/ItemPreview";
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
import { Progress } from "@/components/ui/progress";
import * as React from "react";
import { toast } from "sonner";
import type { GithubComRhajizadaGazetteInternalServiceItem as ItemModel } from "../api/data-contracts";
import { useAuth } from "../context/AuthContext";

export default function SuggestedItems() {
  const { api, logout } = useAuth();
  const PAGE_SIZE = 10;
  const CHUNK_SIZE = 100;

  const [allItems, setAllItems] = React.useState<ItemModel[]>([]);
  const [totalCount, setTotalCount] = React.useState(0);
  const [loadedCount, setLoadedCount] = React.useState(0);
  const [preloading, setPreloading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  const [page, setPage] = React.useState(1);
  const [search, setSearch] = React.useState("");
  type SortKey = "title" | "published_parsed";
  const [sortKey, setSortKey] = React.useState<SortKey>("published_parsed");
  const [sortAsc, setSortAsc] = React.useState(false);
  const labelMap: Record<SortKey, string> = {
    title: "title",
    published_parsed: "published",
  };

  React.useEffect(() => {
    let cancelled = false;

    async function loadAllPages() {
      try {
        const first = await api.suggestedList(
          { limit: CHUNK_SIZE, offset: 0 },
          { secure: true, format: "json" },
        );
        if (cancelled) return;

        const total = first.data.total_count ?? 0;
        setTotalCount(total);

        let acc = first.data.items!;
        setAllItems(acc);
        setLoadedCount(acc.length);

        const pages = Math.ceil(total / CHUNK_SIZE);
        for (let p = 2; p <= pages && !cancelled; p++) {
          const res = await api.suggestedList(
            { limit: CHUNK_SIZE, offset: (p - 1) * CHUNK_SIZE },
            { secure: true, format: "json" },
          );
          if (cancelled) break;

          const itemsPage = res.data.items!;
          acc = acc.concat(itemsPage);
          setAllItems([...acc]);
          setLoadedCount(acc.length);
        }
      } catch (err: any) {
        if (err.error === "Unauthorized") {
          logout();
        } else {
          const msg = err.text?.() ?? "failed to load suggested items";
          toast.error(msg);
          setError(msg);
        }
      } finally {
        if (!cancelled) setPreloading(false);
      }
    }

    loadAllPages();
    return () => {
      cancelled = true;
    };
  }, [api, logout]);

  if (error) {
    return <p className="text-center text-red-600">{error}</p>;
  }

  const processed = allItems
    .filter((it) =>
      [it.title, it.description]
        .filter(Boolean)
        .some((s) => s!.toLowerCase().includes(search.toLowerCase())),
    )
    .sort((a, b) => {
      let cmp: number;
      if (sortKey === "title") {
        cmp = a.title!.localeCompare(b.title!);
      } else {
        const ta = a.published_parsed
          ? new Date(a.published_parsed).getTime()
          : 0;
        const tb = b.published_parsed
          ? new Date(b.published_parsed).getTime()
          : 0;
        cmp = ta - tb;
      }
      return sortAsc ? cmp : -cmp;
    });

  const totalPages = Math.ceil(processed.length / PAGE_SIZE);
  const paginated = processed.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE);

  const pagesArr: (number | "ellipsis")[] = [];
  if (totalPages > 0) {
    pagesArr.push(1);
    if (page > 3) pagesArr.push("ellipsis");

    const start = Math.max(2, page - 2);
    const end = Math.min(totalPages - 1, page + 2);
    for (let p = start; p <= end; p++) {
      pagesArr.push(p);
    }

    if (page < totalPages - 2) pagesArr.push("ellipsis");
    if (totalPages > 1) pagesArr.push(totalPages);
  }

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <main className="flex-grow bg-gray-50 p-6">
        <div className="container mx-auto px-4">
          <h1 className="text-4xl font-bold text-black mb-4">Suggested</h1>

          <div className="mb-6 flex items-center gap-2">
            <Input
              placeholder="Search..."
              value={search}
              onChange={(e) => {
                setSearch(e.currentTarget.value);
                setPage(1);
              }}
              className="flex-1"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                if (sortKey === "title" && sortAsc) {
                  setSortAsc(false);
                } else if (sortKey === "title") {
                  setSortKey("published_parsed");
                  setSortAsc(true);
                } else if (sortAsc) {
                  setSortAsc(false);
                } else {
                  setSortKey("title");
                  setSortAsc(true);
                }
                setPage(1);
              }}
            >
              sort by {labelMap[sortKey]} {sortAsc ? "↑" : "↓"}
            </Button>
          </div>

          {preloading && (
            <div className="mb-8">
              <Progress
                value={
                  totalCount > 0
                    ? Math.round((loadedCount / totalCount) * 100)
                    : 0
                }
                className="w-full"
              />
              <p className="mt-1 text-sm text-gray-600">
                Loading ({loadedCount}/{totalCount})
              </p>
            </div>
          )}

          {!preloading && (
            <>
              {paginated.length === 0 ? (
                <p className="text-center">No suggested items.</p>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {paginated.map((item) => (
                    <ItemPreview key={item.id} item={item} />
                  ))}
                </div>
              )}

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
                      {pagesArr.map((p, i) =>
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
      <Footer />
    </div>
  );
}
