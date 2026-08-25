import React, { useEffect, useState } from "react";
import { TagIncludingLinkCount } from "@linkwarden/types/global";
import Modal from "../Modal";
import { Button } from "@/components/ui/button";
import { useTranslation } from "next-i18next";
import toast from "react-hot-toast";
import { useRemoveTag } from "@linkwarden/router/tags";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
  activeTag: TagIncludingLinkCount;
};

export default function DeleteTagModal({ onClose, activeTag }: Props) {
  const { t } = useTranslation();
  const [tag, setTag] = useState<TagIncludingLinkCount>(activeTag);

  const deleteTag = useRemoveTag();

  useEffect(() => {
    setTag(activeTag);
  }, []);

  const submit = async () => {
    const load = toast.loading(t("deleting"));

    await deleteTag.mutateAsync(tag.id as number, {
      onSettled: (data, error) => {
        toast.dismiss(load);

        if (error) {
          toast.error(error.message);
        } else {
          toast.success(t("deleted"));
          onClose();
        }
      },
    });
  };

  return (
    <Modal toggleModal={onClose}>
      <GlazeModalFrame
        title={t("delete_tag")}
        icon="bi-tag"
        tone="destructive"
        footer={
          <>
            <Button type="button" variant="ghost" onClick={() => onClose()}>
              {t("cancel")}
            </Button>
            <Button type="button" variant="destructive" onClick={submit}>
              <i className="bi-trash" />
              {t("delete")}
            </Button>
          </>
        }
      >
        <p>{t("tag_deletion_confirmation_message")}</p>
      </GlazeModalFrame>
    </Modal>
  );
}
