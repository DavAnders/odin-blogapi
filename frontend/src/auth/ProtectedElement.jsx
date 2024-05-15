import { useContext, useEffect } from "react";
import PropTypes from "prop-types";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "./AuthContext";

const ProtectedElement = ({
  component: Component,
  isPublic = false,
  restricted = false,
}) => {
  const { isAuthenticated, loading } = useContext(AuthContext);
  const navigate = useNavigate();

  useEffect(() => {
    if (!loading) {
      if (!isAuthenticated && !isPublic) {
        navigate("/login");
      } else if (isAuthenticated && restricted) {
        navigate("/dashboard");
      }
    }
  }, [isAuthenticated, loading, navigate, isPublic, restricted]);

  if (loading) return <div>Loading...</div>;
  return <Component />;
};

ProtectedElement.propTypes = {
  component: PropTypes.elementType.isRequired,
  isPublic: PropTypes.bool,
  restricted: PropTypes.bool,
};

export default ProtectedElement;
