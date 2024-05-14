import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";
import Comments from "./Comments";

const PostDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [post, setPost] = useState(null);
  const [error, setError] = useState("");
  const [user, setUser] = useState(null);

  useEffect(() => {
    const fetchPost = async () => {
      const token = localStorage.getItem("token");
      try {
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

    const fetchUser = () => {
      const token = localStorage.getItem("token");
      if (!token) return;
      const decodedToken = jwtDecode(token);
      setUser(decodedToken);
    };

    fetchPost();
    fetchUser();
  }, [id]);

  const handleEdit = () => {
    navigate(`/edit-post/${id}`);
  };

  if (error) return <p style={{ color: "red" }}>{error}</p>;

  return (
    <div>
      {post ? (
        <div>
          <h2>{post.title}</h2>
          <p>{post.content}</p>
          <p>
            Published by {post.authorUsername} on{" "}
            {new Date(post.publishedAt).toLocaleString()}
          </p>
          {user && user.userId === post.authorId && (
            <button onClick={handleEdit}>Edit</button>
          )}
          <Comments postId={id} />
        </div>
      ) : (
        <p>Loading post...</p>
      )}
    </div>
  );
};

export default PostDetail;
