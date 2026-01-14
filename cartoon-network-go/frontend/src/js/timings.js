const BASE_URL = "http://localhost:8000";

/* ==== COMMON IMAGE PICKER (same as home.js) ==== */
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

/* ==== LOAD CARTOON TIMINGS ==== */
async function loadTimings() {
  try {
    const res = await fetch(`${BASE_URL}/cartoons/timings`);
    const data = await res.json();

    const list = document.getElementById("timingList");
    list.innerHTML = "";

    data.forEach(c => {

      const thumb = getCartoonImage(c);

      const row = document.createElement("div");
      row.className = "timing-row";

      row.innerHTML = `
        <img src="${thumb}">
        <div class="info">${c.Name}</div>
        <div class="time">${c.ShowTime}</div>
      `;

      list.appendChild(row);
    });

  } catch (err) {
    console.error("Timing load error:", err);
  }
}

/* ==== INIT ==== */
window.onload = loadTimings;
