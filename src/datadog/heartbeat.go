package main

const (
	EventKindHeartbeat = "heartbeat"
	EventKindAlert     = "alert"
)

type Event struct {
	// Shared
	ID   string
	Kind string

	Deployment string

	// Heartbeat
	Timestamp  int64
	Job        string
	Index      string
	InstanceID string `json:"instance_id"`
	JobState   string `json:"job_state"`
	AgentID    string `json:"agent_id"`
	Teams      []string
	Metrics    []Metric
	Vitals     Vitals

	// Alert
	Severity  int
	Category  string
	Title     string
	Summary   string
	Source    string
	CreatedAt int `json:"created_at"`
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

func SeverityToString(severity int) string { // todo enum
	switch severity {
	case 1:
		return "alert"
	case 2:
		return "critical"
	case 3:
		return "error"
	case 4:
		return "warning"
		// case -1: return "ignored"
	default:
		return "ignored"
	}
}
