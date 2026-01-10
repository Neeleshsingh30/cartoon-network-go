const BASE_URL = "http://localhost:8000";

const id = new URLSearchParams(window.location.search).get("id");

async function loadCartoon(){
  const res = await fetch(`${BASE_URL}/cartoon/${id}`);
  const data = await res.json();

  const banner = data.Images.find(i => i.ImageType === "banner")?.ImageURL;

  document.getElementById("banner").style.backgroundImage = `url(${banner})`;

  document.getElementById("cartoonInfo").innerHTML = `
    <div class="info-box">Description: ${data.Description}</div>
    <div class="info-box">Genre: ${data.Genre}</div>
    <div class="info-box">Age Group: ${data.AgeGroup}</div>
    <div class="info-box">Universe: ${data.Universe}</div>
  `;

  const table = document.getElementById("charTable");

  data.Characters.forEach(ch => {
    const row = document.createElement("tr");
    row.innerHTML = `
      <td>${ch.Name}</td>
      <td>${ch.Role}</td>
      <td>${ch.Power}</td>
      <td>${ch.Description}</td>
    `;
    table.appendChild(row);
  });
}

window.onload = loadCartoon;
