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
import { Spinner } from "@/components/ui/spinner";
import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";
import type { GithubComRhajizadaGazetteInternalServiceItem as ItemModel } from "../api/data-contracts";
import { useAuth } from "../context/AuthContext";

export default function SuggestedItems() {
  const { api, logout } = useAuth();
  const PAGE_SIZE = 10;

  const [items, setItems] = useState<ItemModel[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [search, setSearch] = useState("");
  type SortKey = "title" | "published_parsed";
  const [sortKey, setSortKey] = useState<SortKey>("published_parsed");
  const [sortAsc, setSortAsc] = useState(false);
  const labelMap: Record<SortKey, string> = {
    title: "title",
    published_parsed: "published",
  };

  const fetchItems = useCallback(() => {
    setLoading(true);
    api
      .suggestedList(
        { limit: PAGE_SIZE, offset: (page - 1) * PAGE_SIZE },
        { secure: true, format: "json" },
      )
      .then((res) => {
        setItems(res.data.items ?? []);
        setTotal(res.data.total_count ?? 0);
      })
      .catch((err: any) => {
        if (err.error === "Unauthorized") logout();
        else {
          const msg = err.text?.() ?? "Failed to load suggested items";
          toast.error(msg);
          setError(msg);
        }
      })
      .finally(() => setLoading(false));
  }, [api, logout, page]);

  useEffect(() => {
    fetchItems();
  }, [fetchItems]);

  if (error) return <p className="text-center text-red-600">{error}</p>;

  const filtered = items
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

  const totalPages = Math.ceil(total / PAGE_SIZE);
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

          {loading ? (
            <div className="flex justify-center py-20">
              <Spinner />
            </div>
          ) : filtered.length === 0 ? (
            <p className="text-center">No suggested items.</p>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {filtered.map((item) => (
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
        </div>
      </main>
      <Footer />
    </div>
  );
}
