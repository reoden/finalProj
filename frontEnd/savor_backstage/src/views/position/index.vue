<template>
  <div class="container">
    <MenuBar />
    <div class="content">
      <div class="header">
        <div class="menus">
          <div
            class="menu"
            :class="{ selected: menu === 'shop' }"
            @click="changeMenu('shop')"
          >
            {{ $t('menu.shop.subtitle') }}
          </div>
          <div
            class="menu"
            :class="{ selected: menu === 'banner' }"
            @click="changeMenu('banner')"
          >
            {{ $t('menu.position.subtitle') }}
          </div>
          <div
            class="menu"
            :class="{ selected: menu === 'logo' }"
            @click="changeMenu('logo')"
          >
            {{ $t('menu.logo.subtitle') }}
          </div>
        </div>
        <a-button v-if="menu === 'shop'" @click="addPosition"
          >增加展示位</a-button
        >
      </div>
      <div v-if="menu === 'shop'" class="nameEdit">
        <div>景点名修改</div>
        <a-input v-model="appName" class="input" :disabled="!isEdit"></a-input>
        <a-button :loading="editLoading" @click="editName">{{
          isEdit ? '保存' : '编辑'
        }}</a-button>
      </div>
      <div class="table">
        <a-table
          v-if="menu === 'banner'"
          :loading="loading"
          :columns="bannerColumns"
          :data="bannerList"
          row-key="id"
          :border="false"
          :pagination="false"
        >
          <template #index="{ rowIndex }">
            {{ rowIndex + 1 }}
          </template>
          <template #picture="{ record }">
            <img
              :src="record.pic_url"
              style="width: 150px; height: 100px; object-fit: contain"
            />
          </template>
          <template #operate="{ record }">
            <div style="cursor: pointer" @click="showPic(record, 'banner')">
              <IconEdit style="margin-right: 3px" />编辑
            </div>
          </template>
        </a-table>
        <a-table
          v-else-if="menu === 'logo'"
          :loading="loading"
          :columns="logoColumns"
          :data="logoList"
          row-key="id"
          :border="false"
          :pagination="false"
        >
          <template #id="{ record }">
            {{
              Number(record.id) === 1
                ? '左上角logo'
                : Number(record.id) === 2
                ? '首页主logo'
                : record.id
            }}
          </template>
          <template #pic="{ record }">
            <img
              :src="record.pic_url"
              style="width: 150px; height: 100px; object-fit: contain"
            />
          </template>
          <template #operate="{ record }">
            <div style="cursor: pointer" @click="showPic(record, 'logo')">
              <IconEdit style="margin-right: 3px" />更换图片
            </div>
          </template>
        </a-table>
        <a-table
          v-else
          :loading="loading"
          :columns="columns"
          :data="positionList"
          row-key="id"
          :border="false"
          :pagination="false"
        >
          <template #index="{ rowIndex }">
            {{ rowIndex + 1 }}
          </template>
          <template #operate="{ record }">
            <icon-delete class="btn-del" title="移除" @click="del(record)" />
            <img
              class="handle"
              :src="iconDropImage"
              style="width: 15px; height: 15px; object-fit: contain"
            />
          </template>
        </a-table>
      </div>
    </div>
    <a-modal
      v-model:visible="visible"
      title="选择景点"
      :closable="false"
      :width="431"
      :footer="false"
      modal-class="custom-modal"
      :render-to-body="false"
      @cancel="handleCancel"
    >
      <a-form
        :model="form"
        :label-col-props="{ span: 0 }"
        :wrapper-col-props="{ span: 24 }"
      >
        <div style="margin-bottom: 16px"
          >展示位号：{{ positionList.length + 1 }}</div
        >
        <div style="margin-bottom: 16px">请输入景点ID或景点名</div>
        <a-form-item field="app_id" label="">
          <a-select
            v-model="form.id"
            placeholder="ID是8位数字，或者景点名称"
            allow-search
            :filter-option="filterOption"
            @search="queryList"
          >
            <a-option
              v-for="shop in shopList"
              :key="shop.id"
              :value="shop.id"
              :label="shop.name"
              >{{ shop.name }}</a-option
            >
          </a-select>
        </a-form-item>
        <a-button
          class="primary"
          :disabled="!form.id"
          :loading="btnLoading"
          @click="handleOk"
          >确认</a-button
        >
      </a-form>
    </a-modal>

    <a-modal
      v-model:visible="visiblePic"
      title="编辑"
      :closable="false"
      :width="431"
      :footer="false"
      modal-class="custom-modal"
      :render-to-body="false"
      @cancel="handleCancelPic"
    >
      <a-form
        :model="form"
        :label-col-props="{ span: 4 }"
        :wrapper-col-props="{ span: 20 }"
      >
        <a-form-item label="图片">
          <a-upload
            action="/"
            :file-list="file"
            :show-file-list="false"
            style="margin: 0 auto"
            :custom-request="customRequest"
          >
            <template #upload-button>
              <div style="position: relative">
                <img
                  :src="file?.[0]?.url"
                  style="width: 150px; height: 150px; object-fit: contain"
                />
                <div
                  class="arco-upload-list-picture-mask"
                  style="line-height: 150px; margin-bottom: 15px"
                >
                  <IconEdit />
                </div>
              </div>
            </template>
          </a-upload>
        </a-form-item>
        <a-form-item field="redirect_url" label="链接" validate-trigger="blur">
          <a-input v-model="form.redirect_url" placeholder="请输入链接" />
        </a-form-item>
        <a-button
          class="primary"
          :disabled="!curId"
          :loading="btnLoading"
          @click="handleOkPic"
          >确认</a-button
        >
      </a-form>
    </a-modal>

    <a-modal
      v-model:visible="visibleLogo"
      title="更换图片"
      :closable="false"
      :width="431"
      :footer="false"
      modal-class="custom-modal"
      :render-to-body="false"
      @cancel="handleCancelPic"
    >
      <a-form
        :model="form"
        :label-col-props="{ span: 0 }"
        :wrapper-col-props="{ span: 24 }"
      >
        <a-upload
          action="/"
          :file-list="file"
          :show-file-list="false"
          style="margin: 0 auto"
          :custom-request="customRequest"
        >
          <template #upload-button>
            <div style="position: relative">
              <img
                :src="file?.[0]?.url"
                style="width: 150px; height: 150px; object-fit: contain"
              />
              <div
                class="arco-upload-list-picture-mask"
                style="line-height: 150px; margin-bottom: 15px"
              >
                <IconEdit />
              </div>
            </div>
          </template>
        </a-upload>
        <a-button
          class="primary"
          :disabled="!curId"
          :loading="btnLoading"
          @click="handleOkLogo"
          >确认</a-button
        >
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
  import SavorApi from '@/api/savor';
import iconDropImage from '@/assets/images/icon-drop.png';
import MenuBar from '@/components/menu-bar/index.vue';
import useLoading from '@/hooks/loading';
import { Message, Modal } from '@arco-design/web-vue';
import Sortable from 'sortablejs';
import { computed, onMounted, reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';

  const { t } = useI18n();
  const positionList = ref<any>([]);
  const { loading, toggle: toggleLoading } = useLoading(false);
  const menu = ref<any>('shop');
  const bannerList = ref<any>([]);
  const logoList = ref<any>([]);
  const { loading: editLoading, toggle: toggleEditLoading } = useLoading(false);
  const appName = ref('景点收集');
  const isEdit = ref(false);

  const columns = computed(() => {
    return [
      {
        title: t('position.colunm.id'),
        slotName: 'index',
        dataIndex: 'index',
      },
      {
        title: t('examine.colunm.id'),
        dataIndex: 'id',
      },
      {
        title: t('examine.colunm.name'),
        dataIndex: 'name',
      },
      {
        dataIndex: 'operate',
        slotName: 'operate',
        title: t('position.colunm.operate'),
      },
    ];
  });

  const bannerColumns = computed(() => {
    return [
      {
        title: t('position.colunm.id'),
        slotName: 'index',
        dataIndex: 'index',
      },
      {
        dataIndex: 'pic_url',
        slotName: 'picture',
        title: t('settled.detail.pic'),
      },
      {
        dataIndex: 'redirect_url',
        slotName: 'redirect_url',
        title: '链接',
      },
      {
        dataIndex: 'operate',
        slotName: 'operate',
        title: t('settled.detail.operate'),
      },
    ];
  });

  const logoColumns = computed(() => {
    return [
      {
        title: '展示位',
        slotName: 'id',
        dataIndex: 'id',
      },
      {
        dataIndex: 'pic_url',
        slotName: 'pic',
        title: t('settled.detail.pic'),
      },
      {
        dataIndex: 'operate',
        slotName: 'operate',
        title: t('settled.detail.operate'),
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
      const {
        data: { code, body },
      } = await SavorApi.positionList();
      if (code === 0) {
        positionList.value = body?.apps || [];
        const tbody: any = document.querySelector('.table .arco-table tbody');
        const drag = new Sortable(tbody, {
          handle: '.handle',
          animation: 150,
          // 需要在odEnd方法中处理原始eltable数据, 使原始数据与显示数据保持顺序一致
          onEnd: async (evt) => {
            const { newIndex, oldIndex }: any = evt;
            // 拖拽后接口保存排序后的位置
            const tableData = JSON.parse(JSON.stringify(positionList.value));
            const targetRow = tableData[oldIndex];
            tableData.splice(oldIndex, 1);
            tableData.splice(newIndex, 0, targetRow);
            positionList.value = tableData;
            try {
              const appIds = tableData.map((item: any) => item.id);
              const {
                data: { code: codeSec },
              } = await SavorApi.changePosition({ app_ids: appIds });
              if (codeSec !== 0) {
                // 保存失败后，撤回本次拖拽
                returnDrop(evt);
              }
            } catch (error) {
              returnDrop(evt);
            }
          },
        });
      }
    } finally {
      toggleLoading();
    }
  };

  const visible = ref(false);
  const form = reactive({ id: undefined, redirect_url: '' });
  const { loading: btnLoading, toggle: toggleBtnLoading } = useLoading(false);

  const addPosition = () => {
    visible.value = true;
  };

  const handleOk = async () => {
    toggleBtnLoading();
    try {
      const {
        data: { code },
      } = await SavorApi.addPosition({ app_id: Number(form.id) });
      if (code === 0) {
        visible.value = false;
        Message.success('添加成功');
        fetchData();
      }
    } finally {
      toggleBtnLoading();
    }
  };

  const handleCancel = () => {
    visible.value = false;
    form.id = undefined;
    form.redirect_url = '';
  };

  const shopList = ref<any>([]);
  const queryList = async (val?: any) => {
    try {
      if (!val) {
        return;
      }
      const {
        data: { code, body },
      } = await SavorApi.searchShop({ query: val, page: 1, page_size: 20 });
      if (code === 0) {
        shopList.value = body?.apps;
      }
    } catch {
      shopList.value = [];
    }
  };

  const filterOption = (value: any, option: any) => {
    return (
      option.value.toString().indexOf(value) >= 0 ||
      option.label.indexOf(value) >= 0
    );
  };

  const del = (row?: any) => {
    Modal.confirm({
      content: `确认将商家 ${row.name} 移除广告位吗?`,
      okText: '确定',
      async onOk() {
        const {
          data: { code },
        } = await SavorApi.deletePosition(row.id);
        if (code === 0) {
          Message.success('移除成功');
          fetchData();
        }
      },
    });
  };

  const fetchBannerData = async () => {
    toggleLoading();
    try {
      const {
        data: { code, body },
      } = await SavorApi.staticBanner();
      if (code === 0) {
        bannerList.value = body?.apps;
      }
    } finally {
      toggleLoading();
    }
  };

  const fetchLogoData = async () => {
    toggleLoading();
    try {
      const {
        data: { code, body },
      } = await SavorApi.staticLogo();
      if (code === 0) {
        logoList.value = body?.apps;
      }
    } finally {
      toggleLoading();
    }
  };

  const changeMenu = (val: string) => {
    menu.value = val;
    if (val === 'banner') {
      fetchBannerData();
    } else if (val === 'logo') {
      fetchLogoData();
    } else {
      fetchData();
    }
  };

  const visiblePic = ref(false);
  const visibleLogo = ref(false);
  const file = ref<any>([]);
  const curId = ref('');
  const modalType = ref('');

  const showPic = (record: any, type: any) => {
    curId.value = record.id;
    modalType.value = type;
    if (type === 'logo') {
      visibleLogo.value = true;
      file.value = [{ url: record.pic_url, name: record.pic_file }];
    } else {
      form.redirect_url = record.redirect_url;
      visiblePic.value = true;
      file.value = [{ url: record.pic_url, name: record.pic_file }];
    }
  };

  const customRequest = (option: any) => {
    const { onProgress, onError, onSuccess, fileItem, name } = option;
    const xhr = new XMLHttpRequest();
    if (xhr.upload) {
      xhr.upload.onprogress = function (event) {
        let percent;
        if (event.total > 0) {
          // 0 ~ 1
          percent = event.loaded / event.total;
        }
        onProgress(percent, event);
      };
    }
    xhr.onerror = function error(e) {
      onError(e);
    };
    xhr.onload = function onload() {
      if (xhr.status < 200 || xhr.status >= 300) {
        return onError(xhr.responseText);
      }
      onSuccess(xhr.response);
      const res = JSON.parse(xhr.response);
      const fielUrl = res?.body?.file_url;
      const fielName = res?.body?.file_name;
      const curFile = {
        url: fielUrl,
        name: fielName,
      };
      file.value = [curFile];
      return null;
    };

    const formData = new FormData();
    formData.append(name || 'file', fileItem.file);
    const url = `${import.meta.env.VITE_API_BASE_URL}restaurant/v1/apps/pic`;
    xhr.open('post', url, true);
    // xhr.open('post', 'restaurant/v1/apps/pic', true);
    xhr.send(formData);

    return {
      abort() {
        xhr.abort();
      },
    };
  };

  const handleCancelPic = () => {
    visiblePic.value = false;
    visibleLogo.value = false;
    file.value = [];
    curId.value = '';
  };

  const handleOkPic = async () => {
    toggleBtnLoading();
    try {
      const {
        data: { code },
      } = await SavorApi.changeStaticBanner(curId.value, {
        picture: file.value?.[0]?.name,
        redirect_url: form.redirect_url,
      });
      if (code === 0) {
        handleCancelPic();
        Message.success('修改成功');
        fetchBannerData();
      }
    } finally {
      toggleBtnLoading();
    }
  };

  const handleOkLogo = async () => {
    toggleBtnLoading();
    try {
      const {
        data: { code },
      } = await SavorApi.changeStaticLogo({
        id: curId.value,
        pic: file.value?.[0]?.name,
      });
      if (code === 0) {
        handleCancelPic();
        Message.success('修改成功');
        fetchLogoData();
      }
    } finally {
      toggleBtnLoading();
    }
  };

  const fetchSwitch = async () => {
    try {
      const {
        data: { code, body },
      } = await SavorApi.prizeConfig();
      if (code === 0) {
        appName.value = body.app_name;
      }
    } finally {
      console.log();
    }
  };

  const editName = async () => {
    if (isEdit.value) {
      toggleEditLoading();
      try {
        const {
          data: { code },
        } = await SavorApi.renameMenu({ app_name: appName.value });
        if (code === 0) {
          isEdit.value = false;
        }
      } finally {
        toggleEditLoading();
      }
    } else {
      isEdit.value = true;
    }
  };

  onMounted(() => {
    fetchData();
    queryList();
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
        border-bottom: 1px solid #4a5069;
        padding: 30px 38px 0px;
        font-size: 16px;
        font-weight: 400;
        display: flex;
        justify-content: space-between;

        .menus {
          display: flex;
        }

        .menu {
          padding-bottom: 30px;
          margin-right: 40px;
          cursor: pointer;

          &.selected {
            border-bottom: 2px solid #d2b276;
          }
        }

        :deep(.arco-btn) {
          border: 1px solid #c3c3c3;
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
        padding: 4px 48px;

        :deep(.arco-table) {
          .arco-table-container {
            border: unset;
          }

          .arco-table-th {
            background-color: #272a37;
            border: unset;
            font-size: 16px;
            font-weight: 600;
            color: #fff;
          }

          .arco-table-td {
            background-color: #272a37;
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
        }
      }
      .nameEdit {
        display: flex;
        align-items: center;
        padding: 30px 48px;

        .input {
          width: 200px;
          margin: 0px 10px;
        }
      }
    }

    :deep(.arco-modal-header) {
      padding: 40px 20px 20px;
      border-bottom: unset;
    }
    :deep(.arco-btn.primary) {
      width: 183px;
      height: 40px;
      margin: 0px auto;
      border: unset;
      color: #fff;
      background: #d1b276;
    }
  }
</style>

<style lang="less">
  .custom-modal {
    color: #fff;
    .arco-form-item-label {
      color: #4e5969 !important;
    }
    .arco-input-wrapper {
      color: #1d2129 !important;
      background: #f2f3f5 !important;
      border: 1px solid #f2f3f5 !important;
    }
  }
</style>
