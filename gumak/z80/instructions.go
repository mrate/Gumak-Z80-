package z80

import (
	"mutex/gumak/helpers"
	"mutex/gumak/log"
)

type InstrOp func(cpu *CPU) int
type InstrOpIXIY func(cpu *CPU, idx *uint16) int
type InstrOpIXIY_cb func(cpu *CPU, idx *uint16, n int) int

var OpCodes []InstrOp

func init() {
	OpCodes = make([]InstrOp, 256)

	instrOpIXIY_Undocumented := func(cpu *CPU, idx *uint16) int {
		panic("Undocumented & unimplemented OP")
	}

	OpCodes_ED := make([]InstrOp, 256)
	OpCodes_IX_IY := make([]InstrOpIXIY, 256)
	OpCodes_IX_IY_cb := make([]InstrOpIXIY_cb, 256)

	for i := 0; i < 256; i++ {
		invalidOp := func(cpu *CPU) int {
			panic("Invalid instruction")
		}

		OpCodes[i], OpCodes_ED[i] = invalidOp, invalidOp

		OpCodes_IX_IY[i] = func(cpu *CPU, idx *uint16) int {
			panic("Invalid instruction")
		}

		OpCodes_IX_IY_cb[i] = func(cpu *CPU, idx *uint16, n int) int {
			panic("Invalid instruction")
		}
	}

	OpCodes[0x0] = NOP
	OpCodes[0x01] = /* ld bc,nn */ func(cpu *CPU) int { return LD_RR_nn(cpu, &cpu.Reg.B, &cpu.Reg.C) }
	OpCodes[0x02] = /* ld (bc),a */ func(cpu *CPU) int { return LD_mem_R_16(cpu, &cpu.Reg.B, &cpu.Reg.C, &cpu.Reg.A) }
	OpCodes[0x03] = /* inc bc */ func(cpu *CPU) int { return INC16(cpu, &cpu.Reg.B, &cpu.Reg.C) }
	OpCodes[0x04] = /* inc b */ func(cpu *CPU) int { return INC(cpu, &cpu.Reg.B) }
	OpCodes[0x05] = /* dec b */ func(cpu *CPU) int { return DEC(cpu, &cpu.Reg.B) }
	OpCodes[0x06] = /* n ld b,n */ func(cpu *CPU) int { return LD_R_n(cpu, &cpu.Reg.B, FetchOperand8(cpu)) }
	OpCodes[0x07] = /* rlca */ func(cpu *CPU) int { Alu_RLC_A(cpu); return 4 }
	OpCodes[0x08] = /* ex af,af_ */ EX_AF_AF_
	OpCodes[0x09] = /* add hl,bc */ func(cpu *CPU) int { return ADD16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.B, &cpu.Reg.C) }
	OpCodes[0x0a] = /* ld a,(bc) */ func(cpu *CPU) int { return LD_R_mem_16(cpu, &cpu.Reg.B, &cpu.Reg.C, &cpu.Reg.A) }
	OpCodes[0x0b] = /* dec bc */ func(cpu *CPU) int { return DEC16(cpu, &cpu.Reg.B, &cpu.Reg.C) }
	OpCodes[0x0c] = /* inc c */ func(cpu *CPU) int { return INC(cpu, &cpu.Reg.C) }
	OpCodes[0x0d] = /* dec c */ func(cpu *CPU) int { return DEC(cpu, &cpu.Reg.C) }
	OpCodes[0x0e] = /* n ld c,n */ func(cpu *CPU) int { return LD_R_n(cpu, &cpu.Reg.C, FetchOperand8(cpu)) }
	OpCodes[0x0f] = /* rrca */ func(cpu *CPU) int {
		Alu_RRC_A(cpu)
		log.Trace(2, "RRC")
		return 4
	}
	OpCodes[0x10] = /* djnz $+2 */ DJNZ_e
	OpCodes[0x11] = /* ld de,nn */ func(cpu *CPU) int { return LD_RR_nn(cpu, &cpu.Reg.D, &cpu.Reg.E) }
	OpCodes[0x12] = /* ld (de),a */ func(cpu *CPU) int { return LD_mem_R_16(cpu, &cpu.Reg.D, &cpu.Reg.E, &cpu.Reg.A) }
	OpCodes[0x13] = /* inc de */ func(cpu *CPU) int { return INC16(cpu, &cpu.Reg.D, &cpu.Reg.E) }
	OpCodes[0x14] = /* inc d */ func(cpu *CPU) int { return INC(cpu, &cpu.Reg.D) }
	OpCodes[0x15] = /* dec d */ func(cpu *CPU) int { return DEC(cpu, &cpu.Reg.D) }
	OpCodes[0x16] = /* n ld d,n */ func(cpu *CPU) int { return LD_R_n(cpu, &cpu.Reg.D, FetchOperand8(cpu)) }
	OpCodes[0x17] = /* rla */ func(cpu *CPU) int {
		Alu_RL_A(cpu)
		log.Trace(2, "RLA")
		return 4
	}
	OpCodes[0x18] = /* jr $+2 */ JR_e
	OpCodes[0x19] = /* add hl,de */ func(cpu *CPU) int { return ADD16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.D, &cpu.Reg.E) }
	OpCodes[0x1a] = /* ld a,(de) */ func(cpu *CPU) int { return LD_R_mem_16(cpu, &cpu.Reg.D, &cpu.Reg.E, &cpu.Reg.A) }
	OpCodes[0x1b] = /* dec de */ func(cpu *CPU) int { return DEC16(cpu, &cpu.Reg.D, &cpu.Reg.E) }
	OpCodes[0x1c] = /* inc e */ func(cpu *CPU) int { return INC(cpu, &cpu.Reg.E) }
	OpCodes[0x1d] = /* dec e */ func(cpu *CPU) int { return DEC(cpu, &cpu.Reg.E) }
	OpCodes[0x1e] = /* n ld e,n */ func(cpu *CPU) int { return LD_R_n(cpu, &cpu.Reg.E, FetchOperand8(cpu)) }
	OpCodes[0x1f] = /* rra */ func(cpu *CPU) int {
		Alu_RR_A(cpu)
		log.Trace(2, "RRA")
		return 4
	}
	OpCodes[0x20] = /* jr nz,$+2 */ func(cpu *CPU) int {
		return JR_FLAG_e(cpu, FLAG_ZERO, false)
	}
	OpCodes[0x21] = /* ld hl,nn */ func(cpu *CPU) int { return LD_R_nn(cpu, &cpu.Reg.H, &cpu.Reg.L) }
	OpCodes[0x22] = /* ld (nn),hl */ func(cpu *CPU) int { return LD_nn_HL_mem(cpu) }
	OpCodes[0x23] = /* inc hl */ func(cpu *CPU) int { return INC16(cpu, &cpu.Reg.H, &cpu.Reg.L) }
	OpCodes[0x24] = /* inc h */ func(cpu *CPU) int { return INC(cpu, &cpu.Reg.H) }
	OpCodes[0x25] = /* dec h */ func(cpu *CPU) int { return DEC(cpu, &cpu.Reg.H) }
	OpCodes[0x26] = /* n ld h,n */ func(cpu *CPU) int { return LD_R_n(cpu, &cpu.Reg.H, FetchOperand8(cpu)) }
	OpCodes[0x27] = /* daa */ DAA
	OpCodes[0x28] = /* jr z,$+2 */ func(cpu *CPU) int {
		return JR_FLAG_e(cpu, FLAG_ZERO, true)
	}
	OpCodes[0x29] = /* add hl,hl */ func(cpu *CPU) int { return ADD16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.H, &cpu.Reg.L) }
	OpCodes[0x2a] = /* ld hl,(nn) */ func(cpu *CPU) int { return LD_HL_nn_mem(cpu) }
	OpCodes[0x2b] = /* dec hl */ func(cpu *CPU) int { return DEC16(cpu, &cpu.Reg.H, &cpu.Reg.L) }
	OpCodes[0x2c] = /* inc l */ func(cpu *CPU) int { return INC(cpu, &cpu.Reg.L) }
	OpCodes[0x2d] = /* dec l */ func(cpu *CPU) int { return DEC(cpu, &cpu.Reg.L) }
	OpCodes[0x2e] = /* n ld l,n */ func(cpu *CPU) int { return LD_R_n(cpu, &cpu.Reg.L, FetchOperand8(cpu)) }
	OpCodes[0x2f] = /* cpl */ func(cpu *CPU) int {
		Alu_CPL_A(cpu)
		log.Trace(2, "CPL")
		return 4
	}
	OpCodes[0x30] = /* jr nc,$+2 */ func(cpu *CPU) int {
		return JR_FLAG_e(cpu, FLAG_CARRY, false)
	}
	OpCodes[0x31] = /* ld sp,nn */ func(cpu *CPU) int {
		nn := FetchOperand16(cpu)
		cpu.Reg.SP = nn

		log.Trace(2, "LD SP, $%04x", nn)
		return 10
	}
	OpCodes[0x32] = /* ld (nn),a */ func(cpu *CPU) int { return LD_mem_R(cpu, &cpu.Reg.A) }
	OpCodes[0x33] = /* inc sp */ func(cpu *CPU) int {
		cpu.Reg.SP++
		log.Trace(2, "INC SP")
		return 4
	}
	OpCodes[0x34] = /* inc (hl) */ INC_HL
	OpCodes[0x35] = /* dec (hl) */ DEC_HL
	OpCodes[0x36] = /* ld (hl),n */ func(cpu *CPU) int {
		n := FetchOperand8(cpu)
		MEM_HL_W(cpu, n)
		log.Trace(2, "LD (HL), $%02x", n)
		return 10
	}
	OpCodes[0x37] = /* scf */ func(cpu *CPU) int {
		Alu_SCF(cpu)
		log.Trace(2, "SCF")
		return 4
	}
	OpCodes[0x38] = /* jr c,$+2 */ func(cpu *CPU) int {
		return JR_FLAG_e(cpu, FLAG_CARRY, true)
	}
	OpCodes[0x39] = /* add hl,sp */ func(cpu *CPU) int {
		cpu.Reg.HL_write(Alu_ADD16(cpu, cpu.Reg.HL(), cpu.Reg.SP))
		log.Trace(2, "ADD HL, SP")
		return 11
	}
	OpCodes[0x3a] = /* ld a,(nn) */ LD_A_nn_mem
	OpCodes[0x3b] = /* dec sp */ func(cpu *CPU) int {
		cpu.Reg.SP--
		log.Trace(2, "DEC SP")
		return 4
	}
	OpCodes[0x3c] = /* inc a */ func(cpu *CPU) int { return INC(cpu, &cpu.Reg.A) }
	OpCodes[0x3d] = /* dec a */ func(cpu *CPU) int { return DEC(cpu, &cpu.Reg.A) }
	OpCodes[0x3e] = /* ld a,n */ func(cpu *CPU) int { return LD_R_n(cpu, &cpu.Reg.A, FetchOperand8(cpu)) }
	OpCodes[0x3f] = /* ccf */ func(cpu *CPU) int {
		Alu_CCF(cpu)
		log.Trace(2, "CCF")
		return 4
	}
	OpCodes[0x40] = /* ld b,b */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.B, &cpu.Reg.B) }
	OpCodes[0x41] = /* ld b,c */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.B, &cpu.Reg.C) }
	OpCodes[0x42] = /* ld b,d */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.B, &cpu.Reg.D) }
	OpCodes[0x43] = /* ld b,e */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.B, &cpu.Reg.E) }
	OpCodes[0x44] = /* ld b,h */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.B, &cpu.Reg.H) }
	OpCodes[0x45] = /* ld b,l */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.B, &cpu.Reg.L) }
	OpCodes[0x46] = /* ld b,(hl) */ func(cpu *CPU) int { return LD_R_mem_16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.B) }
	OpCodes[0x47] = /* ld b,a */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.B, &cpu.Reg.A) }
	OpCodes[0x48] = /* ld c,b */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.C, &cpu.Reg.B) }
	OpCodes[0x49] = /* ld c,c */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.C, &cpu.Reg.C) }
	OpCodes[0x4a] = /* ld c,d */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.C, &cpu.Reg.D) }
	OpCodes[0x4b] = /* ld c,e */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.C, &cpu.Reg.E) }
	OpCodes[0x4c] = /* ld c,h */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.C, &cpu.Reg.H) }
	OpCodes[0x4d] = /* ld c,l */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.C, &cpu.Reg.L) }
	OpCodes[0x4e] = /* ld c,(hl) */ func(cpu *CPU) int { return LD_R_mem_16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.C) }
	OpCodes[0x4f] = /* ld c,a */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.C, &cpu.Reg.A) }
	OpCodes[0x50] = /* ld d,b */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.D, &cpu.Reg.B) }
	OpCodes[0x51] = /* ld d,c */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.D, &cpu.Reg.C) }
	OpCodes[0x52] = /* ld d,d */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.D, &cpu.Reg.D) }
	OpCodes[0x53] = /* ld d,e */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.D, &cpu.Reg.E) }
	OpCodes[0x54] = /* ld d,h */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.D, &cpu.Reg.H) }
	OpCodes[0x55] = /* ld d,l */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.D, &cpu.Reg.L) }
	OpCodes[0x56] = /* ld d,(hl) */ func(cpu *CPU) int { return LD_R_mem_16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.D) }
	OpCodes[0x57] = /* ld d,a */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.D, &cpu.Reg.A) }
	OpCodes[0x58] = /* ld e,b */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.E, &cpu.Reg.B) }
	OpCodes[0x59] = /* ld e,c */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.E, &cpu.Reg.C) }
	OpCodes[0x5a] = /* ld e,d */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.E, &cpu.Reg.D) }
	OpCodes[0x5b] = /* ld e,e */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.E, &cpu.Reg.E) }
	OpCodes[0x5c] = /* ld e,h */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.E, &cpu.Reg.H) }
	OpCodes[0x5d] = /* ld e,l */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.E, &cpu.Reg.L) }
	OpCodes[0x5e] = /* ld e,(hl) */ func(cpu *CPU) int { return LD_R_mem_16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.E) }
	OpCodes[0x5f] = /* ld e,a */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.E, &cpu.Reg.A) }
	OpCodes[0x60] = /* ld h,b */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.H, &cpu.Reg.B) }
	OpCodes[0x61] = /* ld h,c */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.H, &cpu.Reg.C) }
	OpCodes[0x62] = /* ld h,d */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.H, &cpu.Reg.D) }
	OpCodes[0x63] = /* ld h,e */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.H, &cpu.Reg.E) }
	OpCodes[0x64] = /* ld h,h */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.H, &cpu.Reg.H) }
	OpCodes[0x65] = /* ld h,l */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.H, &cpu.Reg.L) }
	OpCodes[0x66] = /* ld h,(hl) */ func(cpu *CPU) int { return LD_R_mem_16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.H) }
	OpCodes[0x67] = /* ld h,a */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.H, &cpu.Reg.A) }
	OpCodes[0x68] = /* ld l,b */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.L, &cpu.Reg.B) }
	OpCodes[0x69] = /* ld l,c */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.L, &cpu.Reg.C) }
	OpCodes[0x6a] = /* ld l,d */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.L, &cpu.Reg.D) }
	OpCodes[0x6b] = /* ld l,e */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.L, &cpu.Reg.E) }
	OpCodes[0x6c] = /* ld l,h */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.L, &cpu.Reg.H) }
	OpCodes[0x6d] = /* ld l,l */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.L, &cpu.Reg.L) }
	OpCodes[0x6e] = /* ld l,(hl) */ func(cpu *CPU) int { return LD_R_mem_16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.L) }
	OpCodes[0x6f] = /* ld l,a */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.L, &cpu.Reg.A) }
	OpCodes[0x70] = /* ld (hl),b */ func(cpu *CPU) int { return LD_mem_R_16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.B) }
	OpCodes[0x71] = /* ld (hl),c */ func(cpu *CPU) int { return LD_mem_R_16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.C) }
	OpCodes[0x72] = /* ld (hl),d */ func(cpu *CPU) int { return LD_mem_R_16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.D) }
	OpCodes[0x73] = /* ld (hl),e */ func(cpu *CPU) int { return LD_mem_R_16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.E) }
	OpCodes[0x74] = /* ld (hl),h */ func(cpu *CPU) int { return LD_mem_R_16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.H) }
	OpCodes[0x75] = /* ld (hl),l */ func(cpu *CPU) int { return LD_mem_R_16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.L) }
	OpCodes[0x76] = /* halt */ HALT
	OpCodes[0x77] = /* ld (hl),a */ func(cpu *CPU) int { return LD_mem_R_16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.A) }
	OpCodes[0x78] = /* ld a,b */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.A, &cpu.Reg.B) }
	OpCodes[0x79] = /* ld a,c */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.A, &cpu.Reg.C) }
	OpCodes[0x7a] = /* ld a,d */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.A, &cpu.Reg.D) }
	OpCodes[0x7b] = /* ld a,e */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.A, &cpu.Reg.E) }
	OpCodes[0x7c] = /* ld a,h */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.A, &cpu.Reg.H) }
	OpCodes[0x7d] = /* ld a,l */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.A, &cpu.Reg.L) }
	OpCodes[0x7e] = /* ld a,(hl) */ func(cpu *CPU) int { return LD_R_mem_16(cpu, &cpu.Reg.H, &cpu.Reg.L, &cpu.Reg.A) }
	OpCodes[0x7f] = /* ld a,a */ func(cpu *CPU) int { return LD_R_R(cpu, &cpu.Reg.A, &cpu.Reg.A) }
	OpCodes[0x80] = /* add a,b */ func(cpu *CPU) int { return ADD_A(cpu, &cpu.Reg.B) }
	OpCodes[0x81] = /* add a,c */ func(cpu *CPU) int { return ADD_A(cpu, &cpu.Reg.C) }
	OpCodes[0x82] = /* add a,d */ func(cpu *CPU) int { return ADD_A(cpu, &cpu.Reg.D) }
	OpCodes[0x83] = /* add a,e */ func(cpu *CPU) int { return ADD_A(cpu, &cpu.Reg.E) }
	OpCodes[0x84] = /* add a,h */ func(cpu *CPU) int { return ADD_A(cpu, &cpu.Reg.H) }
	OpCodes[0x85] = /* add a,l */ func(cpu *CPU) int { return ADD_A(cpu, &cpu.Reg.L) }
	OpCodes[0x86] = /* add a,(hl) */ func(cpu *CPU) int {
		Alu_ADD_A(cpu, MEM_HL(cpu))
		log.Trace(2, "ADD A, (HL)")
		return 7
	}
	OpCodes[0x87] = /* add a,a */ func(cpu *CPU) int { return ADD_A(cpu, &cpu.Reg.A) }
	OpCodes[0x88] = /* adc a,b */ func(cpu *CPU) int { return ADC_A(cpu, &cpu.Reg.B) }
	OpCodes[0x89] = /* adc a,c */ func(cpu *CPU) int { return ADC_A(cpu, &cpu.Reg.C) }
	OpCodes[0x8a] = /* adc a,d */ func(cpu *CPU) int { return ADC_A(cpu, &cpu.Reg.D) }
	OpCodes[0x8b] = /* adc a,e */ func(cpu *CPU) int { return ADC_A(cpu, &cpu.Reg.E) }
	OpCodes[0x8c] = /* adc a,h */ func(cpu *CPU) int { return ADC_A(cpu, &cpu.Reg.H) }
	OpCodes[0x8d] = /* adc a,l */ func(cpu *CPU) int { return ADC_A(cpu, &cpu.Reg.L) }
	OpCodes[0x8e] = /* adc a,(hl) */ func(cpu *CPU) int {
		Alu_ADC_A(cpu, MEM_HL(cpu))
		log.Trace(2, "ADD A, (HL)")
		return 7
	}
	OpCodes[0x8f] = /* adc a,a */ func(cpu *CPU) int { return ADC_A(cpu, &cpu.Reg.A) }
	OpCodes[0x90] = /* sub b */ func(cpu *CPU) int { return SUB_A(cpu, &cpu.Reg.B) }
	OpCodes[0x91] = /* sub c */ func(cpu *CPU) int { return SUB_A(cpu, &cpu.Reg.C) }
	OpCodes[0x92] = /* sub d */ func(cpu *CPU) int { return SUB_A(cpu, &cpu.Reg.D) }
	OpCodes[0x93] = /* sub e */ func(cpu *CPU) int { return SUB_A(cpu, &cpu.Reg.E) }
	OpCodes[0x94] = /* sub h */ func(cpu *CPU) int { return SUB_A(cpu, &cpu.Reg.H) }
	OpCodes[0x95] = /* sub l */ func(cpu *CPU) int { return SUB_A(cpu, &cpu.Reg.L) }
	OpCodes[0x96] = /* sub (hl) */ func(cpu *CPU) int {
		Alu_SUB_A(cpu, MEM_HL(cpu))
		log.Trace(2, "SUB A, (HL)")
		return 7
	}
	OpCodes[0x97] = /* sub a */ func(cpu *CPU) int { return SUB_A(cpu, &cpu.Reg.A) }
	OpCodes[0x98] = /* sbc b */ func(cpu *CPU) int { return SBC_A(cpu, &cpu.Reg.B) }
	OpCodes[0x99] = /* sbc c */ func(cpu *CPU) int { return SBC_A(cpu, &cpu.Reg.C) }
	OpCodes[0x9a] = /* sbc d */ func(cpu *CPU) int { return SBC_A(cpu, &cpu.Reg.D) }
	OpCodes[0x9b] = /* sbc e */ func(cpu *CPU) int { return SBC_A(cpu, &cpu.Reg.E) }
	OpCodes[0x9c] = /* sbc h */ func(cpu *CPU) int { return SBC_A(cpu, &cpu.Reg.H) }
	OpCodes[0x9d] = /* sbc l */ func(cpu *CPU) int { return SBC_A(cpu, &cpu.Reg.L) }
	OpCodes[0x9e] = /* sbc (hl) */ func(cpu *CPU) int {
		Alu_SBC_A(cpu, MEM_HL(cpu))
		log.Trace(2, "SBC A, (HL)")
		return 7
	}
	OpCodes[0x9f] = /* sbc a */ func(cpu *CPU) int { return SBC_A(cpu, &cpu.Reg.A) }
	OpCodes[0xa0] = /* and b */ func(cpu *CPU) int { return AND_A(cpu, &cpu.Reg.B) }
	OpCodes[0xa1] = /* and c */ func(cpu *CPU) int { return AND_A(cpu, &cpu.Reg.C) }
	OpCodes[0xa2] = /* and d */ func(cpu *CPU) int { return AND_A(cpu, &cpu.Reg.D) }
	OpCodes[0xa3] = /* and e */ func(cpu *CPU) int { return AND_A(cpu, &cpu.Reg.E) }
	OpCodes[0xa4] = /* and h */ func(cpu *CPU) int { return AND_A(cpu, &cpu.Reg.H) }
	OpCodes[0xa5] = /* and l */ func(cpu *CPU) int { return AND_A(cpu, &cpu.Reg.L) }
	OpCodes[0xa6] = /* and (hl) */ func(cpu *CPU) int {
		Alu_AND_A(cpu, MEM_HL(cpu))
		log.Trace(2, "AND A, (HL)")
		return 7
	}
	OpCodes[0xa7] = /* and a */ func(cpu *CPU) int { return AND_A(cpu, &cpu.Reg.A) }
	OpCodes[0xa8] = /* xor b */ func(cpu *CPU) int { return XOR_A(cpu, &cpu.Reg.B) }
	OpCodes[0xa9] = /* xor c */ func(cpu *CPU) int { return XOR_A(cpu, &cpu.Reg.C) }
	OpCodes[0xaa] = /* xor d */ func(cpu *CPU) int { return XOR_A(cpu, &cpu.Reg.D) }
	OpCodes[0xab] = /* xor e */ func(cpu *CPU) int { return XOR_A(cpu, &cpu.Reg.E) }
	OpCodes[0xac] = /* xor h */ func(cpu *CPU) int { return XOR_A(cpu, &cpu.Reg.H) }
	OpCodes[0xad] = /* xor l */ func(cpu *CPU) int { return XOR_A(cpu, &cpu.Reg.L) }
	OpCodes[0xae] = /* xor (hl) */ func(cpu *CPU) int {
		Alu_XOR_A(cpu, MEM_HL(cpu))
		log.Trace(2, "XOR A, (HL)")
		return 7
	}
	OpCodes[0xaf] = /* xor a */ func(cpu *CPU) int { return XOR_A(cpu, &cpu.Reg.A) }
	OpCodes[0xb0] = /* or b */ func(cpu *CPU) int { return OR_A(cpu, &cpu.Reg.B) }
	OpCodes[0xb1] = /* or c */ func(cpu *CPU) int { return OR_A(cpu, &cpu.Reg.C) }
	OpCodes[0xb2] = /* or d */ func(cpu *CPU) int { return OR_A(cpu, &cpu.Reg.D) }
	OpCodes[0xb3] = /* or e */ func(cpu *CPU) int { return OR_A(cpu, &cpu.Reg.E) }
	OpCodes[0xb4] = /* or h */ func(cpu *CPU) int { return OR_A(cpu, &cpu.Reg.H) }
	OpCodes[0xb5] = /* or l */ func(cpu *CPU) int { return OR_A(cpu, &cpu.Reg.L) }
	OpCodes[0xb6] = /* or (hl) */ func(cpu *CPU) int {
		Alu_OR_A(cpu, MEM_HL(cpu))
		log.Trace(2, "OR A, (HL)")
		return 7
	}
	OpCodes[0xb7] = /* or a */ func(cpu *CPU) int { return OR_A(cpu, &cpu.Reg.A) }
	OpCodes[0xb8] = /* cp b */ func(cpu *CPU) int { return CP_A(cpu, &cpu.Reg.B) }
	OpCodes[0xb9] = /* cp c */ func(cpu *CPU) int { return CP_A(cpu, &cpu.Reg.C) }
	OpCodes[0xba] = /* cp d */ func(cpu *CPU) int { return CP_A(cpu, &cpu.Reg.D) }
	OpCodes[0xbb] = /* cp e */ func(cpu *CPU) int { return CP_A(cpu, &cpu.Reg.E) }
	OpCodes[0xbc] = /* cp h */ func(cpu *CPU) int { return CP_A(cpu, &cpu.Reg.H) }
	OpCodes[0xbd] = /* cp l */ func(cpu *CPU) int { return CP_A(cpu, &cpu.Reg.L) }
	OpCodes[0xbe] = /* cp (hl) */ func(cpu *CPU) int {
		Alu_CP_A(cpu, MEM_HL(cpu))
		log.Trace(2, "CP (HL)")
		return 7
	}
	OpCodes[0xbf] = /* cp a */ func(cpu *CPU) int { return CP_A(cpu, &cpu.Reg.A) }
	OpCodes[0xc0] = /* ret nz */ func(cpu *CPU) int { return RET_cc(cpu, FLAG_Z, false) }
	OpCodes[0xc1] = /* pop bc */ func(cpu *CPU) int { return POP(cpu, &cpu.Reg.B, &cpu.Reg.C) }
	OpCodes[0xc2] = /* jp nz,$+3 */ func(cpu *CPU) int {
		return JP_FLAG_nn(cpu, FLAG_ZERO, false)
	}
	OpCodes[0xc3] = /* jp $+3 */ JP_nn
	OpCodes[0xc4] = /* call nz,nn */ func(cpu *CPU) int { return CALL_FLAG_nn(cpu, FLAG_ZERO, false) }
	OpCodes[0xc5] = /* push bc */ func(cpu *CPU) int { return PUSH(cpu, &cpu.Reg.B, &cpu.Reg.C) }
	OpCodes[0xc6] = /* n add a,n */ func(cpu *CPU) int { return ADD_A_n(cpu, FetchOperand8(cpu)) }
	OpCodes[0xc7] = /* rst 0 */ func(cpu *CPU) int { return RST(cpu, 0) }
	OpCodes[0xc8] = /* ret z */ func(cpu *CPU) int { return RET_cc(cpu, FLAG_Z, true) }
	OpCodes[0xc9] = /* ret */ RET
	OpCodes[0xca] = /* jp z,$+3 */ func(cpu *CPU) int {
		return JP_FLAG_nn(cpu, FLAG_ZERO, true)
	}

	OpCodes[0xcb] = Instr_0xCB

	OpCodes[0xcc] = /* call z,nn */ func(cpu *CPU) int { return CALL_FLAG_nn(cpu, FLAG_ZERO, true) }
	OpCodes[0xcd] = /* call nn */ CALL_nn
	OpCodes[0xce] = /* adc a,n */ func(cpu *CPU) int {
		n := FetchOperand8(cpu)
		Alu_ADC_A(cpu, n)

		log.Trace(2, "ADD A, $%02x", n)
		return 7
	}
	OpCodes[0xcf] = /* rst 8h */ func(cpu *CPU) int { return RST(cpu, 0x8) }
	OpCodes[0xd0] = /* ret nc */ func(cpu *CPU) int { return RET_cc(cpu, FLAG_CARRY, false) }
	OpCodes[0xd1] = /* pop de */ func(cpu *CPU) int { return POP(cpu, &cpu.Reg.D, &cpu.Reg.E) }
	OpCodes[0xd2] = /* jp nc,$+3 */ func(cpu *CPU) int {
		return JP_FLAG_nn(cpu, FLAG_CARRY, false)
	}
	OpCodes[0xd3] = /* out (n),a */ OUT_n_A
	OpCodes[0xd4] = /* call nc,nn */ func(cpu *CPU) int { return CALL_FLAG_nn(cpu, FLAG_CARRY, false) }
	OpCodes[0xd5] = /* push de */ func(cpu *CPU) int { return PUSH(cpu, &cpu.Reg.D, &cpu.Reg.E) }
	OpCodes[0xd6] = /* sub n */ func(cpu *CPU) int {
		n := FetchOperand8(cpu)
		Alu_SUB_A(cpu, n)

		log.Trace(2, "SUB A, $%02x", n)
		return 7
	}
	OpCodes[0xd7] = /* rst 10h */ func(cpu *CPU) int { return RST(cpu, 0x10) }
	OpCodes[0xd8] = /* ret c */ func(cpu *CPU) int { return RET_cc(cpu, FLAG_CARRY, true) }
	OpCodes[0xd9] = /* exx */ EXX
	OpCodes[0xda] = /* jp c,$+3 */ func(cpu *CPU) int {
		return JP_FLAG_nn(cpu, FLAG_CARRY, true)
	}
	OpCodes[0xdb] = /* in a,(n) */ IN_A_n
	OpCodes[0xdc] = /* call c,nn */ func(cpu *CPU) int { return CALL_FLAG_nn(cpu, FLAG_CARRY, true) }

	OpCodes[0xdd] = func(cpu *CPU) int {
		op := FetchInstruction(cpu)
		return OpCodes_IX_IY[op](cpu, &cpu.Reg.IX)
	}

	OpCodes[0xde] = /* sbc a,n */ func(cpu *CPU) int {
		n := FetchOperand8(cpu)
		Alu_SBC_A(cpu, n)

		log.Trace(2, "SBC A, $%02x", n)
		return 7
	}
	OpCodes[0xdf] = /* rst 18h */ func(cpu *CPU) int { return RST(cpu, 0x18) }
	OpCodes[0xe0] = /* ret po */ func(cpu *CPU) int { return RET_cc(cpu, FLAG_PARTY_OVERFLOW, false) }
	OpCodes[0xe1] = /* pop hl */ func(cpu *CPU) int { return POP(cpu, &cpu.Reg.H, &cpu.Reg.L) }
	OpCodes[0xe2] = /* jp po,$+3 */ func(cpu *CPU) int {
		return JP_FLAG_nn(cpu, FLAG_PARTY_OVERFLOW, false)
	}
	OpCodes[0xe3] = /* ex (sp),hl */ EX_ISPI_HL
	OpCodes[0xe4] = /* call po,nn */ func(cpu *CPU) int { return CALL_FLAG_nn(cpu, FLAG_PARTY_OVERFLOW, false) }
	OpCodes[0xe5] = /* push hl */ func(cpu *CPU) int { return PUSH(cpu, &cpu.Reg.H, &cpu.Reg.L) }
	OpCodes[0xe6] = /* and n */ func(cpu *CPU) int {
		n := FetchOperand8(cpu)
		Alu_AND_A(cpu, n)

		log.Trace(2, "AND A, $%02x", n)
		return 7
	}
	OpCodes[0xe7] = /* rst 20h */ func(cpu *CPU) int { return RST(cpu, 0x20) }
	OpCodes[0xe8] = /* ret pe */ func(cpu *CPU) int { return RET_cc(cpu, FLAG_PARTY_OVERFLOW, true) }
	OpCodes[0xe9] = /* jp (hl) */ JP_IHLI
	OpCodes[0xea] = /* jp pe,$+3 */ func(cpu *CPU) int {
		return JP_FLAG_nn(cpu, FLAG_PARTY_OVERFLOW, true)
	}
	OpCodes[0xeb] = /* ex de,hl */ EX_DE_HL
	OpCodes[0xec] = /* call pe,nn */ func(cpu *CPU) int { return CALL_FLAG_nn(cpu, FLAG_PARTY_OVERFLOW, false) }

	OpCodes[0xed] = func(cpu *CPU) int {
		op := FetchInstruction(cpu)
		return OpCodes_ED[op](cpu)
	}

	OpCodes[0xee] = /* xor n 		*/ func(cpu *CPU) int {
		n := FetchOperand8(cpu)
		Alu_XOR_A(cpu, n)

		log.Trace(2, "XOR A, $%02x", n)
		return 7
	}

	OpCodes[0xef] = /* rst 28h 	*/ func(cpu *CPU) int { return RST(cpu, 0x28) }
	OpCodes[0xf0] = /* ret p 		*/ func(cpu *CPU) int { return RET_cc(cpu, FLAG_SIGN, false) }
	OpCodes[0xf1] = /* pop af 		*/ func(cpu *CPU) int { return POP(cpu, &cpu.Reg.A, &cpu.Reg.F) }
	OpCodes[0xf2] = /* jp p,$+3 	*/ func(cpu *CPU) int { return JP_FLAG_nn(cpu, FLAG_SIGN, false) }
	OpCodes[0xf3] = /* di 			*/ DI
	OpCodes[0xf4] = /* call p,nn 	*/ func(cpu *CPU) int { return CALL_FLAG_nn(cpu, FLAG_SIGN, false) }
	OpCodes[0xf5] = /* push af 	*/ func(cpu *CPU) int { return PUSH(cpu, &cpu.Reg.A, &cpu.Reg.F) }
	OpCodes[0xf6] = /* or n 		*/ func(cpu *CPU) int {
		n := FetchOperand8(cpu)
		Alu_OR_A(cpu, n)

		log.Trace(2, "OR A, $%02x", n)
		return 7
	}
	OpCodes[0xf7] = /* rst 30h 	*/ func(cpu *CPU) int { return RST(cpu, 0x30) }
	OpCodes[0xf8] = /* ret m 		*/ func(cpu *CPU) int { return RET_cc(cpu, FLAG_SIGN, true) }
	OpCodes[0xf9] = /* ld sp,hl 	*/ func(cpu *CPU) int {
		cpu.Reg.SP = cpu.Reg.HL()
		log.Trace(2, "LD SP, HL")
		return 6
	}
	OpCodes[0xfa] = /* jp m,$+3 	*/ func(cpu *CPU) int { return JP_FLAG_nn(cpu, FLAG_SIGN, true) }
	OpCodes[0xfb] = /* ei 			*/ EI
	OpCodes[0xfc] = /* call m,nn 	*/ func(cpu *CPU) int { return CALL_FLAG_nn(cpu, FLAG_SIGN, true) }

	OpCodes[0xfd] = func(cpu *CPU) int {
		op := FetchInstruction(cpu)
		return OpCodes_IX_IY[op](cpu, &cpu.Reg.IY)
	}

	OpCodes[0xfe] = /* cp n */ func(cpu *CPU) int {
		n := FetchOperand8(cpu)
		Alu_CP_A(cpu, n)

		log.Trace(2, "CP A, $%02x", n)
		return 7
	}
	OpCodes[0xff] = /* rst 38h */ func(cpu *CPU) int { return RST(cpu, 0x38) }

	OpCodes_ED[0x40] = /* in b,(c)     */ func(cpu *CPU) int { return IN_R_C(cpu, &cpu.Reg.B) }
	OpCodes_ED[0x41] = /* out (c),b    */ func(cpu *CPU) int { return OUT_C_R(cpu, &cpu.Reg.B) }
	OpCodes_ED[0x42] = /* sbc hl,bc    */ func(cpu *CPU) int { return SBC_HL_16(cpu, &cpu.Reg.B, &cpu.Reg.C) }
	OpCodes_ED[0x43] = /* ld (nn),bc   */ func(cpu *CPU) int { return LD_nn_RR_mem(cpu, &cpu.Reg.B, &cpu.Reg.C) }
	OpCodes_ED[0x44] = /* neg          */ func(cpu *CPU) int {
		Alu_NEG_A(cpu)
		log.Trace(2, "NEG")
		return 8
	}
	OpCodes_ED[0x45] = /* retn         */ RETN
	OpCodes_ED[0x46] = /* im 0         */ IM0
	OpCodes_ED[0x47] = /* ld i,a       */ LD_I_A
	OpCodes_ED[0x48] = /* in c,(c)     */ func(cpu *CPU) int { return IN_R_C(cpu, &cpu.Reg.C) }
	OpCodes_ED[0x49] = /* out (c),c    */ func(cpu *CPU) int { return OUT_C_R(cpu, &cpu.Reg.C) }
	OpCodes_ED[0x4a] = /* adc hl,bc    */ func(cpu *CPU) int { return ADC_HL_16(cpu, &cpu.Reg.B, &cpu.Reg.C) }
	OpCodes_ED[0x4b] = /* ld bc,(nn)   */ func(cpu *CPU) int { return LD_RR_nn_mem(cpu, &cpu.Reg.B, &cpu.Reg.C) }
	OpCodes_ED[0x4d] = /* reti         */ RETI
	OpCodes_ED[0x4f] = /* ld r,a       */ LD_R_A
	OpCodes_ED[0x50] = /* in d,(c)     */ func(cpu *CPU) int { return IN_R_C(cpu, &cpu.Reg.D) }
	OpCodes_ED[0x51] = /* out (c),d    */ func(cpu *CPU) int { return OUT_C_R(cpu, &cpu.Reg.D) }
	OpCodes_ED[0x52] = /* sbc hl,de    */ func(cpu *CPU) int { return SBC_HL_16(cpu, &cpu.Reg.D, &cpu.Reg.E) }
	OpCodes_ED[0x53] = /* ld (nn),de   */ func(cpu *CPU) int { return LD_nn_RR_mem(cpu, &cpu.Reg.D, &cpu.Reg.E) }
	OpCodes_ED[0x56] = /* im 1         */ IM1
	OpCodes_ED[0x57] = /* ld a,i       */ LD_A_I
	OpCodes_ED[0x58] = /* in e,(c)     */ func(cpu *CPU) int { return IN_R_C(cpu, &cpu.Reg.E) }
	OpCodes_ED[0x59] = /* out (c),e    */ func(cpu *CPU) int { return OUT_C_R(cpu, &cpu.Reg.E) }
	OpCodes_ED[0x5a] = /* adc hl,de    */ func(cpu *CPU) int { return ADC_HL_16(cpu, &cpu.Reg.D, &cpu.Reg.E) }
	OpCodes_ED[0x5b] = /* ld de,(nn)   */ func(cpu *CPU) int { return LD_RR_nn_mem(cpu, &cpu.Reg.D, &cpu.Reg.E) }
	OpCodes_ED[0x5e] = /* im 2         */ IM2
	OpCodes_ED[0x5f] = /* ld a,r       */ LD_A_R
	OpCodes_ED[0x60] = /* in h,(c)     */ func(cpu *CPU) int { return IN_R_C(cpu, &cpu.Reg.H) }
	OpCodes_ED[0x61] = /* out (c),h    */ func(cpu *CPU) int { return OUT_C_R(cpu, &cpu.Reg.H) }
	OpCodes_ED[0x62] = /* sbc hl,hl    */ func(cpu *CPU) int { return SBC_HL_16(cpu, &cpu.Reg.H, &cpu.Reg.L) }
	OpCodes_ED[0x67] = /* rrd          */ RRD
	OpCodes_ED[0x68] = /* in l,(c)     */ func(cpu *CPU) int { return IN_R_C(cpu, &cpu.Reg.L) }
	OpCodes_ED[0x69] = /* out (c),l    */ func(cpu *CPU) int { return OUT_C_R(cpu, &cpu.Reg.L) }
	OpCodes_ED[0x6a] = /* adc hl,hl    */ func(cpu *CPU) int { return ADC_HL_16(cpu, &cpu.Reg.H, &cpu.Reg.L) }
	OpCodes_ED[0x6f] = /* rld          */ RLD
	OpCodes_ED[0x72] = /* sbc hl,sp    */ func(cpu *CPU) int {
		cpu.Reg.HL_write(Alu_SBC16(cpu, cpu.Reg.HL(), cpu.Reg.SP))
		log.Trace(2, "SBC HL, SP")
		return 11
	}
	OpCodes_ED[0x73] = /* ld (nn),sp   */ func(cpu *CPU) int {
		nn := FetchOperand16(cpu)
		MemoryWrite16(cpu, nn, cpu.Reg.SP)

		log.Trace(2, "LD ($%04x), SP", nn)
		return 20
	}
	OpCodes_ED[0x78] = /* in a,(c)     */ func(cpu *CPU) int { return IN_R_C(cpu, &cpu.Reg.A) }
	OpCodes_ED[0x79] = /* out (c),a    */ func(cpu *CPU) int { return OUT_C_R(cpu, &cpu.Reg.A) }
	OpCodes_ED[0x7a] = /* adc hl,sp    */ func(cpu *CPU) int {
		cpu.Reg.HL_write(Alu_ADC16(cpu, cpu.Reg.HL(), cpu.Reg.SP))
		log.Trace(2, "ADC HL, SP")
		return 11
	}
	OpCodes_ED[0x7b] = /* ld sp,(nn)   */ func(cpu *CPU) int {
		nn := FetchOperand16(cpu)
		cpu.Reg.SP = MemoryRead16(cpu, nn)

		log.Trace(2, "LD SP, ($%04x)", nn)
		return 20
	}
	OpCodes_ED[0xa0] = /* ldi          */ LDI
	OpCodes_ED[0xa1] = /* cpi          */ CPI
	OpCodes_ED[0xa2] = /* ini          */ INI
	OpCodes_ED[0xa3] = /* outi         */ OUTI
	OpCodes_ED[0xa8] = /* ldd          */ LDD
	OpCodes_ED[0xa9] = /* cpd          */ CPD
	OpCodes_ED[0xaa] = /* ind          */ IND
	OpCodes_ED[0xab] = /* outd         */ OUTD
	OpCodes_ED[0xb0] = /* ldir         */ LDIR
	OpCodes_ED[0xb1] = /* cpir         */ CPIR
	OpCodes_ED[0xb2] = /* inir         */ INIR
	OpCodes_ED[0xb3] = /* otir         */ OUTIR
	OpCodes_ED[0xb8] = /* lddr         */ LDDR
	OpCodes_ED[0xb9] = /* cpdr         */ CPDR
	OpCodes_ED[0xba] = /* indr         */ INDR
	OpCodes_ED[0xbb] = /* otdr         */ OTDR

	OpCodes_IX_IY[0x09] = /* add ix/iy,bc */ func(cpu *CPU, idx *uint16) int {
		*idx = Alu_ADD16(cpu, *idx, cpu.Reg.BC())
		log.Trace(2, "ADD %s, BC", cpu.Reg.Name16(idx))
		return 15
	}
	OpCodes_IX_IY[0x19] = /* add ix/iy,de      */ func(cpu *CPU, idx *uint16) int {
		*idx = Alu_ADD16(cpu, *idx, cpu.Reg.DE())
		log.Trace(2, "ADD %s, DE", cpu.Reg.Name16(idx))
		return 15
	}
	OpCodes_IX_IY[0x21] = /* ld ix/iy,nn       */ LD_IXIY_nn
	OpCodes_IX_IY[0x22] = /* ld (nn),ix/iy     */ LD_nn_IXIY_mem
	OpCodes_IX_IY[0x23] = /* inc ix/iy         */ INC_IX_IY
	OpCodes_IX_IY[0x29] = /* add ix/iy,ix/iy   */ func(cpu *CPU, idx *uint16) int {
		*idx = Alu_ADD16(cpu, *idx, *idx)
		log.Trace(2, "ADD %s, %s", cpu.Reg.Name16(idx), cpu.Reg.Name16(idx))
		return 15
	}
	OpCodes_IX_IY[0x2a] = /* ld ix/iy,(nn)     */ LD_IXIY_nn_mem
	OpCodes_IX_IY[0x2b] = /* dec ix/iy         */ DEC_IX_IY
	OpCodes_IX_IY[0x34] = /* inc (ix/iy+n)     */ INC_IXIYd
	OpCodes_IX_IY[0x35] = /* dec (ix/iy+n)     */ DEC_IXIYd
	OpCodes_IX_IY[0x36] = /* ld (ix/iy+n),n    */ LD_IXIYd_n
	OpCodes_IX_IY[0x39] = /* add ix/iy,sp      */ func(cpu *CPU, idx *uint16) int {
		*idx = Alu_ADD16(cpu, *idx, cpu.Reg.SP)
		log.Trace(2, "ADD %s, SP", cpu.Reg.Name16(idx))
		return 15
	}

	OpCodes_IX_IY[0x46] = /* ld b,(ix/iy+n)    */ func(cpu *CPU, idx *uint16) int { return LD_R_IXIYd(cpu, idx, &cpu.Reg.B) }
	OpCodes_IX_IY[0x4e] = /* ld c,(ix/iy+n)    */ func(cpu *CPU, idx *uint16) int { return LD_R_IXIYd(cpu, idx, &cpu.Reg.C) }
	OpCodes_IX_IY[0x56] = /* ld d,(ix/iy+n)    */ func(cpu *CPU, idx *uint16) int { return LD_R_IXIYd(cpu, idx, &cpu.Reg.D) }
	OpCodes_IX_IY[0x5e] = /* ld e,(ix/iy+n)    */ func(cpu *CPU, idx *uint16) int { return LD_R_IXIYd(cpu, idx, &cpu.Reg.E) }
	OpCodes_IX_IY[0x66] = /* ld h,(ix/iy+n)    */ func(cpu *CPU, idx *uint16) int { return LD_R_IXIYd(cpu, idx, &cpu.Reg.H) }
	OpCodes_IX_IY[0x6e] = /* ld l,(ix/iy+n)    */ func(cpu *CPU, idx *uint16) int { return LD_R_IXIYd(cpu, idx, &cpu.Reg.L) }

	OpCodes_IX_IY[0x70] = /* ld (ix/iy+n),b    */ func(cpu *CPU, idx *uint16) int { return LD_IXIYd_R(cpu, idx, &cpu.Reg.B) }
	OpCodes_IX_IY[0x71] = /* ld (ix/iy+n),c    */ func(cpu *CPU, idx *uint16) int { return LD_IXIYd_R(cpu, idx, &cpu.Reg.C) }
	OpCodes_IX_IY[0x72] = /* ld (ix/iy+n),d    */ func(cpu *CPU, idx *uint16) int { return LD_IXIYd_R(cpu, idx, &cpu.Reg.D) }
	OpCodes_IX_IY[0x73] = /* ld (ix/iy+n),e    */ func(cpu *CPU, idx *uint16) int { return LD_IXIYd_R(cpu, idx, &cpu.Reg.E) }
	OpCodes_IX_IY[0x74] = /* ld (ix/iy+n),h    */ func(cpu *CPU, idx *uint16) int { return LD_IXIYd_R(cpu, idx, &cpu.Reg.H) }
	OpCodes_IX_IY[0x75] = /* ld (ix/iy+n),l    */ func(cpu *CPU, idx *uint16) int { return LD_IXIYd_R(cpu, idx, &cpu.Reg.L) }

	OpCodes_IX_IY[0x77] = /* ld (ix/iy+n),a    */ func(cpu *CPU, idx *uint16) int { return LD_IXIYd_R(cpu, idx, &cpu.Reg.A) }
	OpCodes_IX_IY[0x7e] = /* ld a,(ix/iy+n)    */ func(cpu *CPU, idx *uint16) int { return LD_R_IXIYd(cpu, idx, &cpu.Reg.A) }

	OpCodes_IX_IY[0x86] = /* add a,(ix/iy+n)   */ func(cpu *CPU, idx *uint16) int {
		d := FetchOperand8Compl(cpu)
		Alu_ADD_A(cpu, MemoryRead(cpu, uint16(int(*idx)+d)))

		log.Trace(2, "ADD A, (%s+%d)", cpu.Reg.Name16(idx), d)
		return 19
	}
	OpCodes_IX_IY[0x8e] = /* adc a,(ix/iy+n)   */ func(cpu *CPU, idx *uint16) int {
		d := FetchOperand8Compl(cpu)
		Alu_ADC_A(cpu, MemoryRead(cpu, uint16(int(*idx)+d)))

		log.Trace(2, "ADC A, (%s+%d)", cpu.Reg.Name16(idx), d)
		return 19
	}
	OpCodes_IX_IY[0x96] = /* sub (ix/iy+n)     */ func(cpu *CPU, idx *uint16) int {
		d := FetchOperand8Compl(cpu)
		Alu_SUB_A(cpu, MemoryRead(cpu, uint16(int(*idx)+d)))

		log.Trace(2, "SUB A, (%s+%d)", cpu.Reg.Name16(idx), d)
		return 19
	}
	OpCodes_IX_IY[0x9e] = /* sbc a,(ix/iy+n)   */ func(cpu *CPU, idx *uint16) int {
		d := FetchOperand8Compl(cpu)
		Alu_SBC_A(cpu, MemoryRead(cpu, uint16(int(*idx)+d)))

		log.Trace(2, "SBC A, (%s+%d)", cpu.Reg.Name16(idx), d)
		return 19
	}
	OpCodes_IX_IY[0xa6] = /* and (ix/iy+n)     */ func(cpu *CPU, idx *uint16) int {
		d := FetchOperand8Compl(cpu)
		Alu_AND_A(cpu, MemoryRead(cpu, uint16(int(*idx)+d)))

		log.Trace(2, "AND A, (%s+%d)", cpu.Reg.Name16(idx), d)
		return 19
	}
	OpCodes_IX_IY[0xae] = /* xor (ix/iy+n)     */ func(cpu *CPU, idx *uint16) int {
		d := FetchOperand8Compl(cpu)
		Alu_XOR_A(cpu, MemoryRead(cpu, uint16(int(*idx)+d)))

		log.Trace(2, "XOR A, (%s+%d)", cpu.Reg.Name16(idx), d)
		return 19
	}
	OpCodes_IX_IY[0xb6] = /* or (ix/iy+n)      */ func(cpu *CPU, idx *uint16) int {
		d := FetchOperand8Compl(cpu)
		Alu_OR_A(cpu, MemoryRead(cpu, uint16(int(*idx)+d)))

		log.Trace(2, "OR A, (%s+%d)", cpu.Reg.Name16(idx), d)
		return 19
	}
	OpCodes_IX_IY[0xbe] = /* cp (ix/iy+n)      */ func(cpu *CPU, idx *uint16) int {
		d := FetchOperand8Compl(cpu)
		Alu_CP_A(cpu, MemoryRead(cpu, uint16(int(*idx)+d)))

		log.Trace(2, "CP A, (%s+%d)", cpu.Reg.Name16(idx), d)
		return 19
	}

	OpCodes_IX_IY[0xcb] = func(cpu *CPU, idx *uint16) int {
		d := FetchOperand8Compl(cpu)
		op := FetchOperand8(cpu)
		return OpCodes_IX_IY_cb[op](cpu, idx, d)
	}

	OpCodes_IX_IY[0xe1] = /* pop ix/iy         */ POP_IX_IY
	OpCodes_IX_IY[0xe3] = /* ex (sp),ix/iy     */ EX_SP_IX_IY
	OpCodes_IX_IY[0xe5] = /* push ix/iy        */ PUSH_IX_IY
	OpCodes_IX_IY[0xe9] = /* jp (ix/iy)        */ JP_IX_IY
	OpCodes_IX_IY[0xf9] = /* ld sp,ix/iy       */ LD_SP_IX_IY

	// Undocumented instructions:
	OpCodes_IX_IY[0x24] = /* INC  IXH      */ func(cpu *CPU, reg *uint16) int { return UN_INC_HL(cpu, reg, false) }
	OpCodes_IX_IY[0x25] = /* DEC  IXH      */ func(cpu *CPU, reg *uint16) int { return UN_DEC_HL(cpu, reg, false) }
	OpCodes_IX_IY[0x26] = /* LD   IXH,nn   */ func(cpu *CPU, reg *uint16) int {
		n := FetchOperand8(cpu)
		l, _ := helpers.To8(*reg)
		*reg = helpers.To16(l, n)
		log.Trace(2, "LD %sh, %02x", cpu.Reg.Name16(reg), n)
		return 4 // TODO: CHECK
	}
	OpCodes_IX_IY[0x2C] = /* INC  IXL      */ func(cpu *CPU, reg *uint16) int { return UN_INC_HL(cpu, reg, true) }
	OpCodes_IX_IY[0x2D] = /* DEC  IXL      */ func(cpu *CPU, reg *uint16) int { return UN_DEC_HL(cpu, reg, true) }
	OpCodes_IX_IY[0x2E] = /* LD   IXL,nn   */ func(cpu *CPU, reg *uint16) int {
		n := FetchOperand8(cpu)
		_, h := helpers.To8(*reg)
		*reg = helpers.To16(n, h)
		log.Trace(2, "LD %sl, %02x", cpu.Reg.Name16(reg), n)
		return 4 // TODO: CHECK
	}
	OpCodes_IX_IY[0x44] = /* LD   B,IXH    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_R_HL(cpu, &cpu.Reg.B, reg, false) }
	OpCodes_IX_IY[0x45] = /* LD   B,IXL    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_R_HL(cpu, &cpu.Reg.B, reg, true) }
	OpCodes_IX_IY[0x4C] = /* LD   C,IXH    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_R_HL(cpu, &cpu.Reg.C, reg, false) }
	OpCodes_IX_IY[0x4D] = /* LD   C,IXL    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_R_HL(cpu, &cpu.Reg.C, reg, true) }
	OpCodes_IX_IY[0x54] = /* LD   D,IXH    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_R_HL(cpu, &cpu.Reg.D, reg, false) }
	OpCodes_IX_IY[0x55] = /* LD   D,IXL    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_R_HL(cpu, &cpu.Reg.D, reg, true) }
	OpCodes_IX_IY[0x5C] = /* LD   E,IXH    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_R_HL(cpu, &cpu.Reg.E, reg, false) }
	OpCodes_IX_IY[0x5D] = /* LD   E,IXL    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_R_HL(cpu, &cpu.Reg.E, reg, true) }
	OpCodes_IX_IY[0x60] = /* LD   IXH,B    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_HL_R(cpu, reg, &cpu.Reg.B, false) }
	OpCodes_IX_IY[0x61] = /* LD   IXH,C    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_HL_R(cpu, reg, &cpu.Reg.C, false) }
	OpCodes_IX_IY[0x62] = /* LD   IXH,D    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_HL_R(cpu, reg, &cpu.Reg.D, false) }
	OpCodes_IX_IY[0x63] = /* LD   IXH,E    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_HL_R(cpu, reg, &cpu.Reg.E, false) }
	OpCodes_IX_IY[0x64] = /* LD   IXH,IXH  */ func(cpu *CPU, reg *uint16) int {
		// TODO:
		l, h := helpers.To8(*reg)
		*reg = helpers.To16(l, h)
		log.Trace(2, "LD %sh, %sl", cpu.Reg.Name16(reg), cpu.Reg.Name16(reg))
		return 4 // TODO: CHECK
	}
	OpCodes_IX_IY[0x65] = /* LD   IXH,IXL  */ func(cpu *CPU, reg *uint16) int {
		l, _ := helpers.To8(*reg)
		*reg = helpers.To16(l, l)
		log.Trace(2, "LD %sh, %sl", cpu.Reg.Name16(reg), cpu.Reg.Name16(reg))
		return 4 // TODO: CHECK
	}
	OpCodes_IX_IY[0x67] = /* LD   IXH,A    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_HL_R(cpu, reg, &cpu.Reg.A, false) }
	OpCodes_IX_IY[0x68] = /* LD   IXL,B    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_HL_R(cpu, reg, &cpu.Reg.B, true) }
	OpCodes_IX_IY[0x69] = /* LD   IXL,C    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_HL_R(cpu, reg, &cpu.Reg.C, true) }
	OpCodes_IX_IY[0x6A] = /* LD   IXL,D    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_HL_R(cpu, reg, &cpu.Reg.D, true) }
	OpCodes_IX_IY[0x6B] = /* LD   IXL,E    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_HL_R(cpu, reg, &cpu.Reg.E, true) }
	OpCodes_IX_IY[0x6C] = /* LD   IXL,IXH  */ func(cpu *CPU, reg *uint16) int {
		_, h := helpers.To8(*reg)
		*reg = helpers.To16(h, h)
		log.Trace(2, "LD %sl, %sh", cpu.Reg.Name16(reg), cpu.Reg.Name16(reg))
		return 4 // TODO: CHECK
	}
	OpCodes_IX_IY[0x6D] = /* LD   IXL,IXL  */ func(cpu *CPU, reg *uint16) int {
		// TODO:
		l, h := helpers.To8(*reg)
		*reg = helpers.To16(l, h)
		log.Trace(2, "LD %sl, %sh", cpu.Reg.Name16(reg), cpu.Reg.Name16(reg))
		return 4 // TODO: CHECK
	}
	OpCodes_IX_IY[0x6F] = /* LD   IXL,A    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_HL_R(cpu, reg, &cpu.Reg.A, true) }
	OpCodes_IX_IY[0x7C] = /* LD   A,IXH    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_R_HL(cpu, &cpu.Reg.A, reg, false) }
	OpCodes_IX_IY[0x7D] = /* LD   A,IXL    */ func(cpu *CPU, reg *uint16) int { return UN_LD_R_R_HL(cpu, &cpu.Reg.A, reg, true) }
	OpCodes_IX_IY[0x84] = /* ADD  A,IXH    */ instrOpIXIY_Undocumented
	OpCodes_IX_IY[0x85] = /* ADD  A,IXL    */ instrOpIXIY_Undocumented
	OpCodes_IX_IY[0x8C] = /* ADC  A,IXH    */ instrOpIXIY_Undocumented
	OpCodes_IX_IY[0x8D] = /* ADC  A,IXL    */ instrOpIXIY_Undocumented
	OpCodes_IX_IY[0x94] = /* SUB  IXH      */ instrOpIXIY_Undocumented
	OpCodes_IX_IY[0x95] = /* SUB  IXL      */ instrOpIXIY_Undocumented
	OpCodes_IX_IY[0x9C] = /* SBC  A,IXH    */ instrOpIXIY_Undocumented
	OpCodes_IX_IY[0x9D] = /* SBC  A,IXL    */ instrOpIXIY_Undocumented
	OpCodes_IX_IY[0xA4] = /* AND  IXH      */ instrOpIXIY_Undocumented
	OpCodes_IX_IY[0xA5] = /* AND  IXL      */ instrOpIXIY_Undocumented
	OpCodes_IX_IY[0xAC] = /* XOR  IXH      */ instrOpIXIY_Undocumented
	OpCodes_IX_IY[0xAD] = /* XOR  IXL      */ instrOpIXIY_Undocumented
	OpCodes_IX_IY[0xB4] = /* OR   IXH      */ instrOpIXIY_Undocumented
	OpCodes_IX_IY[0xB5] = /* OR   IXL      */ instrOpIXIY_Undocumented
	OpCodes_IX_IY[0xBC] = /* CP   IXH      */ instrOpIXIY_Undocumented
	OpCodes_IX_IY[0xBD] = /* CP   IXL      */ instrOpIXIY_Undocumented

	OpCodes_IX_IY_cb[0x06] = /* rlc (ix/iy+n)     */ func(cpu *CPU, idx *uint16, n int) int {
		addr := uint16(int(*idx) + n)
		MemoryWrite(cpu, addr, Alu_RLC(cpu, MemoryRead(cpu, addr)))
		return 23
	}
	OpCodes_IX_IY_cb[0x0e] = /* rrc (ix/iy+n)     */ func(cpu *CPU, idx *uint16, n int) int {
		addr := uint16(int(*idx) + n)
		MemoryWrite(cpu, addr, Alu_RRC(cpu, MemoryRead(cpu, addr)))
		return 23
	}
	OpCodes_IX_IY_cb[0x16] = /* rl (ix/iy+n)      */ func(cpu *CPU, idx *uint16, n int) int {
		addr := uint16(int(*idx) + n)
		MemoryWrite(cpu, addr, Alu_RL(cpu, MemoryRead(cpu, addr)))
		return 23
	}
	OpCodes_IX_IY_cb[0x1e] = /* rr (ix/iy+n)      */ func(cpu *CPU, idx *uint16, n int) int {
		addr := uint16(int(*idx) + n)
		MemoryWrite(cpu, addr, Alu_RR(cpu, MemoryRead(cpu, addr)))
		return 23
	}
	OpCodes_IX_IY_cb[0x26] = /* sla (ix/iy+n)     */ func(cpu *CPU, idx *uint16, n int) int {
		addr := uint16(int(*idx) + n)
		MemoryWrite(cpu, addr, Alu_SLA(cpu, MemoryRead(cpu, addr)))
		return 23
	}
	OpCodes_IX_IY_cb[0x2e] = /* sra (ix/iy+n)     */ func(cpu *CPU, idx *uint16, n int) int {
		addr := uint16(int(*idx) + n)
		MemoryWrite(cpu, addr, Alu_SRA(cpu, MemoryRead(cpu, addr)))
		return 23
	}
	OpCodes_IX_IY_cb[0x46] = /* bit 0,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return BIT_IX_IY_cb(cpu, idx, n, 0) }
	OpCodes_IX_IY_cb[0x4e] = /* bit 1,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return BIT_IX_IY_cb(cpu, idx, n, 1) }
	OpCodes_IX_IY_cb[0x56] = /* bit 2,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return BIT_IX_IY_cb(cpu, idx, n, 2) }
	OpCodes_IX_IY_cb[0x5e] = /* bit 3,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return BIT_IX_IY_cb(cpu, idx, n, 3) }
	OpCodes_IX_IY_cb[0x66] = /* bit 4,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return BIT_IX_IY_cb(cpu, idx, n, 4) }
	OpCodes_IX_IY_cb[0x6e] = /* bit 5,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return BIT_IX_IY_cb(cpu, idx, n, 5) }
	OpCodes_IX_IY_cb[0x76] = /* bit 6,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return BIT_IX_IY_cb(cpu, idx, n, 6) }
	OpCodes_IX_IY_cb[0x7e] = /* bit 7,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return BIT_IX_IY_cb(cpu, idx, n, 7) }
	OpCodes_IX_IY_cb[0x86] = /* res 0,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return RES_IX_IY_cb(cpu, idx, n, 0) }
	OpCodes_IX_IY_cb[0x8e] = /* res 1,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return RES_IX_IY_cb(cpu, idx, n, 1) }
	OpCodes_IX_IY_cb[0x96] = /* res 2,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return RES_IX_IY_cb(cpu, idx, n, 2) }
	OpCodes_IX_IY_cb[0x9e] = /* res 3,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return RES_IX_IY_cb(cpu, idx, n, 3) }
	OpCodes_IX_IY_cb[0xa6] = /* res 4,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return RES_IX_IY_cb(cpu, idx, n, 4) }
	OpCodes_IX_IY_cb[0xae] = /* res 5,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return RES_IX_IY_cb(cpu, idx, n, 5) }
	OpCodes_IX_IY_cb[0xb6] = /* res 6,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return RES_IX_IY_cb(cpu, idx, n, 6) }
	OpCodes_IX_IY_cb[0xbe] = /* res 7,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return RES_IX_IY_cb(cpu, idx, n, 7) }
	OpCodes_IX_IY_cb[0xc6] = /* set 0,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return SET_IX_IY_cb(cpu, idx, n, 0) }
	OpCodes_IX_IY_cb[0xce] = /* set 1,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return SET_IX_IY_cb(cpu, idx, n, 1) }
	OpCodes_IX_IY_cb[0xd6] = /* set 2,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return SET_IX_IY_cb(cpu, idx, n, 2) }
	OpCodes_IX_IY_cb[0xde] = /* set 3,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return SET_IX_IY_cb(cpu, idx, n, 3) }
	OpCodes_IX_IY_cb[0xe6] = /* set 4,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return SET_IX_IY_cb(cpu, idx, n, 4) }
	OpCodes_IX_IY_cb[0xee] = /* set 5,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return SET_IX_IY_cb(cpu, idx, n, 5) }
	OpCodes_IX_IY_cb[0xf6] = /* set 6,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return SET_IX_IY_cb(cpu, idx, n, 6) }
	OpCodes_IX_IY_cb[0xfe] = /* set 7,(ix/iy+n)   */ func(cpu *CPU, idx *uint16, n int) int { return SET_IX_IY_cb(cpu, idx, n, 7) }
}
