import { LinkIncludingShortenedCollectionAndTags } from "@linkwarden/types/global";
import Link from "next/link";
import React, { useEffect, useState } from "react";

function LinkTypeBadge({
  link,
}: {
  link: LinkIncludingShortenedCollectionAndTags;
}) {
  const [url, setUrl] = useState("");

  useEffect(() => {
    if (link.type === "url" && link.url) {
      try {
        setUrl(new URL(link.url).host.toLowerCase());
      } catch (error) {
        console.log(error);
      }
    }
  }, [link]);

  const typeIcon = () => {
    switch (link.type) {
      case "pdf":
        return "bi-file-earmark-pdf";
      case "image":
        return "bi-file-earmark-image";
      default:
        return "bi-link-45deg";
    }
  };

  const badgeClassName =
    "inline-flex max-w-full w-fit items-center gap-1.5 rounded-md bg-base-content/[0.055] px-2 py-1 text-[11px] font-medium leading-none text-base-content/55 ring-1 ring-inset ring-base-content/[0.06] transition-colors";

  return link.url && url ? (
    <Link
      href={link.url || ""}
      target="_blank"
      rel="noreferrer"
      title={link.url || ""}
      onClick={(e) => {
        e.stopPropagation();
      }}
      className={`${badgeClassName} hover:bg-base-content/[0.09] hover:text-base-content/75 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/50`}
    >
      <i className="bi-globe2 shrink-0 text-[10px]" aria-hidden="true" />
      <span className="truncate">{url}</span>
    </Link>
  ) : (
    <div className={badgeClassName}>
      <i className={`${typeIcon()} shrink-0 text-[10px]`} aria-hidden="true" />
      <span className="truncate capitalize">{link.type}</span>
    </div>
  );
}

export default React.memo(LinkTypeBadge);
