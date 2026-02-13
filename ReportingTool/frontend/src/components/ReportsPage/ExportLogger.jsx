
function ExportLogger({ progressLogs = [], exportStatus = 'exporting', onStartNew }) {
  const isComplete = exportStatus === 'success' || exportStatus === 'error';
  const isSuccess = exportStatus === 'success';
  const isError = exportStatus === 'error';

  return (
    <div className={`
      border-2 p-4 w-full max-w-2xl h-96 flex flex-col
      ${isSuccess ? 'border-green-500 bg-green-50' : ''}
      ${isError ? 'border-red-500 bg-red-50' : ''}
      ${!isComplete ? 'border-gray-300 bg-white' : ''}
    `}>
      {isComplete && (
        <div className={`
          flex justify-between items-center mb-4 pb-3 border-b flex-shrink-0
          ${isSuccess ? 'border-green-500' : ''}
          ${isError ? 'border-red-500' : ''}
        `}>
          <div className={`
            font-bold text-lg
            ${isSuccess ? 'text-green-500' : ''}
            ${isError ? 'text-red-500' : ''}
          `}>
            {isSuccess ? '✓ Export Complete!' : '❌ Export Failed'}
          </div>
          <button 
            onClick={onStartNew}
            className={`
              px-4 py-2 text-white rounded-md cursor-pointer font-medium
              ${isSuccess ? 'bg-green-500 hover:bg-green-600' : ''}
              ${isError ? 'bg-red-500 hover:bg-red-600' : ''}
            `}
          >
            {isError ? 'Try Again' : 'Start New Export'}
          </button>
        </div>
      )}
      <div className="overflow-y-auto flex-1">
        {progressLogs.map((log, index) => (
          <div 
            key={index} 
            className={`mb-2 ${log.startsWith('Error:') ? 'text-red-600 font-semibold' : ''}`}
          >
            {log}
          </div>
        ))}
      </div>
    </div>
  );
}
export default ExportLogger;