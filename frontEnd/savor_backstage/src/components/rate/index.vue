<template>
  <div>
    <icon-star-fill v-for="(itm, idx) in Math.ceil(rate)" :key="idx" />
    <span v-if="showText" style="margin-left: 5px;">{{ rate }}</span>
  </div>
</template>

<script lang="ts" setup>
  const props = defineProps({
    rate: {
      type: Object as any,
      default() {
        return 0;
      },
    },
    showText: {
      type: Boolean,
      default() {
        return true;
      },
    }
  });
</script>

<style scoped lang="less">
</style>
