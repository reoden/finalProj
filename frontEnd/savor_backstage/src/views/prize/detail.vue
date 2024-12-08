<template>
  <div class="container">
    <MenuBar />
    <a-spin v-if="loading" :size="32" />
    <div v-else class="content">
      <a-form
        ref="formRef" :rules="rules" :model="form"
        class="form"
        :label-col-props="{ span: 3 }"
        :wrapper-col-props="{ span: 21 }"
        :style="{ width: '100%' }"
        @submit="handleSubmit"
      >
        <div class="title">
          <div class="back">
            <a-button class="form-btn" @click="handleCancel">
              <icon-left style="margin-right: 5px;" />
              返回
            </a-button>
          </div>
          <div class="subtitle">大奖详情</div>
          <a-button class="form-btn primary" html-type="submit" :loading="btnLoading" >保存</a-button>
        </div>
        <a-form-item field="pictures" :label="$t('settled.form.pic')" validate-trigger="blur">
          <a-upload
            list-type="picture-card"
            action="/"
            :file-list="fileList"
            :custom-request="customRequest"
            @before-remove="beforeRemove"
          />
        </a-form-item>
        <a-form-item field="name" :label="$t('news.colunm.name')" validate-trigger="blur">
          <a-input v-model="form.name" :placeholder="$t('news.form.name.placeholder')" />
        </a-form-item>
        <a-form-item field="describe" :label="$t('guest.colunm.describe')" validate-trigger="blur">
          <a-textarea v-model="form.describe" :placeholder="$t('news.form.describe.placeholder')" :auto-size="false" />
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import { useI18n } from 'vue-i18n';
  import { useRouter } from 'vue-router';
  import useLoading from '@/hooks/loading';
  import SavorApi from '@/api/savor';
  import MenuBar from '@/components/menu-bar/index.vue';
  import { Message, Modal } from '@arco-design/web-vue';

  const { t } = useI18n();
  const { loading, toggle: toggleLoading } = useLoading(false);
  const router = useRouter();
  const formRef = ref();
  const initForm = {
    pictures: [],
    name: '',
    describe: '',
  }
  const form = ref<any>(initForm);
  const fileList = ref<any>([]);

  // 获取列表数据
  const fetchData = async () => {
    toggleLoading();
    try {
      const { id } = router.currentRoute.value.query;
      if (!id) {
        return;
      }
      const { data: { code, body } } = await SavorApi.prizeDetail(id);
      if (code === 0) {
        form.value = body;
        fileList.value = body.pic_url && body.pic_url.map((item: any) => { return { url: item } })
      }
    } finally {
      toggleLoading();
    }
  };

  const rules: any = {
    name: [
      {
        required: true,
        message: t('chef.form.name.validator'),
      },
    ],
    describe: [
      {
        required: true,
        message: t('chef.form.describe.validator'),
      },
    ],
    pictures: [
      {
        required: true,
        message: t('chef.form.picture.validator'),
      },
    ]
  };
  const { loading: btnLoading, toggle: toggleBtnLoading } = useLoading(false);

  const handleCancel = () => {
    router.push({ name: 'prize' });
  }

  const handleSubmit = async() => {
    await formRef.value.validate(async (valid: any) => {
		  if (!valid) {
        toggleBtnLoading();
        try {
          const { data: { code } } = await SavorApi.addPrize(form.value);
          if (code === 0) {
            handleCancel();
            Message.success('保存成功');
            fetchData();
          }
        } finally{
          toggleBtnLoading();
        }
      }
    })
  };

  const customRequest = (option: any) => {
    const { onProgress, onError, onSuccess, fileItem, name } = option
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
        url: fielUrl
      }
      fileList.value = fileList.value.length ? fileList.value.concat(curFile): [curFile];
      form.value.pictures = form.value.pictures.length ? form.value.pictures.concat(fielName) : [fielName];
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
        xhr.abort()
      }
    }
  };

  const beforeRemove = (file: any) => {
    return new Promise((resolve, reject) => {
      Modal.confirm({
        content: '确认删除该图片吗？',
        onOk: () => {
          const newFileList = fileList.value.filter((item: any) => item.url !== file.url);
          let pictures = form.value.pictures.filter((item: any) => item !== file.url);
          pictures = form.value.pictures.filter((item: any) =>  file.url.indexOf(item) < 0); 
          fileList.value = newFileList;
          form.value.pictures = pictures;
          resolve(true);
        },
        onCancel: () => {},
      });
    });
  };

  onMounted(() => {
    fetchData();
    console.log('test', import.meta.env.VITE_API_BASE_URL);
  });
</script>

<style lang="less" scoped>
  .container {
    height: 100%;
    width: 100%;
    display: flex;
    overflow: hidden;
    color: #fff;

    :deep(.arco-spin) {
      margin-top: 60px;
    }
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
