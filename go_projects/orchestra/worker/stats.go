package worker

import (
	"log"

	"github.com/c9s/goprocinfo/linux"
)

type CPUStat struct {
	Id string `json:"id"`
	User uint64 `json:"user"`
	Nice uint64 `json:"nice"`
	System uint64 `json:"system"`
	Idle uint64 `json:"idle"`
	IOWait uint64 `json:"iowait"`
	IRQ uint64 `json:"irq"`
	SoftIRQ uint64 `json:"softirq"`
	Steal uint64 `json:"steal"`
	Guest uint64 `json:"guest"`
	GuestNice uint64 `json:"guest_nice"`
}



//this struct keeps of the imp stats of the node the worker is running on
type Stats struct {
	MemStats *linux.MemInfo
	DiskStats *linux.Disk
	CpuStats *linux.CPUStat
	LoadStats *linux.LoadAvg
}

//get the total memory of the worker machine
//@todo: remember to add the kb at the end off the num
func (s *Stats) MemTotalKb() uint64 {
	return s.MemStats.MemTotal
}

func (s *Stats) MemAvailbeKb() uint64 {
	return s.MemStats.MemAvailable
}

func (s *Stats) MemUsedKb() uint64 {
	return s.MemStats.MemTotal - s.MemStats.MemAvailable
}

func (s *Stats) MemUsedPercent() uint64 {
	return s.MemStats.MemAvailable / s.MemStats.MemTotal
}


//steps to change cpu usage
//1. Sum values of idle states
//2. sum values of the non-idle states
//3. sum of total idle and non-idle states
//4. sub the idle from total and divide the result by the total
func (s *Stats) CpuUsage() float64 {

	idle := s.CpuStats.Idle + s.CpuStats.IOWait
	nonIdle := s.CpuStats.User + s.CpuStats.Nice + s.CpuStats.System + s.CpuStats.IRQ + s.CpuStats.SoftIRQ + s.CpuStats.Steal
	total := idle + nonIdle

	if total == 0 {
		return 0.00
	}

	return (float64(total) - float64(idle)) / float64(total)
}

//helper functions to fetch stats
func GetMemoryInfo() *linux.MemInfo {
	memstats, err := linux.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Printf("error reading from /proc/meminfo \n")
		return &linux.MemInfo{}
	}

	return memstats
}

func GetDiskStats() *linux.Disk {
	disk, err := linux.ReadDisk("/")
	if err != nil {
		log.Printf("error while reading from / \n")
		return &linux.Disk{}
	}
	return disk
}

func GetCpuStats() *linux.CPUStat {
	cpustat, err := linux.ReadStat("/proc/stat")
	if err != nil {
		log.Printf("error while getting cpu stats \n")
		return &linux.CPUStat{}
	}

	return &cpustat.CPUStatAll
}

func GetLoadStats() *linux.LoadAvg {
	loadavg, err := linux.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		log.Printf("error while reading form /proc/loadavg")
		return &linux.LoadAvg{}
	}
	return loadavg
}

func GetStats() *Stats {
	return &Stats{
		MemStats: GetMemoryInfo(),
		DiskStats: GetDiskStats(),
		CpuStats: GetCpuStats(),
		LoadStats: GetLoadStats(),
	}
}


