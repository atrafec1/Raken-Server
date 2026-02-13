

function ExportReportsButton({ text, onClick }) {
  return (
    <button
      onClick={onClick}
      className="px-6 py-2.5 bg-primary-600 text-white font-medium rounded-lg 
                 hover:bg-primary-700 active:bg-primary-800 
                 transition-colors duration-200 
                 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2
                 shadow-sm hover:shadow-md"
    >
      {text || "Export Reports"}
    </button>
  );
}

export default ExportReportsButton;