export const saveSettings = (settings) => {
  localStorage.setItem('appSettings', JSON.stringify(settings));
};
