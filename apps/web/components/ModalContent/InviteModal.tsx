import toast from "react-hot-toast";
import Modal from "../Modal";
import TextInput from "../TextInput";
import { FormEvent, useLayoutEffect, useRef, useState } from "react";
import { useTranslation } from "next-i18next";
import { useAddUser } from "@linkwarden/router/users";
import { signIn } from "next-auth/react";
import { Button } from "../ui/button";
import { useConfig } from "@linkwarden/router/config";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
};

type FormData = {
  username?: string;
  email?: string;
  invite: boolean;
};

export default function InviteModal({ onClose }: Props) {
  const { t } = useTranslation();
  const { data: config } = useConfig();
  const emailEnabled = config?.EMAIL_PROVIDER;
  const addUser = useAddUser();

  const [form, setForm] = useState<FormData>({
    username: emailEnabled ? undefined : "",
    email: emailEnabled ? "" : undefined,
    invite: true,
  });
  const [submitLoader, setSubmitLoader] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const requiredValue = emailEnabled ? form.email : form.username;
  const canSubmit = Boolean(requiredValue?.trim()) && !submitLoader;

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (submitLoader) return;

    const value = requiredValue?.trim();
    if (!value) {
      toast.error(t("fill_all_fields_error"));
      return;
    }

    setSubmitLoader(true);

    const payload = emailEnabled
      ? { ...form, email: value }
      : { ...form, username: value };

    await addUser.mutateAsync(payload, {
      onSettled: () => {
        setSubmitLoader(false);
      },
      onSuccess: async () => {
        await signIn("invite", {
          email: payload.email,
          callbackUrl: "/member-onboarding",
          redirect: false,
        });
        onClose();
      },
    });
  }

  useLayoutEffect(() => {
    inputRef.current?.focus();
  }, []);

  return (
    <Modal toggleModal={onClose}>
      <GlazeModalFrame
        title={t("invite_user")}
        description={emailEnabled ? t("invite_user_desc") : undefined}
        icon="bi-person-plus"
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
              type="submit"
              form="invite-user-form"
              variant="primary"
              disabled={!canSubmit}
            >
              <i className="bi-send" aria-hidden="true" />
              {t("send_invitation")}
            </Button>
          </>
        }
      >
        <form id="invite-user-form" onSubmit={submit}>
          <div className="rounded-xl border border-base-300 bg-base-200/50 p-3">
            <label className="mb-2 block text-xs font-medium text-base-content/60">
              {emailEnabled ? t("email") : t("username")}
            </label>
            <TextInput
              ref={inputRef}
              type={emailEnabled ? "email" : "text"}
              placeholder={
                emailEnabled ? t("placeholder_email") : t("placeholder_john")
              }
              className="bg-base-100"
              onChange={(e) =>
                emailEnabled
                  ? setForm({ ...form, email: e.target.value })
                  : setForm({ ...form, username: e.target.value })
              }
              value={emailEnabled ? form.email : form.username}
            />
          </div>
        </form>
      </GlazeModalFrame>
    </Modal>
  );
}
