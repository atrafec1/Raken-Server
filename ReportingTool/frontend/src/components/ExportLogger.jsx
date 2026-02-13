function ExportLogger({ progressLogs = [], isComplete = false, onStartNew }) {
  return (
    <div className={`
      border-2 p-4 w-full max-w-2xl h-96 flex flex-col
      ${isComplete ? 'border-green-500 bg-green-50' : 'border-gray-300 bg-white'}
    `}>
      {isComplete && (
        <div className="flex justify-between items-center mb-4 pb-3 border-b border-green-500 flex-shrink-0">
          <div className="text-green-500 font-bold text-lg">
            ✓ Export Complete!
          </div>
          <button 
            onClick={onStartNew}
            className="px-4 py-2 bg-green-500 text-white rounded-md cursor-pointer font-medium hover:bg-green-600"
          >
            Start New Export
          </button>
        </div>
      )}
      <div className="overflow-y-auto flex-1">
        {progressLogs.map((log, index) => (
          <div key={index} className="mb-2">
            {log}
          </div>
        ))}
      </div>
    </div>
  );
}

export default ExportLogger;