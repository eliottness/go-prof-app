//go:build !cgo

package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ebitengine/purego"
	"github.com/felixge/go-prof-app/shared"
)

var cgoOpSym uintptr
var cgoCpuHopSym uintptr = loadAndGetSymbol()

// Using the embeded shared library file, we dump it into a file and load it using dlopen
func loadLib() (uintptr, error) {

	file, err := os.CreateTemp("", "libhog-*.so")
	if err != nil {
		return uintptr(0), fmt.Errorf("error creating temp file: %w", err)
	}

	if err := os.WriteFile(file.Name(), shared.LibHog, 0400); err != nil {
		return uintptr(0), fmt.Errorf("error writing file: %w", err)
	}

	lib, err := purego.Dlopen(file.Name(), purego.RTLD_GLOBAL|purego.RTLD_NOW)
	if err != nil {
		return uintptr(0), fmt.Errorf("error opening shared library at path '%s'. Reason: %s", file.Name(), err)
	}

	return lib, nil
}

func loadAndGetSymbol() uintptr {
	lib, err := loadLib()
	if err != nil {
		panic(err)
	}

	sym, err := purego.Dlsym(lib, "cpuHog")
	if err != nil {
		panic(err)
	}

	cgoOpSym, err = purego.Dlsym(lib, "cpuHogOp")
	if err != nil {
		panic(err)
	}

	log.Println("Loaded libhog")

	return sym
}

type PostsHandler struct {
	DB          *sql.DB
	CPUDuration time.Duration
	SQLDuration time.Duration
	// CGO determines if cgo is used for simulating the CPUDuration.
	CGO bool
}

func (h *PostsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth(h.DB, w, r)
	if !ok {
		return
	}

	posts, err := h.ioWork(r.Context(), userID)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, "ioWork: %s", err)
		return
	}

	data, err := h.cpuWork(posts)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, "cpuWork: %s", err)
		return
	}
	w.Write(data)
}

func (h *PostsHandler) ioWork(ctx context.Context, userID int) ([]*Post, error) {
	q := `SELECT id, user_id, title, body FROM posts, pg_sleep($1) WHERE user_id = $2`
	var rows *sql.Rows
	rows, err := h.DB.QueryContext(ctx, q, h.SQLDuration.Seconds(), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var p Post
		if err = rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Body); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	return posts, nil
}

func (h *PostsHandler) cpuWork(posts []*Post) ([]byte, error) {
	var (
		wg   sync.WaitGroup
		stop = make(chan struct{})
		data []byte
	)
	wg.Add(1)
	if h.CGO {
		go cgoCPUHog(posts, &data, &wg, stop)
	} else {
		go goCPUHog(posts, &data, &wg, stop)
	}
	time.Sleep(h.CPUDuration)
	close(stop)
	wg.Wait()
	return data, nil
}

//go:noinline
func cgoCPUHog(posts []*Post, data *[]byte, wg *sync.WaitGroup, stop chan struct{}) {
	defer wg.Done()

loop:
	for {
		select {
		case <-stop:
			break loop
		default:
			purego.SyscallN(cgoCpuHopSym, cgoOpSym)
		}
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(posts); err != nil {
		return
	}
	*data = buf.Bytes()
}

//go:noinline
func goCPUHog(posts []*Post, data *[]byte, wg *sync.WaitGroup, stop chan struct{}) {
	defer wg.Done()

	for {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(posts); err != nil {
			return
		}
		select {
		case <-stop:
			*data = buf.Bytes()
			return
		default:
			buf.Reset()
		}
	}
}

type Post struct {
	ID     int
	UserID int
	Title  string
	Body   string
}
