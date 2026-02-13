import {useEffect, useState} from "react";
import DateRangePicker from "../components/DateRangePicker";
import FolderSelector from "../components/FolderSelector";
import ExportReportsButton from "../components/ExportReportsButton";
import Error from "../components/Error";
import '../index.css';
import {ExportReports, ChangeExportDir, GetExportDir} from "../../wailsjs/go/main/App";


function isValidDate(dateString) {
    if (!dateString) return false;
    
    const regex = /^\d{4}-\d{2}-\d{2}$/;
    return regex.test(dateString);
}
function isValidDateRange(fromDate, toDate) {
    if (!isValidDate(fromDate) || !isValidDate(toDate)) {
        alert("Invalid Date range")
        return false;
    }
    if (new Date(fromDate) > new Date(toDate)) {
        alert("From date cannot be after To date")
        return false;
    }
    return true;
}

function ReportsPage() {

    useEffect(() => {
    const loadExportDir = async () => {
        try {
            const exportDir = await GetExportDir();
            setSavePath(exportDir);
        } catch (err) {
            setError(true);
        }
    };
    loadExportDir();
    },  []);


    const exportReportsToComputer = async () => {
        if (isValidDateRange(fromDate, toDate, savePath)) {
            try {
                const result = await ExportReports(fromDate, toDate);
                console.log("Reports exported successfully:", result);
            } catch (error) {
                console.log(error)
            }
        }
    };

    const [fromDate, setFromDate] = useState('');
    const [toDate, setToDate] = useState('');
    const [savePath, setSavePath] = useState('Desktop/Raken Reports');
    const handleDateSelection = (fromDate, toDate) => {
        setFromDate(fromDate)
        setToDate(toDate)
    }
    const [error, setError] = useState(false);

    const handleFolderSelection = async (initialPath) => {
        try {
           const path = await ChangeExportDir(initialPath)
           setSavePath(path)
        } catch (err) {
            setError(true)
        }
    }


 
    return (
    <div className="flex flex-col h-screen">
        <div className="pt-8 text-center">
        <h1 className="text-primary font-bold text-6xl">Raken Report Exporter</h1>
        </div>
        
        <div className="flex-1 flex flex-col gap-4 items-center justify-center">
        <div className="flex items-end gap-3">
            <DateRangePicker fromDate={fromDate} toDate={toDate} onChange={handleDateSelection} />
            <ExportReportsButton onClick={exportReportsToComputer} />
        </div>
        
        <FolderSelector path={savePath} onSelect={handleFolderSelection} />
        </div>
    </div>
);
}

export default ReportsPage