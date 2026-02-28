package types

type ImageName string

const (
	ImageNameUI          ImageName = "mtui"
	ImageNameNginx       ImageName = "nginx"
	ImageNameNodExporter ImageName = "node_exporter"
	ImageNameTraefik     ImageName = "traefik"
	ImageNameIPv6Nat     ImageName = "ipv6nat"
)

type ImageVersion struct {
	Name    ImageName `json:"name" gorm:"primarykey;column:name"`
	Version string    `json:"version" gorm:"column:version"`
}

func (m *ImageVersion) TableName() string {
	return "image_version"
}
