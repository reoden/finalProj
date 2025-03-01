<template>
  <div class="nav">
    <div class="item" @click="goTo">
      {{ nav.settled }}
    </div>
    <div class="item">
      {{ nav.about }}
    </div>
    <div class="item" :class="{ user: nav.login }" @click="login">
      <span v-if="!nav.login">{{ nav.loginText }}</span>
      <template v-else>
        <img :src="userImage" alt="" />
        <span>{{ nav.showPhone }}</span>
      </template>
    </div>
    <img :src="navImage" alt="" @click="changeLangue" />
    <Login ref="loginRef" />
  </div>
</template>

<script lang="ts" setup>
  import navImage from '@/assets/images/language.png';
import userImage from '@/assets/images/user.png';
import Login from '@/components/login/index.vue';
import useLocale from '@/hooks/locale';
import { useUserStore } from '@/store';
import { isLogin } from '@/utils/auth';
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

  const userInfo = useUserStore();
  const { t } = useI18n();
  const router = useRouter();
  const loginRef = ref(null);
  const pattern = /^(\d{3})\d{4}(\d{4})$/;
  const nav = computed(() => {
    return {
      login: isLogin(),
      loginText: t('home.nav.login'),
      settled: t('home.nav.settled'),
      about: t('home.nav.about'),
      showPhone: userInfo.phone && userInfo.phone.replace(pattern, '$1****$2'),
    };
  });
  const { changeLocale, currentLocale } = useLocale();

  const changeLangue = () => {
    const defaultLocale = localStorage.getItem('arco-locale') || 'zh-CN';
    if (defaultLocale === 'zh-CN') {
      changeLocale('en-US');
    } else {
      changeLocale('zh-CN');
    }
  };

  const goTo = () => {
    router.push({ name: 'settled' });
  };

  const login = () => {
    if (loginRef.value) {
      loginRef.value.showModal();
    }
  };
</script>

<style lang="less" scoped>
  .nav {
    display: flex;
    height: 100%;
    font-family: PingFang SC;
    font-size: 16px;
    font-weight: 600;
    line-height: 21.45px;
    text-align: center;
    color: #ffffff;
    align-items: center;

    .item {
      margin-left: 8px;
      padding: 3px 8px;
      cursor: pointer;
      &:hover {
        background: #3a3e4d;
      }
    }

    .user {
      border: 1px solid #ffffff;
      border-radius: 25px;
      padding: 5px 16px;
      display: flex;
      align-items: center;
      img {
        width: 27px;
        height: 27px;
        margin: 0px;
        margin-right: 5px;
      }
      &:hover {
        background: #272a37;
      }
    }

    img {
      width: 22px;
      height: 22px;
      margin-left: 16px;
      cursor: pointer;
    }
  }
</style>
