const BASE_URL = "http://localhost:8000";
const id = new URLSearchParams(window.location.search).get("id");

function renderStars(imdb){
  const stars = Math.round(imdb / 2);
  return "★".repeat(stars) + "☆".repeat(5 - stars);
}

async function loadCartoon(){
  const res = await fetch(`${BASE_URL}/cartoon/${id}`);
  const data = await res.json();

  const banner = data.Images.find(i => i.ImageType === "banner")?.ImageURL;
  document.getElementById("banner").style.backgroundImage = `url(${banner})`;

  document.getElementById("cartoonInfo").innerHTML = `
    <div class="info-box">Description: ${data.Description}</div>
    <div class="info-box">Genre: ${data.Genre}</div>
    <div class="info-box">Universe: ${data.Universe}</div>
    <div class="info-box">Age Group: ${data.AgeGroup}</div>
  `;

  document.getElementById("ratingBox").innerHTML = `
    <div>IMDb: ${data.ImdbRating}/10 
      <span class="rating-stars">${renderStars(data.ImdbRating)}</span>
    </div>
    <div>📅 Aired On: ${new Date(data.AirDate).toDateString()}</div>
    <div>⏰ Show Time: ${data.ShowTime}</div>
  `;

  const table = document.getElementById("charTable");
  table.innerHTML = "";
  data.Characters.forEach(ch => {
    table.innerHTML += `
      <tr>
        <td>${ch.Name}</td>
        <td>${ch.Role}</td>
        <td>${ch.Power}</td>
        <td>${ch.Description}</td>
      </tr>
    `;
  });

  loadRecommendations();
}

async function loadRecommendations(){
  try {
    const res = await fetch(`${BASE_URL}/cartoon/${id}/recommendations`);
    const data = await res.json();

    console.log("RECOMMEND DATA =>", data);

    const grid = document.getElementById("recommendGrid");
    grid.innerHTML = "";

    if (!data || data.length === 0) {
      grid.innerHTML = "<p>No recommendations available.</p>";
      return;
    }

    data.forEach(c => {

      let thumb = "../../images/CN-BG-AUTH.jpg";

      if (c.Images && c.Images.length > 0) {
        const imgObj = c.Images.find(i => i.ImageType === "thumbnail") 
                    || c.Images.find(i => i.ImageType === "poster");

        if (imgObj) thumb = imgObj.ImageURL;
      }

      const card = document.createElement("div");
      card.className = "recommend-card";
      card.onclick = () => window.location.href = `cartoon.html?id=${c.ID}`;

      card.innerHTML = `
        <img src="${thumb}">
        <p>${c.Name}</p>
      `;

      grid.appendChild(card);
    });

  } catch (err) {
    console.error("Recommendation Error:", err);
  }
}

window.onload = loadCartoon;

