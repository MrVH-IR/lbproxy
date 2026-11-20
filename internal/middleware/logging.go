package middleware

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	requests, errors *log.Logger
	mu               sync.Mutex
	dir              string
	retention        int
	day              string
	reqFile, errFile *os.File
}

func NewLogger(dir string, retention int) (*Logger, error) {
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, err
	}
	l := &Logger{dir: dir, retention: retention}
	if err := l.rotate(time.Now()); err != nil {
		return nil, err
	}
	return l, nil
}
func (l *Logger) rotate(now time.Time) error {
	day := now.Format("2006-01-02")
	if day == l.day {
		return nil
	}
	if l.reqFile != nil {
		_ = l.reqFile.Close()
	}
	if l.errFile != nil {
		_ = l.errFile.Close()
	}
	var err error
	l.reqFile, err = os.OpenFile(filepath.Join(l.dir, "requests-"+day+".log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	l.errFile, err = os.OpenFile(filepath.Join(l.dir, "errors-"+day+".log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	l.requests = log.New(l.reqFile, "", 0)
	l.errors = log.New(l.errFile, "", 0)
	l.day = day
	return l.archiveAndPrune(now)
}
func (l *Logger) archiveAndPrune(now time.Time) error {
	// Compress only closed daily logs; the current day's files remain writable.
	entries, err := os.ReadDir(l.dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		n := e.Name()
		if strings.HasSuffix(n, ".log") && len(n) >= 14 {
			date := n[len(n)-14 : len(n)-4]
			if _, e := time.Parse("2006-01-02", date); e == nil && date < l.day {
				src := filepath.Join(l.dir, n)
				dst := src + ".tar.gz"
				if _, e := os.Stat(dst); os.IsNotExist(e) {
					if e := archive(src, dst); e != nil {
						return e
					}
					if e := os.Remove(src); e != nil {
						return e
					}
				}
			}
		}
	}
	entries, err = os.ReadDir(l.dir)
	if err != nil {
		return err
	}
	var archives []os.DirEntry
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".tar.gz") {
			archives = append(archives, e)
		}
	}
	sort.Slice(archives, func(i, j int) bool { return archives[i].Name() > archives[j].Name() })
	for i := l.retention * 2; i < len(archives); i++ {
		if err := os.Remove(filepath.Join(l.dir, archives[i].Name())); err != nil {
			return err
		}
	}
	return nil
}
func archive(src, dst string) error {
	in, e := os.Open(src)
	if e != nil {
		return e
	}
	defer in.Close()
	out, e := os.Create(dst)
	if e != nil {
		return e
	}
	gz := gzip.NewWriter(out)
	tw := tar.NewWriter(gz)
	info, e := in.Stat()
	if e == nil {
		e = tw.WriteHeader(&tar.Header{Name: filepath.Base(src), Mode: 0640, Size: info.Size(), ModTime: info.ModTime()})
	}
	if e == nil {
		_, e = io.Copy(tw, in)
	}
	_ = tw.Close()
	_ = gz.Close()
	ce := out.Close()
	if e != nil {
		_ = os.Remove(dst)
		return e
	}
	return ce
}
func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.reqFile != nil {
		_ = l.reqFile.Close()
	}
	if l.errFile != nil {
		_ = l.errFile.Close()
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) { w.status = code; w.ResponseWriter.WriteHeader(code) }
func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(sw, r)
		line := fmt.Sprintf("%s ip=%q method=%s path=%q status=%d duration=%s", start.Format(time.RFC3339), r.RemoteAddr, r.Method, r.URL.RequestURI(), sw.status, time.Since(start))
		l.mu.Lock()
		defer l.mu.Unlock()
		if err := l.rotate(time.Now()); err != nil {
			log.Printf("log rotation: %v", err)
			return
		}
		l.requests.Println(line)
		if sw.status >= 400 {
			l.errors.Println(line)
		}
	})
}
func (l *Logger) Errorf(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if err := l.rotate(time.Now()); err != nil {
		log.Printf("log rotation: %v", err)
		return
	}
	l.errors.Printf("%s %s", time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))
}
