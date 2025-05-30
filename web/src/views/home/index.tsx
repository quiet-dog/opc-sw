import axios from "axios";
import { MenuOption, NInput, useDialog } from "naive-ui";
import { h, onMounted, ref } from "vue";
import http from "../../http";
import { RouterLink } from "vue-router";

export function useHomeHook() {
    const dialog = useDialog()
    let socket = new WebSocket(`ws://${location.host}/ws`)

    socket.onopen = () => {
        console.log("连接成功")
    }
    socket.onmessage = (e) => {
        // console.log("数据", e.data)
        let data = JSON.parse(e.data)
        table.value.forEach((item: any) => {
            if (item.nodeId === data.nodeId) {
                item.value = data.value
                item.time = data.sourceTimestamp
            }
        })
    }
    socket.onclose = () => {
        socket = new WebSocket(`ws://${location.host}/ws`)
    }
    socket.onerror = () => {
        socket = new WebSocket(`ws://${location.host}/ws`)
    }


    const menus = ref([

    ])

    const table = ref([])

    const selectValue = ref()


    const addMenu = () => {
        const input = ref("");
        dialog.create({
            title: "添加opc服务器",
            content: () => {
                return h(NInput, {
                    placeholder: "请输入opc服务器地址",
                    modelValue: input.value,
                    "onUpdate:value": (val: string) => {
                        input.value = val
                    }
                })
            },
            positiveText: "确定",
            onPositiveClick: async () => {
                await http.post("/service", {
                    opc: input.value
                })
                getMenus()
            }
        })
    }

    const addNode = () => {
        const nodeId = ref("")
        const param = ref("")
        const description = ref("")
        const extend = ref("")
        const serviceId = selectValue.value
        const key = ref("")
        const deviceType = ref("")
        const deviceName = ref("")
        const bmsDeviceName = ref("")
        const bmsArea =ref("")
        const bmsLabel = ref("")
        const emsAare = ref("")

        dialog.create({
            title: "添加节点",
            content: () => {
                return h("div", [
                    h(NInput, {
                        placeholder: "请输入节点id",
                        modelValue: nodeId.value,
                        "onUpdate:value": (val: string) => {
                            nodeId.value = val
                        }
                    }),
                    h(NInput, {
                        placeholder: "请输入参数",
                        modelValue: param.value,
                        "onUpdate:value": (val: string) => {
                            param.value = val
                        }
                    }),
                    // h(NInput, {
                    //     placeholder: "请输入扩展信息",
                    //     modelValue: extend.value,
                    //     "onUpdate:value": (val: string) => {
                    //         extend.value = val
                    //     }
                    // }),
                    // h(NInput, {
                    //     placeholder: "请输入描述",
                    //     modelValue: description.value,
                    //     "onUpdate:value": (val: string) => {
                    //         description.value = val
                    //     }
                    // })
                    h(NInput, {
                        placeholder: "请输入设备类型",
                        modelValue: deviceType.value,
                        "onUpdate:value": (val: string) => {
                            deviceType.value = val
                        }
                    }),

                    h(NInput, {
                        placeholder: "请输入设备名称",
                        modelValue: deviceName.value,
                        "onUpdate:value": (val: string) => {
                            deviceName.value = val
                        }
                    }),

                    h(NInput, {
                        placeholder: "请输入BMS设备名称",
                        modelValue: bmsDeviceName.value,
                        "onUpdate:value": (val: string) => {
                            bmsDeviceName.value = val
                        }
                    }),
                    h(NInput, {
                        placeholder: "请输入BMS设备标签",
                        modelValue: bmsLabel.value,
                        "onUpdate:value": (val: string) => {
                            bmsLabel.value = val
                        }
                    }),
                    h(NInput, {
                        placeholder: "请输入BMS区域",
                        modelValue: bmsArea.value,
                        "onUpdate:value": (val: string) => {
                            bmsArea.value = val
                        }
                    }),
                    h(NInput, {
                        placeholder: "请输入EMS区域",
                        modelValue: emsAare.value,
                        "onUpdate:value": (val: string) => {
                            emsAare.value = val
                        }
                    }),
                    h(NInput, {
                        placeholder: "请输入key",
                        modelValue: key.value,
                        "onUpdate:value": (val: string) => {
                            key.value = val
                        }
                    }),
                ])
            },
            positiveText: "确定",
            onPositiveClick: async () => {
                await http.post("/node", {
                    nodeId: nodeId.value,
                    param: param.value,
                    description: description.value,
                    serviceId,
                    extend: extend.value
                })
                getNodes()
            }
        })
    }

    const getMenus = () => {
        http.get("/service").then(res => {
            const result: any = []
            res.data?.forEach((item: any) => {
                result.push({
                    label: item.opc,
                    value: item.ID,
                })
            })
            menus.value = result
        })
    }

    const getNodes = () => {
        http.post("/node/list", {
            serviceId: selectValue.value
        }).then(res => {
            console.log("res.data", res.data)
            table.value = res.data
        })
    }

    const changeSelect = (value: any) => {
        getNodes()
    }

    onMounted(() => {
        getMenus()
        getNodes()
    })



    return {
        menus,
        addMenu,
        selectValue,
        changeSelect,
        addNode,
        table,
        getNodes
    }
}