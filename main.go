package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "strconv"
    "sync"
    "log"
)

var PORT_SCAN_FILE = [...]string {"/proc/net/raw", "/proc/net/raw6", "/proc/net/tcp", "/proc/net/tcp6", "/proc/net/udp", "/proc/net/udp6", "/proc/net/udplite", "/proc/net/udplite6"}
//var PORT_SCAN_FILE = [...]string {"./raw", "./raw6", "./tcp", "./tcp6", "./udp", "./udp6", "./udplite", "./udplite6"}
var SPECIAL_AVAILABLE_PORTS = [...]int {80, 443, 53, 8080, 2379, 4001}
const AVAILABLE_PORT_START int = 20000
const AVAILABLE_PORT_END int = 20009

type Set struct {
    m map[interface{}]bool
    sync.RWMutex
}

func NewSet() *Set {
    return &Set {
        m: map[interface{}]bool{},
    }
}

func (s *Set) Add(item interface{}) {
    s.Lock()
    defer s.Unlock()
    s.m[item] = true
}

func (s *Set) Remove(item interface{}) {
    s.Lock()
    s.Unlock()
    delete(s.m, item)
}

func (s *Set) Has(item interface{}) bool {
    s.RLock()
    defer s.RUnlock()
    _, ok := s.m[item]
    return ok
}

func (s *Set) Len() int {
    return len(s.List())
}

func (s *Set) Clear() {
    s.m = map[interface{}]bool{}
}

func (s *Set) IsEmpty() bool {
    if s.Len() == 0 {
        return true
    }
    return false
}

func (s *Set) List() []interface{} {
    s.RLock()
    defer s.RUnlock()
    list := []interface{}{}
    for item := range s.m {
        list = append(list, item)
    }
    return list
}

func main() {
    var portSet *Set = NewSet()
    for _, file := range PORT_SCAN_FILE {
        portInfoStr, err := ioutil.ReadFile(file)
        if err != nil {
            log.Println("[ERROR] read file error: ", err.Error())
        }
        portInfo := strings.Split(string(portInfoStr), "\n")
        for idx, line := range portInfo {
            if idx > 0 && len(line) > 3 {
                port16 := strings.Split(strings.Split(line, ":")[2], " ")[0]
                portTmp, convErr := strconv.ParseInt(port16, 16, 32)
                port := int(portTmp)
                if convErr != nil {
                  log.Println("[Warn] Get Available Ports: ", convErr.Error())
                } else if !portSet.Has(port) {
                    portSet.Add(port)
                }
            }
        }
    }
    var availablePorts []int
    aport := AVAILABLE_PORT_START
    for aport <= AVAILABLE_PORT_END {
        if !portSet.Has(aport) {
            availablePorts = append(availablePorts, aport)
        }
        aport++
    }
    for _, bport := range SPECIAL_AVAILABLE_PORTS {
        if !portSet.Has(bport) {
            availablePorts = append(availablePorts, bport)
        }
    }
    for _, port := range availablePorts {
        fmt.Println(port)
    }
}
