import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

export type AuthUser = {
  id: string;
  name: string;
  email: string;
  roles: string[];
};

interface AuthState {
  accessToken: string | null;
  accessExpiresAt: string | null;
  user: AuthUser | null;
  setSession: (token: string, expiresAt: string, user: AuthUser) => void;
  setAccessToken: (token: string, expiresAt: string) => void;
  clear: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      accessToken: null,
      accessExpiresAt: null,
      user: null,
      setSession: (token, expiresAt, user) =>
        set({ accessToken: token, accessExpiresAt: expiresAt, user }),
      setAccessToken: (token, expiresAt) =>
        set((state) => ({ ...state, accessToken: token, accessExpiresAt: expiresAt })),
      clear: () => set({ accessToken: null, accessExpiresAt: null, user: null }),
    }),
    {
      name: "erp-auth",
      storage: createJSONStorage(() => sessionStorage),
    },
  ),
);
