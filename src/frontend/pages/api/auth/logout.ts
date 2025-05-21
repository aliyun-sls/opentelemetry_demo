// src/frontend/pages/api/auth/logout.ts  
import type { NextApiRequest, NextApiResponse } from 'next';  
  
type LogoutResponse = {  
  message: string;  
  error?: string;  
}  
  
export default async function handler(  
  req: NextApiRequest,   
  res: NextApiResponse<LogoutResponse>  
) {  
  if (req.method !== 'POST') {  
    return res.status(405).json({ error: '方法不允许', message: '只允许 POST 请求' });  
  }  
  
  try {  
    const response = await fetch('/user/api/logout', {
      method: 'POST',
      credentials: 'include',
    });

    const data: LogoutResponse = await response.json();

    if (!response.ok) {
      console.error('[Logout] 后端注销失败:', data.error);
      throw new Error(data.error || '注销失败');
    }

    return res.status(200).json(data);
  } catch (error) {
    console.error('[Logout] 请求异常:', error);
    return res.status(500).json({
      message: '注销过程中发生错误',
      error: error instanceof Error ? error.message : '未知错误'
    });
  }
}