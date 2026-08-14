import React, { useLayoutEffect, useRef, useState } from "react";
import Modal from "../Modal";
import { useTranslation } from "next-i18next";
import { useAddRssSubscription } from "@linkwarden/router/rss";
import toast from "react-hot-toast";
import TextInput from "../TextInput";
import CollectionSelection from "../InputSelect/CollectionSelection";
import { Button } from "../ui/button";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
};

export default function NewRssSubscriptionModal({ onClose }: Props) {
  const { t } = useTranslation();
  const addRssSubscription = useAddRssSubscription();
  const [submitLoader, setSubmitLoader] = useState(false);

  const [form, setForm] = useState({
    name: "",
    url: "",
    collectionId: 0,
    collectionName: "",
  });

  const normalizedName = form.name.trim();
  const normalizedUrl = form.url.trim();
  const normalizedCollectionName = form.collectionName.trim();
  const hasCollection = Boolean(form.collectionId || normalizedCollectionName);
  const canSubmit = Boolean(normalizedName && normalizedUrl && hasCollection) && !submitLoader;

  const submit = async () => {
    if (submitLoader) return;

    if (!normalizedName || !normalizedUrl || !hasCollection) {
      toast.error(t("fill_all_fields"));
      return;
    }

    setSubmitLoader(true);
    const load = toast.loading(t("creating"));

    await addRssSubscription.mutateAsync(
      {
        ...form,
        name: normalizedName,
        url: normalizedUrl,
        collectionName: normalizedCollectionName,
      },
      {
        onSettled: (_, error) => {
          setSubmitLoader(false);
          toast.dismiss(load);

          if (error) {
            toast.error(error.message);
          } else {
            onClose();
            toast.success(t("created"));
          }
        },
      }
    );
  };

  const inputRef = useRef<HTMLInputElement>(null);

  useLayoutEffect(() => {
    inputRef.current?.focus();
  }, []);

  return (
    <Modal toggleModal={onClose}>
      <GlazeModalFrame
        title={t("create_rss_subscription")}
        icon="bi-rss"
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
              variant="primary"
              onClick={submit}
              disabled={!canSubmit}
            >
              <i className="bi-rss" aria-hidden="true" />
              {t("create_rss_subscription")}
            </Button>
          </>
        }
      >
        <div className="rounded-xl border border-base-300 bg-base-200/50 p-3">
          <label className="mb-2 block text-xs font-medium text-base-content/60">
            {t("link")}
          </label>
          <TextInput
            ref={inputRef}
            type="url"
            placeholder="https://example.com/rss"
            className="bg-base-100"
            value={form.url}
            onChange={(e) => setForm({ ...form, url: e.target.value })}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                e.preventDefault();
                submit();
              }
            }}
          />
        </div>

        <div className="grid gap-3 sm:grid-cols-2">
          <div className="rounded-xl border border-base-300 bg-base-200/50 p-3">
            <label className="mb-2 block text-xs font-medium text-base-content/60">
              {t("name")}
            </label>
            <TextInput
              type="text"
              placeholder="Sample RSS"
              className="bg-base-100"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
            />
          </div>

          <div className="rounded-xl border border-base-300 bg-base-200/50 p-3">
            <label className="mb-2 block text-xs font-medium text-base-content/60">
              {t("collection")}
            </label>
            <CollectionSelection
              onChange={(option: any) => {
                setForm({
                  ...form,
                  collectionId: option?.__isNew__ ? 0 : option?.value || 0,
                  collectionName: option?.label || "",
                });
              }}
            />
          </div>
        </div>
      </GlazeModalFrame>
    </Modal>
  );
}
