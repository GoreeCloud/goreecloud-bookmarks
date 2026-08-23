import React, { useState } from "react";
import CollectionSelection from "@/components/InputSelect/CollectionSelection";
import TagSelection from "@/components/InputSelect/TagSelection";
import useLinkStore from "@/store/links";
import { LinkIncludingShortenedCollectionAndTags } from "@linkwarden/types/global";
import toast from "react-hot-toast";
import Modal from "../Modal";
import { useTranslation } from "next-i18next";
import { useBulkEditLinks } from "@linkwarden/router/links";
import { Button } from "../ui/button";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
};

export default function BulkEditLinksModal({ onClose }: Props) {
  const { t } = useTranslation();
  const { selectedIds, clearSelected, selectionCount } = useLinkStore();
  const [submitLoader, setSubmitLoader] = useState(false);
  const [removePreviousTags, setRemovePreviousTags] = useState(false);
  const [updatedValues, setUpdatedValues] = useState<
    Pick<LinkIncludingShortenedCollectionAndTags, "tags" | "collectionId">
  >({ tags: [] });

  const updateLinks = useBulkEditLinks();

  const setCollection = (e: any) => {
    const collectionId = e?.value || null;
    setUpdatedValues((prevValues) => ({ ...prevValues, collectionId }));
  };

  const setTags = (e: any) => {
    const tags = e.map((tag: any) => ({ name: tag.label }));
    setUpdatedValues((prevValues) => ({ ...prevValues, tags }));
  };

  const submit = async () => {
    if (submitLoader) return;

    setSubmitLoader(true);
    const load = toast.loading(t("updating"));

    const links = Object.keys(selectedIds).map((k) => ({
      id: Number(k),
    }));

    await updateLinks.mutateAsync(
      {
        links,
        newData: updatedValues,
        removePreviousTags,
      },
      {
        onSettled: (data, error) => {
          setSubmitLoader(false);
          toast.dismiss(load);

          if (error) {
            toast.error(error.message);
          } else {
            clearSelected();
            onClose();
            toast.success(t("updated"));
          }
        },
      }
    );
  };

  return (
    <Modal toggleModal={onClose}>
      <GlazeModalFrame
        title={
          selectionCount === 1
            ? t("edit_link")
            : t("edit_links", { count: selectionCount })
        }
        icon="bi-pencil-square"
        footer={
          <>
            <Button
              type="button"
              variant="ghost"
              onClick={() => onClose()}
              disabled={submitLoader}
            >
              {t("cancel")}
            </Button>
            <Button
              type="button"
              variant="primary"
              onClick={submit}
              disabled={submitLoader}
            >
              <i className="bi-check-lg" aria-hidden="true" />
              {t("save_changes")}
            </Button>
          </>
        }
      >
        <div className="grid gap-3 sm:grid-cols-2">
          <div className="rounded-xl border border-base-300 bg-base-200/50 p-3">
            <p className="mb-2 text-xs font-medium text-base-content/60">
              {t("move_to_collection")}
            </p>
            <CollectionSelection
              showDefaultValue={false}
              onChange={setCollection}
              creatable={false}
            />
          </div>

          <div className="rounded-xl border border-base-300 bg-base-200/50 p-3">
            <p className="mb-2 text-xs font-medium text-base-content/60">
              {t("add_tags")}
            </p>
            <TagSelection onChange={setTags} />
          </div>
        </div>

        <label className="flex cursor-pointer items-start gap-3 rounded-xl border border-base-300 bg-base-200/40 p-3">
          <input
            type="checkbox"
            className="checkbox checkbox-primary mt-0.5"
            checked={removePreviousTags}
            onChange={(e) => setRemovePreviousTags(e.target.checked)}
          />
          <span className="block font-medium text-base-content">
            {t("remove_previous_tags")}
          </span>
        </label>
      </GlazeModalFrame>
    </Modal>
  );
}
