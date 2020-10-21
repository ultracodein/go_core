package mock

// FakeScanner является сканером-заглушкой, предназначенной для тестирования
type FakeScanner func(string, int)

// Scan реализует интерфейс Scanner для сканера-заглушки
func (f FakeScanner) Scan(url string, depth int) (data map[string]string, err error) {
	data = map[string]string{
		"https://habr.com/":     "Лучшие публикации за сутки / Хабр",
		"https://www.cnews.ru/": "Интернет-издание о высоких технологиях - CNews",
	}

	return data, nil
}
