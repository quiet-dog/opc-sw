import { NButton } from "naive-ui";
import { h, ref, defineEmits } from "vue";
import http from "../../http";

export function useOpcHook() {

    const emit = defineEmits(["refresh"])

    const columns = ref([
        {
            title: "节点ID",
            key: "nodeId",
        },
        {
            title: "参数",
            key: "param",
        },
        {
            title: "描述",
            key: "description",
        },
        {
            title: "类型",
            key: "type",
        },
        {
            title: "当前值",
            key: "value",
        },
        {
            title: "设备类型",
            key: "deviceType",
        },
        {
            title: "设备名称",
            key: "deviceName",
        },
        {
            title: "BMS设备名称",
            key: "bmsDeviceName",
        },
        {
            title: "BMS区域",
            key: "bmsArea",
        },
        {
            title: "BMS标签",
            key: "bmsLabel",
        },
        {
            title: "EMS区域",
            key: "emsArea",
        },
        {
            title: "对应key",
            key: "key",
        },
        {
            title: "值时间",
            key: "time",
        },
        {
            title: "扩展",
            key: "extend",
        },
        {
            title: "操作",
            key: "action",
            render(row: any) {
                return h(NButton, {
                    text: true,
                    type: "error",
                    onClick: () => {
                        http.post("/node/delete", {
                            id: row.ID
                        }).then(res => {
                            emit("refresh")
                        })
                    }
                }, "删除")
            }
        }
    ])

    return {
        columns
    }
}