const BASE_URL = "https://cartoon-network-go-1.onrender.com";
const token = localStorage.getItem("admin_token");

if (!token) window.location.href = "admin.html";

function goDashboard() {
  window.location.href = "admin_dashboard.html";
}

function logout() {
  localStorage.clear();
  window.location.href = "admin.html";
}

/* ================= ADD CARTOON ================= */
async function addCartoon() {

  const airDateValue = document.getElementById("air_date").value;

  const body = {
    name: document.getElementById("name").value.trim(),
    genre: document.getElementById("genre").value.trim(),
    age_group: document.getElementById("age_group").value.trim(),
    universe: document.getElementById("universe").value.trim(),
    show_time: document.getElementById("show_time").value.trim(),
    imdb_rating: Number(document.getElementById("imdb_rating").value) || 0,
    description: document.getElementById("description").value.trim(),
    air_date: airDateValue ? airDateValue : null
  };

  console.log("FINAL PAYLOAD:", body);

  const res = await fetch(`${BASE_URL}/admin/cartoon`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`
    },
    body: JSON.stringify(body)
  });

  const data = await res.json();

  if (!res.ok) {
    alert(data.error || "Failed to add cartoon");
    return;
  }

  alert("Cartoon added successfully");
  loadCartoons();
}



/* ================= UPLOAD CARTOON IMAGE (URL) ================= */


async function uploadCartoonImage() {

  const cartoonId = document.getElementById("imgCartoonId").value.trim();
  const imageType = document.getElementById("imgType").value;
  const imageUrl  = document.getElementById("imgUrl").value.trim();

  console.log("DEBUG →", { cartoonId, imageType, imageUrl });

  if (!cartoonId) {
    alert("Cartoon ID required");
    return;
  }

  if (!imageType) {
    alert("Image type required");
    return;
  }

  if (!imageUrl) {
    alert("Image URL required");
    return;
  }

  const body = {
    cartoon_id: Number(cartoonId),
    image_type: imageType,
    image_url: imageUrl
  };

  try {
    const res = await fetch(`${BASE_URL}/admin/cartoon/upload-image`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify(body)
    });

    const data = await res.json();

    if (!res.ok) {
      alert(data.error || "Upload failed");
      return;
    }

    alert("Image uploaded successfully");
  } catch (err) {
    console.error("Upload error:", err);
    alert("Server not reachable");
  }
}

/* ================= ADD CHARACTER ================= */
async function addCharacter() {
  await fetch(`${BASE_URL}/admin/cartoon/${char_cartoon_id.value}/character`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`
    },
    body: JSON.stringify({
      name: char_name.value,
      role: char_role.value,
      power: char_power.value,
      description: char_description.value
    })
  });

  alert("Character added");
}

/* ================= LOAD CARTOONS ================= */
async function loadCartoons() {
  const res = await fetch(`${BASE_URL}/admin/cartoons`, {
    headers: { Authorization: `Bearer ${token}` }
  });

  const data = await res.json();
  cartoonTable.innerHTML = "";

  data.cartoons.forEach(c => {
    cartoonTable.innerHTML += `
      <tr>
        <td>${c.id}</td>
        <td><img src="${c.thumbnail || '../../src/images/CN-BG-AUTH.jpg'}"></td>
        <td>${c.name}</td>
        <td><button onclick="deleteCartoon(${c.id})">Delete</button></td>
      </tr>
    `;
  });
}

/* ================= DELETE CARTOON ================= */
async function deleteCartoon(id) {
  if (!confirm("Delete cartoon?")) return;

  await fetch(`${BASE_URL}/admin/cartoon/${id}`, {
    method: "DELETE",
    headers: { Authorization: `Bearer ${token}` }
  });

  loadCartoons();
}

loadCartoons();
