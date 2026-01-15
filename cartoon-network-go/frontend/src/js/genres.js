const BASE_URL = "https://cartoon-network-go-1.onrender.com";

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

/* ==== LOAD GENRE CARTOONS ==== */
async function loadGenres() {
  try {
    const res = await fetch(`${BASE_URL}/cartoons/by-genre`);
    const data = await res.json();

    const container = document.getElementById("genreContainer");
    container.innerHTML = "";

    Object.keys(data).forEach(genre => {

      const safeId = genre.replace(/\s+/g, "-").toLowerCase();

      const section = document.createElement("div");
      section.className = "genre-section";

      section.innerHTML = `
        <h2 class="genre-title">${genre}</h2>
        <div class="genre-row" id="row-${safeId}"></div>
      `;

      container.appendChild(section);

      const row = document.getElementById(`row-${safeId}`);

      data[genre].forEach(c => {

        const thumb = getCartoonImage(c);

        const card = document.createElement("div");
        card.className = "genre-card";
        card.onclick = () => window.location.href = `cartoon.html?id=${c.ID}`;

        card.innerHTML = `
          <img src="${thumb}">
          <p>${c.Name}</p>
        `;

        row.appendChild(card);
      });
    });

  } catch (err) {
    console.error("Genre load error:", err);
  }
}

window.onload = loadGenres;
