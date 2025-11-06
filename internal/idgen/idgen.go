// idgen 包提供分布式ID生成功能
// 使用类似雪花算法的ID生成策略，生成64位整数ID
package idgen

import (
	"crypto/md5"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	// 机器ID（0-1023，共10位）
	machineID int64

	// 序列号（0-4095，同一毫秒内的递增序列，共12位）
	sequence int64

	// 互斥锁，保证并发安全
	mu sync.Mutex

	// 上次生成ID的时间戳（毫秒）
	lastTimestamp int64

	// 起始时间戳（2024-01-01 00:00:00 UTC的毫秒时间戳）
	epoch int64 = 1704067200000

	// 机器ID的位数（10位，支持1024台机器）
	machineIDBits int64 = 10

	// 序列号的位数（12位，每毫秒可生成4096个ID）
	sequenceBits int64 = 12

	// 机器ID左移位数
	machineIDShift int64 = sequenceBits

	// 时间戳左移位数
	timestampShift int64 = sequenceBits + machineIDBits

	// 序列号掩码（4095）
	sequenceMask int64 = -1 ^ (-1 << sequenceBits)

	// 最大机器ID（1023）
	maxMachineID int64 = -1 ^ (-1 << machineIDBits)
)

func init() {
	// 初始化机器ID
	initMachineID()
}

// initMachineID 初始化机器ID
// 优先从环境变量获取，如果没有则根据主机名和IP地址生成
func initMachineID() {
	// 尝试从环境变量获取机器ID
	if envID := os.Getenv("MACHINE_ID"); envID != "" {
		var id int64
		if _, err := fmt.Sscanf(envID, "%d", &id); err == nil && id >= 0 && id <= maxMachineID {
			machineID = id
			return
		}
	}

	// 根据主机名和IP地址生成机器ID
	hostname, _ := os.Hostname()
	// 使用MD5哈希生成一个稳定的机器ID
	hash := md5.Sum([]byte(hostname))
	// 取hash的前4个字节，转换为int64，然后取模
	machineID = int64(hash[0])<<24 | int64(hash[1])<<16 | int64(hash[2])<<8 | int64(hash[3])
	if machineID < 0 {
		machineID = -machineID
	}
	machineID = machineID % (maxMachineID + 1)
}

// GenerateUserID 生成用户ID
// 返回64位整数ID，格式：时间戳(42位) + 机器ID(10位) + 序列号(12位)
// 特点：
// - 全局唯一
// - 趋势递增（按时间排序）
// - 包含时间信息
// - 支持分布式环境
func GenerateUserID() int64 {
	mu.Lock()
	defer mu.Unlock()

	// 获取当前时间戳（毫秒）
	now := time.Now().UnixMilli()

	// 如果当前时间小于上次时间戳，说明时钟回拨，等待时钟追上
	if now < lastTimestamp {
		// 等待时钟追上
		time.Sleep(time.Duration(lastTimestamp-now) * time.Millisecond)
		now = time.Now().UnixMilli()
	}

	// 如果是同一毫秒内，序列号递增
	if now == lastTimestamp {
		sequence = (sequence + 1) & sequenceMask
		// 如果序列号溢出（超过4095），等待下一毫秒
		if sequence == 0 {
			now = waitNextMillis(lastTimestamp)
		}
	} else {
		// 新的毫秒，序列号重置为0
		sequence = 0
	}

	// 更新上次时间戳
	lastTimestamp = now

	// 生成ID：时间戳(42位) + 机器ID(10位) + 序列号(12位)
	// 时间戳从epoch开始计算，减少ID长度
	timestamp := now - epoch

	id := (timestamp << timestampShift) |
		(machineID << machineIDShift) |
		sequence

	return id
}

// waitNextMillis 等待下一毫秒
func waitNextMillis(lastTimestamp int64) int64 {
	now := time.Now().UnixMilli()
	for now <= lastTimestamp {
		now = time.Now().UnixMilli()
	}
	return now
}

// ParseUserID 解析用户ID，返回时间戳、机器ID和序列号
// 用于调试和分析
func ParseUserID(id int64) (timestamp int64, machineID int64, sequence int64) {
	timestamp = (id >> timestampShift) + epoch
	machineID = (id >> machineIDShift) & maxMachineID
	sequence = id & sequenceMask
	return timestamp, machineID, sequence
}
