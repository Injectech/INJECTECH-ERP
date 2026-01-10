import { apiRequest } from "./apiClient";

export type AuthData = {
  access_token: string;
  access_expires_at: string;
  user: {
    id: string;
    name: string;
    email: string;
    roles: string[];
  };
};

export type LoginPayload = {
  email: string;
  password: string;
};

export type RegisterPayload = {
  email: string;
  password: string;
  name: string;
};

export function login(payload: LoginPayload) {
  return apiRequest<AuthData>("/api/v1/auth/login", {
    method: "POST",
    body: payload,
    auth: false,
  });
}

export function register(payload: RegisterPayload) {
  return apiRequest<AuthData>("/api/v1/auth/register", {
    method: "POST",
    body: payload,
    auth: false,
  });
}
