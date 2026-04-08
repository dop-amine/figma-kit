package restapi

// Component represents a published component in a Figma team/file library.
type Component struct {
	Key             string    `json:"key"`
	FileKey         string    `json:"file_key"`
	NodeID          string    `json:"node_id"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
	User            User      `json:"user"`
	ContainingFrame FrameInfo `json:"containing_frame"`
}

// ComponentSet represents a published component set (variants) in a library.
type ComponentSet struct {
	Key             string    `json:"key"`
	FileKey         string    `json:"file_key"`
	NodeID          string    `json:"node_id"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
	User            User      `json:"user"`
	ContainingFrame FrameInfo `json:"containing_frame"`
}

// Style represents a published style in a library.
type Style struct {
	Key             string    `json:"key"`
	FileKey         string    `json:"file_key"`
	NodeID          string    `json:"node_id"`
	StyleType       string    `json:"style_type"` // FILL, TEXT, EFFECT, GRID
	ThumbnailURL    string    `json:"thumbnail_url"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
	User            User      `json:"user"`
	ContainingFrame FrameInfo `json:"containing_frame"`
}

// User represents a Figma user associated with a library asset.
type User struct {
	ID     string `json:"id"`
	Handle string `json:"handle"`
	ImgURL string `json:"img_url"`
}

// FrameInfo describes the frame that contains a component or style.
type FrameInfo struct {
	NodeID                 string `json:"nodeId"`
	Name                   string `json:"name"`
	BackgroundColor        string `json:"backgroundColor"`
	PageID                 string `json:"pageId"`
	PageName               string `json:"pageName"`
	ContainingComponentSet string `json:"containingComponentSet"`
}

// ComponentsResponse wraps a paginated list of components.
type ComponentsResponse struct {
	Status int            `json:"status"`
	Error  bool           `json:"error"`
	Meta   ComponentsMeta `json:"meta"`
}

// ComponentsMeta holds the components array and pagination cursor.
type ComponentsMeta struct {
	Components []Component `json:"components"`
	Cursor     Cursor      `json:"cursor"`
}

// ComponentSetsResponse wraps a paginated list of component sets.
type ComponentSetsResponse struct {
	Status int               `json:"status"`
	Error  bool              `json:"error"`
	Meta   ComponentSetsMeta `json:"meta"`
}

// ComponentSetsMeta holds the component_sets array and pagination cursor.
type ComponentSetsMeta struct {
	ComponentSets []ComponentSet `json:"component_sets"`
	Cursor        Cursor         `json:"cursor"`
}

// StylesResponse wraps a paginated list of styles.
type StylesResponse struct {
	Status int        `json:"status"`
	Error  bool       `json:"error"`
	Meta   StylesMeta `json:"meta"`
}

// StylesMeta holds the styles array and pagination cursor.
type StylesMeta struct {
	Styles []Style `json:"styles"`
	Cursor Cursor  `json:"cursor"`
}

// SingleComponentResponse wraps a single component lookup.
type SingleComponentResponse struct {
	Status int       `json:"status"`
	Error  bool      `json:"error"`
	Meta   Component `json:"meta"`
}

// SingleComponentSetResponse wraps a single component set lookup.
type SingleComponentSetResponse struct {
	Status int          `json:"status"`
	Error  bool         `json:"error"`
	Meta   ComponentSet `json:"meta"`
}

// SingleStyleResponse wraps a single style lookup.
type SingleStyleResponse struct {
	Status int   `json:"status"`
	Error  bool  `json:"error"`
	Meta   Style `json:"meta"`
}

// Cursor holds pagination state for paginated responses.
type Cursor struct {
	Before string `json:"before"`
	After  string `json:"after"`
}
