import type { Attachment } from "svelte/attachments";

export const pwaHack = (height: string) => (e: Element) => {
  if (e instanceof HTMLElement) {
    setTimeout(() => (e.style.height = height), 10);
  }
};
