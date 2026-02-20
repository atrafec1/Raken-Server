import { useState } from 'react';
import ReportsPage from "./pages/ReportsPage";
import './index.css';
import Navbar from "./components/Navbar";
import { HashRouter as Router, Routes, Route } from "react-router-dom";
import PayrollPage from "./pages/PayrollPage";
import MaterialPage from './pages/MaterialPage';
import Layout from "./components/Layout";

function App() {

    return (
        <Router>
            <Routes>
                <Route element={<Layout />}>
                <Route path="/reports" element={<ReportsPage />} />
                <Route path="/" element={<PayrollPage />} />
                <Route path="/materials" element={<MaterialPage />} />
                </Route>
            </Routes>
        </Router>
    );
}

export default App;
