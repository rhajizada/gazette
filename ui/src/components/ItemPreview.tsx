import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Heart } from "lucide-react";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { toast } from "sonner";
import type { GithubComRhajizadaGazetteInternalServiceItem as ItemModel } from "../api/data-contracts";
import { useAuth } from "../context/AuthContext";

const MIN_WIDTH = 100;
const MIN_HEIGHT = 100;

interface ItemPreviewProps {
  item: ItemModel;
}

export function ItemPreview({ item }: ItemPreviewProps) {
  const { api, logout } = useAuth();
  const [liked, setLiked] = useState<boolean>(item.liked ?? false);
  const [loading, setLoading] = useState<boolean>(false);
  const [showImage, setShowImage] = useState(false);

  const toggleLike = async () => {
    if (loading) return;
    setLoading(true);
    try {
      if (liked) {
        await api.itemsLikeDelete(item.id!, { format: "json" });
        setLiked(false);
      } else {
        await api.itemsLikeCreate(item.id!, { format: "json" });
        setLiked(true);
      }
      toast.success(
        liked ? "removed from liked items" : "added to liked items",
      );
    } catch (err: any) {
      if (err.error === "Unauthorized") logout();
      else {
        const message = await err.text();
        toast.error(message || "failed to like item");
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!item?.image?.url) return;
    const img = new Image();
    img.src = item.image.url;
    img.onload = () => {
      if (img.naturalWidth > MIN_WIDTH && img.naturalHeight > MIN_HEIGHT) {
        setShowImage(true);
      }
    };
  }, [item?.image?.url]);

  return (
    <Card
      className="
       group
       flex flex-col
       h-full
       overflow-hidden
       rounded-2xl
       shadow hover:shadow-lg
       transition-shadow duration-300
       cursor-pointer
     "
    >
      <Link to={`/items/${item.id}`}>
        <CardContent className="flex items-start break-words flex-1">
          {showImage && (
            <img
              src={item.image.url}
              alt={item.image.title ?? item.title}
              className="
          w-1/3
          flex-shrink-0
          rounded-lg
          mr-6
          h-auto
          object-contain
        "
            />
          )}
          <div className="flex-1 break-words">
            <h3 className="text-lg font-semibold text-gray-900 hover:text-primary transition-colors">
              {item.title}
            </h3>

            {item.description && (
              <div
                className="leading-7 [&:not(:first-child)]:mt-6"
                dangerouslySetInnerHTML={{ __html: item.description }}
              />
            )}

            {(item.authors ?? []).filter(
              (a) => a.name?.trim() || a.email?.trim(),
            ).length > 0 && (
              <div className="flex flex-wrap gap-1 mt-2">
                {(item.authors ?? [])
                  .filter((a) => a.name?.trim() || a.email?.trim())
                  .map((a, i) => (
                    <Badge key={i}>{a.name || a.email}</Badge>
                  ))}
              </div>
            )}
          </div>
        </CardContent>
      </Link>

      <CardFooter className="mt-auto px-4 py-2 flex items-center justify-between">
        <button
          onClick={toggleLike}
          disabled={loading}
          aria-label={liked ? "Unlike item" : "Like item"}
          className={`p-2 rounded-full transition-colors ${
            liked
              ? "text-red-500 hover:bg-red-100"
              : "text-gray-500 hover:bg-gray-100"
          } ${loading ? "opacity-50 cursor-not-allowed" : ""}`}
        >
          <Heart
            fill={liked ? "currentColor" : "none"}
            strokeWidth={2}
            className="w-5 h-5"
          />
        </button>
        {item.published_parsed && (
          <span className="text-xs text-gray-500">
            {new Date(item.published_parsed).toLocaleDateString()}
          </span>
        )}
      </CardFooter>
    </Card>
  );
}
