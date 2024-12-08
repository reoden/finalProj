<template>
  <div class="container">
    <Header :back="true" />
    <div class="content search">
      <Detail
        :loading="loading"
        :data="detailData"
        @cancel="handleCancel"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted, watch } from 'vue';
  import { useI18n } from 'vue-i18n';
  import { useRouter } from 'vue-router';
  import { useAppStore } from '@/store';
  import useLoading from '@/hooks/loading';
  import Header from '@/components/header/index.vue';
  import SavorApi from '@/api/savor';
  import Detail from '@/components/detail/index.vue';

  const router = useRouter();
  const { loading, toggle: toggleLoading } = useLoading(false);
  const originDetail = ref<any>({});
  const detailData = ref<any>({});
  const { t } = useI18n();
  const globalStore = useAppStore();

  const dateArr = [{
    label: t('settled.form.label_mon'),
    value: t('settled.form.label_mon_us')
  }, {
    label: t('settled.form.label_tue'),
    value: t('settled.form.label_tue_us')
  }, {
    label: t('settled.form.label_wed'),
    value: t('settled.form.label_wed_us')
  }, {
    label: t('settled.form.label_thur'),
    value: t('settled.form.label_thur_us')
  }, {
    label: t('settled.form.label_fir'),
    value: t('settled.form.label_fir_us')
  }, {
    label: t('settled.form.label_satu'),
    value: t('settled.form.label_satu_us')
  }, {
    label: t('settled.form.label_sun'),
    value: t('settled.form.label_sun_us')
  }];

  const changeLangue = (newVal: any) => {
    if (newVal === 'zh-CN') {
      detailData.value = originDetail.value;
    } else {
      detailData.value = {
        ...originDetail.value,
        ...originDetail.value.en
      };
    }
  }

  const fetchData = async () => {
    toggleLoading();
    try {
      const { id } = router.currentRoute.value.query;
      const { data: { code, body } } = await SavorApi.shopDetail(id);
      if (code === 0) {
        if (body) {
          originDetail.value = {
            ...body,
            work_date_label: body.work_date ? body.work_date.map((item: any) => {
              const map: any = {};
              dateArr.forEach((date: any) => {
                map[date.value] = date.label;
              })
              return map[item];
            }) : [],
          };
          changeLangue(globalStore.langue);
        }
      }
    } finally {
      toggleLoading();
    }
  };

  const handleCancel = () => {
    router.push({
      name: 'homeList'
    });
  }

  onMounted(() => {
    fetchData();
    watch(() => globalStore.langue, (newVal) => {
      changeLangue(newVal);
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
      position: relative;
      display: flex;
      flex-direction: column;
      align-items: center;
      flex: 1;
      overflow-y: auto;

      .title {
        font-size: 60px;
        font-family: PingFang SC;
        line-height: 128.7px;
        text-align: center;
        color: #FFFFFF;
      }

      .search {
        margin-top: 60px;
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

    }
  }
</style>
