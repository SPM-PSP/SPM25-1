<template>
  <div class="uno-card" 
       :class="[color || 'wild', { 
         hoverable,
         selected,
         'in-discard': inDiscard,
         playable, 
         'special': isSpecial 
       }]" 
       @click="onCardClick">
    <div class="card-inner">
      <div class="card-corner top-left">{{ getSymbol }}</div>
      <div class="card-center">{{ getSymbol }}</div>
      <div class="card-corner bottom-right">{{ getSymbol }}</div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'UnoCard',
  props: {
    color: String,
    value: [String, Number],
    type: {
      type: String,
      required: true
    },
    playable: {
      type: Boolean,
      default: false
    },
    hoverable: {
      type: Boolean,
      default: false
    },
    selected: {
      type: Boolean,
      default: false
    },
    inDiscard: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    isSpecial() {
      return this.type !== 'number';
    },
    getSymbol() {
      switch(this.type) {
        case 'number': return this.value;
        case 'skip': return '⊘';
        case 'reverse': return '↺';
        case 'draw_two': return '+2';
        case 'wild': return '★';
        case 'wild_draw_four': return '+4';
        default: return '';
      }
    }
  },
  methods: {
    onCardClick() {
      // 移除 playable 检查，改为发出点击事件
      this.$emit('click', {
        color: this.color,
        value: this.value,
        type: this.type
      });
    }
  }
}
</script>

<style scoped>
.uno-card {
  width: 90px;
  height: 130px;
  background: white;
  border-radius: 8px;
  position: relative;
  transition: all 0.3s ease;
  margin: 0 -15px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  z-index: 1;
  cursor: pointer;  /* 默认可点击 */
}

.playable:hover {
  transform: translateY(-20px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.2);
  z-index: 10;
}

.hoverable:hover {
  transform: translateY(-30px);
  margin: 0 5px;  /* 悬浮时增加间距 */
  z-index: 10;
  cursor: pointer;  /* 添加指针样式 */
}

.selected {
  transform: translateY(-20px) !important;  /* 添加 !important 确保选中状态优先 */
  border: 2px solid #4CAF50;
  box-shadow: 0 4px 12px rgba(76,175,80,0.3);
  z-index: 5;
}

/* 基础卡牌样式 */
.card-inner {
  height: 100%;
  width: 100%;
  border-radius: 8px;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  color: white;
  font-weight: bold;
}

.card-corner {
  position: absolute;
  font-size: 1.2em;
}

.top-left {
  top: 8px;
  left: 8px;
}

.bottom-right {
  bottom: 8px;
  right: 8px;
  transform: rotate(180deg);
}

.card-center {
  font-size: 2.5em;
}

.red .card-inner { 
  background: #f44336;
}

.blue .card-inner { 
  background: #2196f3;
}

.green .card-inner { 
  background: #4caf50;
}

.yellow .card-inner { 
  background: #ffc107;
}

.wild .card-inner { 
  background: linear-gradient(90deg, 
    #f44336 25%, 
    #2196f3 25%, 
    #2196f3 50%, 
    #4caf50 50%, 
    #4caf50 75%, 
    #ffc107 75%
  );
}

.special .card-center {
  font-size: 2em;
}

/* 禁用状态样式 */
.uno-card[disabled] {
  opacity: 0.7;
  cursor: not-allowed;
}

/* 移除不需要的样式 */
.uno-card:not(.playable) {
  cursor: pointer;  /* 修改为可点击样式 */
}
</style>
