package jpush

import "encoding/json"

//Schedule 定时任务
type Schedule struct {
	ScheduleID string      `json:"schedule_id,omitempty"`
	Cid        string      `json:"cid,omitempty"`
	Name       string      `json:"name"`
	Enabled    bool        `json:"enabled"`
	Trigger    interface{} `json:"trigger"`
	Push       PayLoad     `json:"push"`
}

//Single 定时任务
type Single struct {
	Time string `json:"time"`
}

//Periodical 定期任务
type Periodical struct {
	Start     string   `json:"start"`
	End       string   `json:"end"`
	Time      string   `json:"time"`
	TimeUnit  string   `json:"time_unit"`
	Frequency int      `json:"frequency"`
	Point     []string `json:"point"`
}

func NewSchedule(name string) *Schedule {
	schedule := &Schedule{}
	schedule.Name = name
	schedule.Enabled = true
	return schedule
}
func (schedule *Schedule) SetPayload(payload PayLoad) {
	schedule.Push = payload
}
func (schedule *Schedule) SetCid(cid string) {
	schedule.Cid = cid
}

func (schedule *Schedule) SetSingle(single Single) {
	schedule.Trigger = map[string]interface{}{"single": single}
}
func (schedule *Schedule) SetPeriodical(periodical Periodical) {
	schedule.Trigger = map[string]interface{}{"periodical": periodical}
}

func (schedule *Schedule) ToBytes() ([]byte, error) {
	content, err := json.Marshal(schedule)
	if err != nil {
		return nil, err
	}
	return content, nil
}
