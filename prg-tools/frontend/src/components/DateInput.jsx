function DateInput({ onChange, label, subLabel, value }) {
  return (
    <div className="flex flex-col gap-2 mb-4">
      {label && (
        <label className="font-semibold text-sm text-gray-800">
          {label}
        </label>
      )}
      
    
      <input
        type="date"
        value={value}
        onChange={onChange}
        className="px-3 py-2 border border-gray-300 rounded-md text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 hover:border-gray-400 transition-colors"
      />

        {subLabel && (
        <p className="text-sm text-gray-500 -mt-1">
          {subLabel}
        </p>
      )}
    </div>
  );
}

export default DateInput;