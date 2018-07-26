package dbhandler

// PagedResults paged results from db
type PagedResults struct {
	Total           int                      `json:"total"`
	CurrentPage     int                      `json:"currentPage"`
	TotalPage       int                      `json:"totalPage"`
	PageSize        int                      `json:"pageSize"`
	NextPage        int                      `json:"nextPage,omitempty"`
	PreviousPage    int                      `json:"previousPage,omitempty"`
	HasNextPage     bool                     `json:"hasNextPage,omitempty"`
	HasPreviousPage bool                     `json:"hasPreviousPage,omitempty"`
	Items           []map[string]interface{} `json:"items"`
}

// DatabaseConfig provide a uniform struct for storing database configs
type DatabaseConfig struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	User           string `json:"user"`
	Pass           string `json:"pass"`
	Database       string `json:"database"`
	AuthDB         string `json:"auth_db,omitempty"`
	CollectionName string `json:"collection_name"`
}

// DatabaseHandler defines interface for a database handler
type DatabaseHandler interface {
	GetConnection() error
	CloseConnection()
	GetAllItems(dataname string, limit int, page int, orderBy string, sortBy string, filters map[string]interface{}) (PagedResults, error)
	AddNewItem(dataName string, item map[string]interface{}) (map[string]interface{}, error)
	RemoveItemByID(dataName string, id interface{}) error
	FindItemByID(dataName string, id interface{}) (map[string]interface{}, error)
	UpdateBy(dataName string, selector interface{}, update map[string]interface{}) error
	IsConnecting() bool
}
