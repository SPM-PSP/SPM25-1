import Vue from 'vue'
import VueRouter from 'vue-router'
import HomeView from '../views/HomeView.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView
  },
  {
    path: '/create-player',
    name: 'createPlayer',
    component: () => import('../views/CreatePlayer.vue')
  },
  {
    path: '/lobby/:roomId?',
    name: 'lobby',
    component: () => import('../views/Lobby.vue'),
    props: true
  },
  {
    path: '/game/:roomId',
    name: 'game',
    component: () => import('../views/Game.vue')
  },
  {
    path: '/join/:roomId',
    name: 'join',
    component: () => import('../views/CreatePlayer.vue'),
    props: route => ({
      isJoining: true,
      roomToJoin: route.params.roomId
    })
  }
]

const router = new VueRouter({
  routes
})

router.beforeEach(async (to, from, next) => {
  // 延迟检查，确保 Vue 实例和 store 已完全初始化
  if (!router.app || !router.app.$store) {
    next();
    return;
  }

  const gameRoutes = ['lobby', 'game'];
  if (gameRoutes.includes(to.name)) {
    const roomId = to.params.roomId;
    if (!roomId) {
      next('/');
      return;
    }

    const savedPlayer = localStorage.getItem('playerInfo');
    if (!savedPlayer) {
      next({
        name: 'createPlayer',
        query: { redirect: to.fullPath }
      });
      return;
    }

    try {
      const playerInfo = JSON.parse(savedPlayer);
      await router.app.$store.dispatch('unogame/initializeConnection', {
        roomId,
        playerId: playerInfo.nickname
      });
      next();
    } catch (error) {
      console.error('初始化连接失败:', error);
      next('/');
    }
  } else {
    next();
  }
});

export default router
