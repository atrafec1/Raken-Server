import { Outlet } from "react-router-dom";
import Navbar from "./Navbar";
import "../index.css";
function Layout() {
  return (
    <div className="flex h-screen ">
      <Navbar />
      
      {/* Main Content Area */}
      <main className="flex-1 overflow-auto">
        <Outlet />
      </main>
    </div>
  );
}

export default Layout;