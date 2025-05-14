<template>
  <div class="game-container">
    <div class="game-background"></div>
    <!-- 修改为上中下布局 -->
    <div class="top-area">
      <div class="room-info">
        房间号: {{ roomData?.id || '加载中...' }}
      </div>
      <!-- 添加当前玩家提示 -->
      <div class="current-turn">
        当前玩家: {{ getCurrentPlayerName }}
        <span v-if="isMyTurn" class="my-turn">(我的回合)</span>
      </div>
      <div class="direction-indicator" :class="{ reversed: roomData?.direction !== 'clockwise' }">
        {{ roomData?.direction === 'clockwise' ? '↻' : '↺' }}
      </div>
    </div>

    <div class="main-area">
      <!-- 左侧玩家 -->
      <div class="side-players left">
        <div v-for="player in leftPlayers" 
             :key="player.id"
             :class="['player-box', { 
               'active': isCurrentPlayer(player),
               'self': isSelfPlayer(player),
               'creator': isCreator(player)
             }]">
          <div class="avatar">{{ player.avatar || '😊' }}</div>
          <div class="cards-count">{{ player.hand?.length || 0 }}</div>
          <div class="player-name">
            {{ player.id }}
            <span v-if="isSelfPlayer(player)" class="player-tag">(我)</span>
            <span v-if="isCreator(player)" class="player-tag host">(房主)</span>
          </div>
        </div>
      </div>

      <!-- 中央牌区 -->
      <div class="center-area">
        <!-- 牌堆 -->
        <div class="deck">
          <div class="draw-pile" @click="drawCard" :class="{ active: isMyTurn }">
            <div class="card-back"></div>
          </div>
          <div class="discard-pile">
            <UnoCard v-if="currentCard"
                    :color="currentCard.color"
                    :value="currentCard.value"
                    :type="currentCard.type" />
          </div>
        </div>
      </div>

      <!-- 右侧玩家 -->
      <div class="side-players right">
        <div v-for="player in rightPlayers" 
             :key="player.id"
             :class="['player-box', { 
               'active': isCurrentPlayer(player),
               'self': isSelfPlayer(player),
               'creator': isCreator(player)
             }]">
          <div class="avatar">{{ player.avatar || '😊' }}</div>
          <div class="cards-count">{{ player.hand?.length || 0 }}</div>
          <div class="player-name">
            {{ player.id }}
            <span v-if="isSelfPlayer(player)" class="player-tag">(我)</span>
            <span v-if="isCreator(player)" class="player-tag host">(房主)</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 重新设计的底部手牌区域 -->
    <div class="bottom-area">
      <div class="my-cards-wrapper">
        <div class="my-cards">
          <TransitionGroup name="card" tag="div" class="cards-container">
            <UnoCard v-for="(card, index) in myCards"
                    :key="index"
                    :color="card.color"
                    :value="card.value"
                    :type="card.type"
                    :hoverable="true"
                    :selected="selectedCard === card"
                    :playable="isMyTurn"
                    :class="{ 'shake-animation': card === shakingCard }"
                    @click="handleCardClick(card)" />
          </TransitionGroup>
        </div>
      </div>

      <!-- 修改游戏控制按钮区域 -->
      <div class="game-controls">
        <button class="control-btn" 
                @click="drawCard"
                :disabled="!isMyTurn">
          摸牌
        </button>
        <button class="control-btn" 
                @click="confirmPlayCard"
                :disabled="!isMyTurn || !selectedCard">
          出牌
        </button>
      </div>
    </div>

    <!-- 添加颜色选择器弹窗 -->
    <div v-if="showColorPicker" class="color-picker-modal">
      <div class="color-picker">
        <h3>选择颜色</h3>
        <div class="color-grid">
          <button
            v-for="color in ['red', 'blue', 'green', 'yellow']"
            :key="color"
            :class="['color-btn', color]"
            @click="selectColor(color)"
          ></button>
        </div>
      </div>
    </div>

    <!-- 添加获胜弹窗 -->
    <WinnerModal 
      :show="showWinnerModal"
      :winner="winner"
      @beforeLeave="handleWinnerLeave"
    />
  </div>
</template>

<script>
import { mapState } from 'vuex';
import UnoCard from '../components/uno/Card.vue';
import WinnerModal from '../components/uno/WinnerModal.vue';

export default {
  name: 'Game',
  components: {
    UnoCard,
    WinnerModal
  },
  data() {
    return {
      selectedCard: null,
      shakingCard: null,
      showColorPicker: false,
      pendingWildCard: null,
      selectedColor: null,
      showAIThinking: false,  // 添加一个数据属性来显示AI思考中
      winner: null,
      showWinnerModal: false
    }
  },
  computed: {
    ...mapState('unogame', ['roomData', 'player']),
    currentPlayerId() {
      return this.player?.nickname;
    },
    currentCard() {
      return this.roomData?.discardPile?.[this.roomData.discardPile.length - 1];
    },
    myCards() {
      return this.roomData?.players?.find(p => p.id === this.currentPlayerId)?.hand || [];
    },
    isMyTurn() {
      const currentPlayerIndex = this.roomData?.currentPlayerIndex;
      return this.roomData?.players[currentPlayerIndex]?.id === this.currentPlayerId;
    },
    leftPlayers() {
      return this.roomData?.players?.slice(0, Math.ceil(this.roomData.players.length / 2)) || [];
    },
    rightPlayers() {
      return this.roomData?.players?.slice(Math.ceil(this.roomData.players.length / 2)) || [];
    },
    getCurrentPlayerName() {
      const currentPlayerIndex = this.roomData?.currentPlayerIndex;
      const currentPlayer = this.roomData?.players[currentPlayerIndex];
      return currentPlayer?.id || '等待中';
    },
    isRoomCreator() {
      return this.roomData?.creator === this.player?.nickname;
    }
  },
  methods: {
    isCurrentPlayer(player) {
      const currentPlayerIndex = this.roomData?.currentPlayerIndex;
      return this.roomData?.players[currentPlayerIndex]?.id === player.id;
    },
    isSelfPlayer(player) {
      return player.id === this.currentPlayerId;
    },
    isCreator(player) {
      return player.id === this.roomData?.creator;
    },
    canPlayCard(card) {
      const topCard = this.currentCard;
      return card.color === topCard.color || 
             card.value === topCard.value || 
             card.type === 'wild';
    },
    playCard(card) {
      if (this.isMyTurn && this.canPlayCard(card)) {
        this.$store.dispatch('unogame/playCard', card);
      }
    },
    drawCard() {
      if (this.isMyTurn) {
        this.$store.dispatch('unogame/drawCard');
      }
    },
    getPlayerClasses(player, index) {
      return {
        'current-player': this.isCurrentPlayer(index),
        'self': player.id === this.currentPlayerId,
        'creator': player.id === this.roomData?.creator
      };
    },
    async handleCardClick(card) {
      if (!this.isMyTurn) return;

      if (card.type === 'wild' || card.type === 'wild_draw_four') {
        this.pendingWildCard = card;
        this.showColorPicker = true;
      } else {
        this.selectedCard = card;
      }
    },
    async selectColor(color) {
      if (this.pendingWildCard) {
        const cardToPlay = {
          value: this.pendingWildCard.value,
          color: color,
          type: this.pendingWildCard.type
        };
        this.showColorPicker = false;
        this.selectedCard = cardToPlay;
        await this.confirmPlayCard();
      }
    },
    async confirmPlayCard() {
      if (this.selectedCard && this.isMyTurn) {
        try {
          // 统一按照 value, color, type 的顺序发送数据
          const cardData = {
            value: this.selectedCard.value,
            color: this.selectedCard.color,
            type: this.selectedCard.type
          };
          
          const result = await this.$store.dispatch('unogame/playCard', cardData);
          if (result === 'shake') {
            this.shakeCard(this.selectedCard);
          } else {
            this.selectedCard = null;
          }
        } catch (error) {
          console.error('出牌失败:', error);
          this.shakeCard(this.selectedCard);
        }
      }
    },
    shakeCard(card) {
      this.shakingCard = card;
      setTimeout(() => {
        this.shakingCard = null;
      }, 820); // 动画持续时间加上一点延迟
    },
    cancelSelection() {
      this.selectedCard = null;
    },
    async checkForAITurn() {
      if (!this.isRoomCreator) return;

      const currentPlayer = this.roomData?.players[this.roomData.currentPlayerIndex];
      if (currentPlayer?.type === 'bot') {
        try {
          console.log('AI回合，正在计算出牌...');
          
          // 添加延时并显示加载状态
          this.showAIThinking = true;  // 添加一个数据属性来显示AI思考中
          await new Promise(resolve => setTimeout(resolve, 1000));

          const aiMove = await this.$store.dispatch('unogame/handleAITurn', {
            roomId: this.roomData.id,
            aiPlayerIndex: this.roomData.currentPlayerIndex
          }).catch(async error => {
            // 如果失败了，尝试摸牌
            if (error.code === 'ECONNABORTED') {
              console.log('AI出牌超时，改为摸牌');
              await this.$store.dispatch('unogame/drawCard');
            } else {
              throw error;
            }
          });
          
          this.showAIThinking = false;
          console.log('AI决定出牌:', aiMove);
        } catch (error) {
          this.showAIThinking = false;
          console.error('AI出牌失败:', error);
        }
      }
    },
    async handleWinnerLeave() {
      try {
        // 先处理离开房间的逻辑
        await this.$store.dispatch('unogame/leaveRoom');
        
        // 使用 push 并处理导航错误
        await this.$router.push({ 
          path: '/',
          replace: true  // 使用 replace 模式
        }).catch(err => {
          if (err.name === 'NavigationDuplicated') {
            // 忽略重复导航错误
            return;
          }
          throw err;
        });
      } catch (error) {
        console.error('离开房间失败:', error);
        if (error.name !== 'NavigationDuplicated') {
          alert('离开房间失败，请重试');
        }
      }
    }
  },
  watch: {
    // 修改 watch，确保及时响应玩家变化
    'roomData': {
      handler(newData, oldData) {
        // 检查是否有玩家胜利
        const winner = newData?.players?.find(p => p.hand?.length === 0);
        if (winner) {
          this.winner = winner;
          this.showWinnerModal = true;
          // 停止游戏相关的操作
          this.selectedCard = null;
          this.showColorPicker = false;
        }
        if (newData?.currentPlayerIndex !== oldData?.currentPlayerIndex) {
          this.checkForAITurn();
        }
      },
      deep: true,
      immediate: true
    }
  },
  created() {
    // 初始化时获取房间信息
    this.$store.dispatch('unogame/fetchRoomInfo', this.$route.params.roomId);
  }
}
</script>

<style>
/* 全局样式 */
html, body {
  margin: 0;
  padding: 0;
  overflow: hidden;
  position: fixed;
  width: 100%;
  height: 100%;
}
</style>

<style scoped>
.game-container {
  width: 100vw;
  height: 100vh;
  position: fixed;  /* 改为 fixed 定位 */
  top: 0;
  left: 0;
  background: #e8f4ea;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.top-area {
  height: 48px;
  min-height: 48px;
  background: white;
  z-index: 2;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}

.main-area {
  flex: 1;
  display: flex;
  position: relative;
  height: calc(100vh - 250px);
  min-height: 0;
}

.side-players {
  width: 180px;
  padding: 10px;
  overflow: hidden;  /* 改为 hidden */
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.player-box {
  flex-shrink: 0;
  padding: 12px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.center-area {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  min-width: 0;
}

.bottom-area {
  height: 250px;
  min-height: 250px;
  padding: 20px;
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 25px;  /* 增加间距 */
  padding-bottom: 30px;  /* 减小底部内边距 */
}

.my-cards {
  flex: 1;
  position: relative;
  display: flex;
  justify-content: center;
  align-items: flex-start;  /* 改为顶部对齐 */
  padding-top: 20px;  /* 减小顶部内边距 */
}

.cards-container {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  gap: -30px;  /* 负间距让卡牌重叠 */
  padding: 0 60px;
}

.game-controls {
  display: flex;
  justify-content: center;
  gap: 20px;
  padding: 10px;  /* 减小内边距 */
  margin-top: -20px;  /* 向上移动按钮 */
  background: transparent;  /* 修改为透明背景 */
}

.control-btn {
  min-width: 100px;
  padding: 12px 24px;
  border: none;
  border-radius: 20px;
  font-size: 16px;
  color: white;  /* 改为白色文字 */
  cursor: pointer;
  transition: all 0.2s ease;
  background: #90caf9;  /* 使用柔和的蓝色 */
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.control-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: #e0e0e0;
}

.control-btn:not(:disabled):hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.2);
  background: #64b5f6;  /* 悬浮时稍深一点的蓝色 */
}

/* 确保所有内容都不会超出容器 */
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

/* 移除默认的页面滚动 */
:root {
  overflow: hidden;
}

/* 优化小屏幕显示 */
@media (max-height: 700px) {
  .bottom-area {
    height: 180px;
    min-height: 180px;
  }
  
  .top-area {
    height: 40px;
    min-height: 40px;
  }
  
  .main-area {
    height: calc(100vh - 220px);
  }
}

.direction-indicator {
  color: #4CAF50;
  font-size: 24px;
  text-shadow: 0 0 10px rgba(76, 175, 80, 0.5);
}

.avatar {
  font-size: 24px;
  margin-bottom: 8px;
}

.cards-count {
  font-size: 18px;
  color: #666;
  margin: 5px 0;
}

.player-name {
  font-size: 14px;
  color: #333;
}

.player-box.active {
  background: rgba(76, 175, 80, 0.1);
  border: 2px solid #4CAF50;
}

.player-box.self {
  background: rgba(33, 150, 243, 0.1);
  border: 2px solid #2196F3;
}

.player-box.creator {
  border: 2px solid #FFC107;
}

.player-tag {
  font-size: 12px;
  padding: 2px 4px;
  border-radius: 4px;
  margin-left: 4px;
}

.player-tag.host {
  background: #FFC107;
  color: #000;
}

.current-turn {
  color: #333;
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.my-turn {
  color: #4CAF50;
  font-weight: bold;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  10%, 30%, 50%, 70%, 90% { transform: translateX(-5px); }
  20%, 40%, 60%, 80% { transform: translateX(5px); }
}

.shake-animation {
  animation: shake 0.8s cubic-bezier(.36,.07,.19,.97) both;
}

.shake-animation:hover {
  animation: none;  /* 悬浮时停止抖动 */
}

.color-picker-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.color-picker {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
}

.color-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
  margin-top: 15px;
}

.color-btn {
  width: 60px;
  height: 60px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s;
}

.color-btn:hover {
  transform: scale(1.1);
}

.color-btn.red { background: #f44336; }
.color-btn.blue { background: #2196f3; }
.color-btn.green { background: #4caf50; }
.color-btn.yellow { background: #ffc107; }
</style>
