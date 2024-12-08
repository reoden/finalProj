import axios from 'axios';
import qs from 'query-string';
import { API_PREFIX } from '@/config/index';

const getCode = (data: any) => {
  return axios.post('restaurant/v1/sendcode', data)
}

const login = (data: any) => {
  return axios.post('/restaurant/v1/login', data)
}

const loginAdmin = (data: any) => {
  return axios.post('/restaurant/admin/login', data)
}

const getUserInfo = () => {
  return axios.get('restaurant/v1/self', {
    paramsSerializer: (obj) => {
      return qs.stringify(obj);
    },
  });
}

const uploadImage = (data: any) => {
  return axios.post('/restaurant/v1/apps/pic', data)
}

// 商家入驻提审
const submitShop = (data?: any) => {
  return axios.post('/restaurant/admin/apps/submit', data);
}

// 商家入驻保存
const saveShop = (data?: any) => {
  return axios.post('/restaurant/admin/apps/save', data);
}

// 商家餐厅信息
const userShop = () => {
  return axios.get('/restaurant/v1/apps/self');
}

// 首页餐厅展示列表数据

/**
 * 后台接口
 */
// 后台->商家内容审核
const shopList = (params?: any) => {
  return axios.get('/restaurant/admin/apps', {
    params,
    paramsSerializer: (obj) => {
      return qs.stringify(obj);
    },
  });
}

// 后台->商家审核信息
const shopDetail = (id: any) => {
  return axios.get(`/restaurant/v1/apps/${id}`);
}

// 后台->审核
const examineShop = (id: any, status: string) => {
  // todo 该接口有问题，审核不了
  return axios.put(`/restaurant/admin/apps/${id}`, { status });
}

// 展示位搜索商家
const searchShop = (params: any) => {
  return axios.get('/restaurant/v1/apps/search', { params });
}

// 展示位列表
const positionList = () => {
  return axios.get('/restaurant/v1/banner', {
    paramsSerializer: (obj) => {
      return qs.stringify(obj);
    },
  });
}

// 展示位排序
const changePosition = (params?: any) => {
  return axios.put('/restaurant/admin/banner', params);
}

// 增加展示位商家
const addPosition = (data: any) => {
  return axios.post('/restaurant/admin/banner', data);
}

// 删除展示位商家
const deletePosition = (id: any) => {
  return axios.delete(`/restaurant/admin/banner/${id}`);
}

// 厨师列表
const chefList = () => {
  return axios.get('/restaurant/v1/chef', {
    paramsSerializer: (obj) => {
      return qs.stringify(obj);
    },
  });
}

// 厨师排序
const changeChef = (params?: any) => {
  return axios.put('/restaurant/admin/chefShowList', params);
}

// 添加厨师
const addChef = (data: any) => {
  return axios.post('/restaurant/admin/chef/submit', data);
}

// 删除厨师
const deleteChef = (id: any) => {
  return axios.delete(`/restaurant/admin/chef/${id}`);
}

// 厨师详情
const chefDetail = (id: any) => {
  return axios.get(`/restaurant/v1/chef/${id}`);
}

// 嘉宾列表
const guestList = () => {
  return axios.get('/restaurant/v1/guest', {
    paramsSerializer: (obj) => {
      return qs.stringify(obj);
    },
  });
}

// 嘉宾排序
const changeGuest = (params?: any) => {
  return axios.put('/restaurant/admin/guestShowList', params);
}

// 添加嘉宾
const addGuest = (data: any) => {
  return axios.post('/restaurant/admin/guest/submit', data);
}

// 删除嘉宾
const deleteGuest = (id: any) => {
  return axios.delete(`/restaurant/admin/guest/${id}`);
}

// 嘉宾详情
const guestDetail = (id: any) => {
  return axios.get(`/restaurant/v1/guest/${id}`);
}

// 资讯列表
const newsList = () => {
  return axios.get('/restaurant/v1/news', {
    paramsSerializer: (obj) => {
      return qs.stringify(obj);
    },
  });
}

// 资讯排序
const changeNews = (params?: any) => {
  return axios.put('/restaurant/admin/newsShowList', params);
}

// 添加资讯
const addNews = (data: any) => {
  return axios.post('/restaurant/admin/news/submit', data);
}

// 删除资讯
const deleteNews = (id: any) => {
  return axios.delete(`/restaurant/admin/news/${id}`);
}

// 资讯详情
const newsDetail = (id: any) => {
  return axios.get(`/restaurant/v1/news/${id}`);
}

// 大奖列表
const prizeList = (params: any) => {
  return axios.get('/restaurant/v1/prizeShowList', {
    params,
    paramsSerializer: (obj) => {
      return qs.stringify(obj);
    },
  });
}

// 大奖详情
const prizeDetail = (id: any) => {
  return axios.get(`/restaurant/v1/prize/${id}`);
}

const prizeConfig = () => {
  return axios.get('/restaurant/v1/config/self');
}

// 大奖广告展示位列表
const prizeBanner = () => {
  return axios.get('/restaurant/v1/prizeShowList/banner');
}

// banner展示位列表
const staticBanner = () => {
  return axios.get('/restaurant/v1/static_banner_pic_show');
}

// logo列表
const staticLogo = () => {
  return axios.get('/restaurant/v1/logo');
}

export default {
  getCode,
  login,
  loginAdmin,
  getUserInfo,
  uploadImage,

  userShop,
  submitShop,
  saveShop,
  examineShop,

  shopList,
  shopDetail,

  searchShop,
  positionList,
  addPosition,
  changePosition,
  deletePosition,

  chefList,
  changeChef,
  addChef,
  deleteChef,
  chefDetail,

  guestList,
  changeGuest,
  addGuest,
  deleteGuest,
  guestDetail,

  newsList,
  changeNews,
  addNews,
  deleteNews,
  newsDetail,

  prizeConfig,
  prizeList,
  prizeDetail,
  prizeBanner,

  staticBanner,
  staticLogo
}
