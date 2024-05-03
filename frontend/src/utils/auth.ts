import { useAuthStore } from "@/stores/auth";
import router from "@/router";
import { JwtPayload, jwtDecode } from "jwt-decode";
import { baseURL, noAuth } from "./constants";
import { StatusError } from "@/api/utils";
import CryptoJS from "crypto-js";

export function parseToken(token: string) {
  // falsy or malformed jwt will throw InvalidTokenError
  const data = jwtDecode<JwtPayload & { user: IUser }>(token);

  // Setting cookie options
  const cookieOptions = {
    path: "/",
    secure: true, // Ensures the cookie is sent only over HTTPS
    httpOnly: true, // Prevents JavaScript from accessing the cookie
  };

  // Constructing the cookie string
  let cookieString = `auth=${token};`;

  // Adding additional options
  for (const [key, value] of Object.entries(cookieOptions)) {
    cookieString += ` ${key}=${value};`;
  }

  // Setting the cookie
  document.cookie = cookieString;

  // Setting token in localStorage and Vuex store
  localStorage.setItem("jwt", token);

  const authStore = useAuthStore();
  authStore.jwt = token;
  authStore.setUser(data.user);
}

export async function validateLogin() {
  try {
    if (localStorage.getItem("jwt")) {
      await renew(<string>localStorage.getItem("jwt"));
    }
  } catch (error) {
    console.warn("Invalid JWT token in storage"); // eslint-disable-line
    throw error;
  }
}

export async function getProxyFlag() {
  const name = "pyproxy";
  const cookies = document.cookie.split(";");
  for (let i = 0; i < cookies.length; i++) {
    const cookie = cookies[i].trim();
    if (cookie.indexOf(name + "=") === 0) {
      return cookie.substring(name.length + 1, cookie.length);
    }
  }
  return "off";
}

export async function ConvertStringToHex(str: string) {
  const arr = [];
  for (let i = 0; i < str.length; i++) {
    arr[i] = ("00" + str.charCodeAt(i).toString(16)).slice(-4);
  }
  return "\\u" + arr.join("\\u");
}

async function CalculateHash(message: string) {
  const encoder = new TextEncoder();
  const data = encoder.encode(message);
  if (crypto.subtle === undefined) {
    const wordArray = CryptoJS.lib.WordArray.create(data);
    const hash = CryptoJS.SHA512(wordArray);
    // Convert the hash to a hexadecimal string and return it
    return hash.toString(CryptoJS.enc.Hex);
  } else {
    const hashBuffer = await crypto.subtle.digest("SHA-512", data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    // Convert each byte to a hexadecimal string, pad with zeros, and join them to form the final hash
    return hashArray.map((byte) => byte.toString(16).padStart(2, "0")).join("");
  }
}

export async function login(
  username: string,
  password: string,
  recaptcha: string
) {
  const proxy_status = await getProxyFlag();
  let payload;
  if (proxy_status === "on") {
    const hex_user = await ConvertStringToHex(username);
    const signature = await CalculateHash(password);
    const hex_recaptcha = await ConvertStringToHex(recaptcha);
    payload = btoa(hex_user + "," + signature + "," + hex_recaptcha) // eslint-disable-line
  } else {
    console.warn(
      "pyproxy is turned off! auth header will be sent as plain text"
    );
    payload = JSON.stringify({ username, password, recaptcha });
  }

  const res = await fetch(`${baseURL}/api/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: payload,
    },
  });

  const body = await res.text();

  if (res.status === 200) {
    parseToken(body);
  } else {
    throw new StatusError(
      body || `${res.status} ${res.statusText}`,
      res.status
    );
  }
}

export async function renew(jwt: string) {
  const res = await fetch(`${baseURL}/api/renew`, {
    method: "POST",
    headers: {
      "X-Auth": jwt,
    },
  });

  const body = await res.text();

  if (res.status === 200) {
    parseToken(body);
  } else {
    throw new StatusError(
      body || `${res.status} ${res.statusText}`,
      res.status
    );
  }
}

export async function signup(username: string, password: string) {
  const data = { username, password };

  const res = await fetch(`${baseURL}/api/signup`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  if (res.status !== 200) {
    throw new StatusError(`${res.status} ${res.statusText}`, res.status);
  }
}

export function logout() {
  document.cookie = "auth=; Max-Age=0; Path=/; SameSite=Strict;";

  const authStore = useAuthStore();
  authStore.clearUser();

  localStorage.setItem("jwt", "");
  if (noAuth) {
    window.location.reload();
  } else {
    router.push({ path: "/login" });
  }
}
