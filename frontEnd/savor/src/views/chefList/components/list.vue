<template>
  <div class="list" :class="{ mobile: mobile }">
    <a-spin v-if="loading" :size="32" />
    <a-list v-else-if="!mobile" :bordered="false" :scrollbar="false">
      <a-list-item v-for="item of data" :key="item.id" class="list-item">
        <img class="img" :src="item.pic_url?.[0]" alt="">
        <div class="con">
          <div class="name">{{  item.name }}</div>
          <div v-if="type === 'chef'" class="address">
            {{ item.address }}
          </div>
          <div class="describe">
            <span>{{  item.describe }}</span>
          </div>
          <div class="extra">
            <div class="date">{{ item?.updated_at?.slice(0, 10).replaceAll('-', '.') }}</div>
            <a-button type="primary" class="form-btn" @click="handleDetail(item)">{{ $t('shop.detail') }}</a-button>
          </div>
        </div>
      </a-list-item>
    </a-list>
    <div v-else class="chef-con">
      <div v-for="item of data" :key="item.id" class="chef-item" @click="handleDetail(item)">
        <img class="img" :src="item?.pic_url?.[0]" alt="" >
        <div class="detail">
          <div class="name">{{  item.name }}</div>
          <div class="date">{{ item?.updated_at?.slice(0, 10).replaceAll('-', '.') }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { computed } from 'vue';
  import { isMoblie } from '@/utils';

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
      type: Object as any,
      default() {
        return [];
      },
    },
  });
 
  const emit = defineEmits(['showDetail']);

  const handleDetail = (item: any) => {
    emit('showDetail', item);
  }
</script>

<style lang="less" scoped>
 
  .list {
    width: 85%;
    margin: 25px 5px 20px;
    box-sizing: border-box;
    color: #fff;
    font-family: PingFang SC;
    font-size: 12px;
    // font-weight: 600;

    .list-item {
      height: 246px;
      border: unset!important;
      font-family: PingFang SC;
      font-size: 16px;
      line-height: 21.45px;
      text-align: left;
      border-radius: 16px;
      background: #272A37;
      padding: 0px!important;
      margin-bottom: 30px;
      box-shadow: 0px 3px 15px 0px #1C1E28;

      :deep(.arco-list-item-content) {
        display: flex;
        height: 100%;
      }


      .img {
        object-fit: cover;
        width: 246px;
        height: 246px;
        border-top-left-radius: 16px;
        border-bottom-left-radius: 16px;
      }

      .con {
        flex: 1;
        padding: 20px 30px 20px 40px;
        display: flex;
        flex-direction: column;
        color: #fff;

        .name {
          font-size: 18px;
          margin-bottom: 10px;
          font-weight: 600;
        }

        .address {
          font-weight: 600;
        }

        .describe {
          margin-top: 10px;
          flex: 1;
          color: #8B95BC;
          position: relative;
          overflow: hidden;

          span {
            display: -webkit-box;
            -webkit-line-clamp: 5;
            -webkit-box-orient: vertical;
            overflow: hidden;
            text-overflow: ellipsis;
          }
        }

        .extra {
          height: 30px;
          display: flex;
          justify-content: space-between;
          .date {
            color: #8E92BC;
          }

          .form-btn {
            background-color: #272A37!important;
            border-radius: 5px;
          }

        }
      }
    }

    .star {
      margin: 10px 0 5px;
    }

    .price {
      margin-bottom: 5px;
    }

  }

  .mobile {
    width: 100%;
    height: 100%;
    padding: 0px 24px;
    margin: 0px;

    .chef-con {
      width: 100%;
      padding-bottom: 23px;

      .chef-item {
        border-radius: 5px;
        margin-top: 22px;
        position: relative;

        .detail {
          position: absolute;
          width: 100%;
          height: 34px;
          bottom: 2px;
          left: 0px;
          right: 0px;
          background: #00000080;
          padding: 0px 14px;
          border-bottom-left-radius: 5px;
          border-bottom-right-radius: 5px;

          display: flex;
          align-items: center;
          justify-content: space-between;

          .name {
            flex: 1;
            font-size: 14px;
            font-weight: 500;
            color: #fff;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            padding-right: 16px;
          }
        }
      }

      .img {
        width: 100%;
        height: 300px;
        object-fit: cover;
        border-radius: 5px;
      }
    }
  }
</style>
