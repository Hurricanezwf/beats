package clskafka

import (
	"fmt"
	"hash/crc64"
	"math"
	"strings"
	"sync/atomic"
)

func (c *client) generatePackageID() string {
	batchID := atomic.AddUint64(&c.autoIncrBatchID, 1)
	if batchID >= math.MaxUint64-1 {
		batchID = 0
	}
	return strings.ToUpper(fmt.Sprintf("%s-%x", c.producerID, batchID))
}

func generateProducerHash(str string) string {
	table := crc64.MakeTable(crc64.ECMA)
	hash := crc64.Checksum([]byte(str), table)
	hashString := fmt.Sprintf("%016x", hash)
	return strings.ToUpper(hashString)
}
