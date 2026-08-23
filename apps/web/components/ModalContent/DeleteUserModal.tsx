import Modal from "../Modal";
import { Button } from "@/components/ui/button";
import { useTranslation } from "next-i18next";
import { useDeleteUser } from "@linkwarden/router/users";
import { useState } from "react";
import { useSession } from "next-auth/react";
import { useConfig } from "@linkwarden/router/config";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
  userId: number;
};

export default function DeleteUserModal({ onClose, userId }: Props) {
  const { t } = useTranslation();
  const [submitLoader, setSubmitLoader] = useState(false);
  const deleteUser = useDeleteUser();
  const { data } = useSession();
  const { data: config } = useConfig();

  const isAdmin = data?.user?.id === (config?.ADMIN || 1);

  const submit = async () => {
    if (submitLoader) return;

    setSubmitLoader(true);

    await deleteUser.mutateAsync(userId, {
      onSuccess: () => {
        onClose();
      },
      onSettled: () => {
        setSubmitLoader(false);
      },
    });
  };

  return (
    <Modal toggleModal={onClose}>
      <GlazeModalFrame
        title={isAdmin ? t("delete_user") : t("remove_user")}
        description={t("confirm_user_deletion")}
        icon={isAdmin ? "bi-person-x" : "bi-person-dash"}
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
              <i
                className={isAdmin ? "bi-trash" : "bi-person-dash"}
                aria-hidden="true"
              />
              {isAdmin ? t("delete_confirmation") : t("confirm")}
            </Button>
          </>
        }
      >
        <p>{t("confirm_user_removal_desc")}</p>

        {isAdmin && (
          <div
            role="alert"
            className="flex items-start gap-3 rounded-xl border border-warning/25 bg-warning/10 p-3"
          >
            <i
              className="bi-exclamation-triangle mt-0.5 text-warning"
              aria-hidden="true"
            />
            <span>
              <b>{t("warning")}:</b> {t("irreversible_action_warning")}
            </span>
          </div>
        )}
      </GlazeModalFrame>
    </Modal>
  );
}
