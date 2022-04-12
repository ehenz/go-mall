package utils

import "net"

// GetFreePort 动态获取一个可用的端口
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	lis, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer lis.Close()
	return lis.Addr().(*net.TCPAddr).Port, nil
}
