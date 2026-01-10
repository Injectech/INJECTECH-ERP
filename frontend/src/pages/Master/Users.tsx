import { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import PageBreadcrumb from "../../components/common/PageBreadCrumb";
import PageMeta from "../../components/common/PageMeta";
import ComponentCard from "../../components/common/ComponentCard";
import Input from "../../components/form/input/InputField";
import Badge from "../../components/ui/badge/Badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "../../components/ui/table";
import { createUser, fetchUser } from "../../services/erp";

const createSchema = z.object({
  name: z.string().min(2, "Nama minimal 2 karakter"),
  email: z.string().email("Email tidak valid"),
  password: z.string().min(6, "Password minimal 6 karakter"),
  roles: z.string().min(1, "Roles wajib diisi"),
});

type CreateForm = z.infer<typeof createSchema>;

export default function Users() {
  const [lookupId, setLookupId] = useState("");
  const [searchId, setSearchId] = useState("");

  const userQuery = useQuery({
    queryKey: ["user", searchId],
    queryFn: () => fetchUser(searchId),
    enabled: Boolean(searchId),
  });

  const createForm = useForm<CreateForm>({
    resolver: zodResolver(createSchema),
  });

  const createMutation = useMutation({
    mutationFn: createUser,
    onSuccess: () => {
      createForm.reset();
    },
  });

  const onCreate = (values: CreateForm) => {
    const roles = values.roles
      .split(",")
      .map((role) => role.trim())
      .filter(Boolean);
    createMutation.mutate({
      name: values.name,
      email: values.email,
      password: values.password,
      roles,
    });
  };

  return (
    <>
      <PageMeta
        title="ERP Users | Dashboard"
        description="User management page for ERP backend."
      />
      <PageBreadcrumb pageTitle="Users" />

      <div className="space-y-6">
        <div className="grid gap-4 md:grid-cols-3">
          <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-theme-sm">
            <p className="text-sm text-gray-500">Total Users</p>
            <p className="mt-2 text-2xl font-semibold text-gray-900">N/A</p>
            <p className="mt-1 text-xs text-gray-500">
              List endpoint belum tersedia
            </p>
          </div>
          <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-theme-sm">
            <p className="text-sm text-gray-500">Active Sessions</p>
            <p className="mt-2 text-2xl font-semibold text-gray-900">34</p>
            <p className="mt-1 text-xs text-gray-500">Peak at 11:00</p>
          </div>
          <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-theme-sm">
            <p className="text-sm text-gray-500">Roles Assigned</p>
            <p className="mt-2 text-2xl font-semibold text-gray-900">7</p>
            <p className="mt-1 text-xs text-gray-500">Admin, Manager, Staff</p>
          </div>
        </div>

        <ComponentCard
          title="Lookup User"
          desc="Cari user berdasarkan ID."
        >
          <div className="flex flex-col gap-4 md:flex-row md:items-center">
            <div className="flex-1">
              <Input
                placeholder="Masukkan User ID"
                value={lookupId}
                onChange={(e) => setLookupId(e.target.value)}
              />
            </div>
            <button
              className="inline-flex items-center justify-center rounded-lg bg-brand-500 px-4 py-2 text-sm font-semibold text-white shadow-theme-sm hover:bg-brand-600"
              onClick={() => setSearchId(lookupId.trim())}
            >
              Search
            </button>
          </div>

          {userQuery.isLoading && (
            <p className="text-sm text-gray-500">Loading user...</p>
          )}
          {userQuery.error && (
            <p className="text-sm text-error-500">
              {(userQuery.error as Error).message}
            </p>
          )}
          {userQuery.data && (
            <div className="overflow-hidden rounded-xl border border-gray-200 bg-white">
              <div className="max-w-full overflow-x-auto">
                <Table>
                  <TableHeader className="border-b border-gray-100">
                    <TableRow>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        User
                      </TableCell>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        Roles
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
                    <TableRow key={userQuery.data.id}>
                      <TableCell className="px-5 py-4 text-start">
                        <div>
                          <p className="text-sm font-semibold text-gray-900">
                            {userQuery.data.name}
                          </p>
                          <p className="text-xs text-gray-500">
                            {userQuery.data.email}
                          </p>
                        </div>
                      </TableCell>
                      <TableCell className="px-5 py-4 text-sm text-gray-600">
                        {userQuery.data.roles?.length ? (
                          userQuery.data.roles.join(", ")
                        ) : (
                          <Badge size="sm" color="light">
                            none
                          </Badge>
                        )}
                      </TableCell>
                      <TableCell className="px-5 py-4 text-sm text-gray-600">
                        {userQuery.data.permissions?.length ? (
                          userQuery.data.permissions.join(", ")
                        ) : (
                          <Badge size="sm" color="light">
                            none
                          </Badge>
                        )}
                      </TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </div>
            </div>
          )}
        </ComponentCard>

        <ComponentCard
          title="Create User"
          desc="Tambahkan user baru ke sistem."
        >
          <form
            className="grid gap-4 md:grid-cols-2"
            onSubmit={createForm.handleSubmit(onCreate)}
          >
            <div>
              <Input
                placeholder="Nama lengkap"
                error={Boolean(createForm.formState.errors.name)}
                {...createForm.register("name")}
              />
              {createForm.formState.errors.name && (
                <p className="mt-1 text-xs text-error-500">
                  {createForm.formState.errors.name.message}
                </p>
              )}
            </div>
            <div>
              <Input
                placeholder="Email"
                type="email"
                error={Boolean(createForm.formState.errors.email)}
                {...createForm.register("email")}
              />
              {createForm.formState.errors.email && (
                <p className="mt-1 text-xs text-error-500">
                  {createForm.formState.errors.email.message}
                </p>
              )}
            </div>
            <div>
              <Input
                placeholder="Password"
                type="password"
                error={Boolean(createForm.formState.errors.password)}
                {...createForm.register("password")}
              />
              {createForm.formState.errors.password && (
                <p className="mt-1 text-xs text-error-500">
                  {createForm.formState.errors.password.message}
                </p>
              )}
            </div>
            <div>
              <Input
                placeholder="Roles (comma separated)"
                error={Boolean(createForm.formState.errors.roles)}
                {...createForm.register("roles")}
              />
              {createForm.formState.errors.roles && (
                <p className="mt-1 text-xs text-error-500">
                  {createForm.formState.errors.roles.message}
                </p>
              )}
            </div>
            <div className="md:col-span-2 flex items-center gap-3">
              <button
                className="inline-flex items-center justify-center rounded-lg bg-brand-500 px-4 py-2 text-sm font-semibold text-white shadow-theme-sm hover:bg-brand-600"
                type="submit"
                disabled={createMutation.isPending}
              >
                {createMutation.isPending ? "Creating..." : "Create User"}
              </button>
              {createMutation.isError && (
                <p className="text-sm text-error-500">
                  {(createMutation.error as Error).message}
                </p>
              )}
              {createMutation.isSuccess && (
                <p className="text-sm text-success-500">User created.</p>
              )}
            </div>
          </form>
        </ComponentCard>
      </div>
    </>
  );
}
