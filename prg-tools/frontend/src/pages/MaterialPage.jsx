import { ExportProgressEstimate, GetJobMaterialInfo } from "../../wailsjs/go/main/App";
import Button from "../components/Button";
import { useState, useEffect } from "react";
import DateInput from "../components/DateInput";

function getLastWeekRange() {
    const today = new Date();
    const lastWeek = new Date(today);
    lastWeek.setDate(today.getDate() - 7);
    const toISO = (d) => d.toISOString().slice(0, 10);
    return { from: toISO(lastWeek), to: toISO(today) };
}

function MaterialPage() {
    const [fromDate, setFromDate] = useState("");
    const [toDate, setToDate] = useState("");

    useEffect(() => {
        const { from, to } = getLastWeekRange();
        setFromDate(from);
        setToDate(to);
    }, []);

    const handleProgressEstimateExport = async () => {
        try {
            const jobMaterialInfo = await GetJobMaterialInfo(fromDate, toDate);
            if (!jobMaterialInfo || jobMaterialInfo.length === 0) {
                alert("No job material information found.");
                return;
            }
            await ExportProgressEstimate(jobMaterialInfo);
            alert("Progress estimate exported successfully!");
        } catch (error) {
            console.error("Error exporting progress estimate:", error);
            alert(`Export failed: ${error.message || error}`);
        }
    };

    return (
        <div className="min-h-screen bg-gray-50 p-8 flex flex-col items-center">
            <div className="max-w-3xl w-full bg-white rounded-xl shadow-lg p-8 flex flex-col items-center">
                <h2 className="text-3xl font-bold text-gray-800 mb-2 text-center">
                    Progress Estimate Exporter
                </h2>
                <p className="text-gray-500 text-sm mb-6 text-center">
                    Select a date range to gather material logs and export to excel sheets
                </p>
                

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 w-full mb-6">
                    <DateInput label="From" value={fromDate} onChange={setFromDate} />
                    <DateInput label="To" value={toDate} onChange={setToDate} />
                </div>

                <Button 
                    text={"Export"}
                    onClick={handleProgressEstimateExport} 
                    className="w-full font-semibold text-base"
                >
                    Export Progress Estimate
                </Button>
            </div>
        </div>
    );
}

export default MaterialPage;