import React, { useEffect, useState } from "react";
import Modal from "../Modal";
import { Button } from "@/components/ui/button";
import { useTranslation } from "next-i18next";
import toast from "react-hot-toast";
import { RssSubscription } from "@linkwarden/prisma/client";
import { useDeleteRssSubscription } from "@linkwarden/router/rss";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
  rssSubscription: RssSubscription;
};

export default function DeleteRssSubscriptionModal({
  onClose,
  rssSubscription,
}: Props) {
  const { t } = useTranslation();
  const [subscription, setSubscription] =
    useState<RssSubscription>(rssSubscription);
  const [submitLoader, setSubmitLoader] = useState(false);
  const deleteRssSubscription = useDeleteRssSubscription();

  useEffect(() => {
    setSubscription(rssSubscription);
  }, [rssSubscription]);

  const submit = async () => {
    if (submitLoader) return;

    setSubmitLoader(true);
    const load = toast.loading(t("deleting"));

    await deleteRssSubscription.mutateAsync(subscription.id, {
      onSettled: (_, error) => {
        setSubmitLoader(false);
        toast.dismiss(load);

        if (error) {
          toast.error(error.message);
        } else {
          onClose();
          toast.success(t("rss_subscription_deleted"));
        }
      },
    });
  };

  return (
    <Modal toggleModal={onClose}>
      <GlazeModalFrame
        title={t("rss_subscriptions")}
        description={t("rss_deletion_confirmation")}
        icon="bi-rss"
        tone="destructive"
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
              variant="destructive"
              onClick={submit}
              disabled={submitLoader}
            >
              <i className="bi-trash" aria-hidden="true" />
              {t("delete")}
            </Button>
          </>
        }
      >
        <div className="rounded-xl border border-error/20 bg-error/5 p-3">
          <p className="text-xs font-medium text-base-content/60">
            {t("name")}
          </p>
          <p className="mt-1 break-words font-medium text-base-content">
            {subscription.name}
          </p>
        </div>
      </GlazeModalFrame>
    </Modal>
  );
}
