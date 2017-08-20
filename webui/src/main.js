import Vue from 'vue';
import VueAxios from './http';

import App from './App';
import router from './router';

Vue.use(VueAxios, {
  baseURL: process.env.API_URL,
});
Vue.config.productionTip = false;

new Vue({
  el: '#app',
  router,
  template: '<App/>',
  components: {
    App,
  },
});
