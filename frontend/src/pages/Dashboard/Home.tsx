import PageBreadcrumb from "../../components/common/PageBreadCrumb";
import PageMeta from "../../components/common/PageMeta";
import ErpMetrics from "../../components/erp/ErpMetrics";
import InventoryAlerts from "../../components/erp/InventoryAlerts";
import RecentAudit from "../../components/erp/RecentAudit";
import QuickActions from "../../components/erp/QuickActions";
import ModuleCoverage from "../../components/erp/ModuleCoverage";

export default function Home() {
  return (
    <>
      <PageMeta
        title="ERP Dashboard | Overview"
        description="Operational overview for ERP modules."
      />
      <PageBreadcrumb pageTitle="Dashboard" />
      <div className="grid grid-cols-12 gap-4 md:gap-6">
        <div className="col-span-12">
          <ErpMetrics />
        </div>
        <div className="col-span-12 space-y-6 xl:col-span-7">
          <InventoryAlerts />
          <ModuleCoverage />
        </div>
        <div className="col-span-12 space-y-6 xl:col-span-5">
          <RecentAudit />
          <QuickActions />
        </div>
      </div>
    </>
  );
}
