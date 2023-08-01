//go:build trace

package z80

import (
	"fmt"
	"sort"
)

func SymbolForAddress(cpu *CPU, addr uint16) (string, uint16) {
	if cpu.symbols == nil {
		return "", 0
	}

	if len(cpu.addressCache) == 0 {
		cpu.addressCache = make([]int, len(*cpu.symbols))
		index := 0
		for i := range *cpu.symbols {
			cpu.addressCache[index] = int(i)
			index++
		}

		sort.Ints(cpu.addressCache)
	}

	pos := sort.SearchInts(cpu.addressCache, int(cpu.Reg.PC))

	if pos < len(*cpu.symbols) {
		if cpu.addressCache[pos] > int(cpu.Reg.PC) && pos > 0 {
			pos--
		}

		a := uint16(cpu.addressCache[pos])
		return (*cpu.symbols)[a], a
	}

	return "???", addr
}

func SymbolForAddressRelative(cpu *CPU, reg uint16) string {
	symbol, addr := SymbolForAddress(cpu, reg)
	return fmt.Sprintf("%s+%04x", symbol, reg-addr)
}
