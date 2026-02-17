import { useEffect, useState } from "react";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import ExportLogger from "../components/ReportsPage/ExportLogger";
import DateInput from "../components/DateInput";
import FolderSelector from "../components/FolderSelector";
import ExportReportsButton from "../components/ReportsPage/ExportReportsButton";
import Error from "../components/Error";
import "../index.css";
import {
  ExportDailyReports,
  ChangeExportDir,
  GetExportDir
} from "../../wailsjs/go/main/App";


function isValidDate(dateString) {
  if (!dateString) return false;
  const regex = /^\d{4}-\d{2}-\d{2}$/;
  return regex.test(dateString);
}

function isValidDateRange(fromDate, toDate) {
  if (!isValidDate(fromDate) || !isValidDate(toDate)) {
    alert("Invalid Date range");
    return false;
  }
  if (new Date(fromDate) > new Date(toDate)) {
    alert("From date cannot be after To date");
    return false;
  }
  return true;
}

function ReportsPage() {
  // Function to calculate last Sunday and Monday
  const getLastWeek = () => {
    const lastSunday = new Date();
    lastSunday.setDate(lastSunday.getDate() - lastSunday.getDay());
    const lastMonday = new Date(lastSunday);
    lastMonday.setDate(lastMonday.getDate() - 6);
    return {
      monday: lastMonday.toISOString().split("T")[0],
      sunday: lastSunday.toISOString().split("T")[0]
    };
  };

  const lastWeek = getLastWeek();

  const [log, setLog] = useState([]);
  const [fromDate, setFromDate] = useState(lastWeek.monday);
  const [toDate, setToDate] = useState(lastWeek.sunday);
  const [savePath, setSavePath] = useState("Desktop/Raken Reports");
  const [error, setError] = useState(false);
  const [exportStatus, setExportStatus] = useState('idle'); // 'idle' | 'exporting' | 'success' | 'error'



  // Load export directory on mount
  useEffect(() => {
    const loadExportDir = async () => {
      try {
        const exportDir = await GetExportDir();
        setSavePath(exportDir);
      } catch {
        setError(true);
      }
    };
    loadExportDir();
  }, []);

  // Subscribe to backend export events
  useEffect(() => {
    const unsubscribeStatus = EventsOn("exportProgress", (message) => {
      setLog((prev) => [...prev, message]);
    });

    const unsubscribeComplete = EventsOn("exportComplete", () => {
      setExportStatus('success');
      setLog((prev) => [...prev, "Export complete."]);
    });

    const unsubscribeError = EventsOn("exportError", (message) => {
      setExportStatus('error');
      setLog((prev) => [...prev, `Error: ${message}`]);
    });

    return () => {
      unsubscribeStatus();
      unsubscribeComplete();
      unsubscribeError();
    };
  }, []);

  const exportReportsToComputer = async () => {
    if (!isValidDateRange(fromDate, toDate)) return;

    setLog([]);
    setExportStatus('exporting');

    try {
      await ExportDailyReports(fromDate, toDate);
    } catch (err) {
      setExportStatus('error');
      console.log(err);
    }
  };

  const handleStartNew = () => {
    setExportStatus('idle');
    setLog([]);
    const lastWeek = getLastWeek();
    setFromDate(lastWeek.monday);
    setToDate(lastWeek.sunday);
  };

  const handleFolderSelection = async (initialPath) => {
    try {
      const path = await ChangeExportDir(initialPath);
      setSavePath(path);
    } catch {
      setError(true);
    }
  };

  if (error) {
    return <Error />;
  }

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto mb-8">
        <div className="mb-32">
        <h1 className="text-4xl font-bold text-gray-800 text-center mb-8">
          Daily Report Exporter
        </h1>
        <div className="text-sm text-gray-500 text-center ">
        <p> Export daily reports across all projects.</p>
        <p> Note: Date range should not exceed 30 days</p>
        </div>
        </div> 
        {/* Show controls only when idle */}
        {exportStatus === 'idle' && (
          <>
            <div className="flex gap-4 mb-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 flex-1">
                <DateInput
                  value={fromDate}
                  onChange={(e) => setFromDate(e.target.value)}
                  label="Start Date"
                  subLabel="Beginning of export period"
                />

                <DateInput
                  value={toDate}
                  onChange={(e) => setToDate(e.target.value)}
                  label="End Date"
                  subLabel="End of export period"
                />
              </div>
              
            
            </div>

            <div className="flex items-end mb-6">
              <div className="flex-1">
                <FolderSelector
                  path={savePath}
                  onSelect={handleFolderSelection}
                />
              </div>
              <ExportReportsButton
                onClick={exportReportsToComputer}
                disabled={false}
              />
            </div>
          </>
        )}

        {/* Show logger when exporting, success, or error */}
        {exportStatus !== 'idle' && (
          <ExportLogger 
            progressLogs={log} 
            exportStatus={exportStatus}
            onStartNew={handleStartNew}
          />
        )}
      </div>
    </div>
  );
}

export default ReportsPage;