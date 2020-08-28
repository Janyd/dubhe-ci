package utils

type TimeHandler interface {
	CreatedAt() int64
	UpdatedAt() int64
}
