<template>
  <div class="container">
    <Header :back="true" />
    <div class="content" :class="{ mobile: mobile }">
      <a-spin v-if="loading" :size="32" />
      <div v-else class="form" :class="{ view: isView }">
        <a-form
          v-if="!isView"
          ref="formRef"
          :rules="rules"
          :model="form"
          :label-col-props="{ span: 6 }"
          :wrapper-col-props="{ span: 18 }"
          :style="{ width: '100%' }"
          @submit="handleSubmit"
        >
          <div class="title">
            <span>{{ $t('settled.con.title') }}</span>
            <a-space v-if="!mobile">
              <a-button
                class="form-btn"
                :loading="saveLoading"
                @click="handleSubmit('save')"
                >{{ $t('settled.form.btn1') }}</a-button
              >
              <a-button
                class="form-btn primary"
                html-type="submit"
                :loading="subLoading"
                >{{ $t('settled.form.btn2') }}</a-button
              >
            </a-space>
          </div>
          <a-form-item
            field="pictures"
            :label="$t('settled.form.pic')"
            validate-trigger="blur"
          >
            <a-upload
              list-type="picture-card"
              action="/"
              :limit="10"
              :file-list="fileList"
              :custom-request="customRequest"
              @before-remove="beforeRemove"
            />
          </a-form-item>
          <a-form-item
            field="name"
            :label="$t('settled.form.name')"
            validate-trigger="blur"
          >
            <a-input
              v-model="form.name"
              :placeholder="$t('settled.form.name.placeholder')"
            />
          </a-form-item>
          <a-form-item
            field="introduction"
            :label="$t('settled.form.introduction')"
            validate-trigger="blur"
          >
            <a-textarea
              v-model="form.introduction"
              :placeholder="$t('settled.form.introduction.placeholder')"
              :auto-size="false"
              style="height: 200px"
            />
          </a-form-item>
          <!-- <a-form-item
            field="leixing"
            :label="$t('settled.form.leixing')"
            validate-trigger="blur"
          >
            <a-select
              v-model="form.leixing"
              :placeholder="$t('settled.form.leixing.placeholder')"
              :style="{ width: '300px' }"
            >
              <a-option v-for="item in leixingArr" :key="item">{{
                item
              }}</a-option>
            </a-select>
          </a-form-item>
          <a-divider /> -->
          <!-- <a-form-item
            field="describe"
            :label="$t('settled.form.describe')"
            validate-trigger="blur"
          >
            <a-textarea
              v-model="form.describe"
              :placeholder="$t('settled.form.describe.placeholder')"
              :auto-size="false"
              style="height: 200px"
            />
          </a-form-item> -->
          <!-- <a-form-item
            field="hostman_name"
            :label="$t('settled.form.hostman_name')"
            validate-trigger="blur"
          >
            <a-input
              v-model="form.hostman_name"
              :placeholder="$t('settled.form.hostman_name.placeholder')"
            />
          </a-form-item>
          <a-form-item
            field="hostman_pics"
            :label="$t('settled.form.hostman_pics')"
            validate-trigger="blur"
          >
            <a-upload
              list-type="picture-card"
              action="/"
              :limit="10"
              :file-list="hostmanList"
              :custom-request="uploadFileHost"
              @before-remove="beforeRemoveHost"
            />
          </a-form-item>
          <a-form-item
            field="hostman_think"
            :label="$t('settled.form.hostman_think')"
            validate-trigger="blur"
          >
            <a-textarea
              v-model="form.hostman_think"
              :placeholder="$t('settled.form.hostman_think.placeholder')"
              :auto-size="false"
              style="height: 100px"
            />
          </a-form-item>
          <div v-for="(hot, index) of form.hot" :key="index">
            <a-form-item
              :field="`hot[${index}].hot_name`"
              :label="$t('settled.form.hot_name')"
              :rules="rules.hot_name"
            >
              <a-input
                v-model="hot.hot_name"
                :placeholder="$t('settled.form.hot_name.placeholder')"
              />
              <a-button
                :style="{ marginLeft: '10px' }"
                @click="handleDelete(index)"
                >{{ $t('settled.form.del_hot') }}</a-button
              >
            </a-form-item>
            <a-form-item
              :field="`hot[${index}].hot_pics`"
              :label="$t('settled.form.hot_pics')"
              :rules="rules.hot_pics"
              validate-trigger="blur"
            >
              <a-upload
                list-type="picture-card"
                action="/"
                :limit="10"
                :file-list="hot.hot_pic_urls"
                :custom-request="(file: any) => uploadFileList(file, index)"
                @before-remove="(file: any) => beforeRemoveList(file, index)"
              />
            </a-form-item>
            <a-form-item
              :field="`hot[${index}].hot_desc`"
              :label="$t('settled.form.hot_desc')"
              :rules="rules.hot_desc"
            >
              <a-textarea
                v-model="hot.hot_desc"
                :placeholder="$t('settled.form.hot_desc.placeholder')"
                :auto-size="false"
                style="height: 100px"
              />
            </a-form-item>
          </div> -->
          <!-- <a-form-item>
            <a-button @click="handleAdd"
              ><icon-plus /> {{ $t('settled.form.plus_hot') }}</a-button
            >
          </a-form-item>
          <a-divider /> -->
          <div class="dateWrap">
            <a-form-item
              field="work_date"
              :label="$t('settled.form.time')"
              :label-col-props="{ span: mobile ? 6 : 12 }"
              :wrapper-col-props="{ span: mobile ? 18 : 12 }"
            >
              <a-select
                v-model="form.work_date"
                :placeholder="$t('settled.form.work_date.placeholder')"
                :style="{ width: '300px' }"
                multiple
                allow-clear
              >
                <a-option
                  v-for="item in dateArr"
                  :key="item.value"
                  :value="item.value"
                  >{{ item.label }}</a-option
                >
              </a-select>
            </a-form-item>
            <a-form-item
              field="time"
              label=""
              :label-col-props="{ span: mobile ? 6 : 1 }"
              :wrapper-col-props="{ span: mobile ? 18 : 23 }"
            >
              <a-time-picker
                v-model="form.time"
                type="time-range"
                style="width: 300px"
                format="HH:mm"
                class="test"
                :placeholder="[
                  $t('settled.form.time.placeholder1'),
                  $t('settled.form.time.placeholder2'),
                ]"
              />
            </a-form-item>
          </div>
          <a-form-item
            field="address"
            content-class="address"
            :label="$t('settled.form.address')"
            validate-trigger="blur"
          >
            <div>
              <a-cascader
                v-model="form.post_code"
                :options="addressOption"
                path-mode
                expand-trigger="hover"
                :style="{ width: mobile ? '275px' : '300px' }"
                :placeholder="$t('settled.form.province.placeholder')"
                @change="handleChange"
              />
            </div>
            <div>
              <a-input
                v-model="form.address"
                :placeholder="$t('settled.form.address.placeholder')"
              />
            </div>
          </a-form-item>
          <a-form-item
            field="phone"
            :label="$t('settled.form.phone')"
            validate-trigger="blur"
          >
            <a-input
              v-model="form.phone"
              :placeholder="$t('settled.form.phone.placeholder')"
            />
          </a-form-item>
          <!-- <a-form-item field="yuyue" :label="$t('settled.form.yuyue')">
            <a-radio-group v-model="form.yuyue">
              <a-radio :value="1">{{
                $t('settled.form.have_vege.option1')
              }}</a-radio>
              <a-radio :value="0">{{
                $t('settled.form.have_vege.option2')
              }}</a-radio>
            </a-radio-group>
          </a-form-item> -->
          <a-form-item
            field="per"
            :label="$t('settled.form.per')"
            validate-trigger="blur"
          >
            <a-input
              v-model="form.per"
              :placeholder="$t('settled.form.per.placeholder')"
              style="width: 240px"
            >
              <template #append>{{ $t('settled.form.per.desc') }}</template>
            </a-input>
          </a-form-item>
          <!-- <a-form-item field="have_vege" :label="$t('settled.form.have_vege')">
            <a-radio-group v-model="form.have_vege">
              <a-radio :value="1">{{
                $t('settled.form.have_vege.option1')
              }}</a-radio>
              <a-radio :value="0">{{
                $t('settled.form.have_vege.option2')
              }}</a-radio>
            </a-radio-group>
          </a-form-item> -->
          <a-form-item v-if="mobile" field="per" label="">
            <a-space>
              <a-button
                class="form-btn"
                :loading="saveLoading"
                @click="handleSubmit('save')"
                >{{ $t('settled.form.btn1') }}</a-button
              >
              <a-button
                class="form-btn primary"
                html-type="submit"
                :loading="subLoading"
                >{{ $t('settled.form.btn2') }}</a-button
              >
            </a-space>
          </a-form-item>
        </a-form>
        <Detail v-else :data="form" />
      </div>
      <a-divider
        v-if="!isView && !loading"
        direction="vertical"
        class="divider"
      />
      <div v-if="!isView && !loading" class="example">
        <div class="title">{{ $t('settled.con.example') }}</div>
        <ShopItem
          v-for="(item, index) in exampleList"
          :key="index"
          :data="item"
        />
      </div>
      <div v-if="isView" class="right">
        <div
          >{{ $t('settled.detail.status') }} :{{ STATUS_MAP[form.status] }}</div
        >
        <div class="operate"
          >{{ $t('settled.detail.operate') }}:<a-button
            class="agree"
            @click="edit"
            >{{ $t('settled.detail.operate.btn') }}</a-button
          ></div
        >
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import SavorApi from '@/api/savor';
import Header from '@/components/header/index.vue';
import { STATUS_MAP } from '@/config';
import useLoading from '@/hooks/loading';
import { isMoblie } from '@/utils';
import ShopItem from '@/views/home/components/shop-item.vue';
import { Message, Modal } from '@arco-design/web-vue';
import { pcaa } from 'area-data';
import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import Detail from './detail.vue';

  const { t } = useI18n();
  const { loading, toggle: toggleLoading } = useLoading(false);
  const { loading: saveLoading, toggle: toggleSaveLoading } = useLoading(false);
  const { loading: subLoading, toggle: toggleSubLoading } = useLoading(false);
  const mobile = computed(() => isMoblie());
  // const leixingArr = ['中餐', '西餐', '日料', '其他'];
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

  const validatePhone = (value: any, callback: (err: string) => void) => {
    const reg = new RegExp(/^(\+86)?\s*1[3-9]\d{9}$/);
    if (!reg.test(value)) {
      return callback(t('settled.form.phone.validator1'));
    }
    return true;
  };
  const formRef = ref();
  const initForm: any = {
    pictures: [],
    hostman_pics: [],
    post_code: '',
    name: '',
    leixing: '',
    work_date: [],
    introduction: '',
    address: '',
    hostman_think: '',
    hot: [],
    per: '',
    describe: '',
    hostman_name: '',
    phone: '',
    yuyue: 0,
    time: '',
    have_vege: 0,
  };
  const form = ref<any>(initForm);
  const isView = ref<any>(false);
  const fileList = ref<any>([]);
  const addressLabel = ref<any>('');
  const hostmanList = ref<any>([]);

  const rules: any = {
    name: [
      {
        required: true,
        message: t('settled.form.name.validator'),
      },
    ],
    address: [
      {
        required: true,
        message: t('settled.form.address.validator'),
      },
    ],
    describe: [
      {
        required: true,
        message: t('settled.form.describe.validator'),
      },
    ],
    phone: [
      {
        required: true,
        message: t('settled.form.phone.validator'),
      },
      {
        validator: validatePhone,
      },
    ],
    time: [
      {
        required: true,
        message: t('settled.form.time.validator'),
      },
    ],
    introduction: [
      {
        required: true,
        message: t('settled.form.introduction.validator'),
      },
    ],
    // leixing: [
    //   {
    //     required: true,
    //     message: t('settled.form.leixing.validator'),
    //   },
    // ],
    work_date: [
      {
        required: true,
        message: t('settled.form.work_date.validator'),
      },
    ],
    // hostman_name: [
    //   {
    //     required: true,
    //     message: t('settled.form.hostman_name.validator'),
    //   },
    // ],
    // hostman_pics: [
    //   {
    //     required: true,
    //     message: t('settled.form.pic.validator1'),
    //   },
    // ],
    // hostman_think: [
    //   {
    //     required: true,
    //     message: t('settled.form.hostman_think.validator'),
    //   },
    // ],
    per: [
      {
        required: true,
        message: t('settled.form.per.validator'),
      },
    ],
    pictures: [
      {
        required: true,
        message: t('settled.form.pic.validator1'),
      },
    ],
    // hot_name: [
    //   {
    //     required: true,
    //     message: t('settled.form.hot_name.validator'),
    //   },
    // ],
    // hot_pics: [
    //   {
    //     required: true,
    //     message: t('settled.form.pic.validator1'),
    //   },
    // ],
    // hot_desc: [
    //   {
    //     required: true,
    //     message: t('settled.form.hot_desc.validator'),
    //   },
    // ],
  };

  const exampleList = [
    {
      id: 10,
      name: '一起其他餐厅',
      address: 'Xin',
      describe: 'desc',
      status: 0,
      price: 150,
      option: 5.8,
    },
    {
      id: 8,
      name: '遇外滩SKYLINE',
      address: 'Shanghai',
      describe: 'describe',
      status: 0,
      option: 8,
    },
    {
      id: 7,
      name: 'J Prime牛排海鲜',
      address: 'Shanghai',
      describe: 'describe',
      status: 0,
      option: 7,
    },
  ];

  const fetchData = async () => {
    toggleLoading();
    try {
      const {
        data: { code, body },
      } = await SavorApi.userShop();
      if (code === 0) {
        if (body) {
          isView.value = true;
          form.value = {
            status: body?.status,
            post_code: body?.post_code,
            post_name: body?.post_name,
            work_begin_at: body?.work_begin_at,
            work_end_at: body?.work_end_at,
            pictures: body?.pictures,
            pics_url: body?.pics_url,
            name: body?.name,
            address: body?.address,
            describe: body?.describe,
            introduction: body?.introduction,
            phone: body?.phone,
            time: [body?.work_begin_at, body?.work_end_at],
            have_vege: body?.have_vege,
            hostman_pics: body.hostman_pics || [],
            hostman_urls: body.hostman_urls || [],
            leixing: body.leixing || '',
            work_date: body.work_date || [],
            work_date_label: body.work_date
              ? body.work_date.map((item: any) => {
                  const map: any = {};
                  dateArr.forEach((date: any) => {
                    map[date.value] = date.label;
                  });
                  return map[item];
                })
              : [],
            hostman_name: body.hostman_name || '',
            hostman_think: body.hostman_think || '',
            yuyue: body?.yuyue,
            per: body.per || '',
            hot:
              body.hot.map((item: any) => {
                return {
                  ...item,
                  origin_urls: item.hot_pic_urls,
                  hot_pic_urls: item.hot_pic_urls.map((pic: any) => {
                    return { url: pic };
                  }),
                };
              }) || [],
          };
          addressLabel.value = body?.post_name;
          fileList.value =
            body.pics_url &&
            body.pics_url.map((item: any) => {
              return { url: item };
            });
          hostmanList.value =
            body.hostman_urls &&
            body.hostman_urls.map((item: any) => {
              return { url: item };
            });
        } else {
          isView.value = false;
        }
      }
    } finally {
      toggleLoading();
    }
  };

  const customRequest = (option: any, type?: any, index?: any) => {
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
      };
      if (!type) {
        fileList.value = fileList.value.length
          ? fileList.value.concat(curFile)
          : [curFile];
        form.value.pictures = form.value.pictures.length
          ? form.value.pictures.concat(fielName)
          : [fielName];
      } else if (type === 'hostman') {
        hostmanList.value = hostmanList.value.length
          ? hostmanList.value.concat(curFile)
          : [curFile];
        form.value.hostman_pics = form.value.hostman_pics.length
          ? form.value.hostman_pics.concat(fielName)
          : [fielName];
      } else if (type === 'list') {
        const current = form.value.hot[index];
        current.hot_pic_urls = current.hot_pic_urls.length
          ? current.hot_pic_urls.concat(curFile)
          : [curFile];
        current.hot_pics = current.hot_pics.length
          ? current.hot_pics.concat(fielName)
          : [fielName];
        form.value.hot[index] = current;
      }

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

  //
  const uploadFileHost = (option: any) => {
    customRequest(option, 'hostman');
  };

  const uploadFileList = (option: any, index: any) => {
    customRequest(option, 'list', index);
  };

  const beforeUpload = (file: any) => {
    return new Promise((resolve, reject) => {
      const reader: any = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => {
        const image: any = new Image();
        image.src = reader.result;
        image.onload = () => {
          const w: any = image.width;
          const h: any = image.height;
          const aspectRatio: any = w / h;
          if (Math.abs(aspectRatio - 16 / 9) > 0.01) {
            Message.error(t('settled.form.pic.validator'));
            return resolve(false);
          }
          return resolve(true);
        };
      };
    });
  };

  const beforeRemove = (file: any, type?: any, index?: any) => {
    return new Promise((resolve, reject) => {
      Modal.confirm({
        content: t('settled.form.pic.del'),
        onOk: () => {
          if (!type) {
            const newFileList = fileList.value.filter(
              (item: any) => item.url !== file.url
            );
            let pictures = form.value.pictures.filter(
              (item: any) => item !== file.url
            );
            pictures = form.value.pictures.filter(
              (item: any) => file.url.indexOf(item) < 0
            );
            fileList.value = newFileList;
            form.value.pictures = pictures;
          } else if (type === 'hostman') {
            const newFileList = hostmanList.value.filter(
              (item: any) => item.url !== file.url
            );
            let pictures = form.value.hostman_pics.filter(
              (item: any) => item !== file.url
            );
            pictures = form.value.hostman_pics.filter(
              (item: any) => file.url.indexOf(item) < 0
            );
            hostmanList.value = newFileList;
            form.value.hostman_pics = pictures;
          } else if (type === 'list') {
            const current = form.value.hot[index];
            const newFileList = current.hot_pic_urls.filter(
              (item: any) => item.url !== file.url
            );
            let pictures = form.value.hot_pics.filter(
              (item: any) => item !== file.url
            );
            pictures = form.value.hot_pics.filter(
              (item: any) => file.url.indexOf(item) < 0
            );
            form.value.hot[index].hot_pic_urls = newFileList;
            form.value.hot[index].hot_pics = pictures;
          }
          resolve(true);
        },
        onCancel: () => {},
      });
    });
  };

  const beforeRemoveHost = (file: any) => {
    beforeRemove(file, 'hostman');
  };

  const beforeRemoveList = (file: any, index: any) => {
    beforeRemove(file, 'list', index);
  };

  const handleSubmit = async (type: string) => {
    if (!form.value.post_code) {
      Message.error(t('settled.form.province.validator'));
      return;
    }
    await formRef.value.validate(async (valid: any) => {
      if (!valid) {
        if (type === 'save') {
          toggleSaveLoading();
        } else {
          toggleSubLoading();
        }
        try {
          const hotList: any = [];
          form.value.hot.forEach((item: any) => {
            hotList.push({
              hot_name: item.hot_name,
              hot_pics: item.hot_pics,
              hot_desc: item.hot_desc,
            });
          });
          const data: any = {
            name: form.value.name,
            address: form.value.address,
            describe: form.value.describe,
            introduction: form.value.introduction,
            phone: form.value.phone,
            work_begin_at: form.value.time[0],
            work_end_at: form.value.time[1],
            have_vege: form.value.have_vege,
            pictures: form.value.pictures,
            post_code: Array.isArray(form.value.post_code)
              ? form.value.post_code[form.value.post_code.length - 1]
              : form.value.post_code,
            post_name: addressLabel.value,
            leixing: form.value.leixing,
            work_date: form.value.work_date,
            hostman_name: form.value.hostman_name,
            hostman_pics: form.value.hostman_pics,
            hostman_think: form.value.hostman_think,
            yuyue: form.value.yuyue,
            per: form.value.per,
            hot: hotList,
          };
          if (form.value.id) {
            data.id = form.value.id;
          }
          let res: any = {};
          if (type === 'save') {
            res = await SavorApi.saveShop(data);
          } else {
            res = await SavorApi.submitShop(data);
          }
          const {
            data: { code },
          } = res;
          if (code === 0) {
            Message.success(type === 'save' ? '保存成功' : '提审成功');
            fetchData();
            isView.value = true;
          }
        } finally {
          if (type === 'save') {
            toggleSaveLoading();
          } else {
            toggleSubLoading();
          }
        }
      }
    });
  };

  const edit = () => {
    isView.value = false;
  };

  const addressOption = ref<any>([]);
  const initAddress = () => {
    const proData: any = [];
    const add = pcaa['86'];
    Object.keys(add).forEach((p: any) => {
      const children: any = [];
      Object.keys(pcaa[p]).forEach((c: any) => {
        const child: any = [];
        Object.keys(pcaa[c]).forEach((co) => {
          const v = { value: co, label: pcaa[c][co] };
          child.push(v);
        });
        const v: any = { value: c, label: pcaa[p][c] };
        if (child.length > 0) v.children = child;
        children.push(v);
      });
      const proitem: any = { value: p, label: pcaa['86'][p] };
      if (children.length > 0) {
        proitem.children = children;
      }
      proData.push(proitem);
    });

    addressOption.value = proData;
  };

  const handleChange = (value: any) => {
    const add = pcaa['86'];
    const first = add[value[0]];
    const second = pcaa[value[0]][value[1]];
    const third = pcaa[value[1]][value[2]];
    addressLabel.value = [first, second, third].join('');
  };

  const handleAdd = () => {
    form.value.hot.push({
      hot_name: '',
      hot_pics: [],
      hot_pic_urls: [],
      hot_desc: '',
    });
  };

  const handleDelete = (index: any) => {
    form.value.hot.splice(index, 1);
  };

  onMounted(() => {
    fetchData();
    initAddress();
  });
</script>

<style lang="less" scoped>
  .container {
    height: 100%;
    width: 100%;
    color: #fff;
    display: flex;
    flex-direction: column;
    font-family: PingFang SC;

    .content {
      width: 100%;
      background: #20222c;
      position: relative;
      display: flex;
      flex: 1;
      overflow: hidden;

      .form {
        width: 70%;
        padding: 18px 60px;
        overflow-y: auto;

        &.view {
          width: 80%;
          padding: 18px 60px 0px;

          .form {
            padding: 0px;
          }
        }

        .title {
          font-size: 24px;
          font-weight: 600;
          line-height: 42.9px;
          text-align: center;
          display: flex;
          justify-content: space-between;
          margin-bottom: 16px;
        }

        .form-btn {
          width: 100px;
          border: 1px solid #d1b276;
          background: #20222c;
          color: #d1b276;
          font-weight: 600;
          &.primary {
            background: #d1b276;
            color: #20222c;
          }
        }

        :deep(.address) {
          flex-direction: column;
          align-items: flex-start;

          div {
            width: 100%;
            &:first-child {
              margin-bottom: 5px;
            }
          }
        }
        :deep(.arco-select-view-single),
        :deep(.arco-select-view-multiple),
        :deep(.arco-btn) {
          background: #272a37;
          color: #fff;
          border: 1px solid #c3c3c3;
        }

        .dateWrap {
          display: flex;
          flex-direction: row;

          :deep(.arco-form-item) {
            width: 50%;
          }
        }
      }

      .divider {
        height: 100%;
        border-color: #272a37;
      }

      .example {
        flex: 1;
        padding: 18px 48px 18px 100px;

        .title {
          font-size: 24px;
          font-weight: 600;
          line-height: 42.9px;
        }
      }

      .right {
        flex: 1;
        background: #272a37;
        display: flex;
        flex-direction: column;
        padding-top: 60px;
        padding-left: 30px;
        font-size: 16px;

        .operate {
          margin-top: 20px;
        }
        :deep(.arco-btn) {
          background: #272a37;
          border: 1px solid #fff;
          color: #fff;
        }
      }
    }

    .mobile {
      display: block;
      overflow-y: auto;

      .form.view {
        width: 100%;
        padding: 18px 16px 0px;
      }

      .form {
        width: 100%;
        padding: 18px 16px;
        :deep(.arco-form-item-label-col) {
          padding-right: 20px;
        }

        .dateWrap {
          display: block;
          :deep(.arco-form-item) {
            width: 100%;
          }
        }
      }

      .right {
        padding: 16px 0px 16px 30px;
      }

      .example,
      .divider {
        display: none;
      }
    }

    :deep(.arco-spin) {
      margin-top: 60px;
    }
  }
</style>
