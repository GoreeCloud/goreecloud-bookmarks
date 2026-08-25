import React from "react";
import Modal from "../Modal";
import { Button } from "@/components/ui/button";
import { useTranslation } from "next-i18next";
import toast from "react-hot-toast";
import { useBulkTagDeletion } from "@linkwarden/router/tags";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
  selectedTags: number[];
  setSelectedTags: (tags: number[]) => void;
};

export default function BulkDeleteTagsModal({
  onClose,
  selectedTags,
  setSelectedTags,
}: Props) {
  const { t } = useTranslation();

  const deleteTagsById = useBulkTagDeletion();

  const deleteTag = async () => {
    const load = toast.loading(t("deleting"));

    await deleteTagsById.mutateAsync(
      {
        tagIds: selectedTags,
      },
      {
        onSettled: (data, error) => {
          toast.dismiss(load);

          if (error) {
            toast.error(error.message);
          } else {
            setSelectedTags([]);
            onClose();
            toast.success(t("deleted"));
          }
        },
      }
    );
  };

  return (
    <Modal toggleModal={onClose}>
      <GlazeModalFrame
        title={
          selectedTags.length === 1
            ? t("delete_tag")
            : t("delete_tags", { count: selectedTags.length })
        }
        icon="bi-tags"
        tone="destructive"
        footer={
          <>
            <Button type="button" variant="ghost" onClick={() => onClose()}>
              {t("cancel")}
            </Button>
            <Button type="button" variant="destructive" onClick={deleteTag}>
              <i className="bi-trash" />
              {t("delete")}
            </Button>
          </>
        }
      >
        <p>
          {selectedTags.length === 1
            ? t("tag_deletion_confirmation_message")
            : t("tags_deletion_confirmation_message", {
                count: selectedTags.length,
              })}
        </p>
      </GlazeModalFrame>
    </Modal>
  );
}
