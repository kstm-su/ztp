import Vue from 'vue';
import VueAxios from './http';

import App from './App';
import router from './router';

Vue.use(VueAxios, {
  baseURL: 'http://localhost:5001/',
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
