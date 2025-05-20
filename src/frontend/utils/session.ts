import { NextApiRequest, NextApiResponse } from 'next';
import { serialize } from 'cookie';

interface Session {
  sessionId: string;
  username: string;
  role: string;
}

export function getSession(req: NextApiRequest): Session | null {
  const sessionCookie = req.cookies.session;
  if (!sessionCookie) {
    return null;
  }

  try {
    return JSON.parse(sessionCookie);
  } catch {
    return null;
  }
}

export function setSession(
  req: NextApiRequest,
  res: NextApiResponse,
  session: Session
) {
  const cookie = serialize('session', JSON.stringify(session), {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    maxAge: 24 * 60 * 60, // 24 小时
    path: '/',
  });

  res.setHeader('Set-Cookie', cookie);
}

export function clearSession(res: NextApiResponse) {
  const cookie = serialize('session', '', {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    maxAge: -1,
    path: '/',
  });

  res.setHeader('Set-Cookie', cookie);
} 