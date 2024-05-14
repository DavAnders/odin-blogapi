import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import { AuthProvider } from "./auth/AuthContext";
import { useAuth } from "./auth/UseAuth";
import LoginForm from "./loginForm";
import RegistrationForm from "./registerForm";
import Dashboard from "./dashboard";
import Posts from "./posts";
import UserList from "./userList";
import CreatePost from "./createPost";
import PostDetail from "./PostDetail";
import EditPost from "./EditPost";
import PropTypes from "prop-types";

function App() {
  return (
    <Router>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<LoginForm />} />
          <Route path="/register" element={<RegistrationForm />} />
          <Route
            path="/dashboard"
            element={<ProtectedElement component={Dashboard} />}
          />
          <Route
            path="/users"
            element={<ProtectedElement component={UserList} />}
          />
          <Route
            path="/posts"
            element={<ProtectedElement component={Posts} />}
          />
          <Route path="*" element={<Navigate to="/login" />} />
          <Route
            path="/create-post"
            element={<ProtectedElement component={CreatePost} />}
          />
          <Route
            path="/posts/:id"
            element={<ProtectedElement component={PostDetail} />}
          />
          <Route
            path="/edit-post/:id"
            element={<ProtectedElement component={EditPost} />}
          />
        </Routes>
      </AuthProvider>
    </Router>
  );
}

function ProtectedElement({ component: Component }) {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? <Component /> : <Navigate to="/login" />;
}

ProtectedElement.propTypes = {
  component: PropTypes.elementType.isRequired,
};
export default App;
