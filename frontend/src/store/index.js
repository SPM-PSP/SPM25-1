import Vue from 'vue'
import Vuex from 'vuex'
import unogame from './unogame'

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    unogame
  }
})

export default store
