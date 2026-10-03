// Remove registrations left behind by earlier PWA releases.
export const unregisterServiceWorkers = async () => {
  if (!("serviceWorker" in navigator)) return;

  try {
    const registrations = await navigator.serviceWorker.getRegistrations();
    await Promise.all(registrations.map((registration) => registration.unregister()));
  } catch (error) {
    console.warn("Failed to unregister service workers", error);
  }
};
