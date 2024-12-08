<template>
  <div class="detail">
    <a-button class="back form-btn" @click="handleCancel"><icon-left style="margin-right: 5px;" />{{ $t('shop.back') }}</a-button>
    <a-spin v-if="loading" :size="32" />
    <div v-else class="content">
      <a-carousel
        :style="{
          width: '100%',
          height: '394px',
          marginBottom: '20px'
        }"
        auto-play
      >
        <a-carousel-item v-for="image in data.pictures" :key="image">
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
          <img :src="addrImage" alt="">
          {{ data.address }}
        </div>
        <div class="shop">
          <div class="phone">
            <img :src="phoneImage" alt="">
            {{ data.phone }}
          </div>
          <div class="time">
            <img :src="timeImage" alt="">
            {{ data.work_begin_at }} ~ {{ data.work_end_at}}
          </div>
          <div class="diet">
            <img :src="dietImage" alt="">
            {{ data.have_vege ? $t('settled.form.have_vege.no') : $t('settled.form.have_vege.yes') }}
          </div>
        </div>
      </div>
      <a-divider />
      <div class="describe">
        {{ data.describe }}
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import addrImage from '@/assets/images/addr.png';
  import phoneImage from '@/assets/images/phone.png';
  import timeImage from '@/assets/images/time.png';
  import dietImage from '@/assets/images/diet.png';

  const props = defineProps({
    loading: {
      type: Boolean,
      default() {
        return false;
      },
    },
    data: {
      type: Object as any,
      default() {
        return [];
      },
    },
  });
 
  const emit = defineEmits(['fetchData', 'showDetail', 'cancel']);

  const fetchData = () => {
    emit('fetchData');
  }

  const handleCancel = () => {
    emit('cancel');
  }
</script>

<style lang="less" scoped>
 
  .detail {
    width: 65%;
    margin: 65px 5px 20px;
    box-sizing: border-box;
    color: #fff;
    font-family: PingFang SC;
    font-size: 12px;
    display: flex;

    .content {
      margin-left: 20px;
      border: unset;
      font-family: PingFang SC;
      font-size: 16px;
      line-height: 21.45px;
      text-align: left;
      border-radius: 16px;
      padding: 0px!important;
      margin-bottom: 25px;

      .detailInfo {
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
        .address {
          margin: 10px 0;
        }
        img {
          width: 13px;
          height: 13px;
          margin-right: 2px;
          vertical-align: middle;
        }
      }

      :deep(.arco-divider-horizontal) {
        border-color:  #272A37;
      }
      .image {
        height: 100%;
        object-fit: contain;
      }
    }
  }
</style>
