import { Footer } from "@/components/Footer";
import { Navbar } from "@/components/Navbar";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { PlusCircle, Trash2 } from "lucide-react";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { toast } from "sonner";
import type { GithubComRhajizadaGazetteInternalServiceCollection as CollectionModel } from "../api/data-contracts";
import { useAuth } from "../context/AuthContext";

export default function Collections() {
  const { api, logout } = useAuth();
  const [collections, setCollections] = useState<CollectionModel[]>([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [newName, setNewName] = useState("");
  const [search, setSearch] = useState("");

  type SortKey = "name" | "created_at" | "last_updated";
  const [sortKey, setSortKey] = useState<SortKey>("name");
  const [sortAsc, setSortAsc] = useState(true);

  const labelMap: Record<SortKey, string> = {
    name: "name",
    created_at: "created",
    last_updated: "last updated",
  };

  useEffect(() => {
    (async () => {
      setLoading(true);
      try {
        const pageSize = 100;
        let offset = 0;
        let total = Infinity;
        const acc: CollectionModel[] = [];
        while (offset < total) {
          const res = await api.collectionsList(
            { limit: pageSize, offset },
            { secure: true, format: "json" },
          );
          const chunk = res.data.collections || [];
          acc.push(...chunk);
          total = res.data.total_count ?? chunk.length;
          offset += chunk.length;
          if (!chunk.length) break;
        }
        setCollections(acc);
      } catch (err: any) {
        console.error(err);
        if (err.error === "Unauthorized") logout();
        else {
          const message = await err.text();
          toast.error(message || "failed to load collections");
        }
      } finally {
        setLoading(false);
      }
    })();
  }, [api, logout]);

  const handleCreate = async () => {
    if (!newName.trim()) return;
    setCreating(true);
    try {
      const res = await api.collectionsCreate(
        { name: newName },
        { secure: true, format: "json" },
      );
      setCollections((prev) => [...prev, res.data]);
      setNewName("");
      toast.success("collection created");
    } catch (err: any) {
      const message = await err.text();
      toast.error(message || "failed to create collection");
    } finally {
      setCreating(false);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await api.collectionsDelete(id, { secure: true });
      setCollections((prev) => prev.filter((c) => c.id !== id));
      toast.success("collection deleted");
    } catch (err: any) {
      const message = await err.text();
      toast.error(message || "failed to delete collection");
    }
  };

  const filtered = collections
    .filter((c) => c.name?.toLowerCase().includes(search.toLowerCase()))
    .sort((a, b) => {
      let cmp: number;
      if (sortKey === "name") {
        cmp = a.name!.localeCompare(b.name!);
      } else {
        const aRaw =
          (sortKey === "created_at" ? a.created_at : a.last_updated) ?? "";
        const bRaw =
          (sortKey === "created_at" ? b.created_at : b.last_updated) ?? "";
        const aTime = aRaw ? new Date(aRaw).getTime() : 0;
        const bTime = bRaw ? new Date(bRaw).getTime() : 0;
        cmp = aTime - bTime;
      }
      return sortAsc ? cmp : -cmp;
    });

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <main className="flex-grow container mx-auto px-4 py-6">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6">
          <div className="flex items-center gap-2">
            <Input
              placeholder="New collection"
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              className="flex-1 min-w-0"
            />
            <Button
              onClick={handleCreate}
              disabled={creating}
              size="sm"
              className="flex-shrink-0"
            >
              {creating ? (
                <Spinner size="sm" />
              ) : (
                <PlusCircle className="mr-1" />
              )}
            </Button>
          </div>

          <div className="flex items-center gap-2">
            <Input
              placeholder="Search..."
              value={search}
              onChange={(e) => {
                setSearch(e.target.value);
              }}
              className="flex-1 min-w-0"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                if (sortKey === "name" && sortAsc) {
                  setSortAsc(false);
                } else if (sortKey === "name" && !sortAsc) {
                  setSortKey("created_at");
                  setSortAsc(true);
                } else if (sortKey === "created_at" && sortAsc) {
                  setSortAsc(false);
                } else if (sortKey === "created_at" && !sortAsc) {
                  setSortKey("last_updated");
                  setSortAsc(true);
                } else if (sortKey === "last_updated" && sortAsc) {
                  setSortAsc(false);
                } else {
                  setSortKey("name");
                  setSortAsc(true);
                }
              }}
              className="flex-shrink-0"
            >
              sort by {labelMap[sortKey]} {sortAsc ? "↑" : "↓"}
            </Button>
          </div>
        </div>

        {loading ? (
          <div className="flex justify-center py-10">
            <Spinner />
          </div>
        ) : (
          <div className="shadow-md rounded-lg overflow-hidden">
            <Table className="border-collapse">
              <TableHeader>
                <TableRow>
                  <TableHead className="px-3 py-2">Name</TableHead>
                  <TableHead className="px-3 py-2">Created</TableHead>
                  <TableHead className="px-3 py-2">Last Updated</TableHead>
                  <TableHead className="px-3 py-2"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filtered.map((c) => (
                  <TableRow key={c.id} className="hover:bg-gray-50">
                    <TableCell className="px-3 py-2">
                      <Link
                        to={`/collections/${c.id}`}
                        className="text-blue-600 hover:underline"
                      >
                        {c.name}
                      </Link>
                    </TableCell>
                    <TableCell className="px-3 py-2">
                      {c.created_at
                        ? new Date(c.created_at).toLocaleDateString()
                        : "-"}
                    </TableCell>
                    <TableCell className="px-3 py-2">
                      {c.last_updated
                        ? new Date(c.last_updated).toLocaleDateString()
                        : "-"}
                    </TableCell>
                    <TableCell className="px-3 py-2">
                      <AlertDialog>
                        <AlertDialogTrigger asChild>
                          <Button variant="destructive" size="sm">
                            <Trash2 size={14} />
                          </Button>
                        </AlertDialogTrigger>
                        <AlertDialogContent>
                          <AlertDialogHeader>
                            <AlertDialogTitle>
                              Delete Collection?
                            </AlertDialogTitle>
                            <AlertDialogDescription>
                              This action cannot be undone. Delete "{c.name}"?
                            </AlertDialogDescription>
                          </AlertDialogHeader>
                          <AlertDialogFooter>
                            <AlertDialogCancel>Cancel</AlertDialogCancel>
                            <AlertDialogAction
                              onClick={() => handleDelete(c.id!)}
                            >
                              Delete
                            </AlertDialogAction>
                          </AlertDialogFooter>
                        </AlertDialogContent>
                      </AlertDialog>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}
      </main>
      <div className="mt-auto">
        <Footer />
      </div>
    </div>
  );
}
