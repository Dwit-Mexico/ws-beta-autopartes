package domain

type WebCatalogRoutes struct {
	Name string `gorm:"type:nvarchar(100);not null" json:"name" validate:"required,min=3,max=100"`
	Path string `gorm:"type:nvarchar(300);not null" json:"path" validate:"required,min=3,max=300"`
}

type WebPages struct {
	Routes []WebCatalogRoutes `json:"routes"`
}
