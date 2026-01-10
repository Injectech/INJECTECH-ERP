import { Link } from "react-router";
import ComponentCard from "../common/ComponentCard";

const actions = [
  {
    title: "Create User",
    desc: "Invite a new team member",
    path: "/erp/users",
  },
  {
    title: "Add Product",
    desc: "Register a new SKU",
    path: "/erp/products",
  },
  {
    title: "Adjust Inventory",
    desc: "Update stock balance",
    path: "/erp/inventory",
  },
  {
    title: "Review Roles",
    desc: "Audit access control",
    path: "/erp/roles",
  },
];

export default function QuickActions() {
  return (
    <ComponentCard title="Quick Actions" desc="Common admin tasks.">
      <div className="grid gap-3 sm:grid-cols-2">
        {actions.map((action) => (
          <Link
            key={action.title}
            to={action.path}
            className="rounded-xl border border-gray-200 bg-white p-4 shadow-theme-xs transition hover:border-brand-200 hover:bg-brand-50"
          >
            <p className="text-sm font-semibold text-gray-900">
              {action.title}
            </p>
            <p className="mt-1 text-xs text-gray-500">{action.desc}</p>
          </Link>
        ))}
      </div>
    </ComponentCard>
  );
}
