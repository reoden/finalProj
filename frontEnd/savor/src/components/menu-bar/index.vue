<template>
  <div class="menu">
    <div class="logo-wrap">
      <Logo :en="false" />
    </div>
    <a-menu :selected-keys="selectedKey" :show-collapse-button="false">
      <a-menu-item
        v-for="menu in meunData"
        :key="menu.name"
        @click="goTo(menu.name)"
      >
        {{ menu.title }}
      </a-menu-item>
    </a-menu>
    <div class="user">
      {{ $t('menu.user.title') }}：
      <span>{{ showPhone }}</span>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { useI18n } from 'vue-i18n';
  import { useRoute, useRouter } from 'vue-router';
  import { useUserStore } from '@/store';
  import Logo from '@/components/logo/index.vue';

  const { t } = useI18n();
  const route = useRoute();
  const router = useRouter();
  const userInfo = useUserStore();
  const selectedKey = route.name;
  const pattern = /^(\d{3})\d{4}(\d{4})$/;
  const showPhone =
    userInfo.phone && userInfo.phone.replace(pattern, '$1****$2');
  const meunData = [
    {
      title: t('menu.examine.title'),
      name: 'examine',
    },
    {
      title: t('menu.position.title'),
      name: 'position',
    },
    {
      title: t('menu.chef.title'),
      name: 'chef',
    },
    {
      title: t('menu.guest.title'),
      name: 'guest',
    },
    {
      title: t('menu.news.title'),
      name: 'news',
    },
  ];

  const goTo = (name: string) => {
    router.push({ name });
  };
</script>

<style scoped lang="less">
  .menu {
    width: 231px;
    height: 100%;
    background: #20222c;
    position: relative;

    .logo-wrap {
      height: 80px;
      display: flex;
      align-items: center;
      margin-left: 36px;
    }

    :deep(.arco-menu) {
      background: #20222c;
      .arco-menu-inner {
        padding: 4px 16px;
      }

      .arco-menu-item {
        color: #fff;
        background: transparent;
      }

      .arco-menu-selected {
        background: #272a37;
        color: #fff;
        border-radius: 5px;
      }
    }

    .user {
      position: absolute;
      bottom: 36px;
      height: 36px;
      left: 50%;
      transform: translateX(-50%);
      white-space: nowrap;
      border: 1px solid #ffffff;
      border-radius: 25px;
      padding: 5px 16px;
      display: flex;
      color: #fff;
      align-items: center;
    }
  }
</style>
