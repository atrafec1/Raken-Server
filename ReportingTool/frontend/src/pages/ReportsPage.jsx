import { useEffect, useState } from "react";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import ExportLogger from "../components/ReportsPage/ExportLogger";
import DateRangePicker from "../components/DateRangePicker";
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
  const [log, setLog] = useState([]);
  const [fromDate, setFromDate] = useState("");
  const [toDate, setToDate] = useState("");
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

    setLog([]);          // Clear previous run
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
    setFromDate("");
    setToDate("");
  };

  const handleDateSelection = (fromDate, toDate) => {
    setFromDate(fromDate);
    setToDate(toDate);
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
    <div className="flex flex-col h-screen">
      <div className="pt-8 text-center">
        <h1 className="text-primary font-bold text-6xl">
          Raken Report Exporter
        </h1>
      </div>

      <div className="flex-1 flex flex-col gap-4 items-center justify-center">
        {/* Show controls only when idle */}
        {exportStatus === 'idle' && (
          <>
            <div className="flex items-end gap-3">
              <DateRangePicker
                fromDate={fromDate}
                toDate={toDate}
                onChange={handleDateSelection}
              />
              <ExportReportsButton
                onClick={exportReportsToComputer}
                disabled={false}
              />
            </div>

            <FolderSelector
              path={savePath}
              onSelect={handleFolderSelection}
            />
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