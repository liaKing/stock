package util

import (
	"fmt"
	"sync"
	"time"
)

const (
	twepoch        = int64(1624914000000) // 起始时间戳，这里使用的是2021-06-29 00:00:00的时间戳
	machineBits    = uint(10)             // 机器ID位数，可根据需求调整
	sequenceBits   = uint(12)             // 序列号位数，可根据需求调整
	maxMachineID   = int64(-1) ^ (int64(-1) << machineBits)
	maxSequenceNum = int64(-1) ^ (int64(-1) << sequenceBits)
)

type Snowflake struct {
	machineID     int64
	sequence      int64
	lastTimestamp int64
	mutex         sync.Mutex
}

func NewSnowflake(machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, fmt.Errorf("machine ID out of range: %d", machineID)
	}
	return &Snowflake{
		machineID: machineID,
		sequence:  0,
	}, nil
}

// Generate 获取uuid
func (s *Snowflake) Generate() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	timestamp := time.Now().UnixNano() / 1e6
	if timestamp < s.lastTimestamp {
		// 时钟回拨，等待时间戳追赶上
		timestamp = s.lastTimestamp
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequenceNum
		if s.sequence == 0 {
			// 当前毫秒内的序列号已用完，等待至下一毫秒
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	id := (timestamp-twepoch)<<machineBits | s.machineID<<sequenceBits | s.sequence
	return id
}
