import { useState } from "react";
import Loading from "../components/Loading";
import DateInput from "../components/DateInput";
import Button from "../components/Button";
import { ProcessPayroll, FetchPayrollEntries, ExportPayrollWarnings } from "../../wailsjs/go/main/App";
import WarningModal from "../components/PayrollPage/WarningsModal";

function PayrollHelper() {
  // Calculate default dates (last Monday to last Sunday)
  const lastSunday = new Date();
  lastSunday.setDate(lastSunday.getDate() - lastSunday.getDay());
  const lastMonday = new Date(lastSunday);
  lastMonday.setDate(lastMonday.getDate() - 6);

  const [loading, setLoading] = useState(false);
  const [fromDate, setFromDate] = useState(lastMonday.toISOString().split("T")[0]);
  const [toDate, setToDate] = useState(lastSunday.toISOString().split("T")[0]);
  const [warnings, setWarnings] = useState([]);
  const [warningModalOpen, setWarningModalOpen] = useState(false);
  const [payrollResult, setPayrollResult] = useState(null);

  const processPayroll = async () => {
    setLoading(true);
    try {
      const result = await FetchPayrollEntries(fromDate, toDate);

      if (!result.Entries || result.Entries.length === 0) {
        alert("No payroll entries found for the selected date range.");
        return;
      }

      setPayrollResult(result);

      if (result.Warnings && result.Warnings.length > 0) {
        setWarnings(result.Warnings);
        setWarningModalOpen(true);
      } else {
        await exportPayroll(result);
      }
    } catch (error) {
      console.error("Error processing payroll:", error);
      alert(`Error: ${error.message || error}`);
    } finally {
      setLoading(false);
    }
  };

  const exportPayroll = async (result) => {
    setLoading(true);
    try {
      await ProcessPayroll(result);
      alert("Payroll processed successfully!");

      setPayrollResult(null);
      setWarnings([]);
    } catch (error) {
      console.error("Error exporting payroll:", error);
      alert(`Export failed: ${error.message || error}`);
    } finally {
      setLoading(false);
    }
  };

  const handleContinueWithWarnings = async () => {
    setWarningModalOpen(false);
    await exportPayroll(payrollResult);
  };

  const handleCancelExport = () => {
    setWarningModalOpen(false);
    setWarnings([]);
    setPayrollResult(null);
  };

  const handleExportWarnings = async (warnings) => {
  try {
    setLoading(true);
    await ExportPayrollWarnings(warnings);
    alert("Warnings exported");
  } catch (err) {
    console.error("Export warnings failed:", err);
    alert(`Export failed: ${err.message || err}`);
  } finally {
    setLoading(false);
  }
};

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto mb-8">
        <h1 className="text-4xl font-bold text-gray-800 text-center mb-8">
          Payroll Helper
        </h1>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <DateInput
            value={fromDate}
            onChange={(e) => setFromDate(e.target.value)}
            label="Start Date"
            subLabel="Beginning of payroll period"
          />

          <DateInput
            value={toDate}
            onChange={(e) => setToDate(e.target.value)}
            label="Week Ending Date"
            subLabel="End of payroll period"
          />
        </div>

        <div className="flex justify-center">
          <Button
            text={loading ? "Processing..." : "Process Payroll"}
            onClick={processPayroll}
            disabled={loading}
          />
        </div>

        {loading && <Loading />}
      </div>

      {warningModalOpen && (
        <WarningModal
          warnings={warnings}
          onCancel={handleCancelExport}
          onContinue={handleContinueWithWarnings}
          onExport={handleExportWarnings}
        />
      )}
    </div>
  );
}

export default PayrollHelper;
