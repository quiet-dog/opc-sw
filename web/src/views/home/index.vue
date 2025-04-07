<script lang='ts' setup>
import { NButton, NLayout, NLayoutHeader, NLayoutSider, NMenu, NSelect, NSpace } from 'naive-ui';
import { useHomeHook } from './index'
import OpcTable from "../opc/index.vue";
const { menus, selectValue, changeSelect, addNode, table, getNodes } = useHomeHook();

</script>

<template>
    <div>
        <NButton @click="addNode">添加NodeId</NButton>
        <NSpace vertical>
            <NLayout>
                <NLayoutSider has-sider>
                    <NSelect @update-value="changeSelect" v-model:value="selectValue" :options="menus"></NSelect>
                </NLayoutSider>
                <NLayout>
                    <OpcTable @refresh="() => getNodes()" :tableList="table" />
                </NLayout>
            </NLayout>
        </NSpace>
        <div @mouseenter=""></div>

        <!-- <NLayout has-sider>
            <NLayoutSider bordered content-style="padding: 24px;">
                菜单
            </NLayoutSider>
            <NLayout>
                <NLayoutHeader>

                </NLayoutHeader>
            </NLayout>
        </NLayout> -->

    </div>
</template>

<style scoped></style>
