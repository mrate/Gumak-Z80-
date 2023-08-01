//go:build !trace

package z80

func SymbolForAddressRelative(cpu *CPU, reg uint16) string {
	return ""
}
