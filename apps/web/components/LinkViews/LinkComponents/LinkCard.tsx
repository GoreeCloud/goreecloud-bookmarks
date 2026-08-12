import {
  ArchivedFormat,
  CollectionIncludingMembersAndLinkCount,
  LinkIncludingShortenedCollectionAndTags,
} from "@linkwarden/types/global";
import React, { useRef, useState } from "react";
import unescapeString from "@/lib/client/unescapeString";
import LinkActions from "@/components/LinkViews/LinkComponents/LinkActions";
import LinkDate from "@/components/LinkViews/LinkComponents/LinkDate";
import LinkCollection from "@/components/LinkViews/LinkComponents/LinkCollection";
import Image from "next/image";
import {
  atLeastOneFormatAvailable,
  formatAvailable,
} from "@linkwarden/lib/formatStats";
import LinkIcon from "./LinkIcon";
import toast from "react-hot-toast";
import LinkTypeBadge from "./LinkTypeBadge";
import useLocalSettingsStore from "@/store/localSettings";
import LinkPin from "./LinkPin";
import LinkFormats from "./LinkFormats";
import openLink from "@/lib/client/openLink";
import { useDraggable } from "@dnd-kit/core";
import { cn } from "@/lib/utils";
import { TFunction } from "i18next";

type Props = {
  link: LinkIncludingShortenedCollectionAndTags;
  collection: CollectionIncludingMembersAndLinkCount;
  isPublicRoute: boolean;
  t: TFunction<"translation", undefined>;
  user: any;
  disableDraggable: boolean;
  isSelected: boolean;
  toggleSelected: (id: number) => void;
  imageHeightClass: string;
  editMode?: boolean;
};

function LinkCard({
  link,
  collection,
  isPublicRoute,
  t,
  user,
  disableDraggable,
  isSelected,
  toggleSelected,
  imageHeightClass,
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

  const ref = useRef<HTMLDivElement>(null);
  const [linkModal, setLinkModal] = useState(false);

  return (
    <div
      ref={setNodeRef}
      className={cn(
        "group relative overflow-hidden rounded-2xl border bg-base-100 shadow-sm transition-all duration-200 touch-manipulation select-none",
        "border-base-content/10 hover:-translate-y-0.5 hover:border-base-content/20 hover:shadow-lg focus-within:border-primary/30 focus-within:shadow-md",
        isSelected && "border-primary/60 bg-primary/[0.035] ring-2 ring-primary/20",
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
      <div ref={ref} className="h-full">
        <div
          className="flex h-full cursor-pointer flex-col"
          onClick={() =>
            !editMode && openLink(link, user, () => setLinkModal(true))
          }
          {...listeners}
          {...attributes}
        >
          {show.image && (
            <div className="relative overflow-hidden border-b border-base-content/[0.07] bg-base-200/45">
              <div className={cn("relative overflow-hidden", imageHeightClass)}>
                {formatAvailable(link, "preview") ? (
                  <Image
                    src={`/api/v1/archives/${link.id}?format=${ArchivedFormat.jpeg}&preview=true&updatedAt=${link.updatedAt}`}
                    width={1280}
                    height={720}
                    alt=""
                    className={cn(
                      "z-10 h-full w-full select-none object-cover transition-transform duration-300 group-hover:scale-[1.015]",
                      imageHeightClass
                    )}
                    draggable="false"
                    onError={(e) => {
                      const target = e.target as HTMLElement;
                      target.style.display = "none";
                    }}
                    unoptimized
                  />
                ) : link.preview === "unavailable" ? (
                  <div
                    className={cn(
                      "flex items-center justify-center bg-base-200/70",
                      imageHeightClass
                    )}
                  >
                    <div className="opacity-70">
                      <LinkIcon link={link} />
                    </div>
                  </div>
                ) : (
                  <div className={cn("skeleton rounded-none", imageHeightClass)} />
                )}

                {show.icon && formatAvailable(link, "preview") && (
                  <div className="absolute inset-0 z-10 flex items-center justify-center bg-gradient-to-t from-black/20 via-transparent to-black/[0.03]">
                    <div className="rounded-2xl bg-base-100/90 p-1.5 shadow-lg ring-1 ring-white/30 backdrop-blur-md">
                      <LinkIcon link={link} />
                    </div>
                  </div>
                )}

                {show.preserved_formats &&
                  link.type === "url" &&
                  atLeastOneFormatAvailable(link) && (
                    <div className="absolute bottom-2 right-2 z-20 rounded-lg border border-white/15 bg-base-100/85 px-1.5 py-0.5 text-base-content/70 shadow-sm backdrop-blur-md">
                      <LinkFormats link={link} />
                    </div>
                  )}
              </div>
            </div>
          )}

          <div className="flex min-h-0 flex-1 flex-col">
            <div className="flex flex-1 flex-col gap-2.5 p-3.5 sm:p-4">
              {show.name && (
                <p className="line-clamp-2 w-full text-sm font-semibold leading-5 tracking-[-0.01em] text-base-content">
                  {unescapeString(link.name)}
                </p>
              )}

              {show.link && <LinkTypeBadge link={link} />}

              {show.description && link.description && (
                <p className="line-clamp-2 text-xs leading-5 text-base-content/55">
                  {unescapeString(link.description)}
                </p>
              )}
            </div>

            {(show.collection || show.date) && (
              <div className="border-t border-base-content/[0.07] px-3.5 py-2.5 sm:px-4">
                <div className="flex min-w-0 items-center justify-between gap-2 text-xs text-base-content/50">
                  {show.collection && !isPublicRoute && collection && (
                    <div className="min-w-0 truncate">
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
            )}
          </div>
        </div>

        <LinkActions
          link={link}
          linkModal={linkModal}
          t={t}
          setLinkModal={(e) => setLinkModal(e)}
          className="absolute right-2.5 top-2.5 z-30 h-10 w-10 rounded-lg border border-base-content/10 bg-base-100/90 text-base-content/60 opacity-100 shadow-sm backdrop-blur-md transition-all duration-150 hover:bg-base-100 hover:text-base-content sm:right-3 sm:top-3 sm:h-8 sm:w-8 sm:opacity-0 sm:group-hover:opacity-100 sm:group-focus-within:opacity-100"
        />
        {!isPublicRoute && <LinkPin link={link} />}
      </div>
    </div>
  );
}

export default React.memo(LinkCard);
