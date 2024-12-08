<template>
  <MobileHeader v-if="mobile" :back="back"/>
  <div v-else class="header-pc">
    <Logo />
    <Navigation :lang="lang" />
    <Nav @change="changelang" />
  </div>
</template>

<script lang="ts" setup>
  import { computed, ref } from 'vue';
  import Logo from '@/components/logo/index.vue';
  import Navigation from '@/components/navigation/index.vue';
  import { isMoblie }from '@/utils';
  import Nav from './nav.vue';
  import MobileHeader from './mobile.vue';

  const mobile = computed(() => isMoblie());
  const lang = ref(false);
  const props = defineProps({
    back: {
      type: Boolean,
      default() {
        return false;
      },
    }
  });

  const changelang = () => {
    lang.value = !lang.value
  }
</script>

<style scoped lang="less">
  .header-pc {
    height: 80px;
    display: flex;
    padding: 0px 36px;
    z-index: 2;
    justify-content: space-between;
    background-color: #272A37;
  }
</style>
