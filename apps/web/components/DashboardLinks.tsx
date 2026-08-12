import { LinkIncludingShortenedCollectionAndTags } from "@linkwarden/types/global";
import useLocalSettingsStore from "@/store/localSettings";
import {
  ArchivedFormat,
  CollectionIncludingMembersAndLinkCount,
} from "@linkwarden/types/global";
import { useEffect, useState } from "react";
import unescapeString from "@/lib/client/unescapeString";
import LinkActions from "@/components/LinkViews/LinkComponents/LinkActions";
import LinkDate from "@/components/LinkViews/LinkComponents/LinkDate";
import LinkCollection from "@/components/LinkViews/LinkComponents/LinkCollection";
import Image from "next/image";
import {
  atLeastOneFormatAvailable,
  formatAvailable,
} from "@linkwarden/lib/formatStats";
import { useCollections } from "@linkwarden/router/collections";
import { useUser } from "@linkwarden/router/user";
import { useRouter } from "next/router";
import openLink from "@/lib/client/openLink";
import LinkIcon from "./LinkViews/LinkComponents/LinkIcon";
import LinkFormats from "./LinkViews/LinkComponents/LinkFormats";
import LinkTypeBadge from "./LinkViews/LinkComponents/LinkTypeBadge";
import LinkPin from "./LinkViews/LinkComponents/LinkPin";
import { useDraggable } from "@dnd-kit/core";
import { cn } from "@linkwarden/lib/utils";
import { useTranslation } from "next-i18next";

export function DashboardLinks({
  links,
  isLoading,
  type,
}: {
  links?: LinkIncludingShortenedCollectionAndTags[];
  isLoading?: boolean;
  type?: "collection" | "recent";
}) {
  return (
    <div className="flex w-full gap-3 overflow-x-auto overflow-y-hidden pb-1 hide-scrollbar">
      {isLoading ? (
        <div className="min-w-60 w-60 overflow-hidden rounded-2xl border border-base-content/[0.07] bg-base-100">
          <div className="skeleton h-40 w-full rounded-none" />
          <div className="space-y-2 p-4">
            <div className="skeleton h-4 w-2/3" />
            <div className="skeleton h-3 w-full" />
            <div className="skeleton h-3 w-1/2" />
          </div>
        </div>
      ) : (
        links?.map((e) => <Card key={e.id} link={e} dashboardType={type} />)
      )}
    </div>
  );
}

type Props = {
  link: LinkIncludingShortenedCollectionAndTags;
  editMode?: boolean;
  dashboardType?: "collection" | "recent";
};

export function Card({ link, editMode, dashboardType }: Props) {
  const { t } = useTranslation();

  const { attributes, listeners, setNodeRef, isDragging } = useDraggable({
    id: `${link.id}-${dashboardType}`,
    data: {
      linkId: link.id,
      link,
      dashboardType,
    },
  });
  const { data: collections = [] } = useCollections();
  const { data: user } = useUser();

  const {
    settings: { show },
  } = useLocalSettingsStore();

  const router = useRouter();
  const isPublicRoute = router.pathname.startsWith("/public");

  const [collection, setCollection] =
    useState<CollectionIncludingMembersAndLinkCount>(
      collections.find(
        (e) => e.id === link.collection.id
      ) as CollectionIncludingMembersAndLinkCount
    );

  useEffect(() => {
    setCollection(
      collections.find(
        (e) => e.id === link.collection.id
      ) as CollectionIncludingMembersAndLinkCount
    );
  }, [collections, link]);

  const [linkModal, setLinkModal] = useState(false);

  return (
    <div
      ref={setNodeRef}
      className={cn(
        "group relative min-w-60 w-60 overflow-hidden rounded-2xl border border-base-content/10 bg-base-100 shadow-sm transition-all duration-200 touch-manipulation select-none",
        "hover:-translate-y-0.5 hover:border-base-content/20 hover:shadow-lg focus-within:border-primary/30 focus-within:shadow-md",
        isDragging ? "opacity-30" : "opacity-100"
      )}
    >
      <div
        className="flex h-full cursor-pointer flex-col"
        onClick={() =>
          !editMode && openLink(link, user, () => setLinkModal(true))
        }
        {...listeners}
        {...attributes}
      >
        {show.image && (
          <div className="relative h-40 overflow-hidden border-b border-base-content/[0.07] bg-base-200/45">
            {formatAvailable(link, "preview") ? (
              <Image
                src={`/api/v1/archives/${link.id}?format=${ArchivedFormat.jpeg}&preview=true&updatedAt=${link.updatedAt}`}
                width={1280}
                height={720}
                alt=""
                className="z-10 h-40 w-full select-none object-cover transition-transform duration-300 group-hover:scale-[1.015]"
                draggable="false"
                onError={(e) => {
                  const target = e.target as HTMLElement;
                  target.style.display = "none";
                }}
                unoptimized
              />
            ) : link.preview === "unavailable" ? (
              <div className="flex h-40 items-center justify-center bg-base-200/70">
                <div className="opacity-70">
                  <LinkIcon link={link} />
                </div>
              </div>
            ) : (
              <div className="skeleton h-40 rounded-none" />
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
        )}

        <div className="flex min-h-0 flex-1 flex-col">
          <div className="flex flex-1 flex-col gap-2.5 p-3.5">
            {show.name && (
              <p className="line-clamp-2 text-sm font-semibold leading-5 tracking-[-0.01em] text-base-content">
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
            <div className="border-t border-base-content/[0.07] px-3.5 py-2.5">
              <div className="flex min-w-0 items-center justify-between gap-2 text-xs text-base-content/50">
                {show.collection && !isPublicRoute && collection && (
                  <div className="min-w-0 truncate">
                    <LinkCollection
                      link={link}
                      collection={collection}
                      isPublicRoute={false}
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
        t={t}
        linkModal={linkModal}
        setLinkModal={(e) => setLinkModal(e)}
        className="absolute top-3 right-3 z-30 h-8 w-8 rounded-lg border border-base-content/10 bg-base-100/90 text-base-content/60 opacity-0 shadow-sm backdrop-blur-md transition-all duration-150 hover:bg-base-100 hover:text-base-content group-hover:opacity-100 group-focus-within:opacity-100"
      />
      {!isPublicRoute && <LinkPin link={link} />}
    </div>
  );
}
