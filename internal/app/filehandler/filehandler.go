package filehandler

type FileHandler interface {
	Write(dst string, data []byte) error
}
