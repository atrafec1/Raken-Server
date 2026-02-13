import { useState } from "react";

function PayrollHelper() {
  const [weekEndingDate, setWeekEndingDate] = useState("");

  const setToLastSunday = () => {
    const today = new Date();
    const dayOfWeek = today.getDay();
    const lastSunday = new Date(today);
    lastSunday.setDate(today.getDate() - dayOfWeek);
    setWeekEndingDate(lastSunday.toISOString().split('T')[0]);
  };

  const handleImportFromRaken = () => {
    // TODO: Implement Raken import
    console.log("Import from Raken");
  };

  const handleUploadCSV = () => {
    // TODO: Implement CSV upload
    console.log("Upload CSV");
  };

  const handleGeneratePayroll = () => {
    // TODO: Implement payroll generation
    console.log("Generate Payroll for week ending:", weekEndingDate);
  };

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      {/* Header */}
      <div className="max-w-4xl mx-auto mb-8">
        <h1 className="text-4xl font-bold text-gray-800 text-center">Payroll Helper</h1>
      </div>

      {/* Cards Section */}
      <div className="max-w-4xl mx-auto mb-8 grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Import from Raken Card */}
        <div 
          onClick={handleImportFromRaken}
          className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow cursor-pointer border-2 border-transparent hover:border-blue-500"
        >
          <div className="text-center">
            <div className="text-5xl mb-4">📊</div>
            <h2 className="text-xl font-semibold text-gray-800 mb-2">
              Import Data from Raken
            </h2>
            <p className="text-gray-600 text-sm">
              Fetch timecard data from Raken API
            </p>
          </div>
        </div>

        {/* Upload CSV Card */}
        <div 
          onClick={handleUploadCSV}
          className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow cursor-pointer border-2 border-transparent hover:border-blue-500"
        >
          <div className="text-center">
            <div className="text-5xl mb-4">📁</div>
            <h2 className="text-xl font-semibold text-gray-800 mb-2">
              Upload CSV Files
            </h2>
            <p className="text-gray-600 text-sm">
              Upload existing CSV timecard files
            </p>
          </div>
        </div>
      </div>

      {/* Date Entry Section */}
      <div className="max-w-4xl mx-auto mb-8">
        <div className="bg-white rounded-lg shadow-md p-6">
          <label className="block text-lg font-semibold text-gray-800 mb-3">
            Week Ending Date
          </label>
          
          <div className="flex gap-3 mb-3">
            <input
              type="date"
              value={weekEndingDate}
              onChange={(e) => setWeekEndingDate(e.target.value)}
              className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              onClick={setToLastSunday}
              className="px-6 py-2 bg-gray-200 text-gray-800 rounded-md hover:bg-gray-300 font-medium"
            >
              Last Sunday
            </button>
          </div>

          <p className="text-sm text-gray-600 italic">
            Enter Sunday for given payroll week.
          </p>
        </div>
      </div>

      {/* Generate Button */}
      <div className="max-w-4xl mx-auto">
        <button
          onClick={handleGeneratePayroll}
          disabled={!weekEndingDate}
          className="w-full py-4 bg-blue-600 text-white text-lg font-semibold rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
        >
          Generate Payroll Files
        </button>
      </div>
    </div>
  );
}

export default PayrollHelper;