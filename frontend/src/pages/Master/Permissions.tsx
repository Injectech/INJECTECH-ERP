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
import { fetchPermissions } from "../../services/erp";

const moduleOptions = [
  { value: "all", label: "All modules" },
  { value: "security", label: "Security" },
  { value: "catalog", label: "Catalog" },
  { value: "inventory", label: "Inventory" },
  { value: "audit", label: "Audit" },
];

export default function Permissions() {
  const permissionsQuery = useQuery({
    queryKey: ["permissions"],
    queryFn: fetchPermissions,
  });

  return (
    <>
      <PageMeta
        title="ERP Permissions | Dashboard"
        description="Permission management page for ERP backend."
      />
      <PageBreadcrumb pageTitle="Permissions" />

      <div className="space-y-6">
        <ComponentCard
          title="Permission Catalog"
          desc="All capabilities available in the ERP system."
        >
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div className="grid w-full gap-3 md:grid-cols-2">
              <Input placeholder="Search permission code" />
              <Select
                options={moduleOptions}
                placeholder="Module"
                onChange={() => {}}
              />
            </div>
            <button className="inline-flex items-center justify-center rounded-lg bg-brand-500 px-4 py-2 text-sm font-semibold text-white shadow-theme-sm hover:bg-brand-600">
              Add Permission
            </button>
          </div>

          {permissionsQuery.isLoading && (
            <p className="text-sm text-gray-500">Loading permissions...</p>
          )}
          {permissionsQuery.error && (
            <p className="text-sm text-error-500">
              {(permissionsQuery.error as Error).message}
            </p>
          )}
          {permissionsQuery.data && (
            <div className="overflow-hidden rounded-xl border border-gray-200 bg-white">
              <div className="max-w-full overflow-x-auto">
                <Table>
                  <TableHeader className="border-b border-gray-100">
                    <TableRow>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        Permission Code
                      </TableCell>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        Description
                      </TableCell>
                    </TableRow>
                  </TableHeader>
                  <TableBody className="divide-y divide-gray-100">
                    {permissionsQuery.data.map((permission) => (
                      <TableRow key={permission.id}>
                        <TableCell className="px-5 py-4 text-sm font-semibold text-gray-900">
                          {permission.code}
                        </TableCell>
                        <TableCell className="px-5 py-4 text-sm text-gray-600">
                          {permission.description}
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
