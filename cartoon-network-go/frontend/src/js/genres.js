const BASE_URL = "http://localhost:8000";

async function loadGenres(){
  try{
    const res = await fetch(`${BASE_URL}/cartoons/by-genre`);
    const data = await res.json();

    const container = document.getElementById("genreContainer");
    container.innerHTML = "";

    // 🔥 IMPORTANT FIX
    Object.keys(data).forEach(genre => {

      const section = document.createElement("div");
      section.className = "genre-section";

      section.innerHTML = `
        <h2 class="genre-title">${genre}</h2>
        <div class="genre-row" id="row-${genre}"></div>
      `;

      container.appendChild(section);

      const row = document.getElementById(`row-${genre}`);

      data[genre].forEach(c => {

        let thumb = "../../images/CN-BG-AUTH.jpg";
        if(c.Images && c.Images.length){
          const img = c.Images.find(i => i.ImageType === "thumbnail")
                 || c.Images.find(i => i.ImageType === "poster");
          if(img) thumb = img.ImageURL;
        }

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

  }catch(err){
    console.error("Genre load error:", err);
  }
}

window.onload = loadGenres;
