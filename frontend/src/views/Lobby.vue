<template>
  <div class="lobby">
    <div class="room-header">
      <h2>房间号：{{ roomId }}</h2>
      <div class="share-section">
        <input 
          ref="shareLink" 
          :value="joinLink" 
          readonly 
          class="share-input"
        >
        <button @click="copyLink" class="copy-btn">
          {{ copySuccess ? '已复制' : '复制链接' }}
        </button>
      </div>
      <p class="share-tip">分享此链接邀请好友加入游戏</p>
    </div>

    <div v-if="roomData" class="players-container">
      <div class="player-grid">
        <div v-for="(player, index) in roomData.players" 
             :key="player.id"
             :class="['player-card', { 
               'current-player': player.id === currentNickname,
               'host': player.id === roomData.creator
             }]">
          <span class="emoji-avatar">{{ player.avatar || '😊' }}</span>
          <span class="player-name">
            {{ player.id }}
            <span v-if="player.id === currentNickname" class="self-tag">(我)</span>
            <span v-if="player.id === roomData.creator" class="host-tag">(房主)</span>
          </span>
        </div>
      </div>
    </div>
    <div v-else class="loading">
      加载中...
    </div>

    <div class="action-buttons">
      <button 
        v-if="isRoomCreator"
        class="add-bot-btn" 
        @click="handleAddBot"
        :disabled="!canAddBot">
        添加人机
      </button>
      <button 
        v-if="isRoomCreator"
        class="start-btn" 
        :disabled="!canStart" 
        @click="startGame">
        开始游戏
      </button>
      <button 
        class="leave-btn" 
        @click="handleLeaveRoom">
        离开房间
      </button>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex';

export default {
  name: 'Lobby',
  props: {
    roomId: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      maxPlayers: 4,
      copySuccess: false
    }
  },
  computed: {
    ...mapState('unogame', ['player', 'roomData', 'isConnected']),
    ...mapGetters('unogame', ['currentPlayer']),
    isHost() {
      return this.player && this.player.isHost;
    },
    roomLink() {
      return `http://localhost:8082/join/${this.roomId}`;
    },
    // 删除 remainingSlots 计算属性
    canStart() {
      return this.roomData?.players?.length >= 2;
    },
    joinLink() {
      return `${window.location.origin}/#/join/${this.roomId}`;
    },
    currentNickname() {
      return this.player?.nickname;
    },
    isCurrentPlayer() {
      return this.player && this.roomData?.players.some(p => p.id === this.player.nickname);
    },
    isCreator() {
      return this.roomData?.creator === this.player?.nickname;
    },
    canAddBot() {
      return this.roomData?.players?.length < 4;
    },
    // 修改判断条件，确保是当前房间的创建者
    isRoomCreator() {
      return this.roomData?.creator === this.player?.nickname;
    }
  },
  async created() {
    const roomId = this.$route.params.roomId;
    if (!roomId) {
      console.error('未找到房间号');
      this.$router.push('/');
      return;
    }

    try {
      const roomInfo = await this.$store.dispatch('unogame/fetchRoomInfo', roomId);
      
      // 检查游戏状态
      if (roomInfo.status === 'playing') {
        // 如果玩家在游戏中，跳转到游戏页面
        const playerInfo = JSON.parse(localStorage.getItem('playerInfo'));
        if (roomInfo.players.some(p => p.id === playerInfo.nickname)) {
          this.$router.replace(`/game/${roomId}`);
          return;
        } else {
          alert('游戏已经开始，无法加入');
          this.$router.push('/');
          return;
        }
      }
      
      // 检查玩家是否在房间中
      const playerInfo = JSON.parse(localStorage.getItem('playerInfo'));
      const isInRoom = roomInfo.players.some(p => p.id === playerInfo.nickname);
      
      if (!isInRoom && roomInfo.creator !== playerInfo.nickname) {
        console.warn('玩家不在房间中');
        this.$router.push(`/join/${roomId}`);
        return;
      }
    } catch (error) {
      console.error('获取房间信息失败:', error);
      alert('房间不存在或已失效');
      this.$router.push('/');
    }
  },
  methods: {
    copyLink() {
      const linkInput = this.$refs.shareLink;
      linkInput.select();
      document.execCommand('copy');
      this.copySuccess = true;
      setTimeout(() => {
        this.copySuccess = false;
      }, 2000);
    },
    async startGame() {
      try {
        await this.$store.dispatch('unogame/startGame');
        // 删除这里的路由跳转，由 WebSocket 消息处理统一处理
      } catch (error) {
        console.error('开始游戏失败:', error);
        alert('开始游戏失败，请重试');
      }
    },
    async handleAddBot() {
      if (!this.canAddBot) return;
      try {
        const roomId = this.$route.params.roomId;
        if (!roomId) throw new Error('房间ID不存在');
        
        // 添加加载状态
        const btn = document.querySelector('.add-bot-btn');
        if (btn) btn.textContent = '添加中...';
        
        await this.$store.dispatch('unogame/addBot', { roomId });
        
        // 恢复按钮文字
        if (btn) btn.textContent = '添加人机';
      } catch (error) {
        console.error('添加人机失败:', error);
        alert('添加人机失败');
        // 恢复按钮文字
        const btn = document.querySelector('.add-bot-btn');
        if (btn) btn.textContent = '添加人机';
      }
    },
    async handleLeaveRoom() {
      try {
        if (this.isRoomCreator) {
          if (!confirm('离开房间将导致房间解散，确定要离开吗？')) {
            return;
          }
        } else {
          if (!confirm('确定要离开房间吗？')) {
            return;
          }
        }

        await this.$store.dispatch('unogame/leaveRoom');
        // 使用 replace 而不是 push 来避免导航重复
        this.$router.replace('/').catch(err => {
          if (err.name !== 'NavigationDuplicated') {
            console.error('导航错误:', err);
          }
        });
      } catch (error) {
        console.error('离开房间失败:', error);
        alert('离开房间失败，请重试');
      }
    },
  }
}
</script>

<style scoped>
.lobby {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.room-header {
  text-align: center;
  margin-bottom: 30px;
}

.room-link {
  display: flex;
  gap: 10px;
  justify-content: center;
  margin-top: 10px;
}

room-link input {
  width: 300px;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.player-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin: 30px 0;
}

.player-card {
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  padding: 15px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.player-card.host {
  border-color: #4CAF50;
}

.player-card.empty {
  border-style: dashed;
  justify-content: center;
  color: #999;
}

.emoji-avatar {
  font-size: 24px;
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 30px;
}

.add-bot-btn, .start-btn {
  padding: 10px 20px;
  border-radius: 4px;
  border: none;
  cursor: pointer;
}

.add-bot-btn {
  background-color: #2196F3;
  color: white;
}

.start-btn {
  background-color: #4CAF50;
  color: white;
}

.start-btn:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.share-section {
  margin: 15px 0;
  display: flex;
  justify-content: center;
  gap: 10px;
}

.share-input {
  width: 300px;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.copy-btn {
  padding: 8px 16px;
  background: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.share-tip {
  color: #666;
  font-size: 14px;
}

.current-player {
  border: 2px solid #4CAF50;
  background-color: rgba(76, 175, 80, 0.1);
}

.self-tag {
  color: #4CAF50;
  font-weight: bold;
  margin-left: 4px;
}

.add-bot-btn {
  padding: 8px 16px;
  background: #2196F3;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin: 10px;
}

.add-bot-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.leave-btn {
  padding: 8px 16px;
  background: #f44336;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin: 10px;
}

.leave-btn:hover {
  background: #d32f2f;
}
</style>
