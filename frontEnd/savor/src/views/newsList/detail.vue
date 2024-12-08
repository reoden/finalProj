<template>
  <div class="container">
    <Header :back="true" />
    <div class="content">
      <Detail
        :loading="loading"
        :data="detail"
        type="news"
        @back="handleBack"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted, watch } from 'vue';
  import { useRouter } from 'vue-router';
  import { useAppStore } from '@/store';
  import useLoading from '@/hooks/loading';
  import Header from '@/components/header/index.vue';
  import SavorApi from '@/api/savor';
  import Detail from '../chefList/components/detail.vue';

  const originDetail = ref<any>({});
  const detail = ref<any>({});
  const { loading, toggle: toggleLoading } = useLoading(false);
  const router = useRouter();
  const globalStore = useAppStore();

  const changeLangue = (newVal: any) => {
    if (newVal === 'zh-CN') {
      detail.value = originDetail.value;
    } else {
      detail.value = {
        ...originDetail.value,
        ...originDetail.value.en
      };
    }
  }

  // 获取厨师列表数据
  const fetchData = async () => {
    toggleLoading();
    const { id } = router.currentRoute.value.query;
    try {
      const { data: { code, body } } = await SavorApi.newsDetail(id);
      if (code === 0) {
        originDetail.value = body;
        changeLangue(globalStore.langue);
      }
    } finally {
      toggleLoading();
    }
  };
  
  const handleBack = () => {
    router.push({ name: 'newsList' });
  }

  onMounted(() => {
    fetchData();
    watch(() => globalStore.langue, (newVal) => {
      changeLangue(newVal);
    });
  });
</script>

<style lang="less" scoped>
  .container {
    height: 100%;
    width: 100%;
    color: #fff;
    font-family: PingFang SC;
    display: flex;
    flex-direction: column;

    .content {
      width: 100%;
      background: #20222C;
      position: relative;
      display: flex;
      justify-content: center;
      flex: 1;
      overflow-y: auto;
    }
  }
</style>
