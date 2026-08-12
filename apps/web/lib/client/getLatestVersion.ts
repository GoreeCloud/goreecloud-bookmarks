export default async function getLatestVersion(
  setShowAnnouncement: Function
) {
  // GoreeCloud Bookmarks does not poll the upstream Linkwarden announcement feed.
  // A GoreeCloud-controlled release/announcement source can be added later.
  void setShowAnnouncement;
  return;
}
