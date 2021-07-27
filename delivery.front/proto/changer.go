package proto

// Тестовый вариант ввода данных для изменения статуса курьера
type ChangeStat struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
}
