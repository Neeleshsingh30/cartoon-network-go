const BASE_URL = "http://localhost:8000";
const token = localStorage.getItem("admin_token");

if (!token) {
  window.location.href = "admin.html";
}

function goDashboard() {
  window.location.href = "admin_dashboard.html";
}

function logout() {
  localStorage.clear();
  window.location.href = "admin.html";
}

/* ================= ADD CARTOON ================= */
async function addCartoon() {
  const imageFile = document.getElementById("cartoonImage").files[0];

  let imageURL = "";

  if (imageFile) {
    const formData = new FormData();
    formData.append("image", imageFile);

    const imgRes = await fetch(`${BASE_URL}/admin/upload-image`, {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
      body: formData
    });

    const imgData = await imgRes.json();
    imageURL = imgData.image_url;
  }

  const body = {
    name: cartoonName.value,
    genre: genre.value,
    ageGroup: ageGroup.value,
    description: description.value
  };

  await fetch(`${BASE_URL}/admin/cartoon`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`
    },
    body: JSON.stringify(body)
  });

  loadCartoons();
}

/* ================= ADD CHARACTER ================= */
async function addCharacter() {
  const cartoonId = cartoonSelect.value;

  const body = {
    name: charName.value,
    role: charRole.value,
    power: charPower.value,
    description: charDesc.value
  };

  await fetch(`${BASE_URL}/admin/cartoon/${cartoonId}/character`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`
    },
    body: JSON.stringify(body)
  });

  alert("Character added successfully");
}

/* ================= LOAD CARTOONS ================= */
async function loadCartoons() {
  const res = await fetch(`${BASE_URL}/admin/cartoons`, {
    headers: { Authorization: `Bearer ${token}` }
  });

  const data = await res.json();
  const table = document.getElementById("cartoonTable");
  const select = document.getElementById("cartoonSelect");

  table.innerHTML = "";
  select.innerHTML = "";

  data.cartoons.forEach(c => {
    table.innerHTML += `
      <tr>
        <td>${c.id}</td>
        <td><img src="${c.thumbnail || '../../src/images/CN-BG-AUTH.jpg'}"></td>
        <td>${c.name}</td>
        <td><button onclick="deleteCartoon(${c.id})">Delete</button></td>
      </tr>
    `;

    select.innerHTML += `<option value="${c.id}">${c.name}</option>`;
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
