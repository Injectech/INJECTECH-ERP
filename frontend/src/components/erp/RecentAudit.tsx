import ComponentCard from "../common/ComponentCard";
import Badge from "../ui/badge/Badge";

const logs = [
  {
    id: "AUD-091",
    actor: "Nadia Lestari",
    action: "user.create",
    detail: "Created account for bimo@erp.local",
    severity: "info",
    time: "10:24",
  },
  {
    id: "AUD-092",
    actor: "System",
    action: "inventory.adjust",
    detail: "Thermal Receipt Paper -12",
    severity: "warning",
    time: "10:41",
  },
  {
    id: "AUD-093",
    actor: "Rafi Pratama",
    action: "product.update",
    detail: "Updated pricing PRD-1001",
    severity: "info",
    time: "11:05",
  },
  {
    id: "AUD-094",
    actor: "System",
    action: "auth.login.failed",
    detail: "tania@erp.local",
    severity: "error",
    time: "11:17",
  },
];

export default function RecentAudit() {
  return (
    <ComponentCard title="Recent Audit Logs" desc="Latest security events.">
      <div className="space-y-4">
        {logs.map((log) => (
          <div
            key={log.id}
            className="flex flex-col gap-3 rounded-xl border border-gray-100 bg-gray-50/60 p-4 sm:flex-row sm:items-center sm:justify-between"
          >
            <div>
              <p className="text-sm font-semibold text-gray-900">
                {log.action}
              </p>
              <p className="mt-1 text-xs text-gray-500">{log.detail}</p>
              <p className="mt-2 text-xs text-gray-400">
                {log.actor} · {log.time}
              </p>
            </div>
            <Badge
              size="sm"
              color={
                log.severity === "error"
                  ? "error"
                  : log.severity === "warning"
                  ? "warning"
                  : "info"
              }
            >
              {log.severity}
            </Badge>
          </div>
        ))}
      </div>
    </ComponentCard>
  );
}
