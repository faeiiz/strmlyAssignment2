let token = "";

// Handle Signup
async function signup() {
  const name = document.querySelector(
    '.flip-card__back input[placeholder="Name"]'
  ).value;
  const email = document.querySelector(
    '.flip-card__back input[name="email"]'
  ).value;
  const password = document.querySelector(
    '.flip-card__back input[name="password"]'
  ).value;

  const res = await fetch("http://localhost:8080/signup", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, email, password }),
  });

  const data = await res.json();
  alert("Signup: " + JSON.stringify(data));
}

// Login
async function login() {
  const email = document.querySelector(
    '.flip-card__front input[name="email"]'
  ).value;
  const password = document.querySelector(
    '.flip-card__front input[name="password"]'
  ).value;

  console.log(email, password);

  const res = await fetch("http://localhost:8080/login", {
    method: "POST",
      headers: { "Content-Type": "application/json" },
     credentials: 'include', 
    body: JSON.stringify({ email, password }),
    credentials: "include", //sends httpOnly cookies to back
  });

  const data = await res.json();
  if (res.ok) {
    window.location.href = "home.html";
  } else {
    alert("Login failed");
  }
}
document
  .querySelector(".flip-card__front button")
  .addEventListener("click", (e) => {
    e.preventDefault();
    login();
  });
document
  .querySelector(".flip-card__back button")
  .addEventListener("click", (e) => {
    e.preventDefault();
    signup();
  });
