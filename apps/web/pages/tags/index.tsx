import MainLayout from "@/layouts/MainLayout";
import PageHeader from "@/components/PageHeader";
import getServerSideProps from "@/lib/client/getServerSideProps";
import { useTranslation } from "next-i18next";
import { useTags } from "@linkwarden/router/tags";
import TagCard from "@/components/TagCard";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { ReactElement, useEffect, useState } from "react";
import NewTagModal from "@/components/ModalContent/NewTagModal";
import BulkDeleteTagsModal from "@/components/ModalContent/BulkDeleteTagsModal";
import MergeTagsModal from "@/components/ModalContent/MergeTagsModal";
import { NextPageWithLayout } from "../_app";
import { TagSort } from "@linkwarden/types/global";
import { useInView } from "react-intersection-observer";

const Page: NextPageWithLayout = () => {
  const { t } = useTranslation();
  const [sortBy, setSortBy] = useState<TagSort>(TagSort.DateNewestFirst);
  const { ref, inView } = useInView();
  const {
    data: tags = [],
    isLoading,
    hasNextPage,
    isFetchingNextPage,
    fetchNextPage,
  } = useTags(undefined, {
    sort: sortBy,
  });

  const [newTagModal, setNewTagModal] = useState(false);
  const [bulkDeleteModal, setBulkDeleteModal] = useState(false);
  const [mergeTagsModal, setMergeTagsModal] = useState(false);
  const [editMode, setEditMode] = useState(false);
  const [selectedTags, setSelectedTags] = useState<number[]>([]);

  useEffect(() => {
    if (!inView || !hasNextPage || isFetchingNextPage) return;
    fetchNextPage();
  }, [inView, hasNextPage, isFetchingNextPage, fetchNextPage]);

  return (
    <div className="flex w-full flex-1 flex-col gap-6 p-3 sm:p-5">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <PageHeader
          icon="bi-hash"
          title={t("tags")}
          description={
            tags.length === 1
              ? t("showing_count_result", { count: tags.length })
              : t("showing_count_results", { count: tags.length })
          }
        />

        <div className="flex items-center gap-1.5 self-start rounded-xl border border-neutral-content bg-base-100 p-1.5 shadow-sm sm:self-auto">
          <Button
            variant="ghost"
            size="icon"
            className={editMode ? "bg-primary/15 hover:bg-primary/15" : ""}
            onClick={() => {
              setEditMode(!editMode);
              setSelectedTags([]);
            }}
            aria-label={t("edit")}
          >
            <i className="bi-check2-square text-neutral" />
          </Button>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" aria-label="Sort tags">
                <i className="bi-arrow-down-up text-neutral" />
              </Button>
            </DropdownMenuTrigger>

            <DropdownMenuContent sideOffset={4} align="end">
              <DropdownMenuRadioGroup
                value={sortBy.toString()}
                onValueChange={(value) => setSortBy(Number(value) as TagSort)}
              >
                <DropdownMenuRadioItem value={TagSort.DateNewestFirst.toString()}>
                  {t("date_newest_first")}
                </DropdownMenuRadioItem>
                <DropdownMenuRadioItem value={TagSort.DateOldestFirst.toString()}>
                  {t("date_oldest_first")}
                </DropdownMenuRadioItem>
                <DropdownMenuRadioItem value={TagSort.NameAZ.toString()}>
                  {t("name_az")}
                </DropdownMenuRadioItem>
                <DropdownMenuRadioItem value={TagSort.NameZA.toString()}>
                  {t("name_za")}
                </DropdownMenuRadioItem>
                <DropdownMenuRadioItem value={TagSort.LinkCountHighLow.toString()}>
                  {t("link_count_high_low")}
                </DropdownMenuRadioItem>
                <DropdownMenuRadioItem value={TagSort.LinkCountLowHigh.toString()}>
                  {t("link_count_low_high")}
                </DropdownMenuRadioItem>
              </DropdownMenuRadioGroup>
            </DropdownMenuContent>
          </DropdownMenu>

          <Button
            variant="primary"
            className="h-9"
            onClick={() => setNewTagModal(true)}
          >
            <i className="bi-tag" />
            <span className="hidden sm:inline">{t("new_tag")}</span>
          </Button>
        </div>
      </div>

      {tags.length > 0 && editMode && (
        <div className="flex flex-col gap-3 rounded-xl border border-neutral-content bg-base-200/50 p-3 sm:flex-row sm:items-center sm:justify-between">
          <label className="flex items-center gap-3 text-sm">
            <input
              type="checkbox"
              className="checkbox checkbox-primary"
              onChange={() => {
                if (selectedTags.length === tags.length) setSelectedTags([]);
                else setSelectedTags(tags.map((tag) => tag.id));
              }}
              checked={selectedTags.length === tags.length && tags.length > 0}
            />
            <span>
              {selectedTags.length > 0
                ? selectedTags.length === 1
                  ? t("tag_selected")
                  : t("tags_selected", { count: selectedTags.length })
                : t("nothing_selected")}
            </span>
          </label>

          <div className="flex items-center gap-2">
            <Button
              onClick={() => setMergeTagsModal(true)}
              variant="simple"
              size="sm"
              disabled={selectedTags.length < 2}
            >
              <i className="bi-intersect" />
              {t("merge_tags")}
            </Button>
            <Button
              onClick={() => setBulkDeleteModal(true)}
              variant="ghost"
              size="sm"
              disabled={selectedTags.length === 0}
              className="text-error"
            >
              <i className="bi-trash" />
              {t("delete")}
            </Button>
          </div>
        </div>
      )}

      {!isLoading && tags.length === 0 ? (
        <div className="flex min-h-[22rem] flex-1 items-center justify-center rounded-2xl border border-dashed border-neutral-content bg-base-200/40 p-8">
          <div className="flex max-w-md flex-col items-center text-center">
            <div className="mb-4 flex h-14 w-14 items-center justify-center rounded-2xl border border-neutral-content bg-base-100 text-primary shadow-sm">
              <i className="bi-tags text-2xl" />
            </div>
            <h2 className="text-lg font-semibold">{t("create_your_first_tag")}</h2>
            <p className="mt-2 text-sm leading-relaxed text-neutral">
              {t("create_your_first_tag_desc")}
            </p>
            <Button
              className="mt-5"
              variant="primary"
              onClick={() => setNewTagModal(true)}
            >
              <i className="bi-plus-lg" />
              {t("new_tag")}
            </Button>
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
          {tags.map((tag: any) => (
            <TagCard
              key={tag.id}
              tag={tag}
              selected={selectedTags.includes(tag.id)}
              editMode={editMode}
              onSelect={(id: number) => {
                if (selectedTags.includes(id))
                  setSelectedTags((previous) =>
                    previous.filter((selectedId) => selectedId !== id)
                  );
                else setSelectedTags((previous) => [...previous, id]);
              }}
            />
          ))}
          {isLoading && !tags.length && <TagCardSkeleton />}
          {isFetchingNextPage && <TagCardSkeleton />}
        </div>
      )}

      {hasNextPage && <div ref={ref} className="h-1 w-full" />}

      {newTagModal && <NewTagModal onClose={() => setNewTagModal(false)} />}
      {bulkDeleteModal && (
        <BulkDeleteTagsModal
          onClose={() => {
            setBulkDeleteModal(false);
            setEditMode(false);
          }}
          selectedTags={selectedTags}
          setSelectedTags={setSelectedTags}
        />
      )}
      {mergeTagsModal && (
        <MergeTagsModal
          onClose={() => {
            setMergeTagsModal(false);
            setEditMode(false);
          }}
          selectedTags={selectedTags}
          setSelectedTags={setSelectedTags}
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

const TagCardSkeleton = () => {
  return (
    <div className="min-h-[9.5rem] rounded-xl border border-neutral-content bg-base-100 p-4">
      <div className="flex gap-3">
        <div className="skeleton h-9 w-9 rounded-lg" />
        <div className="flex-1">
          <div className="skeleton h-4 w-3/4" />
          <div className="skeleton mt-2 h-3 w-24" />
        </div>
      </div>
    </div>
  );
};
