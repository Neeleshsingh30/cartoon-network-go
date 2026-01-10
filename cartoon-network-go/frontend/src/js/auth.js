const BASE_URL = "http://localhost:8000";

let isLogin = true;

function toggle(){
  isLogin = !isLogin;

  document.getElementById("title").innerText = isLogin ? "Login" : "Signup";
  document.querySelector(".auth-box button").innerText = isLogin ? "Login" : "Signup";
  document.getElementById("confirm").style.display = isLogin ? "none" : "block";
  document.getElementById("email").style.display = isLogin ? "none" : "block";

}

async function submitForm(){

  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value.trim();
  const confirm  = document.getElementById("confirm").value.trim();
  const email = document.getElementById("email").value.trim();



  if(!isLogin && (!username || !email || !password || !confirm)){
  alert("All fields are required");
  return;
}

if(isLogin && (!username || !password)){
  alert("Username and Password required");
  return;
}

 

  if(!isLogin && password !== confirm){
    alert("Passwords do not match");
    return;
  }

  const url = isLogin ? "/login" : "/signup";

const payload = isLogin
  ? { username, password }
  : {
      username,
      email,
      password,
      confirm_password: confirm
    };


  const res = await fetch(BASE_URL + url, {
    method: "POST",
    headers: {"Content-Type":"application/json"},
    body: JSON.stringify(payload)
  });

  const data = await res.json();

  if(!res.ok){
    alert(data.detail || "Invalid request");
    return;
  }

  localStorage.setItem("token", data.access_token);
  window.location.href = "home.html";
}
