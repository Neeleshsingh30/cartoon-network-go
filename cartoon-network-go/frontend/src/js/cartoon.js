const BASE_URL = "http://localhost:8000";
const id = new URLSearchParams(window.location.search).get("id");

let isLiked = false;

/* ---------- AUTH HEADER ---------- */
function authHeader(){
  const token = localStorage.getItem("token");
  return token ? { "Authorization": `Bearer ${token}` } : {};
}

/* ---------- STAR RENDER ---------- */
function renderStars(imdb){
  const stars = Math.round(imdb / 2);
  return "★".repeat(stars) + "☆".repeat(5 - stars);
}

/* ---------- LOAD CARTOON ---------- */
async function loadCartoon(){
  const res = await fetch(`${BASE_URL}/cartoon/${id}`);
  const data = await res.json();

 const banner = data.Images?.find(i => i.ImageType === "banner")?.ImageURL;
if (banner) {
  document.getElementById("bannerImg").src = banner;
}


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
  

  await checkLiked();
  loadRecommendations();
  addView();
}

async function addView(){
  const token = localStorage.getItem("token");
  if(!token) return;

  try{
    await fetch(`${BASE_URL}/cartoon/${id}/view`,{
      method: "POST",
      headers:{
        "Authorization": `Bearer ${token}`
      }
    });
  }catch(err){
    console.warn("View not counted");
  }
}

/* ---------- LIKE STATUS CHECK ---------- */
async function checkLiked(){
  try{
    const res = await fetch(`${BASE_URL}/user/favourites`,{
      headers: authHeader()
    });

    if(!res.ok) return;

    const favs = await res.json();
    isLiked = favs.some(f => String(f.Cartoon.ID) === String(id));
    updateHeart();
  }catch(err){
    console.warn("Like check skipped");
  }
}

/* ---------- HEART UI ---------- */
function updateHeart(){
  const btn = document.getElementById("favBtn");
  if(!btn) return;
  btn.innerHTML = isLiked ? "❤️" : "🤍";
  btn.classList.toggle("liked", isLiked);
}

/* ---------- LIKE / UNLIKE ---------- */
async function toggleFavourite(){

  if(!localStorage.getItem("token")){
    alert("Please login again");
    window.location.href = "index.html";
    return;
  }

  const method = isLiked ? "DELETE" : "POST";

  const res = await fetch(`${BASE_URL}/cartoon/${id}/like`,{
    method,
    headers:{
      ...authHeader(),
      "Content-Type": "application/json"
    }
  });

  if(!res.ok){
    console.error("LIKE FAILED", await res.text());
    return;
  }

  isLiked = !isLiked;
  updateHeart();
}

/* ---------- RECOMMENDATIONS ---------- */
async function loadRecommendations(){
  try{
    const res = await fetch(`${BASE_URL}/cartoon/${id}/recommendations`);
    const data = await res.json();

    const grid = document.getElementById("recommendGrid");
    grid.innerHTML = "";

    data.forEach(c => {
      let thumb = "../../images/CN-BG-AUTH.jpg";

      if(c.Images && c.Images.length){
        const img = c.Images.find(i => i.ImageType === "thumbnail") 
               || c.Images.find(i => i.ImageType === "poster");
        if(img) thumb = img.ImageURL;
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
  }catch(err){
    console.error("Recommendation Error:", err);
  }
}

window.onload = loadCartoon;
window.toggleFavourite = toggleFavourite;
