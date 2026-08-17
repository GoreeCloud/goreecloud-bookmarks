import { MigrationFormat, MigrationRequest } from "@linkwarden/types/global";
import { toast } from "react-hot-toast";
import JSZip from "jszip";

const processOmnivoreZipFile = async (zip: JSZip): Promise<string> => {
  const metadataFiles = Object.keys(zip.files).filter((filePath) => {
    const file = zip.files[filePath];
    return filePath.startsWith("metadata_") && !file.dir;
  });

  if (metadataFiles.length === 0) {
    throw new Error("No Omnivore metadata files were found in the archive.");
  }

  const allMetadataArrays = await Promise.all(
    metadataFiles.map(async (filePath) => {
      const fileContent = await zip.files[filePath].async("string");
      const parsedContent = JSON.parse(fileContent);

      if (!Array.isArray(parsedContent)) {
        throw new Error(`Invalid Omnivore metadata in ${filePath}.`);
      }

      return parsedContent;
    })
  );

  return JSON.stringify(allMetadataArrays.flat());
};

const readErrorMessage = async (response: Response): Promise<string> => {
  try {
    const errorData = await response.json();
    if (typeof errorData?.response === "string" && errorData.response.trim()) {
      return errorData.response;
    }
  } catch (error) {
    console.error("Failed to parse import error response", error);
  }

  return "Failed to import bookmarks. Please try again.";
};

const importBookmarks = async (
  e: React.ChangeEvent<HTMLInputElement>,
  format: MigrationFormat
) => {
  const file = e.target.files?.[0];

  if (!file) return;

  const reader = new FileReader();

  if (format === MigrationFormat.omnivore) reader.readAsArrayBuffer(file);
  else reader.readAsText(file, "UTF-8");

  reader.onload = async (event) => {
    const load = toast.loading("Importing bookmarks...");

    try {
      let request = event.target?.result;

      if (request === null || request === undefined) {
        throw new Error("The selected import file could not be read.");
      }

      if (format === MigrationFormat.omnivore) {
        if (!(request instanceof ArrayBuffer)) {
          throw new Error("The Omnivore archive could not be read as binary data.");
        }

        const zip = await JSZip.loadAsync(request);
        request = await processOmnivoreZipFile(zip);
      }

      if (typeof request !== "string") {
        throw new Error("The selected import file could not be read as text.");
      }

      const body: MigrationRequest = {
        format,
        data: request,
      };

      const response = await fetch("/api/v1/migration", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
      });

      if (!response.ok) {
        toast.error(await readErrorMessage(response));
        return;
      }

      await response.json();
      toast.success("Bookmarks imported successfully. Reloading the page...");

      setTimeout(() => {
        location.reload();
      }, 2000);
    } catch (error) {
      console.error("Bookmark import failed", error);
      toast.error(
        error instanceof Error && error.message
          ? error.message
          : "An error occurred while importing bookmarks. Please try again."
      );
    } finally {
      toast.dismiss(load);
    }
  };

  reader.onerror = (error) => {
    console.error("Failed to read bookmark import file", error);
    toast.error(
      "Failed to read the file. Please make sure the file is correct and try again."
    );
  };
};

export default importBookmarks;
