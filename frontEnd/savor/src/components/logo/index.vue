<template>
  <div class="logo_cn" @click="goTo">
    <img :src="userInfo.logoSmall || logoImage" alt="">
    <!-- <span>{{ en ? logoInfo.title : '景点发布' }}</span> -->
  </div>
</template>

<script lang="ts" setup>
  import logoImage from '@/assets/images/logo.png';
import { useUserStore } from '@/store';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

  const { t } = useI18n();
  const router = useRouter();
  const userInfo = useUserStore();
  const logoInfo = computed(() => {
    return {
      title: t('logo.title'),
    };
  });
  const props = defineProps({
    en: {
      type: Boolean,
      default() {
        return true;
      },
    },
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
      width: 171.72px;
      height: 57.04px;
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
