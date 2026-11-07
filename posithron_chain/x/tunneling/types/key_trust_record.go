package types

import "cosmossdk.io/collections"

// TrustRecordKey is the prefix to retrieve all TrustRecord
var TrustRecordKey = collections.NewPrefix("trustRecord/value/")
