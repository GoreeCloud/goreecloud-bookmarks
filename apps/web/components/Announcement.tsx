import Link from "next/link";
import React, { MouseEventHandler } from "react";
import { Button } from "@/components/ui/button";

type Props = {
  toggleAnnouncementBar: MouseEventHandler<HTMLButtonElement>;
};

export default function Announcement({ toggleAnnouncementBar }: Props) {
  const announcementId = localStorage.getItem("announcementId");
  const announcementMessage = localStorage.getItem("announcementMessage");

  if (!announcementId && !announcementMessage) return null;

  return (
    <div className="fixed mx-auto bottom-20 sm:bottom-10 w-full pointer-events-none p-5 z-30">
      <div className="mx-auto pointer-events-auto p-2 flex justify-between gap-2 items-center border border-primary shadow-xl rounded-xl bg-base-300 backdrop-blur-sm bg-opacity-80 max-w-md">
        <i
          className="bi-stars text-xl text-yellow-600 dark:text-yellow-500"
          aria-hidden="true"
        ></i>
        <p className="w-4/5 text-center text-sm sm:text-base">
          {announcementId ? (
            <>
              See what&apos;s new in{" "}
              <Link
                href="https://github.com/GoreeCloud/goreecloud-bookmarks/releases"
                target="_blank"
                className="underline decoration-dotted underline-offset-4 hover:text-primary duration-100"
              >
                GoreeCloud Bookmarks {announcementId}
              </Link>
            </>
          ) : (
            announcementMessage
          )}
        </p>
        <Button
          variant="ghost"
          size="icon"
          onClick={toggleAnnouncementBar}
          aria-label="Dismiss announcement"
        >
          <i className="bi-x text-xl" aria-hidden="true"></i>
        </Button>
      </div>
    </div>
  );
}
