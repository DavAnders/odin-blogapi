import { useState, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "./auth/AuthContext";

function LoginForm() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const { login } = useContext(AuthContext);

  const handleLogin = async (event) => {
    event.preventDefault();
    setError("");
    try {
      const response = await fetch(import.meta.env.VITE_API_URL + "/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });

      if (response.ok) {
        const data = await response.json();
        console.log("Login Successful:", data);
        login(data.token);
      } else {
        const errorData = await response.text();
        console.error("Login Failed:", errorData);
        setError(`Login failed: ${errorData}`);
      }
    } catch (error) {
      console.error("Error during login:", error);
      setError("Network error, please try again later.");
    }
  };

  return (
    <form onSubmit={handleLogin}>
      <div className="form-group">
        <label>
          Username:
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </label>
      </div>
      <div className="form-group">
        <label>
          Password:
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </label>
      </div>
      <button type="submit">Login</button>
      <button type="button" onClick={() => navigate("/register")}>
        Register
      </button>
      {error && <p style={{ color: "red" }}>{error}</p>}
    </form>
  );
}

export default LoginForm;
