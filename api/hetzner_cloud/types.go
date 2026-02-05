package hetzner_cloud

type LocationType string

const (
	LocationNuernberg   LocationType = "nbg1"
	LocationFalkenstein LocationType = "fsn1"
	LocationHelsinki    LocationType = "hel1"
	LocationAhsburn     LocationType = "ash"
	LocationHillsboro   LocationType = "hil"
)

type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error *ErrorMessage `json:"error"`
}

type PublicNet struct {
	EnableIPv4 bool `json:"enable_ipv4"`
	EnableIPv6 bool `json:"enable_ipv6"`
}

type CreateServerRequest struct {
	Image            string            `json:"image"`
	Labels           map[string]string `json:"labels"`
	Location         string            `json:"location"`
	Name             string            `json:"name"`
	PublicNet        *PublicNet        `json:"public_net"`
	ServerType       string            `json:"server_type"`
	SSHKeys          []string          `json:"ssh_keys"`
	StartAfterCreate bool              `json:"start_after_create"`
}

type PublicNetEntry struct {
	IP string `json:"ip"`
}

type PublicNetInfo struct {
	IPv4 *PublicNetEntry `json:"ipv4"`
	IPv6 *PublicNetEntry `json:"ipv6"`
}

type ServerInfo struct {
	ID        int            `json:"id"`
	Status    string         `json:"status"`
	PublicNet *PublicNetInfo `json:"public_net"`
}

type CreateServerResponse struct {
	Server *ServerInfo `json:"server"`
}

type ServerTypesResponse struct {
	ServerTypes []*ServerType `json:"server_types"`
}

type ServerType struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Cores      int    `json:"cores"`
	Memory     int    `json:"memory"`
	Disk       int    `json:"disk"`
	Deprecated bool   `json:"deprecated"`
}

type TimeSeriesData struct {
	Values map[string][2]float64 `json:"values"`
}

type TimeSeries struct {
	CPU                 *TimeSeriesData `json:"cpu"`
	Disk0BandwithRead   *TimeSeriesData `json:"disk.0.bandwidth.read"`
	Disk0BandwithWrite  *TimeSeriesData `json:"disk.0.bandwidth.write"`
	Network0BandwithIn  *TimeSeriesData `json:"network.0.bandwidth.in"`
	Network0BandwithOut *TimeSeriesData `json:"network.0.bandwidth.out"`
}

type Metrics struct {
	Start      string      `json:"start"`
	End        string      `json:"end"`
	Step       int         `json:"step"`
	TimeSeries *TimeSeries `json:"time_series"`
}

type MetricsResponse struct {
	Metrics *Metrics `json:"metrics"`
}
