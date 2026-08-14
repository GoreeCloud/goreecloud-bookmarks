import React from "react";
import useLinkStore from "@/store/links";
import Modal from "../Modal";
import { Button } from "@/components/ui/button";
import { useTranslation } from "next-i18next";
import { useBulkDeleteLinks } from "@linkwarden/router/links";
import toast from "react-hot-toast";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
};

export default function BulkDeleteLinksModal({ onClose }: Props) {
  const { t } = useTranslation();
  const { selectedIds, clearSelected, selectionCount } = useLinkStore();

  const deleteLinksById = useBulkDeleteLinks();

  const deleteLink = async () => {
    const load = toast.loading(t("deleting"));
    const ids = Object.keys(selectedIds).map(Number);

    await deleteLinksById.mutateAsync(ids, {
      onSettled: (data, error) => {
        toast.dismiss(load);

        if (error) {
          toast.error(error.message);
        } else {
          clearSelected();
          onClose();
          toast.success(t("deleted"));
        }
      },
    });
  };

  return (
    <Modal toggleModal={onClose}>
      <GlazeModalFrame
        title={
          selectionCount === 1
            ? t("delete_link")
            : t("delete_links", { count: selectionCount })
        }
        icon="bi-trash"
        tone="destructive"
        footer={
          <>
            <Button type="button" variant="ghost" onClick={() => onClose()}>
              {t("cancel")}
            </Button>
            <Button type="button" variant="destructive" onClick={deleteLink}>
              <i className="bi-trash" />
              {t("delete")}
            </Button>
          </>
        }
      >
        <p>
          {selectionCount === 1
            ? t("link_deletion_confirmation_message")
            : t("links_deletion_confirmation_message", {
                count: selectionCount,
              })}
        </p>

        <div
          role="alert"
          className="flex items-start gap-3 rounded-xl border border-warning/25 bg-warning/10 p-3"
        >
          <i
            className="bi-exclamation-triangle mt-0.5 text-base text-warning"
            aria-hidden="true"
          />
          <span>{t("warning_irreversible")}</span>
        </div>

        <p className="text-base-content/60">{t("shift_key_tip")}</p>
      </GlazeModalFrame>
    </Modal>
  );
}
