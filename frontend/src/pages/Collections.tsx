import { useEffect, useState } from "react";
import { Navbar } from "@/components/Navbar";
import { Footer } from "@/components/Footer";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { AlertDialog, AlertDialogTrigger, AlertDialogContent, AlertDialogHeader, AlertDialogTitle, AlertDialogDescription, AlertDialogFooter, AlertDialogCancel, AlertDialogAction } from "@/components/ui/alert-dialog";
import { Trash2, PlusCircle } from "lucide-react";
import { toast } from "sonner";
import type { GithubComRhajizadaGazetteInternalServiceCollection as CollectionModel } from "../api/data-contracts";
import { useAuth } from "../context/AuthContext";
import { Link } from "react-router-dom";


export default function Collections() {
  const { api, logout } = useAuth();
  const [collections, setCollections] = useState<CollectionModel[]>([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [newName, setNewName] = useState("");
  const [search, setSearch] = useState("");
  const [sortKey, setSortKey] = useState<"name" | "created_at" | "last_updated">("name");
  const [sortAsc, setSortAsc] = useState(true);

  useEffect(() => {
    async function fetchAll() {
      setLoading(true);
      try {
        const pageSize = 100;
        let offset = 0;
        let total = Infinity;
        const all: CollectionModel[] = [];
        while (offset < total) {
          const res = await api.collectionsList({ limit: pageSize, offset }, { secure: true, format: "json" });
          const cols = res.data.collections || [];
          all.push(...cols);
          total = res.data.total_count ?? cols.length;
          offset += cols.length;
          if (!cols.length) break;
        }
        setCollections(all);
      } catch (err: any) {
        console.error(err);
        if (err.error === "Unauthorized") logout();
        else toast.error("Failed to load collections");
      } finally {
        setLoading(false);
      }
    }
    fetchAll();
  }, [api, logout]);

  const handleCreate = async () => {
    if (!newName.trim()) return;
    setCreating(true);
    try {
      const res = await api.collectionsCreate({ name: newName }, { secure: true, format: "json" });
      setCollections(prev => [...prev, res.data]);
      setNewName("");
      toast.success("Collection created");
    } catch (err) {
      console.error(err);
      toast.error("Failed to create collection");
    } finally {
      setCreating(false);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await api.collectionsDelete(id, { secure: true });
      setCollections(prev => prev.filter(c => c.id !== id));
      toast.success("Collection deleted");
    } catch (err) {
      console.error(err);
      toast.error("Failed to delete collection");
    }
  };

  const filtered = collections
    .filter(c => c.name?.toLowerCase().includes(search.toLowerCase() ?? ""))
    .sort((a, b) => {
      const aVal = (a[sortKey] || "").toString();
      const bVal = (b[sortKey] || "").toString();
      if (aVal < bVal) return sortAsc ? -1 : 1;
      if (aVal > bVal) return sortAsc ? 1 : -1;
      return 0;
    });

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <main className="flex-grow container mx-auto px-4 py-6">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-xl font-semibold">Collections</h1>
          <div className="flex items-center gap-2">
            <Input placeholder="New collection name" value={newName} onChange={e => setNewName(e.target.value)} className="w-48" />
            <Button onClick={handleCreate} disabled={creating} size="sm">
              {creating ? <Spinner size="sm" /> : <PlusCircle className="mr-1" />}
              Create
            </Button>
          </div>
        </div>
        <div className="flex justify-between items-center mb-4">
          <div className="flex items-center gap-2">
            <Input placeholder="Search..." value={search} onChange={e => setSearch(e.target.value)} className="w-48" />
            <Button variant="outline" size="sm" onClick={() => { setSortKey(k => k === "name" ? "last_updated" : k === "last_updated" ? "created_at" : "name"); setSortAsc(a => !a); }}>
              Sort by {sortKey.replace('_', ' ')} {sortAsc ? '↑' : '↓'}
            </Button>
          </div>
        </div>
        {loading ? (
          <div className="flex justify-center py-10"><Spinner /></div>
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
                {filtered.map(c => (
                  <TableRow key={c.id} className="hover:bg-gray-50">
                    <TableCell className="px-3 py-2">
                      <Link to={`/collections/${c.id}`} className="text-blue-600 hover:underline">
                        {c.name}
                      </Link>
                    </TableCell>
                    <TableCell className="px-3 py-2">{c.created_at ? new Date(c.created_at).toLocaleDateString() : '-'}</TableCell>
                    <TableCell className="px-3 py-2">{c.last_updated ? new Date(c.last_updated).toLocaleDateString() : '-'}</TableCell>
                    <TableCell className="px-3 py-2">
                      <AlertDialog>
                        <AlertDialogTrigger asChild>
                          <Button variant="destructive" size="sm"><Trash2 size={14} /></Button>
                        </AlertDialogTrigger>
                        <AlertDialogContent>
                          <AlertDialogHeader>
                            <AlertDialogTitle>Delete Collection?</AlertDialogTitle>
                            <AlertDialogDescription>
                              This action cannot be undone. Delete "{c.name}"?
                            </AlertDialogDescription>
                          </AlertDialogHeader>
                          <AlertDialogFooter>
                            <AlertDialogCancel>Cancel</AlertDialogCancel>
                            <AlertDialogAction onClick={() => handleDelete(c.id!)}>Delete</AlertDialogAction>
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
      <Footer />
    </div>
  );
}

