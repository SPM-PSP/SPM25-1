import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import Toast from './components/common/Toast.vue'

Vue.use(ElementUI)

// 创建 Toast 构造器
const ToastConstructor = Vue.extend(Toast);
const toast = new ToastConstructor().$mount();
document.body.appendChild(toast.$el);

// 添加到 Vue 原型
Vue.prototype.$toast = {
  show(message, type) {
    toast.show(message, type);
  }
};

// 添加全局路由守卫来处理websocket连接
router.afterEach((to) => {
  const gameRoutes = ['lobby', 'game'];
  if (gameRoutes.includes(to.name)) {
    const savedPlayer = localStorage.getItem('playerInfo');
    if (savedPlayer && to.params.roomId) {
      const playerInfo = JSON.parse(savedPlayer);
      store.dispatch('unogame/initializeConnection', {
        roomId: to.params.roomId,
        playerId: playerInfo.nickname
      }).catch(console.error);
    }
  }
});

Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
