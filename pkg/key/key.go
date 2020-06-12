package key

type SSHKey struct {
	Title string
	Key string `json:"key"`
	ReadOnly bool
}