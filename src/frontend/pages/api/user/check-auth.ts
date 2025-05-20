import type { NextApiRequest, NextApiResponse } from 'next';
import { getSession } from '../../../utils/session';

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method !== 'GET') {
    return res.status(405).json({ error: '方法不允许' });
  }

  try {
    const session = getSession(req);
    if (!session) {
      return res.status(401).json({ error: '未认证' });
    }

    // 验证会话是否有效
    const response = await fetch(`${process.env.USER_SERVICE_URL}/check-auth`, {
      method: 'GET',
      headers: {
        'Cookie': `sessionid=${session.sessionId}`,
      },
    });

    if (!response.ok) {
      return res.status(401).json({ error: '会话已过期' });
    }

    return res.status(200).json({
      username: session.username,
      role: session.role,
    });
  } catch (error) {
    console.error('Auth check error:', error);
    return res.status(500).json({ error: '验证过程中发生错误' });
  }
} 