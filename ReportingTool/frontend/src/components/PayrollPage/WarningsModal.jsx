function WarningModal({ warnings, onCancel, onContinue, onExport}) {
  // Group warnings by type
  const groupedWarnings = warnings.reduce((acc, warning) => {
    if (!acc[warning.WarningType]) {
      acc[warning.WarningType] = [];
    }
    acc[warning.WarningType].push(warning);
    return acc;
  }, {});
  const handleExportClick = async () => {
    if(!onExport) return;
    await onExport(warnings);
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 max-w-2xl w-full max-h-[80vh] overflow-y-auto">
        <h2 className="text-2xl font-bold mb-4">⚠️ Warnings Found</h2>
        
        <p className="mb-4 text-gray-700">
          {warnings.length} issue{warnings.length !== 1 ? 's' : ''} detected:
        </p>

        <div className="space-y-4 mb-6">
          {Object.entries(groupedWarnings).map(([type, items]) => (
            <div key={type} className="border-l-4 border-yellow-500 pl-4">
              <h3 className="font-semibold text-gray-800">
                {type} ({items.length})
              </h3>
              <ul className="list-disc list-inside text-sm text-gray-600 space-y-1">
                {items.map((warning, idx) => (
                  <li key={idx}>{warning.Message}</li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        <p className="mb-6 text-gray-700 font-medium">
          Do you want to continue with export?
        </p>

        <div className="flex gap-4 justify-between items-center">
          <button
            onClick={handleExportClick}
            className="px-4 py-2 border border-blue-500 text-blue-600 rounded-md 
hover:bg-blue-600 hover:text-white transition-colors"
          >
            Export Warnings
          </button>
          
          <div className="flex gap-4">
            <button
              onClick={onCancel}
              className="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={onContinue}
              className="px-4 py-2 bg-yellow-500 text-white rounded-md hover:bg-yellow-600 transition-colors"
            >
              Continue Anyway
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default WarningModal;