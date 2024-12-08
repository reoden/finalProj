<template>
  <div class="logo_cn" @click="goTo">
    <img :src="logoImage" alt="">
    <span>{{ en ? logoInfo.title : '餐厅年鉴' }}</span>
  </div>
</template>

<script lang="ts" setup>
  import { useRouter } from 'vue-router';
  import { useI18n } from 'vue-i18n';
  import { computed } from 'vue';
  import logoImage from '@/assets/images/logo.png';

  const { t } = useI18n();
  const router = useRouter();
  const logoInfo = computed(() =>
    {
      return {
        title: t('logo.title')
      }
    }
  );
  const props = defineProps({
    en: {
      type: Boolean,
      default() {
        return true;
      },
    }
  });
  const goTo = () => {
    router.push({ name: 'home' });
  }
</script>

<style scoped lang="less">
  .logo_cn {
    display: flex;
    align-items: center;
    cursor: pointer;
    img {
      width: 44px;
      height: 44px;
      margin-right: 14px;
    }
    span {
      font-family: PingFang SC;
      font-size: 18px;
      font-weight: 600;
      text-align: center;
      color: #fff;
      white-space: nowrap;
    }
  }
</style>
