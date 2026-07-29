const TOKEN_KEY = "echo-watch-token";

const isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent);
const isStandalone =
  window.matchMedia("(display-mode: standalone)").matches || navigator.standalone === true;

const iosGuide = document.getElementById("ios-guide");
const pairCard = document.getElementById("pair-card");
const tokenInput = document.getElementById("token");
const subscribeBtn = document.getElementById("subscribe-btn");
const testBtn = document.getElementById("test-btn");
const statusEl = document.getElementById("status");

function setStatus(message, isError = false) {
  statusEl.textContent = message;
  statusEl.style.color = isError ? "#ff8f8f" : "";
}

function urlBase64ToUint8Array(base64String) {
  const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding).replace(/-/g, "+").replace(/_/g, "/");
  const raw = atob(base64);
  const output = new Uint8Array(raw.length);
  for (let i = 0; i < raw.length; i += 1) {
    output[i] = raw.charCodeAt(i);
  }
  return output;
}

async function fetchPublicKey() {
  const res = await fetch("/api/vapid-public-key");
  if (!res.ok) {
    throw new Error("VAPID public key unavailable");
  }
  const data = await res.json();
  return data.publicKey;
}

async function registerServiceWorker() {
  if (!("serviceWorker" in navigator)) {
    throw new Error("Service workers are not supported in this browser.");
  }
  return navigator.serviceWorker.register("/sw.js");
}

async function subscribe() {
  const token = tokenInput.value.trim();
  if (!token) {
    setStatus("Enter your pairing token first.", true);
    return;
  }

  setStatus("Requesting notification permission…");
  const permission = await Notification.requestPermission();
  if (permission !== "granted") {
    setStatus("Notification permission denied.", true);
    return;
  }

  const registration = await registerServiceWorker();
  await navigator.serviceWorker.ready;

  const publicKey = await fetchPublicKey();
  const subscription = await registration.pushManager.subscribe({
    userVisibleOnly: true,
    applicationServerKey: urlBase64ToUint8Array(publicKey),
  });

  const res = await fetch("/api/subscribe", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ token, subscription: subscription.toJSON() }),
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || "Subscribe failed");
  }

  localStorage.setItem(TOKEN_KEY, token);
  testBtn.classList.remove("hidden");
  setStatus("Subscribed. You will receive alerts from echo-watch.");
}

async function sendTestPush() {
  const token = tokenInput.value.trim() || localStorage.getItem(TOKEN_KEY) || "";
  if (!token) {
    setStatus("Enter your pairing token first.", true);
    return;
  }

  const res = await fetch("/api/test-push", {
    method: "POST",
    headers: { Authorization: `Bearer ${token}` },
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || "Test push failed");
  }
  setStatus("Test push sent.");
}

function init() {
  const saved = localStorage.getItem(TOKEN_KEY);
  if (saved) {
    tokenInput.value = saved;
    testBtn.classList.remove("hidden");
  }

  if (isIOS && !isStandalone) {
    iosGuide.classList.remove("hidden");
    subscribeBtn.disabled = true;
    setStatus("Install to Home Screen first, then open this app from the icon.");
    return;
  }

  subscribeBtn.addEventListener("click", () => {
    subscribe().catch((err) => setStatus(err.message, true));
  });

  testBtn.addEventListener("click", () => {
    sendTestPush().catch((err) => setStatus(err.message, true));
  });
}

init();
