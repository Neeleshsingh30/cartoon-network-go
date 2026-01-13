/* ================= ADMIN LOGIN ================= */

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
    const response = await fetch("http://localhost:8000/admin/login", {
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

    /* ✅ SAVE AUTH DATA (KEYS MUST MATCH EVERYWHERE) */
    localStorage.setItem("admin_token", data.token);
    localStorage.setItem("admin_role", data.role);

    /* ✅ REDIRECT TO DASHBOARD */
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
