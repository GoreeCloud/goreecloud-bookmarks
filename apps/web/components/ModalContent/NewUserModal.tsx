import toast from "react-hot-toast";
import Modal from "../Modal";
import TextInput from "../TextInput";
import { FormEvent, useState } from "react";
import { useTranslation, Trans } from "next-i18next";
import { useAddUser } from "@linkwarden/router/users";
import { Button } from "../ui/button";
import { useConfig } from "@linkwarden/router/config";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
};

type FormData = {
  name: string;
  username?: string;
  email?: string;
  password: string;
};

export default function NewUserModal({ onClose }: Props) {
  const { t } = useTranslation();
  const { data: config } = useConfig();
  const emailEnabled = Boolean(config?.EMAIL_PROVIDER);
  const addUser = useAddUser();

  const [form, setForm] = useState<FormData>({
    name: "",
    username: "",
    email: emailEnabled ? "" : undefined,
    password: "",
  });
  const [submitLoader, setSubmitLoader] = useState(false);

  const normalizedName = form.name.trim();
  const normalizedUsername = form.username?.trim() || "";
  const normalizedEmail = form.email?.trim() || "";
  const requiredIdentity = emailEnabled ? normalizedEmail : normalizedUsername;
  const canSubmit =
    Boolean(normalizedName && requiredIdentity && form.password.length >= 8) &&
    !submitLoader;

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (submitLoader) return;

    if (form.password.length < 8) {
      toast.error(t("password_length_error"));
      return;
    }

    if (!normalizedName || !requiredIdentity) {
      toast.error(t("fill_all_fields_error"));
      return;
    }

    setSubmitLoader(true);

    await addUser.mutateAsync(
      {
        ...form,
        name: normalizedName,
        username: normalizedUsername || undefined,
        email: emailEnabled ? normalizedEmail : undefined,
      },
      {
        onSuccess: () => {
          onClose();
        },
        onSettled: () => {
          setSubmitLoader(false);
        },
      }
    );
  }

  return (
    <Modal toggleModal={onClose}>
      <GlazeModalFrame
        title={t("create_new_user")}
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
              form="create-user-form"
              variant="primary"
              disabled={!canSubmit}
            >
              <i className="bi-person-plus" aria-hidden="true" />
              {t("create_user")}
            </Button>
          </>
        }
      >
        <form id="create-user-form" onSubmit={submit}>
          <div className="grid gap-3 sm:grid-cols-2">
            <div className="rounded-xl border border-base-300 bg-base-200/50 p-3">
              <label className="mb-2 block text-xs font-medium text-base-content/60">
                {t("display_name")}
              </label>
              <TextInput
                placeholder={t("placeholder_johnny")}
                className="bg-base-100"
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                value={form.name}
              />
            </div>

            {emailEnabled && (
              <div className="rounded-xl border border-base-300 bg-base-200/50 p-3">
                <label className="mb-2 block text-xs font-medium text-base-content/60">
                  {t("email")}
                </label>
                <TextInput
                  type="email"
                  placeholder={t("placeholder_email")}
                  className="bg-base-100"
                  onChange={(e) => setForm({ ...form, email: e.target.value })}
                  value={form.email}
                />
              </div>
            )}

            <div className="rounded-xl border border-base-300 bg-base-200/50 p-3">
              <label className="mb-2 block text-xs font-medium text-base-content/60">
                {t("username")}{" "}
                {emailEnabled && (
                  <span className="font-normal">({t("optional")})</span>
                )}
              </label>
              <TextInput
                placeholder={t("placeholder_john")}
                className="bg-base-100"
                onChange={(e) => setForm({ ...form, username: e.target.value })}
                value={form.username}
              />
            </div>

            <div className="rounded-xl border border-base-300 bg-base-200/50 p-3">
              <label className="mb-2 block text-xs font-medium text-base-content/60">
                {t("password")}
              </label>
              <TextInput
                type="password"
                autoComplete="new-password"
                placeholder="••••••••••••••"
                className="bg-base-100"
                onChange={(e) => setForm({ ...form, password: e.target.value })}
                value={form.password}
              />
            </div>
          </div>

          <div
            role="note"
            className="mt-3 flex items-start gap-3 rounded-xl border border-base-300 bg-base-200/50 p-3 text-sm"
          >
            <i
              className="bi-info-circle mt-0.5 text-primary"
              aria-hidden="true"
            />
            <span>
              <Trans
                i18nKey="password_change_note"
                components={[<b key={0} />]}
              />
            </span>
          </div>
        </form>
      </GlazeModalFrame>
    </Modal>
  );
}
