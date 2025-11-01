package paginator

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

type Page struct {
	Order   int32
	URL     string
	Current bool
}

type Paginator struct {
	Prev          *Page
	Next          *Page
	Total         int32
	PerPage       int32
	TotalShowing  int32
	CurrentPage   int32
	ShowingRange  string
	Pages         []Page
	CountPaginate int32
}

func countPaginate(a, b int32) int32 {
	if a > 0 {
		c := a / b
		if a%b != 0 {
			c = c + 1
		}
		return c
	}
	return 0
}

func NewPaginator(currentPage, perPage int32, total int32, req *http.Request) Paginator {
	p := Paginator{}
	p.PerPage = perPage
	p.Total = total
	p.CurrentPage = currentPage
	lastPage := int32(math.Max(math.Ceil(float64(total)/float64(perPage)), 1))
	if lastPage > 1 {
		numberOfItems := currentPage * perPage
		p.TotalShowing = numberOfItems
		if numberOfItems > total {
			p.TotalShowing = total
		}
		p.ShowingRange = fmt.Sprintf("%d-%d", numberOfItems-perPage+1, p.TotalShowing)
	} else {
		p.ShowingRange = strconv.Itoa(int(total))
		p.TotalShowing = total
	}

	for i := int32(1); i <= lastPage; i++ {
		if i == 1 || i == lastPage || (i >= currentPage-2 && i <= currentPage+2) {
			pageURL := *req.URL
			params := pageURL.Query()
			params.Set("page", strconv.Itoa(int(i)))
			pageURL.RawQuery = params.Encode()

			p.Pages = append(p.Pages, Page{
				Order:   i,
				URL:     pageURL.String(),
				Current: i == currentPage,
			})
		} else if i == currentPage-3 || i == currentPage+3 {
			p.Pages = append(p.Pages, Page{})
		}
	}

	if currentPage != 1 {
		pageURL := *req.URL
		params := pageURL.Query()
		params.Set("page", strconv.Itoa(int(currentPage)-1))
		pageURL.RawQuery = params.Encode()
		p.Prev = &Page{URL: pageURL.String()}
	}

	if currentPage != lastPage {
		pageURL := *req.URL
		params := pageURL.Query()
		params.Set("page", strconv.Itoa(int(currentPage)+1))
		pageURL.RawQuery = params.Encode()
		p.Next = &Page{URL: pageURL.String()}
	}

	p.CountPaginate = countPaginate(p.Total, p.PerPage)

	return p
}
