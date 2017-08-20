import Vue from 'vue';
import Router from 'vue-router';
import Root from '@/components/Root';
import Nodes from '@/components/Nodes';
import Images from '@/components/Images';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Root',
      component: Root,
    },
    {
      path: '/nodes',
      name: 'Nodes',
      component: Nodes,
    },
    {
      path: '/images',
      name: 'Images',
      component: Images,
    },
  ],
});
