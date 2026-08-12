import MainLayout from "@/layouts/MainLayout";
import { ReactElement, ReactNode, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import React from "react";
import { toast } from "react-hot-toast";
import DashboardItem from "@/components/DashboardItem";
import NewLinkModal from "@/components/ModalContent/NewLinkModal";
import getServerSideProps from "@/lib/client/getServerSideProps";
import { useTranslation } from "next-i18next";
import { useCollections } from "@linkwarden/router/collections";
import { useDashboardData } from "@linkwarden/router/dashboardData";
import { useUpdateUser, useUser } from "@linkwarden/router/user";
import SurveyModal from "@/components/ModalContent/SurveyModal";
import ImportDropdown from "@/components/ImportDropdown";
import { Button } from "@/components/ui/button";
import DashboardLayoutDropdown from "@/components/DashboardLayoutDropdown";
import {
  DashboardSection,
  DashboardSectionType,
} from "@linkwarden/prisma/client";
import { DashboardLinks } from "@/components/DashboardLinks";
import { ViewMode } from "@linkwarden/types/global";
import ViewDropdown from "@/components/ViewDropdown";
import clsx from "clsx";
import Icon from "@/components/Icon";
import Droppable from "@/components/Droppable";
import { NextPageWithLayout } from "./_app";

const Page: NextPageWithLayout = () => {
  const { t } = useTranslation();
  const { data: collections = [] } = useCollections();
  const {
    data: {
      links = [],
      numberOfPinnedLinks,
      numberOfTags = 0,
      collectionLinks = {},
    } = {
      links: [],
    },
    ...dashboardData
  } = useDashboardData();

  const { data: user } = useUser();

  const [numberOfLinks, setNumberOfLinks] = useState(0);

  const [dashboardSections, setDashboardSections] = useState<
    DashboardSection[]
  >(user?.dashboardSections || []);

  useEffect(() => {
    setDashboardSections(user?.dashboardSections || []);
  }, [user?.dashboardSections]);

  const [viewMode, setViewMode] = useState<ViewMode>(
    (localStorage.getItem("viewMode") as ViewMode) || ViewMode.Card
  );

  useEffect(() => {
    setNumberOfLinks(
      collections.reduce(
        (accumulator, collection) =>
          accumulator + (collection._count as any).links,
        0
      )
    );
  }, [collections]);

  useEffect(() => {
    if (
      process.env.NEXT_PUBLIC_STRIPE === "true" &&
      user &&
      user.id &&
      user.referredBy === null &&
      // if user is using Linkwarden for more than 3 days
      new Date().getTime() - new Date(user.createdAt).getTime() >
        3 * 24 * 60 * 60 * 1000
    ) {
      setTimeout(() => {
        setShowsSurveyModal(true);
      }, 1000);
    }
  }, [user]);

  const orderedSections = useMemo(() => {
    return [...dashboardSections].sort((a, b) => {
      const orderA = a.order ?? Number.MAX_SAFE_INTEGER;
      const orderB = b.order ?? Number.MAX_SAFE_INTEGER;
      return orderA - orderB;
    });
  }, [dashboardSections]);

  const [newLinkModal, setNewLinkModal] = useState(false);

  const [showSurveyModal, setShowsSurveyModal] = useState(false);

  const updateUser = useUpdateUser();

  const [submitLoader, setSubmitLoader] = useState(false);

  // Function to render the dragged item
  const submitSurvey = async (referer: string, other?: string) => {
    if (submitLoader) return;

    setSubmitLoader(true);

    const load = toast.loading(t("applying"));

    await updateUser.mutateAsync(
      {
        ...user,
        referredBy: referer === "other" ? "Other: " + other : referer,
      },
      {
        onSettled: (data, error) => {
          console.log(data, error);
          setSubmitLoader(false);
          toast.dismiss(load);

          if (error) {
            toast.error(error.message);
          } else {
            toast.success(t("thanks_for_feedback"));
            setShowsSurveyModal(false);
          }
        },
      }
    );
  };

  return (
    <>
      <div className="mx-auto flex w-full max-w-[1600px] grow flex-col gap-6 p-4 sm:p-6 lg:p-7">
        <header className="flex flex-col gap-4 border-b border-neutral-content/70 pb-5 sm:flex-row sm:items-end sm:justify-between">
          <div className="min-w-0">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-neutral-content/70 bg-base-200/70 text-primary">
                <i className="bi-house-fill text-lg" />
              </div>
              <div className="min-w-0">
                <h1 className="text-2xl font-semibold tracking-tight sm:text-3xl">
                  {t("dashboard")}
                </h1>
                <p className="mt-0.5 max-w-2xl text-sm text-neutral">
                  {t("dashboard_desc")}
                </p>
              </div>
            </div>
          </div>

          <div className="flex w-fit items-center gap-1 rounded-xl border border-neutral-content/70 bg-base-200/55 p-1 shadow-sm">
            <DashboardLayoutDropdown />
            <ViewDropdown
              viewMode={viewMode}
              setViewMode={setViewMode}
              dashboard
            />
          </div>
        </header>

        {orderedSections[0] ? (
          <div className="flex flex-col gap-5">
            {orderedSections.map((section, i) => (
              <Section
                key={i}
                sectionData={section}
                t={t}
                collection={collections.find(
                  (c) => c.id === section.collectionId
                )}
                collectionLinks={
                  section.collectionId
                    ? collectionLinks[section.collectionId]
                    : []
                }
                links={links}
                numberOfTags={numberOfTags}
                numberOfLinks={numberOfLinks}
                collectionsLength={collections.length}
                numberOfPinnedLinks={numberOfPinnedLinks}
                dashboardData={dashboardData}
                setNewLinkModal={setNewLinkModal}
              />
            ))}
          </div>
        ) : (
          <DashboardSkeleton />
        )}
      </div>

      {showSurveyModal && (
        <SurveyModal
          submit={submitSurvey}
          onClose={() => {
            setShowsSurveyModal(false);
          }}
        />
      )}
      {newLinkModal && (
        <NewLinkModal
          onClose={() => {
            setNewLinkModal(false);
          }}
        />
      )}
    </>
  );
};

Page.getLayout = function getLayout(page: ReactElement<any>) {
  return <MainLayout>{page}</MainLayout>;
};

export default Page;

export { getServerSideProps };

type SectionProps = {
  sectionData: DashboardSection;
  t: (key: string) => string;
  collection: any;
  collectionsLength: number;
  links: any[];
  numberOfTags: number;
  numberOfLinks: number;
  numberOfPinnedLinks: number;
  dashboardData: any;
  collectionLinks: any[];
  setNewLinkModal: (value: boolean) => void;
};

type DashboardSectionHeaderProps = {
  title: string;
  description?: string;
  icon?: string;
  leading?: ReactNode;
  href?: string;
  viewAllLabel?: string;
};

const DashboardSectionHeader = ({
  title,
  description,
  icon,
  leading,
  href,
  viewAllLabel,
}: DashboardSectionHeaderProps) => (
  <div className="flex items-start justify-between gap-4">
    <div className="flex min-w-0 items-start gap-3">
      {leading || (
        <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
          <i className={clsx(icon, "text-base")} />
        </div>
      )}
      <div className="min-w-0">
        <h2 className="truncate text-base font-semibold tracking-tight sm:text-lg">
          {title}
        </h2>
        {description && (
          <p className="mt-0.5 text-xs text-neutral sm:text-sm">{description}</p>
        )}
      </div>
    </div>

    {href && viewAllLabel && (
      <Link
        href={href}
        className="flex shrink-0 items-center gap-1 rounded-lg px-2.5 py-1.5 text-sm font-medium text-neutral transition-colors hover:bg-base-content/5 hover:text-base-content focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/50"
      >
        {viewAllLabel}
        <i className="bi-chevron-right text-xs" />
      </Link>
    )}
  </div>
);

const DashboardSectionSurface = ({ children }: { children: ReactNode }) => (
  <section className="flex flex-col gap-4 rounded-2xl border border-neutral-content/70 bg-base-100/55 p-4 shadow-sm sm:p-5">
    {children}
  </section>
);

const DashboardEmptyState = ({
  icon,
  title,
  description,
  children,
  minHeight = "min-h-56",
}: {
  icon: string;
  title: string;
  description: string;
  children?: ReactNode;
  minHeight?: string;
}) => (
  <div
    className={clsx(
      "flex w-full flex-col items-center justify-center gap-2 rounded-xl border border-dashed border-neutral-content/80 bg-base-200/35 px-6 py-9 text-center",
      minHeight
    )}
  >
    <div className="mb-1 flex h-11 w-11 items-center justify-center rounded-xl bg-primary/10 text-primary">
      <i className={clsx(icon, "text-xl")} />
    </div>
    <p className="text-base font-semibold tracking-tight">{title}</p>
    <p className="max-w-md text-sm leading-6 text-neutral">{description}</p>
    {children && <div className="mt-3 flex flex-wrap justify-center gap-2">{children}</div>}
  </div>
);

const DashboardSkeleton = () => (
  <div className="flex flex-col gap-5">
    <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-4">
      <div className="skeleton h-24 w-full rounded-2xl" />
      <div className="skeleton h-24 w-full rounded-2xl" />
      <div className="skeleton h-24 w-full rounded-2xl" />
      <div className="skeleton h-24 w-full rounded-2xl" />
    </div>
    {[0, 1, 2].map((item) => (
      <div
        key={item}
        className="rounded-2xl border border-neutral-content/50 bg-base-100/40 p-5"
      >
        <div className="mb-5 flex items-center gap-3">
          <div className="skeleton h-9 w-9 rounded-lg" />
          <div className="flex flex-col gap-2">
            <div className="skeleton h-4 w-32" />
            <div className="skeleton h-3 w-48" />
          </div>
        </div>
        <div className="skeleton h-56 w-full rounded-xl" />
      </div>
    ))}
  </div>
);

const Section = ({
  sectionData,
  t,
  collection,
  links,
  numberOfTags,
  numberOfLinks,
  collectionsLength,
  numberOfPinnedLinks,
  dashboardData,
  collectionLinks,
  setNewLinkModal,
}: SectionProps) => {
  switch (sectionData.type) {
    case DashboardSectionType.STATS:
      return (
        <section className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-4">
          <DashboardItem
            name={numberOfLinks === 1 ? t("link") : t("links")}
            value={numberOfLinks}
            icon={"bi-link-45deg"}
          />

          <DashboardItem
            name={collectionsLength === 1 ? t("collection") : t("collections")}
            value={collectionsLength}
            icon={"bi-folder"}
          />

          <DashboardItem
            name={numberOfTags === 1 ? t("tag") : t("tags")}
            value={numberOfTags}
            icon={"bi-hash"}
          />

          <DashboardItem
            name={t("pinned")}
            value={numberOfPinnedLinks}
            icon={"bi-pin-angle"}
          />
        </section>
      );
    case DashboardSectionType.RECENT_LINKS:
      return (
        <DashboardSectionSurface>
          <DashboardSectionHeader
            icon="bi-clock-history"
            title={t("recent_links")}
            description={t("recent_links_desc")}
            href="/links"
            viewAllLabel={t("view_all")}
          />

          {dashboardData.isLoading ||
          (links && links[0] && !dashboardData.isLoading) ? (
            <DashboardLinks
              type="recent"
              links={links}
              isLoading={dashboardData.isLoading}
            />
          ) : (
            <DashboardEmptyState
              icon="bi-link-45deg"
              title={t("view_added_links_here")}
              description={t("view_added_links_here_desc")}
            >
              <Button
                onClick={() => {
                  setNewLinkModal(true);
                }}
                variant="primary"
              >
                <i className="bi-plus-lg text-lg" />
                {t("add_link")}
              </Button>
              <ImportDropdown />
            </DashboardEmptyState>
          )}
        </DashboardSectionSurface>
      );
    case DashboardSectionType.PINNED_LINKS: {
      const hasPinnedLinks = links?.some(
        (e: any) => e.pinnedBy && e.pinnedBy[0]
      );

      return (
        <DashboardSectionSurface>
          <DashboardSectionHeader
            icon="bi-pin-angle"
            title={t("pinned_links")}
            description={t("pinned_links_desc")}
            href="/links/pinned"
            viewAllLabel={t("view_all")}
          />
          <Droppable
            id="pinned-links-section"
            data={{
              name: "pinned-links",
            }}
            className={clsx(
              !dashboardData.isLoading &&
                !hasPinnedLinks &&
                "grow flex flex-col"
            )}
          >
            {dashboardData.isLoading || hasPinnedLinks ? (
              <DashboardLinks
                links={links.filter((e: any) => e.pinnedBy && e.pinnedBy[0])}
                isLoading={dashboardData.isLoading}
              />
            ) : (
              <DashboardEmptyState
                icon="bi-pin-angle"
                title={t("pin_favorite_links_here")}
                description={t("pin_favorite_links_here_desc")}
              />
            )}
          </Droppable>
        </DashboardSectionSurface>
      );
    }
    case DashboardSectionType.COLLECTION: {
      const hasCollectionLinks = collectionLinks?.length > 0;

      return (
        collection?.id && (
          <DashboardSectionSurface>
            <DashboardSectionHeader
              leading={
                <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-base-200/80">
                  {collection.icon ? (
                    <Icon
                      icon={collection.icon}
                      color={collection.color || "#0ea5e9"}
                    />
                  ) : (
                    <i
                      className="bi-folder-fill text-lg"
                      style={{ color: collection.color || "#0ea5e9" }}
                    />
                  )}
                </div>
              }
              title={collection.name}
              description={collection.description || undefined}
              href={`/collections/${collection.id}`}
              viewAllLabel={t("view_all")}
            />
            <Droppable
              id={`dashboard-${collection.id}`}
              data={{
                id: collection.id,
                name: collection.name,
                ownerId: collection.ownerId,
                type: "collection",
              }}
              className={clsx(
                !dashboardData.isLoading &&
                  !hasCollectionLinks &&
                  "grow flex flex-col"
              )}
            >
              {dashboardData.isLoading || hasCollectionLinks ? (
                <DashboardLinks
                  type="collection"
                  links={collectionLinks}
                  isLoading={dashboardData.isLoading}
                />
              ) : (
                <DashboardEmptyState
                  icon="bi-folder"
                  title={t("no_link_in_collection")}
                  description={t("no_link_in_collection_desc")}
                  minHeight="min-h-64"
                />
              )}
            </Droppable>
          </DashboardSectionSurface>
        )
      );
    }
    default:
      return null;
  }
};
