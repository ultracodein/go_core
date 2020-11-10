package spider

import (
	"gosearch/pkg/crawler"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
)

func indexPage(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`<html><head><title>Index</title></head><body><a href="/news.php">News</a><a href="/about.php">About</a></body></html>`))
}

func newsPage(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`<html><head><title>News</title></head><body><p>News</p></body></html>`))
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`<html><head><title>About</title></head><body><p>About</p><a href="/contacts.php">Contacts</a><a href="/contacts.php">Contacts</a></body></html>`))
}

func contactsPage(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`<html><head><title>Contacts</title></head><body><p>Contacts</p><a href="/contacts.php">Contacts</a></body></html>`))
}

func createFakeServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", indexPage)
	handler.HandleFunc("/news.php", newsPage)
	handler.HandleFunc("/about.php", aboutPage)
	handler.HandleFunc("/contacts.php", contactsPage)
	srv := httptest.NewServer(handler)
	return srv
}

func TestService_Scan(t *testing.T) {
	srv := createFakeServer()
	defer srv.Close()

	s := New()
	url := srv.URL
	got, _ := s.Scan(url, 3)
	want := []crawler.Document{
		{
			ID:    0,
			URL:   url,
			Title: "Index",
			Body:  "",
		},
		{
			ID:    0,
			URL:   url + "/news.php",
			Title: "News",
			Body:  "",
		},
		{
			ID:    0,
			URL:   url + "/about.php",
			Title: "About",
			Body:  "",
		},
		{
			ID:    0,
			URL:   url + "/contacts.php",
			Title: "Contacts",
			Body:  "",
		},
	}

	// из-за рекурсии в Scan порядок элементов в slice
	// может меняться от запуска к запуску:
	// сортируем slice-ы, чтобы порядок не влиял успешное выполнение теста

	sort.SliceStable(got, func(i, j int) bool {
		return len(got[i].URL) < len(got[j].URL)
	})
	sort.SliceStable(want, func(i, j int) bool {
		return len(want[i].URL) < len(want[j].URL)
	})

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
