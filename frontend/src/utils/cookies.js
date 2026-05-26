import Cookies from 'js-cookie';

const COOKIE_NAME = 'auth_token';
const COOKIE_OPTIONS = {
  secure: true,
  sameSite: 'Strict',
  expires: 1,
};
const DEFAULT_SECRET = 'change_this_to_a_stronger_secret';

const encoder = new TextEncoder();
const decoder = new TextDecoder();

const getSecret = () => import.meta.env.VITE_TOKEN_SECRET || DEFAULT_SECRET;

const arrayBufferToBase64 = (buffer) => {
  let binary = '';
  const bytes = new Uint8Array(buffer);
  for (let i = 0; i < bytes.byteLength; i += 1) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary);
};

const base64ToArrayBuffer = (base64) => {
  const binary = atob(base64);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i += 1) {
    bytes[i] = binary.charCodeAt(i);
  }
  return bytes.buffer;
};

const getCryptoKey = async () => {
  const secret = getSecret();
  const salt = encoder.encode('securetask-token-salt');
  const baseKey = await window.crypto.subtle.importKey(
    'raw',
    encoder.encode(secret),
    'PBKDF2',
    false,
    ['deriveKey']
  );

  return window.crypto.subtle.deriveKey(
    {
      name: 'PBKDF2',
      salt,
      iterations: 250000,
      hash: 'SHA-256',
    },
    baseKey,
    { name: 'AES-GCM', length: 256 },
    false,
    ['encrypt', 'decrypt']
  );
};

const encryptData = async (value) => {
  const key = await getCryptoKey();
  const iv = window.crypto.getRandomValues(new Uint8Array(12));
  const encrypted = await window.crypto.subtle.encrypt(
    { name: 'AES-GCM', iv },
    key,
    encoder.encode(value)
  );

  return `${arrayBufferToBase64(iv)}.${arrayBufferToBase64(encrypted)}`;
};

const decryptData = async (value) => {
  const [ivBase64, dataBase64] = value.split('.');
  if (!ivBase64 || !dataBase64) return null;

  const key = await getCryptoKey();
  const iv = new Uint8Array(base64ToArrayBuffer(ivBase64));
  const encrypted = base64ToArrayBuffer(dataBase64);

  try {
    const decrypted = await window.crypto.subtle.decrypt(
      { name: 'AES-GCM', iv },
      key,
      encrypted
    );
    return decoder.decode(decrypted);
  } catch (error) {
    console.error('Token decryption failed:', error);
    return null;
  }
};

export const setToken = async (token) => {
  const encryptedToken = await encryptData(token);
  Cookies.set(COOKIE_NAME, encryptedToken, COOKIE_OPTIONS);
};

export const getToken = async () => {
  const encryptedToken = Cookies.get(COOKIE_NAME);
  if (!encryptedToken) return null;

  return decryptData(encryptedToken);
};

export const removeToken = () => {
  Cookies.remove(COOKIE_NAME, { secure: true, sameSite: 'Strict' });
};
