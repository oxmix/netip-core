package collector

import (
	"bufio"
	"bytes"
	"context"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

var reWhoQuotes = regexp.MustCompile(`^\s*string\s+"([^"]+)"`)

func (c *Collector) collectWho() {
	c.ChanWhoLogged = make(chan *WhoLogged, 255)
	c.whoSessionMu = sync.Mutex{}
	c.whoSessions = make(map[string]*WhoLogged)

	for {
		cmd := exec.Command("dbus-monitor", "--system", "type='signal',sender='org.freedesktop.login1'")
		stderr, _ := cmd.StderrPipe()
		go func() {
			buf := make([]byte, 1024)
			n, _ := stderr.Read(buf)
			if n > 0 {
				log.Printf("[who] dbus-monitor stderr: %s", string(buf[:n]))
			}
		}()
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println("[who] dbus-monitor stdout err:", err)
			return
		}
		if err := cmd.Start(); err != nil {
			log.Println("[who] failed to start dbus-monitor:", err)
			return
		}

		log.Println("[who] dbus-monitor started")

		scanner := bufio.NewScanner(stdout)

		session := "" // new or removed
		for scanner.Scan() {
			line := scanner.Text()
			l := strings.ToLower(line)
			if strings.Contains(l, "member=sessionnew") {
				session = "new"
				continue
			}
			if strings.Contains(l, "member=sessionremoved") {
				session = "removed"
				continue
			}

			// skip if not allowed session
			if session == "" {
				continue
			}

			// find first line with ID
			if strings.Contains(line, "string") {
				m := reWhoQuotes.FindStringSubmatch(line)
				if len(m) >= 2 {
					go c.handleWhoEvent(session, m[1])
					// flush after
					session = ""
				}
				continue
			}
		}

		// if scanner fell or dbus-monitor ended
		if err := scanner.Err(); err != nil {
			log.Println("[who] scanner err:", err, "trying to restarting dbus-monitor")
		}
		_ = cmd.Wait()

		log.Println("[who] end")

		time.Sleep(8 * time.Second)
	}
}

func (c *Collector) handleWhoEvent(session, sessionID string) {
	if session == "removed" {
		c.whoSessionMu.Lock()
		if info, ok := c.whoSessions[sessionID]; ok {
			w := *info
			w.Session = "removed"
			w.Time = time.Now().UTC()
			c.whoSessionMu.Unlock()
			c.ChanWhoLogged <- &w

			c.whoSessionMu.Lock()
			delete(c.whoSessions, sessionID)
			c.whoSessionMu.Unlock()
			return
		}
		c.whoSessionMu.Unlock()
	}

	w := fetchSessionInfo(session, sessionID)
	if session == "new" && w.User != "unknown" {
		c.whoSessionMu.Lock()
		c.whoSessions[sessionID] = &WhoLogged{
			Session: "new",
			Time:    w.Time,
			User:    w.User,
			Device:  w.Device,
			IP:      w.IP,
		}
		c.whoSessionMu.Unlock()
	}

	if w.Session == "removed" && w.User == "unknown" {
		return
	}

	c.ChanWhoLogged <- w
}

func fetchSessionInfo(session, sessionId string) *WhoLogged {
	w := &WhoLogged{
		Session: session,
		Time:    time.Now().UTC(),
		User:    "unknown",
		Device:  "",
		IP:      "",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "loginctl", "show-session", sessionId,
		"-p", "Name", "-p", "TTY", "-p", "RemoteHost", "-p", "RemoteAddress").CombinedOutput()
	if err != nil {
		if session == "new" {
			log.Printf("[who] loginctl show-session %s err: %v out: %s", sessionId, err, out)
		}
		return w
	}

	parseLoginCtlOutput(bytes.NewReader(out), w)
	return w
}

func parseLoginCtlOutput(b *bytes.Reader, w *WhoLogged) {
	scanner := bufio.NewScanner(b)
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.IndexByte(line, '='); idx != -1 {
			key := line[:idx]
			val := line[idx+1:]
			switch key {
			case "Name":
				if val != "" {
					w.User = val
				}
			case "TTY":
				if val != "" {
					w.Device = val
				}
			case "RemoteHost":
				if val != "" && val != "n/a" {
					w.IP = val
				}
			case "RemoteAddress":
				// if RemoteHost is empty
				if w.IP == "" && val != "" && val != "n/a" {
					w.IP = val
				}
			}
		}
	}
}
