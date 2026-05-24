import Cookies from 'js-cookie';

const COOKIE_NAME = 'auth_token';

const COOKIE_OPTIONS = {
  secure: true,       // only sent over HTTPS
  sameSite: 'Strict', // blocks cross-site request forgery
  expires: 1,         // 1 day
};

export const setToken = (token) => {
  Cookies.set(COOKIE_NAME, token, COOKIE_OPTIONS);
};

export const getToken = () => {
  return Cookies.get(COOKIE_NAME) ?? null;
};

export const removeToken = () => {
  Cookies.remove(COOKIE_NAME, { secure: true, sameSite: 'Strict' });
};
