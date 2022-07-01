package auto_plan

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type AutoSendJob struct {
	ID            int       `json:"id"`              // 任务ID
	JobCode       string    `json:"job_code"`        // 任务号
	Status        int       `json:"status"`          // 状态
	StatusResult  string    `json:"status_result"`   // 任务说明
	CreatedByName string    `json:"created_by_name"` // 创建人
	CreatedAt     time.Time `json:"created_at"`      // 创建时间
	ExecTime      int       `json:"exec_time"`       // 执行时长
}

type AutoSendJobResult struct {
	JobResultId int    `json:"job_result_id"`
	JobCode     string `json:"job_code"`   // 派车任务号
	JobResult   string `json:"job_result"` // 自动派车任务结果
	CarTypes    string `json:"car_types"`  // 车型信息
}

// JobResultPlan 派车任务结果列表元素（派车单）
type JobResultPlan struct {
	PlanId          int                 `json:"plan_id"`      // 派车单ID
	PlanCode        string              `json:"plan_code"`    // 派车单号
	IsLimitRun      int                 `json:"is_limit_run"` // 是否限行 1：限行 其他：不限行
	IsDistanceValid int                 `json:"is_distance_valid"`
	CarId           int                 `json:"car_id"`         // 车辆ID
	CarNumber       string              `json:"car_number"`     // 车牌号
	CarTypeName     string              `json:"car_type_name"`  // 车型
	DriverId        int                 `json:"driver_id"`      // 司机ID
	DriverName      string              `json:"driver_name"`    // 司机姓名
	FullRate        float64             `json:"full_rate"`      // 满载率
	FullRateType    int                 `json:"full_rate_type"` // 满载率类型
	WaybillIds      []int               `json:"waybill_ids"`    // 运单ID数组
	WaybillList     []*JobResultWaybill `json:"waybill_list"`   // 运单数组
	LogisticsIds    []int               `json:"logistics_ids"`  // 所有的物流单ID
}

// JobResultWaybill 派车任务结果列表元素（运单）
type JobResultWaybill struct {
	WaybillOrderId   int    `json:"waybill_order_id"`
	WaybillOrderCode string `json:"waybill_order_code"`
	Weight           int    `json:"weight"`
	Volume           int64  `json:"volume"`
	TotalThickness   int    `json:"total_thickness"`
	MaxLength        int    `json:"max_length"`
	Lmt              int    `json:"lmt"`
	UseCopilot       int    `json:"use_copilot"`
	LogisticsIds     []int  `json:"logistics_ids"`
	TotalPrice       int64  `json:"total_price"`
}

func PlanResult() {
	db, err := gorm.Open("mysql", "sizu:uuC4W9YyZ6@tcp(am-bp1y98h053h3267bh167320o.ads.aliyuncs.com:3306)/xpx_tms?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		print(err)
	}
	startTime := "2021-12-01 00:00:00"
	endStart := "2022-01-01 00:00:00"
	sql := "SELECT id,job_code,`status`,IF(`status` = 3, '有效', '无效' ) as status_result,created_by_name,created_at,timestampdiff(SECOND,created_at,finish_at) as exec_time FROM tms_auto_send_job WHERE created_at BETWEEN ? AND ? AND is_deleted = 1 AND status = 3 ORDER BY created_at ASC"
	autoSendJob := make([]*AutoSendJob, 0)
	db.Raw(sql, startTime, endStart).Scan(&autoSendJob)
	f, err1 := os.OpenFile("auto_info_result", os.O_RDWR, 0666)
	check(err1)
	info := "任务ID" + "\t" + "任务号" + "\t" + "状态" + "\t" + "任务说明" + "\t" + "创建人" + "\t" + "创建时间" + "\t" + "执行时长/秒" + "\t" + "派车单" + "\t" + "满载率" + "\t" + "运单数" + "\t" + "车牌号"
	_, err1 = io.WriteString(f, info) //写入文件(字符串)
	for _, v := range autoSendJob {
		var writeString = ""
		createdA := v.CreatedAt.Format("2006-01-02 15:04:05")
		writeString = strconv.Itoa(v.ID) + "\t" + v.JobCode + "\t" + strconv.Itoa(v.Status) + "\t" + v.StatusResult + "\t" + v.CreatedByName + "\t" + createdA + "\t" + strconv.Itoa(v.ExecTime)
		sql := "SELECT job_result FROM tms_auto_send_job_result WHERE job_code = ?"
		var info string
		var autoSendJobResult AutoSendJobResult
		planInfos := make([]*JobResultPlan, 0)
		db.Raw(sql, v.JobCode).Scan(&autoSendJobResult)
		_ = json.Unmarshal([]byte(autoSendJobResult.JobResult), &planInfos)
		for _, v := range planInfos {
			count := len(v.WaybillIds)
			info = writeString + "\t" + v.PlanCode + "\t" + fmt.Sprintf("%.2f", v.FullRate) + "\t" + strconv.Itoa(count) + "\t" + v.CarNumber + "\n"
			_, err1 := io.WriteString(f, info) //写入文件(字符串)
			check(err1)
		}
	}
	f.Close()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
