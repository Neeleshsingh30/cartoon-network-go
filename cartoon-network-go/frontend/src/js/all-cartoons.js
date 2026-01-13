const BASE_URL = "http://localhost:8000";

async function loadAllCartoons(){
  const res = await fetch(`${BASE_URL}/cartoons`);
  const cartoons = await res.json();

  const grid = document.getElementById("allGrid");

 cartoons.forEach(c => {
  let img = c.Images.find(i => i.ImageType === "banner")?.ImageURL 
            || "../../images/CN-BG-AUTH.jpg";

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

}

window.onload = loadAllCartoons;
