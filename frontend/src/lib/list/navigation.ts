import { navigate } from "../../router";
import { listStore } from "./store";
import { type ListSummary } from "./store/store";

export const navigateToDefaultList = async () => {
  const lists = await new Promise<ListSummary[]>((resolve) =>
    listStore.lists.subscribe(resolve),
  );
  const defaultList = lists.at(0);

  if (lists.length === 1 && defaultList) {
    return navigate("/lists/:id", {
      params: { id: defaultList.id },
    });
  } else {
    return navigate("/lists", { replace: true });
  }
};
