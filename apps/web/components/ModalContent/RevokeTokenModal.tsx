import React, { useEffect, useState } from "react";
import Modal from "../Modal";
import { Button } from "@/components/ui/button";
import { useTranslation } from "next-i18next";
import { AccessToken } from "@linkwarden/prisma/client";
import { useRevokeToken } from "@linkwarden/router/tokens";
import toast from "react-hot-toast";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
  activeToken: AccessToken;
};

export default function RevokeTokenModal({ onClose, activeToken }: Props) {
  const { t } = useTranslation();
  const [token, setToken] = useState<AccessToken>(activeToken);
  const [submitLoader, setSubmitLoader] = useState(false);
  const revokeToken = useRevokeToken();

  useEffect(() => {
    setToken(activeToken);
  }, [activeToken]);

  const revoke = async () => {
    if (submitLoader) return;

    setSubmitLoader(true);
    const load = toast.loading(t("sending_request"));

    await revokeToken.mutateAsync(token.id, {
      onSettled: (data, error) => {
        setSubmitLoader(false);
        toast.dismiss(load);

        if (error) {
          toast.error(error.message);
        } else {
          onClose();
          toast.success(t("token_revoked"));
        }
      },
    });
  };

  return (
    <Modal toggleModal={onClose}>
      <GlazeModalFrame
        title={t("revoke_token")}
        description={t("revoke_confirmation")}
        icon="bi-key"
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
              onClick={revoke}
              disabled={submitLoader}
            >
              <i className="bi-slash-circle" aria-hidden="true" />
              {t("revoke")}
            </Button>
          </>
        }
      >
        <div className="rounded-xl border border-error/20 bg-error/5 p-3">
          <p className="text-xs font-medium text-base-content/60">
            {t("name")}
          </p>
          <p className="mt-1 break-words font-medium text-base-content">
            {token.name}
          </p>
        </div>
      </GlazeModalFrame>
    </Modal>
  );
}
