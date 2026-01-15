const BASE_URL = "https://cartoon-network-go-1.onrender.com";
const token = localStorage.getItem("admin_token");
const role = localStorage.getItem("admin_role");

/* ================= AUTH + ROLE CHECK ================= */
if (!token) {
  window.location.href = "admin.html";
}

if (role !== "super_admin") {
  alert(
    "You cannot manage admins.\nOnly super admin can do this.\nPlease contact super admin."
  );
  window.location.href = "admin_dashboard.html";
}

/* ================= NAVIGATION ================= */
function goDashboard() {
  window.location.href = "admin_dashboard.html";
}

function logout() {
  localStorage.removeItem("admin_token");
  localStorage.removeItem("admin_role");
  window.location.href = "admin.html";
}

/* ================= CREATE ADMIN ================= */
async function createAdmin() {
  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value.trim();
  const roleValue = document.getElementById("role").value;

  const error = document.getElementById("error");
  const success = document.getElementById("success");

  error.innerText = "";
  success.innerText = "";

  if (!username || !password) {
    error.innerText = "Username and password are required";
    return;
  }

  try {
    const res = await fetch(`${BASE_URL}/admin/create-admin`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({ username, password, role: roleValue })
    });

    const data = await res.json();

    if (!res.ok) {
      error.innerText = data.error || "Failed to create admin";
      return;
    }

    success.innerText = "Admin created successfully ✅";
    document.getElementById("username").value = "";
    document.getElementById("password").value = "";

    loadAdmins();

  } catch (err) {
    error.innerText = "Backend server not reachable";
  }
}

/* ================= LOAD ADMINS ================= */
async function loadAdmins() {
  try {
    const res = await fetch(`${BASE_URL}/admin/list`, {
      headers: { Authorization: `Bearer ${token}` }
    });

    const data = await res.json();
    const table = document.getElementById("adminTable");
    table.innerHTML = "";

    (data.admins || []).forEach(a => {
      table.innerHTML += `
        <tr>
          <td>${a.id}</td>
          <td>${a.username}</td>
          <td>${a.role}</td>
          <td>
            ${
              a.role === "super_admin"
                ? "<span>Protected</span>"
                : `<button onclick="deleteAdmin(${a.id})">Delete</button>`
            }
          </td>
        </tr>
      `;
    });

  } catch (err) {
    console.error(err);
  }
}

/* ================= DELETE ADMIN ================= */
async function deleteAdmin(id) {
  if (!confirm("Are you sure you want to delete this admin?")) return;

  try {
    const res = await fetch(`${BASE_URL}/admin/delete-admin/${id}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${token}`
      }
    });

    const data = await res.json();

    if (!res.ok) {
      alert(data.error || "Delete failed");
      return;
    }

    alert("Admin deleted successfully");
    loadAdmins();

  } catch (err) {
    alert("Backend server not reachable");
  }
}

/* ================= INIT ================= */
window.onload = function () {
  loadAdmins();
};
