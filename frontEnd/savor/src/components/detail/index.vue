<template>
  <div v-if="!mobile" class="detail-pc">
    <a-button class="back form-btn" @click="handleCancel"
      ><icon-left style="margin-right: 5px" />{{ $t('shop.back') }}</a-button
    >
    <a-spin v-if="loading" :size="32" />
    <div v-else class="content">
      <a-carousel
        :style="{
          width: '100%',
          height: '394px',
          marginBottom: '20px',
        }"
        auto-play
      >
        <a-carousel-item v-for="image in data.pics_url" :key="image">
          <img
            class="image"
            :src="image"
            :style="{
              width: '100%',
            }"
          />
        </a-carousel-item>
      </a-carousel>
      <div class="detailInfo">
        <div>{{ data.name }}</div>
        <div class="address">
          <img :src="addrImage" alt="" />
          {{ data.post_name }}{{ data.address }}
        </div>
        <div class="shop">
          <div class="phone">
            <img :src="phoneImage" alt="" />
            {{ data.phone }}
          </div>
          <div class="time">
            <img :src="timeImage" alt="" />
            {{ data.work_begin_at }} ~ {{ data.work_end_at }}
          </div>
          <!-- <div class="diet">
            <img :src="dietImage" alt="" />
            {{
              data.have_vege
                ? $t('settled.form.have_vege.yes')
                : $t('settled.form.have_vege.no')
            }}
          </div> -->
        </div>
      </div>
      <a-divider />
      <div class="describe">
        {{ data.describe }}
      </div>
      <!-- <div class="extra">
        <a-carousel :style="{
          width: '100%',
          height: '399px',
          marginBottom: '20px',
        }" show-arrow="never" :auto-play="true">
          <a-carousel-item v-for="image in data.hostman_urls" :key="image">
            <a-image :src="image" :style="{
              'width': '100%',
              'height': '100%',
              'object-fit': 'contain',
            }" :preview-props="{
                actionsLayout: [],
              }"></a-image>
          </a-carousel-item>
        </a-carousel>
        <div class="hostman_name">{{ $t('settled.form.hostman') }}： {{ data.hostman_name }}</div>
        <div class="describe">{{ data.hostman_think }}</div>
      </div>
      <div v-for="(hot, index) of data.hot" :key="index">
        <div class="hostman_name">{{ $t('settled.form.hot') }}： {{ hot.hot_name }}</div>
        <a-carousel :style="{
          width: '100%',
          height: '377px',
          marginBottom: '20px',
        }" show-arrow="never" :auto-play="true">
          <a-carousel-item v-for="image in hot.hot_pic_urls" :key="image">
            <a-image :src="image" :style="{
              'width': '100%',
              'height': '100%',
              'object-fit': 'contain',
            }" :preview-props="{
                actionsLayout: [],
              }"></a-image>
          </a-carousel-item>
        </a-carousel>
        <div class="describe">{{ hot.hot_desc }}</div>
      </div> -->
      <div class="info">
        <div class="hostman_name">{{ $t('settled.form.detail_lable') }}</div>
        <div class="info-detail">
          <div class="item">
            <div class="label">{{ $t('settled.form.phone_lable') }}</div>
            <div class="desc">{{ data.phone }}</div>
          </div>
          <!-- <div class="item">
            <div class="label">{{ $t('settled.form.yuyue') }}</div>
            <div class="desc">
              {{
                Number(data.yuyue) === 1
                  ? $t('settled.form.yuyue.option1')
                  : $t('settled.form.yuyue.option2')
              }}
            </div>
          </div> -->
          <div class="item">
            <div class="label">{{ $t('settled.form.address') }}</div>
            <div class="desc">{{ data.post_name || '' }}{{ data.address }}</div>
          </div>
          <div class="item">
            <div class="label">{{ $t('settled.form.time') }}</div>
            <div class="desc"
              >{{ data.work_begin_at }} ~ {{ data.work_end_at }}</div
            >
          </div>
          <div class="item">
            <div class="label">{{ $t('settled.form.open_lable') }}</div>
            <div class="desc">
              {{ data.work_date_label && data.work_date_label.join('、') }}</div
            >
          </div>
          <div class="item">
            <div class="label">{{ $t('settled.form.per') }}</div>
            <div class="desc">
              {{ data.per }}{{ $t('settled.form.per.desc') }}</div
            >
          </div>
          <!-- <div class="item">
            <div class="label">{{ $t('settled.form.have_vege') }}</div>
            <div class="desc">{{
              data.have_vege
                ? $t('settled.form.have_vege.option1')
                : $t('settled.form.have_vege.option2')
            }}</div>
          </div> -->
        </div>
      </div>
    </div>
  </div>
  <Mobile v-else :detail-data="data" />
</template>

<script lang="ts" setup>
import addrImage from '@/assets/images/addr.png';
import phoneImage from '@/assets/images/phone.png';
import timeImage from '@/assets/images/time.png';
import { isMoblie } from '@/utils';
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import Mobile from './mobile.vue';

const mobile = computed(() => isMoblie());
defineProps({
  loading: {
    type: Boolean,
    default() {
      return false;
    },
  },
  data: {
    type: Object as any,
    default() {
      return { items: [] };
    },
  },
});
const router = useRouter();

defineEmits(['showDetail', 'cancel']);

const handleCancel = () => {
  router.back();
};
</script>

<style lang="less" scoped>
.detail-pc {
  width: 65%;
  margin: 20px 5px;
  box-sizing: border-box;
  color: #8b95bc;
  font-family: PingFang SC;
  font-size: 12px;
  display: flex;

  .content {
    flex: 1;
    margin-left: 20px;
    border: unset;
    font-family: PingFang SC;
    font-size: 16px;
    line-height: 21.45px;
    text-align: left;
    border-radius: 16px;
    padding: 0px !important;
    margin-bottom: 25px;

    .detailInfo {
      color: #fff;

      .shop {
        display: flex;

        div {
          margin-right: 16px;

          &:first-child {
            img {
              height: 18px;
            }
          }
        }
      }

      .address,
      .phone,
      .time,
      .diet {
        display: flex;
        align-items: center;
      }

      .address {
        margin: 10px 0;

        img {
          height: 16px;
          margin-right: 5px;
        }
      }

      img {
        width: 13px;
        height: 13px;
        margin-right: 5px;
      }
    }

    :deep(.arco-divider-horizontal) {
      border-color: #272a37;
    }

    .image {
      height: 100%;
      object-fit: contain;
    }

    .describe {
      line-height: 24px;
      margin-bottom: 20px;
    }

    :deep(.arco-image-img) {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .hostman_name {
      font-size: 18px;
      color: #caa97f;
      margin-bottom: 10px;
    }

    .info {
      margin-top: 24px;
      color: #fff;

      .info-detail {
        background-color: #3a3e4d;
        padding: 12px 12px 6px 12px;
        border-radius: 5px;

        .item {
          font-size: 14px;
          display: flex;
          padding: 10px 0;
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
  }
}
</style>