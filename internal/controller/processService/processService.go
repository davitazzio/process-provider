package processservice

// // A ProcessService does nothing.
// type ProcessService struct {
// 	ProcessName string
// 	Active      bool
// }

// type ProcessServiceList struct {
// 	items []*ProcessService
// }

// func (pl *ProcessServiceList) GetItems() []*ProcessService {

// 	return pl.items
// }

// func (pl *ProcessServiceList) AddService(service *ProcessService) {

// 	pl.items = append(pl.items, service)

// }

// var instances *ProcessServiceList

// func (pl *ProcessServiceList) GetInstance(processName string) *ProcessService {

// 	if instances == nil {
// 		instances = &ProcessServiceList{}
// 	}
// 	for _, process := range instances.GetItems() {
// 		if process.GetName() == processName {
// 			return process
// 		}
// 	}
// 	newProcess := &ProcessService{ProcessName: processName, Active: false}
// 	instances.AddService(newProcess)
// 	return newProcess
// }

// func (p *ProcessService) UpdateActive(condition bool) {
// 	p.Active = condition
// }
// func (p *ProcessService) GetActive() bool {
// 	return p.Active
// }

// func (p *ProcessService) GetName() string {
// 	return p.ProcessName
// }

// func (p *ProcessService) SetName(name string) {
// 	p.ProcessName = name
// }

// func DeleteProcessService(processName string) {
// 	if instances == nil {
// 		return
// 	}
// 	for i, process := range instances.GetItems() {
// 		if process.GetName() == processName {
// 			remove(instances.items, i)
// 		}
// 	}
// }
// func remove(s []*ProcessService, i int) []*ProcessService {
// 	s[i] = s[len(s)-1]
// 	return s[:len(s)-1]
// }
