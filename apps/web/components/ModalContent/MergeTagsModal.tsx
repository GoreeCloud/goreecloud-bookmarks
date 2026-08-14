import React, { useState } from "react";
import Modal from "../Modal";
import { Button } from "@/components/ui/button";
import { useTranslation } from "next-i18next";
import toast from "react-hot-toast";
import { useMergeTags } from "@linkwarden/router/tags";
import TextInput from "../TextInput";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
  selectedTags: number[];
  setSelectedTags: (tags: number[]) => void;
};

export default function MergeTagsModal({
  onClose,
  selectedTags,
  setSelectedTags,
}: Props) {
  const { t } = useTranslation();
  const [newTagName, setNewTagName] = useState("");
  const [submitLoader, setSubmitLoader] = useState(false);
  const mergeTags = useMergeTags();

  const merge = async () => {
    const trimmedTagName = newTagName.trim();
    if (submitLoader || !trimmedTagName) return;

    setSubmitLoader(true);
    const load = toast.loading(t("merging"));

    await mergeTags.mutateAsync(
      {
        tagIds: selectedTags,
        newTagName: trimmedTagName,
      },
      {
        onSettled: (data, error) => {
          setSubmitLoader(false);
          toast.dismiss(load);

          if (error) {
            toast.error(error.message);
          } else {
            setSelectedTags([]);
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
        title={t("merge_count_tags", { count: selectedTags.length })}
        description={t("rename_tag_instruction")}
        icon="bi-intersect"
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
              onClick={merge}
              disabled={submitLoader || !newTagName.trim()}
            >
              <i className="bi-intersect" aria-hidden="true" />
              {t("merge_tags")}
            </Button>
          </>
        }
      >
        <div className="rounded-xl border border-base-300 bg-base-200/50 p-3">
          <p className="mb-2 text-xs font-medium text-base-content/60">
            {t("tag")}
          </p>
          <TextInput
            value={newTagName}
            onChange={(e) => setNewTagName(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                e.preventDefault();
                merge();
              }
            }}
            placeholder={t("tag_name_placeholder")}
            className="bg-base-100"
          />
        </div>
      </GlazeModalFrame>
    </Modal>
  );
}
