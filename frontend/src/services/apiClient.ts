import { useAuthStore } from "../stores/authStore";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

type ApiEnvelope<T> = {
  success: boolean;
  message: string;
  data: T;
};

type RequestOptions = {
  method?: string;
  body?: unknown;
  auth?: boolean;
  retry?: boolean;
};

let refreshPromise: Promise<ApiEnvelope<{ access_token: string; access_expires_at: string }>> | null =
  null;

async function parseJson<T>(res: Response): Promise<T> {
  const text = await res.text();
  if (!text) {
    throw new Error("Empty response from server");
  }
  try {
    return JSON.parse(text) as T;
  } catch {
    throw new Error(text);
  }
}

async function refreshToken() {
  if (!refreshPromise) {
    refreshPromise = fetch(`${API_BASE_URL}/api/v1/auth/refresh`, {
      method: "POST",
      credentials: "include",
    }).then(async (res) => {
      const payload = await parseJson<ApiEnvelope<{
        access_token: string;
        access_expires_at: string;
      }>>(res);
      if (!res.ok || !payload.success) {
        throw new Error(payload?.message || "Unauthorized");
      }
      return payload;
    });
  }
  try {
    return await refreshPromise;
  } finally {
    refreshPromise = null;
  }
}

export async function apiRequest<T>(path: string, options: RequestOptions = {}) {
  const { method = "GET", body, auth = true, retry = true } = options;
  const accessToken = useAuthStore.getState().accessToken;
  const headers: Record<string, string> = {};

  if (body) {
    headers["Content-Type"] = "application/json";
  }

  if (auth && accessToken) {
    headers.Authorization = `Bearer ${accessToken}`;
  }

  const res = await fetch(`${API_BASE_URL}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
    credentials: "include",
  });

  if (auth && res.status === 401 && retry) {
    try {
      const refreshed = await refreshToken();
      useAuthStore
        .getState()
        .setAccessToken(
          refreshed.data.access_token,
          refreshed.data.access_expires_at,
        );
      return apiRequest<T>(path, { ...options, retry: false });
    } catch (err) {
      useAuthStore.getState().clear();
      throw err;
    }
  }

  const payload = await parseJson<ApiEnvelope<T>>(res);
  if (!res.ok || !payload.success) {
    throw new Error(payload?.message || "Request failed");
  }
  return payload.data;
}
