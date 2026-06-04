package main

import (
	"cos518project/chubby/api"
	"fmt"
	"net/rpc"
	"time"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

func main() {
	nodes := []string{"127.0.0.1:5379", "127.0.0.1:6379", "127.0.0.1:7379"}

	// Vòng lặp vô hạn để cập nhật liên tục
	for {
		// Lệnh dọn dẹp màn hình (Clear Screen)
		fmt.Print("\033[H\033[2J")

		timeNow := time.Now().Format("15:04:05")

		fmt.Println(Bold + Cyan + "╔═══════════════════════════════════════════════════════════════════════════════════════════╗" + Reset)
		fmt.Printf(Bold+Cyan+"║                  CHUBBY LIVE DASHBOARD - LIVE 🟢 (%s)                               ║\n"+Reset, timeNow)
		fmt.Println(Bold + Cyan + "╠════════════════╦══════════════╦══════════╦═════════════════════╦══════════════════════════╣" + Reset)
		fmt.Println(Bold + Cyan + "║    NODE (IP)   ║     ROLE     ║   TERM   ║    LEADER ADDRESS   ║     ACTIVE SESSIONS      ║" + Reset)
		fmt.Println(Bold + Cyan + "╠════════════════╬══════════════╬══════════╬═════════════════════╬══════════════════════════╣" + Reset)

		for _, addr := range nodes {
			client, err := rpc.Dial("tcp", addr)
			if err != nil {
				fmt.Printf("║ %-14s ║ "+Bold+Red+"%-12s"+Reset+Cyan+" ║ %-8s ║ %-19s ║ %-24s ║\n", addr, "OFFLINE", "-", "-", "-")
				continue
			}

			var req api.StatusRequest
			var resp api.StatusResponse

			err = client.Call("Handler.Status", req, &resp)
			if err != nil {
				fmt.Printf("║ %-14s ║ "+Bold+Red+"%-12s"+Reset+Cyan+" ║ %-8s ║ %-19s ║ %-24s ║\n", addr, "ERROR", "-", "-", "-")
				client.Close()
				continue
			}

			// Xử lý màu sắc
			rolePad := fmt.Sprintf("%-12s", resp.NodeRole)
			if resp.NodeRole == "Leader" {
				rolePad = Bold + Green + rolePad + Reset + Cyan
			} else {
				rolePad = Yellow + rolePad + Reset + Cyan
			}

			termStr := resp.Term
			if termStr == "" {
				termStr = "N/A"
			}
			leaderAddr := resp.LeaderAddress
			if leaderAddr == "" {
				leaderAddr = "N/A"
			}

			// XỬ LÝ TỰ ĐỘNG XUỐNG DÒNG CHO KHÁCH HÀNG
			clientsList := resp.ActiveClients
			if len(clientsList) == 0 {
				clientsList = []string{"Trống"} // Nếu không có ai thì gán 1 phần tử là "Trống"
			}

			for i, clientName := range clientsList {
				// Cắt chuỗi nếu tên quá dài để không làm vỡ khung
				if len(clientName) > 24 {
					clientName = clientName[:21] + "..."
				}

				if i == 0 {
					// Vòng lặp đầu tiên: In đầy đủ thông tin Node + Khách hàng đầu tiên
					fmt.Printf("║ %-14s ║ %s ║ %-8s ║ %-19s ║ %-24s ║\n", addr, rolePad, termStr, leaderAddr, clientName)
				} else {
					// Các khách hàng tiếp theo: Bỏ trống các cột IP, Role, Term, Leader
					fmt.Printf("║ %-14s ║ %-12s ║ %-8s ║ %-19s ║ %-24s ║\n", "", "", "", "", clientName)
				}
			}

			client.Close()
		}

		fmt.Println(Bold + Cyan + "╚════════════════╩══════════════╩══════════╩═════════════════════╩══════════════════════════╝" + Reset)
		fmt.Println("\n⏳ Đang làm mới dữ liệu sau mỗi 10 giây... (Bấm Ctrl+C để thoát)")

		// Cho tiến trình ngủ 10 giây trước khi quét lại
		time.Sleep(10 * time.Second)
	}
}
