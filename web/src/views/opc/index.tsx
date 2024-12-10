import { ref } from "vue";

export function useOpcHook() {
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
        }
    ])

    return {
        columns
    }
}