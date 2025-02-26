<template>
  <div class="container">
    <Header />
    <div class="content" :class="{ mobile: mobile }">
      <div class="title" :class="{ svg: !userInfo.logoBig }">
        <img v-if="userInfo.logoBig" :src="userInfo.logoBig" alt="" srcset="" />
      </div>
      <div class="search">
        <a-input-search
          v-model:model-value="inputVal"
          :placeholder="$t('home.con.search')"
        >
          <template #suffix>
            <img
              class="search-icon"
              :src="SearchIcon"
              alt=""
              @click="searchVal()"
            />
          </template>
        </a-input-search>
      </div>
      <div v-if="userInfo.bannerList" class="banner">
        <div class="cnt">
          <div v-for="item in userInfo.bannerList" :key="item.picture">
            <a
              v-if="item.redirect_url"
              :href="item.redirect_url"
              target="_blank"
            >
              <img :src="item.pic_url" alt="" srcset="" />
            </a>
            <img v-else :src="item.pic_url" alt="" srcset="" />
          </div>
        </div>
      </div>
      <div class="city" :class="{ mobile: mobile }">
        <div class="city-title">{{ $t('nav.list.city.title') }}</div>
        <div class="list">
          <div
            v-for="city in hotCity"
            :key="city.cnLabel"
            class="item"
            :class="{ selected: city.cnLabel === selectedCity }"
            @click="selectCity(city.cnLabel)"
          >
            {{ city.label }}
          </div>
          <a-popover
            :popup-visible="visible"
            trigger="click"
            :content-class="contentClass"
            arrow-class="city-arrow"
            position="bottom"
          >
            <a-button @click="showCityMore">{{
              $t('nav.list.city.more')
            }}</a-button>
            <template #content>
              <div
                v-for="item in addressOption"
                :key="item.label"
                class="city_label"
                :class="{ selected: item.label === selectedCity }"
                @click="selectCity(item.label)"
              >
                <span>{{ item.label }}</span>
              </div>
            </template>
          </a-popover>
        </div>
      </div>
      <div class="shop-list">
        <a-spin v-if="loading" :size="32" />
        <template v-if="!mobile">
          <template v-if="selectedCity">
            <ShopItem
              v-for="item in searchList"
              :key="item.id"
              :data="item"
              @show-detail="showDetail"
            />
          </template>
          <template v-else>
            <ShopItem
              v-for="item in shopList"
              :key="item.id"
              :data="item"
              @show-detail="showDetail"
            />
          </template>
          <div
            v-if="selectedCity ? searchList.length : shopList.length"
            class="shop hidden"
          />
          <div
            v-if="selectedCity ? searchList.length : shopList.length"
            class="shop hidden"
          />
        </template>
        <div v-else class="shop-con">
          <template v-if="selectedCity">
            <div
              v-for="item of searchList"
              :key="item.id"
              class="list-item"
              @click="showDetail(item)"
            >
              <img class="img" :src="item?.pics_url?.[0]" alt="" />
              <div class="detail">
                <div class="name">{{ item.name }}</div>
              </div>
            </div>
          </template>
          <template v-else>
            <div
              v-for="item of shopList"
              :key="item.id"
              class="list-item"
              @click="showDetail(item)"
            >
              <img class="img" :src="item?.pics_url?.[0]" alt="" />
              <div class="detail">
                <div class="name">{{ item.name }}</div>
              </div>
            </div>
          </template>
        </div>
        <a-empty
          v-if="selectedCity ? !searchList.length : !shopList.length"
          style="margin-bottom: 20px"
        ></a-empty>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted, computed, watch } from 'vue';
  import { useRouter } from 'vue-router';
  import { useUserStore } from '@/store';
  import useLoading from '@/hooks/loading';
  import { useAppStore } from '@/store';
  import Header from '@/components/header/index.vue';
  import SavorApi from '@/api/savor';
  import SearchIcon from '@/assets/images/search_icon.png';
  import { isMoblie } from '@/utils';
  import { pca } from 'area-data';
  import { pinyin } from 'pinyin-pro';
  import ShopItem from './components/shop-item.vue';

  const mobile = computed(() => isMoblie());
  const userInfo = useUserStore();
  const router = useRouter();
  const inputVal = ref('');
  const shopList = ref<any>([]);
  const originShopList = ref<any>([]);
  const searchList = ref<any>([]);
  const originSearchList = ref<any>([]);
  const { loading, toggle: toggleLoading } = useLoading(false);
  const originCity = ref<any>(['北京', '上海', '杭州', '深圳', '广州', '成都']);
  const hotCity = ref<any>(originCity.value);
  const selectedCity = ref<any>('');
  const globalStore = useAppStore();

  const contentClass = computed(() => {
    if (mobile.value) {
      return 'city-modal mobile';
    }
    return 'city-modal';
  });

  const changeLangue = (newVal: any) => {
    if (newVal === 'zh-CN') {
      shopList.value = originShopList.value;
    } else {
      shopList.value = originShopList.value.map((item: any) => {
        return {
          ...item,
          ...item.en,
        };
      });
    }
  };

  const changeLangueSecond = (newVal: any) => {
    if (newVal === 'zh-CN') {
      searchList.value = originSearchList.value;
    } else {
      searchList.value = originSearchList.value.map((item: any) => {
        return {
          ...item,
          ...item.en,
        };
      });
    }
  };

  // 获取首页餐厅列表数据
  const fetchData = async () => {
    toggleLoading();
    try {
      const {
        data: { code, body },
      } = await SavorApi.positionList();
      if (code === 0) {
        if (body && body.apps) {
          originShopList.value = body.apps;
          changeLangue(globalStore.langue);
        } else {
          shopList.value = [];
        }
      }
    } finally {
      toggleLoading();
    }
  };

  const searchVal = () => {
    router.push({
      name: 'homeList',
      query: {
        search: inputVal.value,
      },
    });
  };

  const showDetail = (data: any) => {
    router.push({
      name: 'homeDetail',
      query: {
        id: data.id,
      },
    });
  };

  const addressOption = ref<any>([]);
  const visible = ref(false);
  const initAddress = (newVal?: any) => {
    let proData: any = [];
    const add = pca['86'];
    const headerList: any = [];
    hotCity.value = originCity.value.map((item: any) => {
      return {
        cnLabel: item,
        label: newVal === 'en-US' ? pinyin(item, { toneType: 'none' }) : item,
      };
    });
    Object.keys(add).forEach((p: any) => {
      const children: any = [];
      Object.keys(pca[p]).forEach((c: any) => {
        const v: any = { value: c, label: pca[p][c], cnLabel: pca[p][c] };
        v.label =
          newVal === 'en-US' ? pinyin(v.label, { toneType: 'none' }) : v.label;
        children.push(v);
      });
      const proitem: any = { value: p, label: pca['86'][p] };
      if (
        [
          '北京市',
          '天津市',
          '上海市',
          '重庆市',
          '台湾省',
          '香港特别行政区',
          '澳门特别行政区',
        ].indexOf(proitem.cnLabel) < 0
      ) {
        if (children.length) {
          proData = proData.concat(children);
        }
      } else {
        const item = {
          ...proitem,
          label: proitem.label
            .replaceAll('市', '')
            .replaceAll('特别行政区', '')
            .replaceAll('省', ''),
          cnLabel: proitem.label
            .replaceAll('市', '')
            .replaceAll('特别行政区', '')
            .replaceAll('省', ''),
        };
        item.label =
          newVal === 'en-US'
            ? pinyin(item.label, { toneType: 'none' })
            : item.label;
        headerList.push(item);
      }
    });
    proData = headerList.concat(proData);
    addressOption.value = proData;
  };

  const showCityMore = () => {
    visible.value = !visible.value;
  };

  const searchData = async () => {
    toggleLoading();
    try {
      const params = {
        page: 1,
        page_size: 2000,
        query: inputVal.value,
        search: selectedCity.value,
      };
      const {
        data: { code, body },
      } = await SavorApi.searchShop(params);
      if (code === 0) {
        if (body && body.apps) {
          originSearchList.value = body.apps;
        } else {
          originSearchList.value = [];
        }
        changeLangueSecond(globalStore.langue);
      }
    } finally {
      toggleLoading();
    }
  };

  const selectCity = (val: any) => {
    visible.value = false;
    if (val === selectedCity.value) {
      selectedCity.value = '';
    } else {
      selectedCity.value = val;
    }
    if (selectedCity.value) {
      searchData();
    } else {
      fetchData();
    }
  };

  onMounted(() => {
    fetchData();
    initAddress(globalStore.langue);
    watch(
      () => globalStore.langue,
      (newVal) => {
        changeLangue(newVal);
        changeLangueSecond(newVal);
        initAddress(newVal);
      }
    );
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
      overflow-y: auto;
      padding-top: 70px;
      position: relative;

      .title {
        width: 359.26px;
        height: 193.03px;
        margin: 0 auto;

        img {
          width: 100%;
          height: 100%;
        }

        &.svg {
          background-image: url(../../assets/images/content.svg);
          background-size: 100% 100%;
          background-repeat: no-repeat;
          color: transparent;
        }
      }

      .search {
        width: 87%;
        margin: 60px auto 0;

        .search-icon {
          width: 30px;
          height: 30px;
          cursor: pointer;
        }

        .arco-input-wrapper {
          height: 36px;
          background: transparent;
          border: 1px solid #ffffff;
          border-radius: 50px;
          color: #fff;
          padding-right: 3px;
        }

        .arco-input-focus {
          border: 1px solid #ffffff;
        }

        :deep(.arco-icon-hover) {
          display: none;
        }
      }

      .city {
        width: 85%;
        color: #fff;
        display: flex;
        align-items: center;
        margin: 20px auto 0;

        .city-title {
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
              color: #8b95bc;
            }
            &.selected {
              border-bottom: 2px solid #fff;
            }
          }
        }
      }

      .city.mobile {
        flex-direction: column;
        align-items: flex-start;

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

      .banner {
        width: 87%;
        margin: 36px auto 0;

        .cnt {
          display: flex;
          justify-content: space-between;

          div {
            width: 30%;
            height: 212px;
          }

          img {
            width: 100%;
            height: 212px;
            object-fit: cover;
          }
        }
      }

      .shop-list {
        width: 87%;
        margin: 20px auto 0;
        display: flex;
        display: flex;
        flex-wrap: wrap;
        justify-content: space-between;
      }

      &.search {
        height: 100%;
        background-color: #20222c;
      }

      .shop.hidden {
        width: 300px;
        height: 120px;
        visibility: hidden;
      }

      :deep(.arco-spin) {
        margin-top: 60px;
      }
    }

    .content.mobile {
      padding-top: 30px;

      .title {
        width: 87%;
        margin-bottom: 16px;
      }

      .search {
        margin: 0px auto;
      }

      .banner {
        margin: 22px auto 0;
        padding: 0px 23px;
        width: 100%;

        .cnt {
          width: 100%;
          overflow-x: auto;
          overflow-y: hidden;
          display: flex;
          height: 144px;

          div {
            width: 233px;
            height: 144px;
            border-radius: 10px;
            margin-right: 8px;
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
        margin: 0px;

        .shop-con {
          width: 100%;
          padding: 0px 23px;

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
                white-space: nowrap;
                overflow: hidden;
                text-overflow: ellipsis;
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
  }
</style>

<style lang="less">
  .city-modal {
    width: 1000px;
    background: #272a37;
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
        color: #8b95bc;
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
    background-color: #272a37;
    border: unset;
  }
</style>
