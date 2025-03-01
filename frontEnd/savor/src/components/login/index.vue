<template>
  <a-modal
    v-model:visible="visible"
    :width="modalWidth"
    class="login-modal"
    :closable="false"
    :footer="false"
  >
    <div class="title">{{ $t('savor.login.title') }}</div>
    <a-form
      class="form"
      :model="form"
      :label-col-props="{ span: 0 }"
      :wrapper-col-props="{ span: 24 }"
    >
      <div class="type">{{ $t('savor.login.type') }}</div>
      <a-form-item field="phone">
        <a-input size="large" v-model="form.phone" :placeholder="$t('savor.login.form.phone.placeholder')">
          <template #prefix>
            +86
          </template>
        </a-input>
      </a-form-item>
      <a-form-item field="code" class="code">
        <a-input v-model="form.code" :placeholder="$t('savor.login.form.code.placeholder')">
          <template #append>
            <span v-if="coutdown">{{ coutdown }}{{  $t('savor.login.code') }}</span>
            <a-spin v-else :loading="loading">
              <span @click="getCode">{{  $t('savor.login.form.code.get') }}</span>
            </a-spin>
          </template>
        </a-input>
      </a-form-item>
      <a-button class="login-btn" :loading="loginLoading" @click="handleLogin">{{  $t('savor.login.btn') }}</a-button>
    </a-form>
  </a-modal>
</template>

<script lang="ts" setup>
  import { ref, reactive, computed, defineExpose } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import { useI18n } from 'vue-i18n';
  import { useUserStore } from '@/store';
  import useLoading from '@/hooks/loading';
  import SavorApi from '@/api/savor';

  const { t } = useI18n();
  const visible = ref(false);
  const form = reactive({
    phone: '',
    code: ''
  });
  const coutdown = ref<any>(null);
  const interval = ref<any>(null);
  const userStore = useUserStore();
  const { loading, toggle: toggleLoading } = useLoading(false);
  const { loading: loginLoading, toggle: toggleLoginLoading } = useLoading(false);
  const emit = defineEmits(['success']);

  const props = defineProps({
    mobile: {
      type: Boolean,
      default() {
        return false;
      },
    }
  });

  const modalWidth = computed(() => {
    return props.mobile ? '90%' : 431;
  });

  const showModal = () => {
    visible.value = true;
  }

  const getCode = async() => {
    const reg = new RegExp(/^1[3-9]\d{9}$/);
    if (!form.phone) {
      Message.error(t('savor.login.form.phone.placeholder'));
    } else if (!reg.test(form.phone)) {
      Message.error(t('settled.form.phone.validator1'));
    } else {
      toggleLoading();
      try {
        const { data: { code } } = await SavorApi.getCode({ phone: form.phone });
        if (code === 0) {
          coutdown.value = 60;
          interval.value = setInterval(() => {
            if (coutdown.value === 0) {
              clearInterval(interval.value);
              coutdown.value = null;
              interval.value = null;
            } else {
              coutdown.value -= 1;
            }
          }, 1000);
        } else {
          Message.error(t('savor.login.form.code.error'));
        }
      } finally {
        toggleLoading();
      }
    }
  }

  const handleLogin = async () => {
    const regPhone = new RegExp(/^1[3-9]\d{9}$/);
    const regCode = new RegExp(/^\d{6}$/);
    if (!form.phone) {
      Message.error(t('savor.login.form.phone.placeholder'));
    } else if (!regPhone.test(form.phone)) {
      Message.error(t('settled.form.phone.validator1'));
    } else if (!form.code) {
      Message.error(t('savor.login.form.code.placeholder'));
    } else if(!regCode.test(form.code)) {
      Message.error(t('savor.login.form.code.validator'));
    } else {
      toggleLoginLoading();
      try {
        await userStore.login({ phone: form.phone, code: form.code });
        Message.success(t('savor.login.success'));
        visible.value = false;
        emit('success');
        await userStore.info();
      } catch (err) {
        Message.error(t('savor.login.fail'));
      } finally {
        toggleLoginLoading();
      }
    }
  };

  defineExpose({
    showModal
  })
</script>

<style scoped lang="less">
  .login-modal {
    font-family: PingFang SC;
    .title {
      font-size: 20px;
      font-weight: 600;
      line-height: 28px;
      text-align: center;
      margin-bottom: 20px;
    }

    .form {
      .type {
        font-size: 14px;
        font-weight: 400;
        line-height: 21px;
        text-align: left;
        margin-bottom: 10px;
      }


      :deep(.arco-form-item-label-required-symbol) {
        display: none;
      }

      :deep(.arco-form-item-label), :deep(.arco-radio-label), :deep(.arco-picker-input) {
        color: #fff;
      }

      :deep(.arco-input) {
        height: 48px;
        color: #000;
      }

      :deep(.arco-input-wrapper), :deep(.arco-textarea-wrapper), :deep(.arco-picker) {
        color: #979797!important;
        background: #fff!important;
        border: 1px solid #979797!important;
        border-radius: 5px;
      }

      .code {
        :deep(.arco-input-wrapper) {
          border-top-right-radius: 0px;
          border-bottom-right-radius: 0px;
        }
      }

      :deep(.arco-input-append) {
        text-align: center;
        border: 1px solid #979797;
        border-left: unset;
        color: #D1B276;
        background-color: #fff;
        border-radius: 5px;
        border-top-left-radius: 0px;
        border-bottom-left-radius: 0px;
        cursor: pointer;
        padding: 0px;

        :deep(.arco-spin) {
          height: 100%;
        }

        span {
          align-items: center;
          display: flex;
          justify-content: center;
          width: 100px;
          height: 100%;
          text-align: center;
          color: #D1B276;
          font-weight: 500;
        }
      }

      .login-btn {
        width: 100%;
        height: 56px;
        border-radius: 5px;
        background: #D1B276;
        color: #fff;
        font-size: 16px;
        font-weight: 600;
      }
    }
  }
 
</style>
