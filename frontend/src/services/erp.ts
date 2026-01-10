import { apiRequest } from "./apiClient";

export type User = {
  id: string;
  email: string;
  name: string;
  roles: string[];
  permissions?: string[];
  created_at: string;
  updated_at: string;
};

export type Role = {
  id: string;
  name: string;
  description: string;
  permissions: string[];
  created_at: string;
};

export type Permission = {
  id: string;
  code: string;
  description: string;
  created_at: string;
};

export type Product = {
  id: string;
  sku: string;
  name: string;
  description: string;
  price: number;
  created_at: string;
};

export type Inventory = {
  id: string;
  product_id: string;
  quantity: number;
  location: string;
  created_at: string;
};

export type Location = {
  id: string;
  name: string;
  description: string;
  is_default: boolean;
  created_at: string;
};

export type AuditLog = {
  id: string;
  actor_id: string;
  action: string;
  resource: string;
  metadata: Record<string, unknown>;
  created_at: string;
};

export type CreateUserPayload = {
  email: string;
  password: string;
  name: string;
  roles: string[];
};

export type CreateProductPayload = {
  sku: string;
  name: string;
  description: string;
  price: number;
};

export type CreateInventoryPayload = {
  product_id: string;
  quantity: number;
  location: string;
};

export type CreateLocationPayload = {
  name: string;
  description: string;
};

export type CreateRolePayload = {
  name: string;
  description: string;
  permissions: string[];
};

export function fetchUser(id: string) {
  return apiRequest<User>(`/api/v1/users/${id}`);
}

export function createUser(payload: CreateUserPayload) {
  return apiRequest<User>("/api/v1/users", {
    method: "POST",
    body: payload,
  });
}

export function fetchRoles() {
  return apiRequest<Role[]>("/api/v1/roles");
}

export function createRole(payload: CreateRolePayload) {
  return apiRequest<Role>("/api/v1/roles", {
    method: "POST",
    body: payload,
  });
}

export function fetchPermissions() {
  return apiRequest<Permission[]>("/api/v1/permissions");
}

export function fetchProducts() {
  return apiRequest<Product[]>("/api/v1/products");
}

export function createProduct(payload: CreateProductPayload) {
  return apiRequest<Product>("/api/v1/products", {
    method: "POST",
    body: payload,
  });
}

export type UpdateProductPayload = CreateProductPayload & { id: string };

export function updateProduct(payload: UpdateProductPayload) {
  return apiRequest<Product>(`/api/v1/products/${payload.id}`, {
    method: "PUT",
    body: {
      sku: payload.sku,
      name: payload.name,
      description: payload.description,
      price: payload.price,
    },
  });
}

export function deleteProduct(id: string) {
  return apiRequest<{ id: string }>(`/api/v1/products/${id}`, {
    method: "DELETE",
  });
}

export function createInventory(payload: CreateInventoryPayload) {
  return apiRequest<Inventory>("/api/v1/inventory", {
    method: "POST",
    body: payload,
  });
}

export function fetchInventory(productId: string) {
  return apiRequest<Inventory[]>(`/api/v1/inventory/product/${productId}`);
}

export function fetchAllInventory() {
  return apiRequest<Inventory[]>("/api/v1/inventory");
}

export function adjustInventory(id: string, delta: number) {
  return apiRequest<void>(`/api/v1/inventory/${id}/adjust`, {
    method: "PATCH",
    body: { delta },
  });
}

export function updateInventoryLocation(id: string, location: string) {
  return apiRequest<void>(`/api/v1/inventory/${id}/location`, {
    method: "PATCH",
    body: { location },
  });
}

export function fetchLocations() {
  return apiRequest<Location[]>("/api/v1/locations");
}

export function createLocation(payload: CreateLocationPayload) {
  return apiRequest<Location>("/api/v1/locations", {
    method: "POST",
    body: payload,
  });
}

export function fetchAuditLogs(actorId?: string) {
  const query = actorId ? `?actor_id=${encodeURIComponent(actorId)}` : "";
  return apiRequest<AuditLog[]>(`/api/v1/audit/logs${query}`);
}
