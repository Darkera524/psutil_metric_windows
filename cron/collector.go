package cron

import (
	"github.com/StackExchange/wmi"
	"fmt"
	"github.com/open-falcon/common/model"
	"github.com/Darkera524/psutil_metric_windows/g"
	"time"
)

type Win32_PerfRawData_perfProc_Process struct {
	Name string
	IDProcess int
	PercentProcessorTime int
}

func Collect(){
	var interval int64 = g.Config().Transfer.Interval
	var ticker = time.NewTicker(time.Duration(interval) * time.Second)

	for {
		<-ticker.C
		query_wmi()
	}
}

func query_wmi(){
	var dst []Win32_PerfRawData_perfProc_Process
	q := wmi.CreateQuery(&dst, "")
	err := wmi.Query(q, &dst)
	if err != nil {
		err.Error()
	}

	proc_metrics,_ := convirtProcessInfoToMetrics(dst)

	g.SendMetrics(proc_metrics)


}

func convirtProcessInfoToMetrics(procInfo []Win32_PerfRawData_perfProc_Process)(metrics []*model.MetricValue, err error){
	var total_time int
	for _,proc := range procInfo{
		if proc.Name == "_Total"{
			total_time = proc.PercentProcessorTime
		}
	}

	hostname, _ := g.Hostname()
	now := time.Now().Unix()
	var tags string
	var attachtags = g.Config().AttachTags
	var interval int64 = g.Config().Transfer.Interval
	if attachtags != "" {
		tags = attachtags
	}

	for i:=0;i<len(procInfo);i++{
		cmdline := procInfo[i].Name
		cpu_percent := float64(procInfo[i].PercentProcessorTime) / float64(total_time)
		var tag string
		if tags != "" {
			tag = fmt.Sprintf("%s,pid=%d,cmdline=%s", tags, procInfo[i].IDProcess, cmdline)
		} else {
			tag = fmt.Sprintf("pid=%d,cmdline=%s", procInfo[i].IDProcess, cmdline)
		}
		singleMetric := &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "proc.cpu.percent",
			Value:     cpu_percent,
			Timestamp: now,
			Step:      interval,
			Type:      "GAUGE",
			Tags:      tag,
		}
		metrics = append(metrics, singleMetric)

		/*singleMetric = &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "proc.mem.percent",
			Value:     procInfo[i].MemPercent,
			Timestamp: now,
			Step:      interval,
			Type:      "GAUGE",
			Tags:      tag,
		}
		metrics = append(metrics, singleMetric)

		singleMetric = &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "proc.fd.num",
			Value:     procInfo[i].FileDescriptorNum,
			Timestamp: now,
			Step:      interval,
			Type:      "GAUGE",
			Tags:      tag,
		}
		metrics = append(metrics, singleMetric)

		singleMetric = &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "proc.thread.num",
			Value:     procInfo[i].ThreadNum,
			Timestamp: now,
			Step:      interval,
			Type:      "GAUGE",
			Tags:      tag,
		}
		metrics = append(metrics, singleMetric)*/
	}
	return metrics,nil
}
