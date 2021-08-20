package ReplicatedServiceClient

type Args struct {

}

type Reply struct {

}
type ReplicatedClient interface {
	Init(servers []string, callOne func(string,Args) Reply)

	Call(args Args) Reply
}
