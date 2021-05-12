package tradingdb2task

// TaskGroup - TaskGroup
type TaskGroup struct {
	StartTs      int64 `json:"StartTs"`
	TaskGroupID  int   `json:"TaskGroupID"`
	RunningTime  int64 `json:"RunningTime"`
	LastTaskNums int   `json:"LastTaskNums"`
	MaxTaskNums  int   `json:"MaxTaskNums"`
	LastTime     int64 `json:"LastTime"`
}
