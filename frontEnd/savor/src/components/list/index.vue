<template>
  <div class="list">
    <a-spin v-if="loading" :size="32" />
    <a-list v-else :bordered="false" :scrollbar="false">
      <template #scroll-loading>
        <!-- <div v-if="bottom"></div>
        <a-spin v-else /> -->
      </template>
      <a-list-item v-for="item of data" :key="item.id" class="list-item">
        <img class="img" :src="item?.pics_url?.[0]" alt="">
        <div class="con">
          <div class="name">{{  item.name }}</div>
          <div class="shop">
            <div class="address">
              <img :src="addrImage" alt="">
              <span>{{ item.post_name }}{{ item.address }}</span>
            </div>
            <div class="phone">
              <img :src="phoneImage" alt="">
              {{ item.phone  }}
            </div>
            <div class="time">
              <img :src="timeImage" alt="">
              {{ item.work_begin_at }} ~ {{ item.work_end_at}}
            </div>
            <div class="diet">
              <img :src="dietImage" alt="">
              {{ item.have_vege ? $t('settled.form.have_vege.yes') : $t('settled.form.have_vege.no') }}
            </div>
          </div>
          <div class="describe">
            <span>
              {{ item.describe }}
            </span>
          </div>
          <div class="extra">
            <div class="date">{{ item?.updated_at?.slice(0, 10).replaceAll('-', '.') }}</div>
            <a-button type="primary" class="form-btn" @click="handleDetail(item)">{{ $t('shop.detail') }}</a-button>
          </div>
        </div>
      </a-list-item>
    </a-list>
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
    bottom: {
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
  const emit = defineEmits(['showDetail']);

  const handleDetail = (item: any) => {
    emit('showDetail', item);
  };
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
      border: unset;
      font-family: PingFang SC;
      font-size: 16px;
      line-height: 21.45px;
      text-align: left;
      border-radius: 16px;
      background: #272a37;
      padding: 0px!important;
      margin-bottom: 25px;
      box-shadow: 0px 3px 15px 0px #1c1e28;

      :deep(.arco-list-item-content) {
        display: flex;
      }


      .img {
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
        }

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
          img {
            width: 13px;
            height: 13px;
            margin-right: 2px;
          }
        }

        .describe {
          margin-top: 10px;
          flex: 1;
          color: #8B95BC;
          span {
            display: -webkit-box;
            -webkit-line-clamp: 5;
            -webkit-box-orient: vertical;
            overflow: hidden;
            text-overflow: ellipsis;
          }
        }

        .address {
          flex: 1;
          display: flex;
          span {
            margin-left: 5px;
            flex: 1;
            display: -webkit-box;
            -webkit-line-clamp: 1;
            -webkit-box-orient: vertical;
            overflow: hidden;
            text-overflow: ellipsis;
          }
        }

        .phone {
          width: 180px;
        }

        .time {
          width: 150px;
        }

        .diet {
          width: 130px;
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
  }
</style>
