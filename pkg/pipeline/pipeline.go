package pipeline

import (
	"time"

	"github.com/EdmundMartin/trackpix/pkg/server"
)

// DataRecord are the objects to be saved in Clickhouse
type DataRecord struct {
	IPAddr      string
	VisitedPage string
	Datetime    time.Time
	UserAgent   string
	Language    string
}

// TrackingPipeline defines the interface required by a data pipeline
type TrackingPipeline interface {
	Run(chan *server.TrackingData, chan *DataRecord)
}

// TrxPipe takes records from server into Clickhouse suitable format
type TrxPipe struct {
	Functions []func(*server.TrackingData, *DataRecord) *DataRecord
}

// Run takes events from pixel and pushes them through data pipeline
func (trxP TrxPipe) Run(events chan *server.TrackingData, dataObjs chan *DataRecord) {
	for {
		event := <-events
		record := &DataRecord{}
		for _, function := range trxP.Functions {
			record = function(event, record)
		}
		dataObjs <- record
	}
}
