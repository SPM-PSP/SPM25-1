<template>
  <div class="home">
    <h1 class="title">UNO 游戏</h1>
    <div class="actions">
      <button class="btn create-room" @click="createRoom">创建房间</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HomeView',
  methods: {
    createRoom() {
      this.$router.push('/create-player');
    }
  }
}
</script>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background-color: #f5f5f5;
}

.title {
  font-size: 48px;
  margin-bottom: 40px;
  color: #333;
}

.btn {
  padding: 12px 24px;
  font-size: 18px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  transition: all 0.3s;
}

.create-room {
  background-color: #4CAF50;
  color: white;
}

.create-room:hover {
  background-color: #45a049;
}
</style>
