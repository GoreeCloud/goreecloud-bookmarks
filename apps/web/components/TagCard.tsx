import React, { useState } from "react";
import { useTranslation } from "next-i18next";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
} from "@/components/ui/dropdown-menu";
import { Button } from "./ui/button";
import { Checkbox } from "./ui/checkbox";
import { TagIncludingLinkCount } from "@linkwarden/types/global";
import DeleteTagModal from "./ModalContent/DeleteTagModal";
import { cn } from "@/lib/utils";
import { useRouter } from "next/router";

export default function TagCard({
  tag,
  editMode,
  selected,
  onSelect,
}: {
  tag: TagIncludingLinkCount;
  editMode: boolean;
  selected: boolean;
  onSelect: (tagId: number) => void;
}) {
  const { t } = useTranslation();
  const router = useRouter();
  const [deleteTagModal, setDeleteTagModal] = useState(false);

  const formattedDate = new Date(tag.createdAt).toLocaleString(t("locale"), {
    year: "numeric",
    month: "short",
    day: "numeric",
  });

  const activateTag = () => {
    if (editMode) onSelect(tag.id);
    else router.push(`/tags/${tag.id}`);
  };

  return (
    <div
      className={cn(
        "group relative min-h-[9.5rem] rounded-xl border bg-base-100 p-4 shadow-sm transition duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40",
        editMode ? "cursor-pointer" : "cursor-pointer hover:-translate-y-0.5 hover:shadow-md",
        selected
          ? "border-primary ring-2 ring-primary/15"
          : "border-neutral-content hover:border-primary/35"
      )}
      onClick={activateTag}
      role={editMode ? "button" : "link"}
      tabIndex={0}
      onKeyDown={(event) => {
        if (editMode && (event.key === "Enter" || event.key === " ")) {
          event.preventDefault();
          onSelect(tag.id);
        } else if (!editMode && event.key === "Enter") {
          event.preventDefault();
          router.push(`/tags/${tag.id}`);
        }
      }}
    >
      {editMode ? (
        <Checkbox
          checked={selected}
          className="absolute right-3 top-3 z-20 pointer-events-none"
          aria-label={tag.name}
        />
      ) : (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              type="button"
              variant="ghost"
              size="icon"
              className="absolute right-2 top-2 z-20 h-10 w-10 rounded-full opacity-80 hover:opacity-100 sm:right-2.5 sm:top-2.5 sm:h-8 sm:w-8 sm:opacity-65"
              onClick={(event) => event.stopPropagation()}
              aria-label={t("more")}
            >
              <i className="bi-three-dots text-lg text-neutral" aria-hidden="true" />
            </Button>
          </DropdownMenuTrigger>

          <DropdownMenuContent
            sideOffset={4}
            side="bottom"
            align="end"
            className="z-[30]"
            onClick={(event) => event.stopPropagation()}
          >
            <DropdownMenuItem
              onSelect={() => setDeleteTagModal(true)}
              className="text-error"
            >
              <i className="bi-trash" />
              {t("delete_tag")}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      )}

      <div className="flex h-full flex-col">
        <div className="flex min-w-0 items-start gap-3 pr-10 sm:pr-8">
          <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg border border-neutral-content bg-base-200 text-primary">
            <i className="bi-hash text-lg" aria-hidden="true" />
          </div>
          <div className="min-w-0 flex-1 pt-1">
            <h2
              className="truncate text-sm font-semibold leading-tight"
              title={tag.name}
            >
              {tag.name}
            </h2>
            <p className="mt-1 text-[11px] text-neutral">{formattedDate}</p>
          </div>
        </div>

        <div className="mt-auto flex items-end justify-between gap-3 pt-5">
          <span className="text-xs text-neutral">{t("links")}</span>
          <span className="inline-flex items-center gap-1.5 rounded-full border border-neutral-content bg-base-200 px-2.5 py-1 text-xs font-semibold tabular-nums">
            <i className="bi-link-45deg text-sm text-neutral" aria-hidden="true" />
            {tag._count?.links || 0}
          </span>
        </div>
      </div>

      {deleteTagModal && (
        <DeleteTagModal
          onClose={() => setDeleteTagModal(false)}
          activeTag={tag}
        />
      )}
    </div>
  );
}
