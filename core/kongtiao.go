package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sw/global"
	"sw/model/node"
	"sw/opc"
	"time"

	"github.com/go-resty/resty/v2"
)

type KongTiaoDTO struct {
	DeviceSn                       string  `json:"deviceSn"`
	ZhiBanGongKuanYaLiSheDing      float64 `json:"zhiBanGongKuanYaLiSheDing"`
	ZhiBanGongKuanFengLiangSheDing float64 `json:"zhiBanGongKuanFengLiangSheDing"`
	FengFaWenDingZhuangTai         int16   `json:"fengFaWenDingZhuangTai"`
	FaWeiFanKuan                   int16   `json:"faWeiFanKuan"`
	QiangZhiFaWeiSheDing           int16   `json:"qiangZhiFaWeiSheDing"`
	QiangZhiMoShiKaiGuan           int16   `json:"qiangZhiMoShiKaiGuan"`
	PidKongZhiJiFenXiShu           int16   `json:"pidKongZhiJiFenXiShu"`
	PodKongZhiBiLiXiShu            int16   `json:"podKongZhiBiLiXiShu"`
	FengLiangFanKui                int16   `json:"fengLiangFanKui"`
	FangJianShiJiYaLi              float64 `json:"fangJianShiJiYaLi"`
	GongKuangMoShi                 int16   `json:"gongKuangMoShi"`
	ShuangGongKuangQieHuanShiJian  int16   `json:"shuangGongKuangQieHuanShiJian"`
	FengLiangSheDing               int16   `json:"fengLiangSheDing"`
	YaLiSheDing                    float64 `json:"yaLiSheDing"`
}

func InitKongTiao() {
	var kongTiaoNodes []node.NodeModel
	global.DB.Where("device_type = ?", "空调设备").Find(&kongTiaoNodes)

	data := []map[string]interface{}{}
	for _, n := range kongTiaoNodes {
		p := strings.Split(n.Param, "-")
		if len(p) == 2 {
			deviceSn := p[0]
			// 判断data中是否存在deviceSn
			exist := false
			for _, d := range data {
				if d["deviceSn"] == deviceSn {
					exist = true
					break
				}
			}
			if !exist {
				data = append(data, map[string]interface{}{"deviceSn": deviceSn})
			}
		}
	}

	client := resty.New().SetTimeout(3 * time.Second)

	for {
		select {
		case <-context.Background().Done():
			return
		case <-time.After(3 * time.Second):
			{
				for _, d := range kongTiaoNodes {
					p := strings.Split(d.Param, "-")
					if len(p) == 2 {
						deviceSn := p[0]
						key := p[1]
						var msg opc.Data
						global.Redis.Get(global.Ctx, fmt.Sprintf("%d", d.ID)).Scan(&msg)
						// 寻找data相同的sn,赋值key
						for _, v := range data {
							if v["deviceSn"] == deviceSn {
								v[key] = msg.Value
								break
							}
						}
					}
				}

				jsonByte, err := json.Marshal(data)
				if err != nil {
					continue
				}

				client.R().SetHeader("Authorization", "Bearer MASTER_TOKEN_123456").SetBody(jsonByte).Post(fmt.Sprintf("http://%s:%s/manage/kongTiaoData/receive", global.Config.Sw.Host, global.Config.Sw.Port))
			}
		}
	}

}
