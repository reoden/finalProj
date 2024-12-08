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
      redirect: 'home',
    },
    {
      path: '/home',
      name: 'home',
      component: () => import('@/views/home/index.vue'),
      meta: {
        requiresAuth: false,
      },
    },
    {
      path: '/home/list',
      name: 'homeList',
      component: () => import('@/views/home/list.vue'),
      meta: {
        requiresAuth: false,
        parentName: 'home',
      },
    },
    {
      path: '/home/detail',
      name: 'homeDetail',
      component: () => import('@/views/home/detail.vue'),
      meta: {
        requiresAuth: false,
        parentName: 'home',
      },
    },
    {
      path: '/settled',
      name: 'settled',
      component: () => import('@/views/settled/index.vue'),
      meta: {
        requiresAuth: false,
      },
    },
    {
      path: '/chef',
      name: 'chefList',
      component: () => import('@/views/chefList/index.vue'),
      meta: {
        requiresAuth: false,
      },
    },
    {
      path: '/chef/detail',
      name: 'chefDetail',
      component: () => import('@/views/chefList/detail.vue'),
      meta: {
        requiresAuth: false,
        parentName: 'chefList',
      },
    },
    {
      path: '/guest',
      name: 'guestList',
      component: () => import('@/views/guestList/index.vue'),
      meta: {
        requiresAuth: false,
      },
    },
    {
      path: '/guest/detail',
      name: 'guestDetail',
      component: () => import('@/views/guestList/detail.vue'),
      meta: {
        requiresAuth: false,
        parentName: 'guestList',
      },
    },
    {
      path: '/news',
      name: 'newsList',
      component: () => import('@/views/newsList/index.vue'),
      meta: {
        requiresAuth: false,
      },
    },
    {
      path: '/news/detail',
      name: 'newsDetail',
      component: () => import('@/views/newsList/detail.vue'),
      meta: {
        requiresAuth: false,
        parentName: 'newsList',
      },
    },
    {
      path: '/prize',
      name: 'prizeList',
      component: () => import('@/views/prizeList/index.vue'),
      meta: {
        requiresAuth: false,
      },
    },
    {
      path: '/prize/detail',
      name: 'prizeDetail',
      component: () => import('@/views/prizeList/detail.vue'),
      meta: {
        requiresAuth: false,
        parentName: 'prizeList',
      },
    },
    // {
    //   path: '/mobile',
    //   name: 'mobileIndex',
    //   component: () => import('@/views/mobile/index.vue'),
    //   meta: {
    //     requiresAuth: false,
    //   },
    // },
    // {
    //   path: '/mobile/list',
    //   name: 'mobileList',
    //   component: () => import('@/views/mobile/list.vue'),
    //   meta: {
    //     requiresAuth: false,
    //   },
    // },
    // {
    //   path: '/mobile/detail',
    //   name: 'mobileDetail',
    //   component: () => import('@/views/mobile/detail.vue'),
    //   meta: {
    //     requiresAuth: false,
    //   },
    // },
    REDIRECT_MAIN,
    NOT_FOUND_ROUTE,
  ],
  scrollBehavior() {
    return { top: 0 };
  },
});

createRouteGuard(router);

export default router;
