package processservice

// A ProcessService does nothing.
type ProcessService struct {
	ProcessName string
	Executed    bool
}

type ProcessServiceList struct {
	items []*ProcessService
}

func (pl *ProcessServiceList) GetItems() []*ProcessService {
	return pl.items
}

func (pl *ProcessServiceList) AddService(service *ProcessService) {
	pl.items = append(pl.items, service)

}

var instances *ProcessServiceList

func GetInstance(processName string) *ProcessService {

	// Pattern singleton for each "process" created
	if instances == nil {
		instances = &ProcessServiceList{}
	}
	for _, process := range instances.GetItems() {
		if process.GetName() == processName {
			return process
		}
	}
	newProcess := &ProcessService{ProcessName: processName, Executed: false}
	instances.AddService(newProcess)
	return newProcess
}

func (p *ProcessService) UpdateExecuted(condition bool) {
	p.Executed = condition
}
func (p *ProcessService) GetExecuted() bool {
	return p.Executed
}

func (p *ProcessService) GetName() string {
	return p.ProcessName
}

func (p *ProcessService) SetName(name string) {
	p.ProcessName = name
}

func DeleteProcessService(processName string) {
	if instances == nil {
		return
	}
	for i, process := range instances.GetItems() {
		if process.GetName() == processName {
			remove(instances.items, i)
		}
	}
}
func remove(s []*ProcessService, i int) []*ProcessService {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
