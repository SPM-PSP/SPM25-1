<template>
  <div class="create-player">
    <h2>{{ isJoining ? '加入房间' : '创建角色' }}</h2>
    <div class="create-form">
      <div class="input-group">
        <label>你的昵称</label>
        <input v-model="nickname" type="text" placeholder="请输入昵称">
      </div>
      
      <div class="avatar-group">
        <label>选择表情头像</label>
        <div class="avatar-list">
          <div 
            v-for="(emoji, index) in avatars" 
            :key="index"
            :class="['avatar-item', { selected: selectedAvatar === emoji }]"
            @click="selectAvatar(emoji)"
          >
            <span class="emoji">{{ emoji }}</span>
          </div>
        </div>
      </div>

      <button 
        class="create-btn" 
        :disabled="!isValid"
        @click="createPlayer"
      >
        创建
      </button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'CreatePlayer',
  props: {
    isJoining: {
      type: Boolean,
      default: false
    },
    roomToJoin: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      nickname: '',
      selectedAvatar: null,
      avatars: [
        '😊', '😎', '🤠', '🤓', 
        '😄', '🐱', '🐶', '🐼',
        '🦊', '🐯', '🦁', '🐸'
      ]
    }
  },
  computed: {
    isValid() {
      return this.nickname && this.selectedAvatar;
    }
  },
  methods: {
    selectAvatar(emoji) {
      this.selectedAvatar = emoji;
    },
    async createPlayer() {
      if (this.isValid) {
        try {
          // 创建玩家信息
          const playerInfo = {
            nickname: this.nickname,
            avatar: this.selectedAvatar
          };

          if (this.isJoining && this.roomToJoin) {
            // 验证名字是否重复
            await this.$store.dispatch('unogame/validatePlayerName', {
              roomId: this.roomToJoin,
              nickname: this.nickname
            });
          }

          await this.$store.dispatch('unogame/createPlayer', playerInfo);

          if (this.isJoining && this.roomToJoin) {
            // 如果是加入模式，直接使用现有房间号
            await this.$store.dispatch('unogame/joinRoom', {
              roomId: this.roomToJoin,
              playerInfo
            });
            this.$router.push({
              name: 'lobby',
              params: { roomId: this.roomToJoin }
            });
          } else {
            // 创建新房间
            const result = await this.$store.dispatch('unogame/createRoom');
            if (result.success) {
              this.$router.push({
                name: 'lobby',
                params: { roomId: result.roomId }
              });
            }
          }
        } catch (error) {
          console.error('操作失败:', error);
          alert(error.message || '操作失败，请重试');
        }
      }
    }
  }
}
</script>

<style scoped>
.create-player {
  max-width: 600px;
  margin: 40px auto;
  padding: 20px;
}

.input-group {
  margin-bottom: 20px;
}

.input-group input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin-top: 8px;
}

.avatar-list {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
  margin-top: 10px;
}

.avatar-item {
  cursor: pointer;
  padding: 10px;
  border: 2px solid transparent;
  border-radius: 8px;
  transition: all 0.3s;
  display: flex;
  justify-content: center;
  align-items: center;
}

.emoji {
  font-size: 32px;
}

.avatar-item.selected {
  border-color: #4CAF50;
}

.create-btn {
  width: 100%;
  padding: 12px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 20px;
}

.create-btn:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}
</style>
