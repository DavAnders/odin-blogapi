import { useContext } from "react";
import { Navigate } from "react-router-dom";
import PropTypes from "prop-types";
import { AuthContext } from "./AuthContext";

function ProtectedElement({ component: Component }) {
  const { isAuthenticated, loading } = useContext(AuthContext);
  if (loading) {
    return <div>Loading...</div>;
  }
  return isAuthenticated ? <Component /> : <Navigate to="/login" />;
}

ProtectedElement.propTypes = {
  component: PropTypes.elementType.isRequired,
};

export default ProtectedElement;
