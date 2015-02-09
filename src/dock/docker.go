package dock

//A Docker handles incoming communication.
type Docker interface {
	Launch(agent DockAgent, id string) error //Launches a docker
}
