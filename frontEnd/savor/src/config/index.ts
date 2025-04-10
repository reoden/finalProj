export const TABLE_BORDERED = {
  wrapper: true,
  cell: false,
  headerCell: true,
  bodyCell: false
};

export const BASE_PAGINATION = {
  page: 1,
  pageSize: 20,
  showTotal: true,
  showPageSize: true,
  pageSizeOptions: [10, 20, 50, 100]
};

export type DrawerControll = 'edit' | 'view' | 'add' | 'copy' | 'wx_edit' | 'upload' | 'auth' | 'bind' | 'uploadOne' | 'uploadError'; 
export type winControll = {
  visible: boolean;
  operaType: DrawerControll;
  rowData: Record<string, any>;
  pointer?: number;
};

export const API_PREFIX = '/zeus';
export const API_PREFIX_PRO = '/ads_center';

export const STATUS_MAP: any = {
  0: '待审核',
  1: '审核未通过',
  2: '审核通过',
  3: '已删除',
  4: '草稿',
  5: '存在敏感词',
  6: '图像虚假',
  7: '文本不匹配',
  8: '图文不匹配',
}