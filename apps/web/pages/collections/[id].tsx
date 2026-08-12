import {
  AccountSettings,
  CollectionIncludingMembersAndLinkCount,
  Sort,
  ViewMode,
} from "@linkwarden/types/global";
import { useRouter } from "next/router";
import React, { ReactElement, useEffect, useState } from "react";
import MainLayout from "@/layouts/MainLayout";
import ProfilePhoto from "@/components/ProfilePhoto";
import usePermissions from "@/hooks/usePermissions";
import NoLinksFound from "@/components/NoLinksFound";
import getPublicUserData from "@/lib/client/getPublicUserData";
import EditCollectionModal from "@/components/ModalContent/EditCollectionModal";
import EditCollectionSharingModal from "@/components/ModalContent/EditCollectionSharingModal";
import DeleteCollectionModal from "@/components/ModalContent/DeleteCollectionModal";
import NewCollectionModal from "@/components/ModalContent/NewCollectionModal";
import getServerSideProps from "@/lib/client/getServerSideProps";
import { useTranslation } from "next-i18next";
import LinkListOptions from "@/components/LinkListOptions";
import { useCollections } from "@linkwarden/router/collections";
import { useUser } from "@linkwarden/router/user";
import { useLinks } from "@linkwarden/router/links";
import Links from "@/components/LinkViews/Links";
import Icon from "@/components/Icon";
import CollectionCard from "@/components/CollectionCard";
import { IconWeight } from "@phosphor-icons/react";
import PageHeader from "@/components/PageHeader";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { NextPageWithLayout } from "../_app";

const Page: NextPageWithLayout = () => {
  const { t } = useTranslation();
  const router = useRouter();
  const { data: collections = [] } = useCollections();
  const { data: user } = useUser();

  const [sortBy, setSortBy] = useState<Sort>(
    Number(localStorage.getItem("sortBy")) ?? Sort.DateNewestFirst
  );
  const [viewMode, setViewMode] = useState<ViewMode>(
    (localStorage.getItem("viewMode") as ViewMode) || ViewMode.Card
  );
  const [activeCollection, setActiveCollection] =
    useState<CollectionIncludingMembersAndLinkCount>();
  const [collectionOwner, setCollectionOwner] = useState<
    Partial<AccountSettings>
  >({});
  const [editCollectionModal, setEditCollectionModal] = useState(false);
  const [newCollectionModal, setNewCollectionModal] = useState(false);
  const [editCollectionSharingModal, setEditCollectionSharingModal] =
    useState(false);
  const [deleteCollectionModal, setDeleteCollectionModal] = useState(false);
  const [editMode, setEditMode] = useState(false);

  const { links, data } = useLinks({
    sort: sortBy,
    collectionId: Number(router.query.id),
  });

  const permissions = usePermissions(activeCollection?.id as number);

  useEffect(() => {
    setActiveCollection(
      collections.find((collection) => collection.id === Number(router.query.id))
    );
  }, [router, collections]);

  useEffect(() => {
    const fetchOwner = async () => {
      if (activeCollection && activeCollection.ownerId !== user?.id) {
        const owner = await getPublicUserData(activeCollection.ownerId as number);
        setCollectionOwner(owner);
      } else if (activeCollection && activeCollection.ownerId === user?.id) {
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
  }, [activeCollection, user]);

  useEffect(() => {
    if (editMode) setEditMode(false);
  }, [router]);

  const subcollections = collections.filter(
    (collection) => collection.parentId === activeCollection?.id
  );
  const members = activeCollection
    ? [...activeCollection.members].sort(
        (a, b) => (a.userId as number) - (b.userId as number)
      )
    : [];
  const visibleMembers = members.slice(0, 3);
  const extraMembers = Math.max(0, members.length - visibleMembers.length);

  return (
    <div className="flex w-full flex-col gap-5 p-3 sm:p-5">
      {activeCollection && (
        <section className="overflow-hidden rounded-2xl border border-neutral-content bg-base-100 shadow-sm">
          <div className="h-1.5 w-full" style={{ backgroundColor: activeCollection.color }} />
          <div className="p-4 sm:p-5">
            <div className="flex items-start justify-between gap-4">
              <div className="flex min-w-0 items-start gap-3">
                <div
                  className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl border border-neutral-content bg-base-200"
                  style={{ color: activeCollection.color }}
                >
                  {activeCollection.icon ? (
                    <Icon
                      icon={activeCollection.icon}
                      size={29}
                      weight={
                        (activeCollection.iconWeight || "regular") as IconWeight
                      }
                      color={activeCollection.color}
                    />
                  ) : (
                    <i className="bi-folder-fill text-2xl" />
                  )}
                </div>

                <div className="min-w-0 pt-0.5">
                  <div className="flex min-w-0 flex-wrap items-center gap-2">
                    <h1 className="break-words text-2xl font-semibold sm:text-3xl">
                      {activeCollection.name}
                    </h1>
                    {activeCollection.isPublic && (
                      <span className="inline-flex items-center gap-1 rounded-full border border-neutral-content bg-base-200 px-2 py-0.5 text-[10px] font-medium text-neutral">
                        <i className="bi-globe2" />
                        Public
                      </span>
                    )}
                  </div>
                  {activeCollection.description && (
                    <p className="mt-2 max-w-3xl text-sm leading-relaxed text-neutral">
                      {activeCollection.description}
                    </p>
                  )}
                </div>
              </div>

              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="shrink-0 rounded-full"
                    title={t("more")}
                  >
                    <i className="bi-three-dots text-lg text-neutral" />
                  </Button>
                </DropdownMenuTrigger>

                <DropdownMenuContent sideOffset={4} align="end">
                  <DropdownMenuItem
                    onClick={() => {
                      let blocked = false;

                      for (const link of links) {
                        if (!link.url) continue;
                        const openedWindow = window.open(link.url, "_blank");
                        if (openedWindow) openedWindow.blur();
                        else blocked = true;
                      }

                      window.focus();
                      if (blocked) alert(t("popups_blocked_open_all_links"));
                    }}
                  >
                    <i className="bi-box-arrow-up-right" />
                    {t("open_all_links")}
                  </DropdownMenuItem>

                  {permissions === true && (
                    <DropdownMenuItem onClick={() => setEditCollectionModal(true)}>
                      <i className="bi-pencil-square" />
                      {t("edit_collection_info")}
                    </DropdownMenuItem>
                  )}

                  <DropdownMenuItem
                    onClick={() => setEditCollectionSharingModal(true)}
                  >
                    <i className="bi-people" />
                    {permissions === true
                      ? t("share_and_collaborate")
                      : t("view_team")}
                  </DropdownMenuItem>

                  {(permissions === true ||
                    (permissions?.canCreate &&
                      permissions?.canUpdate &&
                      permissions?.canDelete)) && (
                    <DropdownMenuItem onClick={() => setNewCollectionModal(true)}>
                      <i className="bi-folder-plus" />
                      {t("create_subcollection")}
                    </DropdownMenuItem>
                  )}

                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    onClick={() => setDeleteCollectionModal(true)}
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
            </div>

            <div className="mt-5 flex flex-col gap-3 border-t border-neutral-content pt-4 sm:flex-row sm:items-center sm:justify-between">
              <button
                type="button"
                className="flex w-fit items-center gap-3 rounded-lg text-left focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40"
                onClick={() => setEditCollectionSharingModal(true)}
              >
                <div className="flex items-center">
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
                </div>
                <span className="text-xs text-neutral">
                  {activeCollection.members.length > 0
                    ? activeCollection.members.length === 1
                      ? t("by_author_and_other", {
                          author: collectionOwner.name,
                          count: activeCollection.members.length,
                        })
                      : t("by_author_and_others", {
                          author: collectionOwner.name,
                          count: activeCollection.members.length,
                        })
                    : t("by_author", { author: collectionOwner.name })}
                </span>
              </button>

              <div className="flex items-center gap-4 text-xs text-neutral">
                <span className="inline-flex items-center gap-1.5">
                  <i className="bi-folder2" />
                  <span className="font-semibold tabular-nums text-base-content">
                    {subcollections.length}
                  </span>
                </span>
                <span className="inline-flex items-center gap-1.5">
                  <i className="bi-link-45deg" />
                  <span className="font-semibold tabular-nums text-base-content">
                    {activeCollection._count?.links || 0}
                  </span>
                </span>
              </div>
            </div>
          </div>
        </section>
      )}

      {subcollections.length > 0 && (
        <section className="flex flex-col gap-3 rounded-2xl border border-neutral-content bg-base-200/30 p-4 sm:p-5">
          <PageHeader
            icon="bi-folder"
            title={t("collections")}
            description={
              subcollections.length === 1
                ? t("showing_count_result", { count: subcollections.length })
                : t("showing_count_results", { count: subcollections.length })
            }
            sm
          />
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4">
            {subcollections.map((collection) => (
              <CollectionCard key={collection.id} collection={collection} />
            ))}
          </div>
        </section>
      )}

      <LinkListOptions
        t={t}
        viewMode={viewMode}
        setViewMode={setViewMode}
        sortBy={sortBy}
        setSortBy={setSortBy}
        editMode={
          permissions === true || permissions?.canUpdate || permissions?.canDelete
            ? editMode
            : undefined
        }
        setEditMode={
          permissions === true || permissions?.canUpdate || permissions?.canDelete
            ? setEditMode
            : undefined
        }
        links={links}
      >
        <PageHeader
          icon="bi-link-45deg"
          title={t("links")}
          description={
            activeCollection?._count?.links === 1
              ? t("showing_count_result", {
                  count: activeCollection?._count?.links,
                })
              : t("showing_count_results", {
                  count: activeCollection?._count?.links || 0,
                })
          }
          sm={subcollections.length > 0}
        />
      </LinkListOptions>

      <Links
        editMode={editMode}
        links={links}
        layout={viewMode}
        useData={data}
      />
      {!data.isLoading && links.length === 0 && <NoLinksFound />}

      {activeCollection && (
        <>
          {editCollectionModal && (
            <EditCollectionModal
              onClose={() => setEditCollectionModal(false)}
              activeCollection={activeCollection}
            />
          )}
          {editCollectionSharingModal && (
            <EditCollectionSharingModal
              onClose={() => setEditCollectionSharingModal(false)}
              activeCollection={activeCollection}
            />
          )}
          {newCollectionModal && (
            <NewCollectionModal
              onClose={() => setNewCollectionModal(false)}
              parent={activeCollection}
            />
          )}
          {deleteCollectionModal && (
            <DeleteCollectionModal
              onClose={() => setDeleteCollectionModal(false)}
              activeCollection={activeCollection}
            />
          )}
        </>
      )}
    </div>
  );
};

Page.getLayout = function getLayout(page: ReactElement<any>) {
  return <MainLayout>{page}</MainLayout>;
};

export default Page;

export { getServerSideProps };
