import React, { useEffect, useState } from "react";
import CollectionSelection from "@/components/InputSelect/CollectionSelection";
import TagSelection from "@/components/InputSelect/TagSelection";
import TextInput from "@/components/TextInput";
import unescapeString from "@/lib/client/unescapeString";
import { useRouter } from "next/router";
import toast from "react-hot-toast";
import Modal from "../Modal";
import { useTranslation } from "next-i18next";
import { useCollections } from "@linkwarden/router/collections";
import { useUploadFile } from "@linkwarden/router/links";
import { PostLinkSchemaType } from "@linkwarden/lib/schemaValidation";
import { useConfig } from "@linkwarden/router/config";
import { Button } from "@/components/ui/button";
import GlazeModalFrame from "./GlazeModalFrame";

type Props = {
  onClose: Function;
};

export default function UploadFileModal({ onClose }: Props) {
  const { t } = useTranslation();
  const { data: config } = useConfig();

  const initial = {
    name: "",
    description: "",
    type: "url",
    tags: [],
    collection: { id: undefined, name: "" },
  } as PostLinkSchemaType;

  const [link, setLink] = useState<PostLinkSchemaType>(initial);
  const [file, setFile] = useState<File>();
  const [submitLoader, setSubmitLoader] = useState(false);
  const [optionsExpanded, setOptionsExpanded] = useState(false);

  const uploadFile = useUploadFile();
  const router = useRouter();
  const { data: collections = [] } = useCollections();

  const setCollection = (option: any) => {
    setLink({
      ...link,
      collection: {
        id: option?.__isNew__ ? undefined : option?.value,
        name: option?.label || "",
      },
    });
  };

  const setTags = (options: any[]) => {
    const tagNames = options.map((option: any) => ({ name: option.label }));
    setLink({ ...link, tags: tagNames });
  };

  useEffect(() => {
    setOptionsExpanded(false);

    if (router.pathname.startsWith("/collections/") && router.query.id) {
      const currentCollection = collections.find(
        (collection) => collection.id === Number(router.query.id)
      );

      if (
        currentCollection &&
        currentCollection.ownerId &&
        router.asPath.startsWith("/collections/")
      ) {
        setLink({
          ...initial,
          collection: {
            id: currentCollection.id,
            name: currentCollection.name,
          },
        });
        return;
      }
    }

    setLink({ ...initial, collection: { name: "Unorganized" } });
  }, [router.pathname, router.query.id, router.asPath, collections]);

  const submit = async () => {
    if (submitLoader || !file) return;

    setSubmitLoader(true);
    const load = toast.loading(t("creating"));

    await uploadFile.mutateAsync(
      { link, file },
      {
        onSettled: (data, error) => {
          setSubmitLoader(false);
          toast.dismiss(load);

          if (error) {
            toast.error(error.message);
          } else {
            onClose();
            toast.success(t("created_success"));
          }
        },
      }
    );
  };

  return (
    <Modal toggleModal={onClose}>
      <GlazeModalFrame
        title={t("upload_file")}
        icon="bi-file-earmark-arrow-up"
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
              disabled={!file || submitLoader}
            >
              <i className="bi-upload" aria-hidden="true" />
              {t("upload_file")}
            </Button>
          </>
        }
      >
        <div className="grid gap-3 sm:grid-cols-5">
          <div className="rounded-xl border border-base-300 bg-base-200/50 p-3 sm:col-span-3">
            <label className="mb-2 block text-xs font-medium text-base-content/60">
              {t("file")}
            </label>
            <label className="flex min-h-11 w-full cursor-pointer items-center rounded-lg border border-base-300 bg-base-100 px-3 duration-150 hover:bg-base-200">
              <input
                type="file"
                accept=".pdf,.png,.jpg,.jpeg"
                className="custom-file-input w-full cursor-pointer"
                onChange={(e) => setFile(e.target.files?.[0])}
              />
            </label>
            <p className="mt-2 text-xs text-base-content/60">
              {t("file_types", { size: config?.MAX_FILE_BUFFER || 10 })}
            </p>
          </div>

          <div className="rounded-xl border border-base-300 bg-base-200/50 p-3 sm:col-span-2">
            <label className="mb-2 block text-xs font-medium text-base-content/60">
              {t("collection")}
            </label>
            {link.collection?.name && (
              <CollectionSelection
                onChange={setCollection}
                defaultValue={{
                  value: link.collection?.id,
                  label: link.collection?.name || "Unorganized",
                }}
              />
            )}
          </div>
        </div>

        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="w-fit px-2 text-sm"
          onClick={() => setOptionsExpanded(!optionsExpanded)}
          aria-expanded={optionsExpanded}
        >
          {optionsExpanded ? t("hide_options") : t("more_options")}
          <i
            className={`bi-chevron-${optionsExpanded ? "up" : "down"}`}
            aria-hidden="true"
          />
        </Button>

        {optionsExpanded && (
          <div className="grid gap-3 sm:grid-cols-2">
            <div className="rounded-xl border border-base-300 bg-base-200/50 p-3">
              <label className="mb-2 block text-xs font-medium text-base-content/60">
                {t("name")}
              </label>
              <TextInput
                value={link.name}
                onChange={(e) => setLink({ ...link, name: e.target.value })}
                placeholder={t("example_link")}
                className="bg-base-100"
              />
            </div>

            <div className="rounded-xl border border-base-300 bg-base-200/50 p-3">
              <label className="mb-2 block text-xs font-medium text-base-content/60">
                {t("tags")}
              </label>
              <TagSelection
                onChange={setTags}
                defaultValue={link.tags?.map((tag) => ({
                  value: tag.id,
                  label: tag.name,
                }))}
              />
            </div>

            <div className="rounded-xl border border-base-300 bg-base-200/50 p-3 sm:col-span-2">
              <label className="mb-2 block text-xs font-medium text-base-content/60">
                {t("description")}
              </label>
              <textarea
                value={unescapeString(link.description || "") || ""}
                onChange={(e) =>
                  setLink({ ...link, description: e.target.value })
                }
                placeholder={t("description_placeholder")}
                className="h-32 w-full resize-none rounded-lg border border-base-300 bg-base-100 p-3 text-sm outline-none duration-100 focus:border-primary"
              />
            </div>
          </div>
        )}
      </GlazeModalFrame>
    </Modal>
  );
}
