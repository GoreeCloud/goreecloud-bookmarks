import CollectionCard from "@/components/CollectionCard";
import { ReactElement, useMemo, useState } from "react";
import MainLayout from "@/layouts/MainLayout";
import { useSession } from "next-auth/react";
import SortDropdown from "@/components/SortDropdown";
import { Sort } from "@linkwarden/types/global";
import NewCollectionModal from "@/components/ModalContent/NewCollectionModal";
import PageHeader from "@/components/PageHeader";
import getServerSideProps from "@/lib/client/getServerSideProps";
import { useTranslation } from "next-i18next";
import { useCollections } from "@linkwarden/router/collections";
import { Button } from "@/components/ui/button";
import { NextPageWithLayout } from "../_app";

const Page: NextPageWithLayout = () => {
  const { t } = useTranslation();
  const { data: collections = [], isLoading } = useCollections();
  const [sortBy, setSortBy] = useState<Sort>(Sort.DateNewestFirst);
  const [newCollectionModal, setNewCollectionModal] = useState(false);
  const { data } = useSession();

  const sortKey: Sort =
    typeof sortBy === "string" ? (Number(sortBy) as Sort) : sortBy;

  const compare = useMemo(() => {
    switch (sortKey) {
      case Sort.NameAZ:
        return (a: any, b: any) => a.name.localeCompare(b.name);
      case Sort.NameZA:
        return (a: any, b: any) => b.name.localeCompare(a.name);
      case Sort.DateOldestFirst:
        return (a: any, b: any) =>
          new Date(a.createdAt as string).getTime() -
          new Date(b.createdAt as string).getTime();
      case Sort.DateNewestFirst:
      default:
        return (a: any, b: any) =>
          new Date(b.createdAt as string).getTime() -
          new Date(a.createdAt as string).getTime();
    }
  }, [sortKey]);

  const sortedCollections = useMemo(
    () => [...collections].sort(compare),
    [collections, compare]
  );

  const ownedCollections = sortedCollections.filter(
    (collection) =>
      collection.ownerId === data?.user.id && collection.parentId === null
  );
  const sharedCollections = sortedCollections.filter(
    (collection) => collection.ownerId !== data?.user.id
  );

  return (
    <div className="flex h-full w-full flex-col gap-6 p-3 sm:p-5">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <PageHeader
          icon="bi-folder"
          title={t("collections")}
          description={t("collections_you_own")}
        />

        <div className="flex items-center gap-2 self-start rounded-xl border border-neutral-content bg-base-100 p-1.5 shadow-sm sm:self-auto">
          <SortDropdown sortBy={sortBy} setSort={setSortBy} t={t} />
          <Button
            variant="primary"
            className="h-9"
            onClick={() => setNewCollectionModal(true)}
          >
            <i className="bi-folder-plus" />
            <span className="hidden sm:inline">{t("new_collection")}</span>
          </Button>
        </div>
      </div>

      {!isLoading && collections.length === 0 ? (
        <div className="flex min-h-[22rem] flex-1 items-center justify-center rounded-2xl border border-dashed border-neutral-content bg-base-200/40 p-8">
          <div className="flex max-w-md flex-col items-center text-center">
            <div className="mb-4 flex h-14 w-14 items-center justify-center rounded-2xl border border-neutral-content bg-base-100 text-primary shadow-sm">
              <i className="bi-folder-plus text-2xl" />
            </div>
            <h2 className="text-lg font-semibold">{t("create_your_first_collection")}</h2>
            <p className="mt-2 text-sm leading-relaxed text-neutral">
              {t("create_your_first_collection_desc")}
            </p>
            <Button
              className="mt-5"
              variant="primary"
              onClick={() => setNewCollectionModal(true)}
            >
              <i className="bi-plus-lg" />
              {t("new_collection")}
            </Button>
          </div>
        </div>
      ) : (
        <>
          <section className="flex flex-col gap-3">
            <div className="flex items-center justify-between gap-3">
              <div>
                <h2 className="text-sm font-semibold">{t("collections_you_own")}</h2>
                <p className="text-xs text-neutral">
                  {ownedCollections.length === 1
                    ? t("showing_count_result", { count: ownedCollections.length })
                    : t("showing_count_results", { count: ownedCollections.length })}
                </p>
              </div>
            </div>

            {ownedCollections.length > 0 ? (
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4">
                {ownedCollections.map((collection) => (
                  <CollectionCard key={collection.id} collection={collection} />
                ))}
              </div>
            ) : (
              <div className="rounded-xl border border-dashed border-neutral-content bg-base-200/30 px-5 py-8 text-center text-sm text-neutral">
                {t("you_have_no_collections")}
              </div>
            )}
          </section>

          {sharedCollections.length > 0 && (
            <section className="flex flex-col gap-3 border-t border-neutral-content pt-5">
              <PageHeader
                icon="bi-people"
                title={t("other_collections")}
                description={t("other_collections_desc")}
                sm
              />

              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4">
                {sharedCollections.map((collection) => (
                  <CollectionCard key={collection.id} collection={collection} />
                ))}
              </div>
            </section>
          )}
        </>
      )}

      {newCollectionModal && (
        <NewCollectionModal onClose={() => setNewCollectionModal(false)} />
      )}
    </div>
  );
};

Page.getLayout = function getLayout(page: ReactElement<any>) {
  return <MainLayout>{page}</MainLayout>;
};

export default Page;

export { getServerSideProps };
