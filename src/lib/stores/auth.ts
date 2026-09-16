import { writable } from "svelte/store";
import { browser } from "$app/environment";

export interface User {
  username: string;
  id: string;
}

function createAuthStore() {
  const { subscribe, set } = writable<{
    user: User | null;
    token: string | null;
    isAuthenticated: boolean;
  }>({
    user: null,
    token: null,
    isAuthenticated: false,
  });

  return {
    subscribe,
    login: (user: User, token: string) => {
      if (browser) {
        localStorage.setItem("token", token);
        localStorage.setItem("user", JSON.stringify(user));
      }
      set({ user, token, isAuthenticated: true });
    },
    logout: () => {
      if (browser) {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
      }
      set({ user: null, token: null, isAuthenticated: false });
    },
    initialize: () => {
      if (browser) {
        const token = localStorage.getItem("token");
        const userStr = localStorage.getItem("user");
        if (token && userStr) {
          try {
            const user = JSON.parse(userStr);
            // Optional: Validate token via API here if strict
            set({ user, token, isAuthenticated: true });
          } catch (e) {
            localStorage.removeItem("user");
            localStorage.removeItem("token");
            set({ user: null, token: null, isAuthenticated: false });
          }
        }
      }
    },
  };
}

export const auth = createAuthStore();
