import { NButton } from "naive-ui";
import { h, ref } from "vue";
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
            title: "当前值",
            key: "value",
        },
        {
            title: "值时间",
            key: "time",
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
                })
            }
        }
    ])

    return {
        columns
    }
}