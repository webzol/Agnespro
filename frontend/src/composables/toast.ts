// Lightweight global toast system
import { reactive } from "vue";

export interface Toast { id: number; text: string; kind: "ok"|"err"|"info"; }
const state = reactive({ items: [] as Toast[] });
let nextId = 1;

export function useToast() {
  function push(text: string, kind: Toast["kind"] = "info") {
    const id = nextId++;
    state.items.push({ id, text, kind });
    setTimeout(() => {
      const i = state.items.findIndex(t => t.id === id);
      if (i >= 0) state.items.splice(i, 1);
    }, 3500);
  }
  return {
    state,
    ok:   (text: string) => push(text, "ok"),
    err:  (text: string) => push(text, "err"),
    info: (text: string) => push(text, "info"),
  };
}
