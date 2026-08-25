import React from "react";
import importBookmarks from "@/lib/client/importBookmarks";
import { MigrationFormat } from "@linkwarden/types/global";
import { useTranslation } from "next-i18next";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
} from "@/components/ui/dropdown-menu";
import { Button } from "./ui/button";

const ImportDropdown = () => {
  const { t } = useTranslation();

  const importOptions = [
    {
      id: "import-linkwarden-file",
      format: MigrationFormat.linkwarden,
      label: t("from_linkwarden"),
      accept: ".json,application/json",
    },
    {
      id: "import-html-file",
      format: MigrationFormat.htmlFile,
      label: t("from_html"),
      accept: ".html,text/html",
    },
    {
      id: "import-pocket-file",
      format: MigrationFormat.pocket,
      label: t("from_pocket"),
      accept: ".csv,text/csv",
    },
    {
      id: "import-wallabag-file",
      format: MigrationFormat.wallabag,
      label: t("from_wallabag"),
      accept: ".json,application/json",
    },
    {
      id: "import-omnivore-file",
      format: MigrationFormat.omnivore,
      label: t("from_omnivore"),
      accept: ".zip,application/zip",
    },
  ];

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="metal">
          <i className="bi-cloud-upload text-xl" aria-hidden="true"></i>
          {t("import_links")}
        </Button>
      </DropdownMenuTrigger>

      <DropdownMenuContent side="bottom" align="start">
        {importOptions.map((item) => (
          <DropdownMenuItem
            asChild
            key={item.id}
            onSelect={(event) => event.preventDefault()}
          >
            <label
              htmlFor={item.id}
              className="w-full cursor-pointer whitespace-nowrap"
            >
              {item.label}
              <input
                type="file"
                id={item.id}
                accept={item.accept}
                className="hidden"
                aria-label={`${t("import_links")}: ${item.label}`}
                onChange={(event) => {
                  void importBookmarks(event, item.format);
                  event.currentTarget.value = "";
                }}
              />
            </label>
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default ImportDropdown;
