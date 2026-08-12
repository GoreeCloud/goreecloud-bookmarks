import React, { FormEvent, useLayoutEffect, useRef, useState } from "react";
import TextInput from "@/components/TextInput";
import toast from "react-hot-toast";
import Modal from "../Modal";
import { Button } from "@/components/ui/button";
import { useTranslation } from "next-i18next";
import { Separator } from "../ui/separator";
import { useUpsertTags } from "@linkwarden/router/tags";

type Props = {
  onClose: Function;
};

export default function NewTagModal({ onClose }: Props) {
  const { t } = useTranslation();
  const upsertTags = useUpsertTags();
  const [tag, setTag] = useState({ label: "" });
  const [submitLoader, setSubmitLoader] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  useLayoutEffect(() => {
    inputRef.current?.focus();
  }, []);

  const submit = async (event?: FormEvent) => {
    event?.preventDefault();
    if (submitLoader || !tag.label.trim()) return;

    setSubmitLoader(true);
    const load = toast.loading(t("creating"));

    await upsertTags.mutateAsync([tag], {
      onSettled: (data, error) => {
        setSubmitLoader(false);
        toast.dismiss(load);
        if (error) {
          toast.error(t(error.message));
        } else {
          onClose();
          toast.success(t("created"));
        }
      },
    });
  };

  return (
    <Modal toggleModal={onClose}>
      <div className="flex items-start gap-3 pr-9">
        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-neutral-content bg-base-200 text-primary">
          <i className="bi-tag" />
        </div>
        <div>
          <h2 className="text-lg font-semibold">{t("create_new_tag")}</h2>
          <p className="mt-1 text-xs text-neutral">{t("tags")}</p>
        </div>
      </div>

      <Separator className="my-4" />

      <form onSubmit={submit} className="flex flex-col gap-5">
        <div>
          <p className="mb-2 text-xs font-medium text-neutral">{t("name")}</p>
          <TextInput
            ref={inputRef}
            value={tag.label}
            onChange={(e) => setTag({ ...tag, label: e.target.value })}
            className="h-11 bg-base-200"
            placeholder={t("tag_name_placeholder")}
          />
        </div>

        <div className="flex items-center justify-end gap-2 border-t border-neutral-content pt-4">
          <Button type="button" variant="ghost" onClick={() => onClose()}>
            {t("cancel")}
          </Button>
          <Button
            type="submit"
            variant="primary"
            disabled={submitLoader || !tag.label.trim()}
          >
            <i className="bi-plus-lg" />
            {t("create_new_tag")}
          </Button>
        </div>
      </form>
    </Modal>
  );
}
