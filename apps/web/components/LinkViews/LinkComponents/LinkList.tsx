import {
  CollectionIncludingMembersAndLinkCount,
  LinkIncludingShortenedCollectionAndTags,
} from "@linkwarden/types/global";
import React, { useState } from "react";
import unescapeString from "@/lib/client/unescapeString";
import LinkActions from "@/components/LinkViews/LinkComponents/LinkActions";
import LinkDate from "@/components/LinkViews/LinkComponents/LinkDate";
import LinkCollection from "@/components/LinkViews/LinkComponents/LinkCollection";
import LinkIcon from "@/components/LinkViews/LinkComponents/LinkIcon";
import { cn, isPWA } from "@/lib/utils";
import toast from "react-hot-toast";
import LinkTypeBadge from "./LinkTypeBadge";
import useLocalSettingsStore from "@/store/localSettings";
import LinkPin from "./LinkPin";
import { atLeastOneFormatAvailable } from "@linkwarden/lib/formatStats";
import LinkFormats from "./LinkFormats";
import openLink from "@/lib/client/openLink";
import { useDraggable } from "@dnd-kit/core";
import { TFunction } from "i18next";

type Props = {
  link: LinkIncludingShortenedCollectionAndTags;
  collection: CollectionIncludingMembersAndLinkCount;
  isPublicRoute: boolean;
  t: TFunction<"translation", undefined>;
  disableDraggable: boolean;
  user: any;
  isSelected: boolean;
  toggleSelected: (id: number) => void;
  count: number;
  className?: string;
  editMode?: boolean;
};

function LinkList({
  link,
  collection,
  isPublicRoute,
  t,
  disableDraggable,
  user,
  isSelected,
  toggleSelected,
  editMode,
}: Props) {
  const { attributes, listeners, setNodeRef, isDragging } = useDraggable({
    id: link.id?.toString() ?? "",
    data: {
      linkId: link.id,
      link,
    },
    disabled: disableDraggable,
  });

  const {
    settings: { show },
  } = useLocalSettingsStore();

  const [linkModal, setLinkModal] = useState(false);

  return (
    <div
      ref={setNodeRef}
      className={cn(
        "group relative mb-2 flex min-h-[4.5rem] items-center overflow-hidden rounded-xl border bg-base-100 shadow-sm transition-all duration-150 touch-manipulation select-none",
        isSelected
          ? "border-primary/60 bg-primary/[0.035] ring-2 ring-primary/15"
          : "border-base-content/[0.08] hover:border-base-content/15 hover:bg-base-200/35 hover:shadow-md",
        !isPWA() ? "px-3 py-2.5" : "px-2 py-2.5",
        isDragging ? "opacity-30" : "opacity-100"
      )}
      onClick={() =>
        editMode
          ? toggleSelected(link.id as number)
          : editMode
            ? toast.error(t("link_selection_error"))
            : undefined
      }
    >
      <div
        className="flex min-w-0 flex-1 cursor-pointer items-center gap-3 pr-24 sm:pr-20"
        onClick={() =>
          !editMode && openLink(link, user, () => setLinkModal(true))
        }
        {...attributes}
        {...listeners}
      >
        {show.icon && (
          <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl border border-base-content/[0.07] bg-base-200/55 shadow-sm">
            <LinkIcon link={link} hideBackground />
          </div>
        )}

        <div className="min-w-0 flex-1">
          <div className="flex min-w-0 items-center gap-2">
            {show.name && (
              <p className="truncate text-sm font-semibold tracking-[-0.01em] text-base-content sm:text-[15px]">
                {unescapeString(link.name)}
              </p>
            )}

            {show.preserved_formats &&
              link.type === "url" &&
              atLeastOneFormatAvailable(link) && (
                <div className="hidden shrink-0 rounded-md bg-base-content/[0.045] px-1.5 py-0.5 text-base-content/50 sm:block">
                  <LinkFormats link={link} />
                </div>
              )}
          </div>

          <div className="mt-1.5 flex min-w-0 flex-wrap items-center gap-x-2.5 gap-y-1 text-xs text-base-content/50">
            {show.link && <LinkTypeBadge link={link} />}
            {show.collection && collection && !isPublicRoute && (
              <div className="min-w-0 max-w-48 truncate">
                <LinkCollection
                  link={link}
                  collection={collection}
                  isPublicRoute={isPublicRoute}
                />
              </div>
            )}
            {show.date && (
              <div className="shrink-0">
                <LinkDate link={link} />
              </div>
            )}
          </div>
        </div>
      </div>

      {!isPublicRoute && <LinkPin link={link} />}
      <LinkActions
        link={link}
        linkModal={linkModal}
        t={t}
        setLinkModal={(e) => setLinkModal(e)}
        className="absolute right-2.5 top-2.5 z-20 h-10 w-10 rounded-lg border border-base-content/10 bg-base-100/90 text-base-content/60 opacity-100 shadow-sm backdrop-blur-md transition-all duration-150 hover:bg-base-100 hover:text-base-content sm:right-3 sm:top-3 sm:h-8 sm:w-8 sm:opacity-0 sm:group-hover:opacity-100 sm:group-focus-within:opacity-100"
      />
    </div>
  );
}

export default React.memo(LinkList);
