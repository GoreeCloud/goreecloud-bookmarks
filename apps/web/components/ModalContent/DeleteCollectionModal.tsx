import React, { useEffect, useState } from "react";
import { CollectionIncludingMembersAndLinkCount } from "@linkwarden/types/global";
import { useRouter } from "next/router";
import usePermissions from "@/hooks/usePermissions";
import Modal from "../Modal";
import { Button } from "@/components/ui/button";
import { useTranslation } from "next-i18next";
import { useDeleteCollection } from "@linkwarden/router/collections";
import toast from "react-hot-toast";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
  activeCollection: CollectionIncludingMembersAndLinkCount;
};

export default function DeleteCollectionModal({
  onClose,
  activeCollection,
}: Props) {
  const { t } = useTranslation();
  const [collection, setCollection] =
    useState<CollectionIncludingMembersAndLinkCount>(activeCollection);
  const router = useRouter();
  const permissions = usePermissions(collection.id as number);

  useEffect(() => {
    setCollection(activeCollection);
  }, []);

  const deleteCollection = useDeleteCollection({ toast, t });

  const submit = async () => {
    if (!collection) return null;

    deleteCollection.mutateAsync(collection.id as number);

    onClose();
    router.push("/collections");
  };

  const canDelete = permissions === true;

  return (
    <Modal toggleModal={onClose}>
      <GlazeModalFrame
        title={canDelete ? t("delete_collection") : t("leave_collection")}
        icon={canDelete ? "bi-folder-x" : "bi-box-arrow-right"}
        tone={canDelete ? "destructive" : "default"}
        footer={
          <>
            <Button type="button" variant="ghost" onClick={() => onClose()}>
              {t("cancel")}
            </Button>
            <Button type="button" onClick={submit} variant="destructive">
              <i className={canDelete ? "bi-trash" : "bi-box-arrow-right"} />
              {canDelete ? t("delete") : t("leave")}
            </Button>
          </>
        }
      >
        {canDelete ? (
          <>
            <p>{t("collection_deletion_prompt")}</p>
            <div
              role="alert"
              className="flex items-start gap-3 rounded-xl border border-warning/25 bg-warning/10 p-3"
            >
              <i
                className="bi-exclamation-triangle mt-0.5 text-base text-warning"
                aria-hidden="true"
              />
              <span>
                <b>{t("warning")}: </b>
                {t("deletion_warning")}
              </span>
            </div>
          </>
        ) : (
          <p>{t("leave_prompt")}</p>
        )}
      </GlazeModalFrame>
    </Modal>
  );
}
