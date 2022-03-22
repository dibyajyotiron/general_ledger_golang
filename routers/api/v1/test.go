package v1

import (
	"general_ledger_golang/pkg/app"
	"general_ledger_golang/pkg/e"
	"math"
	"net/http"
	"strconv"

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
		Hoststat interface{} `json:"host_stat"`
		CPUStat  interface{} `json:"cpu_stat"`
		VmStat   interface{} `json:"vm_stat"`
		DiskStat interface{} `json:"disk_stat"`
	}
	is_cpu_info_needed := false
	is_host_info_needed := false
	is_vm_info_needed := false
	is_disk_info_needed := false

	cpu_info := c.Query("cpu_info")
	host_info := c.Query("host_info")
	vm_info := c.Query("vm_info")
	disk_info := c.Query("disk_info")

	if cpu_info != "" {
		var err error
		is_cpu_info_needed, err = strconv.ParseBool(cpu_info)
		if err != nil {
			appGin.Response(http.StatusInternalServerError, e.INVALID_PARAMS, map[string]bool{"success": false})
			return
		}
	}
	if host_info != "" {
		var err error
		is_host_info_needed, err = strconv.ParseBool(host_info)
		if err != nil {
			appGin.Response(http.StatusInternalServerError, e.INVALID_PARAMS, map[string]bool{"success": false})
			return
		}
	}
	if disk_info != "" {
		var err error
		is_disk_info_needed, err = strconv.ParseBool(disk_info)
		if err != nil {
			appGin.Response(http.StatusInternalServerError, e.INVALID_PARAMS, map[string]bool{"success": false})
			return
		}
	}
	if vm_info != "" {
		var err error
		is_vm_info_needed, err = strconv.ParseBool(vm_info)
		if err != nil {
			appGin.Response(http.StatusInternalServerError, e.INVALID_PARAMS, map[string]bool{"success": false})
			return
		}
	}

	info := new(SysInfo)
	info.Response = "Test Route!"

	if is_cpu_info_needed {
		cpuStat, err := cpu.Info()
		if err != nil {
			panic(err)
		}
		info.CPUStat = cpuStat
	}
	if is_host_info_needed {
		hostStat, _ := host.Info()
		info.Hoststat = hostStat
	}
	if is_vm_info_needed {
		vmStat, _ := mem.VirtualMemory()
		vmStat.Total = vmStat.Total / uint64(math.Pow(1024, 3))
		vmStat.Used = vmStat.Used / uint64(math.Pow(1024, 3))
		info.VmStat = vmStat
	}
	if is_disk_info_needed {
		diskStat, _ := disk.Usage("/")
		info.DiskStat = diskStat
	}

	if is_cpu_info_needed || is_disk_info_needed || is_host_info_needed || is_vm_info_needed {
		info.Response = "System related information fetched!"
	}
	appGin.Response(http.StatusOK, e.SUCCESS, info)
	return
}
