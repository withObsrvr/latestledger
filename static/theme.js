function themeIcon(theme) {
  return theme === "dark" ? "☀" : "☾";
}

function setTheme(theme) {
  document.documentElement.dataset.theme = theme;
  localStorage.setItem("theme", theme);
  document.querySelectorAll("[data-theme-toggle]").forEach((button) => {
    button.textContent = themeIcon(theme);
    button.setAttribute("aria-label", theme === "dark" ? "Switch to light mode" : "Switch to dark mode");
    button.setAttribute("aria-pressed", theme === "dark" ? "true" : "false");
  });
}

function initThemeToggle() {
  const current = document.documentElement.dataset.theme || "light";
  document.querySelectorAll("[data-theme-toggle]").forEach((button) => {
    button.textContent = themeIcon(current);
    button.setAttribute("aria-label", current === "dark" ? "Switch to light mode" : "Switch to dark mode");
    button.setAttribute("aria-pressed", current === "dark" ? "true" : "false");
    button.onclick = () => setTheme((document.documentElement.dataset.theme || "light") === "dark" ? "light" : "dark");
  });
}

document.addEventListener("DOMContentLoaded", initThemeToggle);
document.body.addEventListener("htmx:afterSwap", initThemeToggle);
