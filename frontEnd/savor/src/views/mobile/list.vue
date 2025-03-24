<template>
  <div class="container">
    <HeaderModile />
    <div class="content" ref="scrollRef" id="content" @scroll="debounceScroll">
      <div class="search">
        <a-input-search
          v-model:model-value="inputVal"
          placeholder="景点收集"
        >
          <template #suffix>
            <img class="search-icon" :src="SearchIcon" alt="" @click="searchVal()">
          </template>
        </a-input-search>
      </div>
      <div class="shop-list">
        <a-spin v-if="searchLoading || loading" :size="32" />
        <div v-else class="shop-con">
          <div v-for="item of searchList" :key="item.id" class="list-item" @click="showDetail(item)">
            <img class="img" :src="item?.pics_url?.[0]" alt="">
            <div class="detail">
              <div class="name">{{  item.name }}</div>
            </div>
          </div>
          <a-spin v-if="loadingMore" :size="32" />
        </div>
      </div>
      <a-back-top target-container="#content" :style="{position:'fixed'}" />
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { useDebounceFn } from '@vueuse/core';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

  import SavorApi from '@/api/savor';
import SearchIcon from '@/assets/images/search.png';
import HeaderModile from '@/components/header/mobile.vue';
import useLoading from '@/hooks/loading';

  const router = useRouter();
  const inputVal = ref<any>('');
  const searchList = ref<any>([]);
  const page = ref<any>(0);
  const isBottom = ref<any>(false);
  const { loading, toggle: toggleLoading } = useLoading(false);
  const { loading: searchLoading, toggle: toggleSearchLoading } = useLoading(false);
  const { loading: loadingMore, toggle: toggleLoadingMore } = useLoading(false);

  const searchData = async (type?: any) => {
    if (type) {
      toggleLoadingMore();
    } else {
      toggleSearchLoading();
    }
    try {
      const params = { page: page.value, page_size: 15, search: inputVal.value };
      const { data: { code, body } } = await SavorApi.searchShop(params);
      if (code === 0) {
        if (body && body.apps) {
          searchList.value = page.value === 1 ? body.apps : (searchList.value || []).concat(body.apps);
          isBottom.value = searchList.value === body.total;
        }
      }
    } finally {
      if (type) {
        toggleLoadingMore();
      } else {
        toggleSearchLoading();
      }
    }
  }

  const showDetail = (data: any) => {
    router.push({
      name: 'mobileDetail',
      query: {
        search: inputVal.value,
        id: data.id
      },
    });
  }

  const searchVal = () => {
    page.value = 1;
    searchData();
  }
  const scrollRef = ref();

  const handleScroll = () => {
    const { scrollTop, scrollHeight, clientHeight } = scrollRef.value;
    if (page.value !== 0 && clientHeight + scrollTop >= scrollHeight - 5 && !isBottom.value && !loadingMore.value) {
      page.value += 1;
      searchData('more');
    }
  }
  const debounceScroll = useDebounceFn(handleScroll, 300);

  onMounted(() => {
    page.value += 1;
    const { search } = router.currentRoute.value.query;
    inputVal.value = search || '';
    searchData();
  });

 
</script>

<style lang="less" scoped>
  .container {
    height: 100%;
    width: 100%;
    display: flex;
    flex-direction: column;

    .content {
      flex: 1;
      padding-top: 30px;
      position: relative;
      display: flex;
      flex-direction: column;
      align-items: center;
      overflow-y: auto;
      background: #20222C;

      .title {
        font-size: 60px;
        font-family: PingFang SC;
        line-height: 128.7px;
        text-align: center;
        color: #FFFFFF;
      }

      .search {
        width: 78%;

        .search-icon {
          width: 30px;
          height: 30px;
          cursor: pointer;
          margin-right: 10px;
        }

        .arco-input-wrapper {
          height: 55px;
          background: #3A3E4D;
          border: unset !important;;
          border-radius: 15px;
          color: #fff;
          padding-right: 3px;
        }

        .arco-input-focus {
          border: unset;
        }

        :deep(.arco-icon-hover) {
          display: none;
        }
      }

      

      .shop-list {
        width: 100%;
        padding-bottom: 30px;

        .shop-con {
          width: 100%;
          padding: 0px 23px;
        }

        .list-item {
          border-radius: 5px;
          margin-top: 22px;
          position: relative;

          .detail {
            position: absolute;
            width: 100%;
            height: 34px;
            bottom: 0px;
            left: 0px;
            right: 0px;
            background: #00000080;
            padding: 0px 14px;
            border-bottom-left-radius: 5px;
            border-bottom-right-radius: 5px;

            display: flex;
            align-items: center;

            .name {
              font-size: 14px;
              font-weight: 500;
              color: #fff;
            }
          }
        }

        .img {
          width: 100%;
          height: 193px;
          border-radius: 5px;
        }
      }

      &.search {
        height: 100%;
        background-color: #20222C;
      }

      .shop.hidden {
        width: 300px;
        height: 120px;
        visibility: hidden;
      }
    }
  }
</style>
