export function checkAuth() {
  const cookies = document.cookie;
  return cookies.includes('isLoggedIn=true');
}