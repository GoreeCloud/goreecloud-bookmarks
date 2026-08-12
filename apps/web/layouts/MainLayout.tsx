import Announcement from "@/components/Announcement";
import Sidebar from "@/components/Sidebar";
import { ReactNode, useEffect, useState } from "react";
import getLatestVersion from "@/lib/client/getLatestVersion";
import DragNDrop from "@/components/DragNDrop";
import { LinkIncludingShortenedCollectionAndTags } from "@linkwarden/types/global";
import useSidebarCollapse from "@/hooks/useSidebarCollapse";
import { cn } from "@/lib/utils";

interface Props {
  children: ReactNode;
}

export default function MainLayout({ children }: Props) {
  const showAnnouncementBar = localStorage.getItem("showAnnouncementBar");

  const [showAnnouncement, setShowAnnouncement] = useState(
    showAnnouncementBar === "true"
  );
  const { sidebarIsCollapsed, toggleSidebar } = useSidebarCollapse();

  useEffect(() => {
    getLatestVersion(setShowAnnouncement);
  }, []);

  useEffect(() => {
    localStorage.setItem(
      "showAnnouncementBar",
      showAnnouncement ? "true" : "false"
    );
  }, [showAnnouncement]);

  const toggleAnnouncementBar = () => setShowAnnouncement(!showAnnouncement);

  const [activeLink, setActiveLink] =
    useState<LinkIncludingShortenedCollectionAndTags | null>(null);

  return (
    <DragNDrop activeLink={activeLink} setActiveLink={setActiveLink}>
      <div
        className="flex min-h-screen bg-base-200/50"
        data-testid="dashboard-wrapper"
      >
        {showAnnouncement && (
          <Announcement toggleAnnouncementBar={toggleAnnouncementBar} />
        )}
        <Sidebar
          toggleSidebar={toggleSidebar}
          sidebarIsCollapsed={sidebarIsCollapsed}
        />

        <main
          className={cn(
            "w-[calc(100%-56px)] min-w-0 h-screen overflow-hidden p-0 lg:py-2 lg:pr-2",
            !sidebarIsCollapsed && "lg:w-[calc(100%-288px)]"
          )}
        >
          <div className="h-full overflow-y-auto bg-base-100 sm:pb-0 pb-20 lg:rounded-2xl lg:border lg:border-base-content/10 lg:shadow-sm">
            {children}
          </div>
        </main>
      </div>
    </DragNDrop>
  );
}
