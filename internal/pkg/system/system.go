package system

import (
	"fmt"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"net"
	"net/http"
)

func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("не удалось найти локальный IP-адрес")
}

// Функция для проверки, входит ли IP-адрес в доверенную подсеть
func IsIPInTrustedRange(ip string, trustedRange string) (bool, error) {
	// Преобразуем строку подсети в объект типа *net.IPNet
	_, trustedNet, err := net.ParseCIDR(trustedRange)
	if err != nil {
		return false, fmt.Errorf("не удалось разобрать подсеть: %v", err)
	}

	// Преобразуем IP-адрес в объект типа net.IP
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false, fmt.Errorf("не удалось разобрать IP-адрес: %v", ip)
	}

	// Проверяем, что IP-адрес входит в подсеть
	return trustedNet.Contains(parsedIP), nil
}

func TrustedRangeMiddleware(config interfaces.ConfigServerInterface, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.GetTrustedSubnet() == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Извлекаем IP из заголовка X-Real-IP
		realIP := r.Header.Get("X-Real-IP")
		if realIP == "" {
			http.Error(w, "Отсутствует заголовок X-Real-IP", http.StatusForbidden)
			return
		}

		// Проверяем, входит ли IP-адрес в доверенную подсеть
		isTrusted, err := IsIPInTrustedRange(realIP, config.GetTrustedSubnet())
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка проверки подсети: %v", err), http.StatusForbidden)
			return
		}

		if !isTrusted {
			// Если IP не входит в доверенную подсеть, возвращаем статус 403
			http.Error(w, "403 Forbidden: IP-адрес не входит в доверенную подсеть", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
		return
	})
}
