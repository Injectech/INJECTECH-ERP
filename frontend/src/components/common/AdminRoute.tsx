import type { ReactNode } from "react";
import { Navigate } from "react-router";
import { useAuthStore } from "../../stores/authStore";

export default function AdminRoute({
  children,
}: {
  children: ReactNode;
}) {
  const accessToken = useAuthStore((state) => state.accessToken);
  const user = useAuthStore((state) => state.user);
  const isAdmin = Boolean(user?.roles?.includes("admin"));

  if (!accessToken) {
    return <Navigate to="/signin" replace />;
  }

  if (!isAdmin) {
    return <Navigate to="/" replace />;
  }

  return children;
}
