// import { appRoutes } from '../routes';

const mixinRoutes: any = [];

const appClientMenus = mixinRoutes.map(
  (el: { name: any; path: any; meta: any; redirect: any; children: any }) => {
    const { name, path, meta, redirect, children } = el;
    return {
      name,
      path,
      meta,
      redirect,
      children,
    };
  }
);

export default appClientMenus;
