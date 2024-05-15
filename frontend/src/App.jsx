import { Routes, Route } from "react-router-dom";
import { AuthProvider } from "./auth/AuthContext";
import ProtectedElement from "./auth/ProtectedElement";
import LoginForm from "./loginForm";
import RegistrationForm from "./registerForm";
import Dashboard from "./dashboard";
import Posts from "./posts";
import UserList from "./userList";
import CreatePost from "./createPost";
import PostDetail from "./PostDetail";
import EditPost from "./EditPost";

const App = () => {
  return (
    <AuthProvider>
      <Routes>
        <Route
          path="/login"
          element={
            <ProtectedElement component={LoginForm} isPublic restricted />
          }
        />
        <Route
          path="/register"
          element={
            <ProtectedElement
              component={RegistrationForm}
              isPublic
              restricted
            />
          }
        />
        <Route
          path="/dashboard"
          element={<ProtectedElement component={Dashboard} />}
        />
        <Route
          path="/users"
          element={<ProtectedElement component={UserList} />}
        />
        <Route path="/posts" element={<ProtectedElement component={Posts} />} />
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
        <Route
          path="*"
          element={<ProtectedElement component={LoginForm} isPublic />}
        />
      </Routes>
    </AuthProvider>
  );
};

export default App;
