import { useState } from 'react';
import ReportsPage from "./pages/ReportsPage";
import './index.css';
import Navbar from "./components/Navbar";
import { HashRouter as Router, Routes, Route } from "react-router-dom";
import PayrollPage from "./pages/PayrollPage";
import HomePage from "./pages/HomePage";
import Layout from "./components/Layout";

function App() {

    return (
        <Router>
            <Routes>
                <Route element={<Layout />}>
                <Route path="/" element={<HomePage/>} />
                <Route path="/reports" element={<ReportsPage />} />
                <Route path="/payroll" element={<PayrollPage />} />
                </Route>
            </Routes>
        </Router>
    );
}

export default App;
