import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode"; // Corrected the import statement
import PropTypes from "prop-types";

const Navbar = ({ onLogout }) => {
  return (
    <div
      style={{
        padding: "10px",
        backgroundColor: "#f0f0f0",
        display: "flex",
        justifyContent: "space-between",
      }}
    >
      <span>My Dashboard</span>
      <button onClick={onLogout}>Logout</button>
    </div>
  );
};

Navbar.propTypes = {
  onLogout: PropTypes.func.isRequired,
};

const UserInfo = ({ user }) => {
  return (
    <div style={{ margin: "20px", padding: "10px", border: "1px solid gray" }}>
      <h4>User Information</h4>
      <p>Username: {user.username}</p>
      <p>User ID: {user.userId}</p>
    </div>
  );
};

UserInfo.propTypes = {
  user: PropTypes.object.isRequired,
};

const PostsList = ({ title, posts }) => {
  const navigate = useNavigate();

  return (
    <div style={{ margin: "20px", padding: "10px", border: "1px solid gray" }}>
      <h4>{title}</h4>
      {posts.length > 0 ? (
        <ul>
          {posts.map((post) => (
            <li key={post.id}>
              <span
                style={{ cursor: "pointer", color: "blue" }}
                onClick={() => navigate(`/posts/${post.id}`)}
              >
                {post.title}
              </span>{" "}
              - {new Date(post.publishedAt).toLocaleDateString()}
            </li>
          ))}
        </ul>
      ) : (
        <p>No posts found.</p>
      )}
    </div>
  );
};

PostsList.propTypes = {
  title: PropTypes.string.isRequired,
  posts: PropTypes.array.isRequired,
};

const Dashboard = () => {
  const navigate = useNavigate();
  const [user, setUser] = useState(null);
  const [userPosts, setUserPosts] = useState([]);
  const [recentPosts, setRecentPosts] = useState([]);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      console.log("No token found, navigating to login");
      navigate("/login");
      return;
    }

    try {
      const decodedUser = jwtDecode(token);
      if (decodedUser.exp * 1000 < Date.now()) {
        console.log("Token is expired, navigating to login");
        navigate("/login");
        return;
      }
      setUser(decodedUser);
      fetchUserPosts(token, decodedUser.userId);
      fetchRecentPosts(token);
    } catch (error) {
      console.error("Failed to decode token or token is expired:", error);
      localStorage.removeItem("token");
      navigate("/login");
    }
  }, [navigate]);

  const fetchUserPosts = (token, userId) => {
    fetch(`${import.meta.env.VITE_API_URL}/api/posts/user/${userId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((response) => response.json())
      .then((data) => setUserPosts(Array.isArray(data) ? data : []))
      .catch((error) => {
        console.error("Error fetching user posts:", error);
        setUserPosts([]);
      });
  };

  const fetchRecentPosts = (token) => {
    fetch(`${import.meta.env.VITE_API_URL}/api/posts`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((response) => response.json())
      .then((data) => setRecentPosts(Array.isArray(data) ? data : []))
      .catch((error) => {
        console.error("Error fetching recent posts:", error);
        setRecentPosts([]);
      });
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/login");
  };

  const handleCreatePost = () => {
    navigate("/create-post");
  };

  if (!user) {
    return <div>Loading user info...</div>; // Return a loading indicator if user data is not yet available
  }

  return (
    <div>
      <Navbar onLogout={handleLogout} />
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
        }}
      >
        <UserInfo user={user} />
        <button
          onClick={handleCreatePost}
          style={{ margin: "10px", padding: "10px 20px" }}
        >
          Create New Post
        </button>
        <PostsList title="Your Recent Posts" posts={userPosts} />
        <PostsList title="All Recent Posts" posts={recentPosts} />
      </div>
    </div>
  );
};

export default Dashboard;
