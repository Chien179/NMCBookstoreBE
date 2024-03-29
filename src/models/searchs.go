package models

type Text struct {
	MultiMatch struct {
		Query  string   `json:"query"`
		Fields []string `json:"fields"`
		Type   string   `json:"type"`
	} `json:"multi_match"`
}

type GenreID struct {
	Term struct {
		GenresID int64 `json:"genres_id"`
	} `json:"term"`
}

type Price struct {
	Price struct {
		Price struct {
			Gte float64 `json:"gte"`
			Lte float64 `json:"lte"`
		} `json:"price"`
	} `json:"range,omitempty"`
}

type Rating struct {
	Rating struct {
		Rating struct {
			Lte float64 `json:"lte"`
		} `json:"rating"`
	} `json:"range"`
}

type PriceSort struct {
	Price struct {
		Order string `json:"order"`
	} `json:"price"`
}

type TextSort struct {
	NameKeyword struct {
		Order string `json:"order"`
	} `json:"name.keyword"`
}

type QueryElastic struct {
	Query struct {
		Bool struct {
			Must []interface {
			} `json:"must"`
		} `json:"bool"`
	} `json:"query"`
	Sort [1]struct {
		TextSort
		PriceSort
		Rating struct {
			Order string `json:"order"`
		} `json:"rating"`
	} `json:"sort"`
	Collapse struct {
		Field string `json:"field"`
	} `json:"collapse"`
	Aggs struct {
		UniqueBooks struct {
			Cardinality struct {
				Field string `json:"field"`
			} `json:"cardinality"`
		} `json:"unique_books"`
	} `json:"aggs"`
	From int32 `json:"from"`
	Size int32 `json:"size"`
}

type SearchResponse struct {
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore int `json:"max_score"`
		Hits     []struct {
			Index   string   `json:"_index"`
			Type    string   `json:"_type"`
			ID      string   `json:"_id"`
			Score   int64    `json:"_score"`
			Ignored []string `json:"_ignored"`
			Source  struct {
				GenresID    int64    `json:"genres_id"`
				Quantity    int32    `json:"quantity"`
				Publisher   string   `json:"publisher"`
				Image       []string `json:"image"`
				Description string   `json:"description"`
				Sale        float64  `json:"sale"`
				IsDeleted   bool     `json:"is_deleted"`
				Author      string   `json:"author"`
				Price       float64  `json:"price"`
				Rating      float64  `json:"rating"`
				Name        string   `json:"name"`
				ID          int64    `json:"id"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type Aggs struct {
	Aggregations struct {
		UniqueBooks struct {
			Value float64 `json:"value"`
		} `json:"unique_books"`
	} `json:"aggregations"`
}

type SearchRequest struct {
	PageID       int32   `form:"page_id" binding:"required,min=1"`
	PageSize     int32   `form:"page_size" binding:"required,min=24,max=100"`
	Text         string  `form:"text"`
	GenresID     int64   `form:"genres_id"`
	SubgenresID  int64   `form:"subgenres_id"`
	MinPrice     float64 `form:"min_price"`
	MaxPrice     float64 `form:"max_price"`
	Rating       float64 `form:"rating"`
	IsDeleted    bool    `form:"is_deleted"`
	PriceSortAsc bool    `form:"price_sort_asc,default=true"`
	NameSortAsc  bool    `form:"name_sort_asc,default=true"`
}

type RecommedRequest struct {
	ID   int64 `form:"book_id" binding:"required"`
	Size int64 `form:"size,default=6" binding:"required"`
}

type JustForYouRequest struct {
	UserName string `form:"username"  binding:"required"`
}
