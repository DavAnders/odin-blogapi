import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

function Posts() {
  const [posts, setPosts] = useState([]);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const formatDate = (dateString) => {
    const options = {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    };
    return new Date(dateString).toLocaleDateString(undefined, options);
  };

  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const token = localStorage.getItem("token");
        const response = await fetch(
          import.meta.env.VITE_API_URL + "/api/posts",
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        if (!response.ok) throw new Error("Failed to fetch posts");
        const data = await response.json();
        setPosts(data);
      } catch (error) {
        console.error("Fetch error:", error);
        setError("Failed to load posts");
      }
    };

    fetchPosts();
  }, []);

  return (
    <div>
      <div className="buttons">
        <button onClick={() => navigate("/dashboard")}>Dashboard</button>
        <button onClick={() => navigate("/create-post")}>Create Post</button>
      </div>
      <h2>Posts</h2>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <ul style={{ listStyle: "none", padding: 0 }}>
        {posts.map((post) => (
          <li
            key={post.id}
            style={{
              margin: "20px 0",
              padding: "10px",
              border: "1px solid #ccc",
              borderRadius: "8px",
            }}
          >
            <h3
              style={{ cursor: "pointer", color: "violet" }}
              onClick={() => navigate(`/posts/${post.id}`)}
            >
              {post.title}
            </h3>
            <p>
              {post.content.substring(0, 200)}
              {post.content.length > 200 ? "..." : ""}
            </p>
            <small>
              Published by {post.authorUsername} on{" "}
              {formatDate(post.publishedAt)}
            </small>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default Posts;
