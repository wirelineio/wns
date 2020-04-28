// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

type Extension interface {
	IsExtension()
}

type Account struct {
	Address  string  `json:"address"`
	PubKey   *string `json:"pubKey"`
	Number   BigUInt `json:"number"`
	Sequence BigUInt `json:"sequence"`
	Balance  []Coin  `json:"balance"`
}

type Bond struct {
	ID      string `json:"id"`
	Owner   string `json:"owner"`
	Balance []Coin `json:"balance"`
}

type Bot struct {
	Name      string `json:"name"`
	AccessKey string `json:"accessKey"`
}

func (Bot) IsExtension() {}

type Coin struct {
	Type     string  `json:"type"`
	Quantity BigUInt `json:"quantity"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value Value  `json:"value"`
}

type KeyValueInput struct {
	Key   string     `json:"key"`
	Value ValueInput `json:"value"`
}

type NodeInfo struct {
	ID      string `json:"id"`
	Network string `json:"network"`
	Moniker string `json:"moniker"`
}

type Pad struct {
	Name string `json:"name"`
}

func (Pad) IsExtension() {}

type PeerInfo struct {
	Node       NodeInfo `json:"node"`
	IsOutbound bool     `json:"is_outbound"`
	RemoteIP   string   `json:"remote_ip"`
}

type Protocol struct {
	Name string `json:"name"`
}

func (Protocol) IsExtension() {}

type Record struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Name       string      `json:"name"`
	Version    string      `json:"version"`
	BondID     string      `json:"bondId"`
	CreateTime string      `json:"createTime"`
	ExpiryTime string      `json:"expiryTime"`
	Owners     []*string   `json:"owners"`
	Attributes []*KeyValue `json:"attributes"`
	References []*Record   `json:"references"`
	Extension  Extension   `json:"extension"`
}

type Reference struct {
	ID string `json:"id"`
}

type ReferenceInput struct {
	ID string `json:"id"`
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

type UnknownExtension struct {
	Name string `json:"name"`
}

func (UnknownExtension) IsExtension() {}

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
