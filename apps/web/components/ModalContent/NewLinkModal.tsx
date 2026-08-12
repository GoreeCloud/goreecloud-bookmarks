import React, { useEffect, useLayoutEffect, useRef, useState } from "react";
import CollectionSelection from "@/components/InputSelect/CollectionSelection";
import TagSelection from "@/components/InputSelect/TagSelection";
import TextInput from "@/components/TextInput";
import unescapeString from "@/lib/client/unescapeString";
import { useRouter } from "next/router";
import Modal from "../Modal";
import { useTranslation } from "next-i18next";
import { useCollections } from "@linkwarden/router/collections";
import toast from "react-hot-toast";
import {
  PostLinkSchema,
  PostLinkSchemaType,
} from "@linkwarden/lib/schemaValidation";
import { Button } from "@/components/ui/button";
import { useAddLink } from "@linkwarden/router/links";

type Props = {
  onClose: () => void;
};

export default function NewLinkModal({ onClose }: Props) {
  const { t } = useTranslation();
  const initial = {
    name: "",
    url: "",
    description: "",
    type: "url",
    tags: [],
    collection: {
      id: undefined,
      name: "",
    },
  } as PostLinkSchemaType;

  const addLink = useAddLink({
    toast,
    t,
  });

  const inputRef = useRef<HTMLInputElement>(null);
  const [link, setLink] = useState<PostLinkSchemaType>(initial);
  const [optionsExpanded, setOptionsExpanded] = useState(false);
  const router = useRouter();
  const { data: collections = [] } = useCollections();

  const setCollection = (e: any) => {
    if (e?.__isNew__) e.value = undefined;
    setLink({
      ...link,
      collection: { id: e?.value, name: e?.label },
    });
  };

  const setTags = (selectedOptions: any = []) => {
    const tagNames = selectedOptions.map((option: any) => ({
      name: option.label,
    }));
    setLink({ ...link, tags: tagNames });
  };

  useEffect(() => {
    if (router.pathname.startsWith("/collections/") && router.query.id) {
      const currentCollection = collections.find(
        (e) => e.id == Number(router.query.id)
      );

      if (currentCollection && currentCollection.ownerId)
        setLink({
          ...initial,
          collection: {
            id: currentCollection.id,
            name: currentCollection.name,
          },
        });
    } else
      setLink({
        ...initial,
        collection: { name: "Unorganized" },
      });
  }, []);

  useLayoutEffect(() => {
    inputRef.current?.focus();
  }, []);

  const submit = async () => {
    const dataValidation = PostLinkSchema.safeParse(link);

    if (!dataValidation.success)
      return toast.error(
        `Error: ${
          dataValidation.error.issues[0].message
        } [${dataValidation.error.issues[0].path.join(", ")}]`
      );

    addLink.mutateAsync(link);
    onClose();
  };

  const sourceHost = (() => {
    const value = link.url?.trim();
    if (!value) return "";

    try {
      const parsed = new URL(
        value.includes("://") ? value : `https://${value}`
      );
      return parsed.hostname.replace(/^www\./, "");
    } catch {
      return "";
    }
  })();

  return (
    <Modal toggleModal={onClose} className="sm:!max-w-2xl">
      <form
        onSubmit={(e) => {
          e.preventDefault();
          submit();
        }}
      >
        <div className="pr-10">
          <div className="flex items-center gap-3">
            <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-primary/[0.09] text-primary">
              <i className="bi-bookmark-plus text-lg" aria-hidden="true" />
            </span>
            <div className="min-w-0">
              <h2 className="text-lg font-semibold tracking-[-0.015em] text-base-content sm:text-xl">
                {t("create_new_link")}
              </h2>
              <p className="mt-0.5 text-xs text-base-content/45">
                {t("new_link")}
              </p>
            </div>
          </div>
        </div>

        <div className="mt-5 rounded-2xl border border-base-content/10 bg-base-content/[0.02] p-3.5 sm:p-4">
          <div className="grid gap-4 sm:grid-cols-5">
            <div className="sm:col-span-3">
              <div className="mb-2 flex items-center justify-between gap-2">
                <label
                  htmlFor="new-link-url"
                  className="text-xs font-semibold text-base-content/65"
                >
                  {t("link")}
                </label>
                {sourceHost && (
                  <span className="max-w-[65%] truncate rounded-md bg-primary/[0.07] px-2 py-1 text-[11px] font-medium text-primary/75">
                    {sourceHost}
                  </span>
                )}
              </div>

              <div className="relative">
                <span className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-base-content/35">
                  <i className="bi-link-45deg" aria-hidden="true" />
                </span>
                <TextInput
                  id="new-link-url"
                  ref={inputRef}
                  value={link.url || ""}
                  onChange={(e) => setLink({ ...link, url: e.target.value })}
                  placeholder={t("link_url_placeholder")}
                  className="h-11 border-base-content/10 bg-base-100 pl-10 text-sm shadow-sm transition-all hover:border-base-content/20 focus:border-primary/40 focus:ring-2 focus:ring-primary/10"
                />
              </div>
            </div>

            <div className="sm:col-span-2">
              <p className="mb-2 text-xs font-semibold text-base-content/65">
                {t("collection")}
              </p>
              {link.collection?.name && (
                <div className="rounded-xl bg-base-100">
                  <CollectionSelection
                    onChange={setCollection}
                    defaultValue={{
                      value: link.collection?.id,
                      label: link.collection?.name || "Unorganized",
                    }}
                  />
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="mt-3 flex items-center justify-between gap-3">
          <Button
            type="button"
            variant="ghost"
            size="sm"
            className="h-8 rounded-lg px-2.5 text-xs text-base-content/55 hover:bg-base-content/[0.05] hover:text-base-content"
            aria-expanded={optionsExpanded}
            onClick={() => setOptionsExpanded(!optionsExpanded)}
          >
            <i
              className={`bi-sliders2 mr-1.5 text-xs`}
              aria-hidden="true"
            />
            <span>{optionsExpanded ? t("hide_options") : t("more_options")}</span>
            <i
              className={`bi-chevron-${optionsExpanded ? "up" : "down"} ml-1 text-[10px]`}
              aria-hidden="true"
            />
          </Button>

          {!optionsExpanded && (
            <span className="hidden text-[11px] text-base-content/35 sm:inline">
              Enter
            </span>
          )}
        </div>

        {optionsExpanded && (
          <div className="mt-3 rounded-2xl border border-base-content/[0.08] bg-base-content/[0.018] p-3.5 sm:p-4">
            <div className="grid gap-4 sm:grid-cols-2">
              <div>
                <label
                  htmlFor="new-link-name"
                  className="mb-2 block text-xs font-semibold text-base-content/65"
                >
                  {t("name")}
                </label>
                <TextInput
                  id="new-link-name"
                  value={link.name}
                  onChange={(e) => setLink({ ...link, name: e.target.value })}
                  placeholder={t("link_name_placeholder")}
                  className="h-10 border-base-content/10 bg-base-100 text-sm shadow-sm transition-all hover:border-base-content/20 focus:border-primary/40 focus:ring-2 focus:ring-primary/10"
                />
              </div>

              <div>
                <p className="mb-2 text-xs font-semibold text-base-content/65">
                  {t("tags")}
                </p>
                <TagSelection
                  onChange={setTags}
                  defaultValue={
                    link.tags?.map((e) => ({ label: e.name, value: e.id })) || []
                  }
                />
              </div>

              <div className="sm:col-span-2">
                <label
                  htmlFor="new-link-description"
                  className="mb-2 block text-xs font-semibold text-base-content/65"
                >
                  {t("description")}
                </label>
                <textarea
                  id="new-link-description"
                  value={unescapeString(link.description || "") || ""}
                  onChange={(e) =>
                    setLink({ ...link, description: e.target.value })
                  }
                  placeholder={t("link_description_placeholder")}
                  className="h-28 w-full resize-none rounded-xl border border-base-content/10 bg-base-100 p-3 text-sm text-base-content shadow-sm outline-none transition-all placeholder:text-base-content/35 hover:border-base-content/20 focus:border-primary/40 focus:ring-2 focus:ring-primary/10"
                />
              </div>
            </div>
          </div>
        )}

        <div className="mt-5 flex items-center justify-end gap-2 border-t border-base-content/[0.07] pt-4">
          <Button
            type="button"
            variant="ghost"
            className="rounded-lg text-base-content/60"
            onClick={onClose}
          >
            {t("cancel")}
          </Button>
          <Button
            type="submit"
            variant="primary"
            className="min-w-28 rounded-lg shadow-sm"
          >
            <i className="bi-bookmark-plus" aria-hidden="true" />
            {t("create_link")}
          </Button>
        </div>
      </form>
    </Modal>
  );
}
