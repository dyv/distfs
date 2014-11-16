package paxos

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
)

// If we record all the messages we receive then we replicate state
type MsgLog struct {
	Log   []Msg
	Fpath string
	fd    *os.File
}

// returns whether this file existed before: if it did then recover the log
// from that file
func NewMsgLog(sz int, path string, a *Agent, try_recover bool) (*MsgLog, error) {
	l := &MsgLog{}
	l.Log = make([]Msg, sz)
	l.Fpath = path
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0777)
	}
	if _, err := os.Stat(path); os.IsExist(err) {
		// file exists therefore recover from it
		if !try_recover {
			// if we aren't supposed to use it to recover
			// assume that we should delete it
			// this is good for testing or running clean instances of Paxos
			err = os.Remove(path)
			if err != nil {
				log.Print("Error Deleting Log File")
			}
		} else {
			log.Print("File Exists")
			f, err := os.Open(path)
			if err != nil {
				log.Print("Error Opening File Exists")
				return nil, err
			}
			err = l.Recover(a, f)
			if err != nil {
				return nil, err
			}
			err = f.Close()
			if err != nil {
				return nil, err
			}
		}
	}
	log.Print("File Does not Exist: ", path)
	var err error
	l.fd, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Print("Failed to Open File")
		return nil, err
	}
	return l, nil
}

func (l *MsgLog) Recover(a *Agent, f *os.File) error {
	dec := json.NewDecoder(f)
	for {
		var m Msg
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		log.Print("Recovered: ", m)
		// append the recovered message to the in memory log
		l.Log = append(l.Log, m)
		if m.Type == ClientRequest {
			a.StartRequest(m.Round, m.Value, m.Request, false)
		} else {
			a.handleMessage(m, false)
		}
	}
	// after recovering never assume that I am the leader
	a.isLeader = false
	a.leader = nil
	return nil
}

func (l *MsgLog) Resize(n int) {
	tl := make([]Msg, n*2)
	copy(tl, l.Log)
	l.Log = tl
}

func (l *MsgLog) Flush() {
	err := l.fd.Sync()
	if err != nil {
		log.Fatal("Failed to Flush File: ", err)
	}
}

func (l *MsgLog) Append(m Msg) {
	l.Log = append(l.Log, m)
	// Append this one message to the file
	by, err := json.Marshal(m)
	if err != nil {
		log.Fatal("Error Appending To Log:", err)
	}
	_, err = l.fd.Write(by)
	if err != nil {
		log.Fatal("Error Appending To Log:", err)
	}
	l.Flush()
}

type ValueEntry struct {
	Committed bool
	Val       Value
}

// A ValueLog is a sequence of Values that the the Paxos node has accepted
// in the order that it has accepted it
type ValueLog struct {
	Log []ValueEntry
}

func NewValueLog(sz int) *ValueLog {
	l := &ValueLog{}
	l.Log = make([]ValueEntry, 0, sz)
	return l
}

func (l *ValueLog) InsertAt(i int, v Value) {
	if i >= len(l.Log) {
		t := make([]ValueEntry, ((i+1)*3)/2)
		copy(t, l.Log)
		l.Log = t

	}
	l.Log[i] = ValueEntry{true, v}
}

func (l *ValueLog) Append(v Value) {
	l.Log = append(l.Log, ValueEntry{true, v})
}
