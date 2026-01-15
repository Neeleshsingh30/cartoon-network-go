/* ================= ADMIN LOGIN ================= */

const BASE_URL = "https://cartoon-network-go-1.onrender.com"; // loaded from .env at build time

async function adminLogin() {
  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value.trim();
  const error = document.getElementById("error");

  error.innerText = "";

  if (!username || !password) {
    error.innerText = "Username and password required";
    return;
  }

  try {
    const response = await fetch(`${BASE_URL}/admin/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ username, password })
    });

    const data = await response.json();

    if (!response.ok) {
      error.innerText = data.error || "Invalid admin credentials";
      return;
    }

    /*  SAVE AUTH DATA */
    localStorage.setItem("admin_token", data.token);
    localStorage.setItem("admin_role", data.role);

    /*  REDIRECT */
    window.location.href = "admin_dashboard.html";

  } catch (err) {
    console.error(err);
    error.innerText = "Backend server not reachable";
  }
}

/* ================= NAVIGATION (OPTIONAL REUSE) ================= */

function goToManage() {
  window.location.href = "manage_cartoon.html";
}

function goToCreateAdmin() {
  window.location.href = "create_admin.html"; // if exists
}

function logout() {
  localStorage.removeItem("admin_token");
  localStorage.removeItem("admin_role");
  window.location.href = "admin.html";
}
