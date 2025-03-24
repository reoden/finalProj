<template>
  <div class="container">
    <HeaderModile />
    <div class="content" ref="scrollRef" id="content">
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
      <div class="banner" v-if="userInfo.bannerList">
        <div class="cnt">
          <div v-for="item in userInfo.bannerList" :key="item.picture">
            <img :src="item.picture" alt="" srcset="">
          </div>
        </div>
      </div>
      <div class="shop-list">
        <a-spin v-if="loading" :size="32" />
        <div v-else class="shop-con">
          <div v-for="item of shopList" :key="item.id" class="list-item" @click="showDetail(item)">
            <img class="img" :src="item?.pics_url?.[0]" alt="">
            <div class="detail">
              <div class="name">{{  item.name }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import SavorApi from '@/api/savor';
import SearchIcon from '@/assets/images/search_icon.png';
import HeaderModile from '@/components/header/mobile.vue';
import useLoading from '@/hooks/loading';
import { useUserStore } from '@/store';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

  const userInfo = useUserStore();
  const router = useRouter();
  const inputVal = ref('');
  const shopList = ref<any>([]);
  const { loading, toggle: toggleLoading } = useLoading(false);

  // 获取首页餐厅列表数据
  const fetchData = async () => {
    toggleLoading();
    try {
      const { data: { code, body } } = await SavorApi.positionList();
      if (code === 0) {
        if (body && body.apps) {
          shopList.value = body.apps;
        }
      }
    } finally {
      toggleLoading();
    }
  };

  const searchVal = () => {
    router.push({
      name: 'mobileList',
      query: {
        search: inputVal.value,
      },
    });
  }

  const showDetail = (data: any) => {
    router.push({
      name: 'mobileDetail',
      query: {
        id: data.id,
      },
    });
  }

  onMounted(() => {
    fetchData();
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

      .banner {
        width: 100%;
        margin: 22px auto 0;
        padding-left: 23px;

        .cnt {
          overflow-x: auto;
          display: flex;

          div {
            width: 233px;
            height: 144px;
            border-radius: 10px;
            margin-right: 8px;
            &:last-child {
              margin-right: 23px;
            }
          }

          img {
            width: 233px;
            height: 144px;
            border-radius: 10px;
            object-fit: cover;
          }
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
