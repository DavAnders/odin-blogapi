import { useState, useEffect } from "react";
import PropTypes from "prop-types";
import { jwtDecode } from "jwt-decode";

function Comments({ postId }) {
  const [comments, setComments] = useState([]);
  const [newComment, setNewComment] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState(null);

  useEffect(() => {
    const fetchComments = async () => {
      try {
        const token = localStorage.getItem("token");
        const response = await fetch(
          `${import.meta.env.VITE_API_URL}/api/comments/${postId}`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        if (!response.ok) throw new Error("Failed to fetch comments");
        const data = await response.json();
        setComments(Array.isArray(data) ? data : []);
      } catch (error) {
        console.error("Fetch error:", error);
        setError("Failed to load comments");
      } finally {
        setLoading(false);
      }
    };

    const fetchUser = async () => {
      const token = localStorage.getItem("token");
      if (!token) return;
      const decodedToken = jwtDecode(token);
      setUser(decodedToken);
    };

    fetchComments();
    fetchUser();
  }, [postId]);

  const handleCommentSubmit = async (e) => {
    e.preventDefault();
    const token = localStorage.getItem("token");

    try {
      const response = await fetch(
        `${import.meta.env.VITE_API_URL}/api/comments`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            postId,
            content: newComment,
          }),
        }
      );
      if (!response.ok) throw new Error("Failed to submit comment");
      const comment = await response.json();
      setComments((prevComments) =>
        Array.isArray(prevComments) ? [...prevComments, comment] : [comment]
      );
      setNewComment("");
    } catch (error) {
      console.error("Submit error:", error);
      setError("Failed to submit comment");
    }
  };

  const handleDeleteComment = async (commentId) => {
    const token = localStorage.getItem("token");

    try {
      const response = await fetch(
        `${import.meta.env.VITE_API_URL}/api/comments/${commentId}`,
        {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      if (!response.ok) throw new Error("Failed to delete comment");
      setComments((prevComments) =>
        prevComments.filter((comment) => comment.id !== commentId)
      );
    } catch (error) {
      console.error("Delete error:", error);
      setError("Failed to delete comment");
    }
  };

  if (loading) return <p>Loading comments...</p>;
  if (error) return <p style={{ color: "red" }}>{error}</p>;

  return (
    <div>
      <h3>Comments</h3>
      <ul style={{ listStyle: "none", padding: 0 }}>
        {Array.isArray(comments) &&
          comments.map((comment) => (
            <li key={comment.id} style={{ margin: "10px 0" }}>
              <p>{comment.content}</p>
              <small>
                {comment.author} on{" "}
                {new Date(comment.createdAt).toLocaleString()}
              </small>
              {user && user.userId === comment.authorId && (
                <button onClick={() => handleDeleteComment(comment.id)}>
                  Delete
                </button>
              )}
            </li>
          ))}
      </ul>
      <form onSubmit={handleCommentSubmit}>
        <textarea
          value={newComment}
          onChange={(e) => setNewComment(e.target.value)}
          rows="3"
          style={{ width: "100%" }}
        />
        <button type="submit">Add Comment</button>
      </form>
    </div>
  );
}

Comments.propTypes = {
  postId: PropTypes.string.isRequired,
};

export default Comments;
