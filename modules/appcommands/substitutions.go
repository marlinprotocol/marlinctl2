// THIS IS MESSY SYSTEM. DO NOT USE IT IF YOU DO NOT ABSOLUTELY HAVE TO

package appcommands

import "github.com/marlinprotocol/ctl2/modules/util"

// ----------------- RELAY ETH -------------------------------------

func (a *app) relayEthCreateSubstitutions(runnerID string) {
	if a.ProjectID != "relay_eth" {
		return
	}

	runtimeArgs := a.CreateCmd.getStringToStringFromArgStoreOrDie("runtime-args")
	if len(runtimeArgs) == 0 {
		runtimeArgs["DiscoveryAddrs"] = a.CreateCmd.getStringFromArgStoreOrDie("discovery-addrs")
		runtimeArgs["HeartbeatAddrs"] = a.CreateCmd.getStringFromArgStoreOrDie("heartbeat-addrs")
		runtimeArgs["DataDir"] = util.ExpandTilde(a.CreateCmd.getStringFromArgStoreOrDie("datadir"))
		runtimeArgs["DiscoveryPort"] = a.CreateCmd.getStringFromArgStoreOrDie("discovery-port")
		runtimeArgs["PubsubPort"] = a.CreateCmd.getStringFromArgStoreOrDie("pubsub-port")
		runtimeArgs["Address"] = a.CreateCmd.getStringFromArgStoreOrDie("address")
		runtimeArgs["Name"] = a.CreateCmd.getStringFromArgStoreOrDie("name")
		runtimeArgs["AbciVersion"] = a.CreateCmd.getStringFromArgStoreOrDie("abci-version")
		runtimeArgs["SyncMode"] = a.CreateCmd.getStringFromArgStoreOrDie("sync-mode")

		a.CreateCmd.ArgStore["runtime-args"] = runtimeArgs
	}
}

// --------------------- BEACON ------------------------------------

func (a *app) beaconCreateSusbstitutions(runnerID string) {
	if a.ProjectID != "beacon" {
		return
	}

	runtimeArgs := a.CreateCmd.getStringToStringFromArgStoreOrDie("runtime-args")
	if len(runtimeArgs) == 0 {
		runtimeArgs["DiscoveryAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("discovery-addr")
		runtimeArgs["HeartbeatAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("heartbeat-addr")
		runtimeArgs["BootstrapAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("bootstrap-addr")
		runtimeArgs["KeystorePath"] = util.ExpandTilde(a.CreateCmd.getStringFromArgStoreOrDie("keystore-path"))
		runtimeArgs["KeystorePassPath"] = util.ExpandTilde(a.CreateCmd.getStringFromArgStoreOrDie("keystore-pass-path"))

		a.CreateCmd.ArgStore["runtime-args"] = runtimeArgs
	}
}

// --------------------- GATEWAY_DOT -------------------------------

func (a *app) gatewayDotCreateSubstitutions(runnerID string) {
	if a.ProjectID != "gateway_dot" {
		return
	}

	runtimeArgs := a.CreateCmd.getStringToStringFromArgStoreOrDie("runtime-args")
	if len(runtimeArgs) == 0 {
		runtimeArgs["ChainIdentity"] = util.ExpandTilde(a.CreateCmd.getStringFromArgStoreOrDie("chain-identity"))
		runtimeArgs["ListenAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("listen-addr")
		runtimeArgs["DiscoveryAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("discovery-addr")
		runtimeArgs["PubsubAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("pubsub-addr")
		runtimeArgs["BootstrapAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("bootstrap-addr")
		runtimeArgs["InternalListenAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("internal-listen-addr")
		runtimeArgs["KeystorePath"] = util.ExpandTilde(a.CreateCmd.getStringFromArgStoreOrDie("keystore-path"))
		runtimeArgs["KeystorePassPath"] = util.ExpandTilde(a.CreateCmd.getStringFromArgStoreOrDie("keystore-pass-path"))
		runtimeArgs["Contracts"] = a.CreateCmd.getStringFromArgStoreOrDie("contracts")

		a.CreateCmd.ArgStore["runtime-args"] = runtimeArgs
	}
}

// --------------------- GATEWAY_MATICBOR -------------------------------

func (a *app) gatewayMaticBorCreateSubstitutions(runnerID string) {
	if a.ProjectID != "gateway_maticbor" {
		return
	}
	runtimeArgs := a.CreateCmd.getStringToStringFromArgStoreOrDie("runtime-args")
	if len(runtimeArgs) == 0 {

		runtimeArgs["DiscoveryAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("discovery-addr")
		runtimeArgs["PubsubAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("pubsub-addr")
		runtimeArgs["BootstrapAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("bootstrap-addr")
		runtimeArgs["KeystorePath"] = util.ExpandTilde(a.CreateCmd.getStringFromArgStoreOrDie("keystore-path"))
		runtimeArgs["KeystorePassPath"] = util.ExpandTilde(a.CreateCmd.getStringFromArgStoreOrDie("keystore-pass-path"))

		a.CreateCmd.ArgStore["runtime-args"] = runtimeArgs
	}
}

// --------------------- GATEWAY_NEAR -------------------------------

func (a *app) gatewayNearCreateSubstitutions(runnerID string) {
	if a.ProjectID != "gateway_near" {
		return
	}
	runtimeArgs := a.CreateCmd.getStringToStringFromArgStoreOrDie("runtime-args")
	if len(runtimeArgs) == 0 {
		runtimeArgs["ChainIdentity"] = util.ExpandTilde(a.CreateCmd.getStringFromArgStoreOrDie("chain-identity"))
		runtimeArgs["ListenAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("listen-addr")
		runtimeArgs["DiscoveryAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("discovery-addr")
		runtimeArgs["PubsubAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("pubsub-addr")
		runtimeArgs["BootstrapAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("bootstrap-addr")
		runtimeArgs["KeystorePath"] = util.ExpandTilde(a.CreateCmd.getStringFromArgStoreOrDie("keystore-path"))
		runtimeArgs["KeystorePassPath"] = util.ExpandTilde(a.CreateCmd.getStringFromArgStoreOrDie("keystore-pass-path"))
		runtimeArgs["Contracts"] = a.CreateCmd.getStringFromArgStoreOrDie("contracts")

		a.CreateCmd.ArgStore["runtime-args"] = runtimeArgs
	}
}

// --------------------- GATEWAY_IRIS -------------------------------

func (a *app) gatewayIrisCreateSubstitutions(runnerID string) {
	if a.ProjectID != "gateway_iris" {
		return
	}
	runtimeArgs := a.CreateCmd.getStringToStringFromArgStoreOrDie("runtime-args")
	if len(runtimeArgs) == 0 {
		runtimeArgs["DiscoveryAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("discovery-addr")
		runtimeArgs["PubsubAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("pubsub-addr")
		runtimeArgs["BridgeBootstrapAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("bootstrap-addr")
		runtimeArgs["InternalListenAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("internal-listen-addr")
		runtimeArgs["KeystorePath"] = util.ExpandTilde(a.CreateCmd.getStringFromArgStoreOrDie("keystore-path"))
		runtimeArgs["KeystorePassPath"] = util.ExpandTilde(a.CreateCmd.getStringFromArgStoreOrDie("keystore-pass-path"))
		runtimeArgs["Contracts"] = a.CreateCmd.getStringFromArgStoreOrDie("contracts")
		runtimeArgs["GatewayListenPortPeer"] = a.CreateCmd.getStringFromArgStoreOrDie("gateway-listen-port-peer")
		runtimeArgs["GatewayDirection"] = a.CreateCmd.getStringFromArgStoreOrDie("gateway-direction")

		a.CreateCmd.ArgStore["runtime-args"] = runtimeArgs
	}
}

// --------------------- GATEWAY_COSMOS -------------------------------

func (a *app) gatewayCosmosCreateSubstitutions(runnerID string) {
	if a.ProjectID != "gateway_cosmos" {
		return
	}
	runtimeArgs := a.CreateCmd.getStringToStringFromArgStoreOrDie("runtime-args")
	if len(runtimeArgs) == 0 {
		runtimeArgs["DiscoveryAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("discovery-addr")
		runtimeArgs["PubsubAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("pubsub-addr")
		runtimeArgs["BridgeBootstrapAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("bootstrap-addr")
		runtimeArgs["InternalListenAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("internal-listen-addr")
		runtimeArgs["KeystorePath"] = util.ExpandTilde(a.CreateCmd.getStringFromArgStoreOrDie("keystore-path"))
		runtimeArgs["KeystorePassPath"] = util.ExpandTilde(a.CreateCmd.getStringFromArgStoreOrDie("keystore-pass-path"))
		runtimeArgs["Contracts"] = a.CreateCmd.getStringFromArgStoreOrDie("contracts")
		runtimeArgs["GatewayListenPortPeer"] = a.CreateCmd.getStringFromArgStoreOrDie("gateway-listen-port-peer")
		runtimeArgs["GatewayDirection"] = a.CreateCmd.getStringFromArgStoreOrDie("gateway-direction")

		a.CreateCmd.ArgStore["runtime-args"] = runtimeArgs
	}
}

// --------------------- RELAY_IRIS -------------------------------

func (a *app) relayIrisCreateSubstitutions(runnerID string) {
	if a.ProjectID != "relay_iris" {
		return
	}
	runtimeArgs := a.CreateCmd.getStringToStringFromArgStoreOrDie("runtime-args")
	if len(runtimeArgs) == 0 {
		runtimeArgs["DiscoveryAddrs"] = a.CreateCmd.getStringFromArgStoreOrDie("discovery-addrs")
		runtimeArgs["HeartbeatAddrs"] = a.CreateCmd.getStringFromArgStoreOrDie("heartbeat-addrs")
		runtimeArgs["DiscoveryBindAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("discovery-bind-addr")
		runtimeArgs["PubsubBindAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("pubsub-bind-addr")

		a.CreateCmd.ArgStore["runtime-args"] = runtimeArgs
	}
}

// --------------------- RELAY_COSMOS -------------------------------

func (a *app) relayCosmosCreateSubstitutions(runnerID string) {
	if a.ProjectID != "relay_cosmos" {
		return
	}
	runtimeArgs := a.CreateCmd.getStringToStringFromArgStoreOrDie("runtime-args")
	if len(runtimeArgs) == 0 {
		runtimeArgs["DiscoveryAddrs"] = a.CreateCmd.getStringFromArgStoreOrDie("discovery-addrs")
		runtimeArgs["HeartbeatAddrs"] = a.CreateCmd.getStringFromArgStoreOrDie("heartbeat-addrs")
		runtimeArgs["DiscoveryBindAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("discovery-bind-addr")
		runtimeArgs["PubsubBindAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("pubsub-bind-addr")

		a.CreateCmd.ArgStore["runtime-args"] = runtimeArgs
	}
}

// --------------------- RELAY_DOT -------------------------------

func (a *app) relayDotCreateSubstitutions(runnerID string) {
	if a.ProjectID != "relay_dot" {
		return
	}
	runtimeArgs := a.CreateCmd.getStringToStringFromArgStoreOrDie("runtime-args")
	if len(runtimeArgs) == 0 {
		runtimeArgs["DiscoveryAddrs"] = a.CreateCmd.getStringFromArgStoreOrDie("discovery-addrs")
		runtimeArgs["HeartbeatAddrs"] = a.CreateCmd.getStringFromArgStoreOrDie("heartbeat-addrs")
		runtimeArgs["DiscoveryBindAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("discovery-bind-addr")
		runtimeArgs["PubsubBindAddr"] = a.CreateCmd.getStringFromArgStoreOrDie("pubsub-bind-addr")

		a.CreateCmd.ArgStore["runtime-args"] = runtimeArgs
	}
}
