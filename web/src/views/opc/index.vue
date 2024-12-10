<script lang='ts' setup>
import { NDataTable } from 'naive-ui';
import { useOpcHook } from '.';
const { columns } = useOpcHook();
 
const { tableList }= defineProps<{
    tableList: any[]
}>()

</script>

<template>
  <div>
    <NDataTable :data="tableList" :columns="columns" />
  </div>
</template>

<style scoped>

</style>
