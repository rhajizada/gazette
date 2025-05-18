import { ItemPreview } from "@/components/ItemPreview";
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
import { Navigate } from "react-router-dom";
import { toast } from "sonner";
import type { GithubComRhajizadaGazetteInternalServiceItem as ItemModel } from "../api/data-contracts";
import { useAuth } from "../context/AuthContext";

const PAGE_SIZE = 10;
const CHUNK_SIZE = 100;

type SortKey = "title" | "published_parsed";
const labelMap: Record<SortKey, string> = {
  title: "title",
  published_parsed: "published",
};

export function LikedItems() {
  const { api, logout } = useAuth();
  const [allItems, setAllItems] = useState<ItemModel[]>([]);
  const [preloading, setPreloading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [sortKey, setSortKey] = useState<SortKey>("published_parsed");
  const [sortAsc, setSortAsc] = useState(false);

  useEffect(() => {
    let cancelled = false;

    async function loadAll() {
      try {
        const first = await api.itemsList(
          { limit: CHUNK_SIZE, offset: 0 },
          { secure: true, format: "json" },
        );
        if (cancelled) return;

        const total = first.data.total_count ?? 0;
        let acc = first.data.items ?? [];
        setAllItems(acc.slice(0, total));

        let offset = CHUNK_SIZE;
        while (!cancelled && acc.length < total) {
          const res = await api.itemsList(
            { limit: CHUNK_SIZE, offset },
            { secure: true, format: "json" },
          );
          if (cancelled) break;

          const next = res.data.items ?? [];
          if (!next.length) break;

          acc = acc.concat(next).slice(0, total);
          setAllItems(acc);
          offset += CHUNK_SIZE;
        }
      } catch (err: any) {
        if (err.error === "Unauthorized") {
          logout();
        } else {
          const msg = err.text?.() ?? "failed to load liked items";
          toast.error(msg);
          setError(msg);
        }
      } finally {
        if (!cancelled) setPreloading(false);
      }
    }

    loadAll();
    return () => {
      cancelled = true;
    };
  }, [api, logout]);

  if (error) return <Navigate to="/*" replace />;

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

  // build pages array
  const pagesArr: (number | "ellipsis")[] = [];
  if (totalPages > 0) {
    pagesArr.push(1);
    if (page > 3) pagesArr.push("ellipsis");
    for (
      let p = Math.max(2, page - 2);
      p <= Math.min(totalPages - 1, page + 2);
      p++
    ) {
      pagesArr.push(p);
    }
    if (page < totalPages - 2) pagesArr.push("ellipsis");
    if (totalPages > 1) pagesArr.push(totalPages);
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-2">
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
            if (sortKey === "title" && sortAsc) setSortAsc(false);
            else if (sortKey === "title") {
              setSortKey("published_parsed");
              setSortAsc(true);
            } else if (sortAsc) setSortAsc(false);
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

      {preloading ? (
        <div className="flex justify-center py-10">
          <Spinner />
        </div>
      ) : paginated.length === 0 ? (
        <p className="text-center text-gray-500">No liked items found.</p>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {paginated.map((item) => (
            <ItemPreview key={item.id} item={item} />
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
