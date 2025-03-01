<template>
  <div class="header">
    <div class="left">
      <img v-if="back" class="back" :src="backImage" alt="" @click="handleBack">
      <a-trigger
        :trigger="['click', 'hover']"
        clickToClose
        position="top"
        v-model:popupVisible="popupVisible"
      >
        <div :class="`button-trigger ${popupVisible ? 'button-trigger-active' : ''}`">
          <icon-menu-unfold v-if="popupVisible" class="menu" />
          <icon-menu-fold v-else class="menu" />
        </div>
        <template #content>
          <a-menu
            :style="{ marginBottom: '-4px' }"
            :tooltipProps="{ position: 'left' }"
            :selected-keys="selectedKey"
            theme="dark"
          >
            <a-menu-item v-for="menu in meunData" :key="menu.name" @click="goTo(menu.name)">
              {{ t(menu.lable) }}
            </a-menu-item>
          </a-menu>
        </template>
      </a-trigger>
      <!-- <a-button
        class="menu"
        :style="{ padding: '0 12px', height: '30px', lineHeight: '30px', marginBottom: '4px' }"
        type="primary"
        @click="toggleCollapse"
      >
        <icon-menu-unfold v-if="collapsed" />
        <icon-menu-fold v-else />
      </a-button> -->
      <!-- <a-popover :popup-visible="visible" trigger="click" content-class="city-modal" arrow-class="city-arrow" position="bottom">
        <icon-menu class="menu" />
        <template #content>
          <Navigation />
        </template>
      </a-popover> -->
      <img :src="homeImage" alt="" @click="handleHome">
    </div>
    <span>{{ $t('logo.title') }}</span>
    <Nav />
  </div>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { useI18n } from 'vue-i18n';
  import { useRoute, useRouter } from 'vue-router';
  import { useUserStore } from '@/store';
  import backImage from '@/assets/images/back.png';
  import homeImage from '@/assets/images/home.png';
  import Nav from './nav.vue';


  const { t, locale } = useI18n();
  const route = useRoute();
  const router = useRouter();
  const userInfo = useUserStore();
  const parentNodName = route?.meta?.parentName;
  const selectedKey = parentNodName || route.name;
  const emit = defineEmits(['back']);
  const popupVisible = ref(false);
  const isCn = locale.value === 'zh-CN';

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
      lable: userInfo.menuInfo?.prize_name,
      name: 'prizeList',
    });
  }

  const goTo = (name: string) => {
    router.push({ name });
  }
 
  const props = defineProps({
    back: {
      type: Boolean,
      default() {
        return false;
      },
    }
  });

  const handleHome = () => {
    router.push({ name: 'home' })
  }

  const handleBack = () => {
    // emit('back');
    router.back();
  }

</script>

<style scoped lang="less">
  .header {
    height: 46px;
    padding: 0px 16px;
    background-color: #272A37;
    text-align: center;
    position: relative;
    display: flex;
    justify-content: space-between;
    align-items: center;

    span {
      font-family: PingFang SC;
      font-size: 17px;
      font-weight: 600;
      vertical-align: middle;
      color: #F2F2F2;
    }

    .left {
      display: flex;

      .menu {
        height: 24px;
        font-size: 20px;
        vertical-align: 0;
        margin-right: 6px;
        color: #fff;
      }

      .back {
        margin-right: 6px;
      }

      img {
        width: 12px;
        height: 24px;

        &:last-child {
          width: 24px;
        }
      }
    }
  }
</style>
