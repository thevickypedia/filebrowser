import { useAuthStore } from "@/stores/auth";
import router from "@/router";
import { jwtDecode, JwtPayload } from "jwt-decode";
import { baseURL, noAuth } from "./constants";
import { StatusError } from "@/api/utils";

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

export async function ConvertStringToHex(str: string) {
  const arr = [];
  for (let i = 0; i < str.length; i++) {
    arr[i] = ("00" + str.charCodeAt(i).toString(16)).slice(-4);
  }
  return "\\u" + arr.join("\\u");
}

export async function login(
  username: string,
  password: string,
  recaptcha: string
) {
  const hex_user = await ConvertStringToHex(username);
  const hex_pass = await ConvertStringToHex(password);
  const hex_recaptcha = await ConvertStringToHex(recaptcha);
  let payload = btoa(hex_user + "," + hex_pass + "," + hex_recaptcha) // eslint-disable-line
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
