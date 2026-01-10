import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import PageBreadcrumb from "../../components/common/PageBreadCrumb";
import PageMeta from "../../components/common/PageMeta";
import ComponentCard from "../../components/common/ComponentCard";
import Input from "../../components/form/input/InputField";
import Select from "../../components/form/Select";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "../../components/ui/table";
import { fetchAuditLogs } from "../../services/erp";

const severityOptions = [
  { value: "all", label: "All severity" },
  { value: "info", label: "Info" },
  { value: "warning", label: "Warning" },
  { value: "error", label: "Error" },
];

export default function AuditLogs() {
  const [actorId, setActorId] = useState("");
  const [date, setDate] = useState("");

  const logsQuery = useQuery({
    queryKey: ["audit-logs", actorId],
    queryFn: () => fetchAuditLogs(actorId || undefined),
  });

  return (
    <>
      <PageMeta
        title="ERP Audit Logs | Dashboard"
        description="Audit logs page for ERP backend."
      />
      <PageBreadcrumb pageTitle="Audit Logs" />

      <div className="space-y-6">
        <ComponentCard
          title="Audit Trail"
          desc="Review critical actions and security events."
        >
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div className="grid w-full gap-3 md:grid-cols-3">
              <Input
                placeholder="Actor ID (optional)"
                value={actorId}
                onChange={(e) => setActorId(e.target.value)}
              />
              <Input
                type="date"
                value={date}
                onChange={(e) => setDate(e.target.value)}
              />
              <Select
                options={severityOptions}
                placeholder="Severity"
                onChange={() => {}}
              />
            </div>
            <button className="inline-flex items-center justify-center rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-semibold text-gray-700 shadow-theme-xs hover:bg-gray-50">
              Export Logs
            </button>
          </div>

          {logsQuery.isLoading && (
            <p className="text-sm text-gray-500">Loading logs...</p>
          )}
          {logsQuery.error && (
            <p className="text-sm text-error-500">
              {(logsQuery.error as Error).message}
            </p>
          )}
          {logsQuery.data && (
            <div className="overflow-hidden rounded-xl border border-gray-200 bg-white">
              <div className="max-w-full overflow-x-auto">
                <Table>
                  <TableHeader className="border-b border-gray-100">
                    <TableRow>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        Actor ID
                      </TableCell>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        Action
                      </TableCell>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        Resource
                      </TableCell>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        Time
                      </TableCell>
                    </TableRow>
                  </TableHeader>
                  <TableBody className="divide-y divide-gray-100">
                    {logsQuery.data.map((log) => (
                      <TableRow key={log.id}>
                        <TableCell className="px-5 py-4 text-sm font-semibold text-gray-900">
                          {log.actor_id || "system"}
                        </TableCell>
                        <TableCell className="px-5 py-4 text-sm text-gray-600">
                          {log.action}
                        </TableCell>
                        <TableCell className="px-5 py-4 text-sm text-gray-600">
                          {log.resource}
                        </TableCell>
                        <TableCell className="px-5 py-4 text-sm text-gray-500">
                          {new Date(log.created_at).toLocaleString("id-ID")}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            </div>
          )}
        </ComponentCard>
      </div>
    </>
  );
}
