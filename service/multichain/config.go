package multichain

import (
	"github.com/mikeydub/go-gallery/service/persist"
)

type ProviderLookup map[persist.Chain]any

type ChainProvider struct {
	Ethereum *EthereumProvider
	Tezos    *TezosProvider
	Optimism *OptimismProvider
	Arbitrum *ArbitrumProvider
	Poap     *PoapProvider
	Zora     *ZoraProvider
	Base     *BaseProvider
	Polygon  *PolygonProvider
}

type EthereumProvider struct {
	ContractRefresher
	ContractFetcher
	ContractsOwnerFetcher
	TokenDescriptorsFetcher
	TokenMetadataFetcher
	TokensIncrementalContractFetcher
	TokensIncrementalOwnerFetcher
	TokenIdentifierOwnerFetcher
	TokenMetadataBatcher
	Verifier
}

type TezosProvider struct {
	ContractsOwnerFetcher
	TokenDescriptorsFetcher
	TokenMetadataFetcher
	TokensIncrementalContractFetcher
	TokensIncrementalOwnerFetcher
	TokenIdentifierOwnerFetcher
	Verifier
}

type OptimismProvider struct {
	TokenDescriptorsFetcher
	TokenMetadataFetcher
	TokensIncrementalContractFetcher
	TokensIncrementalOwnerFetcher
	TokenIdentifierOwnerFetcher
	TokenMetadataBatcher
}

type ArbitrumProvider struct {
	TokenDescriptorsFetcher
	TokenMetadataFetcher
	TokensIncrementalContractFetcher
	TokensIncrementalOwnerFetcher
	TokenIdentifierOwnerFetcher
	TokenMetadataBatcher
}

type PoapProvider struct {
	TokenDescriptorsFetcher
	TokenMetadataFetcher
	TokensIncrementalOwnerFetcher
	TokenIdentifierOwnerFetcher
}

type ZoraProvider struct {
	ContractFetcher
	ContractsOwnerFetcher
	TokenDescriptorsFetcher
	TokenMetadataFetcher
	TokensIncrementalContractFetcher
	TokensIncrementalOwnerFetcher
	TokenIdentifierOwnerFetcher
	TokenMetadataBatcher
}

type BaseProvider struct {
	TokenDescriptorsFetcher
	TokenMetadataFetcher
	TokensIncrementalContractFetcher
	TokensIncrementalOwnerFetcher
	TokenIdentifierOwnerFetcher
	TokenMetadataBatcher
}

type PolygonProvider struct {
	TokenDescriptorsFetcher
	TokenMetadataFetcher
	TokensIncrementalContractFetcher
	TokensIncrementalOwnerFetcher
	TokenIdentifierOwnerFetcher
	TokenMetadataBatcher
}
