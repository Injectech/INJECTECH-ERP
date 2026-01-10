import { useQuery } from "@tanstack/react-query";
import { BoxCubeIcon, TableIcon, TaskIcon, UserCircleIcon } from "../../icons";
import Badge from "../ui/badge/Badge";
import {
  fetchAuditLogs,
  fetchPermissions,
  fetchProducts,
  fetchRoles,
} from "../../services/erp";

export default function ErpMetrics() {
  const rolesQuery = useQuery({ queryKey: ["roles"], queryFn: fetchRoles });
  const permissionsQuery = useQuery({
    queryKey: ["permissions"],
    queryFn: fetchPermissions,
  });
  const productsQuery = useQuery({
    queryKey: ["products"],
    queryFn: fetchProducts,
  });
  const auditQuery = useQuery({
    queryKey: ["audit-logs", "recent"],
    queryFn: () => fetchAuditLogs(),
  });

  const metrics = [
    {
      title: "Roles",
      value: rolesQuery.data?.length ?? 0,
      trend: "RBAC",
      tone: "primary" as const,
      note: "Defined roles",
      icon: TaskIcon,
    },
    {
      title: "Permissions",
      value: permissionsQuery.data?.length ?? 0,
      trend: "Access",
      tone: "light" as const,
      note: "Policy rules",
      icon: UserCircleIcon,
    },
    {
      title: "Products",
      value: productsQuery.data?.length ?? 0,
      trend: "Catalog",
      tone: "success" as const,
      note: "Active SKUs",
      icon: BoxCubeIcon,
    },
    {
      title: "Audit Events",
      value: auditQuery.data?.length ?? 0,
      trend: "Logs",
      tone: "warning" as const,
      note: "Latest events",
      icon: TableIcon,
    },
  ];

  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
      {metrics.map((item) => (
        <div
          key={item.title}
          className="rounded-2xl border border-gray-200 bg-white p-5 shadow-theme-sm"
        >
          <div className="flex items-center justify-between">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-brand-50 text-brand-500">
              <item.icon className="size-6" />
            </div>
            <Badge size="sm" color={item.tone}>
              {item.trend}
            </Badge>
          </div>
          <div className="mt-4">
            <p className="text-sm text-gray-500">{item.title}</p>
            <p className="mt-2 text-2xl font-semibold text-gray-900">
              {item.value}
            </p>
            <p className="mt-1 text-xs text-gray-500">{item.note}</p>
          </div>
        </div>
      ))}
    </div>
  );
}
