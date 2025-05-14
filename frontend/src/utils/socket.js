import config from '../config';

class WebSocketService {
  constructor() {
    this.ws = null;
    this.callbacks = {};
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectTimeout = null;
    this.wsParams = null;
  }

  connect({ roomId, playerId }) {
    if (this.ws) {
      // 如果已经有连接，先断开
      this.disconnect();
    }

    console.log('准备建立WebSocket连接:', { roomId, playerId });
    this.wsParams = { roomId, playerId };
    this.establishConnection();
  }

  establishConnection() {
    if (!this.wsParams) return;

    const wsUrl = `${config.websocket.url}?roomId=${this.wsParams.roomId}&playerId=${this.wsParams.playerId}`;
    console.log('WebSocket连接URL:', wsUrl);
    this.ws = new WebSocket(wsUrl);
    
    this.ws.onopen = () => {
      console.log('WebSocket连接成功, 参数:', this.wsParams);
      this.reconnectAttempts = 0;
    };
    
    this.ws.onclose = () => {
      console.log('WebSocket连接关闭，参数:', this.wsParams);
      this.reconnect();
    };
    
    this.ws.onerror = (error) => {
      console.error('WebSocket连接错误:', {
        error,
        params: this.wsParams,
        readyState: this.ws.readyState
      });
    };

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        console.log('收到WebSocket消息:', message);

        // 如果消息是一个对象并且包含 type 字段
        if (message && typeof message === 'object' && message.type) {
          console.log(`处理 ${message.type} 类型消息:`, message.data);
          if (this.callbacks[message.type]) {
            this.callbacks[message.type](message.data);
            return;
          }
        }

        // 如果消息没有正确的类型或处理器，当作普通消息处理
        if (this.callbacks['message']) {
          console.log('使用通用消息处理器:', message);
          this.callbacks['message'](message);
        }
      } catch (error) {
        console.error('WebSocket消息处理错误:', error);
      }
    };
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.wsParams = null;
    this.reconnectAttempts = 0;
    clearTimeout(this.reconnectTimeout);
    console.log('WebSocket连接已断开');
  }

  reconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('WebSocket 重连失败次数过多，停止重连');
      return;
    }

    this.reconnectAttempts++;
    clearTimeout(this.reconnectTimeout);
    
    this.reconnectTimeout = setTimeout(() => {
      console.log(`第 ${this.reconnectAttempts} 次尝试重连...`);
      this.establishConnection();
    }, 3000 * this.reconnectAttempts);
  }

  isConnected() {
    return this.ws && this.ws.readyState === WebSocket.OPEN;
  }

  on(type, callback) {
    console.log('注册WebSocket消息处理器:', {
      type,
      existingHandlers: Object.keys(this.callbacks)
    });
    this.callbacks[type] = callback;
  }

  emit(type, data) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      const message = {
        type,  // 不再固定为 system
        data
      };
      console.log('发送WebSocket消息:', message);
      this.ws.send(JSON.stringify(message));
    } else {
      console.warn('WebSocket未连接，无法发送消息');
    }
  }
}

export default new WebSocketService();
