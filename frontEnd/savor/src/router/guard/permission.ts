import type { Router } from 'vue-router';
import NProgress from 'nprogress'; // progress bar

import usePermission from '@/hooks/permission';
import { NOT_FOUND } from '../constants';

export default function setupPermissionGuard(router: Router) {
  router.beforeEach(async (to: any, from, next) => {
    let isMoblie = false;
    if ((navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i))) {
      isMoblie = true;
    } 
    const Permission = usePermission();
    const permissionsAllow = Permission.accessRouter(to);
   
    // eslint-disable-next-line no-lonely-if
    if (permissionsAllow) {
      // console.log('permissionsAllow', permissionsAllow, isMoblie, to);
      // if (isMoblie) {
      //   if (['/mobile', '/mobile/list', '/mobile/detail'].includes(to.path)) {
      //     next();
      //   } else {
      //     next({ name: 'mobileIndex' });
      //   }
      // } else if (['/mobile','/mobile/list', '/mobile/detai;'].includes(to.path)) {
      //   next({ name: 'home' });
      // } else {
      //   console.log('permissionsAllow1', to);
      //   next();
      // }
      next();
    } else {
      const destination =
        // Permission.findFirstPermissionRoute(appRoutes, userStore.role) ||
        NOT_FOUND;
      console.log(destination);
      next(destination);
    }
    NProgress.done();
  });
}
