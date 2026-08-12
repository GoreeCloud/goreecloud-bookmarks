import { LinkIncludingShortenedCollectionAndTags } from "@linkwarden/types/global";
import { useRouter } from "next/router";
import clsx from "clsx";
import usePinLink from "@/lib/client/pinLink";
import { Button } from "@/components/ui/button";

type Props = {
  link: LinkIncludingShortenedCollectionAndTags;
  btnStyle?: string;
};

export default function LinkPin({ link, btnStyle }: Props) {
  const pinLink = usePinLink();
  const router = useRouter();

  const isPublicRoute = router.pathname.startsWith("/public") ? true : false;
  const isAlreadyPinned = link?.pinnedBy && link.pinnedBy[0] ? true : false;

  return (
    <Button
      variant="simple"
      size="icon"
      className={clsx(
        "absolute top-3 right-[3.25rem] h-8 w-8 rounded-lg border border-base-content/10 bg-base-100/90 text-base-content/60 shadow-sm backdrop-blur-md transition-all duration-150 hover:bg-base-100 hover:text-base-content focus-visible:opacity-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/50",
        isAlreadyPinned
          ? "opacity-100 text-primary"
          : "opacity-0 group-hover:opacity-100 group-focus-within:opacity-100",
        btnStyle
      )}
      onClick={() => pinLink(link)}
      aria-label={isAlreadyPinned ? "Unpin bookmark" : "Pin bookmark"}
      title={isAlreadyPinned ? "Unpin bookmark" : "Pin bookmark"}
    >
      <i
        aria-hidden="true"
        className={clsx(
          "text-sm",
          isAlreadyPinned ? "bi-pin-fill" : "bi-pin"
        )}
      />
    </Button>
  );
}
