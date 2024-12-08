<template>
  <div class="container">
    <HeaderModile :back="true" @back="handleCancel" />
    <div class="detail">
      <a-spin v-if="loading" :size="32" />
      <div v-else class="content">
        <a-carousel
          :style="{
            width: '100%',
            height: '193px',
            marginBottom: '20px'
          }"
          show-arrow="never"
          :auto-play="true"
        >
          <a-carousel-item v-for="image in deatilData.pics_url" :key="image">
            <img
              class="image"
              :src="image"
              :style="{
                width: '100%',
              }"
            />
          </a-carousel-item>
        </a-carousel>
        <div class="name-wrapper" :class="{ 'wrapper-desc' : !!deatilData.introduction }">
          <div class="name">{{  deatilData.name }}</div>
          <div class="desc">{{  deatilData.introduction }}</div>
        </div>
        <div class="shopInfo">
          <div class="address">
            <img :src="addrImage" alt="">
            {{ deatilData.post_name || '' }}{{ deatilData.address }}
          </div>
          <div class="time">
            <img :src="timeImage" alt="">
            {{ deatilData.work_begin_at }} ~ {{ deatilData.work_end_at}}
          </div>
          <a-divider />
          <div class="describe">
            <div class="title">餐厅介绍</div>
            <div class="desc-content">
              <div class="desc" :class="{ all: !showIcon }">
                {{  descValue }}
                <icon-down class="iconDown" v-if="showIcon" @click="showMore" />
              </div>
            </div>
          </div>
          <div class="info">
            <div class="title">详细信息</div>
            <div class="info-detail">
              <div class="item">
                <div class="label">电话</div>
                <div class="desc">{{ deatilData.phone }}</div>
              </div>
              <div class="item">
                <div class="label">地址</div>
                <div class="desc">{{ deatilData.post_name || '' }}{{ deatilData.address }}</div>
              </div>
              <div class="item">
                <div class="label">营业时间</div>
                <div class="desc">{{ deatilData.work_begin_at }} ~ {{ deatilData.work_end_at}}</div>
              </div>
              <div class="item">
                <div class="label">是否有素食</div>
                <div class="desc">{{ deatilData.have_vege ? $t('settled.form.have_vege.option1') : $t('settled.form.have_vege.option2') }}</div>
            </div>
            </div>
          </div>
        </div> 
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import { useRouter } from 'vue-router';
  import useLoading from '@/hooks/loading';
  import HeaderModile from '@/components/header/mobile.vue';
  import SavorApi from '@/api/savor';
  import addrImage from '@/assets/images/addr.png';
  import timeImage from '@/assets/images/time.png';

  const router = useRouter();
  const { loading, toggle: toggleLoading } = useLoading(false);
  const deatilData = ref<any>({});
  const descValue = ref<any>('');
  const showIcon = ref<any>(false);

  const fetchData = async () => {
    toggleLoading();
    try {
      const { id } = router.currentRoute.value.query;
      const { data: { code, body } } = await SavorApi.shopDetail(id);
      if (code === 0) {
        if (body) {
          deatilData.value = body;
          if (body.describe.length > 110) {
            descValue.value = `${body.describe.slice(0, 110)}...`;
            showIcon.value = true;
          } else {
            descValue.value = body.describe;
            showIcon.value = false;
          }
        }
      }
    } finally {
      toggleLoading();
    }
  };

  const handleCancel = () => {
    const { search } = router.currentRoute.value.query;
    router.push({
      name: 'mobileList',
      query: {
        search
      },
    });
  }

  const showMore = () => {
    descValue.value = deatilData.value.describe;
    showIcon.value = false;
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

    .detail {
      padding: 27px 0px;
      position: relative;
      display: flex;
      flex-direction: column;
      align-items: center;
      flex: 1;
      overflow-y: auto;
      color: #fff;
      font-weight: 600;


      .name-wrapper {
        position: absolute;
        width: 100%;
        top: 176px;
        left: 0px;
        right: 0px;
        z-index: 9;
        background: linear-gradient(rgba(0, 0, 0, 0), rgba(0, 0, 0, 0.1), rgba(0, 0, 0, 0.5), #000000);
        padding: 0px 14px 12px;

        .name {
          margin-top: 10px;
          font-family: PingFang SC;
          font-size: 16px;
          text-align: left;
          line-height: 22px;
          color: #fff;
        }

        .desc {
          font-family: PingFang SC;
          font-size: 12px;
          font-weight: 400;
          line-height: 18px;
          color: #FFFFFFCC;
        }
      }

      .wrapper-desc {
        top: 158px;
      }

      .shopInfo {
        padding: 0 12px;
      }

      .address, .time {
        margin-bottom: 5px;
        line-height: 20px;
        display: flex;
        align-items: center;

        img {
          width: 13px;
          height: 13px;
          margin-right: 5px;
        }
      }

      .address {
        display: flex;
        align-items: flex-start;

        img {
          margin-top: 2px;
          height: 16px;
        }
      }

      .describe {
        font-weight: normal;
        line-height: 23px;

        .title {
          color: #D1B276;
          font-size: 16px;
          margin-bottom: 8px;
          font-weight: 500;
        }

        .desc-content {
          display: flex;
        }

        .desc {
          position: relative;
        }

        .iconDown {
          position: absolute;
          right: 12px;
          bottom: 5px;
          color: #D1B276;
        }

      }

      .info {
        margin-top: 24px;
        .title {
          color: #D1B276;
          font-size: 16px;
          margin-bottom: 8px;
          font-weight: 500;
        }

        .info-detail {
          background-color: #3A3E4D;
          padding: 20px 8px 40px 8px;

          .item {
            display: flex;
            padding: 20px 0;
            position: relative;

            .label {
              width: 110px;
            }

            .desc {
              flex: 1;
            }

            &::after {
              content: '';
              position: absolute;
              bottom: 0;
              width: 100%;
              height: 1px;
              background-color: #5555;
            }
          }
        }
      }

      :deep(.arco-divider-horizontal) {
        border: 0.5px solid #666666;
      }

      :deep(.arco-carousel-slide.arco-carousel-horizontal) {
        border-radius: 5px;
      }
    }
  }
</style>
