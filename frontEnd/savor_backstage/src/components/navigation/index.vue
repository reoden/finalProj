<template>
  <div class="menu">
    <a-menu
      mode="horizontal"
      :selected-keys="selectedKey"
      :show-collapse-button="false"
    >
      <a-menu-item v-for="menu in meunData" :key="menu.name" @click="goTo(menu.name)">
        {{ menu.title }}
      </a-menu-item>
    </a-menu>
  </div>
</template>

<script lang="ts" setup>
  import { useI18n } from 'vue-i18n';
  import { useRoute, useRouter } from 'vue-router'

  const { t } = useI18n();
  const route = useRoute();
  const router = useRouter();
  const parentNodName = route?.meta?.parentName;
  const selectedKey = parentNodName || route.name;
  console.log('route.name', route);
  const meunData = [{
    title: t('nav.restaurant.title'),
    name: 'home',
  }, {
    title: t('nav.chef.title'),
    name: 'chefList',
  }, {
    title: t('nav.guest.title'),
    name: 'guestList',
  }, {
    title: t('nav.news.title'),
    name: 'newsList',
  }];

  const goTo = (name: string) => {
    router.push({ name });
  }

</script>

<style scoped lang="less">
  .menu {
    flex: 1;
    height: 100%;
    position: relative;
    display: flex;

    .logo-wrap {
      height: 80px;
      display: flex;
      align-items: center;
      margin-left: 36px;
    }

    :deep(.arco-menu) {
      background: #272A37;
      .arco-menu-inner {
        padding: 4px 16px;
      }

      .arco-menu-overflow-wrap {
        text-align: right;
        margin-right: 90px;
      }

      .arco-menu-item {
        color: #fff;
        background: transparent;
      }

      .arco-menu-selected {
        background: #272A37;
        color: #fff;
        border-radius: 5px;
        &:hover {
          background: #272A37;
        }
      }

      .arco-menu-selected-label {
        background-color: #fff;
      }
    }

    .user {
      position: absolute;
      bottom: 36px;
      height: 36px;
      left: 50%;
      transform: translateX(-50%);
      white-space: nowrap;
      border: 1px solid #FFFFFF;
      border-radius: 25px;
      padding: 5px 16px;
      display: flex;
      color: #fff;
      align-items: center;
    }
  }
</style>
