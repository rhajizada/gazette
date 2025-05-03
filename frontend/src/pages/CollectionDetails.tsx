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
import { Navigate, useParams } from "react-router-dom";
import { toast } from "sonner";
import type {
  GithubComRhajizadaGazetteInternalServiceCollection as CollectionType,
  GithubComRhajizadaGazetteInternalServiceItem as ItemModel,
} from "../api/data-contracts";
import { useAuth } from "../context/AuthContext";

export default function CollectionDetail() {
  const { collectionID } = useParams<{ collectionID: string }>();
  const { api, logout } = useAuth();
  const PAGE_SIZE = 10;

  const [collection, setCollection] = useState<CollectionType | null>(null);
  const [items, setItems] = useState<ItemModel[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [itemsLoading, setItemsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [notFound, setNotFound] = useState(false);

  const [itemSearch, setItemSearch] = useState("");
  type SortKey = "title" | "published_parsed";
  const [itemSortKey, setItemSortKey] = useState<SortKey>("published_parsed");
  const [itemSortAsc, setItemSortAsc] = useState(false);
  const labelMap: Record<SortKey, string> = {
    title: "title",
    published_parsed: "published",
  };

  useEffect(() => {
    if (!collectionID) return;
    setLoading(true);
    api
      .collectionsDetail(collectionID, { secure: true, format: "json" })
      .then((res) => setCollection(res.data))
      .catch((err) => {
        if (err.status === 404) setNotFound(true);
        else {
          const message = err.text();
          toast.error(message || "failed to load collections");
          setError("Failed to load collections");
          if (err.error === "Unauthorized") logout();
        }
      })
      .finally(() => setLoading(false));
  }, [api, collectionID, logout]);

  const fetchItems = useCallback(() => {
    if (!collectionID) return;
    setItemsLoading(true);
    api
      .collectionsItemsList(
        collectionID,
        { limit: PAGE_SIZE, offset: (page - 1) * PAGE_SIZE },
        { secure: true, format: "json" },
      )
      .then((res) => {
        setItems(res.data.items ?? []);
        setTotal(res.data.total_count ?? 0);
      })
      .catch((err) => {
        if (err.error === "Unauthorized") logout();
        else console.error(err);
      })
      .finally(() => setItemsLoading(false));
  }, [api, collectionID, page, logout]);

  useEffect(() => {
    fetchItems();
  }, [fetchItems]);

  if (loading) return <Loader />;
  if (notFound || !collection) return <Navigate to="*" replace />;
  if (error || !collection)
    return (
      <>
        <Navbar />
        <div className="p-6 text-red-600">
          {error || "Collection not available"}
        </div>
        <Footer />
      </>
    );

  const displayed = items
    .filter((it) =>
      [it.title, it.description]
        .filter(Boolean)
        .some((text) => text!.toLowerCase().includes(itemSearch.toLowerCase())),
    )
    .sort((a, b) => {
      let cmp: number;
      if (itemSortKey === "title") {
        cmp = a.title!.localeCompare(b.title!);
      } else {
        const aRaw = a.published_parsed;
        const bRaw = b.published_parsed;
        const aTime = aRaw ? new Date(aRaw).getTime() : 0;
        const bTime = bRaw ? new Date(bRaw).getTime() : 0;
        cmp = aTime - bTime;
      }
      return itemSortAsc ? cmp : -cmp;
    });

  const totalPages = Math.ceil(total / PAGE_SIZE);
  const pages: (number | "ellipsis")[] = [];
  if (totalPages <= 7) {
    for (let i = 1; i <= totalPages; i++) pages.push(i);
  } else if (page <= 4) pages.push(1, 2, 3, 4, 5, "ellipsis", totalPages);
  else if (page > totalPages - 4)
    pages.push(
      1,
      "ellipsis",
      totalPages - 4,
      totalPages - 3,
      totalPages - 2,
      totalPages - 1,
      totalPages,
    );
  else
    pages.push(1, "ellipsis", page - 1, page, page + 1, "ellipsis", totalPages);

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <main className="flex-grow mx-auto max-w-5xl px-6 py-8">
        <div className="mb-6">
          <h1 className="text-4xl font-bold">{collection.name}</h1>
          {collection.last_updated && (
            <p className="text-sm text-gray-500 mt-1">
              Last updated:{" "}
              {new Date(collection.last_updated).toLocaleDateString()}
            </p>
          )}
        </div>

        <div className="mt-8 grid grid-cols-1 gap-4 mb-6">
          <div className="flex items-center gap-2">
            <Input
              placeholder="Search..."
              value={itemSearch}
              onChange={(e) => {
                setItemSearch(e.currentTarget.value);
                setPage(1);
              }}
              className="flex-1 min-w-0"
            />
            <Button
              variant="outline"
              size="sm"
              className="flex-shrink-0"
              onClick={() => {
                if (itemSortKey === "title" && itemSortAsc) {
                  setItemSortAsc(false);
                } else if (itemSortKey === "title" && !itemSortAsc) {
                  setItemSortKey("published_parsed");
                  setItemSortAsc(true);
                } else if (itemSortKey === "published_parsed" && itemSortAsc) {
                  setItemSortAsc(false);
                } else {
                  setItemSortKey("title");
                  setItemSortAsc(true);
                }
                setPage(1);
              }}
            >
              sort by {labelMap[itemSortKey]} {itemSortAsc ? "↑" : "↓"}
            </Button>
          </div>
        </div>

        {itemsLoading ? (
          <Loader />
        ) : displayed.length === 0 ? (
          <p className="text-center">No items found.</p>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {displayed.map((item) => (
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
                {pages.map((p, idx) =>
                  p === "ellipsis" ? (
                    <PaginationItem key={idx}>
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
      </main>
      <div className="mt-auto">
        <Footer />
      </div>
    </div>
  );
}

function Loader() {
  return (
    <div className="flex items-center justify-center min-h-[300px]">
      <Spinner />
    </div>
  );
}
