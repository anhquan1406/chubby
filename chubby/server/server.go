// Define Chubby server application.

// Adapted from Leto server file:
// https://github.com/yongman/leto/blob/master/server/server.go

// Copyright (C) 2018 YanMing <yming0221@gmail.com>
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package server

import (
	"cos518project/chubby/api"
	"cos518project/chubby/config"
	"cos518project/chubby/store"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

type App struct {
	listener net.Listener

	// wrapper and manager for db instance
	store *store.Store

	logger *log.Logger

	// Current Node's Address
	address string

	// In-memory struct of handles.
	// Maps handle IDs to handle metadata.
	// handles map[int]Handle

	// In-memory struct of locks.
	// Maps filepaths to Lock structs.
	locks map[api.FilePath]*Lock

	// In-memory struct of sessions.
	sessions map[api.ClientID]*Session
}

// No choice but to make this variable package-level :(
var app *App

func Run(conf *config.Config) {
	var err error

	// Init app struct.
	app = &App{
		logger:   log.New(os.Stderr, "[server] ", log.LstdFlags),
		store:    store.New(conf.RaftDir, conf.RaftBind, conf.InMem),
		address:  conf.Listen,
		locks:    make(map[api.FilePath]*Lock),
		sessions: make(map[api.ClientID]*Session),
	}

	// Open the store.
	bootstrap := conf.Join == ""
	err = app.store.Open(bootstrap, conf.NodeID)
	if err != nil {
		log.Fatal(err)
	}

	if !bootstrap {
		// Set up TCP connection.
		client, err := rpc.Dial("tcp", conf.Join)
		if err != nil {
			log.Fatal(err)
		}

		app.logger.Printf("set up connection to %s", conf.Join)

		var req JoinRequest
		var resp JoinResponse

		req.RaftAddr = conf.RaftBind
		req.NodeID = conf.NodeID

		err = client.Call("Handler.Join", req, &resp)
		if err != nil {
			log.Fatal(err)
		}
		if resp.Error != nil {
			log.Fatal(err)
		}
	}

	// Listen for client connections.
	handler := new(Handler)
	err = rpc.Register(handler)

	app.listener, err = net.Listen("tcp", conf.Listen)
	app.logger.Printf("server listen in %s", conf.Listen)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Accept connections.
	rpc.Accept(app.listener)
}

// ---------------------------------------------------------
// TÍNH NĂNG MỚI 1: API Status - Xử lý truy vấn trạng thái
// ---------------------------------------------------------

func (h *Handler) Status(req api.StatusRequest, resp *api.StatusResponse) error {
	// 1. Lấy trạng thái hiện tại (Leader, Follower, Candidate)
	resp.NodeRole = app.store.Raft.State().String()

	// 2. Lấy thông số nhiệm kỳ (term) từ bộ thống kê
	stats := app.store.Raft.Stats()
	resp.Term = stats["term"]

	// 3. Lấy địa chỉ của Leader hiện tại trong cụm
	resp.LeaderAddress = string(app.store.Raft.Leader())

	// 4. BƠM DỮ LIỆU: Lật sổ hộ khẩu lấy tên khách hàng đang Active
	var clients []string
	for clientID := range app.sessions {
		clients = append(clients, string(clientID))
	}
	resp.ActiveClients = clients

	return nil
}
