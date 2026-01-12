const BASE_URL = "http://localhost:8000";

/* ---------- AUTH HEADER ---------- */
function authHeader(){
  const token = localStorage.getItem("token");
  return token ? { "Authorization": `Bearer ${token}` } : null;
}

/* ---------- LOAD FAVOURITES ---------- */
async function loadFavourites(){

  const headers = authHeader();
  if(!headers){
    window.location.href = "index.html";
    return;
  }

  try{
    const res = await fetch(`${BASE_URL}/user/favourites`,{
      headers: headers
    });

    if(res.status === 401){
      localStorage.removeItem("token");
      window.location.href = "index.html";
      return;
    }

    const data = await res.json();
    console.log("FAV DATA =>", data);

    const grid = document.getElementById("favGrid");
    grid.innerHTML = "";

    if(!Array.isArray(data) || data.length === 0){
      grid.innerHTML = "<p>No favourite cartoons yet ❤️</p>";
      return;
    }

    data.forEach(f => {
      if(!f.Cartoon) return;
      const c = f.Cartoon;

      let thumb = "../../images/CN-BG-AUTH.jpg";
      if(c.Images && c.Images.length){
        const img = c.Images.find(i => i.ImageType === "thumbnail") 
                 || c.Images.find(i => i.ImageType === "poster");
        if(img) thumb = img.ImageURL;
      }

      const card = document.createElement("div");
      card.className = "fav-card";
      card.onclick = () => window.location.href = `cartoon.html?id=${c.ID}`;

      card.innerHTML = `
        <img src="${thumb}">
        <p>${c.Name}</p>
      `;

      grid.appendChild(card);
    });

  }catch(err){
    console.error("Fav Load Error:", err);
  }
}

/* ---------- INIT (DELAYED) ---------- */
window.onload = function(){
  setTimeout(loadFavourites, 100);
};
