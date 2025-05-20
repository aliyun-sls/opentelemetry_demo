import type { NextApiRequest, NextApiResponse } from 'next';
import { clearSession } from '../../../utils/session';

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: '方法不允许' });
  }

  try {
    const session = req.cookies.session;
    if (session) {
      // 调用用户服务的登出接口
      await fetch(`${process.env.USER_SERVICE_URL}/logout`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Cookie': `sessionid=${session}`,
        },
      });
    }

    // 清除会话 cookie
    clearSession(res);

    return res.status(200).json({ message: '登出成功' });
  } catch (error) {
    console.error('Logout error:', error);
    return res.status(500).json({ error: '登出过程中发生错误' });
  }
} 