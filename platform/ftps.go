package platform

import (
	"os"
)

// FTPSClient represents SakuraCloud FTPS Client
type FTPSClient interface {
	Connect(string, int) error
	Login(string, string) error
	StoreFile(string, *os.File) error
	Quit() error
}
