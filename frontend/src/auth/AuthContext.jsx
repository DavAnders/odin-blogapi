// authcontext.jsx
import { createContext, useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import PropTypes from "prop-types";
import { jwtDecode } from "jwt-decode";

const checkTokenValidity = (token) => {
  if (!token) return false;

  try {
    const decoded = jwtDecode(token);
    const currentTime = Date.now() / 1000;
    return decoded.exp > currentTime;
  } catch (error) {
    console.error("Token decoding error:", error);
    return false;
  }
};

const AuthContext = createContext(null);

const AuthProvider = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const initializeAuth = () => {
      const token = localStorage.getItem("token");
      console.log("AuthProvider useEffect called");
      console.log("Token from localStorage on mount:", token);
      if (checkTokenValidity(token)) {
        console.log("Token is valid");
        const decodedUser = jwtDecode(token);
        setUser(decodedUser);
        setIsAuthenticated(true);
      } else {
        console.log("Token is invalid or expired");
        setIsAuthenticated(false);
        setUser(null);
        localStorage.removeItem("token");
        navigate("/login");
      }
      setLoading(false);
    };

    initializeAuth();
  }, [navigate]);

  const login = (token) => {
    console.log("Login called with token:", token);
    if (checkTokenValidity(token)) {
      localStorage.setItem("token", token);
      const decodedUser = jwtDecode(token);
      setUser(decodedUser);
      setIsAuthenticated(true);
      navigate("/dashboard");
    }
  };

  const logout = () => {
    console.log("Logout called");
    localStorage.removeItem("token");
    setIsAuthenticated(false);
    setUser(null);
    navigate("/login");
  };

  return (
    <AuthContext.Provider
      value={{ isAuthenticated, user, login, logout, loading }}
    >
      {children}
    </AuthContext.Provider>
  );
};

AuthProvider.propTypes = {
  children: PropTypes.node.isRequired,
};

export { AuthContext, AuthProvider };
