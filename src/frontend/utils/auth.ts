import SessionGateway from '../gateways/Session.gateway';

export interface AuthResult {
  isAuthenticated: boolean;
  isAdmin?: boolean;
}

export function checkAuth(): boolean {
  if (typeof window !== 'undefined') {
      const localSid = localStorage.getItem('sid');
      const session = SessionGateway.getSession();
      const sessionSid = session.sid;
      return !!localSid && !!sessionSid && localSid === sessionSid;
    }
  return false;
}


export function isAdmin(): boolean {
  const userRole = localStorage.getItem('userRole');
  return userRole === '0'; // ROLE_ADMIN = 0
}


export function getAuthStatus(): AuthResult {
  const sessionId = localStorage.getItem('sid');

  if (!sessionId) {
    return { isAuthenticated: false };
  }

  const userRole = localStorage.getItem('uid');
  const isAdminUser = userRole === '0'; // ROLE_ADMIN = 0

  return {
    isAuthenticated: true,
    isAdmin: isAdminUser
  };
}


export async function logout(): Promise<boolean> {
  try {
    const response = await fetch('/user/api/logout', {
      method: 'POST',
      credentials: 'include', // 包含 cookie
    });

    if (!response.ok) {
      throw new Error('注销失败');
    }

    SessionGateway.removeSessionValue('sid');
    localStorage.clear();

    return true;
  } catch (error) {
    console.error('注销错误:', error);
    localStorage.clear();
    throw error;
  }
}