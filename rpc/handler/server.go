package handler

func Register(group, addr string) error {
	sgm := GetSGMInst()
	sgm.Register(group, addr)
	//TODO activate sync process
	return nil
}

func Unregister(group, addr string) error {
	sgm := GetSGMInst()
	sgm.Unregister(group, addr)
	//TODO activate sync process
	return nil
}

func SyncSrvGroups(srvgroups []byte) error {
	//TODO handle sync request
	return nil
}
