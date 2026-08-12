import { useLinks } from "@linkwarden/router/links";
import MainLayout from "@/layouts/MainLayout";
import { Sort, ViewMode } from "@linkwarden/types/global";
import { useRouter } from "next/router";
import React, { ReactElement, useEffect, useState } from "react";
import LinkListOptions from "@/components/LinkListOptions";
import getServerSideProps from "@/lib/client/getServerSideProps";
import { useTranslation } from "next-i18next";
import Links from "@/components/LinkViews/Links";
import SearchBar from "@/components/SearchBar";
import { NextPageWithLayout } from "./_app";

const Page: NextPageWithLayout = () => {
  const { t } = useTranslation();
  const router = useRouter();

  const query =
    typeof router.query.q === "string"
      ? decodeURIComponent(router.query.q)
      : "";

  const [viewMode, setViewMode] = useState<ViewMode>(
    (localStorage.getItem("viewMode") as ViewMode) || ViewMode.Card
  );

  const [sortBy, setSortBy] = useState<Sort>(
    Number(localStorage.getItem("sortBy")) ?? Sort.DateNewestFirst
  );

  const [editMode, setEditMode] = useState(false);

  useEffect(() => {
    if (editMode) return setEditMode(false);
  }, [router]);

  const { links, data } = useLinks({
    sort: sortBy,
    searchQueryString: query,
  });

  return (
    <div className="flex h-full w-full flex-col gap-5 p-4 sm:p-5 lg:p-6">
      <section className="rounded-2xl border border-base-content/10 bg-base-100 p-4 shadow-sm sm:p-5">
        <div className="flex items-start gap-3.5">
          <span className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-primary/[0.09] text-primary">
            <i className="bi-search text-lg" aria-hidden="true" />
          </span>
          <div className="min-w-0">
            <h1 className="text-xl font-semibold tracking-[-0.02em] text-base-content sm:text-2xl">
              {t("search_results")}
            </h1>
            <p className="mt-1 truncate text-xs text-base-content/45 sm:text-sm">
              {query ? `${t("results_for")} “${query}”` : t("search_for_links")}
            </p>
          </div>
        </div>

        <div className="mt-4 max-w-3xl">
          <SearchBar fullWidth placeholder={t("search_for_links")} />
        </div>

        <div className="mt-4 border-t border-base-content/[0.07] pt-3">
          <LinkListOptions
            t={t}
            viewMode={viewMode}
            setViewMode={setViewMode}
            sortBy={sortBy}
            setSortBy={setSortBy}
            editMode={editMode}
            setEditMode={setEditMode}
            links={links}
          >
            <div className="min-w-0">
              <p className="text-xs font-semibold uppercase tracking-[0.08em] text-base-content/40">
                {query || t("search_results")}
              </p>
            </div>
          </LinkListOptions>
        </div>
      </section>

      {!data.isLoading && links && !links[0] ? (
        <div className="flex min-h-56 flex-col items-center justify-center rounded-2xl border border-dashed border-base-content/15 bg-base-content/[0.018] px-6 py-10 text-center">
          <span className="flex h-12 w-12 items-center justify-center rounded-2xl bg-base-content/[0.045] text-base-content/40">
            <i className="bi-search text-xl" aria-hidden="true" />
          </span>
          <p className="mt-4 text-sm font-semibold text-base-content/70">
            {t("nothing_found")}
          </p>
          {query && (
            <p className="mt-1 max-w-md text-xs leading-5 text-base-content/40">
              “{query}”
            </p>
          )}
        </div>
      ) : (
        <Links
          editMode={editMode}
          links={links}
          layout={viewMode}
          useData={data}
        />
      )}
    </div>
  );
};

Page.getLayout = function getLayout(page: ReactElement<any>) {
  return <MainLayout>{page}</MainLayout>;
};

export default Page;

export { getServerSideProps };
