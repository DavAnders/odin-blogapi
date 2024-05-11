import { useState } from "react";
import { useNavigate } from "react-router-dom";

function RegistrationForm() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [email, setEmail] = useState("");
  const navigate = useNavigate();

  const handleRegister = async (event) => {
    event.preventDefault();
    try {
      const response = await fetch(import.meta.env.VITE_API_URL + "/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password, email }),
      });
      const data = await response.json();
      if (response.ok) {
        console.log("Registration Successful:", data);
        localStorage.setItem("token", data.token); // Store the token in localStorage
        navigate("/dashboard"); // Redirect to the dashboard
      } else {
        console.error("Registration Failed:", data);
      }
    } catch (error) {
      console.error("Error during registration:", error);
    }
  };

  return (
    <form onSubmit={handleRegister}>
      <label>
        Username:
        <input
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
      </label>
      <label>
        Password:
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
      </label>
      <label>
        Email:
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
      </label>
      <button type="submit">Register</button>
    </form>
  );
}

export default RegistrationForm;
