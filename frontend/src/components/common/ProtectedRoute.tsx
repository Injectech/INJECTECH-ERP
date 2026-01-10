import type { ReactNode } from "react";
import { Navigate } from "react-router";
import { useAuthStore } from "../../stores/authStore";

export default function ProtectedRoute({
  children,
}: {
  children: ReactNode;
}) {
  const accessToken = useAuthStore((state) => state.accessToken);

  if (!accessToken) {
    return <Navigate to="/signin" replace />;
  }

  return children;
}
