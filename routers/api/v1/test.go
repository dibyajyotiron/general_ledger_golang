package v1

import (
	"math"
	"net/http"
	"strconv"

	"general_ledger_golang/pkg/app"
	"general_ledger_golang/pkg/e"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func TestAppStatus(c *gin.Context) {
	appGin := app.Gin{C: c}
	// SysInfo saves the basic system information
	type SysInfo struct {
		Response string      `json:"response"`
		HostStat interface{} `json:"host_stat"`
		CPUStat  interface{} `json:"cpu_stat"`
		VmStat   interface{} `json:"vm_stat"`
		DiskStat interface{} `json:"disk_stat"`
	}
	isCpuInfoNeeded := false
	isHostInfoNeeded := false
	isVmInfoNeeded := false
	isDiskInfoNeeded := false

	cpuInfo := c.Query("cpu_info")
	hostInfo := c.Query("host_info")
	vmInfo := c.Query("vm_info")
	diskInfo := c.Query("disk_info")

	if cpuInfo != "" {
		var err error
		isCpuInfoNeeded, err = strconv.ParseBool(cpuInfo)
		if err != nil {
			appGin.Response(http.StatusInternalServerError, e.INVALID_PARAMS, map[string]bool{"success": false})
			return
		}
	}
	if hostInfo != "" {
		var err error
		isHostInfoNeeded, err = strconv.ParseBool(hostInfo)
		if err != nil {
			appGin.Response(http.StatusInternalServerError, e.INVALID_PARAMS, map[string]bool{"success": false})
			return
		}
	}
	if diskInfo != "" {
		var err error
		isDiskInfoNeeded, err = strconv.ParseBool(diskInfo)
		if err != nil {
			appGin.Response(http.StatusInternalServerError, e.INVALID_PARAMS, map[string]bool{"success": false})
			return
		}
	}
	if vmInfo != "" {
		var err error
		isVmInfoNeeded, err = strconv.ParseBool(vmInfo)
		if err != nil {
			appGin.Response(http.StatusInternalServerError, e.INVALID_PARAMS, map[string]bool{"success": false})
			return
		}
	}

	info := new(SysInfo)
	info.Response = "Test Route!"

	if isCpuInfoNeeded {
		cpuStat, err := cpu.Info()
		if err != nil {
			panic(err)
		}
		info.CPUStat = cpuStat
	}
	if isHostInfoNeeded {
		hostStat, _ := host.Info()
		info.HostStat = hostStat
	}
	if isVmInfoNeeded {
		vmStat, _ := mem.VirtualMemory()
		vmStat.Total = vmStat.Total / uint64(math.Pow(1024, 3))
		vmStat.Used = vmStat.Used / uint64(math.Pow(1024, 3))
		info.VmStat = vmStat
	}
	if isDiskInfoNeeded {
		diskStat, _ := disk.Usage("/")
		info.DiskStat = diskStat
	}

	if isCpuInfoNeeded || isDiskInfoNeeded || isHostInfoNeeded || isVmInfoNeeded {
		info.Response = "System related information fetched!"
	}
	appGin.Response(http.StatusOK, e.SUCCESS, info)
	return
}
