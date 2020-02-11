package hiera

// ProtoVersion is the protocol version used in the initial negotionation between Hiera and a plugin
const ProtoVersion = 1

// MagicCookie is a value that must be set in the environment variable HIERA_MAGIC_COOKIE in order to run the
// plugin. If it is not set, the plugin will terminate with a message informing the user that it isn't intended for
// normal execution.
const MagicCookie = 0xBEBAC0DE
