import { Footer } from "@/components/Footer";
import { Navbar } from "@/components/Navbar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Spinner } from "@/components/ui/spinner";
import { cn } from "@/lib/utils";
import { Check, ChevronsUpDown, Heart as HeartIcon } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import { Link, Navigate, useParams } from "react-router-dom";
import { toast } from "sonner";
import type {
  GithubComRhajizadaGazetteInternalServiceCollection as CollectionModel,
  GithubComRhajizadaGazetteInternalServiceItem as ItemModel,
} from "../api/data-contracts";
import { useAuth } from "../context/AuthContext";

export default function ItemDetails() {
  const { itemID } = useParams<{ itemID: string }>();
  const { api, logout } = useAuth();

  const [item, setItem] = useState<ItemModel | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [notFound, setNotFound] = useState(false);

  const [liked, setLiked] = useState(false);
  const [likeLoading, setLikeLoading] = useState(false);

  const [collections, setCollections] = useState<CollectionModel[]>([]);
  const [userCollections, setMyCollections] = useState<CollectionModel[]>([]);
  const [included, setIncluded] = useState<Record<string, boolean>>({});
  const [collectionsLoading, setCollectionsLoading] = useState(false);

  const [open, setOpen] = useState(false);

  const getColor = (id: string) => {
    let hash = 0;
    for (let i = 0; i < id.length; i++) {
      hash = id.charCodeAt(i) + ((hash << 5) - hash);
    }
    const hue = Math.abs(hash) % 360;
    return `hsl(${hue}, 70%, 50%)`;
  };

  // fetch item detail
  useEffect(() => {
    if (!itemID) return;
    setLoading(true);
    api
      .itemsDetail(itemID, { secure: true, format: "json" })
      .then((res) => {
        setItem(res.data);
        setLiked(res.data.liked ?? false);
      })
      .catch((err) => {
        if (err.error === "Unauthorized") logout();
        if (err.status === 404 || err.status === 400) {
          setNotFound(true);
        } else {
          const message = err.text();
          toast.error(message || "failed to load item");
          setError("failed to load item");
        }
      })
      .finally(() => setLoading(false));
  }, [api, itemID, logout]);

  // fetch collections logic
  const fetchCollections = useCallback(async () => {
    if (!itemID) return;
    setCollectionsLoading(true);
    try {
      const [allRes, firstMyRes] = await Promise.all([
        api.collectionsList(
          { limit: 50, offset: 0 },
          { secure: true, format: "json" },
        ),
        api.itemsCollectionsList(
          itemID,
          { limit: 50, offset: 0 },
          { secure: true, format: "json" },
        ),
      ]);

      const allCols = allRes.data.collections || [];
      setCollections(allCols);

      const perPage = 50;
      const total =
        firstMyRes.data.total_count ?? firstMyRes.data.collections?.length ?? 0;
      let myCols = firstMyRes.data.collections || [];

      const pages = Math.ceil(total / perPage);
      for (let page = 1; page < pages; page++) {
        const pageRes = await api.itemsCollectionsList(
          itemID,
          { limit: perPage, offset: page * perPage },
          { secure: true, format: "json" },
        );
        myCols = myCols.concat(pageRes.data.collections || []);
      }

      setMyCollections(myCols);
      setIncluded((prev) => ({
        ...prev,
        ...Object.fromEntries(myCols.map((c) => [c.id!, true])),
      }));
    } catch (err: any) {
      if (err.error === "Unauthorized") logout();
      else {
        const message = err.text();
        toast.error(message || "failed to load collections");
      }
    } finally {
      setCollectionsLoading(false);
    }
  }, [api, itemID]);

  useEffect(() => {
    if (item && !error && !notFound) {
      fetchCollections();
    }
  }, [item, error, notFound, fetchCollections]);

  const handleOpenChange = (isOpen: boolean) => {
    setOpen(isOpen);
    if (isOpen && !collections.length && !collectionsLoading) {
      fetchCollections();
    }
  };

  const toggleLike = async () => {
    if (!item || likeLoading) return;
    setLikeLoading(true);
    try {
      if (liked) {
        await api.itemsLikeDelete(item.id!, { secure: true, format: "json" });
      } else {
        await api.itemsLikeCreate(item.id!, { secure: true, format: "json" });
      }
      setLiked((prev) => !prev);
      toast.success(
        liked ? "removed from liked items" : "added to liked items",
      );
    } catch (err: any) {
      if (err.error === "Unauthorized") logout();
      else {
        const message = await err.text();
        toast.error(message || "failed to update item");
      }
    } finally {
      setLikeLoading(false);
    }
  };

  const toggleCollection = async (colId: string, colName: string) => {
    try {
      if (included[colId]) {
        await api.collectionsItemDelete(colId, itemID!, { secure: true });
        setIncluded((prev) => ({ ...prev, [colId]: false }));
        setMyCollections((prev) => prev.filter((c) => c.id !== colId));
      } else {
        await api.collectionsItemCreate(colId, itemID!, { secure: true });
        setIncluded((prev) => ({ ...prev, [colId]: true }));
        const added = collections.find((c) => c.id === colId);
        if (added) setMyCollections((prev) => [...prev, added]);
      }
      toast.success(
        included[colId]
          ? `removed from collection ${colName}`
          : `added to collection ${colName}`,
      );
    } catch (err: any) {
      if (err.error === "Unauthorized") logout();
      else {
        const message = await err.text();
        toast.error(message || "failed to update collection");
      }
    }
  };

  if (loading) {
    return (
      <div className="flex flex-col min-h-screen">
        <Navbar />
        <div className="flex-grow flex items-center justify-center">
          <Spinner />
        </div>
        <Footer />
      </div>
    );
  }

  if (notFound) {
    return <Navigate to="/*" />;
  }

  if (error || !item) {
    return (
      <div className="flex flex-col min-h-screen">
        <Navbar />
        <main className="flex-grow p-6 text-red-600">
          {error || "Item not available"}
        </main>
        <Footer />
      </div>
    );
  }

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <main className="flex-grow mx-auto max-w-5xl px-6 py-8 grid grid-cols-1 lg:grid-cols-[3fr_1fr] gap-8">
        <div>
          <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
            {item.title}
          </h1>
          {item.authors && item.authors.length > 0 && (
            <div className="flex flex-wrap gap-2 mt-4">
              {item.authors.map((a, i) => (
                <Badge key={i}>{a.name || a.email}</Badge>
              ))}
            </div>
          )}
          <p className="leading-7 [&:not(:first-child)]:mt-6 text-sm text-gray-500">
            Published:{" "}
            {item.published_parsed &&
              new Date(item.published_parsed).toLocaleString()}
            <br />
            Updated:{" "}
            {item.updated_parsed &&
              new Date(item.updated_parsed).toLocaleString()}
          </p>
          <br />
          {item.content && (
            <div
              className="mt-6 leading-7 max-w-prose mx-auto"
              dangerouslySetInnerHTML={{ __html: item.content }}
            />
          )}
          {item.link && (
            <div className="mt-6 border-l-2 pl-4 italic">
              <Button asChild variant="link">
                <a href={item.link} target="_blank" rel="noopener noreferrer">
                  Read Original Source
                </a>
              </Button>
            </div>
          )}
          <div className="mt-8">
            <button
              onClick={toggleLike}
              disabled={likeLoading}
              aria-label={liked ? "Unlike" : "Like"}
              className={cn(
                "p-2 rounded-full transition-colors",
                liked
                  ? "text-red-500 hover:bg-red-100"
                  : "text-gray-500 hover:bg-gray-100",
                likeLoading && "opacity-50 cursor-not-allowed",
              )}
            >
              <HeartIcon
                fill={liked ? "currentColor" : "none"}
                strokeWidth={2}
                className="w-6 h-6"
              />
            </button>
          </div>
        </div>
        <aside>
          {item.categories && item.categories.length > 0 && (
            <div className="mt-8 clear-left">
              <h2 className="text-3xl font-semibold">Categories</h2>
              <br />
              <div className="flex flex-wrap gap-2">
                {item.categories.slice(0, 3).map((cat) => (
                  <Badge key={cat}>{cat}</Badge>
                ))}
              </div>
            </div>
          )}
          {item.enclosures?.length > 0 && (
            <div className="mt-8">
              <h2 className="text-3xl font-semibold">Attachments</h2>
              <div className="mt-4 flex flex-col gap-4">
                {item.enclosures.map((enc: any, i: number) =>
                  enc.type?.startsWith("audio/") ? (
                    <audio key={i} controls src={enc.url} className="w-full" />
                  ) : enc.type?.startsWith("video/") ? (
                    <video key={i} controls src={enc.url} className="w-full" />
                  ) : (
                    <a
                      key={i}
                      href={enc.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="font-medium underline underline-offset-4"
                    >
                      Download ({enc.type || "file"})
                    </a>
                  ),
                )}
              </div>
            </div>
          )}
          <br />
          <h2 className="text-3xl font-semibold">Collections</h2>
          <br />
          <Popover open={open} onOpenChange={handleOpenChange}>
            <PopoverTrigger asChild>
              <Button
                variant="outline"
                role="combobox"
                aria-expanded={open}
                className="w-full justify-between"
              >
                Add to collection
                <ChevronsUpDown className="ml-2 h-4 w-4 opacity-50" />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-[250px] p-0">
              {collectionsLoading ? (
                <Spinner className="m-4" />
              ) : (
                <Command>
                  <CommandInput placeholder="Search collections..." />
                  <CommandList>
                    <CommandEmpty>No collections.</CommandEmpty>
                    <CommandGroup>
                      {collections.map((c) => (
                        <CommandItem
                          key={c.id}
                          value={c.id!}
                          onSelect={() => toggleCollection(c.id!, c.name!)}
                        >
                          <Check
                            className={cn(
                              "mr-2 h-4 w-4",
                              included[c.id!] ? "opacity-100" : "opacity-0",
                            )}
                          />
                          {c.name}
                        </CommandItem>
                      ))}
                    </CommandGroup>
                  </CommandList>
                </Command>
              )}
            </PopoverContent>
          </Popover>
          <br />
          <div className="mt-4 flex flex-wrap gap-2">
            {userCollections.map((c) => (
              <Badge
                key={c.id}
                variant="outline"
                style={{ backgroundColor: getColor(c.id!), color: "white" }}
              >
                <Link to={`/collections/${c.id}`}>{c.name}</Link>
              </Badge>
            ))}
          </div>
        </aside>
      </main>
      <Footer />
    </div>
  );
}
