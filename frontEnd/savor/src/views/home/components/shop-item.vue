<template>
  <div class="shop" @click="showDetail">
    <div class="detail">
      <div class="title">{{ data.name }}</div>
      <div v-if="data.price" class="price">¥ {{ data.price }}/{{ $t('settled.detail.person') }}</div>
      <div class="address"><span>{{ data.post_name }}{{ data.address }}</span></div>
    </div>
    <img :src="data?.pics_url?.[0]" alt="">
  </div>
</template>

<script lang="ts" setup>
  const emit = defineEmits(['showDetail']);
  const props = defineProps({
    data: {
      type: Object as any,
      default() {
        return {};
      },
    },
  });

  const showDetail = () => {
    emit('showDetail', props.data);
  }
</script>

<style lang="less" scoped>
 
  .shop {
    width: 300px;
    height: 120px;
    margin: 5px 5px 20px;
    box-sizing: border-box;
    padding: 16px 24px;
    background: #3A3E4D;
    border-radius: 5px;
    display: flex;
    justify-content: space-between;
    color: #fff;
    font-family: PingFang SC;
    font-size: 12px;
    font-weight: 600;
    cursor: pointer;

    .detail {
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      width: 154px;
    }

    .title {
      font-family: PingFang SC;
      font-size: 16px;
      line-height: 21.45px;
      text-align: left;
      white-space: nowrap;
      -webkit-line-clamp: 1;
      -webkit-box-orient: vertical;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .star {
      margin: 10px 0 5px;
    }

    .price {
      margin-bottom: 5px;
    }

    img {
      flex-shrink: 0;
      width: 78px;
      height: 78px;
      margin-left: 20px;
    }

    .address {
      margin-bottom: 8px;
      span {
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
        text-overflow: ellipsis;
        line-height: 16px;
      }
    }
  }
</style>
