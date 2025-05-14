export default {
  // API配置
  api: {
    baseURL: 'http://localhost:8082',  // 确保这是正确的基础URL
    AUTH_TOKEN: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcyNjk4NzAsInVzZXJuYW1lIjoiYWRtaW4zMzMifQ.rqB8jnqbVPjJdBGply2ttGozl6y8H2PtMPTl92SBSpw',  // 固定的API认证token
    getToken: () => localStorage.getItem('token') || '',  // API认证token
    timeout: 15000                          // 增加超时时间到15秒
  },
  
  // WebSocket配置
  websocket: {
    url: 'ws://localhost:8082/ws/joinRoom',  // WebSocket连接路径
    reconnectInterval: 3000,
    maxReconnectAttempts: 5
  }
}
