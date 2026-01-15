const BASE_URL = "https://cartoon-network-go-1.onrender.com";

/* ---------- AUTH HEADER ---------- */
function authHeader(){
  const token = localStorage.getItem("token");
  return token ? { "Authorization": `Bearer ${token}` } : null;
}

/* ---------- REQUIRE LOGIN ---------- */
function requireAuth(){
  const token = localStorage.getItem("token");
  if(!token){
    window.location.href = "index.html";
    return false;
  }
  return true;
}

/* ---------- LOAD HOME CARTOONS ---------- */
async function loadHomeCartoons() {
  try{
    const res = await fetch(`${BASE_URL}/cartoons`);
    const data = await res.json();

    const grid = document.getElementById("cartoonGrid");
    grid.innerHTML = "";

   data.forEach(c => {

  console.log("FULL CARTOON:", c);
  console.log("IMAGES ARRAY:", c.Images);

  let thumb = "/src/images/CN-BG-AUTH.jpg";

  if (c.Images && c.Images.length) {
    console.log("FIRST IMAGE OBJ:", c.Images[0]);
  }

      if (c.Images && c.Images.length) {
  const img =
    c.Images.find(i => i.image_type === "thumbnail") ||
    c.Images.find(i => i.image_type === "poster");

  if (img && img.image_url) thumb = img.image_url;
}


      const card = document.createElement("div");
      card.className = "cartoon-card";
      card.onclick = () => window.location.href = `cartoon.html?id=${c.ID}`;

      card.innerHTML = `
        <img src="${thumb}">
        <p>${c.Name}</p>
      `;

      grid.appendChild(card);
    });

  }catch(err){
    console.error("Home load error:", err);
  }
}

/* ---------- LOAD USER FAV COUNT ---------- */
async function loadFavCount(){
  if(!requireAuth()) return;

  const headers = authHeader();
  if(!headers) return;

  try{
    const res = await fetch(`${BASE_URL}/user/favourites`,{
      headers: headers
    });

    if(res.status === 401){
      localStorage.removeItem("token");
      window.location.href = "index.html";
      return;
    }

    const favs = await res.json();
    console.log("Favourite cartoons:", favs.length);

  }catch(err){
    console.warn("Fav count skipped");
  }
}

async function loadTrendingCartoons(){
  try{
    const res = await fetch(`${BASE_URL}/cartoons/trending`);
    const data = await res.json();

    const grid = document.getElementById("trendingGrid");
    grid.innerHTML = "";

     data.forEach(c => {

  console.log("FULL CARTOON:", c);
  console.log("IMAGES ARRAY:", c.Images);

  let thumb = "/src/images/CN-BG-AUTH.jpg";

  if (c.Images && c.Images.length) {
    console.log("FIRST IMAGE OBJ:", c.Images[0]);
  }

      if (c.Images && c.Images.length) {
  const img =
    c.Images.find(i => i.image_type === "thumbnail") ||
    c.Images.find(i => i.image_type === "poster");

  if (img && img.image_url) thumb = img.image_url;
}

      const card = document.createElement("div");
      card.className = "cartoon-card";
      card.onclick = () => window.location.href = `cartoon.html?id=${c.ID}`;

      card.innerHTML = `
        <img src="${thumb}">
        <p>${c.Name}</p>
        <span style="color:gold;">🔥 ${c.view_count} views</span>
      `;

      grid.appendChild(card);
    });

  }catch(err){
    console.error("Trending load error:", err);
  }
}

/* ---------- LOGOUT ---------- */
function logout(){
  localStorage.removeItem("token");
  window.location.href = "index.html";
}

/* ---------- INIT ---------- */
window.onload = function(){
  // wait for token to settle after redirect
  setTimeout(() => {

    const token = localStorage.getItem("token");
    if(!token){
      window.location.href = "index.html";
      return;
    }
    loadTrendingCartoons();
    loadHomeCartoons();
    loadFavCount();

  }, 150);   // 150ms delay = race condition solved
};
function doSearch(){
  const q = document.getElementById("searchInput").value.trim();
  if(!q) return;
  window.location.href = `search.html?query=${q}`;
}