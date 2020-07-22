// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

type Account struct {
	Address  string  `json:"address"`
	PubKey   *string `json:"pubKey"`
	Number   string  `json:"number"`
	Sequence string  `json:"sequence"`
	Balance  []Coin  `json:"balance"`
}

type AuthorityRecord struct {
	OwnerAddress   string `json:"ownerAddress"`
	OwnerPublicKey string `json:"ownerPublicKey"`
	Height         string `json:"height"`
}

type AuthorityResult struct {
	Meta    ResultMeta         `json:"meta"`
	Records []*AuthorityRecord `json:"records"`
}

type Bond struct {
	ID      string `json:"id"`
	Owner   string `json:"owner"`
	Balance []Coin `json:"balance"`
}

type Coin struct {
	Type     string `json:"type"`
	Quantity string `json:"quantity"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value Value  `json:"value"`
}

type KeyValueInput struct {
	Key   string     `json:"key"`
	Value ValueInput `json:"value"`
}

type NameRecord struct {
	Latest  NameRecordEntry    `json:"latest"`
	History []*NameRecordEntry `json:"history"`
}

type NameRecordEntry struct {
	ID     string `json:"id"`
	Height string `json:"height"`
}

type NameResult struct {
	Meta    ResultMeta    `json:"meta"`
	Records []*NameRecord `json:"records"`
}

type NodeInfo struct {
	ID      string `json:"id"`
	Network string `json:"network"`
	Moniker string `json:"moniker"`
}

type PeerInfo struct {
	Node       NodeInfo `json:"node"`
	IsOutbound bool     `json:"is_outbound"`
	RemoteIP   string   `json:"remote_ip"`
}

type Record struct {
	ID         string      `json:"id"`
	Names      []string    `json:"names"`
	BondID     string      `json:"bondId"`
	CreateTime string      `json:"createTime"`
	ExpiryTime string      `json:"expiryTime"`
	Owners     []*string   `json:"owners"`
	Attributes []*KeyValue `json:"attributes"`
	References []*Record   `json:"references"`
}

type RecordResult struct {
	Meta    ResultMeta `json:"meta"`
	Records []*Record  `json:"records"`
}

type Reference struct {
	ID string `json:"id"`
}

type ReferenceInput struct {
	ID string `json:"id"`
}

type ResultMeta struct {
	Height string `json:"height"`
}

type Status struct {
	Version    string           `json:"version"`
	Node       NodeInfo         `json:"node"`
	Sync       SyncInfo         `json:"sync"`
	Validator  *ValidatorInfo   `json:"validator"`
	Validators []*ValidatorInfo `json:"validators"`
	NumPeers   string           `json:"num_peers"`
	Peers      []*PeerInfo      `json:"peers"`
	DiskUsage  string           `json:"disk_usage"`
}

type SyncInfo struct {
	LatestBlockHash   string `json:"latest_block_hash"`
	LatestBlockHeight string `json:"latest_block_height"`
	LatestBlockTime   string `json:"latest_block_time"`
	CatchingUp        bool   `json:"catching_up"`
}

type ValidatorInfo struct {
	Address          string  `json:"address"`
	VotingPower      string  `json:"voting_power"`
	ProposerPriority *string `json:"proposer_priority"`
}

type Value struct {
	Null      *bool      `json:"null"`
	Int       *int       `json:"int"`
	Float     *float64   `json:"float"`
	String    *string    `json:"string"`
	Boolean   *bool      `json:"boolean"`
	JSON      *string    `json:"json"`
	Reference *Reference `json:"reference"`
	Values    []*Value   `json:"values"`
}

type ValueInput struct {
	Null      *bool           `json:"null"`
	Int       *int            `json:"int"`
	Float     *float64        `json:"float"`
	String    *string         `json:"string"`
	Boolean   *bool           `json:"boolean"`
	Reference *ReferenceInput `json:"reference"`
	Values    []*ValueInput   `json:"values"`
}
