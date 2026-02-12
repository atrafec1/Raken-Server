import React, { useState } from 'react';

function DateRangePicker({ fromDate, toDate, onSubmit, buttonText }) {
  const [startDate, setStartDate] = useState(fromDate || '');
  const [endDate, setEndDate] = useState(toDate || '');

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(startDate, endDate);
  };

  const handleStartDateChange = (e) => setStartDate(e.target.value);
  const handleEndDateChange = (e) => setEndDate(e.target.value);

  return (
    <div>
      <form onSubmit={handleSubmit} className="flex items-center gap-2">
        <input
          type="date"
          value={startDate}
          onChange={handleStartDateChange}
          className="border rounded px-2 py-1"
        />
        <input
          type="date"
          value={endDate}
          onChange={handleEndDateChange}
          className="border rounded px-2 py-1"
        />
        <button
          type="submit"
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          {buttonText}
        </button>
      </form>
    </div>
  );
}

export default DateRangePicker;
