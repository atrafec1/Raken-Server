import { Link, useLocation } from "react-router-dom";

function Navbar() {
  const location = useLocation();
  
  const isActive = (path) => location.pathname === path;
  
  const navItems = [
    { name: "Payroll", path: "/", disabled: false },
    { name: "Reports", path: "/reports", disabled: false },
  ];
  
  return (
    <aside className="w-56 bg-neutral-900 border-r border-neutral-800 flex flex-col h-screen">
      {/* Logo Section */}
      <div className="px-4 py-5 border-b border-neutral-800">
        <Link to="/" className="flex items-center gap-2">
        
          <span className="text-white font-semibold text-base">PRG Tools</span>
        </Link>
      </div>

      {/* Main Navigation */}
      <nav className="flex-1 px-2 py-4 space-y-1">
        {navItems.map((item) => (
          item.disabled ? (
            <div
              key={item.name}
              className="px-3 py-2 text-neutral-600 text-sm cursor-not-allowed"
            >
              {item.name}
            </div>
          ) : (
            <Link
              key={item.name}
              to={item.path}
              className={`
                block px-3 py-2 rounded text-sm transition-colors
                ${isActive(item.path)
                  ? "bg-neutral-800 text-white font-medium"
                  : "text-neutral-400 hover:text-white hover:bg-neutral-800/50"
                }
              `}
            >
              {item.name}
            </Link>
          )
        ))}
      </nav>
    </aside>
  );
}

export default Navbar;