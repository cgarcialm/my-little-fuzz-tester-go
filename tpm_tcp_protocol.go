package main

const (
	MaxBufferSize = 1048576
)

type TcpTpmCommand uint32

const (
	SignalPowerOn TcpTpmCommand = iota + 1
	SignalPowerOff
	SignalPPOn
	SignalPPOff
	SignalHashStart
	SignalHashData
	SignalHashEnd
	SendCommand
	SignalCancelOn
	SignalCancelOff
	SignalNvOn
	SignalNvOff
	SignalKeyCacheOn
	SignalKeyCacheOff
	RemoteHandshake
	SetAlternativeResult
	SignalReset
	SignalRestart
	_
	SessionEnd
	Stop
	_
	_
	_
	GetCommandResponseSizes
	ActGetSignaled
	_
	_
	_
	TestFailureMode
	_
	_
	_
	_
	SetFwHash
	SetFwSvn
)

type TpmEndPointInfo uint

const (
	// Platform hierarchy is enabled, and hardware platform functionality (such
	// as SignalHashStart/Data/End) is available.
	PlatformAvailable TpmEndPointInfo = 0x01

	// The device is TPM Resource Manager (TRM), rather than a raw TPM.
	// This means context management commands are unavailable, and the handle values
	// returned to the client are virtualized.
	UsesTbs TpmEndPointInfo = 0x02

	// The TRM is in raw mode (i.e. no actual resourse virtualization is performed).
	InRawMode TpmEndPointInfo = 0x04

	// Phisical presence signals (SignalPPOn/Off) are supported.
	SupportsPP TpmEndPointInfo = 0x08
)
