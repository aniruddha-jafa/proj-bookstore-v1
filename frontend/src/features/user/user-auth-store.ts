import { create } from "zustand";
import { UserResponse } from "@/features/user/user-schema";
import { refreshToken } from "@/features/auth/auth-api";

interface UserStore {
    user: UserResponse | null;
    accessToken: string | null;
    isLoading: boolean;
    setUser: (user: UserResponse) => void;
    setAccessToken: (accessToken: string) => void;
    refreshSession: () => Promise<void>;
   logout: () => void;
}

export const useUserStore = create<UserStore>((set) => ({
    user: null,
    accessToken: null,
    isLoading: true,
    setUser: (user: UserResponse) => set({ user, isLoading: false }),
    setAccessToken: (accessToken: string) => set({ accessToken }),
    refreshSession: async () => {
        const res = await refreshToken();
        if (!res.ok) {
            set({ user: null, accessToken: null, isLoading: false });
            return;
        }
        set({ accessToken: res.data.token, isLoading: true });
    },
    logout: () => set({ user: null, accessToken: null }),
}));
