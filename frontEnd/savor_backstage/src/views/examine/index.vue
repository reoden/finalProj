<template>
  <div class="container">
    <MenuBar />
    <div class="content">
      <div class="title">
        <span @click="showList">{{ $t('menu.examine.title') }}</span
        >{{ type === 'detail' ? ` > ${detailInfo.name}` : '' }}
      </div>
      <a-divider class="divider" />
      <div v-if="type === 'list'" class="table">
        <a-table
          :loading="loading"
          :columns="columns"
          :data="shopList"
          row-key="id"
          :border="false"
          :pagination="{
            total: pagination.total,
            current: pagination.page,
            pageSize: pagination.size,
          }"
          @page-change="changePage"
        >
          <template #status="{ record }">
            {{ STATUS_MAP[record.status || 0] }}
          </template>
          <template #address="{ record }">
            {{ record?.post_name }} {{ record.address }}
          </template>
          <template #operate="{ record }">
            <a-link class="link" :hoverable="false" @click="detail(record)"
              >查看详情</a-link
            >
          </template>
        </a-table>
      </div>
      <div v-if="type === 'detail'" class="detail">
        <a-spin v-if="detailLoading" :size="32" />
        <template v-else>
          <div class="left">
            <Detail :data="detailInfo" />
          </div>
          <div class="right">
            <div class="btns">
              <a-button
                v-if="detailInfo.status !== 2"
                class="agree"
                :loading="agreeLoading"
                @click="agree"
                >通过</a-button
              >
              <a-button
                v-if="detailInfo.status !== 2"
                class="refuse"
                :loading="refuseLoading"
                @click="refuse"
                >拒绝</a-button
              >
            </div>
            <a-button class="form-btn down" :loading="downLoading" @click="down"
              >下载商家资料</a-button
            >
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import SavorApi from '@/api/savor';
import MenuBar from '@/components/menu-bar/index.vue';
import { STATUS_MAP } from '@/config';
import useLoading from '@/hooks/loading';
import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import Detail from './detail.vue';

  const { t } = useI18n();
  const shopList = ref<any>([]);
  const { loading, toggle: toggleLoading } = useLoading(false);
  const { loading: detailLoading, toggle: toggleDetailLoading } =
    useLoading(false);
  const pagination = ref<any>({
    page: 1,
    size: 20,
    total: 0
  });
  const type = ref<any>('list');
  const detailInfo = ref<any>({});
  const dateArr = [
    {
      label: t('settled.form.label_mon'),
      value: t('settled.form.label_mon_us'),
    },
    {
      label: t('settled.form.label_tue'),
      value: t('settled.form.label_tue_us'),
    },
    {
      label: t('settled.form.label_wed'),
      value: t('settled.form.label_wed_us'),
    },
    {
      label: t('settled.form.label_thur'),
      value: t('settled.form.label_thur_us'),
    },
    {
      label: t('settled.form.label_fir'),
      value: t('settled.form.label_fir_us'),
    },
    {
      label: t('settled.form.label_satu'),
      value: t('settled.form.label_satu_us'),
    },
    {
      label: t('settled.form.label_sun'),
      value: t('settled.form.label_sun_us'),
    },
  ];

  const columns = computed(() => {
    return [
      {
        title: t('examine.colunm.id'),
        dataIndex: 'id',
        width: 100,
      },
      {
        title: t('examine.colunm.name'),
        dataIndex: 'name',
        width: 120,
      },
      {
        title: t('examine.colunm.address'),
        slotName: 'address',
        dataIndex: 'address',
        ellipsis: true,
      },
      {
        dataIndex: 'describe',
        title: t('examine.colunm.describe'),
        ellipsis: true,
      },
      {
        dataIndex: 'status',
        slotName: 'status',
        width: 100,
        title: t('examine.colunm.status'),
      },
      {
        dataIndex: 'operate',
        slotName: 'operate',
        width: 90,
        title: t('examine.colunm.operate'),
      },
    ];
  });

  const detail = async (record: any) => {
    type.value = 'detail';
    toggleDetailLoading();
    try {
      const {
        data: { code, body },
      } = await SavorApi.shopDetail(record.id);
      if (code === 0) {
        if (body) {
          detailInfo.value = {
            ...body,
            work_date_label: body.work_date
              ? body.work_date.map((item: any) => {
                  const map: any = {};
                  dateArr.forEach((date: any) => {
                    map[date.value] = date.label;
                  });
                  return map[item];
                })
              : [],
          };
        }
      }
    } finally {
      toggleDetailLoading();
    }
  };

  const showList = () => {
    if (type.value === 'list') {
      return;
    }
    type.value = 'list';
    detailInfo.value = {};
  };

  // 列表数据
  const fetchData = async () => {
    toggleLoading();
    try {
      const params = {
        page: pagination.value.page,
        size: pagination.value.size,
      };
      const {
        data: { code, body },
      } = await SavorApi.shopList(params);
      if (code === 0) {
        pagination.value.total = body?.total;
        if (body && body.apps) {
          shopList.value = body.apps;
        }
      }
    } finally {
      toggleLoading();
    }
  };

  const changePage = (size: any) => {
    pagination.value.page = size;
    fetchData();
  };

  const { loading: agreeLoading, toggle: toggleAgreeLoading } =
    useLoading(false);
  const { loading: refuseLoading, toggle: toggleRefuseLoading } =
    useLoading(false);
  const { loading: downLoading, toggle: toggleDownLoading } = useLoading(false);

  const agree = async () => {
    toggleAgreeLoading();
    try {
      const {
        data: { code },
      } = await SavorApi.examineShop({
        id: detailInfo.value.id,
        status: '通过',
      });
      if (code === 0) {
        showList();
      }
    } finally {
      toggleAgreeLoading();
    }
  };

  const refuse = async () => {
    toggleRefuseLoading();
    try {
      const {
        data: { code },
      } = await SavorApi.examineShop({
        id: detailInfo.value.id,
        status: '拒绝',
      });
      if (code === 0) {
        showList();
      }
    } finally {
      toggleRefuseLoading();
    }
  };

  const down = async () => {
    toggleDownLoading();
    try {
      const {
        data: { code, body },
      } = await SavorApi.downFile(detailInfo.value.id);
      if (code === 0 && body) {
        const a = document.createElement('a');
        a.href = body;
        a.style.display = 'none';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
      }
    } finally {
      toggleDownLoading();
    }
  };

  onMounted(() => {
    fetchData();
  });
</script>

<style lang="less" scoped>
  .container {
    height: 100%;
    width: 100%;
    display: flex;
    overflow: hidden;

    .content {
      flex: 1;
      position: relative;
      display: flex;
      flex-direction: column;

      .title {
        font-size: 18px;
        font-weight: 600;
        font-family: PingFang SC;
        color: #ffffff;
        padding: 30px 48px;
        span {
          cursor: pointer;
        }
      }

      .divider {
        width: 100%;
        margin: 0px;
        border-color: #4a5069;
      }

      .table {
        margin: 20px;
        padding: 4px 48px;

        :deep(.arco-table) {
          .arco-table-container {
            border: unset;
          }

          .arco-table-th {
            background-color: transparent;
            border: unset;
            font-size: 16px;
            font-weight: 600;
            color: #fff;
          }

          .arco-table-td {
            background-color: transparent;
            border-color: #4a5069;
            color: #fff;
          }

          .arco-table-tr {
            &:hover {
              .arco-table-td {
                background-color: transparent !important;
              }
            }
          }

          .link {
            color: #d2b276;
            text-decoration: underline;
          }
        }
      }

      .detail {
        flex: 1;
        display: flex;
        overflow: hidden;

        .left {
          flex: 1;
          padding: 20px;
          overflow-y: auto;
        }

        .right {
          width: 30%;
          height: 100%;
          background-color: #20222c;
          display: flex;
          flex-direction: column;
          align-items: center;
          padding-top: 50px;
          justify-content: space-between;

          .btns {
            display: flex;
            flex-direction: column;
            align-items: center;
          }

          :deep(.arco-btn) {
            width: 200px;
            color: #fff;

            &.agree {
              background: #d1b276;
              margin-bottom: 20px;
            }
            &.refuse {
              background: #f56c6c;
            }
            &.down {
              margin-bottom: 50px;
            }
          }
        }
      }
    }
  }
</style>
