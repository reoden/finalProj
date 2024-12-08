<template>
  <a-form
    class="form"
    ref="formRef" :model="data"
    :label-col-props="{ span: 5 }"
    :wrapper-col-props="{ span: 19 }"
    :style="{ width: '100%' }"
  >
  <a-form-item content-class="picCon" field="pictures" :label="$t('settled.detail.pic')" validate-trigger="blur">
      <img class="pic" v-for="pic in data.pics_url" :key="pic" :src="pic" alt="">
    </a-form-item>
    <a-form-item field="name" :label="$t('settled.form.name')" validate-trigger="blur">
      {{ data.name }}
    </a-form-item>
    <a-form-item field="introduction" :label="$t('settled.form.introduction')" validate-trigger="blur">
      {{ data.introduction }}
    </a-form-item>
    <a-form-item field="leixing" :label="$t('settled.form.leixing')" validate-trigger="blur">
      {{ data.leixing }}
    </a-form-item>
    <a-form-item field="time" :label="$t('settled.form.time')">
      {{ data.work_date_label.join('、') }} {{ data.work_begin_at }} ~ {{ data.work_end_at }}
    </a-form-item>
    <a-divider />
    <a-form-item field="describe" :label="$t('settled.form.describe')" validate-trigger="blur">
      {{ data.describe }}
    </a-form-item>
    <a-form-item field="hostman_name" :label="$t('settled.form.hostman_name')" validate-trigger="blur">
      {{ data.hostman_name }}
    </a-form-item>
    <a-form-item content-class="picCon" field="hostman_urls" :label="$t('settled.form.hostman_pics')" validate-trigger="blur">
      <img class="pic" v-for="pic in data.hostman_urls" :key="pic" :src="pic" alt="">
    </a-form-item>
    <a-form-item field="hostman_think" :label="$t('settled.form.hostman_think')" validate-trigger="blur">
      {{ data.hostman_think }}
    </a-form-item>
    <div v-for="(hot, index) of data.hot" :key="index">
      <a-form-item :field="`hot[${index}].hot_name`" :label="$t('settled.form.hot_name')">
        {{ hot.hot_name }}
      </a-form-item>
      <a-form-item content-class="picCon" field="`hot[${index}].hot_pics`" :label="$t('settled.form.hot_pics')" validate-trigger="blur">
        <img class="pic" v-for="pic in hot.hot_pic_urls" :key="pic" :src="pic" alt="">
      </a-form-item>
      <a-form-item :field="`hot[${index}].hot_desc`" :label="$t('settled.form.hot_desc')">
        {{ hot.hot_desc }}
      </a-form-item>
    </div>
    <a-divider />
    <a-form-item field="phone" :label="$t('settled.form.phone')" validate-trigger="blur">
      {{ data.phone }}
    </a-form-item>
    <a-form-item field="yuyue" :label="$t('settled.form.yuyue')">
      {{  Number(data.yuyue) === 1 ? $t('settled.form.have_vege.option1') : $t('settled.form.have_vege.option2') }}
    </a-form-item>
    <a-form-item field="address" :label="$t('settled.form.address')" validate-trigger="blur">
     {{ data.post_name }} {{ data.address }}
    </a-form-item>
    <a-form-item field="have_vege" :label="$t('settled.form.have_vege')">
      {{  Number(data.have_vege) === 1 ? $t('settled.form.have_vege.option1') : $t('settled.form.have_vege.option2') }}
    </a-form-item>
    <a-form-item field="per" :label="$t('settled.form.per')" validate-trigger="blur">
      {{ data.per }} {{ $t('settled.form.per.desc') }}
    </a-form-item>
  </a-form>
</template>

<script lang="ts" setup>

  const props = defineProps({
    data: {
      type: Object as any,
      default() {
        return [];
      },
    },
  });
</script>

<style lang="less" scoped>
  .form {
    color: #fff;

    :deep(.arco-form-item-label-col) {
      padding-right: 60px;
    }

    :deep(.picCon) {
      flex-wrap: wrap;
    }

    .pic {
      width: 135px;
      height: 135px;
      margin-right: 24px;
      margin-bottom: 10px;
    }
  }
</style>
