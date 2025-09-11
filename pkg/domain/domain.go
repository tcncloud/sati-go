package domain

type Domain struct{}

func NewDomain() *Domain {
	return &Domain{}
}

func (d *Domain) StartConfigWatcher() {
}

func (d *Domain) StartGateClient() {

}

func (d *Domain) StartPollEvents() {

}

func (d *Domain) StartStreamJobs() {

}
