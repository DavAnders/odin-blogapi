import { createContext, useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import PropTypes from "prop-types";

const AuthContext = createContext(null);

const AuthProvider = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    // Check local storage for token and update authentication status
    const token = localStorage.getItem("token");
    setIsAuthenticated(!!token);
  }, []);

  // Function to handle login, setting token
  const login = (token) => {
    localStorage.setItem("token", token);
    setIsAuthenticated(true);
    navigate("/dashboard");
  };

  // Function to handle logout
  const logout = () => {
    localStorage.removeItem("token");
    setIsAuthenticated(false);
    navigate("/login");
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

AuthProvider.propTypes = {
  children: PropTypes.node.isRequired,
};

export { AuthContext, AuthProvider };
