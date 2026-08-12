export default async function getLatestVersion(
  setShowAnnouncement: Function
) {
  // GoreeCloud Bookmarks does not poll the upstream Linkwarden announcement feed.
  // Remove any announcement metadata left by an earlier upstream-derived session so
  // stale Linkwarden release messages cannot reappear in the GoreeCloud interface.
  localStorage.removeItem("announcementId");
  localStorage.removeItem("announcementMessage");
  setShowAnnouncement(false);

  // A GoreeCloud-controlled release/announcement source can be added later.
  return;
}
