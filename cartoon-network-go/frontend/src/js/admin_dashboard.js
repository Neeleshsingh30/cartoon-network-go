const BASE_URL = "http://localhost:8000";
const token = localStorage.getItem("admin_token");
const role = localStorage.getItem("admin_role");

/* ================= AUTH CHECK ================= */
if (!token) {
  window.location.href = "admin.html";
}

/* ================= LOAD ADMIN CARTOONS ================= */
async function loadAdminCartoons() {
  try {
    const res = await fetch(`${BASE_URL}/admin/cartoons`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });

    // 🔒 Token expired / invalid
    if (res.status === 401) {
      localStorage.removeItem("admin_token");
      localStorage.removeItem("admin_role");
      window.location.href = "admin.html";
      return;
    }

    const data = await res.json();
    console.log("ADMIN CARTOONS:", data);

    const cartoons = data.cartoons || [];
    const grid = document.getElementById("cartoonGrid");
    grid.innerHTML = "";

    if (!cartoons.length) {
      grid.innerHTML = "<p>No cartoons found</p>";
      return;
    }

    cartoons.forEach(c => {
      /* ✅ BACKEND DIRECT FIELD */
      const image =
        c.thumbnail && c.thumbnail.trim() !== ""
          ? c.thumbnail
          : "../../src/images/CN-BG-AUTH.jpg";

      const name = c.name || "Untitled Cartoon";

      const card = document.createElement("div");
      card.className = "cartoon-card";

      card.innerHTML = `
        <img src="${image}" alt="${name}" loading="lazy">
        <p>${name}</p>
      `;

      grid.appendChild(card);
    });

  } catch (err) {
    console.error("Admin cartoons load error:", err);
  }
}

/* ================= LOAD ADMIN LOGS ================= */
async function loadAdminLogs() {
  try {
    const res = await fetch(`${BASE_URL}/admin/logs`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });

    if (res.status === 401) {
      localStorage.removeItem("admin_token");
      localStorage.removeItem("admin_role");
      window.location.href = "admin.html";
      return;
    }

    const data = await res.json();
    console.log("ADMIN LOGS:", data);

    const logs = data.logs || [];
    const box = document.getElementById("adminLogs");
    box.innerHTML = "";

    if (!logs.length) {
      box.innerHTML = "<p>No activity yet</p>";
      return;
    }

    logs.forEach(l => {
      const div = document.createElement("div");
      div.className = "log-item";
      div.innerText =
        l.action ||
        l.Action ||
        l.message ||
        "Admin Activity";

      box.appendChild(div);
    });

  } catch (err) {
    console.error("Admin logs error:", err);
  }
}

/* ================= ACTIONS ================= */
function logout() {
  localStorage.removeItem("admin_token");
  localStorage.removeItem("admin_role");
  window.location.href = "admin.html";
}

/* ✅ NOW REAL CONNECTION */
function goToManage() {
  window.location.href = "manage_cartoon.html";
}

/* OPTIONAL – SUPER ADMIN PAGE */
function goToCreateAdmin() {
  if (role !== "super_admin") {
    alert("Only Super Admin can create admins");
    return;
  }
  window.location.href = "create_admin.html";
}

/* ================= INIT ================= */
window.onload = function () {
  loadAdminCartoons();
  loadAdminLogs();
};

