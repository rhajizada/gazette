import { Footer } from "@/components/Footer";
import { ItemPreview } from "@/components/ItemPreview";
import { Navbar } from "@/components/Navbar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
} from "@/components/ui/popover";
import {
  Command,
  CommandInput,
  CommandList,
  CommandEmpty,
  CommandGroup,
  CommandItem,
} from "@/components/ui/command";
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
import { Check, ChevronsUpDown } from "lucide-react";
import { cn } from "@/lib/utils";
import { useCallback, useEffect, useState, ChangeEvent } from "react";
import { useSearchParams } from "react-router-dom";
import { toast } from "sonner";
import type { GithubComRhajizadaGazetteInternalServiceItem as ItemModel } from "../api/data-contracts";
import { useAuth } from "../context/AuthContext";

export default function CategoriesPage() {
  const { api, logout } = useAuth();
  const UI_PAGE_SIZE = 15;

  const [searchParams, setSearchParams] = useSearchParams();
  const initialNames = searchParams.getAll("name");
  const [selectedCategories, setSelectedCategories] = useState<string[]>(
    initialNames || [],
  );

  const [allCategories, setAllCategories] = useState<string[]>([]);
  const [catsLoading, setCatsLoading] = useState(true);
  const [catSearch, setCatSearch] = useState("");
  const [open, setOpen] = useState(false);

  const [rawItems, setRawItems] = useState<ItemModel[]>([]);
  const [itemsLoading, setItemsLoading] = useState(false);

  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  type SortKey = "title" | "published_parsed";
  const [sortKey, setSortKey] = useState<SortKey>("published_parsed");
  const [sortAsc, setSortAsc] = useState(false);
  const labelMap: Record<SortKey, string> = {
    title: "title",
    published_parsed: "published",
  };

  const joined = selectedCategories.join(", ");
  const displayText = selectedCategories.length
    ? joined.length > 16
      ? `${joined.substring(0, 16)}…`
      : joined
    : "Select categories...";

  useEffect(() => {
    const params = new URLSearchParams();
    selectedCategories.forEach((n) => params.append("name", n));
    setSearchParams(params, { replace: true });
    setPage(1);
    fetchAllItems();
  }, [selectedCategories]);

  useEffect(() => {
    async function loadAllCategories() {
      setCatsLoading(true);
      try {
        const chunkSize = 100;
        let offset = 0;
        const all: string[] = [];
        while (true) {
          const res = await api.categoriesList(
            { limit: chunkSize, offset },
            { secure: true, format: "json" },
          );
          const chunk = res.data.categories ?? [];
          all.push(...chunk);
          if (chunk.length < chunkSize) break;
          offset += chunkSize;
        }
        setAllCategories(all);
      } catch (err: any) {
        if (err.status === 401 || err.error === "Unauthorized") logout();
        else toast.error("Failed to load categories");
      } finally {
        setCatsLoading(false);
      }
    }
    loadAllCategories();
  }, [api, logout]);

  const toggleCategory = useCallback((name: string) => {
    setSelectedCategories((prev) =>
      prev.includes(name) ? prev.filter((c) => c !== name) : [...prev, name],
    );
  }, []);

  const fetchAllItems = useCallback(async () => {
    if (!selectedCategories.length) {
      setRawItems([]);
      return;
    }
    setItemsLoading(true);
    try {
      const chunkSize = 100;
      let offset = 0;
      const all: ItemModel[] = [];
      while (true) {
        const res = await api.categoriesItemsList(
          { name: selectedCategories, limit: chunkSize, offset },
          { secure: true, format: "json" },
        );
        const chunk = res.data.items ?? [];
        all.push(...chunk);
        if (chunk.length < chunkSize) break;
        offset += chunkSize;
      }
      setRawItems(all);
    } catch (err: any) {
      if (err.status === 401 || err.error === "Unauthorized") logout();
      else toast.error("Failed to load items");
    } finally {
      setItemsLoading(false);
    }
  }, [api, logout, selectedCategories]);

  const filteredCats = allCategories.filter((c) =>
    c.toLowerCase().includes(catSearch.toLowerCase()),
  );

  const displayed = rawItems
    .filter(
      (it) =>
        !search ||
        it.title?.toLowerCase().includes(search.toLowerCase()) ||
        it.description?.toLowerCase().includes(search.toLowerCase()),
    )
    .sort((a, b) => {
      const aKey =
        sortKey === "title"
          ? a.title || ""
          : new Date(a.published_parsed || 0).getTime();
      const bKey =
        sortKey === "title"
          ? b.title || ""
          : new Date(b.published_parsed || 0).getTime();
      const cmp =
        sortKey === "title"
          ? (aKey as string).localeCompare(bKey as string)
          : (aKey as number) - (bKey as number);
      return sortAsc ? cmp : -cmp;
    });
  const totalPages = Math.ceil(displayed.length / UI_PAGE_SIZE);
  const pageItems = displayed.slice(
    (page - 1) * UI_PAGE_SIZE,
    page * UI_PAGE_SIZE,
  );

  const includedMap = Object.fromEntries(
    selectedCategories.map((c) => [c, true]),
  );

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <main className="flex-grow mx-auto max-w-5xl px-6 py-8">
        <div className="flex flex-col md:flex-row items-center gap-4 mb-6">
          <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
              <Button
                variant="outline"
                role="combobox"
                aria-expanded={open}
                className="w-[200px] justify-between"
              >
                {displayText}
                <ChevronsUpDown className="opacity-50 ml-2 h-4 w-4" />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-[200px] p-2">
              {catsLoading ? (
                <Spinner className="m-4" />
              ) : (
                <Command>
                  <CommandInput
                    placeholder="Search categories..."
                    className="h-9"
                    value={catSearch}
                    onValueChange={setCatSearch}
                  />
                  <CommandList>
                    <CommandEmpty>No categories found.</CommandEmpty>
                    <CommandGroup>
                      {filteredCats.map((c) => (
                        <CommandItem
                          key={c}
                          value={c}
                          onSelect={() => {
                            toggleCategory(c);
                            setOpen(false);
                          }}
                        >
                          {c}
                          <Check
                            className={cn(
                              "ml-auto h-4 w-4",
                              includedMap[c] ? "opacity-100" : "opacity-0",
                            )}
                          />
                        </CommandItem>
                      ))}
                    </CommandGroup>
                  </CommandList>
                </Command>
              )}
            </PopoverContent>
          </Popover>

          <Input
            placeholder="Search items…"
            value={search}
            onChange={(e: ChangeEvent<HTMLInputElement>) => {
              setSearch(e.currentTarget.value);
              setPage(1);
            }}
            className="flex-1"
          />

          <Button
            variant="outline"
            onClick={() => {
              if (sortKey === "title" && sortAsc) {
                setSortKey("published_parsed");
                setSortAsc(true);
              } else if (sortKey === "published_parsed" && sortAsc) {
                setSortAsc(false);
              } else if (sortKey === "published_parsed") {
                setSortKey("title");
                setSortAsc(true);
              } else {
                setSortAsc(false);
              }
              setPage(1);
            }}
          >
            sort by {labelMap[sortKey]} {sortAsc ? "↑" : "↓"}
          </Button>
        </div>

        <div className="flex flex-wrap gap-2 mb-6">
          {selectedCategories.map((c) => (
            <Badge
              key={c}
              className="cursor-pointer"
              onClick={() => toggleCategory(c)}
            >
              {c} ×
            </Badge>
          ))}
        </div>

        {itemsLoading ? (
          <div className="flex justify-center py-20">
            <Spinner />
          </div>
        ) : pageItems.length === 0 ? (
          <p className="text-center text-gray-600">No items to display.</p>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {pageItems.map((it) => (
              <ItemPreview key={it.id} item={it} />
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
                {(() => {
                  const pages: (number | "ellipsis")[] = [1];
                  if (page > 3) pages.push("ellipsis");
                  const start = Math.max(2, page - 2);
                  const end = Math.min(totalPages - 1, page + 2);
                  for (let p = start; p <= end; p++) pages.push(p);
                  if (page < totalPages - 2) pages.push("ellipsis");
                  if (totalPages > 1) pages.push(totalPages);
                  return pages.map((p, idx) =>
                    p === "ellipsis" ? (
                      <PaginationItem key={`e${idx}`}>
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
                  );
                })()}
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
      <Footer />
    </div>
  );
}
