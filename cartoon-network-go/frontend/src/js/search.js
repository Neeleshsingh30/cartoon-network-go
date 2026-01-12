const BASE_URL = "http://localhost:8000";
const query = new URLSearchParams(window.location.search).get("query");

async function loadResults(){
  const res = await fetch(`${BASE_URL}/cartoons/search?query=${query}`);
  const data = await res.json();

  const grid = document.getElementById("resultGrid");
  grid.innerHTML = "";

  if(data.length === 0){
    grid.innerHTML = "<p>No cartoons found 😢</p>";
    return;
  }

  data.forEach(c => {
    let thumb = "../../images/CN-BG-AUTH.jpg";

    if(c.Images && c.Images.length){
      const img = c.Images.find(i => i.ImageType === "thumbnail") 
             || c.Images.find(i => i.ImageType === "poster");
      if(img) thumb = img.ImageURL;
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
}

window.onload = loadResults;
