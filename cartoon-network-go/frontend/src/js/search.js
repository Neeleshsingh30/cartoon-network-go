const BASE_URL = "https://cartoon-network-go-1.onrender.com";
const query = new URLSearchParams(window.location.search).get("query");

/* ==== COMMON IMAGE PICKER ==== */
function getCartoonImage(c) {
  let thumb = "/src/images/CN-BG-AUTH.jpg";

  if (c.Images && c.Images.length) {
    const img =
      c.Images.find(i => i.image_type === "thumbnail") ||
      c.Images.find(i => i.image_type === "poster");

    if (img && img.image_url) thumb = img.image_url;
  }
  return thumb;
}

/* ==== LOAD SEARCH RESULTS ==== */
async function loadResults() {
  try {
    const res = await fetch(`${BASE_URL}/cartoons/search?query=${encodeURIComponent(query)}`);
    const data = await res.json();

    const grid = document.getElementById("resultGrid");
    grid.innerHTML = "";

    if (!data.length) {
      grid.innerHTML = "<p>No cartoons found </p>";
      return;
    }

    data.forEach(c => {
      const thumb = getCartoonImage(c);

      const card = document.createElement("div");
      card.className = "cartoon-card";
      card.onclick = () => window.location.href = `cartoon.html?id=${c.ID}`;

      card.innerHTML = `
        <img src="${thumb}">
        <p>${c.Name}</p>
      `;

      grid.appendChild(card);
    });

  } catch (err) {
    console.error("Search load error:", err);
  }
}

window.onload = loadResults;
