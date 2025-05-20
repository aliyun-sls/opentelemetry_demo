import type { NextApiRequest, NextApiResponse } from 'next';
import { getSession, setSession } from '../../../utils/session';

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: '方法不允许' });
  }

  try {
    const { username, password } = req.body;

    const response = await fetch(`${process.env.USER_SERVICE_URL}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    const data = await response.json();

    if (!response.ok) {
      return res.status(response.status).json({ error: data.error });
    }

    // 设置会话
    setSession(req, res, {
      sessionId: data.sessionid,
      username: username,
      role: data.role,
    });

    return res.status(200).json({
      username: username,
      role: data.role,
    });
  } catch (error) {
    console.error('Login error:', error);
    return res.status(500).json({ error: '登录过程中发生错误' });
  }
} 