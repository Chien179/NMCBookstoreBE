package helper

import (
	"encoding/json"
	"fmt"

	"github.com/Chien179/NMCBookstoreBE/src/models"
)

func QueryElastic(search models.SearchRequest) (string, error) {
	var query models.QueryElastic
	arrange := "asc"

	if search.Text != "" {
		var text models.Text
		text.MultiMatch.Query = search.Text
		text.MultiMatch.Fields = []string{"name", "author"}
		text.MultiMatch.Type = "phrase_prefix"

		query.Query.Bool.Must = append(query.Query.Bool.Must, text)
	}

	if !search.NameSortAsc {
		arrange = "desc"
	}

	var textSort models.TextSort
	textSort.NameKeyword.Order = arrange

	query.Sort[0].TextSort = textSort

	arrange = "asc"

	if search.GenresID > 0 {
		var genreID models.GenreID
		genreID.Term.GenresID = search.GenresID

		query.Query.Bool.Must = append(query.Query.Bool.Must, genreID)
	}

	var price models.Price
	price.Price.Price.Gte = search.MinPrice

	if search.MaxPrice != 0 {
		price.Price.Price.Lte = search.MaxPrice
	} else {
		price.Price.Price.Lte = 500
	}

	query.Query.Bool.Must = append(query.Query.Bool.Must, price)

	if !search.PriceSortAsc {
		arrange = "desc"
	}

	var priceSort models.PriceSort
	priceSort.Price.Order = arrange

	query.Sort[0].PriceSort = priceSort

	arrange = "asc"

	if search.Rating > 0 {
		var rating models.Rating
		rating.Rating.Rating.Lte = search.Rating

		query.Query.Bool.Must = append(query.Query.Bool.Must, rating)
	}

	query.Sort[0].Rating.Order = arrange

	query.From = search.PageID - 1
	query.Size = search.PageSize

	query.Aggs.UniqueBooks.Cardinality.Field = "id"

	query.Collapse.Field = "id"

	queryByte, err := json.Marshal(query)
	if err != nil {
		return "", err
	}
	result := json.RawMessage(queryByte)
	fmt.Println(string(result))

	return string(result), nil
}
