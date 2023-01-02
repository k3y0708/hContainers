package types

type FlagsServer struct {
	Sku         string `long:"sku" description:"Server SKU" default:"cx11" choice:"cx11" choice:"cpx11" choice:"cx21" choice:"cpx21" choice:"cx31" choice:"cpx31" choice:"cx41" choice:"cpx41" choice:"cx51" choice:"cpx51"`
	Location    string `long:"location" description:"Server location" default:"Falkenstein" choice:"Falkenstein" choice:"Nuremberg" choice:"Helsinki"`
	DisableIPv6 bool   `long:"disable-ipv6" description:"Disable IPv6"`
}

type HelpRow struct {
	Command     string
	Description string
}
