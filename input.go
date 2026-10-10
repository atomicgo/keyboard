package keyboard

import (
	"encoding/hex"
	"unicode/utf8"

	"atomicgo.dev/keyboard/internal"
	"atomicgo.dev/keyboard/keys"
)

// Sequence mappings.
var sequences = map[string]keys.Key{
	// Arrow keys
	"\x1b[A":     {Code: keys.Up},
	"\x1b[B":     {Code: keys.Down},
	"\x1b[C":     {Code: keys.Right},
	"\x1b[D":     {Code: keys.Left},
	"\x1b[1;2A":  {Code: keys.ShiftUp},
	"\x1b[1;2B":  {Code: keys.ShiftDown},
	"\x1b[1;2C":  {Code: keys.ShiftRight},
	"\x1b[1;2D":  {Code: keys.ShiftLeft},
	"\x1b[OA":    {Code: keys.ShiftUp},
	"\x1b[OB":    {Code: keys.ShiftDown},
	"\x1b[OC":    {Code: keys.ShiftRight},
	"\x1b[OD":    {Code: keys.ShiftLeft},
	"\x1b[a":     {Code: keys.ShiftUp},
	"\x1b[b":     {Code: keys.ShiftDown},
	"\x1b[c":     {Code: keys.ShiftRight},
	"\x1b[d":     {Code: keys.ShiftLeft},
	"\x1b[1;3A":  {Code: keys.Up, AltPressed: true},
	"\x1b[1;3B":  {Code: keys.Down, AltPressed: true},
	"\x1b[1;3C":  {Code: keys.Right, AltPressed: true},
	"\x1b[1;3D":  {Code: keys.Left, AltPressed: true},
	"\x1b\x1b[A": {Code: keys.Up, AltPressed: true},
	"\x1b\x1b[B": {Code: keys.Down, AltPressed: true},
	"\x1b\x1b[C": {Code: keys.Right, AltPressed: true},
	"\x1b\x1b[D": {Code: keys.Left, AltPressed: true},
	"\x1b[1;4A":  {Code: keys.ShiftUp, AltPressed: true},
	"\x1b[1;4B":  {Code: keys.ShiftDown, AltPressed: true},
	"\x1b[1;4C":  {Code: keys.ShiftRight, AltPressed: true},
	"\x1b[1;4D":  {Code: keys.ShiftLeft, AltPressed: true},
	"\x1b\x1b[a": {Code: keys.ShiftUp, AltPressed: true},
	"\x1b\x1b[b": {Code: keys.ShiftDown, AltPressed: true},
	"\x1b\x1b[c": {Code: keys.ShiftRight, AltPressed: true},
	"\x1b\x1b[d": {Code: keys.ShiftLeft, AltPressed: true},
	"\x1b[1;5A":  {Code: keys.CtrlUp},
	"\x1b[1;5B":  {Code: keys.CtrlDown},
	"\x1b[1;5C":  {Code: keys.CtrlRight},
	"\x1b[1;5D":  {Code: keys.CtrlLeft},
	"\x1b[Oa":    {Code: keys.CtrlUp, AltPressed: true},
	"\x1b[Ob":    {Code: keys.CtrlDown, AltPressed: true},
	"\x1b[Oc":    {Code: keys.CtrlRight, AltPressed: true},
	"\x1b[Od":    {Code: keys.CtrlLeft, AltPressed: true},
	"\x1b[1;6A":  {Code: keys.CtrlShiftUp},
	"\x1b[1;6B":  {Code: keys.CtrlShiftDown},
	"\x1b[1;6C":  {Code: keys.CtrlShiftRight},
	"\x1b[1;6D":  {Code: keys.CtrlShiftLeft},
	"\x1b[1;7A":  {Code: keys.CtrlUp, AltPressed: true},
	"\x1b[1;7B":  {Code: keys.CtrlDown, AltPressed: true},
	"\x1b[1;7C":  {Code: keys.CtrlRight, AltPressed: true},
	"\x1b[1;7D":  {Code: keys.CtrlLeft, AltPressed: true},
	"\x1b[1;8A":  {Code: keys.CtrlShiftUp, AltPressed: true},
	"\x1b[1;8B":  {Code: keys.CtrlShiftDown, AltPressed: true},
	"\x1b[1;8C":  {Code: keys.CtrlShiftRight, AltPressed: true},
	"\x1b[1;8D":  {Code: keys.CtrlShiftLeft, AltPressed: true},

	// Miscellaneous keys
	"\x1b[Z":      {Code: keys.ShiftTab},
	"\x1b[3~":     {Code: keys.Delete},
	"\x1b[3;3~":   {Code: keys.Delete, AltPressed: true},
	"\x1b[1~":     {Code: keys.Home},
	"\x1b[1;3H~":  {Code: keys.Home, AltPressed: true},
	"\x1b[4~":     {Code: keys.End},
	"\x1b[1;3F~":  {Code: keys.End, AltPressed: true},
	"\x1b[5~":     {Code: keys.PgUp},
	"\x1b[5;3~":   {Code: keys.PgUp, AltPressed: true},
	"\x1b[6~":     {Code: keys.PgDown},
	"\x1b[6;3~":   {Code: keys.PgDown, AltPressed: true},
	"\x1b[7~":     {Code: keys.Home},
	"\x1b[8~":     {Code: keys.End},
	"\x1b\x1b[3~": {Code: keys.Delete, AltPressed: true},
	"\x1b\x1b[5~": {Code: keys.PgUp, AltPressed: true},
	"\x1b\x1b[6~": {Code: keys.PgDown, AltPressed: true},
	"\x1b\x1b[7~": {Code: keys.Home, AltPressed: true},
	"\x1b\x1b[8~": {Code: keys.End, AltPressed: true},

	// Function keys
	"\x1bOP":     {Code: keys.F1},
	"\x1bOQ":     {Code: keys.F2},
	"\x1bOR":     {Code: keys.F3},
	"\x1bOS":     {Code: keys.F4},
	"\x1b[15~":   {Code: keys.F5},
	"\x1b[17~":   {Code: keys.F6},
	"\x1b[18~":   {Code: keys.F7},
	"\x1b[19~":   {Code: keys.F8},
	"\x1b[20~":   {Code: keys.F9},
	"\x1b[21~":   {Code: keys.F10},
	"\x1b[23~":   {Code: keys.F11},
	"\x1b[24~":   {Code: keys.F12},
	"\x1b[1;2P":  {Code: keys.F13},
	"\x1b[1;2Q":  {Code: keys.F14},
	"\x1b[1;2R":  {Code: keys.F15},
	"\x1b[1;2S":  {Code: keys.F16},
	"\x1b[15;2~": {Code: keys.F17},
	"\x1b[17;2~": {Code: keys.F18},
	"\x1b[18;2~": {Code: keys.F19},
	"\x1b[19;2~": {Code: keys.F20},

	// Function keys with the alt modifier
	"\x1b[1;3P":  {Code: keys.F1, AltPressed: true},
	"\x1b[1;3Q":  {Code: keys.F2, AltPressed: true},
	"\x1b[1;3R":  {Code: keys.F3, AltPressed: true},
	"\x1b[1;3S":  {Code: keys.F4, AltPressed: true},
	"\x1b[15;3~": {Code: keys.F5, AltPressed: true},
	"\x1b[17;3~": {Code: keys.F6, AltPressed: true},
	"\x1b[18;3~": {Code: keys.F7, AltPressed: true},
	"\x1b[19;3~": {Code: keys.F8, AltPressed: true},
	"\x1b[20;3~": {Code: keys.F9, AltPressed: true},
	"\x1b[21;3~": {Code: keys.F10, AltPressed: true},
	"\x1b[23;3~": {Code: keys.F11, AltPressed: true},
	"\x1b[24;3~": {Code: keys.F12, AltPressed: true},

	// Function keys, urxvt
	"\x1b[11~": {Code: keys.F1},
	"\x1b[12~": {Code: keys.F2},
	"\x1b[13~": {Code: keys.F3},
	"\x1b[14~": {Code: keys.F4},
	"\x1b[25~": {Code: keys.F13},
	"\x1b[26~": {Code: keys.F14},
	"\x1b[28~": {Code: keys.F15},
	"\x1b[29~": {Code: keys.F16},
	"\x1b[31~": {Code: keys.F17},
	"\x1b[32~": {Code: keys.F18},
	"\x1b[33~": {Code: keys.F19},
	"\x1b[34~": {Code: keys.F20},

	// Function keys with the alt modifier, urxvt
	"\x1b\x1b[11~": {Code: keys.F1, AltPressed: true},
	"\x1b\x1b[12~": {Code: keys.F2, AltPressed: true},
	"\x1b\x1b[13~": {Code: keys.F3, AltPressed: true},
	"\x1b\x1b[14~": {Code: keys.F4, AltPressed: true},
	"\x1b\x1b[25~": {Code: keys.F13, AltPressed: true},
	"\x1b\x1b[26~": {Code: keys.F14, AltPressed: true},
	"\x1b\x1b[28~": {Code: keys.F15, AltPressed: true},
	"\x1b\x1b[29~": {Code: keys.F16, AltPressed: true},
	"\x1b\x1b[31~": {Code: keys.F17, AltPressed: true},
	"\x1b\x1b[32~": {Code: keys.F18, AltPressed: true},
	"\x1b\x1b[33~": {Code: keys.F19, AltPressed: true},
	"\x1b\x1b[34~": {Code: keys.F20, AltPressed: true},
}

var hexCodes = map[string]keys.Key{
	"1b0d": {Code: keys.Enter, AltPressed: true},
	"1b7f": {Code: keys.Backspace, AltPressed: true},
	// support other backspace variants
	"1b08": {Code: keys.Backspace, AltPressed: true},
	"08":   {Code: keys.Backspace},

	// Powershell
	"1b4f41": {Code: keys.Up, AltPressed: false},
	"1b4f42": {Code: keys.Down, AltPressed: false},
	"1b4f43": {Code: keys.Right, AltPressed: false},
	"1b4f44": {Code: keys.Left, AltPressed: false},
	"1b5b46": {Code: keys.End, AltPressed: false},
	"1b5b48": {Code: keys.Home, AltPressed: false},
}

// escapeSequences holds every multi-byte input sequence that maps to a single key,
// keyed by its raw bytes. It merges sequences and hexCodes.
var escapeSequences, escapePrefixes, maxSequenceLen = buildSequenceTables()

func buildSequenceTables() (map[string]keys.Key, map[string]bool, int) {
	all := make(map[string]keys.Key, len(sequences)+len(hexCodes))
	for seq, key := range sequences {
		all[seq] = key
	}

	for code, key := range hexCodes {
		raw, err := hex.DecodeString(code)
		if err != nil {
			panic("keyboard: invalid hex code " + code)
		}

		all[string(raw)] = key
	}

	prefixes := make(map[string]bool)
	longest := 0

	for seq := range all {
		longest = max(longest, len(seq))

		for i := 1; i < len(seq); i++ {
			prefixes[seq[:i]] = true
		}
	}

	return all, prefixes, longest
}

// parseKeys splits raw input into key presses.
//
// Escape sequences and control characters always become keys of their own. If
// groupRunes is true, consecutive printable runes are combined into a single
// RuneKey, which is how pasted text arrives from a terminal. Otherwise every rune
// becomes its own key.
//
// If atEOF is false and the input ends in the middle of an escape sequence or a
// UTF-8 encoded rune, parsing stops before it. The returned count tells how many
// bytes were consumed; the rest has to be passed again together with more input.
// Bytes that are not valid UTF-8 are skipped.
func parseKeys(input []byte, atEOF, groupRunes bool) ([]keys.Key, int) {
	var (
		parsed []keys.Key
		group  []rune
		pos    int
	)

	flush := func() {
		switch {
		case len(group) == 0:
			return
		case len(group) == 1 && group[0] == ' ':
			parsed = append(parsed, keys.Key{Code: keys.Space, Runes: group})
		default:
			parsed = append(parsed, keys.Key{Code: keys.RuneKey, Runes: group})
		}

		group = nil
	}

	for pos < len(input) {
		rest := input[pos:]

		if key, n, ok := matchSequence(rest); ok {
			flush()

			parsed = append(parsed, key)
			pos += n

			continue
		}

		if !atEOF && (escapePrefixes[string(rest)] || !utf8.FullRune(rest)) {
			break
		}

		if rest[0] == internal.KeyEscape && len(rest) > 1 {
			// An escape followed by another key means the key was typed with alt held down.
			if r, size := utf8.DecodeRune(rest[1:]); r != utf8.RuneError {
				flush()

				key := runeKey(r)
				key.AltPressed = true
				parsed = append(parsed, key)
				pos += 1 + size

				continue
			}
		}

		r, size := utf8.DecodeRune(rest)
		pos += size

		if r == utf8.RuneError && size <= 1 {
			continue
		}

		key := runeKey(r)
		if groupRunes && (key.Code == keys.RuneKey || key.Code == keys.Space) {
			group = append(group, r)

			continue
		}

		flush()

		parsed = append(parsed, key)
	}

	flush()

	return parsed, pos
}

// matchSequence returns the key for the longest escape sequence at the start of input.
func matchSequence(input []byte) (keys.Key, int, bool) {
	for n := min(len(input), maxSequenceLen); n >= 1; n-- {
		if key, ok := escapeSequences[string(input[:n])]; ok {
			return key, n, true
		}
	}

	return keys.Key{}, 0, false
}

// runeKey returns the key for a single rune.
func runeKey(r rune) keys.Key {
	switch {
	case r == ' ':
		return keys.Key{Code: keys.Space, Runes: []rune{r}}
	case r <= internal.KeyUnitSeparator || r == internal.KeyDelete:
		return keys.Key{Code: keys.KeyCode(r)}
	default:
		return keys.Key{Code: keys.RuneKey, Runes: []rune{r}}
	}
}
