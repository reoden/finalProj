<template>
  <div class="nav" :class="{ mobile: mobile }">
    <div class="item" @click="goTo">
      {{ nav.settled }}
    </div>
    <a-trigger
      v-if="nav.login"
      trigger="click"
      position="top"
      :popup-translate="[0, 5]"
    >
      <div class="item user">
        <img :src="userImage" alt="" />
        <span v-if="!mobile">{{ nav.showPhone }}</span>
      </div>
      <template #content>
        <div class="trigger-demo-nest">
          <a-button @click="logout">{{ $t('savor.logout.title') }}</a-button>
        </div>
      </template>
    </a-trigger>
    <div v-else class="item" @click="login">
      <span v-if="!nav.login">{{ nav.loginText }}</span>
    </div>
    <img class="lang" :src="navImage" alt="" @click="changeLangue" />
    <Login ref="loginRef" :mobile="mobile" />
  </div>
</template>

<script lang="ts" setup>
  import { useRouter } from 'vue-router';
  import { ref, computed, provide } from 'vue';
  import { useI18n } from 'vue-i18n';
  import useLocale from '@/hooks/locale';
  import { isLogin } from '@/utils/auth';
  import navImage from '@/assets/images/language.png';
  import userImage from '@/assets/images/user.png';
  import Login from '@/components/login/index.vue';
  import { useUserStore, useAppStore } from '@/store';
  import { clearToken } from '@/utils/auth';
  import { isMoblie } from '@/utils';

  const userInfo = useUserStore();
  const useApp = useAppStore();
  const { t } = useI18n();
  const router = useRouter();
  const loginRef = ref(null);
  const pattern = /^(\d{3})\d{4}(\d{4})$/;
  const mobile = computed(() => isMoblie());
  const nav = computed(() => {
    return {
      login: isLogin(),
      loginText: t('home.nav.login'),
      settled: t('home.nav.settled'),
      showPhone: userInfo.phone && userInfo.phone.replace(pattern, '$1****$2'),
    };
  });
  const { changeLocale, currentLocale } = useLocale();
  const emit = defineEmits(['change']);

  const changeLangue = () => {
    const defaultLocale = localStorage.getItem('arco-locale') || 'zh-CN';
    let cn = defaultLocale;
    if (defaultLocale === 'zh-CN') {
      cn = 'en-US';
    } else {
      cn = 'zh-CN';
    }
    changeLocale(cn);
    emit('change');
    useApp.toggleLangue(cn);
  };

  const login = () => {
    if (isLogin()) {
      return;
    }
    if (loginRef.value) {
      loginRef.value.showModal();
    }
  };

  const goTo = () => {
    if (!isLogin()) {
      login();
      return;
    }
    router.push({ name: 'settled' });
  };

  const logout = () => {
    clearToken();
    if (router.currentRoute.value.name === 'home') {
      window.location.reload();
    } else {
      router.push({ name: 'home' });
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
      &:first-child {
        margin-left: 0px;
        padding-left: 0px;
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

  .mobile {
    .item {
      font-size: 14px;
      padding: 3px 0px;
      &:hover {
        background: #272a37;
      }
    }
    .user {
      border: none;
    }
    .lang {
      margin-left: 6px;
    }
  }
</style>
