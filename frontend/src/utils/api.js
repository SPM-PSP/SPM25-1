import axios from 'axios';
import config from '../config';

const api = axios.create({
  baseURL: config.api.baseURL,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
    'Authorization': config.api.AUTH_TOKEN
  },
  timeout: config.api.timeout
});

// 请求拦截器
api.interceptors.request.use(request => {
  console.log('请求配置:', {
    url: request.url,
    method: request.method,
    headers: request.headers,
    data: request.data
  });
  return request;
});

// 响应拦截器
api.interceptors.response.use(
  response => {
    console.log('响应成功:', {
      status: response.status,
      headers: response.headers,
      data: response.data
    });
    return response;
  },
  error => {
    console.error('API错误详情:', {
      config: error.config,
      status: error.response?.status,
      statusText: error.response?.statusText,
      headers: error.response?.headers,
      data: error.response?.data,
      message: error.message
    });
    throw error;
  }
);

export const createRoom = async ({ roomId, playerInfo }) => {
  try {
    const response = await api.post('/createRoom', {
      room_id: roomId,
      creator: {
        id: playerInfo.nickname,
        avatar: playerInfo.avatar
      }
    });
    
    console.log('创建房间响应:', response.data);
    return response.data;
  } catch (error) {
    console.error('创建房间失败:', error.response?.data || error);
    throw error;
  }
};

export const joinRoom = async ({ roomId, playerInfo }) => {
  if (!roomId) {
    throw new Error('房间ID不能为空');
  }
  
  try {
    const response = await api.post('/ws/joinRoom', {  // 修改为正确的HTTP接口路径
      room_id: roomId.toString(), // 确保转换为字符串
      player: {
        id: playerInfo.id || playerInfo.nickname,
        avatar: playerInfo.avatar,
        type: playerInfo.type || ''  // 如果是机器人则传入 'bot'，否则为空
      }
    });
    
    return response.data;
  } catch (error) {
    console.error('加入房间失败:', error);
    throw error;
  }
};

export const getRoomById = async (roomId) => {
  try {
    const response = await api.post('/getRoomById', {
      room_id: roomId
    });
    
    if (response.data.error) {
      throw new Error(response.data.error);
    }
    return response.data;
  } catch (error) {
    console.error('服务器错误:', error.response?.data || error);
    throw new Error(error.response?.data?.error || '获取房间信息失败');
  }
};

export const startUnoGame = async (roomId) => {
  try {
    const response = await api.post('/StartUno', {
      room_id: roomId
    });
    console.log('开始游戏响应:', response.data);
    return response.data;
  } catch (error) {
    console.error('开始游戏失败:', error);
    throw error;
  }
};

export const checkCard = async (roomId, playerIndex, card) => {
  try {
    console.log('发送检查卡牌请求:', {
      room_id: roomId,
      player_index: playerIndex,
      card: card
    });
    
    const response = await api.post('/Uno/checkCard', {
      room_id: roomId,
      player_index: playerIndex,
      card: card
    });
    
    console.log('检查卡牌响应:', response.data);
    return response.data;
  } catch (error) {
    console.error('检查卡牌失败:', error.response?.data || error);
    throw error;
  }
};

export const playCard = async (roomId, card, choose = '') => {
  try {
    console.log('发送出牌请求:', {
      room_id: roomId,
      card: card,
      choose: choose
    });
    
    const response = await api.post('/Uno/handleSpecial', {
      room_id: roomId,
      card: card,
      choose: choose
    });
    
    console.log('出牌响应:', response.data);
    return response.data;
  } catch (error) {
    console.error('出牌失败:', error.response?.data || error);
    throw error;
  }
};

export const drawNewCard = async (roomId, playerId, number = 1) => {
  try {
    console.log('发送摸牌请求:', { 
      room_id: roomId,
      player_id: playerId,
      number: number  
    });
    
    const response = await api.post('/Uno/draw', {
      room_id: roomId,
      player_id: playerId,
      number: number
    });
    
    console.log('摸牌响应:', response.data);
    return response.data;
  } catch (error) {
    console.error('摸牌失败:', error.response?.data || error);
    throw error;
  }
};

export const acceptPenalty = async (roomId) => {
  try {
    console.log('接受惩罚摸牌:', { room_id: roomId });
    
    const response = await api.post('/Uno/accept', {
      room_id: roomId
    });
    
    console.log('接受惩罚响应:', response.data);
    return response.data;
  } catch (error) {
    console.error('接受惩罚失败:', error.response?.data || error);
    throw error;
  }
};

export const leaveRoom = async (roomId, playerId) => {
  try {
    const response = await api.post('/ws/leaveRoom', {
      room_id: roomId,
      player: {
        id: playerId
      }
    });
    return response.data;
  } catch (error) {
    console.error('离开房间失败:', error);
    throw error;
  }
};

export const getAIMove = async (roomId, aiPlayerIndex) => {
  try {
    console.log('请求AI出牌:', {
      room_id: roomId,
      ai_player_index: aiPlayerIndex
    });
    
    // 添加延长的超时时间
    const response = await api.post('/Uno/chat', {
      room_id: roomId,
      ai_player_index: aiPlayerIndex
    }, {
      timeout: 20000  // 为AI操作特别设置20秒超时
    });
    
    console.log('AI出牌响应:', response.data);
    return response.data;
  } catch (error) {
    if (error.code === 'ECONNABORTED') {
      console.error('AI思考时间过长，正在重试...');
      // 如果是超时错误，自动重试一次
      try {
        const retryResponse = await api.post('/Uno/chat', {
          room_id: roomId,
          ai_player_index: aiPlayerIndex
        }, {
          timeout: 30000  // 重试时使用更长的超时时间
        });
        return retryResponse.data;
      } catch (retryError) {
        console.error('AI操作重试失败:', retryError);
        throw retryError;
      }
    }
    console.error('AI出牌失败:', error);
    throw error;
  }
};

export default api;
