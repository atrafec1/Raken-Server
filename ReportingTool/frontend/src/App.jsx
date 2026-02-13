import { useState } from 'react';
import ReportsPage from "./pages/ReportsPage";
import './index.css';
import Navbar from "./components/Navbar";
import { HashRouter as Router, Routes, Route } from "react-router-dom";
import PayrollPage from "./pages/PayrollPage";
import HomePage from "./pages/HomePage";
function App() {

    return (
        <Router>
            <Navbar />
        <div className="flex-1">
            <Routes>
                <Route path="/" element={<HomePage/>} />
                <Route path="/reports" element={<ReportsPage />} />
                <Route path="/payroll" element={<PayrollPage />} />
            </Routes>
        </div>
        </Router>
    );
}

export default App;
