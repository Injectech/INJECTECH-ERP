import { useMemo, useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { Controller, useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import PageBreadcrumb from "../../components/common/PageBreadCrumb";
import PageMeta from "../../components/common/PageMeta";
import ComponentCard from "../../components/common/ComponentCard";
import Badge from "../../components/ui/badge/Badge";
import Input from "../../components/form/input/InputField";
import Select from "../../components/form/Select";
import Label from "../../components/form/Label";
import TextArea from "../../components/form/input/TextArea";
import MultiSelect from "../../components/form/MultiSelect";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "../../components/ui/table";
import { createRole, fetchPermissions, fetchRoles } from "../../services/erp";
import { useAuthStore } from "../../stores/authStore";

const moduleOptions = [
  { value: "all", label: "All modules" },
  { value: "security", label: "Security" },
  { value: "inventory", label: "Inventory" },
  { value: "finance", label: "Finance" },
];

const createSchema = z.object({
  name: z.string().min(2, "Nama role minimal 2 karakter"),
  description: z.string().min(3, "Deskripsi minimal 3 karakter"),
});

type CreateForm = z.infer<typeof createSchema>;

export default function Roles() {
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [selectedPermissions, setSelectedPermissions] = useState<string[]>([]);
  const [permissionError, setPermissionError] = useState<string | null>(null);
  const user = useAuthStore((state) => state.user);
  const isAdmin = Boolean(user?.roles?.includes("admin"));

  const rolesQuery = useQuery({
    queryKey: ["roles"],
    queryFn: fetchRoles,
  });
  const permissionsQuery = useQuery({
    queryKey: ["permissions"],
    queryFn: fetchPermissions,
  });

  const permissionOptions = useMemo(
    () =>
      (permissionsQuery.data || []).map((permission) => ({
        value: permission.code,
        text: permission.code,
      })),
    [permissionsQuery.data],
  );

  const createForm = useForm<CreateForm>({
    resolver: zodResolver(createSchema),
    defaultValues: {
      description: "",
    },
  });

  const createMutation = useMutation({
    mutationFn: createRole,
    onSuccess: () => {
      createForm.reset();
      setSelectedPermissions([]);
      setShowCreateForm(false);
      rolesQuery.refetch();
    },
  });

  const onCreate = (values: CreateForm) => {
    if (selectedPermissions.length === 0) {
      setPermissionError("Pilih minimal 1 permission");
      return;
    }
    setPermissionError(null);
    createMutation.mutate({
      name: values.name.trim(),
      description: values.description.trim(),
      permissions: selectedPermissions,
    });
  };

  return (
    <>
      <PageMeta
        title="ERP Roles | Dashboard"
        description="Role management page for ERP backend."
      />
      <PageBreadcrumb pageTitle="Roles & Permissions" />

      <div className="space-y-6">
        <div className="grid gap-4 md:grid-cols-3">
          <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-theme-sm">
            <p className="text-sm text-gray-500">Total Roles</p>
            <p className="mt-2 text-2xl font-semibold text-gray-900">
              {rolesQuery.data?.length ?? 0}
            </p>
            <p className="mt-1 text-xs text-gray-500">Termasuk role sistem</p>
          </div>
          <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-theme-sm">
            <p className="text-sm text-gray-500">Critical Permissions</p>
            <p className="mt-2 text-2xl font-semibold text-gray-900">5</p>
            <p className="mt-1 text-xs text-gray-500">Requires approval</p>
          </div>
          <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-theme-sm">
            <p className="text-sm text-gray-500">Users Assigned</p>
            <p className="mt-2 text-2xl font-semibold text-gray-900">68</p>
            <p className="mt-1 text-xs text-gray-500">Across all roles</p>
          </div>
        </div>

        <ComponentCard
          title="Role Directory"
          desc="Design access boundaries by assigning permissions to roles."
        >
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div className="grid w-full gap-3 md:grid-cols-2">
              <Input placeholder="Search role name" />
              <Select
                options={moduleOptions}
                placeholder="Filter module"
                onChange={() => {}}
              />
            </div>
            {isAdmin && (
              <button
                className="inline-flex items-center justify-center rounded-lg bg-brand-500 px-4 py-2 text-sm font-semibold text-white shadow-theme-sm hover:bg-brand-600"
                onClick={() => setShowCreateForm((prev) => !prev)}
              >
                {showCreateForm ? "Tutup Form" : "Create Role"}
              </button>
            )}
          </div>

          {showCreateForm && isAdmin && (
            <div className="mt-6 rounded-xl border border-gray-200 bg-white p-5">
              <form
                className="grid gap-4 md:grid-cols-2"
                onSubmit={createForm.handleSubmit(onCreate)}
              >
                <div>
                  <Label>
                    Nama Role <span className="text-error-500">*</span>
                  </Label>
                  <Input
                    placeholder="contoh: manager"
                    error={Boolean(createForm.formState.errors.name)}
                    {...createForm.register("name")}
                  />
                  {createForm.formState.errors.name && (
                    <p className="mt-1 text-xs text-error-500">
                      {createForm.formState.errors.name.message}
                    </p>
                  )}
                </div>
                <div className="md:col-span-2">
                  <Label>
                    Deskripsi <span className="text-error-500">*</span>
                  </Label>
                  <Controller
                    control={createForm.control}
                    name="description"
                    render={({ field }) => (
                      <TextArea
                        placeholder="Deskripsi singkat role"
                        rows={3}
                        error={Boolean(createForm.formState.errors.description)}
                        onChange={field.onChange}
                        value={field.value || ""}
                      />
                    )}
                  />
                  {createForm.formState.errors.description && (
                    <p className="mt-1 text-xs text-error-500">
                      {createForm.formState.errors.description.message}
                    </p>
                  )}
                </div>
                <div className="md:col-span-2">
                  <MultiSelect
                    label="Permissions"
                    options={permissionOptions}
                    value={selectedPermissions}
                    onChange={setSelectedPermissions}
                    placeholder="Pilih permission"
                    disabled={permissionsQuery.isLoading}
                  />
                  {permissionError && (
                    <p className="mt-1 text-xs text-error-500">
                      {permissionError}
                    </p>
                  )}
                  {permissionsQuery.error && (
                    <p className="mt-1 text-xs text-error-500">
                      {(permissionsQuery.error as Error).message}
                    </p>
                  )}
                </div>
                <div className="md:col-span-2 flex items-center gap-3">
                  <button
                    className="inline-flex items-center justify-center rounded-lg bg-brand-500 px-4 py-2 text-sm font-semibold text-white shadow-theme-sm hover:bg-brand-600"
                    type="submit"
                    disabled={createMutation.isPending}
                  >
                    {createMutation.isPending ? "Menyimpan..." : "Simpan Role"}
                  </button>
                  {createMutation.isError && (
                    <p className="text-sm text-error-500">
                      {(createMutation.error as Error).message}
                    </p>
                  )}
                  {createMutation.isSuccess && (
                    <p className="text-sm text-success-500">
                      Role berhasil dibuat.
                    </p>
                  )}
                </div>
              </form>
            </div>
          )}

          {rolesQuery.isLoading && (
            <p className="text-sm text-gray-500">Loading roles...</p>
          )}
          {rolesQuery.error && (
            <p className="text-sm text-error-500">
              {(rolesQuery.error as Error).message}
            </p>
          )}
          {rolesQuery.data && (
            <div className="overflow-hidden rounded-xl border border-gray-200 bg-white">
              <div className="max-w-full overflow-x-auto">
                <Table>
                  <TableHeader className="border-b border-gray-100">
                    <TableRow>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        Role
                      </TableCell>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        Permissions
                      </TableCell>
                    </TableRow>
                  </TableHeader>
                  <TableBody className="divide-y divide-gray-100">
                    {rolesQuery.data.map((role) => (
                      <TableRow key={role.id}>
                        <TableCell className="px-5 py-4 text-start">
                          <p className="text-sm font-semibold text-gray-900">
                            {role.name}
                          </p>
                          <p className="text-xs text-gray-500">
                            {role.description}
                          </p>
                        </TableCell>
                        <TableCell className="px-5 py-4 text-sm text-gray-600">
                          <Badge size="sm" color="primary">
                            {role.permissions?.length || 0} perms
                          </Badge>
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
