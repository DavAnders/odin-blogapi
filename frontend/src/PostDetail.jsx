import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import Comments from "./Comments";

function PostDetail() {
  const { id } = useParams();
  const [post, setPost] = useState(null);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    const fetchPost = async () => {
      try {
        const token = localStorage.getItem("token");
        const response = await fetch(
          `${import.meta.env.VITE_API_URL}/api/posts/${id}`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        if (!response.ok) throw new Error("Failed to fetch post");
        const data = await response.json();
        setPost(data);
      } catch (error) {
        console.error("Fetch error:", error);
        setError("Failed to load post");
      }
    };

    fetchPost();
  }, [id]);

  if (error) return <p style={{ color: "red" }}>{error}</p>;
  if (!post) return <p>Loading...</p>;

  return (
    <div>
      <button onClick={() => navigate("/posts")}>Back to Posts</button>
      <h2>{post.title}</h2>
      <p>{post.content}</p>
      <small>
        Published by {post.authorUsername} on{" "}
        {new Date(post.publishedAt).toLocaleString()}
      </small>
      <Comments postId={id} />
    </div>
  );
}

export default PostDetail;
