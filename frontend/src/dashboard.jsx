import { useEffect, useState, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "./auth/AuthContext";
import PropTypes from "prop-types";

const Navbar = ({ onLogout }) => {
  const navigate = useNavigate();

  return (
    <div
      style={{
        padding: "10px",
        backgroundColor: "#f0f0f0",
        display: "flex",
        justifyContent: "space-between",
        alignItems: "center",
        color: "grey",
      }}
    >
      <h2 style={{ color: "black", marginRight: "5vw" }}>Dashboard</h2>
      <div style={{ marginLeft: "auto" }}>
        <button
          onClick={() => navigate("/profile")}
          style={{ marginRight: "10px" }}
        >
          Profile
        </button>
        <button onClick={onLogout}>Logout</button>
      </div>
    </div>
  );
};

Navbar.propTypes = {
  onLogout: PropTypes.func.isRequired,
};

const UserInfo = ({ user }) => {
  if (!user) {
    return null;
  }
  return (
    <div style={{ margin: "20px", textAlign: "center" }}>
      <h2>Welcome, {user.username}!</h2>
      {user.profilePicUrl && (
        <img
          src={user.profilePicUrl}
          alt="Profile"
          style={{ width: "150px", height: "150px", borderRadius: "50%" }}
        />
      )}
      {user.bio && <p>{user.bio}</p>}
    </div>
  );
};

UserInfo.propTypes = {
  user: PropTypes.object.isRequired,
};

const PostsList = ({ title, posts }) => {
  const navigate = useNavigate();
  return (
    <div
      style={{
        margin: "20px",
        minWidth: "200px",
        border: "1px solid gray",
        padding: "20px",
      }}
    >
      <h4>{title}</h4>
      {posts.length > 0 ? (
        <ul style={{ listStyle: "none", paddingInlineStart: "0" }}>
          {posts.map((post) => (
            <li key={post.id}>
              <span
                style={{
                  cursor: "pointer",
                  color: "violet",
                }}
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
  const { user, logout, loading } = useContext(AuthContext);
  const [userPosts, setUserPosts] = useState([]);
  const [recentPosts, setRecentPosts] = useState([]);
  const [profile, setProfile] = useState({ bio: "", profilePicUrl: "" });

  useEffect(() => {
    if (user) {
      fetchUserPosts(user.userId);
      fetchRecentPosts();
      fetchUserProfile();
    }
  }, [user]);

  const fetchUserPosts = (userId) => {
    const token = localStorage.getItem("token");
    fetch(`${import.meta.env.VITE_API_URL}/api/posts/user/${userId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((response) => response.json())
      .then((data) => {
        setUserPosts(Array.isArray(data) ? data : []);
      })
      .catch((error) => {
        console.error("Error fetching user posts:", error);
        setUserPosts([]);
      });
  };

  const fetchRecentPosts = () => {
    const token = localStorage.getItem("token");
    fetch(`${import.meta.env.VITE_API_URL}/api/posts`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((response) => response.json())
      .then((data) => {
        setRecentPosts(Array.isArray(data) ? data : []);
      })
      .catch((error) => {
        console.error("Error fetching recent posts:", error);
        setRecentPosts([]);
      });
  };

  const fetchUserProfile = () => {
    const token = localStorage.getItem("token");
    fetch(`${import.meta.env.VITE_API_URL}/api/profile`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((response) => response.json())
      .then((data) => {
        setProfile({
          bio: data.bio || "",
          profilePicUrl: data.profilePicUrl || "",
        });
      })
      .catch((error) => {
        console.error("Error fetching user profile:", error);
        setProfile({ bio: "", profilePicUrl: "" });
      });
  };

  const handleLogout = () => {
    logout();
  };

  const handleCreatePost = () => {
    navigate("/create-post");
  };

  if (loading) {
    return <div>Loading user info...</div>;
  }

  return (
    <div>
      <Navbar onLogout={handleLogout} />
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          padding: "20px",
        }}
      >
        <UserInfo user={{ ...user, ...profile }} />
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
