const shortBtn = document.querySelector(".btn");
const inputElement = document.querySelector(".url-input");
const statsElement = document.querySelector(".stats-url");
const outputElement = document.querySelector(".short-url");

const apiUrl = "http://localhost:8080";
var currentUrl = inputElement.value;

async function shortUrl(url) {
  const response = await fetch(`${apiUrl}/api/short`, {
    method: "POST",
    body: JSON.stringify({ url: url }),
  });

  if (!response.ok) {
    outputElement.textContent = "Произошла ошибка при запросе";
    outputElement.href = "#";
  }

  body = await response.json();

  slug = `${body["text"]}_${body["number"]}`;

  slugUrl = `${apiUrl}/s/${slug}`;
  statsUrl = `${apiUrl}/statistics.html?slug=${slug}`;

  currentUrl = slugUrl;
  outputElement.textContent = slugUrl;
  outputElement.href = slugUrl;

  statsElement.textContent = statsUrl;
  statsElement.href = statsUrl;
}

function onClickShortButton() {
  if (currentUrl != inputElement.value) shortUrl(inputElement.value);
}

shortUrl(currentUrl);

shortBtn.addEventListener("click", onClickShortButton, false);
