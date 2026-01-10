import ComponentCard from "../common/ComponentCard";
import Badge from "../ui/badge/Badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "../ui/table";

const alerts = [
  {
    sku: "PRD-1002",
    product: "Thermal Receipt Paper",
    location: "Bandung Warehouse",
    onHand: 6,
    min: 20,
    status: "low",
  },
  {
    sku: "PRD-1004",
    product: "Barcode Scanner",
    location: "Surabaya Warehouse",
    onHand: 0,
    min: 5,
    status: "out",
  },
  {
    sku: "PRD-1011",
    product: "Shipping Label Roll",
    location: "Jakarta Warehouse",
    onHand: 12,
    min: 30,
    status: "low",
  },
];

export default function InventoryAlerts() {
  return (
    <ComponentCard
      title="Inventory Alerts"
      desc="Items below safety stock or out of stock."
    >
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
                  Location
                </TableCell>
                <TableCell
                  isHeader
                  className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                >
                  On Hand
                </TableCell>
                <TableCell
                  isHeader
                  className="px-5 py-3 text-start text-theme-xs font-medium text-gray-500"
                >
                  Status
                </TableCell>
              </TableRow>
            </TableHeader>
            <TableBody className="divide-y divide-gray-100">
              {alerts.map((item) => (
                <TableRow key={item.sku}>
                  <TableCell className="px-5 py-4 text-sm font-semibold text-gray-900">
                    {item.sku}
                  </TableCell>
                  <TableCell className="px-5 py-4 text-sm text-gray-700">
                    {item.product}
                  </TableCell>
                  <TableCell className="px-5 py-4 text-sm text-gray-600">
                    {item.location}
                  </TableCell>
                  <TableCell className="px-5 py-4 text-sm text-gray-600">
                    {item.onHand} / min {item.min}
                  </TableCell>
                  <TableCell className="px-5 py-4 text-sm">
                    <Badge
                      size="sm"
                      color={item.status === "out" ? "error" : "warning"}
                    >
                      {item.status}
                    </Badge>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </div>
    </ComponentCard>
  );
}
