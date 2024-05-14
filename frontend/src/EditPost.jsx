import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";

const EditPost = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [post, setPost] = useState({ title: "", content: "" });
  const [error, setError] = useState("");

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

    fetchPost();
  }, [id]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setPost({ ...post, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const token = localStorage.getItem("token");

    try {
      const response = await fetch(
        `${import.meta.env.VITE_API_URL}/api/posts/${id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(post),
        }
      );
      if (!response.ok) throw new Error("Failed to update post");
      navigate(`/post/${id}`);
    } catch (error) {
      console.error("Submit error:", error);
      setError("Failed to update post");
    }
  };

  return (
    <div>
      <h2>Edit Post</h2>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <form onSubmit={handleSubmit}>
        <div>
          <label>
            Title:
            <input
              type="text"
              name="title"
              value={post.title}
              onChange={handleInputChange}
            />
          </label>
        </div>
        <div>
          <label>
            Content:
            <textarea
              name="content"
              value={post.content}
              onChange={handleInputChange}
            />
          </label>
        </div>
        <button type="submit">Save Changes</button>
      </form>
    </div>
  );
};

export default EditPost;
