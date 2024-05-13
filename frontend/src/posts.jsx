import { useState, useEffect } from "react";

function Posts() {
  const [posts, setPosts] = useState([]);
  const [error, setError] = useState("");

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
            /* Inline for now */
          >
            <h3>{post.title}</h3>
            <p>
              {post.content.substring(0, 200)}
              {post.content.length > 200 ? "..." : ""}
            </p>{" "}
            {/* Display only the first 200 characters */}
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
