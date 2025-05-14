<template>
  <div class="winner-modal" v-if="show">
    <div class="modal-content">
      <div class="confetti">🎉</div>
      <h2>游戏结束</h2>
      <div class="winner-info">
        <span class="winner-avatar">{{ winner.avatar || '😊' }}</span>
        <h3>{{ winner.id }} 获胜！</h3>
      </div>
      <div class="actions">
        <button @click="handleLeave" class="leave-btn">离开房间</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'WinnerModal',
  props: {
    show: Boolean,
    winner: Object
  },
  methods: {
    handleLeave() {
      // 先触发事件，让父组件处理离开逻辑
      this.$emit('beforeLeave');
      // 使用导航守卫确保安全离开
      this.$router.push({ path: '/' }).catch(err => {
        if (err.name !== 'NavigationDuplicated') {
          console.error('导航错误:', err);
          // 如果不是重复导航错误，再次触发离开事件
          this.$emit('leave');
        }
      });
    }
  }
}
</script>

<style scoped>
.winner-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 2000;
}

.modal-content {
  background: white;
  padding: 40px;
  border-radius: 16px;
  text-align: center;
  animation: pop-in 0.5s ease-out;
}

.winner-info {
  margin: 20px 0;
}

.winner-avatar {
  font-size: 48px;
  display: block;
  margin-bottom: 10px;
}

.confetti {
  font-size: 32px;
  margin-bottom: 20px;
}

.leave-btn {
  padding: 12px 24px;
  background: #4CAF50;
  color: white;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  font-size: 16px;
  transition: all 0.2s;
}

.leave-btn:hover {
  background: #45a049;
  transform: translateY(-2px);
}

@keyframes pop-in {
  0% { transform: scale(0.5); opacity: 0; }
  100% { transform: scale(1); opacity: 1; }
}
</style>
