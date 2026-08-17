import React, { useEffect, useState } from "react";
import { LinkIncludingShortenedCollectionAndTags } from "@linkwarden/types/global";
import Modal from "../Modal";
import { useRouter } from "next/router";
import { Button } from "@/components/ui/button";
import { useTranslation } from "next-i18next";
import { useDeleteLink } from "@linkwarden/router/links";
import toast from "react-hot-toast";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
  activeLink: LinkIncludingShortenedCollectionAndTags;
};

export default function DeleteLinkModal({ onClose, activeLink }: Props) {
  const { t } = useTranslation();
  const [link, setLink] =
    useState<LinkIncludingShortenedCollectionAndTags>(activeLink);

  const deleteLink = useDeleteLink({ toast, t });
  const router = useRouter();

  useEffect(() => {
    setLink(activeLink);
  }, []);

  const submit = async () => {
    deleteLink.mutateAsync(link.id as number);

    if (
      router.pathname.startsWith("/links/[id]") ||
      router.pathname.startsWith("/preserved/[id]")
    ) {
      router.push("/dashboard");
    }
    onClose();
  };

  return (
    <Modal toggleModal={onClose}>
      <GlazeModalFrame
        title={t("delete_link")}
        icon="bi-trash"
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
        <p>{t("link_deletion_confirmation_message")}</p>

        <div
          role="note"
          className="flex items-start gap-3 rounded-xl border border-base-300 bg-base-200/60 p-3"
        >
          <i
            className="bi-info-circle mt-0.5 text-base text-primary"
            aria-hidden="true"
          />
          <span>
            <b>{t("tip")}:</b> {t("shift_key_tip")}
          </span>
        </div>
      </GlazeModalFrame>
    </Modal>
  );
}
