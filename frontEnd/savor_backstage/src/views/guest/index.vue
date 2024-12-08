<template>
  <div class="container">
    <MenuBar />
    <div class="content">
      <div class="header">
        <div class="menu selected">
          {{ $t('menu.guest.title') }}
        </div>
        <a-button @click="addGuest">增加嘉宾</a-button>
      </div>
      <div class="nameEdit">
        <div>菜单名修改</div>
        <a-input class="input" v-model="guestName" :disabled="!isEdit"></a-input>
        <a-button @click="editName" :loading="editLoading">{{ isEdit ? '保存' : '编辑' }}</a-button>
      </div>
      <div class="table">
        <a-table
          :loading="loading"
          :columns="columns"
          :data="guestList"
          row-key="id"
          :border="false"
        >
          <template #index="{ rowIndex }">
            {{ rowIndex + 1 }}
          </template>
          <template #operate="{ record }">
            <icon-edit class="btn-del" title="编辑" @click="edit(record)" />
            <icon-delete class="btn-del" title="删除" @click="del(record)" />
            <img
              class="handle"
              :src="iconDropImage"
              style="width: 15px; height: 15px; object-fit: contain" />
          </template>
        </a-table>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted, computed } from 'vue';
  import { useI18n } from 'vue-i18n';
  import Sortable from 'sortablejs'
  import useLoading from '@/hooks/loading';
  import SavorApi from '@/api/savor';
  import MenuBar from '@/components/menu-bar/index.vue';
  import iconDropImage from '@/assets/images/icon-drop.png';
  import { Message, Modal } from '@arco-design/web-vue';
  import { useRouter } from 'vue-router';

  const router = useRouter();
  const { t } = useI18n();
  const guestList = ref<any>([]);
  const { loading, toggle: toggleLoading } = useLoading(false);
  const { loading: editLoading, toggle: toggleEditLoading } = useLoading(false);
  const guestName = ref('嘉宾顾问');
  const isEdit = ref(false);

  const columns = computed(() => {
    return [
    {
        title: t('position.colunm.id'),
        slotName: 'index',
        dataIndex: 'index',
        width: 120
      },
      {
        title: t('guest.colunm.name'),
        dataIndex: 'name',
        width: 150,
      },
      {
        title: t('guest.colunm.describe'),
        dataIndex: 'describe',
        ellipsis: true
      },
      {
        dataIndex: 'operate',
        slotName: 'operate',
        width: 160,
        title: t('position.colunm.operate'),
      },
    ];
  });

  const returnDrop = (evt: any) => {
    // 拖拽失败时的处理函数
    const item: any = evt?.item; // 被拖拽的元素
    const startIndex: any = evt?.oldIndex; // 开始时的索引
    if (startIndex !== null) {
      // 将元素移动到原来的位置
      item.parentNode.insertBefore(item.parentNode.children[startIndex], item);
    }
  };

  // 获取列表数据
  const fetchData = async () => {
    toggleLoading();
    try {
      const { data: { code, body } } = await SavorApi.guestList();
      if (code === 0) {
        guestList.value = body?.apps;
        const tbody: any = document.querySelector('.table .arco-table tbody');
        const drag = new Sortable(tbody, {
          handle: '.handle',
          animation: 150,
          // 需要在odEnd方法中处理原始eltable数据, 使原始数据与显示数据保持顺序一致
          onEnd: async (evt) => {
            const { newIndex, oldIndex } : any = evt;
            // 拖拽后接口保存排序后的位置
            const tableData = JSON.parse(JSON.stringify(guestList.value));
            const targetRow = tableData[oldIndex];
            tableData.splice(oldIndex, 1)
            tableData.splice (newIndex, 0, targetRow);
            guestList.value = tableData;
            try {
              const appIds = tableData.map((item: any) => item.id);
              const { data: { code: codeSec  } } = await SavorApi.changeGuest({ app_ids: appIds });
              if (codeSec !== 0) {
                // 保存失败后，撤回本次拖拽
                returnDrop(evt);
              }
            } catch (error) {
              returnDrop(evt);
            }
          }
        });
      }
    } finally {
      toggleLoading();
    }
  };

  const addGuest = () => {
    router.push({
      name: 'guestBackDetail',
    });
  }

  const edit = (record: any) => {
    router.push({
      name: 'guestBackDetail',
      query: {
        id: record.id
      },
    });
  }

  const del = (row?: any) => {
    Modal.confirm({
      content: `确认删除嘉宾 ${row.name} 吗?`,
      okText: '确定',
      async onOk() {
        const { data: { code } } = await SavorApi.deleteGuest(row.id);
        if (code === 0) {
          Message.success('删除成功');
          fetchData();
        }
      },
    });
  }

  const fetchSwitch = async() => {
    try {
      const { data: { code, body } } = await SavorApi.prizeConfig();
      if (code === 0) {
        guestName.value = body.guest_name;
      }
    } finally { 
      console.log()
    }
  }

  const editName = async() => {
    if (isEdit.value) {
      toggleEditLoading();
      try {
        const { data: { code } } = await SavorApi.renameMenu({ guest_name: guestName.value });
        if (code === 0) {
          isEdit.value = false;
        }
      } finally {
        toggleEditLoading();
      }
    } else {
      isEdit.value = true;
    }
  }

  onMounted(() => {
    fetchData();
    fetchSwitch();
  });
</script>

<style lang="less" scoped>
  .container {
    height: 100%;
    width: 100%;
    display: flex;
    overflow: hidden;
    color: #fff;

    .content {
      flex: 1;
      position: relative;
      display: flex;
      flex-direction: column;

      .header {
        border-bottom: 1px solid #4A5069;
        padding: 30px 38px 0px;
        font-size: 16px;
        font-weight: 400;
        display: flex;
        align-items: flex-start;
        justify-content: space-between;

        .menu {
          padding-bottom: 30px;
          cursor: pointer;
        
          &.selected {
            border-bottom: 2px solid #D2B276;
          }
        }

        :deep(.arco-btn) {
          border: 1px solid #C3C3C3;
          color: #fff;
          background-color: transparent;
        }
      }

      .btn-del {
        font-size: 20px;
        margin-right: 10px;
        cursor: pointer;
      }

      .handle {
        cursor: row-resize;
      }

      .table {
        margin-top: 32px;
        padding: 4px 48px;

        :deep(.arco-table) {
          .arco-table-container {
            border: unset;
          }

          .arco-table-th {
            background-color: #272A37;
            border: unset;
            font-size: 16px;
            font-weight: 600;
            color: #fff;
          }

          .arco-table-td {
            background-color: #272A37;
            border-color: #4A5069;
            color: #fff;
          }

          .arco-table-tr {
            &:hover {
              .arco-table-td {
                background-color: transparent !important;;
              }
            }
          }
        }
      }
      :deep(.arco-form-item-label), :deep(.arco-radio-label), :deep(.arco-picker-input) {
        color: #fff;
      }

      .nameEdit {
        display: flex;
        align-items: center;
        padding: 30px 48px 0;

        .input {
          width: 200px;
          margin: 0px 10px;
        }
      }
    }

    .form {
      padding: 30px 100px 30px 30px;
      overflow-y: auto;

      .title {
        width: 100%;
        display: flex;
        align-items: center;
        margin-bottom: 20px;

        .back {
          width: 12.5%;
          text-align: right;
          background-color: #272A37!important;
        }
        .subtitle {
          margin-left: 20px;
        }
        .primary {
          margin-left: auto;
        }
      }
    }

   

    :deep(.arco-modal-header) {
      padding: 40px 20px 20px;
      border-bottom: unset;
    }
  }

</style>
