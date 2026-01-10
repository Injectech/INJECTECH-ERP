import ComponentCard from "../common/ComponentCard";

const modules = [
  {
    name: "Security & Access",
    coverage: 92,
    owner: "IT Security",
  },
  {
    name: "Inventory Operations",
    coverage: 78,
    owner: "Warehouse Ops",
  },
  {
    name: "Product Catalog",
    coverage: 64,
    owner: "Merchandising",
  },
  {
    name: "Audit Compliance",
    coverage: 88,
    owner: "Risk & Audit",
  },
];

export default function ModuleCoverage() {
  return (
    <ComponentCard title="Module Coverage" desc="Adoption status by team.">
      <div className="space-y-4">
        {modules.map((module) => (
          <div key={module.name}>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-semibold text-gray-900">
                  {module.name}
                </p>
                <p className="text-xs text-gray-500">{module.owner}</p>
              </div>
              <span className="text-sm font-semibold text-gray-700">
                {module.coverage}%
              </span>
            </div>
            <div className="mt-2 h-2 w-full rounded-full bg-gray-100">
              <div
                className="h-2 rounded-full bg-brand-500"
                style={{ width: `${module.coverage}%` }}
              />
            </div>
          </div>
        ))}
      </div>
    </ComponentCard>
  );
}
