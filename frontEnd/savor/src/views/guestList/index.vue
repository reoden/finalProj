<template>
  <div class="container">
    <Header />
    <div class="content">
      <List
        :loading="loading"
        :data="guestList"
        type="guest"
        @showDetail="showDetail"
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
  import List from '../chefList/components/list.vue';

  const router = useRouter();
  const guestList = ref<any>([]);
  const originGuestList = ref<any>([]);
  const { loading, toggle: toggleLoading } = useLoading(false);
  const globalStore = useAppStore();

  const changeLangue = (newVal: any) => {
    if (newVal === 'zh-CN') {
      guestList.value = originGuestList.value;
    } else {
      guestList.value = originGuestList.value.map(((item: any) => {
        return {
          ...item,
          ...item.en
        };
      }));
    }
  }

  // 获取列表数据
  const fetchData = async () => {
    toggleLoading();
    try {
      const { data: { code, body } } = await SavorApi.guestList();
      if (code === 0) {
        originGuestList.value = body?.apps;
        changeLangue(globalStore.langue);
      }
    } finally {
      toggleLoading();
    }
  };

  const showDetail = (item: any) => {
    router.push({
      name: 'guestDetail',
      query: {
        id: item.id
      },
    });
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
    width: 100%;
    height: 100%;
    color: #fff;
    font-family: PingFang SC;
    display: flex;
    flex-direction: column;

    .content {
      flex: 1;
      overflow-y: auto;
      width: 100%;
      // height: calc(100% - 80px);
      background: #20222C;
      position: relative;
      display: flex;
      justify-content: center;
    }
  }
</style>./components/list.vue
