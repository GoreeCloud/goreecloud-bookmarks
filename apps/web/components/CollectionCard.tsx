import Link from "next/link";
import {
  AccountSettings,
  CollectionIncludingMembersAndLinkCount,
} from "@linkwarden/types/global";
import React, { useEffect, useState } from "react";
import ProfilePhoto from "./ProfilePhoto";
import usePermissions from "@/hooks/usePermissions";
import getPublicUserData from "@/lib/client/getPublicUserData";
import EditCollectionModal from "./ModalContent/EditCollectionModal";
import EditCollectionSharingModal from "./ModalContent/EditCollectionSharingModal";
import DeleteCollectionModal from "./ModalContent/DeleteCollectionModal";
import { useTranslation } from "next-i18next";
import { useUser } from "@linkwarden/router/user";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu";
import { Button } from "./ui/button";
import Icon from "./Icon";
import { IconWeight } from "@phosphor-icons/react";

export default function CollectionCard({
  collection,
}: {
  collection: CollectionIncludingMembersAndLinkCount;
}) {
  const { t } = useTranslation();
  const { data: user } = useUser();

  const formattedDate = new Date(collection.createdAt as string).toLocaleString(
    t("locale"),
    {
      year: "numeric",
      month: "short",
      day: "numeric",
    }
  );

  const permissions = usePermissions(collection.id as number);

  const [collectionOwner, setCollectionOwner] = useState<
    Partial<AccountSettings>
  >({});

  useEffect(() => {
    const fetchOwner = async () => {
      if (collection.ownerId !== user?.id) {
        const owner = await getPublicUserData(collection.ownerId as number);
        setCollectionOwner(owner);
      } else if (collection.ownerId === user?.id) {
        setCollectionOwner({
          id: user?.id as number,
          name: user?.name,
          username: user?.username as string,
          image: user?.image as string,
          archiveAsScreenshot: user?.archiveAsScreenshot as boolean,
          archiveAsMonolith: user?.archiveAsMonolith as boolean,
          archiveAsPDF: user?.archiveAsPDF as boolean,
        });
      }
    };

    fetchOwner();
  }, [collection, user]);

  const [editCollectionModal, setEditCollectionModal] = useState(false);
  const [editCollectionSharingModal, setEditCollectionSharingModal] =
    useState(false);
  const [deleteCollectionModal, setDeleteCollectionModal] = useState(false);

  const members = [...collection.members].sort(
    (a, b) => (a.userId as number) - (b.userId as number)
  );
  const visibleMembers = members.slice(0, 3);
  const extraMembers = Math.max(0, members.length - visibleMembers.length);

  return (
    <div className="relative group">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            type="button"
            variant="ghost"
            size="icon"
            className="absolute top-3 right-3 z-20 h-8 w-8 rounded-full bg-base-100/80 opacity-70 shadow-sm backdrop-blur hover:opacity-100"
            aria-label={t("more")}
          >
            <i className="bi-three-dots text-lg text-neutral" />
          </Button>
        </DropdownMenuTrigger>

        <DropdownMenuContent sideOffset={4} side="bottom" align="end">
          {permissions === true && (
            <DropdownMenuItem onSelect={() => setEditCollectionModal(true)}>
              <i className="bi-pencil-square" />
              {t("edit_collection_info")}
            </DropdownMenuItem>
          )}

          <DropdownMenuItem
            onSelect={() => setEditCollectionSharingModal(true)}
          >
            <i className="bi-people" />
            {permissions === true ? t("share_and_collaborate") : t("view_team")}
          </DropdownMenuItem>

          <DropdownMenuSeparator />

          <DropdownMenuItem
            onSelect={() => setDeleteCollectionModal(true)}
            className="text-error"
          >
            {permissions === true ? (
              <>
                <i className="bi-trash" />
                {t("delete_collection")}
              </>
            ) : (
              <>
                <i className="bi-box-arrow-left" />
                {t("leave_collection")}
              </>
            )}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <Link
        href={`/collections/${collection.id}`}
        className="flex min-h-[12rem] flex-col rounded-xl border border-neutral-content bg-base-100 p-4 pr-12 shadow-sm transition duration-150 hover:-translate-y-0.5 hover:border-primary/40 hover:shadow-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40"
        style={{ borderTopColor: collection.color, borderTopWidth: "3px" }}
      >
        <div className="flex min-w-0 items-start gap-3">
          <div
            className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl border border-neutral-content bg-base-200"
            style={{ color: collection.color }}
          >
            {collection.icon ? (
              <Icon
                icon={collection.icon}
                size={25}
                weight={(collection.iconWeight || "regular") as IconWeight}
                color={collection.color}
              />
            ) : (
              <i className="bi-folder-fill text-xl" />
            )}
          </div>

          <div className="min-w-0 flex-1 pt-0.5">
            <div className="flex min-w-0 items-center gap-2">
              <h2 className="truncate text-base font-semibold" title={collection.name}>
                {collection.name}
              </h2>
              {collection.isPublic && (
                <span
                  className="inline-flex shrink-0 items-center gap-1 rounded-full border border-neutral-content bg-base-200 px-2 py-0.5 text-[10px] font-medium text-neutral"
                  title={t("collection_publicly_shared")}
                >
                  <i className="bi-globe2" />
                  Public
                </span>
              )}
            </div>
            {collection.description && (
              <p className="mt-1 line-clamp-2 text-xs leading-relaxed text-neutral">
                {collection.description}
              </p>
            )}
          </div>
        </div>

        <div className="mt-auto flex items-end justify-between gap-3 pt-5">
          <button
            type="button"
            className="flex items-center rounded-full focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40"
            onClick={(e) => {
              e.preventDefault();
              e.stopPropagation();
              setEditCollectionSharingModal(true);
            }}
            aria-label={permissions === true ? t("share_and_collaborate") : t("view_team")}
          >
            {collectionOwner.id && (
              <ProfilePhoto
                src={collectionOwner.image || undefined}
                name={collectionOwner.name}
              />
            )}
            {visibleMembers.map((member, index) => (
              <ProfilePhoto
                key={member.userId || index}
                src={member.user.image || undefined}
                name={member.user.name}
                className="-ml-2.5"
              />
            ))}
            {extraMembers > 0 && (
              <div className="-ml-2.5 flex h-8 w-8 items-center justify-center rounded-full border-2 border-base-100 bg-base-200 text-[10px] font-semibold text-neutral">
                +{extraMembers}
              </div>
            )}
          </button>

          <div className="flex items-center gap-3 text-xs text-neutral">
            <span className="inline-flex items-center gap-1.5">
              <i className="bi-link-45deg text-sm" />
              <span className="font-medium tabular-nums text-base-content">
                {collection._count?.links || 0}
              </span>
            </span>
            <span className="hidden items-center gap-1.5 sm:inline-flex">
              <i className="bi-calendar3" />
              {formattedDate}
            </span>
          </div>
        </div>
      </Link>

      {editCollectionModal && (
        <EditCollectionModal
          onClose={() => setEditCollectionModal(false)}
          activeCollection={collection}
        />
      )}
      {editCollectionSharingModal && (
        <EditCollectionSharingModal
          onClose={() => setEditCollectionSharingModal(false)}
          activeCollection={collection}
        />
      )}
      {deleteCollectionModal && (
        <DeleteCollectionModal
          onClose={() => setDeleteCollectionModal(false)}
          activeCollection={collection}
        />
      )}
    </div>
  );
}
