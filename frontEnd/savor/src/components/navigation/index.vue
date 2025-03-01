<template>
  <div class="menu">
    <a-menu
      mode="horizontal"
      :selected-keys="selectedKey"
      :show-collapse-button="false"
    >
      <a-menu-item v-for="menu in meunData" :key="menu.name" @click="goTo(menu.name)">
        {{ t(menu.lable) }}
      </a-menu-item>
    </a-menu>
  </div>
</template>

<script lang="ts" setup>
  import {  watch, ref } from 'vue';
  import { useI18n } from 'vue-i18n';
  import { useRoute, useRouter } from 'vue-router';
  import { useUserStore } from '@/store';

  const { t, locale } = useI18n();
  const route = useRoute();
  const router = useRouter();
  const userInfo = useUserStore();
  const parentNodName = route?.meta?.parentName;
  const selectedKey = parentNodName || route.name;
  const isCn = locale.value === 'zh-CN';

  const props = defineProps({
    lang: {
      type: Boolean,
      default() {
        return false;
      }
    }
  });
  const meunData = ref<any>([{
    key: 'app_name',
    title: 'nav.restaurant.title',
    lable: isCn ? userInfo.menuInfo?.app_name : 'nav.restaurant.title',
    name: 'home',
  }, {
    key: 'chef_name',
    title: 'nav.chef.title',
    lable: isCn ? userInfo.menuInfo?.chef_name : 'nav.chef.title',
    name: 'chefList',
  }, {
    key: 'guest_name',
    title: 'nav.guest.title',
    lable: isCn ? userInfo.menuInfo?.guest_name : 'nav.guest.title',
    name: 'guestList',
  }, {
    key: 'news_name',
    title: 'nav.news.title',
    lable: isCn ? userInfo.menuInfo?.news_name : 'nav.news.title',
    name: 'newsList',
  }]);
  
  if (userInfo.showPrize) {
    meunData.value.splice(1, 0, {
      key: 'prize_name',
      title: 'nav.prize.title',
      lable: isCn ? userInfo.menuInfo?.prize_name : 'nav.prize.title',
      name: 'prizeList',
    });
  }

  const checkLang = () => {
    const isCns = locale.value === 'zh-CN';
    meunData.value = meunData.value.map((item: any) => {
      return {
        key: item.key,
        title: item.title,
        lable: isCns ? userInfo.menuInfo[item.key] : item.title,
        name: item.name,
      }
    });
  }

  watch(
    () => props.lang,
    () => {
      checkLang();
    }
  );

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
        display: flex;
        justify-content: center;
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
          background: #272A37!important;
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
