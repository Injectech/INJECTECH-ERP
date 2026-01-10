import { useMemo, useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import PageBreadcrumb from "../../components/common/PageBreadCrumb";
import PageMeta from "../../components/common/PageMeta";
import ComponentCard from "../../components/common/ComponentCard";
import Input from "../../components/form/input/InputField";
import Select from "../../components/form/Select";
import Label from "../../components/form/Label";
import Badge from "../../components/ui/badge/Badge";
import { Modal } from "../../components/ui/modal";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "../../components/ui/table";
import {
  adjustInventory,
  fetchAllInventory,
  fetchLocations,
  fetchProducts,
} from "../../services/erp";
import { useAuthStore } from "../../stores/authStore";

const adjustSchema = z.object({
  delta: z.preprocess(
    (value) => Number(value),
    z.number().refine((val) => val !== 0, "Delta tidak boleh 0"),
  ),
});

type AdjustForm = z.infer<typeof adjustSchema>;

export default function Inventory() {
  const [searchTerm, setSearchTerm] = useState("");
  const [productFilters, setProductFilters] = useState<Record<string, string>>(
    {},
  );
  const [isAdjustOpen, setAdjustOpen] = useState(false);
  const [selectedStock, setSelectedStock] = useState<null | {
    id: string;
    location: string;
    productId: string;
    productLabel: string;
    quantity: number;
  }>(null);
  const user = useAuthStore((state) => state.user);
  const isAdmin = Boolean(user?.roles?.includes("admin"));

  const productsQuery = useQuery({
    queryKey: ["products"],
    queryFn: fetchProducts,
  });

  const inventoryQuery = useQuery({
    queryKey: ["inventory-all"],
    queryFn: fetchAllInventory,
  });

  const locationsQuery = useQuery({
    queryKey: ["locations"],
    queryFn: fetchLocations,
    enabled: isAdmin,
  });

  const productMap = useMemo(() => {
    const map = new Map<string, string>();
    (productsQuery.data || []).forEach((product) => {
      map.set(product.id, `${product.sku} - ${product.name}`);
    });
    return map;
  }, [productsQuery.data]);

  const normalizedLocation = (value: string | null | undefined) =>
    value && value.trim() ? value.trim() : "Default";

  const lowStockThreshold = 5;

  const activeInventory = useMemo(
    () => inventoryQuery.data || [],
    [inventoryQuery.data],
  );

  const totalItems = useMemo(
    () =>
      (inventoryQuery.data || []).reduce(
        (sum, item) => sum + item.quantity,
        0,
      ),
    [inventoryQuery.data],
  );

  const lowStockCount = useMemo(
    () =>
      (inventoryQuery.data || []).filter(
        (item) => item.quantity > 0 && item.quantity <= lowStockThreshold,
      ).length,
    [inventoryQuery.data, lowStockThreshold],
  );

  const outOfStockCount = useMemo(
    () => (inventoryQuery.data || []).filter((item) => item.quantity <= 0).length,
    [inventoryQuery.data],
  );

  const inventoryByLocation = useMemo(() => {
    const grouped = new Map<string, typeof inventoryQuery.data>();
    activeInventory.forEach((item) => {
      const location = normalizedLocation(item.location);
      if (!grouped.has(location)) {
        grouped.set(location, []);
      }
      grouped.get(location)?.push(item);
    });
    return grouped;
  }, [activeInventory]);

  const locationKeys = useMemo(() => {
    const fromInventory = Array.from(inventoryByLocation.keys());
    const merged = new Set<string>([...fromInventory]);
    if (merged.size === 0) {
      merged.add("Default");
    }
    const defaultLocation =
      locationsQuery.data?.find((loc) => loc.is_default)?.name || "Default";
    return Array.from(merged).sort((a, b) => {
      if (a === defaultLocation) return -1;
      if (b === defaultLocation) return 1;
      return a.localeCompare(b);
    });
  }, [locationsQuery.data, inventoryByLocation]);

  const adjustForm = useForm<AdjustForm>({
    resolver: zodResolver(adjustSchema),
  });

  const adjustMutation = useMutation({
    mutationFn: ({ id, delta }: { id: string; delta: number }) =>
      adjustInventory(id, delta),
    onSuccess: () => {
      setAdjustOpen(false);
      setSelectedStock(null);
      adjustForm.reset();
      inventoryQuery.refetch();
    },
  });

  const onAdjust = (values: AdjustForm) => {
    if (!selectedStock) return;
    adjustMutation.mutate({ id: selectedStock.id, delta: values.delta });
  };

  const openAdjust = (item: {
    id: string;
    location: string;
    productId: string;
    quantity: number;
  }) => {
    const label = productMap.get(item.productId) || item.productId;
    setSelectedStock({
      id: item.id,
      location: normalizedLocation(item.location),
      productId: item.productId,
      productLabel: label,
      quantity: item.quantity,
    });
    adjustForm.reset({ delta: 0 });
    setAdjustOpen(true);
  };

  return (
    <>
      <PageMeta
        title="ERP Inventory | Dashboard"
        description="Inventory management page for ERP backend."
      />
      <PageBreadcrumb pageTitle="Inventory" />

      <div className="space-y-6">
        <div className="grid gap-4 md:grid-cols-3">
          <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-theme-sm">
            <p className="text-sm text-gray-500">Total Items</p>
            <p className="mt-2 text-2xl font-semibold text-gray-900">
              {totalItems.toLocaleString("id-ID")}
            </p>
            <p className="mt-1 text-xs text-gray-500">Across all locations</p>
          </div>
          <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-theme-sm">
            <p className="text-sm text-gray-500">Low Stock</p>
            <p className="mt-2 text-2xl font-semibold text-gray-900">
              {lowStockCount.toLocaleString("id-ID")}
            </p>
            <p className="mt-1 text-xs text-gray-500">Need reorder</p>
          </div>
          <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-theme-sm">
            <p className="text-sm text-gray-500">Out of Stock</p>
            <p className="mt-2 text-2xl font-semibold text-gray-900">
              {outOfStockCount.toLocaleString("id-ID")}
            </p>
            <p className="mt-1 text-xs text-gray-500">Critical items</p>
          </div>
        </div>

        <ComponentCard title="Stock Levels" desc="Monitor quantities by location.">
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div className="w-full md:max-w-sm">
              <Input
                placeholder="Search SKU or product"
                value={searchTerm}
                onChange={(event) => setSearchTerm(event.target.value)}
              />
            </div>
          </div>

          {inventoryQuery.isLoading && (
            <p className="text-sm text-gray-500">Loading inventory...</p>
          )}
          {inventoryQuery.error && (
            <p className="text-sm text-error-500">
              {(inventoryQuery.error as Error).message}
            </p>
          )}
          {!inventoryQuery.isLoading && activeInventory.length === 0 && (
            <p className="text-sm text-gray-500">
              Tidak ada stok untuk ditampilkan.
            </p>
          )}
          {!inventoryQuery.isLoading &&
            locationKeys.map((location) => {
              const stocks = inventoryByLocation.get(location) || [];
              const productOptions = [
                { value: "all", label: "Semua produk" },
                ...Array.from(
                  new Map(
                    stocks.map((item) => [
                      item.product_id,
                      productMap.get(item.product_id) || item.product_id,
                    ]),
                  ),
                ).map(([value, label]) => ({
                  value,
                  label,
                })),
              ];
              const selectedFilter = productFilters[location];
              const filteredByProduct =
                !selectedFilter || selectedFilter === "all"
                  ? stocks
                  : stocks.filter((item) => item.product_id === selectedFilter);
              const filtered = searchTerm
                ? filteredByProduct.filter((item) => {
                    const label =
                      productMap.get(item.product_id) || item.product_id;
                    return label.toLowerCase().includes(searchTerm.toLowerCase());
                  })
                : filteredByProduct;

              return (
                <div key={location} className="mt-6 rounded-2xl border border-gray-200 bg-white p-5 shadow-theme-sm">
                  <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                    <div>
                      <h3 className="text-lg font-semibold text-gray-900">
                        {location}
                      </h3>
                      <p className="text-xs text-gray-500">
                        {filtered.length} item tersedia
                      </p>
                    </div>
                    <div className="w-full md:max-w-xs">
                      <Select
                        options={productOptions}
                        placeholder="Filter product"
                        defaultValue="all"
                        onChange={(value) =>
                          setProductFilters((prev) => ({
                            ...prev,
                            [location]: value,
                          }))
                        }
                      />
                    </div>
                  </div>

                  {filtered.length === 0 ? (
                    <p className="mt-4 text-sm text-gray-500">
                      Tidak ada stok untuk lokasi ini.
                    </p>
                  ) : (
                    <div className="mt-4 overflow-hidden rounded-xl border border-gray-200 bg-white">
                      <div className="max-w-full overflow-x-auto">
                        <Table>
                          <TableHeader className="border-b border-gray-100">
                            <TableRow>
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
                                Quantity
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
                                  Action
                                </TableCell>
                              )}
                            </TableRow>
                          </TableHeader>
                          <TableBody className="divide-y divide-gray-100">
                            {filtered.map((stock) => (
                              <TableRow key={stock.id}>
                                <TableCell className="px-5 py-4 text-sm text-gray-700">
                                  {productMap.get(stock.product_id) || stock.product_id}
                                </TableCell>
                                <TableCell className="px-5 py-4 text-sm text-gray-600">
                                  {stock.quantity}
                                </TableCell>
                                <TableCell className="px-5 py-4 text-sm">
                                  <Badge
                                    size="sm"
                                    color={stock.quantity === 0 ? "error" : "success"}
                                  >
                                    {stock.quantity === 0 ? "out" : "healthy"}
                                  </Badge>
                                </TableCell>
                                {isAdmin && (
                                  <TableCell className="px-5 py-4 text-sm">
                                    <button
                                      className="text-brand-600 hover:text-brand-700"
                                      onClick={() =>
                                        openAdjust({
                                          id: stock.id,
                                          location,
                                          productId: stock.product_id,
                                          quantity: stock.quantity,
                                        })
                                      }
                                    >
                                      Adjust Stock
                                    </button>
                                  </TableCell>
                                )}
                              </TableRow>
                            ))}
                          </TableBody>
                        </Table>
                      </div>
                    </div>
                  )}
                </div>
              );
            })}
        </ComponentCard>
      </div>

      <Modal
        isOpen={isAdjustOpen && isAdmin}
        onClose={() => setAdjustOpen(false)}
        className="mx-4 w-[420px] max-w-[92vw]"
      >
        <div className="border-b border-gray-200 px-5 py-4 dark:border-gray-800">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white/90">
            Adjust Stock
          </h3>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {selectedStock?.productLabel} • {selectedStock?.location}
          </p>
        </div>
        <form className="grid gap-4 p-5" onSubmit={adjustForm.handleSubmit(onAdjust)}>
          <div>
            <Label>Stok Saat Ini</Label>
            <Input
              value={selectedStock?.quantity ?? 0}
              disabled
              hint="Nilai saat ini"
            />
          </div>
          <div>
            <Label>
              Delta <span className="text-error-500">*</span>
            </Label>
            <Input
              type="number"
              placeholder="Contoh: 10 atau -5"
              error={Boolean(adjustForm.formState.errors.delta)}
              {...adjustForm.register("delta", { valueAsNumber: true })}
            />
            {adjustForm.formState.errors.delta && (
              <p className="mt-1 text-xs text-error-500">
                {adjustForm.formState.errors.delta.message}
              </p>
            )}
          </div>
          <div className="flex items-center justify-end gap-3 border-t border-gray-200 pt-3 dark:border-gray-800">
            <button
              type="button"
              className="inline-flex items-center justify-center rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-semibold text-gray-700 shadow-theme-sm hover:bg-gray-50 dark:border-gray-800 dark:bg-gray-900 dark:text-gray-200"
              onClick={() => setAdjustOpen(false)}
            >
              Batal
            </button>
            <button
              className="inline-flex items-center justify-center rounded-lg bg-brand-500 px-4 py-2 text-sm font-semibold text-white shadow-theme-sm hover:bg-brand-600"
              type="submit"
              disabled={adjustMutation.isPending}
            >
              {adjustMutation.isPending ? "Menyimpan..." : "Simpan"}
            </button>
          </div>
          {adjustMutation.isError && (
            <p className="text-sm text-error-500">
              {(adjustMutation.error as Error).message}
            </p>
          )}
        </form>
      </Modal>
    </>
  );
}
