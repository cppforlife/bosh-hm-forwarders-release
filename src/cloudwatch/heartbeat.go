package main

const (
	EventKindHeartbeat = "heartbeat"
	EventKindAlert     = "alert"
)

type Event struct {
	ID        string
	Kind      string
	Timestamp int64

	Job        string
	Index      string
	InstanceID string `json:"instance_id"`
	JobState   string `json:"job_state"`
	Deployment string
	AgentID    string `json:"agent_id"`

	Metrics []Metric
	Vitals  Vitals
}

type Metric struct {
	Tags      map[string]string
	Timestamp int64
	Value     string
	Name      string
}

type Vitals struct {
	Swap PercentKB
	Mem  PercentKB
	Load []string
	Disk DiskVitals
}

type DiskVitals struct {
	System     PercentInode
	Persistent PercentInode
	Ephemeral  PercentInode
}

type CPUVitals struct {
	Wait string
	User string
	Sys  string
}

type PercentKB struct {
	Percent string
	KB      string
}

type PercentInode struct {
	Percent      string
	InodePercent string `json:"inode_percent"`
}
