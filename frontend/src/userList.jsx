import { useState, useEffect } from "react";

function UserList() {
  const [users, setUsers] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const token = localStorage.getItem("token");
        const response = await fetch(
          import.meta.env.VITE_API_URL + "/api/users",
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        if (!response.ok) throw new Error("Failed to fetch");
        const data = await response.json();
        setUsers(data);
      } catch (error) {
        console.error("Fetch error:", error);
        setError("Failed to load users");
      }
    };

    fetchUsers();
  }, []);

  return (
    <div>
      <h2>Check out all of our awesome users!</h2>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <ul>
        {users.map((user) => (
          <li key={user.createdAt}>{user.username}</li>
        ))}
      </ul>
    </div>
  );
}

export default UserList;
