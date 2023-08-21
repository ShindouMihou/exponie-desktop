package utils

func Has[T comparable](slice []T, item T) bool {
	for _, val := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// cc:https://stackoverflow.com/a/72408490
func Chunk[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}
