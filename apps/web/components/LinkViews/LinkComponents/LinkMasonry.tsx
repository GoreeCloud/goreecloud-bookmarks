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
import Link from "next/link";
import LinkIcon from "./LinkIcon";
import toast from "react-hot-toast";
import LinkTypeBadge from "./LinkTypeBadge";
import useLocalSettingsStore from "@/store/localSettings";
import clsx from "clsx";
import LinkPin from "./LinkPin";
import LinkFormats from "./LinkFormats";
import openLink from "@/lib/client/openLink";
import { useDraggable } from "@dnd-kit/core";
import { cn } from "@linkwarden/lib/utils";
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
  imageHeightClass: string;
  editMode?: boolean;
};

function LinkMasonry({
  link,
  collection,
  isPublicRoute,
  t,
  disableDraggable,
  user,
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
      <div ref={ref}>
        <div
          className="cursor-pointer"
          onClick={() =>
            !editMode && openLink(link, user, () => setLinkModal(true))
          }
          {...listeners}
          {...attributes}
        >
          {show.image && formatAvailable(link, "preview") && (
            <div className="relative overflow-hidden border-b border-base-content/[0.07] bg-base-200/45">
              <Image
                src={`/api/v1/archives/${link.id}?format=${ArchivedFormat.jpeg}&preview=true&updatedAt=${link.updatedAt}`}
                width={1280}
                height={720}
                alt=""
                className={cn(
                  "z-10 w-full select-none object-cover transition-transform duration-300 group-hover:scale-[1.015]",
                  imageHeightClass
                )}
                draggable="false"
                onError={(e) => {
                  const target = e.target as HTMLElement;
                  target.style.display = "none";
                }}
                unoptimized
              />

              {show.icon && (
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
          )}

          <div className="flex flex-col gap-2.5 p-3.5 sm:p-4">
            {show.name && (
              <p className="hyphens-auto text-sm font-semibold leading-5 tracking-[-0.01em] text-base-content">
                {unescapeString(link.name)}
              </p>
            )}

            <div className="flex flex-wrap items-center gap-2">
              {show.link && <LinkTypeBadge link={link} />}
              {!show.image &&
                show.preserved_formats &&
                link.type === "url" &&
                atLeastOneFormatAvailable(link) && (
                  <div className="rounded-md bg-base-content/[0.045] px-1.5 py-0.5 text-base-content/50">
                    <LinkFormats link={link} />
                  </div>
                )}
            </div>

            {show.description && link.description && (
              <p className="hyphens-auto line-clamp-4 text-xs leading-5 text-base-content/55">
                {unescapeString(link.description)}
              </p>
            )}

            {show.tags && link.tags && link.tags[0] && (
              <div className="flex flex-wrap items-center gap-1.5 pt-0.5">
                {link.tags.map((e) => (
                  <Link
                    key={e.id}
                    href={`/tags/${e.id}`}
                    onClick={(event) => {
                      event.stopPropagation();
                    }}
                    className="max-w-full truncate rounded-md bg-primary/[0.07] px-2 py-1 text-[11px] font-medium text-primary/80 transition-colors hover:bg-primary/12 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40"
                  >
                    #{e.name}
                  </Link>
                ))}
              </div>
            )}
          </div>

          {(show.collection || show.date) && (
            <div className="border-t border-base-content/[0.07] px-3.5 py-2.5 sm:px-4">
              <div className="flex min-w-0 flex-wrap items-center justify-between gap-x-2 gap-y-1 text-xs text-base-content/50">
                {!isPublicRoute && show.collection && collection && (
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

        <LinkActions
          link={link}
          linkModal={linkModal}
          t={t}
          setLinkModal={(e) => setLinkModal(e)}
          className="absolute top-3 right-3 z-30 h-8 w-8 rounded-lg border border-base-content/10 bg-base-100/90 text-base-content/60 opacity-0 shadow-sm backdrop-blur-md transition-all duration-150 hover:bg-base-100 hover:text-base-content group-hover:opacity-100 group-focus-within:opacity-100"
        />
        {!isPublicRoute && <LinkPin link={link} />}
      </div>
    </div>
  );
}

export default React.memo(LinkMasonry);
