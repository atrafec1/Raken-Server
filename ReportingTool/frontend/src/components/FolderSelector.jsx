import React from "react";
import { SelectFolder } from "../../wailsjs/go/main/App";

function FolderSelector({ path, onSelect }) {
  const handleSelect = async () => {
    try {
      const selectedPath = await SelectFolder();
      if (selectedPath) {
        onSelect(selectedPath);
      }
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div className="flex items-center gap-3">
      <button
        onClick={handleSelect}
        className="px-4 py-2 bg-primary-600 text-white font-medium rounded-lg 
                   hover:bg-primary-700 active:bg-primary-800 
                   transition-colors duration-200 
                   focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2
                   shadow-sm hover:shadow-md
                   flex-shrink-0"
      >
        Select Folder
      </button>
      
      {path && (
        <div className="flex items-center gap-2 px-4 py-3 bg-neutral-50 border border-neutral-200 rounded-lg flex-1">
          <svg className="w-5 h-5 text-neutral-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
          </svg>
          <p className="text-sm text-neutral-700 font-mono truncate" title={path}>
            {path}
          </p>
        </div>
      )}
    </div>
  );
}

export default FolderSelector;