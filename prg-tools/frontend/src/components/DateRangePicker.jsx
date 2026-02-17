function DateRangePicker({ fromDate, toDate, onChange }) {
  return (
    <div className="flex items-center gap-3">
      <div className="flex flex-col gap-1">
        <label className="text-sm font-medium text-secondary-200">
          From
        </label>
        <input
          type="date"
          value={fromDate}
          onChange={(e) => onChange(e.target.value, toDate)}
          className="px-4 py-2 border border-neutral-300 rounded-lg
                     text-neutral-900 bg-white
                     hover:border-neutral-400
                     focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent
                     transition-colors duration-200
                     cursor-pointer"
        />
      </div>
      
      <div className="flex flex-col gap-1">
        <label className="text-sm font-medium text-secondary-200">
          To
        </label>
        <input
          type="date"
          value={toDate}
          onChange={(e) => onChange(fromDate, e.target.value)}
          className="px-4 py-2 border border-neutral-300 rounded-lg
                     text-neutral-900 bg-white
                     hover:border-neutral-400
                     focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent
                     transition-colors duration-200
                     cursor-pointer"
        />
      </div>
    </div>
  );
}

export default DateRangePicker