<template>
  <div class="container">
    <Header :back="true" />
    <div ref="scrollRef" class="content" @scroll="debounceScroll">
      <div class="banner" :class="{ mobile: mobile }" v-if="bannerInfo">
        <div class="cnt">
          <div>
            <a v-if="bannerInfo.redirect_url" :href="bannerInfo.redirect_url" target="_blank">
              <img :src="bannerInfo.pic_url" alt="" srcset="">
            </a>
            <img v-else :src="bannerInfo.pic_url" alt="" srcset="">
          </div>
        </div>
      </div>
      <div class="city" :class="{ mobile: mobile }">
        <div class="title">{{ $t('nav.list.city.title') }}</div>
        <div class="list">
          <div v-for="city in hotCity" :key="city.cnLabel" class="item" :class="{ selected : city.cnLabel === selectedCity }" @click="selectCity(city.cnLabel)">
            {{ city.label }}
          </div>
          <a-popover :popup-visible="visible" trigger="click" :content-class="contentClass" arrow-class="city-arrow" position="bottom">
            <a-button @click="showCityMore">{{ $t('nav.list.city.more') }}</a-button>
            <template #content>
              <div
                v-for="item in addressOption" 
                :key="item.label" class="city_label"
                :class="{ selected : item.cnLabel === selectedCity }"
                @click="selectCity(item.cnLabel)">
                <span>{{ item.label }}</span>
              </div>
            </template>
          </a-popover>
        </div>
      </div>
      <div class="search">
        <a-input-search
          v-model:model-value="inputVal"
          :placeholder="$t('home.con.search1')"
        >
          <template #suffix>
            <img class="search-icon" :src="SearchIcon" alt="" @click="searchVal()">
          </template>
        </a-input-search>
      </div>
      <List
        v-if="!mobile"
        :loading="searchLoading"
        :data="searchList"
        :bottom="isBottom"
        type="shop"
        @fetchData="searchData"
        @showDetail="showDetail"
      />
      <div v-else class="shop-con">
        <div v-for="item of searchList" :key="item.id" class="list-item" @click="showDetail(item)">
          <img class="img" :src="item?.pics_url?.[0]" alt="">
          <div class="detail">
            <div class="name">{{  item.name }}</div>
          </div>
        </div>
      </div>
      <a-empty style="margin-bottom: 20px;" v-if="!loadingMore && !searchList.length"></a-empty>
      <a-spin v-if="loadingMore" :size="32" />
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted, computed, watch } from 'vue';
  import { useRouter } from 'vue-router';
  import { pca } from 'area-data';
  import { useAppStore } from '@/store';
  import useLoading from '@/hooks/loading';
  import { useDebounceFn } from '@vueuse/core';
  import Header from '@/components/header/index.vue';
  import SavorApi from '@/api/savor';
  import SearchIcon from '@/assets/images/search_icon.png';
  import List from '@/components/list/index.vue';
  import { isMoblie } from '@/utils';
  import { pinyin } from 'pinyin-pro';

  const router = useRouter();
  const inputVal = ref<any>('');
  const globalStore = useAppStore();
  const searchList = ref<any>([]);
  const originSearchList = ref<any>([]);
  const page = ref<any>(1);
  const isBottom = ref<any>(false);
  const { loading: searchLoading, toggle: toggleSearchLoading } = useLoading(false);
  const { loading: loadingMore, toggle: toggleLoadingMore } = useLoading(false);
  const originCity = ref<any>(['北京', '上海', '杭州', '深圳', '广州', '成都']);
  const hotCity = ref<any>(originCity.value);
  const selectedCity = ref<any>('');
  const mobile = computed(() => isMoblie());
  const bannerInfo = ref<any>(null);

  const contentClass = computed(() => {
    if (mobile.value) {
      return 'city-modal mobile';
    }
    return 'city-modal';
  });

  const changeLangue = (newVal: any) => {
    if (newVal === 'zh-CN') {
      searchList.value = originSearchList.value;
    } else {
      searchList.value = originSearchList.value.map(((item: any) => {
        return {
          ...item,
          ...item.en
        };
      }));
    }
  }

  const searchData = async (type?: any) => {
    if (type) {
      toggleLoadingMore();
    } else {
      toggleSearchLoading();
    }
    try {
      const params = { search: selectedCity.value };
      const { data: { code, body } } = await SavorApi.prizeList(params);
      if (code === 0) {
        if (body && body.apps) {
          originSearchList.value = body?.apps || [];
          changeLangue(globalStore.langue);
          isBottom.value = searchList.value.length === (body?.total || 0);
        } else {
          originSearchList.value = [];
          changeLangue(globalStore.langue);
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

  const deatilData = ref<any>({});
  const showDetail = (data: any) => {
    deatilData.value = data;
    router.push({
      name: 'prizeDetail',
      query: {
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

  const addressOption = ref<any>([]);
  const visible = ref(false);
  const initAddress = (newVal?: any) => {
    let proData: any = [];
    const add = pca['86'];
    const headerList: any = [];
    hotCity.value = originCity.value.map((item: any) => {
      return {
        cnLabel: item,
        label: newVal === 'en-US' ? pinyin(item, { toneType: 'none' }) : item
      };
    });
    Object.keys(add).forEach((p: any) => {
      const children: any = [];
      Object.keys(pca[p]).forEach((c: any) => {
        const v: any = { value: c, label: pca[p][c], cnLabel: pca[p][c] };
        v.label = newVal === 'en-US' ? pinyin(v.label, { toneType: 'none' }) : v.label;
        children.push(v)
      });
      const proitem: any = { value: p, label: pca['86'][p] };
      if (['北京市', '天津市', '上海市',  '重庆市', '台湾省', '香港特别行政区', '澳门特别行政区'].indexOf(proitem.cnLabel) < 0) {
        if (children.length) {
          proData = proData.concat(children);
        }
      } else {
        const item = { ...proitem, label: proitem.label.replaceAll('市', '').replaceAll('特别行政区', '').replaceAll('省', ''), cnLabel: proitem.label.replaceAll('市', '').replaceAll('特别行政区', '').replaceAll('省', '') };
        item.label = newVal === 'en-US' ? pinyin(item.label, { toneType: 'none' }) : item.label;
        headerList.push(item);
      }
    });
    proData = headerList.concat(proData);
    addressOption.value = proData;
  }

  const showCityMore = () => {
    visible.value = !visible.value;
  }

  const selectCity = (val: any) => {
    visible.value = false;
    if (val === selectedCity.value) {
      selectedCity.value = '';
    } else {
      selectedCity.value = val;
    }
    searchData();
  }

  const getBanner = async () => {
    try {
      const { data: { code, body } } = await SavorApi.prizeBanner();
      if (code === 0) {
        if (body) {
          bannerInfo.value = body;
        } else {
          bannerInfo.value = null;
        }
      }
    } finally {
      console.log();
    }
  }

  onMounted(() => {
    const { search } = router.currentRoute.value.query;
    inputVal.value = search || '';
    searchData();
    initAddress(globalStore.langue);
    getBanner();
    watch(() => globalStore.langue, (newVal) => {
      changeLangue(newVal);
      searchData();
      initAddress(newVal);
    });
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
      padding-top: 28px;
      position: relative;
      display: flex;
      flex-direction: column;
      align-items: center;
      overflow-y: auto;
      background: #20222C;


      .banner {
        width: 85%;
        margin: 0px auto 28px;

        .cnt {
          display: flex;
          justify-content: space-between;

          div {
            width: 100%;
            height: 212px;
          }

          img {
            width: 100%;
            height: 212px;
            object-fit: cover;
            border-radius: 10px;
          }
        }
      }

      .banner.mobile {
        padding: 0px 23px;
        width: 100%;

        .cnt {
          width: 100%;
          overflow-x: auto;
          overflow-y: hidden;
          display: flex;
          height: 144px;

          div {
            width: 100%;
            height: 144px;
            border-radius: 10px;
          }

          img {
            width: 100%;
            height: 144px;
            border-radius: 10px;
            object-fit: cover;
          }
        }
      }

      .title {
        font-size: 60px;
        font-family: PingFang SC;
        line-height: 128.7px;
        text-align: center;
        color: #FFFFFF;
      }

      .city {
        width: 75%;
        color: #fff;
        display: flex;
        align-items: center;
        margin-bottom: 40px;

        .title {
          font-size: 16px;
          font-weight: 400;
          line-height: 18.15px;
          margin-right: 24px;
        }

        .list {
          display: flex;
          align-items: center;

          .item {
            margin-right: 30px;
            cursor: pointer;

            &:hover {
              color: #8B95BC;
            }
            &.selected {
              border-bottom: 2px solid #fff;
            }
          }
        }
      }

      .mobile {
        flex-direction: column;
        align-items: flex-start;
        margin-bottom: 16px;

        .title {
          margin-bottom: 16px;
        }

        .list {
          flex-wrap: wrap;

          .item {
            margin-right: 10px;
          }
        }
      }

      .search {
        width: 78%;

        .search-icon {
          width: 30px;
          height: 30px;
          cursor: pointer;
        }

        .arco-input-wrapper {
          height: 36px;
          background: transparent;
          border: 1px solid #FFFFFF;
          border-radius: 50px;
          color: #fff;
          padding-right: 3px;
        }

        .arco-input-focus {
          border: 1px solid #FFFFFF
        }

        :deep(.arco-icon-hover) {
          display: none;
        }
      }

      .shop-list {
        margin-top: 70px;
        width: 90%;
        display: flex;
        display: flex;
        flex-wrap: wrap;
        justify-content: space-between;
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

      .shop-con {
        width: 100%;
        padding: 23px 23px 0;

        .list-item {
          border-radius: 5px;
          margin-bottom: 22px;
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
    }
  }
</style>

<style lang="less">
.city-modal {
  width: 1000px;
  background: #272A37;
  border: unset;
  color: #fff;

  .arco-popover-content {
    display: flex;
    flex-wrap: wrap;
    height: 400px;
    overflow-y: auto;
  }

  .city_label {
    width: 160px;
    // white-space: nowrap;
    margin-top: 8px;
    cursor: pointer;

    &:hover {
      color: #8B95BC;
    }
  }
  .city_label.selected span {
    padding-bottom: 5px;
    border-bottom: 2px solid #fff;
  }
}

.city-modal.mobile {
  width: 100%;
}
.city-arrow {
  background-color: #272A37;
  border: unset;
}
</style>
