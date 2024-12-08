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
      if (isMoblie) {
        console.log(to);
        if (['mobileList', 'mobileDetail'].includes(to.name)) {
          next();
        } else {
          next({ name: 'mobileList' });
        }
      } else {
        next();
      }
    } else {
      const destination =
        // Permission.findFirstPermissionRoute(appRoutes, userStore.role) ||
        NOT_FOUND;
      next(destination);
    }
    NProgress.done();
  });
}
