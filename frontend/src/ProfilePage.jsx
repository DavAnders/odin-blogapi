import { useState, useEffect, useContext } from "react";
import { AuthContext } from "./auth/AuthContext";
import { useNavigate } from "react-router-dom";

function ProfilePage() {
  const { isAuthenticated, logout } = useContext(AuthContext);
  const [profile, setProfile] = useState({ bio: "", profilePicUrl: "" });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    const fetchProfile = async () => {
      if (!isAuthenticated) {
        setError("You need to be logged in to view your profile.");
        setLoading(false);
        return;
      }

      try {
        const response = await fetch(
          `${import.meta.env.VITE_API_URL}/api/profile`,
          {
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
          }
        );

        if (response.ok) {
          const data = await response.json();
          setProfile({
            bio: data.bio || "",
            profilePicUrl: data.profilePicUrl || "",
          });
          setLoading(false);
        } else {
          const errorText = await response.text();
          console.error("Failed to fetch profile:", errorText);
          setError(errorText || "Failed to fetch profile");
          setLoading(false);
        }
      } catch (error) {
        console.error("Error fetching profile:", error);
        setError("An unexpected error occurred. Please try again later.");
        setLoading(false);
      }
    };

    fetchProfile();
  }, [isAuthenticated]);

  const handleUpdateProfile = async (event) => {
    event.preventDefault();

    try {
      const response = await fetch(
        `${import.meta.env.VITE_API_URL}/api/profile`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
          body: JSON.stringify(profile),
        }
      );

      if (response.ok) {
        const data = await response.json();
        setProfile({
          bio: data.bio || "",
          profilePicUrl: data.profilePicUrl || "",
        });
        setError("");
      } else {
        const errorText = await response.text();
        console.error("Failed to update profile:", errorText);
        setError(errorText || "Failed to update profile");
      }
    } catch (error) {
      console.error("Error updating profile:", error);
      setError("An unexpected error occurred. Please try again later.");
    }
  };

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      <h1>Profile Page</h1>
      {error && <p style={{ color: "red" }}>{error}</p>}
      {isAuthenticated ? (
        <form onSubmit={handleUpdateProfile}>
          <div className="form-group">
            <label>
              Bio:
              <textarea
                value={profile.bio}
                onChange={(e) =>
                  setProfile({ ...profile, bio: e.target.value })
                }
              />
            </label>
          </div>
          <div className="form-group">
            <label>
              Profile Picture URL:
              <input
                type="text"
                value={profile.profilePicUrl}
                onChange={(e) =>
                  setProfile({ ...profile, profilePicUrl: e.target.value })
                }
              />
            </label>
          </div>
          <div className="updateButton" style={{ marginBottom: "10px" }}>
            <button type="submit">Update Profile</button>
          </div>
        </form>
      ) : (
        <p>Please log in to view and edit your profile.</p>
      )}
      <div className="navButtons">
        {isAuthenticated && <button onClick={logout}>Logout</button>}
        {isAuthenticated && (
          <button onClick={() => navigate("/dashboard")}>Dashboard</button>
        )}
      </div>
    </div>
  );
}

export default ProfilePage;
