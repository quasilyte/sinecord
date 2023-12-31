package synthdb

type Note struct {
	Name string
	Code byte
}

const (
	MinNoteCode      = 24
	MaxNoteCode      = 119
	Octave1StartCode = 36
	Ocvate4EndCode   = 83
)

var NoteTab = [256]Note{
	24:  {Name: "C0", Code: 24},
	25:  {Name: "C#0", Code: 25},
	26:  {Name: "D0", Code: 26},
	27:  {Name: "D#0", Code: 27},
	28:  {Name: "E0", Code: 28},
	29:  {Name: "F0", Code: 29},
	30:  {Name: "F0", Code: 30},
	31:  {Name: "G0", Code: 31},
	32:  {Name: "G#0", Code: 32},
	33:  {Name: "A0", Code: 33},
	34:  {Name: "A#0", Code: 34},
	35:  {Name: "B0", Code: 35},
	36:  {Name: "C1", Code: 36},
	37:  {Name: "C#1", Code: 37},
	38:  {Name: "D1", Code: 38},
	39:  {Name: "D#1", Code: 39},
	40:  {Name: "E1", Code: 40},
	41:  {Name: "F1", Code: 41},
	42:  {Name: "F1", Code: 42},
	43:  {Name: "G1", Code: 43},
	44:  {Name: "G#1", Code: 44},
	45:  {Name: "A1", Code: 45},
	46:  {Name: "A#1", Code: 46},
	47:  {Name: "B1", Code: 47},
	48:  {Name: "C2", Code: 48},
	49:  {Name: "C#2", Code: 49},
	50:  {Name: "D2", Code: 50},
	51:  {Name: "D#2", Code: 51},
	52:  {Name: "E2", Code: 52},
	53:  {Name: "F2", Code: 53},
	54:  {Name: "F2", Code: 54},
	55:  {Name: "G2", Code: 55},
	56:  {Name: "G#2", Code: 56},
	57:  {Name: "A2", Code: 57},
	58:  {Name: "A#2", Code: 58},
	59:  {Name: "B2", Code: 59},
	60:  {Name: "C3", Code: 60},
	61:  {Name: "C#3", Code: 61},
	62:  {Name: "D3", Code: 62},
	63:  {Name: "D#3", Code: 63},
	64:  {Name: "E3", Code: 64},
	65:  {Name: "F3", Code: 65},
	66:  {Name: "F3", Code: 66},
	67:  {Name: "G3", Code: 67},
	68:  {Name: "G#3", Code: 68},
	69:  {Name: "A3", Code: 69},
	70:  {Name: "A#3", Code: 70},
	71:  {Name: "B3", Code: 71},
	72:  {Name: "C4", Code: 72},
	73:  {Name: "C#4", Code: 73},
	74:  {Name: "D4", Code: 74},
	75:  {Name: "D#4", Code: 75},
	76:  {Name: "E4", Code: 76},
	77:  {Name: "F4", Code: 77},
	78:  {Name: "F4", Code: 78},
	79:  {Name: "G4", Code: 79},
	80:  {Name: "G#4", Code: 80},
	81:  {Name: "A4", Code: 81},
	82:  {Name: "A#4", Code: 82},
	83:  {Name: "B4", Code: 83},
	84:  {Name: "C5", Code: 84},
	85:  {Name: "C#5", Code: 85},
	86:  {Name: "D5", Code: 86},
	87:  {Name: "D#5", Code: 87},
	88:  {Name: "E5", Code: 88},
	89:  {Name: "F5", Code: 89},
	90:  {Name: "F5", Code: 90},
	91:  {Name: "G5", Code: 91},
	92:  {Name: "G#5", Code: 92},
	93:  {Name: "A5", Code: 93},
	94:  {Name: "A#5", Code: 94},
	95:  {Name: "B5", Code: 95},
	96:  {Name: "C6", Code: 96},
	97:  {Name: "C#6", Code: 97},
	98:  {Name: "D6", Code: 98},
	99:  {Name: "D#6", Code: 99},
	100: {Name: "E6", Code: 100},
	101: {Name: "F6", Code: 101},
	102: {Name: "F6", Code: 102},
	103: {Name: "G6", Code: 103},
	104: {Name: "G#6", Code: 104},
	105: {Name: "A6", Code: 105},
	106: {Name: "A#6", Code: 106},
	107: {Name: "B6", Code: 107},
	108: {Name: "C7", Code: 108},
	109: {Name: "C#7", Code: 109},
	110: {Name: "D7", Code: 110},
	111: {Name: "D#7", Code: 111},
	112: {Name: "E7", Code: 112},
	113: {Name: "F7", Code: 113},
	114: {Name: "F7", Code: 114},
	115: {Name: "G7", Code: 115},
	116: {Name: "G#7", Code: 116},
	117: {Name: "A7", Code: 117},
	118: {Name: "A#7", Code: 118},
	119: {Name: "B7", Code: 119},
}
