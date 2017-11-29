package device

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/toni-moreno/snmpcollector/pkg/agent/selfmon"
)

// DevStatType a device stat type
type DevStatType uint

const (
	//SnmpGetQueries num Get Queries on last gather cicle
	SnmpGetQueries = 0
	//SnmpWalkQueries num Walk Queries on last gather cicle
	SnmpWalkQueries = 1
	// SnmpGetErrors num Get Errors on last gather cicle
	SnmpGetErrors = 2
	// SnmpWalkErrors num Walk Errors on last gather cicle
	SnmpWalkErrors = 3
	// SnmpQueryTimeouts timeous happened while doing snmp queries
	SnmpQueryTimeouts = 4
	// SnmpOIDGetAll Snmp  all OID based gathered metrics
	SnmpOIDGetAll = 5
	// SnmpOIDGetProcessed only those which match filters
	SnmpOIDGetProcessed = 6
	// SnmpOIDGetErrors OIDs wich has errors
	SnmpOIDGetErrors = 7
	// EvalMetricsAll all Evaluated Metrics
	EvalMetricsAll = 8
	// EvalMetricsOk evaluated OK
	EvalMetricsOk = 9
	// EvalMetricsErrors those evalutaed metrics with some errors
	EvalMetricsErrors = 10
	// MetricSent all values had been sent (measurment fields -- could be from OID's or from computed, evaluated, sources)
	MetricSent = 11
	// MetricSentErrors values that has errors when trying to add to a measurement
	MetricSentErrors = 12
	// MeasurementSent all measurements sent to the influx backend
	MeasurementSent = 13
	// MeasurementSentErrors all measurements with errors
	MeasurementSentErrors = 14
	// CicleGatherStartTime Time which begins the last Gather Cicle
	CicleGatherStartTime = 15
	// CicleGatherDuration Time taken in complete the last gather and sent cicle
	CicleGatherDuration = 16
	// FilterStartTime Time which begins the last filter update
	FilterStartTime = 17
	// FilterDuration Time taken in complete the filtering process
	FilterDuration = 18
	// BackEndSentStartTime Time witch begins the last sent process
	BackEndSentStartTime = 19
	// BackEndSentDuration Time taken in complete the data sent process
	BackEndSentDuration = 20
	// DevStatTypeSize special value to set the last stat position
	DevStatTypeSize = 21
)

// DevStat minimal info to show users
type DevStat struct {
	//ID
	id     string
	TagMap map[string]string
	//Control
	log     *logrus.Logger
	selfmon *selfmon.SelfMon
	mutex   sync.Mutex

	//Counter Statistics
	Counters []interface{}

	//device state
	ReloadLoopsPending int
	DeviceActive       bool
	DeviceConnected    bool
	//extra measurement statistics
	NumMeasurements int
	SysDescription  string
	NumMetrics      int
}

// Init initializes the device stat object
func (s *DevStat) Init(id string, tm map[string]string, l *logrus.Logger) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.id = id
	s.TagMap = tm
	s.log = l
	s.Counters = make([]interface{}, DevStatTypeSize)
	s.Counters[SnmpGetQueries] = 0
	s.Counters[SnmpWalkQueries] = 0
	s.Counters[SnmpGetErrors] = 0
	s.Counters[SnmpWalkErrors] = 0
	s.Counters[SnmpQueryTimeouts] = 0
	s.Counters[SnmpOIDGetAll] = 0
	s.Counters[SnmpOIDGetProcessed] = 0
	s.Counters[SnmpOIDGetErrors] = 0
	s.Counters[EvalMetricsAll] = 0
	s.Counters[EvalMetricsOk] = 0
	s.Counters[EvalMetricsErrors] = 0
	s.Counters[MetricSent] = 0
	s.Counters[MeasurementSent] = 0
	s.Counters[MetricSentErrors] = 0
	s.Counters[MeasurementSentErrors] = 0
	s.Counters[CicleGatherStartTime] = 0
	s.Counters[CicleGatherDuration] = 0.0
	s.Counters[FilterStartTime] = 0
	s.Counters[FilterDuration] = 0.0
	s.Counters[BackEndSentStartTime] = 0
	s.Counters[BackEndSentDuration] = 0.0
}

func (s *DevStat) reset() {
	for k, val := range s.Counters {
		switch v := val.(type) {
		case string:
			s.Counters[k] = ""
		case int32, int64, int:
			s.Counters[k] = 0
		case float64, float32:
			s.Counters[k] = 0.0
		default:
			s.log.Warnf("unknown typpe for counter %#v", v)
		}
	}
}

// GetCounter get Counter for stats
func (s *DevStat) GetCounter(stat DevStatType) interface{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.Counters[stat]
}

func (s *DevStat) getMetricFields() map[string]interface{} {
	fields := map[string]interface{}{
		/*0*/ //"snmp_get_queries": s.Counters[SnmpGetQueries],
		/*1*/ //"snmp_walk_queries": s.Counters[SnmpWalkQueries],
		/*2*/ //"snmp_get_errors": s.Counters[SnmpGetErrors],
		/*3*/ //"snmp_walk_errors": s.Counters[SnmpWalkErrors],
		/*4*/ //"snmp_query_timeouts": s.Counters[SnmpQueryTimeouts],
		/*5*/ "snmp_oid_get_all": s.Counters[SnmpOIDGetAll],
		/*6*/ "snmp_oid_get_processed": s.Counters[SnmpOIDGetProcessed],
		/*7*/ "snmp_oid_get_errors": s.Counters[SnmpOIDGetErrors],
		/*8*/ //"eval_metric_all": s.Counters[EvalMetricsAll],
		/*9*/ //"eval_metric_ok": s.Counters[EvalMetricsOk],
		/*10*/ //"eval_metric_errors": s.Counters[EvalMetricsErrors],
		/*11*/ "metric_sent": s.Counters[MetricSent],
		/*12*/ "metric_sent_errors": s.Counters[MetricSentErrors],
		/*13*/ "measurement_sent": s.Counters[MeasurementSent],
		/*14*/ "measurement_sent_errors": s.Counters[MeasurementSentErrors],
		/*15*/ "cicle_gather_start_time": s.Counters[CicleGatherStartTime],
		/*16*/ "cicle_gather_duration": s.Counters[CicleGatherDuration],
		/*17*/ "filter_start_time": s.Counters[FilterStartTime],
		/*18*/ "filter_duration": s.Counters[FilterDuration],
		/*19*/ "backend_sent_start_time": s.Counters[BackEndSentStartTime],
		/*20*/ "backend_sent_duration": s.Counters[BackEndSentDuration],
	}
	return fields
}

// SetSelfMonitoring set the output device where send monitoring metrics
func (s *DevStat) SetSelfMonitoring(cfg *selfmon.SelfMon) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.selfmon = cfg
}

// ThSafeCopy get a new object with public data copied in thread safe way
func (s *DevStat) ThSafeCopy() *DevStat {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	st := &DevStat{}
	st.Init(s.id, s.TagMap, s.log)
	for k, v := range s.Counters {
		st.Counters[k] = v
	}
	return st
}

// Send send data to the selfmon device
func (s *DevStat) Send() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.log.Infof("STATS SNMP GET: snmp pooling took [%f seconds] SNMP: Gets [%d] , Processed [%d], Errors [%d]", s.Counters[CicleGatherDuration], s.Counters[SnmpOIDGetAll], s.Counters[SnmpOIDGetProcessed], s.Counters[SnmpOIDGetErrors])
	s.log.Infof("STATS SNMP FILTER: filter pooling took [%f seconds] ", s.Counters[FilterDuration])
	s.log.Infof("STATS INFLUX: influx send took [%f seconds]", s.Counters[BackEndSentDuration])
	if s.selfmon != nil {
		s.selfmon.AddDeviceMetrics(s.id, s.getMetricFields(), s.TagMap)
	}
}

// ResetCounters initialize metric counters
func (s *DevStat) ResetCounters() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.reset()
}

// CounterInc n values to the counter set by id
func (s *DevStat) CounterInc(id DevStatType, n int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Counters[id] = s.Counters[id].(int) + int(n)
}

// AddMeasStats add measurement stats to the device stats object
func (s *DevStat) AddMeasStats(mets int64, mete int64, meass int64, mease int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Counters[MetricSent] = s.Counters[MetricSent].(int) + int(mets)
	s.Counters[MetricSentErrors] = s.Counters[MetricSentErrors].(int) + int(mete)
	s.Counters[MeasurementSent] = s.Counters[MeasurementSent].(int) + int(meass)
	s.Counters[MeasurementSentErrors] = s.Counters[MeasurementSentErrors].(int) + int(mease)
}

// UpdateSnmpGetStats update snmp statistics
func (s *DevStat) UpdateSnmpGetStats(g int64, p int64, e int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Counters[SnmpOIDGetAll] = s.Counters[SnmpOIDGetAll].(int) + int(g)
	s.Counters[SnmpOIDGetProcessed] = s.Counters[SnmpOIDGetProcessed].(int) + int(p)
	s.Counters[SnmpOIDGetErrors] = s.Counters[SnmpOIDGetErrors].(int) + int(e)
}

// SetGatherDuration Update Gather Duration stats
func (s *DevStat) SetGatherDuration(start time.Time, duration time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Counters[CicleGatherStartTime] = start.Unix()
	s.Counters[CicleGatherDuration] = duration.Seconds()
}

// AddSentDuration Update Sent Duration stats
func (s *DevStat) AddSentDuration(start time.Time, duration time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	//only register the first start time on concurrent mode
	if s.Counters[BackEndSentStartTime] == 0 {
		s.Counters[BackEndSentStartTime] = start.Unix()
	}
	s.Counters[BackEndSentDuration] = s.Counters[BackEndSentDuration].(float64) + duration.Seconds()
}

// SetFltUpdateStats Set Filter Stats
func (s *DevStat) SetFltUpdateStats(start time.Time, duration time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Counters[FilterStartTime] = start.Unix()
	s.Counters[FilterDuration] = duration.Seconds()
}
