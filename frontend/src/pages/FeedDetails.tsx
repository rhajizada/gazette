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
import { Star } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import { Navigate, useParams } from "react-router-dom";
import { toast } from "sonner";
import type {
  GithubComRhajizadaGazetteInternalServiceFeed as FeedType,
  GithubComRhajizadaGazetteInternalServiceItem as ItemModel,
} from "../api/data-contracts";
import { useAuth } from "../context/AuthContext";

export default function FeedDetails() {
  const { feedID } = useParams<{ feedID: string }>();
  const { api, logout } = useAuth();
  const PAGE_SIZE = 10;

  const [feed, setFeed] = useState<FeedType | null>(null);
  const [items, setItems] = useState<ItemModel[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [itemsLoading, setItemsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [notFound, setNotFound] = useState(false);
  const [subscribed, setSubscribed] = useState(false);
  const [subLoading, setSubLoading] = useState(false);

  const [itemSearch, setItemSearch] = useState("");
  type SortKey = "title" | "published_parsed";
  const [itemSortKey, setItemSortKey] = useState<SortKey>("published_parsed");
  const [itemSortAsc, setItemSortAsc] = useState(false);
  const labelMap: Record<SortKey, string> = {
    title: "title",
    published_parsed: "published",
  };

  useEffect(() => {
    if (!feedID) return;

    setLoading(true);
    api
      .feedsDetail(feedID, { secure: true, format: "json" })
      .then((res) => {
        setFeed(res.data);
        setSubscribed(res.data.subscribed ?? false);
      })
      .catch((err: any) => {
        if (err.status === 404 || err.status === 400) {
          setNotFound(true);
        } else {
          const message = err.text();
          toast.error(message || "failed to load feed");
          setError("Failed to load feed");
          if (err.error === "Unauthorized") logout();
        }
      })
      .finally(() => setLoading(false));
  }, [api, feedID, logout]);

  const fetchItems = useCallback(() => {
    if (!feedID) return;
    setItemsLoading(true);
    api
      .feedsItemsList(
        feedID,
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
  }, [api, feedID, page, logout]);

  useEffect(() => {
    fetchItems();
  }, [fetchItems]);

  const toggleSub = async () => {
    if (!feed || subLoading) return;
    setSubLoading(true);
    try {
      if (subscribed) {
        await api.feedsSubscribeDelete(feed.id!, { format: "json" });
        setSubscribed(false);
      } else {
        await api.feedsSubscribeUpdate(feed.id!, { format: "json" });
        setSubscribed(true);
      }
    } catch (e) {
      console.error(e);
    } finally {
      setSubLoading(false);
    }
  };

  if (loading) return <Loader />;
  if (notFound || !feed) return <Navigate to="/*" />;
  if (error || !feed)
    return (
      <div className="flex flex-col min-h-screen">
        <Navbar />
        <div className="flex-grow flex items-center justify-center">
          <p className="text-red-600">{error || "Feed not available"}</p>
        </div>
        <Footer />
      </div>
    );

  const displayed = items
    .filter(
      (it) =>
        it.title?.toLowerCase().includes(itemSearch.toLowerCase()) ||
        it.description?.toLowerCase().includes(itemSearch.toLowerCase()),
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
          <div className="max-w-4xl mx-auto">
            <h1 className="text-4xl font-bold text-black">{feed.title}</h1>
            {feed.description && (
              <p className="text-lg text-gray-700 mt-2">{feed.description}</p>
            )}
            {feed.updated_parsed && (
              <p className="text-sm text-gray-500 mt-1">
                Last updated:{" "}
                {new Date(feed.updated_parsed).toLocaleDateString()}
              </p>
            )}
            <div className="mt-4">
              <Button
                size="sm"
                onClick={toggleSub}
                disabled={subLoading}
                className="inline-flex items-center bg-white text-gray-800 px-3 py-1 shadow hover:bg-gray-100 transition disabled:opacity-50"
              >
                <Star
                  fill={subscribed ? "currentColor" : "none"}
                  className={`w-5 h-5 mr-2 ${
                    subscribed ? "text-yellow-400" : "text-gray-800"
                  }`}
                />
                {subscribed ? "Unsubscribe" : "Subscribe"}
              </Button>
            </div>
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
                  } else if (
                    itemSortKey === "published_parsed" &&
                    itemSortAsc
                  ) {
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
