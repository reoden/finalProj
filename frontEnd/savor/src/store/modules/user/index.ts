import { defineStore } from 'pinia';
import {
  logout as userLogout,
  getUserInfo,
  LoginData,
} from '@/api/user';
import SavorApi from '@/api/savor';
import { setToken, clearToken } from '@/utils/auth';
import { removeRouteListener } from '@/utils/route-listener';
import { UserState } from './types';
import useAppStore from '../app';

const useUserStore = defineStore('user', {
  state: (): UserState => ({
    name: undefined,
    avatar: undefined,
    job: undefined,
    organization: undefined,
    location: undefined,
    email: undefined,
    introduction: undefined,
    personalWebsite: undefined,
    jobName: undefined,
    organizationName: undefined,
    locationName: undefined,
    phone: undefined,
    registrationDate: undefined,
    accountId: undefined,
    certification: undefined,
    role: '',
    showPrize: false,
    menuInfo: {},
    bannerList: [],
    logoSmall: '',
    logoBig: '',
  }),

  getters: {
    userInfo(state: UserState): UserState {
      return { ...state };
    },
  },

  actions: {
    switchRoles() {
      return new Promise((resolve) => {
        this.role = this.role === 'user' ? 'admin' : 'user';
        resolve(this.role);
      });
    },
    // Set user's information
    setInfo(partial: Partial<UserState>) {
      this.$patch(partial);
    },

    // Reset user's information
    resetInfo() {
      this.$reset();
    },

    // Get user's information
    async info() {
      const { data: { body } } : any = await SavorApi.getUserInfo();
      console.log('info', body)
      this.setInfo(body);
    },

    async getPrize() {
      const { data: { body } } : any = await SavorApi.prizeConfig();
      this.setInfo({ showPrize: !!body.prize, menuInfo: body });
    },

    async getBanner() {
      const { data: { body } } : any = await SavorApi.staticBanner();
      this.setInfo({ bannerList: body.apps });
    },

    async getLogo() {
      const { data: { body } } : any = await SavorApi.staticLogo();
      const list = body?.apps || [];
      this.setInfo({
        logoSmall: list.filter((item: any) => Number(item.id) === 1)?.[0]?.pic_url,
        logoBig: list.filter((item: any) => Number(item.id) === 2)?.[0]?.pic_url
       });
    },
    
    // Login
    async login(loginForm: LoginData) {
      try {
        const { data: { body } } = await SavorApi.login(loginForm);
        setToken(body.access_token);
      } catch (err) {
        clearToken();
        throw err;
      }
    },
    logoutCallBack() {
      const appStore = useAppStore();
      this.resetInfo();
      clearToken();
      removeRouteListener();
      appStore.clearServerMenu();
    },
    // Logout
    async logout() {
      try {
        await userLogout();
      } finally {
        this.logoutCallBack();
      }
    },
  },
});

export default useUserStore;
