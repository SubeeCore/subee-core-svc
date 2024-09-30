package entities_categories_v1

type CategoriesRecap struct {
	Categories []*CategoryRecap `json:"categories"`
}

type CategoryRecap struct {
	Name       string `json:"name"`
	Percentage string `json:"percentage"`
}
