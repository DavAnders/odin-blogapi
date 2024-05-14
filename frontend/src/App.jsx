import {
  BrowserRouter as Router,
  Route,
  Routes,
  Navigate,
} from "react-router-dom";
import LoginForm from "./loginForm";
import RegistrationForm from "./registerForm";
import Dashboard from "./dashboard";
import "./App.css";
import Posts from "./posts";
import UserList from "./userList";
import CreatePost from "./createPost";
import PostDetail from "./PostDetail";

function App() {
  // Simple function to check if user is authenticated
  const isAuthenticated = () => !!localStorage.getItem("token");

  return (
    <Router>
      <div>
        <Routes>
          <Route path="/login" element={<LoginForm />} />
          <Route path="/register" element={<RegistrationForm />} />
          <Route
            path="/dashboard"
            element={
              isAuthenticated() ? <Dashboard /> : <Navigate to="/login" />
            }
          />
          <Route path="/users" element={<UserList />} />
          <Route path="/posts" element={<Posts />} />
          <Route path="*" element={<Navigate to="/login" />} />
          <Route path="/create-post" element={<CreatePost />} />
          <Route path="/posts/:id" element={<PostDetail />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
