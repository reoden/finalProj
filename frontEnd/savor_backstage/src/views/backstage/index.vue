<template>
  <div class="container">
    <div class="left">
      <div class="logo-wrap">
        <Logo />
      </div>
      <div class="con">
        <div class="title">{{ $t('backstage.login.title') }}</div>
        <a-button class="btn" @click="handleLogin">{{
          $t('backstage.login.btn')
        }}</a-button>
      </div>
    </div>
    <img class="bg-img" :src="BgImage" alt="" />
    <Login ref="loginRef" @success="handleLoginSucc" />
  </div>
</template>

<script lang="ts" setup>
  import BgImage from '@/assets/images/bg.png';
import Login from '@/components/login/index.vue';
import Logo from '@/components/logo/index.vue';
import { isLogin } from '@/utils/auth';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

  const router = useRouter();
  const login = isLogin();
  const loginRef = ref(null);

  onMounted(() => {
    if (login) {
      router.push({ name: 'examine' });
    }
  });

  const handleLogin = () => {
    if (loginRef.value) {
      loginRef.value.showModal();
    }
  };

  const handleLoginSucc = () => {
    router.push({ name: 'examine' });
  };
</script>

<style lang="less" scoped>
  .container {
    height: 100%;
    width: 100%;
    background: #272a37;
    display: flex;

    .left {
      width: 50%;
      padding: 0px 36px;
      display: flex;
      flex-direction: column;

      .logo-wrap {
        height: 80px;
        display: flex;
        align-items: center;
      }

      .con {
        flex: 1;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        width: 400px;

        .title {
          font-family: PingFang SC;
          font-size: 50px;
          font-weight: 600;
          line-height: 70px;
          text-align: center;
          color: #fff;
          margin-bottom: 30px;
        }

        .btn {
          width: 100px;
          border: 1px solid #d1b276;
          font-weight: 600;
          background: #d1b276;
          color: #20222c;
          margin-left: -190px;
        }
      }
    }

    .bg-img {
      width: 50%;
      height: 100%;
    }
  }
</style>
