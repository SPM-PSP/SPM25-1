import Vue from 'vue';
import Vuex from 'vuex';
import socket from '../utils/socket';
import { createRoom, getRoomById, joinRoom, startUnoGame, checkCard, playCard, drawNewCard, leaveRoom, getAIMove, acceptPenalty } from '../utils/api';  // 添加 drawNewCard 和 leaveRoom
import router from '../router';  // 添加 router 导入

Vue.use(Vuex);

export default {
  namespaced: true,
  state: {
    players: [],
    currentPlayer: null,
    currentCard: null,
    myCards: [],
    roomId: null,
    gameStatus: 'waiting', // waiting, playing, ended
    player: null,  // 添加玩家信息
    isConnected: false,
    playerRooms: {}, // 存储玩家和房间的关联 { playerName: roomId }
    roomData: null,  // 存储房间详细数据
  },
  getters: {
    currentPlayer: (state) => {
      if (!state.player || !state.roomData || !state.roomData.players) {
        return null;
      }
      return state.roomData.players.find(p => p.id === state.player.nickname);
    }
  },
  mutations: {
    setPlayers(state, players) {
      state.players = players;
    },
    setCurrentPlayer(state, player) {
      state.currentPlayer = player;
    },
    setMyCards(state, cards) {
      state.myCards = cards;
    },
    setCurrentCard(state, card) {
      state.currentCard = card;
    },
    setRoomId(state, roomId) {
      state.roomId = roomId;
    },
    setGameStatus(state, status) {
      state.gameStatus = status;
    },
    setPlayer(state, playerInfo) {
      state.player = playerInfo;
      // 保存到 localStorage
      localStorage.setItem('playerInfo', JSON.stringify(playerInfo));
    },
    setConnected(state, status) {
      state.isConnected = status;
    },
    updatePlayers(state, players) {
      state.players = players;
    },
    setPlayerRoom(state, { playerName, roomId }) {
      state.playerRooms[playerName] = roomId;
      // 同时保存到 localStorage
      localStorage.setItem('playerRooms', JSON.stringify(state.playerRooms));
    },
    setRoomData(state, data) {
      state.roomData = data;
    }
  },
  actions: {
    async createRoom({ commit, state, dispatch }) {
      try {
        const roomId = `room-${Date.now().toString(36)}`;
        const result = await createRoom({ 
          roomId,
          playerInfo: state.player
        });
        
        if (result) {
          commit('setRoomId', roomId);
          // 房主只建立WebSocket连接，不调用joinRoom
          dispatch('connectWebSocket', { 
            roomId, 
            playerId: state.player.nickname 
          });
          commit('setConnected', true);
          
          // 更新房间信息
          await dispatch('fetchRoomInfo', roomId);
          return { success: true, roomId };
        }
        throw new Error('创建房间失败');
      } catch (error) {
        console.error('创建房间失败:', error);
        throw error;
      }
    },

    async joinRoom({ commit, dispatch, state }, { roomId, playerInfo }) {
      // 如果是房主或机器人，不需要建立websocket连接
      if (state.roomData?.creator === playerInfo.nickname || playerInfo.type === 'bot') {
        return { success: true };
      }

      try {
        const response = await joinRoom({ roomId, playerInfo });
        
        if (response) {
          commit('setRoomId', roomId);
          dispatch('connectWebSocket', { 
            roomId, 
            playerId: playerInfo.id || playerInfo.nickname 
          });
          commit('setConnected', true);

          // 使用 system 类型发送加入消息
          socket.emit('message', `玩家 ${playerInfo.nickname} 加入了房间`);
          
          return { success: true };
        }
        throw new Error('加入房间失败');
      } catch (error) {
        console.error('加入房间失败:', error);
        throw error;
      }
    },

    // 修改 WebSocket 消息处理
    initWebSocket({ dispatch, state, commit }) {
      // 通用消息处理
      socket.on('message', async (message) => {
        console.log('处理通用消息:', message);
        const { type, data } = message;

        // 根据消息类型处理
        if (type === 'system' && data === '游戏开始了') {
          console.log('收到游戏开始消息，准备跳转');
          commit('setGameStatus', 'playing');
          await dispatch('fetchRoomInfo', state.roomId);
          
          Vue.nextTick(() => {
            router.replace(`/game/${state.roomId}`).catch(err => {
              if (err.name !== 'NavigationDuplicated') {
                console.error('跳转失败:', err);
              }
            });
          });
          return;
        }

        // 如果是房主离开，其他玩家直接返回首页
        if (type === 'playerLeave') {
          if (data.isCreator) {
            alert('房主已离开，房间解散');
            dispatch('disconnectWebSocket');
            commit('setRoomId', null);
            commit('setRoomData', null);
            router.replace('/').catch(err => {
              if (err.name !== 'NavigationDuplicated') {
                console.error('导航错误:', err);
              }
            });
            return;
          }

          // 如果是普通玩家离开，尝试更新房间信息
          try {
            await dispatch('fetchRoomInfo', state.roomId);
          } catch (error) {
            // 忽略房间不存在的错误
            if (!error.message?.includes('房间未找到')) {
              console.error('更新房间信息失败:', error);
            }
          }
        }
      });

      // 修改系统消息处理
      socket.on('system', async (message) => {
        console.log('处理系统消息:', message);
        
        try {
          if (typeof message === 'string') {
            if (message.includes('加入了房间')) {
              await dispatch('fetchRoomInfo', state.roomId);
            }
          } else if (message.data === '游戏开始了') {
            console.log('收到游戏开始消息，准备跳转');
            commit('setGameStatus', 'playing');
            await dispatch('fetchRoomInfo', state.roomId);
            
            Vue.nextTick(() => {
              router.replace(`/game/${state.roomId}`).catch(err => {
                if (err.name !== 'NavigationDuplicated') {
                  console.error('跳转失败:', err);
                }
              });
            });
          }
        } catch (error) {
          console.error('处理系统消息错误:', error);
        }
      });

      // 游戏相关消息处理
      socket.on('gameStart', ({ roomId, gameState }) => {
        // ...existing code...
      });

      // 添加摸牌消息处理
      socket.on('cardDrawn', async (data) => {
        console.log('收到摸牌消息:', data);
        // 更新房间信息
        await dispatch('fetchRoomInfo', state.roomId);
      });

      // 添加出牌消息处理
      socket.on('cardPlayed', async (message) => {
        console.log('收到出牌消息:', message);
        // 更新房间信息
        await dispatch('fetchRoomInfo', state.roomId);
      });

      // 添加玩家离开消息处理
      socket.on('playerLeave', async (data) => {
        console.log('收到玩家离开消息:', data);
        try {
          // 先更新房间信息
          const roomInfo = await dispatch('fetchRoomInfo', state.roomId);
          
          // 如果是房主离开，其他玩家也要离开
          if (data.playerId === roomInfo.creator) {
            alert('房主已离开，房间解散');
            dispatch('disconnectWebSocket');
            commit('setRoomId', null);
            commit('setRoomData', null);
            router.push('/');
          }
        } catch (error) {
          console.error('处理玩家离开消息失败:', error);
        }
      });

      // 添加玩家加入消息处理
      socket.on('playerJoin', async (data) => {
        console.log('收到玩家加入消息:', data);
        await dispatch('fetchRoomInfo', state.roomId);
      });
    },

    // 修改 connectWebSocket 方法
    async connectWebSocket({ commit, dispatch }, { roomId, playerId }) {
      try {
        // 先断开现有连接
        if (socket.isConnected()) {
          socket.disconnect();
        }

        // 建立新连接
        socket.connect({ roomId, playerId });
        commit('setConnected', true);

        // 初始化消息处理
        dispatch('initWebSocket');
        
        // 等待连接成功
        await new Promise(resolve => {
          const checkConnection = () => {
            if (socket.isConnected()) {
              resolve();
            } else {
              setTimeout(checkConnection, 100);
            }
          };
          checkConnection();
        });

        return true;
      } catch (error) {
        console.error('WebSocket连接失败:', error);
        commit('setConnected', false);
        throw error;
      }
    },

    disconnectWebSocket({ commit }) {
      socket.disconnect();
      commit('setConnected', false);
      commit('setRoomData', null);
    },

    async playCard({ state, dispatch }, card) {
      if (state.roomData?.players[state.roomData?.currentPlayerIndex]?.id !== state.player?.nickname) {
        // 使用 toast 提示
        Vue.prototype.$toast.show('还没轮到你出牌', 'warning');
        return 'shake';
      }
      
      try {
        const playerIndex = state.roomData.currentPlayerIndex;
        const checkResult = await checkCard(state.roomId, playerIndex, {
          value: card.value,
          color: card.color,
          type: card.type
        });
        
        if (!checkResult.valid) {
          // 使用 toast 提示
          Vue.prototype.$toast.show('不能打出这张牌', 'warning');
          return 'shake';
        }

        console.log('开始出牌:', card);  // 添加日志
        // 调用出牌API，使用choose传递选择的颜色
        const playResult = await playCard(
          state.roomId,
          {
            value: card.value,
            color: card.color,
            type: card.type
          }
        );
        
        if (playResult.error) {
          throw new Error(playResult.error);
        }

        // 发送出牌消息
        socket.emit('message', {
          type: 'cardPlayed',
          data: {
            playerId: state.player.nickname,
            roomId: state.roomId,
            card: card
          }
        });

        // 出牌成功，更新房间信息
        await dispatch('fetchRoomInfo', state.roomId);
        return true;
      } catch (error) {
        console.error('出牌错误:', error);
        Vue.prototype.$toast.show(error.message || '出牌失败', 'error');
        throw error;
      }
    },

    async drawCard({ state, dispatch }) {
      const currentPlayerIndex = state.roomData?.currentPlayerIndex;
      const isMyTurn = state.roomData?.players[currentPlayerIndex]?.id === state.player?.nickname;
      
      if (!isMyTurn) {
        Vue.prototype.$toast.show('还没轮到你摸牌', 'warning');
        return;
      }

      try {
        // 检查是否有惩罚需要接受
        if (state.roomData.drawCount > 0) {
          console.log('接受惩罚摸牌:', {
            roomId: state.roomId,
            drawCount: state.roomData.drawCount
          });
          
          // 调用接受惩罚 API
          const result = await acceptPenalty(state.roomId);
          console.log('接受惩罚结果:', result);
        } else {
          console.log('普通摸牌:', {
            roomId: state.roomId,
            playerId: state.player.nickname
          });
          
          // 调用普通摸牌 API
          const result = await drawNewCard(state.roomId, state.player.nickname, 1);
          console.log('摸牌结果:', result);
        }

        // 发送websocket消息通知其他玩家
        socket.emit('message', {
          type: 'cardDrawn',
          data: {
            playerId: state.player.nickname,
            roomId: state.roomId
          }
        });

        // 更新房间信息
        await dispatch('fetchRoomInfo', state.roomId);
        return true;
      } catch (error) {
        console.error('摸牌错误:', error);
        alert(error.message || '摸牌失败');
        throw error;
      }
    },

    async leaveRoom({ state, commit, dispatch }) {
      try {
        const isCreator = state.roomData?.creator === state.player?.nickname;

        // 先通知其他玩家
        socket.emit('message', {
          type: 'playerLeave',
          data: {
            playerId: state.player.nickname,
            roomId: state.roomId,
            isCreator
          }
        });

        // 发送系统消息
        socket.emit('system', `玩家 ${state.player.nickname} 离开了房间`);

        // 调用离开房间 API
        await leaveRoom(state.roomId, state.player.nickname);
        
        // 断开连接并清理数据
        dispatch('disconnectWebSocket');
        commit('setRoomId', null);
        commit('setRoomData', null);
        
        return true;
      } catch (error) {
        console.error('离开房间失败:', error);
        throw error;
      }
    },

    createPlayer({ commit }, playerInfo) {
      commit('setPlayer', playerInfo);
      return Promise.resolve();
    },
    validatePlayerRoom({ state }, { playerName, roomId }) {
      return state.playerRooms[playerName] === roomId;
    },
    async addBot({ state, dispatch }, { roomId }) {
      try {
        if (!roomId) {
          throw new Error('房间ID不能为空');
        }

        const botInfo = {
          id: `Bot_${Date.now().toString(36)}`,
          nickname: `Bot_${Date.now().toString(36)}`,
          avatar: '🤖',
          type: 'bot'
        };

        const response = await joinRoom({ roomId, playerInfo: botInfo });
        if (response) {
          // 先更新房间信息
          await dispatch('fetchRoomInfo', roomId);
          
          // 再发送消息通知其他玩家
          socket.emit('message', {
            type: 'playerJoin',
            data: {
              player: botInfo,
              roomId: roomId
            }
          });
          
          return true;
        }
        throw new Error('添加人机失败');
      } catch (error) {
        console.error('添加人机失败:', error);
        throw error;
      }
    },
    async startGame({ state, commit, dispatch }) {
      if (state.roomData?.players?.length >= 2) {
        try {
          console.log('开始游戏，房间ID:', state.roomId);
          const result = await startUnoGame(state.roomId);
          
          // 更新状态和数据
          commit('setGameStatus', 'playing');
          await dispatch('fetchRoomInfo', state.roomId);

          // 使用 message 类型发送系统消息
          socket.emit('message', {
            type: 'system',
            data: '游戏开始了'
          });

          return result;
        } catch (error) {
          console.error('开始游戏失败:', error.message || error);
          throw new Error('开始游戏失败: ' + (error.message || '未知错误'));
        }
      } else {
        throw new Error('玩家数量不足，无法开始游戏');
      }
    },
    async fetchRoomInfo({ commit }, roomId) {
      try {
        const roomInfo = await getRoomById(roomId);
        if (roomInfo.error) {
          throw new Error(roomInfo.error);
        }
        commit('setRoomData', roomInfo);
        return roomInfo;
      } catch (error) {
        console.error('获取房间信息失败:', error);
        commit('setRoomData', null); // 清空房间数据
        
        // 如果是房间不存在，则不需要抛出错误
        if (error.message?.includes('房间未找到')) {
          return null;
        }
        throw error;
      }
    },
    async validatePlayerName({ state }, { roomId, nickname }) {
      try {
        const roomInfo = await getRoomById(roomId);
        if (roomInfo.players && roomInfo.players.some(p => p.id === nickname)) {
          throw new Error('该名字已被使用，请换一个名字');
        }
        return true;
      } catch (error) {
        throw error;
      }
    },
    // 添加新的初始化连接方法
    async initializeConnection({ commit, dispatch }, { roomId, playerId }) {
      try {
        // 设置基础状态
        commit('setRoomId', roomId);
        
        // 建立 WebSocket 连接
        await dispatch('connectWebSocket', { roomId, playerId });
        
        // 获取房间最新数据
        await dispatch('fetchRoomInfo', roomId);
        
        return true;
      } catch (error) {
        console.error('初始化连接失败:', error);
        throw error;
      }
    },

    // 修改现有方法,移除重复的初始化逻辑
    initializeStore({ commit }) {
      // 只恢复玩家信息
      const savedPlayer = localStorage.getItem('playerInfo');
      if (savedPlayer) {
        commit('setPlayer', JSON.parse(savedPlayer));
      }
    },

    async handleAITurn({ state, dispatch }, { roomId, aiPlayerIndex }) {
      try {
        console.log('开始AI出牌流程:', { roomId, aiPlayerIndex });
        
        // 获取AI的出牌选择
        const aiMove = await getAIMove(roomId, aiPlayerIndex);
        console.log('AI决策结果:', aiMove);

        // 检查结果是否有效
        if (!aiMove) {
          throw new Error('AI决策结果无效');
        }

        // 如果返回的牌有具体的type值，说明是要出这张牌
        if (aiMove.type) {
          console.log('AI选择出牌:', aiMove);
          // 整理卡牌数据顺序
          const cardToPlay = {
            value: aiMove.value,
            color: aiMove.color,
            type: aiMove.type
          };
          
          // 调用出牌API
          const playResult = await playCard(roomId, cardToPlay);
          if (playResult.error) {
            throw new Error(playResult.error);
          }

          // 发送出牌消息
          socket.emit('message', {
            type: 'cardPlayed',
            data: {
              playerId: state.roomData.players[aiPlayerIndex].id,
              roomId: roomId,
              card: cardToPlay
            }
          });
        } else {
          // 否则选择摸牌
          console.log('AI选择摸牌');
          await drawNewCard(roomId, state.roomData.players[aiPlayerIndex].id, 1);
          
          // 发送摸牌消息
          socket.emit('message', {
            type: 'cardDrawn',
            data: {
              playerId: state.roomData.players[aiPlayerIndex].id,
              roomId: roomId
            }
          });
        }

        // 更新房间信息
        await dispatch('fetchRoomInfo', roomId);
        return aiMove;
      } catch (error) {
        console.error('AI操作执行失败:', error);
        throw error;
      }
    },
  }
};
