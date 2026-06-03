package main

import (
	"cos518project/chubby/api"
	"fmt"
	"net/rpc"
)

// Khai báo bảng màu ANSI xịn sò cho Terminal
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

func main() {
	// Danh bạ 3 Node trong hệ thống
	nodes := []string{"127.0.0.1:5379", "127.0.0.1:6379", "127.0.0.1:7379"}

	// Vẽ Header bảng
	fmt.Println(Bold + Cyan + "\n╔════════════════════════════════════════════════════════════════╗" + Reset)
	fmt.Println(Bold + Cyan + "║         CHUBBY DISTRIBUTED LOCK SERVICE - CLUSTER STATUS       ║" + Reset)
	fmt.Println(Bold + Cyan + "╠════════════════╦══════════════╦══════════╦═════════════════════╣" + Reset)
	fmt.Println(Bold + Cyan + "║    NODE (IP)   ║     ROLE     ║   TERM   ║    LEADER ADDRESS   ║" + Reset)
	fmt.Println(Bold + Cyan + "╠════════════════╬══════════════╬══════════╬═════════════════════╣" + Reset)

	// Đi gõ cửa từng Node một
	for _, addr := range nodes {
		client, err := rpc.Dial("tcp", addr)

		// Nếu Node bị sập (Tắt Terminal), báo OFFLINE màu đỏ cực ngầu
		if err != nil {
			fmt.Printf("║ %-14s ║ "+Bold+Red+"%-12s"+Reset+Cyan+" ║ %-8s ║ %-19s ║\n", addr, "OFFLINE", "-", "-")
			continue
		}

		var req api.StatusRequest
		var resp api.StatusResponse

		err = client.Call("Handler.Status", req, &resp)
		if err != nil {
			fmt.Printf("║ %-14s ║ "+Bold+Red+"%-12s"+Reset+Cyan+" ║ %-8s ║ %-19s ║\n", addr, "ERROR", "-", "-")
			client.Close()
			continue
		}

		// Xử lý màu sắc cho Role (Leader thì tô Xanh lá, Follower tô Vàng)
		rolePad := fmt.Sprintf("%-12s", resp.NodeRole)
		if resp.NodeRole == "Leader" {
			rolePad = Bold + Green + rolePad + Reset + Cyan
		} else {
			rolePad = Yellow + rolePad + Reset + Cyan
		}

		// Xử lý chuỗi rỗng
		termStr := resp.Term
		if termStr == "" {
			termStr = "N/A"
		}
		leaderAddr := resp.LeaderAddress
		if leaderAddr == "" {
			leaderAddr = "N/A"
		}

		// In ra từng hàng dữ liệu
		fmt.Printf("║ %-14s ║ %s ║ %-8s ║ %-19s ║\n", addr, rolePad, termStr, leaderAddr)
		client.Close()
	}

	// Vẽ Footer bảng
	fmt.Println(Bold + Cyan + "╚════════════════╩══════════════╩══════════╩═════════════════════╝\n" + Reset)
}
