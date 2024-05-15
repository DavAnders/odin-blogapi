import { useState, useContext } from "react";
import { AuthContext } from "./auth/AuthContext";
import DOMPurify from "dompurify";

function RegistrationForm() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [email, setEmail] = useState("");
  const { login } = useContext(AuthContext);
  const [error, setError] = useState("");

  const sanitizeInput = (input) => {
    return DOMPurify.sanitize(input);
  };

  const handleRegister = async (event) => {
    event.preventDefault();

    const sanitizedUsername = sanitizeInput(username);
    const sanitizedPassword = sanitizeInput(password);
    const sanitizedEmail = sanitizeInput(email);

    // Regex patterns
    const usernamePattern = /^[a-zA-Z0-9_]{3,15}$/;
    const passwordPattern = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/;
    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

    // Validation
    if (!usernamePattern.test(sanitizedUsername)) {
      setError(
        "Invalid username. It must be 3-15 characters long and can contain letters, numbers, and underscores."
      );
      return;
    }
    if (!passwordPattern.test(sanitizedPassword)) {
      setError(
        "Invalid password. It must be at least 8 characters long and contain at least one letter and one number."
      );
      return;
    }
    if (!emailPattern.test(sanitizedEmail)) {
      setError("Invalid email format.");
      return;
    }

    try {
      const response = await fetch(import.meta.env.VITE_API_URL + "/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username: sanitizedUsername,
          password: sanitizedPassword,
          email: sanitizedEmail,
        }),
        credentials: "include", // Necessary to ensure cookies are sent with the request
      });

      if (response.ok) {
        const data = await response.json();
        console.log("Registration Successful:", data);
        login(data.token); // Use the token returned in the response body
        setError(""); // Clear any previous error
      } else {
        const errorData = await response.json();
        setError(DOMPurify.sanitize(errorData.error)); // Set the error message from the response
      }
    } catch (error) {
      console.error("Error during registration:", error);
      setError("An unexpected error occurred. Please try again later.");
    }
  };

  return (
    <form onSubmit={handleRegister}>
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
      <div className="form-group">
        <label>
          Email:
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </label>
      </div>
      <button type="submit">Register</button>
      {error && <p style={{ color: "red" }}>{error}</p>}{" "}
      {/* Display error message */}
    </form>
  );
}

export default RegistrationForm;
