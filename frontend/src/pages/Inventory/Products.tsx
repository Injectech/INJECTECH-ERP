import { useEffect, useMemo, useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { Controller, useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Link } from "react-router";
import PageBreadcrumb from "../../components/common/PageBreadCrumb";
import PageMeta from "../../components/common/PageMeta";
import ComponentCard from "../../components/common/ComponentCard";
import Input from "../../components/form/input/InputField";
import Label from "../../components/form/Label";
import TextArea from "../../components/form/input/TextArea";
import Badge from "../../components/ui/badge/Badge";
import { Modal } from "../../components/ui/modal";
import Select from "../../components/form/Select";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "../../components/ui/table";
import {
  createInventory,
  createProduct,
  deleteProduct,
  fetchLocations,
  fetchProducts,
  fetchInventory,
  updateProduct,
  updateInventoryLocation,
} from "../../services/erp";
import { useAuthStore } from "../../stores/authStore";

const categoryOptions = [
  { value: "all", label: "All categories" },
  { value: "hardware", label: "Hardware" },
  { value: "consumable", label: "Consumable" },
  { value: "service", label: "Service" },
];

const createSchema = z.object({
  sku: z.string().min(2, "SKU minimal 2 karakter"),
  name: z.string().min(2, "Nama produk minimal 2 karakter"),
  description: z.string().min(3, "Deskripsi minimal 3 karakter"),
  price: z.preprocess(
    (value) => Number(value),
    z.number().positive("Harga harus lebih dari 0"),
  ),
  stock: z.preprocess(
    (value) => Number(value),
    z.number().min(0, "Stok tidak boleh negatif"),
  ),
});

type CreateForm = z.infer<typeof createSchema>;

const editSchema = z.object({
  sku: z.string().min(2, "SKU minimal 2 karakter"),
  name: z.string().min(2, "Nama produk minimal 2 karakter"),
  description: z.string().min(3, "Deskripsi minimal 3 karakter"),
  price: z.preprocess(
    (value) => Number(value),
    z.number().positive("Harga harus lebih dari 0"),
  ),
});

type EditForm = z.infer<typeof editSchema>;

export default function Products() {
  const [isDialogOpen, setDialogOpen] = useState(false);
  const [isEditOpen, setEditOpen] = useState(false);
  const [isDeleteOpen, setDeleteOpen] = useState(false);
  const [selectedLocation, setSelectedLocation] = useState("default");
  const [editLocation, setEditLocation] = useState("default");
  const [initialEditLocation, setInitialEditLocation] = useState("default");
  const [editInventoryId, setEditInventoryId] = useState<string | null>(null);
  const [selectedProduct, setSelectedProduct] = useState<null | {
    id: string;
    sku: string;
    name: string;
  }>(null);
  const user = useAuthStore((state) => state.user);
  const isAdmin = Boolean(user?.roles?.includes("admin"));

  const productsQuery = useQuery({
    queryKey: ["products"],
    queryFn: fetchProducts,
  });

  const locationsQuery = useQuery({
    queryKey: ["locations"],
    queryFn: fetchLocations,
    enabled: isAdmin,
  });

  const locationOptions = useMemo(() => {
    const base = [{ value: "default", label: "Default (auto)" }];
    const locations = locationsQuery.data || [];
    return [
      ...base,
      ...locations.map((loc) => ({
        value: loc.name,
        label: loc.is_default ? `${loc.name} (default)` : loc.name,
      })),
    ];
  }, [locationsQuery.data]);

  const createForm = useForm<CreateForm>({
    resolver: zodResolver(createSchema),
    defaultValues: {
      description: "",
      stock: 0,
    },
  });

  const createMutation = useMutation({
    mutationFn: async (payload: CreateForm) => {
      const product = await createProduct({
        sku: payload.sku.trim(),
        name: payload.name.trim(),
        description: payload.description.trim(),
        price: payload.price,
      });
      await createInventory({
        product_id: product.id,
        quantity: payload.stock,
        location: selectedLocation === "default" ? "" : selectedLocation,
      });
      return product;
    },
    onSuccess: () => {
      createForm.reset();
      setSelectedLocation("default");
      setDialogOpen(false);
      productsQuery.refetch();
    },
  });

  const editForm = useForm<EditForm>({
    resolver: zodResolver(editSchema),
    defaultValues: {
      description: "",
    },
  });

  const updateMutation = useMutation({
    mutationFn: async (payload: {
      id: string;
      sku: string;
      name: string;
      description: string;
      price: number;
      inventoryId?: string | null;
      location?: string;
    }) => {
      const updated = await updateProduct(payload);
      const nextLocation = payload.location ?? "default";
      let targetInventoryId = payload.inventoryId ?? null;
      if (!targetInventoryId) {
        const items = await fetchInventory(payload.id);
        targetInventoryId = items[0]?.id ?? null;
      }
      if (targetInventoryId) {
        await updateInventoryLocation(
          targetInventoryId,
          nextLocation === "default" ? "" : nextLocation,
        );
      }
      return updated;
    },
    onSuccess: () => {
      setEditOpen(false);
      setSelectedProduct(null);
      setEditInventoryId(null);
      productsQuery.refetch();
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteProduct,
    onSuccess: () => {
      setDeleteOpen(false);
      setSelectedProduct(null);
      productsQuery.refetch();
    },
  });

  const inventoryQuery = useQuery({
    queryKey: ["inventory", selectedProduct?.id],
    queryFn: () => fetchInventory(selectedProduct!.id),
    enabled: Boolean(selectedProduct?.id && isEditOpen),
  });

  const inventoryTotal =
    inventoryQuery.data?.reduce((sum, item) => sum + item.quantity, 0) ?? 0;

  useEffect(() => {
    if (!isEditOpen) {
      setEditLocation("default");
      setInitialEditLocation("default");
      setEditInventoryId(null);
      return;
    }
    const first = inventoryQuery.data?.[0];
    if (first) {
      const next = first.location?.trim() ? first.location : "default";
      setEditLocation(next);
      setInitialEditLocation(next);
      setEditInventoryId(first.id);
    } else {
      setEditLocation("default");
      setInitialEditLocation("default");
      setEditInventoryId(null);
    }
  }, [inventoryQuery.data, isEditOpen]);

  const onCreate = (values: CreateForm) => {
    createMutation.mutate(values);
  };

  const onUpdate = (values: EditForm) => {
    if (!selectedProduct) return;
    updateMutation.mutate({
      id: selectedProduct.id,
      sku: values.sku.trim(),
      name: values.name.trim(),
      description: values.description.trim(),
      price: values.price,
      inventoryId: editInventoryId,
      location: editLocation,
    });
  };

  const openEdit = (product: {
    id: string;
    sku: string;
    name: string;
    description: string;
    price: number;
  }) => {
    setSelectedProduct({ id: product.id, sku: product.sku, name: product.name });
    editForm.reset({
      sku: product.sku,
      name: product.name,
      description: product.description,
      price: product.price,
    });
    setEditOpen(true);
  };

  const openDelete = (product: { id: string; sku: string; name: string }) => {
    setSelectedProduct(product);
    setDeleteOpen(true);
  };

  return (
    <>
      <PageMeta
        title="ERP Products | Dashboard"
        description="Product management page for ERP backend."
      />
      <PageBreadcrumb pageTitle="Products" />

      <div className="space-y-6">
        <ComponentCard title="Product Catalog" desc="Maintain SKUs and pricing.">
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div className="grid w-full gap-3 md:grid-cols-2">
              <Input placeholder="Search SKU or name" />
              <Select
                options={categoryOptions}
                placeholder="Category"
                onChange={() => {}}
              />
            </div>
            {isAdmin && (
              <button
                className="inline-flex items-center justify-center rounded-lg bg-brand-500 px-4 py-2 text-sm font-semibold text-white shadow-theme-sm hover:bg-brand-600"
                onClick={() => setDialogOpen(true)}
              >
                + Product
              </button>
            )}
          </div>

          <Modal
            isOpen={isDialogOpen && isAdmin}
            onClose={() => setDialogOpen(false)}
            className="mx-4 w-[720px] max-w-[92vw]"
          >
            <div className="border-b border-gray-200 px-5 py-4 dark:border-gray-800 sm:px-6">
              <h3 className="text-xl font-semibold text-gray-900 dark:text-white/90">
                Tambah Produk
              </h3>
              <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                Lengkapi informasi SKU, harga, stok awal, dan deskripsi produk.
              </p>
            </div>
            <form
              className="grid gap-4 p-5 md:grid-cols-2 sm:p-6"
              onSubmit={createForm.handleSubmit(onCreate)}
            >
                <div>
                  <Label>
                    SKU <span className="text-error-500">*</span>
                  </Label>
                  <Input
                    placeholder="SKU-001"
                    error={Boolean(createForm.formState.errors.sku)}
                    {...createForm.register("sku")}
                  />
                  {createForm.formState.errors.sku && (
                    <p className="mt-1 text-xs text-error-500">
                      {createForm.formState.errors.sku.message}
                    </p>
                  )}
                </div>
                <div>
                  <Label>
                    Nama Produk <span className="text-error-500">*</span>
                  </Label>
                  <Input
                    placeholder="Nama produk"
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
                        placeholder="Deskripsi singkat produk"
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
                <div>
                  <Label>
                    Harga <span className="text-error-500">*</span>
                  </Label>
                  <Input
                    type="number"
                    placeholder="0"
                    error={Boolean(createForm.formState.errors.price)}
                    {...createForm.register("price", { valueAsNumber: true })}
                  />
                  {createForm.formState.errors.price && (
                    <p className="mt-1 text-xs text-error-500">
                      {createForm.formState.errors.price.message}
                    </p>
                  )}
                </div>
                <div>
                  <Label>
                    Stok Awal <span className="text-error-500">*</span>
                  </Label>
                  <Input
                    type="number"
                    placeholder="0"
                    error={Boolean(createForm.formState.errors.stock)}
                    {...createForm.register("stock", { valueAsNumber: true })}
                  />
                  {createForm.formState.errors.stock && (
                    <p className="mt-1 text-xs text-error-500">
                      {createForm.formState.errors.stock.message}
                    </p>
                  )}
                  <p className="mt-1 text-xs text-gray-500">
                    Stok awal akan masuk ke inventory default.
                  </p>
                </div>
                <div>
                  <Label>Location</Label>
                  <Select
                    key={selectedLocation}
                    options={locationOptions}
                    placeholder="Pilih lokasi"
                    defaultValue={selectedLocation}
                    onChange={(value) => setSelectedLocation(value)}
                  />
                  <p className="mt-1 text-xs text-gray-500">
                    Jika tidak dipilih, akan menggunakan lokasi default.
                  </p>
                  {locationsQuery.isError && (
                    <p className="mt-1 text-xs text-error-500">
                      {(locationsQuery.error as Error).message}
                    </p>
                  )}
                </div>
                <div className="md:col-span-2 flex items-center justify-end gap-3 border-t border-gray-200 pt-3 dark:border-gray-800">
                  <button
                    type="button"
                    className="inline-flex items-center justify-center rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-semibold text-gray-700 shadow-theme-sm hover:bg-gray-50 dark:border-gray-800 dark:bg-gray-900 dark:text-gray-200"
                    onClick={() => setDialogOpen(false)}
                  >
                    Batal
                  </button>
                  <button
                    className="inline-flex items-center justify-center rounded-lg bg-brand-500 px-4 py-2 text-sm font-semibold text-white shadow-theme-sm hover:bg-brand-600"
                    type="submit"
                    disabled={createMutation.isPending}
                  >
                    {createMutation.isPending ? "Menyimpan..." : "Simpan"}
                  </button>
                </div>
                {createMutation.isError && (
                  <p className="text-sm text-error-500 md:col-span-2">
                    {(createMutation.error as Error).message}
                  </p>
                )}
            </form>
          </Modal>

          {productsQuery.isLoading && (
            <p className="text-sm text-gray-500">Loading products...</p>
          )}
          {productsQuery.error && (
            <p className="text-sm text-error-500">
              {(productsQuery.error as Error).message}
            </p>
          )}
          {productsQuery.data && (
            <div className="overflow-hidden rounded-xl border border-gray-200 bg-white">
              <div className="max-w-full overflow-x-auto">
                <Table>
                  <TableHeader className="border-b border-gray-100">
                    <TableRow>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        SKU
                      </TableCell>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        Product
                      </TableCell>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        Price
                      </TableCell>
                      <TableCell
                        isHeader
                        className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                      >
                        Status
                      </TableCell>
                      {isAdmin && (
                        <TableCell
                          isHeader
                          className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                        >
                          Actions
                        </TableCell>
                      )}
                    </TableRow>
                  </TableHeader>
                  <TableBody className="divide-y divide-gray-100">
                    {productsQuery.data.map((product) => (
                      <TableRow key={product.id}>
                        <TableCell className="px-5 py-4 text-sm font-semibold text-gray-900">
                          {product.sku}
                        </TableCell>
                        <TableCell className="px-5 py-4 text-sm text-gray-700">
                          {product.name}
                        </TableCell>
                        <TableCell className="px-5 py-4 text-sm text-gray-600">
                          Rp {product.price.toLocaleString("id-ID")}
                        </TableCell>
                        <TableCell className="px-5 py-4 text-sm">
                          <Badge size="sm" color="success">
                            active
                          </Badge>
                        </TableCell>
                        {isAdmin && (
                          <TableCell className="px-5 py-4 text-sm">
                            <div className="flex items-center gap-3">
                              <button
                                className="text-brand-600 hover:text-brand-700"
                                onClick={() =>
                                  openEdit({
                                    id: product.id,
                                    sku: product.sku,
                                    name: product.name,
                                    description: product.description,
                                    price: product.price,
                                  })
                                }
                              >
                                Edit
                              </button>
                              <button
                                className="text-error-500 hover:text-error-600"
                                onClick={() =>
                                  openDelete({
                                    id: product.id,
                                    sku: product.sku,
                                    name: product.name,
                                  })
                                }
                              >
                                Delete
                              </button>
                            </div>
                          </TableCell>
                        )}
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            </div>
          )}
        </ComponentCard>
      </div>

      <Modal
        isOpen={isEditOpen && isAdmin}
        onClose={() => setEditOpen(false)}
        className="mx-4 w-[720px] max-w-[92vw]"
      >
        <div className="border-b border-gray-200 px-5 py-4 dark:border-gray-800 sm:px-6">
          <h3 className="text-xl font-semibold text-gray-900 dark:text-white/90">
            Edit Produk
          </h3>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Perbarui informasi SKU, harga, stok, dan deskripsi produk.
          </p>
        </div>
        <form
          className="grid gap-4 p-5 md:grid-cols-2 sm:p-6"
          onSubmit={editForm.handleSubmit(onUpdate)}
        >
          <div>
            <Label>
              SKU <span className="text-error-500">*</span>
            </Label>
            <Input
              placeholder="SKU-001"
              error={Boolean(editForm.formState.errors.sku)}
              {...editForm.register("sku")}
            />
            {editForm.formState.errors.sku && (
              <p className="mt-1 text-xs text-error-500">
                {editForm.formState.errors.sku.message}
              </p>
            )}
          </div>
          <div>
            <Label>
              Nama Produk <span className="text-error-500">*</span>
            </Label>
            <Input
              placeholder="Nama produk"
              error={Boolean(editForm.formState.errors.name)}
              {...editForm.register("name")}
            />
            {editForm.formState.errors.name && (
              <p className="mt-1 text-xs text-error-500">
                {editForm.formState.errors.name.message}
              </p>
            )}
          </div>
          <div className="md:col-span-2">
            <Label>
              Deskripsi <span className="text-error-500">*</span>
            </Label>
            <Controller
              control={editForm.control}
              name="description"
              render={({ field }) => (
                <TextArea
                  placeholder="Deskripsi singkat produk"
                  rows={3}
                  error={Boolean(editForm.formState.errors.description)}
                  onChange={field.onChange}
                  value={field.value || ""}
                />
              )}
            />
            {editForm.formState.errors.description && (
              <p className="mt-1 text-xs text-error-500">
                {editForm.formState.errors.description.message}
              </p>
            )}
          </div>
          <div>
            <Label>
              Harga <span className="text-error-500">*</span>
            </Label>
            <Input
              type="number"
              placeholder="0"
              error={Boolean(editForm.formState.errors.price)}
              {...editForm.register("price", { valueAsNumber: true })}
            />
            {editForm.formState.errors.price && (
              <p className="mt-1 text-xs text-error-500">
                {editForm.formState.errors.price.message}
              </p>
            )}
          </div>
          <div>
            <Label>Location</Label>
            <Select
              key={editLocation}
              options={locationOptions}
              placeholder="Pilih lokasi"
              defaultValue={editLocation}
              onChange={(value) => setEditLocation(value)}
            />
            <p className="mt-1 text-xs text-gray-500">
              Lokasi stok utama. Untuk multi-lokasi gunakan modul Inventory.
            </p>
          </div>
          <div>
            <Label>Stok Saat Ini</Label>
            <Input
              type="number"
              value={inventoryTotal}
              disabled
              hint={
                inventoryQuery.isLoading
                  ? "Memuat stok..."
                  : "Stok diatur lewat inventory"
              }
            />
            <div className="mt-2 flex items-center justify-between text-xs text-gray-500">
              <span>Perubahan stok dilakukan di modul Inventory.</span>
              <Link
                to="/inventory/inventory"
                className="font-semibold text-brand-600 hover:text-brand-700"
              >
                Kelola Inventory
              </Link>
            </div>
          </div>
          <div className="md:col-span-2 flex items-center justify-end gap-3 border-t border-gray-200 pt-3 dark:border-gray-800">
            <button
              type="button"
              className="inline-flex items-center justify-center rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-semibold text-gray-700 shadow-theme-sm hover:bg-gray-50 dark:border-gray-800 dark:bg-gray-900 dark:text-gray-200"
              onClick={() => setEditOpen(false)}
            >
              Batal
            </button>
            <button
              className="inline-flex items-center justify-center rounded-lg bg-brand-500 px-4 py-2 text-sm font-semibold text-white shadow-theme-sm hover:bg-brand-600"
              type="submit"
              disabled={updateMutation.isPending}
            >
              {updateMutation.isPending ? "Menyimpan..." : "Ubah"}
            </button>
          </div>
          {updateMutation.isError && (
            <p className="text-sm text-error-500 md:col-span-2">
              {(updateMutation.error as Error).message}
            </p>
          )}
        </form>
      </Modal>

      <Modal
        isOpen={isDeleteOpen && isAdmin}
        onClose={() => setDeleteOpen(false)}
        className="mx-4 w-fit max-w-[92vw]"
      >
        <div className="px-5 py-6 sm:px-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white/90">
            Hapus
          </h3>
          <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
            Yakin ingin menghapus{" "}
            <span className="font-semibold text-gray-900 dark:text-white/90">
              {selectedProduct?.name || selectedProduct?.sku}
            </span>
            ?
          </p>
          <div className="mt-5 flex items-center justify-end gap-3">
            <button
              type="button"
              className="inline-flex items-center justify-center rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-semibold text-gray-700 shadow-theme-sm hover:bg-gray-50 dark:border-gray-800 dark:bg-gray-900 dark:text-gray-200"
              onClick={() => setDeleteOpen(false)}
            >
              Batal
            </button>
            <button
              className="inline-flex items-center justify-center rounded-lg bg-error-500 px-4 py-2 text-sm font-semibold text-white shadow-theme-sm hover:bg-error-600"
              onClick={() => selectedProduct && deleteMutation.mutate(selectedProduct.id)}
              disabled={deleteMutation.isPending}
            >
              {deleteMutation.isPending ? "Menghapus..." : "Hapus"}
            </button>
          </div>
          {deleteMutation.isError && (
            <p className="mt-3 text-sm text-error-500">
              {(deleteMutation.error as Error).message}
            </p>
          )}
        </div>
      </Modal>
    </>
  );
}
