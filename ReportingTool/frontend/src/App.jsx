import { useState } from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import { Greet, ExportReports } from "../wailsjs/go/main/App";

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
            <img src={logo} id="logo" alt="logo" />
            <div id="result" className="result">{resultText}</div>

            <div id="input" className="input-box">
                <input
                    id="name"
                    className="input"
                    onChange={updateName}
                    autoComplete="off"
                    type="text"
                    placeholder="Enter your name"
                />
                <button className="btn" onClick={greet}>Greet</button>
            </div>

            <div className="date-inputs">
                <input
                    type="text"
                    placeholder="fromDate"
                    value={fromDate}
                    onChange={updateFromDate}
                />
                <input
                    type="text"
                    placeholder="toDate"
                    value={toDate}
                    onChange={updateToDate}
                />
                <button className="btn" onClick={exportReportsToComputer}>Export Reports</button>
            </div>
        </div>
    );
}

export default App;
