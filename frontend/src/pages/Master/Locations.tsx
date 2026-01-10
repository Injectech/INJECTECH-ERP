import { useMutation, useQuery } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import PageBreadcrumb from "../../components/common/PageBreadCrumb";
import PageMeta from "../../components/common/PageMeta";
import ComponentCard from "../../components/common/ComponentCard";
import Input from "../../components/form/input/InputField";
import Label from "../../components/form/Label";
import TextArea from "../../components/form/input/TextArea";
import Badge from "../../components/ui/badge/Badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "../../components/ui/table";
import { createLocation, fetchLocations } from "../../services/erp";

const createSchema = z.object({
  name: z.string().min(2, "Nama lokasi minimal 2 karakter"),
  description: z.string().min(3, "Deskripsi minimal 3 karakter"),
});

type CreateForm = z.infer<typeof createSchema>;

export default function Locations() {
  const locationsQuery = useQuery({
    queryKey: ["locations"],
    queryFn: fetchLocations,
  });

  const createForm = useForm<CreateForm>({
    resolver: zodResolver(createSchema),
    defaultValues: {
      description: "",
    },
  });

  const createMutation = useMutation({
    mutationFn: createLocation,
    onSuccess: () => {
      createForm.reset();
      locationsQuery.refetch();
    },
  });

  const onCreate = (values: CreateForm) => {
    createMutation.mutate({
      name: values.name.trim(),
      description: values.description.trim(),
    });
  };

  return (
    <>
      <PageMeta
        title="ERP Locations | Dashboard"
        description="Location management page for ERP backend."
      />
      <PageBreadcrumb pageTitle="Locations" />

      <div className="space-y-6">
        <ComponentCard
          title="Add Location"
          desc="Tambahkan lokasi gudang atau kantor perusahaan."
        >
          <form
            className="grid gap-4 md:grid-cols-2"
            onSubmit={createForm.handleSubmit(onCreate)}
          >
            <div>
              <Label>
                Nama Lokasi <span className="text-error-500">*</span>
              </Label>
              <Input
                placeholder="Contoh: Warehouse Bandung"
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
              <TextArea
                placeholder="Keterangan singkat lokasi"
                rows={3}
                error={Boolean(createForm.formState.errors.description)}
                onChange={(value) =>
                  createForm.setValue("description", value, {
                    shouldValidate: true,
                  })
                }
                value={createForm.watch("description") || ""}
              />
              {createForm.formState.errors.description && (
                <p className="mt-1 text-xs text-error-500">
                  {createForm.formState.errors.description.message}
                </p>
              )}
            </div>
            <div className="md:col-span-2 flex items-center gap-3">
              <button
                className="inline-flex items-center justify-center rounded-lg bg-brand-500 px-4 py-2 text-sm font-semibold text-white shadow-theme-sm hover:bg-brand-600"
                type="submit"
                disabled={createMutation.isPending}
              >
                {createMutation.isPending ? "Menyimpan..." : "Simpan Lokasi"}
              </button>
              {createMutation.isError && (
                <p className="text-sm text-error-500">
                  {(createMutation.error as Error).message}
                </p>
              )}
              {createMutation.isSuccess && (
                <p className="text-sm text-success-500">
                  Lokasi berhasil dibuat.
                </p>
              )}
            </div>
          </form>
        </ComponentCard>

        <ComponentCard
          title="Location List"
          desc="Daftar lokasi yang tersedia untuk inventory."
        >
          {locationsQuery.isLoading && (
            <p className="text-sm text-gray-500">Loading locations...</p>
          )}
          {locationsQuery.error && (
            <p className="text-sm text-error-500">
              {(locationsQuery.error as Error).message}
            </p>
          )}
          {locationsQuery.data && (
            <div className="overflow-hidden rounded-xl border border-gray-200 bg-white">
              <div className="max-w-full overflow-x-auto">
                <Table>
                  <TableHeader className="border-b border-gray-100">
                    <TableRow>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        Location
                      </TableCell>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        Description
                      </TableCell>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        Default
                      </TableCell>
                    </TableRow>
                  </TableHeader>
                  <TableBody className="divide-y divide-gray-100">
                    {locationsQuery.data.map((location) => (
                      <TableRow key={location.id}>
                        <TableCell className="px-5 py-4 text-sm font-semibold text-gray-900">
                          {location.name}
                        </TableCell>
                        <TableCell className="px-5 py-4 text-sm text-gray-600">
                          {location.description}
                        </TableCell>
                        <TableCell className="px-5 py-4 text-sm">
                          {location.is_default ? (
                            <Badge size="sm" color="primary">
                              default
                            </Badge>
                          ) : (
                            <Badge size="sm" color="light">
                              no
                            </Badge>
                          )}
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
