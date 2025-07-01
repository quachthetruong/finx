package entity

type FlowType string

const (
	FLowTypeDnseOffline FlowType = "001"
	FlowTypeDnseOnline  FlowType = "002"
)

func (f FlowType) String() string {
	return string(f)
}

func FlowTypeFromString(s string) FlowType {
	return FlowType(s)
}

type AssetType string

const (
	AssetTypeUnderlying AssetType = "UNDERLYING"
	AssetTypeDerivative AssetType = "DERIVATIVE"
)

func (a AssetType) String() string {
	return string(a)
}

func AssetTypeFromString(s string) AssetType {
	switch s {
	case "UNDERLYING":
		return AssetTypeUnderlying
	case "DERIVATIVE":
		return AssetTypeDerivative
	default:
		return ""
	}
}

type ActionType string

const (
	Approve                    ActionType = "APPROVE"
	Reject                     ActionType = "REJECT"
	RejectAndSendOtherProposal ActionType = "REJECT_AND_SEND_OTHER_PROPOSAL"
)

const (
	Online  FlowType = "ONLINE"
	Offline FlowType = "OFFLINE"
)

type ProposeType string

const (
	NewLoanPackage      ProposeType = "NEW_LOAN_PACKAGE"
	ExistingLoanPackage ProposeType = "EXISTING_LOAN_PACKAGE"
)

type PromotionCampaignStatus string

const (
	Active   PromotionCampaignStatus = "ACTIVE"
	Inactive PromotionCampaignStatus = "INACTIVE"
)

func (f PromotionCampaignStatus) String() string {
	return string(f)
}

func PromotionCampaignStatusFromString(s string) PromotionCampaignStatus {
	return PromotionCampaignStatus(s)
}
