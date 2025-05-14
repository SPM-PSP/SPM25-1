<template>
  <transition name="toast">
    <div v-if="visible" class="toast" :class="type">
      {{ message }}
    </div>
  </transition>
</template>

<script>
export default {
  name: 'Toast',
  data() {
    return {
      visible: false,
      message: '',
      type: 'info',
      timer: null
    }
  },
  methods: {
    show(message, type = 'info', duration = 2000) {
      clearTimeout(this.timer);
      this.visible = true;
      this.message = message;
      this.type = type;
      
      this.timer = setTimeout(() => {
        this.visible = false;
      }, duration);
    }
  }
}
</script>

<style scoped>
.toast {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  padding: 10px 20px;
  border-radius: 4px;
  color: white;
  font-size: 14px;
  z-index: 9999;
  box-shadow: 0 2px 8px rgba(0,0,0,0.2);
}

.toast.info {
  background: #2196F3;
}

.toast.success {
  background: #4CAF50;
}

.toast.warning {
  background: #FFC107;
  color: #333;
}

.toast.error {
  background: #f44336;
}

.toast-enter-active, .toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter, .toast-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-20px);
}
</style>
