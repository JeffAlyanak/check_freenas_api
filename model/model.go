package model

// Alerts struct holds FreeNAS alert data
type Alerts struct {
	Meta struct {
		Limit      int         `json:"limit"`
		Next       interface{} `json:"next"`
		Offset     int         `json:"offset"`
		Previous   interface{} `json:"previous"`
		TotalCount int         `json:"total_count"`
	} `json:"meta"`
	Objects []struct {
		Dismissed bool   `json:"dismissed"`
		ID        string `json:"id"`
		Level     string `json:"level"`
		Message   string `json:"message"`
		Timestamp int    `json:"timestamp"`
	} `json:"objects"`
}

// Storage struct holds FreeNAS pool data
type Storage []struct {
	Avail    int64 `json:"avail"`
	Children []struct {
		Avail    int64 `json:"avail"`
		Children []struct {
			Avail   int64  `json:"avail"`
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Path    string `json:"path"`
			Status  string `json:"status"`
			Type    string `json:"type"`
			Used    int64  `json:"used"`
			UsedPct int    `json:"used_pct"`
		} `json:"children"`
		ID         int    `json:"id"`
		Mountpoint string `json:"mountpoint"`
		Name       string `json:"name"`
		Path       string `json:"path"`
		Status     string `json:"status"`
		Type       string `json:"type"`
		Used       int64  `json:"used"`
		UsedPct    int    `json:"used_pct"`
	} `json:"children"`
	ID            int    `json:"id"`
	IsDecrypted   bool   `json:"is_decrypted"`
	IsUpgraded    bool   `json:"is_upgraded"`
	Mountpoint    string `json:"mountpoint"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	Used          int64  `json:"used"`
	UsedPct       string `json:"used_pct"`
	VolEncrypt    int    `json:"vol_encrypt"`
	VolEncryptkey string `json:"vol_encryptkey"`
	VolGUID       string `json:"vol_guid"`
	VolName       string `json:"vol_name"`
}
