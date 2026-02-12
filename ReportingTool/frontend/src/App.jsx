import { useState } from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import { Greet, ExportReports } from "../wailsjs/go/main/App";
import DateRangePicker from "./components/DateRangePicker"

function App() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const [fromDate, setFromDate] = useState('');
    const [toDate, setToDate] = useState('');

    const updateName = (e) => setName(e.target.value);
    const updateFromDate = (e) => setFromDate(e.target.value);
    const updateToDate = (e) => setToDate(e.target.value);

    const greet = () => Greet(name).then(setResultText);
    const exportReportsToComputer = () => {
        ExportReports(fromDate, toDate)
            .then(() => console.log("Reports exported successfully"))
            .catch(err => console.error(err));
    };


    return (
        <div id="App">
            <DateRangePicker fromDate={fromDate} toDate={toDate} onFromDateChange={updateFromDate} onToDateChange={updateToDate} />
            <p className="text-blue-500"> Hello</p>
        </div>
    );
}

export default App;
