import { useRouter } from "next/router";
import { FormEvent, ReactElement, useEffect, useState } from "react";
import MainLayout from "@/layouts/MainLayout";
import { Sort, ViewMode } from "@linkwarden/types/global";
import { useLinks } from "@linkwarden/router/links";
import { useTranslation } from "next-i18next";
import getServerSideProps from "@/lib/client/getServerSideProps";
import LinkListOptions from "@/components/LinkListOptions";
import { useRemoveTag, useTag, useUpdateTag } from "@linkwarden/router/tags";
import Links from "@/components/LinkViews/Links";
import toast from "react-hot-toast";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu";
import { NextPageWithLayout } from "../_app";

const Page: NextPageWithLayout = () => {
  const { t } = useTranslation();
  const router = useRouter();
  const parsedTagId = Number(router.query.id);
  const tagId =
    Number.isFinite(parsedTagId) && parsedTagId > 0 ? parsedTagId : undefined;

  const {
    data: activeTag,
    isLoading: isTagLoading,
    isError: isTagError,
  } = useTag(tagId);
  const updateTag = useUpdateTag();
  const removeTag = useRemoveTag();

  const [sortBy, setSortBy] = useState<Sort>(
    Number(localStorage.getItem("sortBy")) ?? Sort.DateNewestFirst
  );
  const [viewMode, setViewMode] = useState<ViewMode>(
    (localStorage.getItem("viewMode") as ViewMode) || ViewMode.Card
  );
  const [renameTag, setRenameTag] = useState(false);
  const [newTagName, setNewTagName] = useState<string>();
  const [editMode, setEditMode] = useState(false);
  const [submitLoader, setSubmitLoader] = useState(false);

  const { links, data } = useLinks({
    sort: sortBy,
    tagId,
  });

  useEffect(() => {
    if (editMode) setEditMode(false);
  }, [router]);

  useEffect(() => {
    if (!router.isReady || isTagLoading) return;
    if (isTagError) router.push("/dashboard");
  }, [router, router.isReady, isTagLoading, isTagError]);

  useEffect(() => {
    setNewTagName(activeTag?.name);
  }, [activeTag]);

  const cancelUpdateTag = () => {
    setNewTagName(activeTag?.name);
    setRenameTag(false);
  };

  const submit = async (event?: FormEvent) => {
    event?.preventDefault();

    if (activeTag?.name === newTagName) return setRenameTag(false);
    if (!newTagName?.trim()) return cancelUpdateTag();

    setSubmitLoader(true);

    if (activeTag) {
      const load = toast.loading(t("applying_changes"));

      await updateTag.mutateAsync(
        {
          ...activeTag,
          name: newTagName.trim(),
        },
        {
          onSettled: (result, error) => {
            setSubmitLoader(false);
            toast.dismiss(load);

            if (error) toast.error(error.message);
            else toast.success(t("tag_renamed"));
          },
        }
      );
    }

    setRenameTag(false);
  };

  const remove = async () => {
    setSubmitLoader(true);

    if (activeTag?.id) {
      const load = toast.loading(t("applying_changes"));

      await removeTag.mutateAsync(activeTag.id, {
        onSettled: (result, error) => {
          toast.dismiss(load);

          if (error) toast.error(error.message);
          else {
            toast.success(t("tag_deleted"));
            router.push("/links");
          }
        },
      });
    }

    setSubmitLoader(false);
    setRenameTag(false);
  };

  return (
    <div className="flex h-full w-full flex-col gap-5 p-3 sm:p-5">
      <div className="rounded-2xl border border-neutral-content bg-base-100 p-4 shadow-sm sm:p-5">
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
          <div className="flex min-w-0 items-center gap-3">
            <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl border border-neutral-content bg-base-200 text-primary">
              <i className="bi-hash text-xl" />
            </div>

            <div className="min-w-0">
              {renameTag ? (
                <form onSubmit={submit} className="flex min-w-0 items-center gap-2">
                  <input
                    type="text"
                    autoFocus
                    className="h-10 min-w-0 flex-1 rounded-lg border border-neutral-content bg-base-200 px-3 text-lg font-semibold outline-none focus:border-primary sm:w-72"
                    value={newTagName || ""}
                    onChange={(e) => setNewTagName(e.target.value)}
                  />
                  <Button
                    type="submit"
                    variant="primary"
                    size="icon"
                    disabled={submitLoader || !newTagName?.trim()}
                    aria-label={t("save_changes")}
                  >
                    <i className="bi-check2" />
                  </Button>
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    onClick={cancelUpdateTag}
                    aria-label={t("cancel")}
                  >
                    <i className="bi-x text-lg" />
                  </Button>
                </form>
              ) : (
                <div className="flex min-w-0 items-center gap-1">
                  <div className="min-w-0">
                    <h1 className="truncate text-xl font-semibold sm:text-2xl">
                      {activeTag?.name}
                    </h1>
                    <p className="mt-0.5 text-xs text-neutral">
                      {links.length === 1
                        ? t("showing_count_result", { count: links.length })
                        : t("showing_count_results", { count: links.length })}
                    </p>
                  </div>

                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="icon" title={t("more")}>
                        <i className="bi-three-dots text-lg text-neutral" />
                      </Button>
                    </DropdownMenuTrigger>

                    <DropdownMenuContent sideOffset={4} align="start">
                      <DropdownMenuItem onClick={() => setRenameTag(true)}>
                        <i className="bi-pencil-square" />
                        {t("rename_tag")}
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem onClick={remove} className="text-error">
                        <i className="bi-trash" />
                        {t("delete_tag")}
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              )}
            </div>
          </div>
        </LinkListOptions>
      </div>

      <Links
        editMode={editMode}
        links={links}
        layout={viewMode}
        useData={data}
      />

      {!data.isLoading && links.length === 0 && (
        <div className="flex min-h-[18rem] flex-1 items-center justify-center rounded-2xl border border-dashed border-neutral-content bg-base-200/40 p-8">
          <div className="max-w-md text-center">
            <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-xl border border-neutral-content bg-base-100 text-primary">
              <i className="bi-tag text-xl" />
            </div>
            <h2 className="text-lg font-semibold">{t("this_tag_has_no_links")}</h2>
            <p className="mt-2 text-sm leading-relaxed text-neutral">
              {t("this_tag_has_no_links_desc")}
            </p>
          </div>
        </div>
      )}
    </div>
  );
};

Page.getLayout = function getLayout(page: ReactElement<any>) {
  return <MainLayout>{page}</MainLayout>;
};

export default Page;

export { getServerSideProps };
