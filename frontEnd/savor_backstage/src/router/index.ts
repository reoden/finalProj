import { createRouter, createWebHistory } from 'vue-router';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';

import { REDIRECT_MAIN, NOT_FOUND_ROUTE } from './routes/base';
import createRouteGuard from './guard';

NProgress.configure({ showSpinner: false }); // NProgress Configuration

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: 'backstage',
    },
    {
      path: '/backstage',
      name: 'backstage',
      component: () => import('@/views/backstage/index.vue'),
      meta: {
        requiresAuth: false,
      },
    },
    {
      path: '/backstage/examine',
      name: 'examine',
      component: () => import('@/views/examine/index.vue'),
      meta: {
        requiresAuth: false,
      },
    },
    {
      path: '/backstage/position',
      name: 'position',
      component: () => import('@/views/position/index.vue'),
      meta: {
        requiresAuth: false,
      },
    },
    {
      path: '/backstage/chef',
      name: 'chef',
      component: () => import('@/views/chef/index.vue'),
      meta: {
        requiresAuth: false,
      },
    },
    {
      path: '/backstage/chef/detail',
      name: 'chefBackDetail',
      component: () => import('@/views/chef/detail.vue'),
      meta: {
        requiresAuth: false,
        parentName: 'chef',
      },
    },
    {
      path: '/backstage/guest',
      name: 'guest',
      component: () => import('@/views/guest/index.vue'),
      meta: {
        requiresAuth: false,
      },
    },
    {
      path: '/backstage/guest/detail',
      name: 'guestBackDetail',
      component: () => import('@/views/guest/detail.vue'),
      meta: {
        requiresAuth: false,
        parentName: 'guest',
      },
    },
    {
      path: '/backstage/news',
      name: 'news',
      component: () => import('@/views/news/index.vue'),
      meta: {
        requiresAuth: false,
      },
    },
    {
      path: '/backstage/news/detail',
      name: 'newsBackDetail',
      component: () => import('@/views/news/detail.vue'),
      meta: {
        requiresAuth: false,
        parentName: 'news',
      },
    },
    {
      path: '/backstage/prize',
      name: 'prize',
      component: () => import('@/views/prize/index.vue'),
      meta: {
        requiresAuth: false,
      },
    },
    {
      path: '/backstage/prize/detail',
      name: 'prizeBackDetail',
      component: () => import('@/views/prize/detail.vue'),
      meta: {
        requiresAuth: false,
      },
    },
    REDIRECT_MAIN,
    NOT_FOUND_ROUTE,
  ],
  scrollBehavior() {
    return { top: 0 };
  },
});

createRouteGuard(router);

export default router;
