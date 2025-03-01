<template>
  <div class="detail" :class="{ mobile: mobile }">
    <a-spin v-if="loading" :size="32" />
    <div v-else-if="!mobile" class="info">
      <a-button class="back form-btn" @click="handleBack"><icon-left style="margin-right: 5px;" />{{ $t('shop.back') }}</a-button>
      <div class="content">
        <a-carousel
          :style="{
            width: '100%',
            height: '394px',
            marginBottom: '20px'
          }"
          auto-play
        >
          <a-carousel-item v-for="image in data.pic_url" :key="image">
            <a-image
              :src="image"
              :style="{
                width: '100%',
                height: '100%',
                'object-fit': 'contain'
              }"
              :preview-props="{
                actionsLayout: [],
              }"
            ></a-image>
          </a-carousel-item>
        </a-carousel>
        <div class="name">{{ data.name }}</div>
        <div v-if="type === 'chef'" class="address">
          <div>{{ data.address }}</div>
        </div>
        <a-divider />
        <div class="describe">
          {{ data.describe }}
        </div>
      </div>
    </div>
    <div v-else class="con">
      <a-carousel
        :style="{
          width: '100%',
          height: '370px',
          marginBottom: '20px'
        }"
        show-arrow="never"
        :auto-play="true"
      >
        <a-carousel-item v-for="image in data.pic_url" :key="image">
          <a-image
            :src="image"
            :style="{
              width: '100%',
              height: '100%',
              'object-fit': 'contain'
            }"
            :preview-props="{
              actionsLayout: [],
            }"
          ></a-image>
        </a-carousel-item>
      </a-carousel>
      <div class="name-wrapper" :class="{ 'wrapper-desc' : !!data.introduction }">
        <div class="name">{{ data.name }}</div>
      </div>
      <div v-if="type === 'chef'" class="address">
        <div>{{ data.address }}</div>
      </div>
      <a-divider />
      <div class="describe">
        {{ data.describe }}
      </div> 
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { isMoblie } from '@/utils';
import { computed } from 'vue';

  const mobile = computed(() => isMoblie());
  const props = defineProps({
    loading: {
      type: Boolean,
      default() {
        return false;
      },
    },
    type: {
      type: String,
      default() {
        return 'chef';
      },
    },
    data: {
      type: Array as any,
      default() {
        return [];
      },
    },
  });
 
  const emit = defineEmits(['back']);

  const handleBack = (item: any) => {
    emit('back');
  }
</script>

<style lang="less" scoped>
 
  .detail {
    width: 60%;
    margin: 25px 5px 20px;
    box-sizing: border-box;
    color: #fff;
    font-family: PingFang SC;
    font-size: 12px;

    .info {
      width: 100%;
      padding-bottom: 50px;
      display: flex;
    }

    .content {
      flex: 1;
    }

    .back {
      margin-right: 30px;
    }

    .name, .address {
      text-align: center;
      font-family: PingFang SC;
      font-size: 18px;
      font-weight: 500;
      margin-bottom: 10px;
    }

    :deep(.arco-divider-horizontal) {
      border-color: #272A37;
    }

    .describe {
      white-space: pre-wrap;
      font-size: 16px;
      line-height: 24px;
    }

    .image {
      height: 100%;
      object-fit: contain;
    }
  }

  .mobile {
    width: 100%;
    padding: 0;
    margin: 0;

    .describe {
      padding: 0 16px 16px;
    }
  }

  :deep(.arco-image-img) {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
</style>
