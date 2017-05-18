package gomipush
import (
	"sync"
	"time"
	"math/rand"
	"strings"
	"strconv"
	    log "github.com/Sirupsen/logrus"

)

type Server struct {
	mu sync.Mutex
    host string
    priority int
    minPriority int
    maxPriority int
    decrStep int
    incrStep int
}
func (s *Server) GetHost() string {
	return s.host
}
func NewServer(host string,  minPriority int, maxPriority int, decrStep int, incrStep int) *Server {
	s := &Server{}
	s.host = host
	s.priority = maxPriority
	s.minPriority = minPriority
	s.maxPriority = maxPriority
	s.decrStep = decrStep
	s.incrStep = incrStep
	return s
}

func (s *Server) changePriority(incr bool, step int) {
	defer s.mu.Unlock()
	s.mu.Lock()

	var changePriority int
	if incr {
		changePriority = s.priority + step
	}else {
		changePriority = s.priority - step
	}
    if changePriority < s.minPriority {
        changePriority = s.minPriority
    }else if changePriority > s.maxPriority {
        changePriority = s.maxPriority
    }
    s.priority = changePriority
}

func (s *Server) IncrPriority() {
	s.changePriority(true, s.incrStep)
}
func (s *Server) DecrPriority() {
	s.changePriority(false, s.decrStep)
}
type ServerSwitch struct {
	feedback *Server
	sandbox *Server
	specified *Server
	emq *Server
	defaultServer *Server
	servers []*Server
	inited bool
	lastRefreshTime time.Time
}

func newServerSwitch() *ServerSwitch {
	ss := &ServerSwitch{}
 	ss.feedback = NewServer(host_production_feedback, 100, 100, 0, 0)
 	ss.sandbox = NewServer(host_sandbox, 100, 100, 0, 0)
 	ss.specified = NewServer(host, 100, 100, 0, 0)
 	ss.emq = NewServer(host_emq, 100, 100, 0, 0)
 	ss.defaultServer = NewServer(host_production, 1, 90, 10, 5)
 	ss.servers = make([]*Server, 0)
 	ss.inited = false
 	ss.lastRefreshTime = time.Now()
 	return ss
}

func (ss *ServerSwitch) NeedRefreshHostList() bool{
     return  !ss.inited || time.Now().Sub(ss.lastRefreshTime).Seconds() >= refresh_server_host_interval
}

func (ss *ServerSwitch) Initialize(hostList string) {
	if !ss.NeedRefreshHostList() {
		return
	}
	vs := strings.Split(hostList, ",")
	for _, s := range vs {
		sp := strings.Split(s, ":")
		if len(sp) < 5 {
			ss.servers = append(ss.servers, ss.defaultServer)
			continue
		}
		minp, err := strconv.Atoi(sp[1])
		maxp, err := strconv.Atoi(sp[2])
		ds, err := strconv.Atoi(sp[3])
		is, err := strconv.Atoi(sp[4])
		if err != nil {
			continue
		}
		ss.servers = append(ss.servers, NewServer(sp[0], minp, maxp, ds, is))
		log.Infof("minP: %v, maxP: %v, DescS: %v, IncrS: %v, serverLenght: %v", minp, maxp, ds, is, len(ss.servers))

	}
	ss.inited = true
	ss.lastRefreshTime = time.Now()
}
func (ss *ServerSwitch) SelectServer(requestPath []string) *Server{
	if len(host) > 0 {
		return ss.specified
	}

	if is_sandbox{
		return ss.sandbox
	}

	if len(requestPath) == 2 {
		if requestPath[1] == "2"{
			return ss.feedback
		}
		if requestPath[1] == "3"{
			return ss.emq
		}
		return ss.selectServer()
	}

	return ss.selectServer()
}
func (ss *ServerSwitch) selectServer() *Server{
	if  !auto_switch_host || !ss.inited{
		return ss.defaultServer
	}
	var allPriority int = 0
	var priority []int = make([]int, 0)

	for  _, i := range ss.servers {
		priority = append(priority, i.priority)
		allPriority += i.priority
	}

	randomPoint := rand.Intn(allPriority)
	var sum int = 0
	for idx, i:= range priority {
		sum += i
		if randomPoint <= sum {
			log.Infof("select server: %v", ss.servers[idx])
			return ss.servers[idx]
		}
	}
	return ss.defaultServer
}

var ss *ServerSwitch
var once sync.Once
func NewServerSwitch() *ServerSwitch{
	once.Do(func(){
		ss = newServerSwitch()
	})
	return ss
}

