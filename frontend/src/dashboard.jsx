import { jwtDecode } from "jwt-decode";
import { useNavigate } from "react-router-dom";

const Navbar = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/login");
  };

  return (
    <div
      style={{
        padding: "10px",
        backgroundColor: "#f0f0f0",
        display: "flex",
        justifyContent: "space-between",
      }}
    >
      <span>My Dashboard</span>
      <button onClick={handleLogout}>Logout</button>
    </div>
  );
};

const UserInfo = () => {
  const token = localStorage.getItem("token");
  const user = token ? jwtDecode(token) : null; // Correctly using jwt_decode

  return (
    <div style={{ margin: "20px", padding: "10px", border: "1px solid gray" }}>
      <h4>User Information</h4>
      {user ? (
        <>
          <p>Username: {user.username}</p>
          <p>User ID: {user.userId}</p>
        </>
      ) : (
        <p>No user data available.</p>
      )}
    </div>
  );
};

const RecentActivities = () => {
  return (
    <div style={{ margin: "20px", padding: "10px", border: "1px solid gray" }}>
      <h4>Recent Activities</h4>
      <ul>
        <li>Logged in</li>
        <li>Viewed dashboard</li>
        <li>Edited profile</li>
      </ul>
    </div>
  );
};

const MainContent = () => {
  return (
    <div style={{ margin: "20px", padding: "20px", border: "1px solid gray" }}>
      <h2>Welcome to Your Dashboard!</h2>
      <p>
        This is the main area where specific dashboard content can be placed.
      </p>
    </div>
  );
};

const Dashboard = () => {
  return (
    <div>
      <Navbar />
      <div style={{ display: "flex", flexDirection: "column" }}>
        <UserInfo />
        <RecentActivities />
        <MainContent />
      </div>
    </div>
  );
};

export default Dashboard;
