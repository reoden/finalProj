<template>
  <div class="detail">
    <div class="content">
      <a-carousel
        :style="{
          width: '100%',
          height: '200px',
          marginBottom: '10px',
        }"
        show-arrow="never"
        :auto-play="true"
      >
        <a-carousel-item v-for="image in detailData.pics_url" :key="image">
          <a-image
            :src="image"
            :style="{
              'width': '100%',
              'height': '100%',
              'object-fit': 'contain',
            }"
            :preview-props="{
              actionsLayout: [],
            }"
          ></a-image>
        </a-carousel-item>
      </a-carousel>
      <div
        class="name-wrapper"
        :class="{ 'wrapper-desc': !!detailData.introduction }"
      >
        <div class="name">{{ detailData.name }}</div>
        <div class="desc">{{ detailData.introduction }}</div>
      </div>
      <div class="shopInfo">
        <a-button class="btn">{{ detailData.leixing }}</a-button>
        <div class="address">
          <icon-location />
          <span style="flex: 1"
            >{{ detailData.post_name || '' }}{{ detailData.address }}</span
          >
        </div>
        <div class="time">
          <icon-clock-circle />
          <span style="flex: 1">
            {{
              detailData.work_date_label &&
              detailData.work_date_label.join('、')
            }}
            {{ detailData.work_begin_at }} ~ {{ detailData.work_end_at }}
          </span>
        </div>
        <a-divider />
        <div class="describe">
          <div class="desc-content">
            <div class="desc" :class="{ all: !showIcon }">
              {{ descValue }}
              <icon-down v-if="showIcon" class="iconDown" @click="showMore" />
            </div>
          </div>
        </div>
        <a-divider />
        <div>
          <div class="hostman_name"
            >{{ $t('settled.form.hostman') }}：
            {{ detailData.hostman_name }}</div
          >
          <a-carousel
            :style="{
              width: '100%',
              height: '380px',
              marginBottom: '20px',
            }"
            show-arrow="never"
            :auto-play="true"
          >
            <a-carousel-item
              v-for="image in detailData.hostman_urls"
              :key="image"
            >
              <a-image
                :src="image"
                :style="{
                  'width': '100%',
                  'height': '100%',
                  'object-fit': 'contain',
                }"
                :preview-props="{
                  actionsLayout: [],
                }"
              ></a-image>
            </a-carousel-item>
          </a-carousel>
          <div>
            <div class="hostman_name think">{{
              $t('settled.form.hostman_think')
            }}</div>
            <div class="hostman_think">{{ detailData.hostman_think }}</div>
          </div>
        </div>
        <div v-for="(hot, index) of detailData.hot" :key="index">
          <div class="hostman_name"
            >{{ $t('settled.form.hot') }}： {{ hot.hot_name }}</div
          >
          <a-carousel
            :style="{
              width: '100%',
              height: '330px',
              marginBottom: '20px',
            }"
            show-arrow="never"
            :auto-play="true"
          >
            <a-carousel-item v-for="image in hot.hot_pic_urls" :key="image">
              <a-image
                :src="image"
                :style="{
                  'width': '100%',
                  'height': '100%',
                  'object-fit': 'contain',
                }"
                :preview-props="{
                  actionsLayout: [],
                }"
              ></a-image>
            </a-carousel-item>
          </a-carousel>
          <div class="hostman_think">{{ hot.hot_desc }}</div>
        </div>
        <div class="info">
          <div class="title">{{ $t('settled.form.detail_lable') }}</div>
          <div class="info-detail">
            <div class="title">{{ $t('settled.form.detail_lable_1') }}</div>
            <div class="item">
              <div class="label">{{ $t('settled.form.phone_lable') }}</div>
              <div class="desc">{{ detailData.phone }}</div>
            </div>
            <div class="item">
              <div class="label">{{ $t('settled.form.yuyue') }}</div>
              <div class="desc">
                {{
                  Number(detailData.yuyue) === 1
                    ? $t('settled.form.yuyue.option1')
                    : $t('settled.form.yuyue.option2')
                }}
              </div>
            </div>
            <div class="item">
              <div class="label">{{ $t('settled.form.address') }}</div>
              <div class="desc"
                >{{ detailData.post_name || '' }}{{ detailData.address }}</div
              >
            </div>
            <div class="item">
              <div class="label">{{ $t('settled.form.time') }}</div>
              <div class="desc"
                >{{ detailData.work_begin_at }} ~
                {{ detailData.work_end_at }}</div
              >
            </div>
            <div class="item">
              <div class="label">{{ $t('settled.form.open_lable') }}</div>
              <div class="desc">
                {{
                  detailData.work_date_label &&
                  detailData.work_date_label.join('、')
                }}</div
              >
            </div>
            <div class="item">
              <div class="label">{{ $t('settled.form.per') }}</div>
              <div class="desc">
                {{ detailData.per }}{{ $t('settled.form.per.desc') }}</div
              >
            </div>
            <div class="item">
              <div class="label">{{ $t('settled.form.have_vege') }}</div>
              <div class="desc">{{
                detailData.have_vege
                  ? $t('settled.form.have_vege.option1')
                  : $t('settled.form.have_vege.option2')
              }}</div>
            </div>
          </div>
          <div class="logo" :class="{ svg: !userInfo.logoBig }">
            <img
              v-if="userInfo.logoBig"
              :src="userInfo.logoBig"
              alt=""
              srcset=""
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted, watch } from 'vue';
  import { useUserStore } from '@/store';

  const showIcon = ref<any>(false);
  const userInfo = useUserStore();

  const props = defineProps({
    detailData: {
      type: Object as any,
      default() {
        return {};
      },
    },
  });
  const descValue = ref<any>('');

  const showMore = () => {
    descValue.value = props.detailData.describe;
    showIcon.value = false;
  };

  const checkDesc = () => {
    if (props?.detailData?.describe?.length > 110) {
      descValue.value = `${props.detailData.describe.slice(0, 110)}...`;
      showIcon.value = true;
    } else {
      descValue.value = props.detailData.describe;
      showIcon.value = false;
    }
  };

  watch(
    () => props.detailData,
    () => {
      checkDesc();
    }
  );

  onMounted(() => {
    checkDesc();
  });
</script>

<style lang="less" scoped>
  .detail {
    width: 100%;
    padding: 0px 0px 27px;
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    flex: 1;
    overflow-y: auto;
    color: #fff;
    font-weight: 600;

    .content {
      width: 100%;
    }

    .name-wrapper {
      position: absolute;
      width: 100%;
      top: 156px;
      left: 0px;
      right: 0px;
      z-index: 9;
      background: linear-gradient(
        rgba(0, 0, 0, 0),
        rgba(0, 0, 0, 0.1),
        rgba(0, 0, 0, 0.5),
        #000000
      );
      padding: 0px 12px 12px;

      .name {
        margin-top: 10px;
        font-family: PingFang SC;
        font-size: 16px;
        text-align: left;
        line-height: 22px;
        color: #fff;
        white-space: nowrap;
        text-overflow: ellipsis;
        overflow: hidden;
      }

      .desc {
        font-family: PingFang SC;
        font-size: 12px;
        font-weight: 400;
        line-height: 18px;
        color: #ffffffcc;
        white-space: nowrap;
        text-overflow: ellipsis;
        overflow: hidden;
      }
    }

    .wrapper-desc {
      top: 138px;
    }

    .shopInfo {
      padding: 0;
    }

    .address,
    .time {
      margin-bottom: 5px;
      line-height: 20px;
      display: flex;
      align-items: center;
      padding: 0 12px;

      :deep(.arco-icon) {
        font-size: 16px;
        color: #d1b276;
        margin-right: 5px;
      }
    }

    .btn {
      margin: 0 12px 12px;
      border-radius: 16px;
      background: #d1b276;
      color: #20222c;
    }

    .describe {
      font-weight: normal;
      line-height: 30px;
      padding: 5px 12px;

      .title {
        color: #d1b276;
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
        color: #d1b276;
      }
    }

    .hostman_name {
      margin: 12px 0 20px 12px;
      font-size: 20px;
      color: #d1b276;

      &.think {
        margin-bottom: 13px;
      }
    }

    .hostman_think {
      margin: 12px;
      line-height: 30px;
      font-weight: normal;
    }

    .logo {
      margin: 20px auto 0;

      img {
        width: 100%;
        height: 100%;
      }
    }

    .info {
      margin-top: 24px;
      .title {
        color: #d1b276;
        font-size: 18px;
        margin-bottom: 8px;
        font-weight: 500;
        padding-left: 12px;
      }

      .info-detail {
        margin: 0 12px;
        background-color: #3a3e4d;
        padding: 20px 12px 6px 12px;
        border-radius: 2px;

        .title {
          font-size: 14px;
          padding-left: 0px;
        }

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

          &:last-child {
            &::after {
              display: none;
            }
          }
        }
      }
    }

    :deep(.arco-image-img) {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    :deep(.arco-divider-horizontal) {
      width: calc(100% - 24px);
      border: 0.5px solid #666666;
    }

    :deep(.arco-carousel-slide.arco-carousel-horizontal) {
      border-radius: 0px;
    }
  }
</style>
