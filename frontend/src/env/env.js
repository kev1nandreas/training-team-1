const required = (key) => {
  const value = import.meta.env[key];
  if (!value) throw new Error(`Missing required env var: ${key}`);
  return value;
};

const env = {
  API_BASE_URL: required('VITE_API_BASE_URL'),
  ADMIN_API_KEY: required('VITE_ADMIN_API_KEY'),
  AWS_ACCESS_KEY: required('VITE_AWS_ACCESS_KEY'),
  AWS_SECRET_KEY: required('VITE_AWS_SECRET_KEY'),
  AWS_BUCKET: required('VITE_AWS_BUCKET'),
  APP_CONFIG: {
    name: 'SecureTask',
    version: '1.0.0',
  },
};

export default env;
