const BASE_URL = "https://cartoon-network-go-1.onrender.com";

/* ==== COMMON IMAGE PICKER ==== */
function getCartoonImage(c) {
  let thumb = "/src/images/CN-BG-AUTH.jpg";

  if (c.Images && c.Images.length) {
    const img =
      c.Images.find(i => i.image_type === "banner") ||
      c.Images.find(i => i.image_type === "thumbnail") ||
      c.Images.find(i => i.image_type === "poster");

    if (img && img.image_url) thumb = img.image_url;
  }
  return thumb;
}

/* ==== LOAD ALL CARTOONS ==== */
async function loadAllCartoons() {
  try {
    const res = await fetch(`${BASE_URL}/cartoons`);
    const cartoons = await res.json();

    const grid = document.getElementById("allGrid");
    grid.innerHTML = "";

    cartoons.forEach(c => {

      const img = getCartoonImage(c);

      const card = document.createElement("div");
      card.className = "cartoon-box";
      card.onclick = () => {
        window.location.href = `cartoon.html?id=${c.ID}`;
      };

      card.innerHTML = `
        <img src="${img}">
        <h4>${c.Name}</h4>
      `;

      grid.appendChild(card);
    });

  } catch (err) {
    console.error("All cartoons load error:", err);
  }
}

window.onload = loadAllCartoons;
