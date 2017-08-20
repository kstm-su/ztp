import Vue from 'vue';
import Router from 'vue-router';
import Root from '@/components/Root';
import Nodes from '@/components/Nodes';
import Node from '@/components/Node';
import Images from '@/components/Images';
import Image from '@/components/Image';

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
      path: '/nodes/:id',
      name: 'Node',
      component: Node,
    },
    {
      path: '/images',
      name: 'Images',
      component: Images,
    },
    {
      path: '/images/:id',
      name: 'Image',
      component: Image,
    },
  ],
});
