<template>
  <div id="app">
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'App',
  async created() {
    // 初始化 store 状态，包括重新建立 WebSocket 连接
    await this.$store.dispatch('unogame/initializeStore');
    
    // 添加页面刷新事件监听
    window.addEventListener('beforeunload', () => {
      // 在页面刷新前保存必要的状态
      if (this.$store.state.unogame.roomId) {
        localStorage.setItem('lastRoomId', this.$store.state.unogame.roomId);
      }
    });
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
