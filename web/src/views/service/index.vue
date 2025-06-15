<script lang='ts' setup>
import { NButton, NCard, NForm, NFormItem, NInput, NMessageProvider, NSelect, NSpace, useMessage } from 'naive-ui';
import { onMounted, ref } from 'vue';
import http from '../../http';
const formRef = ref();
const upformRef = ref();
const form = ref({
    opc: "",
    username: "",
    password: "",
})


const updateForm = ref({
    opc: "",
    username: "",
    password: "",
})
/**
 * http.get("/service").then(res => {
            const result: any = []
            res.data?.forEach((item: any) => {
                result.push({
                    label: item.opc,
                    value: item.ID,
                })
            })
            menus.value = result
        })
 * 
 */

const menus = ref<any[]>([])
const menusValue = ref<any[]>([])
const selectValue = ref()

const getMenus = () => {
    http.get("/service").then(res => {
        const result: any = []
        res.data?.forEach((item: any) => {
            result.push({
                label: item.opc,
                value: item.ID,
            })
        })
        menusValue.value = res.data
        menus.value = result
    })
}

const submit = () => {

    if (form.value.opc == "" || form.value.opc == undefined || form.value.opc == null) {
        alert("opc服务地址不能为空");
        return;
    }

    http.post("/service", {
        opc: form.value.opc,
        username: form.value.username,
        password: form.value.password
    }).then(res => {
        alert("添加成功");
    }).catch(err => {

        alert("添加失败");
    })

}

const update = () => {
    if (updateForm.value.opc == "" || updateForm.value.opc == undefined || updateForm.value.opc == null) {
        alert("opc服务地址不能为空");
        return;
    }
    http.post("/service/update", {
        id: selectValue.value,
        opc: updateForm.value.opc,
        username: updateForm.value.username,
        password: updateForm.value.password
    }).then(res => {
        alert("更新成功");
    }).catch(err => {
        alert("更新失败");
    })
}

const changeSelect = (value: any) => {
    selectValue.value = value;
    const item = menusValue.value.find((item: any) => item.ID == value);
    if (item) {
        updateForm.value.opc = item.opc;
        updateForm.value.username = item.username;
        updateForm.value.password = item.password;
    }
}

const restSys = () => {
    http.get("/service/restart").then(res => {

    }).catch(err => {
    })
}

onMounted(() => {
    getMenus();
})
</script>

<template>
    <div>
        <NSpace vertical>
            <NCard title="opc系统重启" style="width: 300px;">
                <NButton @click="restSys">重启服务</NButton>
            </NCard>

            <!-- <NLayoutSider has-sider> -->
            <NCard title="服务列表" style="width: 300px;">
                <NSelect @update-value="changeSelect" v-model:value="selectValue" :options="menus"></NSelect>
                <NForm v-if="selectValue > 0" ref="upformRef" :model="updateForm">
                    <NFormItem path="opc" label="opc服务地址">
                        <NInput v-model:value="updateForm.opc" />
                    </NFormItem>
                    <NFormItem path="username" label="用户名">
                        <NInput v-model:value="updateForm.username" />
                    </NFormItem>
                    <NFormItem path="password" label="密码">
                        <NInput v-model:value="updateForm.password" />
                    </NFormItem>
                </NForm>
                <NButton v-if="selectValue > 0" @click="update">更新</NButton>
            </NCard>
            <!-- </NLayoutSider> -->
            <NCard title="添加服务" style="width: 300px;">
                <NForm ref="formRef" :model="form">
                    <NFormItem path="opc" label="opc服务地址">
                        <NInput v-model:value="form.opc" />
                    </NFormItem>
                    <NFormItem path="username" label="用户名">
                        <NInput v-model:value="form.username" />
                    </NFormItem>
                    <NFormItem path="password" label="密码">
                        <NInput v-model:value="form.password" />
                    </NFormItem>
                </NForm>
                <NButton @click="submit">提交</NButton>
            </NCard>
        </NSpace>

    </div>
</template>

<style scoped></style>
