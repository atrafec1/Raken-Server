import { Link } from "react-router-dom";

function Navbar() {
  return (
    <nav className="bg-neutral-800 border-b border-neutral-700 px-6 py-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Link to="/"className="text-white font-semibold text-lg">PRG Tools </Link>
        </div>
        <Link to="/payroll"> Payroll </Link> 
        <div className="flex items-center gap-3">
          <span className="text-neutral-400 text-sm">v1.0.0</span>
        </div>
      </div>
    </nav>
  );
}

export default Navbar;