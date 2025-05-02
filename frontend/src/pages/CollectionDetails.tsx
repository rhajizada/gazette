import { useCallback, useEffect, useState } from "react";
import { Navbar } from "@/components/Navbar";
import { Footer } from "@/components/Footer";
import { ItemPreview } from "@/components/ItemPreview";
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

  useEffect(() => {
    if (!collectionID) return;
    setLoading(true);
    api.collectionsDetail(collectionID, { secure: true, format: "json" })
      .then(res => setCollection(res.data))
      .catch(err => {
        if (err.status === 404) setNotFound(true);
        else {
          toast.error("Failed to load collection");
          setError("Failed to load collection");
          if (err.error === "Unauthorized") logout();
        }
      })
      .finally(() => setLoading(false));
  }, [api, collectionID, logout]);

  const fetchItems = useCallback(() => {
    if (!collectionID) return;
    setItemsLoading(true);
    api.collectionsItemsList(
      collectionID,
      { limit: PAGE_SIZE, offset: (page - 1) * PAGE_SIZE },
      { secure: true, format: "json" }
    )
      .then(res => {
        setItems(res.data.items ?? []);
        setTotal(res.data.total_count ?? 0);
      })
      .catch(err => {
        if (err.error === "Unauthorized") logout();
        else console.error(err);
      })
      .finally(() => setItemsLoading(false));
  }, [api, collectionID, page, logout]);

  useEffect(() => { fetchItems(); }, [fetchItems]);

  if (loading) return <Loader />;
  if (notFound || !collection) return <Navigate to="*" replace />;
  if (error || !collection)
    return (
      <>
        <Navbar />
        <div className="p-6 text-red-600">{error || "Collection not available"}</div>
        <Footer />
      </>
    );

  const totalPages = Math.ceil(total / PAGE_SIZE);
  const pages: (number | "ellipsis")[] = [];
  if (totalPages <= 7) {
    for (let i = 1; i <= totalPages; i++) pages.push(i);
  } else if (page <= 4)
    pages.push(1, 2, 3, 4, 5, "ellipsis", totalPages);
  else if (page > totalPages - 4)
    pages.push(
      1,
      "ellipsis",
      totalPages - 4,
      totalPages - 3,
      totalPages - 2,
      totalPages - 1,
      totalPages
    );
  else pages.push(1, "ellipsis", page - 1, page, page + 1, "ellipsis", totalPages);

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <main className="flex-grow mx-auto max-w-5xl px-6 py-8">
        <div className="mb-6">
          <h1 className="text-4xl font-bold">{collection.name}</h1>
          {collection.last_updated && (
            <p className="text-sm text-gray-500 mt-1">
              Last updated: {new Date(collection.last_updated).toLocaleDateString()}
            </p>
          )}
        </div>
        {itemsLoading ? (
          <Loader />
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {items.map(item => (
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
                    onClick={() => setPage(page - 1)}
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
                  )
                )}
                <PaginationItem>
                  <PaginationNext
                    onClick={() => setPage(page + 1)}
                    aria-disabled={page === totalPages}
                  />
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

function Loader() {
  return (
    <div className="flex items-center justify-center min-h-[300px]">
      <Spinner />
    </div>
  );
}

