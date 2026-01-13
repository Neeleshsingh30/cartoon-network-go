const BASE_URL = "http://localhost:8000";

async function loadTimings(){
  const res = await fetch(`${BASE_URL}/cartoons/timings`);
  const data = await res.json();

  const list = document.getElementById("timingList");
  list.innerHTML = "";

  data.forEach(c => {

    let thumb = "../../images/CN-BG-AUTH.jpg";
    if (c.Images && c.Images.length > 0) {
      const img = c.Images.find(i => i.ImageType === "thumbnail");
      if (img) thumb = img.ImageURL;
    }

    const row = document.createElement("div");
    row.className = "timing-row";

    row.innerHTML = `
      <img src="${thumb}">
      <div class="info">${c.Name}</div>
      <div class="time">${c.ShowTime}</div>
    `;

    list.appendChild(row);
  });
}

window.onload = loadTimings;
