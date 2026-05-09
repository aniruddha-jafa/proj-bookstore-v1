"use client";

import { Toaster } from "sonner";

export function ToasterDismissable() {
  return (
    <Toaster
      richColors
      closeButton
      toastOptions={{
        closeButtonAriaLabel: "Dismiss notification",
      }}
    />
  );
}
