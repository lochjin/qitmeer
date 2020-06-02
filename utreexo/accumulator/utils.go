package accumulator

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Qitmeer/qitmeer/log"
	"sort"
)

// "verbose" is a global const to get lots of printfs for debugging
var verbose = false

// takes a slice of dels, removes the twins (in place) and returns a slice
// of parents of twins
func extractTwins(nodes []uint64, row uint8) (parents, dels []uint64) {
	for i := 0; i < len(nodes); i++ {
		if i+1 < len(nodes) && nodes[i]|1 == nodes[i+1] {
			parents = append(parents, parent(nodes[i], row))
			i++ // skip one here
		} else {
			dels = append(dels, nodes[i])
		}
	}
	return
}

// detectSubTreeHight finds the rows of the subtree a given LEAF position and
// the number of leaves (& the forest rows which is redundant)
// This thing is a tricky one.  Makes a weird serpinski fractal thing if
// you map it out in a table.
// Oh wait it's pretty simple.  Go left to right through the bits of numLeaves,
// and subtract that from position until it goes negative.
// (Does not work for nodes not at the bottom)
func detectSubTreeRows(
	position uint64, numLeaves uint64, forestRows uint8) (h uint8) {
	for h = forestRows; position >= (1<<h)&numLeaves; h-- {
		position -= (1 << h) & numLeaves
	}
	return
}

// TODO optimization if it's worth it --
// in many cases detectRow is called often and you're only looking for a
// change in row.  So we could instead have a "higher" function
// where it just looks for a different number of leading 0s.

// detectRow finds the current row of your node given the node
// position and the total forest rows.. counts preceding 1 bits.
func detectRow(position uint64, forestRows uint8) uint8 {
	marker := uint64(1 << forestRows)
	var h uint8
	for h = 0; position&marker != 0; h++ {
		marker >>= 1
	}
	return h
}

// detectOffset takes a node position and number of leaves in forest, and
// returns: which subtree a node is in, the L/R bitfield to descend to the node,
// and the height from node to its tree top (which is the bitfield length).
// we return the opposite of bits, because we always invert em...
func detectOffset(position uint64, numLeaves uint64) (uint8, uint8, uint64) {
	// TODO replace ?
	// similarities to detectSubTreeRows() with more features
	// maybe replace detectSubTreeRows with this

	// th = tree rows
	tr := treeRows(numLeaves)
	// nh = target node row
	nr := detectRow(position, tr) // there's probably a fancier way with bits...

	var biggerTrees uint8

	// add trees until you would exceed position of node

	// This is a bit of an ugly predicate.  The goal is to detect if we've
	// gone past the node we're looking for by inspecting progressively shorter
	// trees; once we have, the loop is over.

	// The predicate breaks down into 3 main terms:
	// A: pos << nh
	// B: mask
	// C: 1<<th & numleaves (treeSize)
	// The predicate is then if (A&B >= C)
	// A is position up-shifted by the row of the node we're targeting.
	// B is the "mask" we use in other functions; a bunch of 0s at the MSB side
	// and then a bunch of 1s on the LSB side, such that we can use bitwise AND
	// to discard high bits.  Together, A&B is shifting position up by nh bits,
	// and then discarding (zeroing out) the high bits.  This is the same as in
	// childMany.  C checks for whether a tree exists at the current tree
	// rows.  If there is no tree at th, C is 0.  If there is a tree, it will
	// return a power of 2: the base size of that tree.
	// The C term actually is used 3 times here, which is ugly; it's redefined
	// right on the next line.
	// In total, what this loop does is to take a node position, and
	// see if it's in the next largest tree.  If not, then subtract everything
	// covered by that tree from the position, and proceed to the next tree,
	// skipping trees that don't exist.

	for ; (position<<nr)&((2<<tr)-1) >= (1<<tr)&numLeaves; tr-- {
		treeSize := (1 << tr) & numLeaves
		if treeSize != 0 {
			position -= treeSize
			biggerTrees++
		}
	}
	return biggerTrees, tr - nr, ^position
}

// child gives you the left child (LSB will be 0)
func child(position uint64, forestRows uint8) uint64 {
	mask := uint64(2<<forestRows) - 1
	return (position << 1) & mask
}

// go down drop times (always left; LSBs will be 0) and return position
func childMany(position uint64, drop, forestRows uint8) uint64 {
	if drop == 0 {
		return position
	}
	if drop > forestRows {
		panic("childMany drop > forestRows")
	}
	mask := uint64(2<<forestRows) - 1
	return (position << drop) & mask
}

// Return the position of the parent of this position
func parent(position uint64, forestRows uint8) uint64 {
	return (position >> 1) | (1 << forestRows)
}

// go up rise times and return the position
func parentMany(position uint64, rise, forestRows uint8) uint64 {
	if rise == 0 {
		return position
	}
	if rise > forestRows {
		panic("parentMany rise > forestRows")
	}
	mask := uint64(2<<forestRows) - 1
	return (position>>rise | (mask << uint64(forestRows-(rise-1)))) & mask
}

// cousin returns a cousin: the child of the parent's sibling.
// you just xor with 2.  Actually there's no point in calling this function but
// it's here to document it.  If you're the left sibling it returns the left
// cousin.
func cousin(position uint64) uint64 {
	return position ^ 2
}

// TODO  inForest can probably be done better a different way.
// do we really need this at all?  only used for error detection in descendToPos

// check if a node is in a forest based on number of leaves.
// go down and right until reaching the bottom, then check if over numleaves
// (same as childmany)
// TODO fix.  says 14 is inforest with 5 leaves...
func inForest(pos, numLeaves uint64, forestRows uint8) bool {
	// quick yes:
	if pos < numLeaves {
		return true
	}
	marker := uint64(1 << forestRows)
	mask := (marker << 1) - 1
	if pos >= mask {
		return false
	}
	for pos&marker != 0 {
		pos = ((pos << 1) & mask) | 1
	}
	return pos < numLeaves
}

// given n leaves, how deep is the tree?
// iterate shifting left until greater than n
func treeRows(n uint64) (e uint8) {
	for ; (1 << e) < n; e++ {
	}
	return
}

// numRoots is just a popCount function, returning the number of 1 bits
func numRoots(n uint64) (r uint8) {
	for n > 0 {
		if n&1 == 1 {
			r++
		}
		n >>= 1
	}
	return
}

// rootPosition: given a number of leaves and a row, find the position of the
// root at that row.  Does not return an error if there's no root at that
// row so watch out and check first.  Checking is easy: leaves & (1<<h)
func rootPosition(leaves uint64, h, forestRows uint8) uint64 {
	mask := uint64(2<<forestRows) - 1
	before := leaves & (mask << (h + 1))
	shifted := (before >> h) | (mask << (forestRows - (h - 1)))
	return shifted & mask
}

// getroots gives you the positions of the tree roots, given a number of leaves.
// LOWEST first (right to left) (blarg change this)
func getRootsReverse(leaves uint64, forestRows uint8) (roots []uint64, rows []uint8) {
	position := uint64(0)

	// go left to right.  But append in reverse so that the roots are low to high
	// run though all bit positions.  if there's a 1, build a tree atop
	// the current position, and move to the right.
	for row := forestRows; position < leaves; row-- {
		if (1<<row)&leaves != 0 {
			// build a tree here
			root := parentMany(position, row, forestRows)
			roots = append([]uint64{root}, roots...)
			rows = append([]uint8{row}, rows...)
			position += 1 << row
		}
	}
	return
}

// subTreePositions takes in a node position and forestRows and returns the
// positions of all children that need to move AND THE NODE ITSELF.  (it works nicer that way)
// Also it returns where they should move to, given the destination of the
// sub-tree root.
// can also be used with the "to" return discarded to just enumerate a subtree
// swap tells whether to activate the sibling swap to try to preserve order
func subTreePositions(
	subroot uint64, moveTo uint64, forestRows uint8) (as []arrow) {

	subRow := detectRow(subroot, forestRows)
	//	fmt.Printf("node %d row %d\n", subroot, subRow)
	rootDelta := int64(moveTo) - int64(subroot)
	// do this with nested loops instead of recursion ... because that's
	// more fun.
	// r is out row in the forest
	// start at the bottom and ascend
	for r := uint8(0); r <= subRow; r++ {
		// find leftmost child on this row; also calculate the
		// delta (movement) for this row
		depth := subRow - r
		leftmost := childMany(subroot, depth, forestRows)
		rowDelta := rootDelta << depth // usually negative
		for i := uint64(0); i < 1<<depth; i++ {
			// loop left to right
			f := leftmost + i
			t := uint64(int64(f) + rowDelta)
			as = append(as, arrow{from: f, to: t})
		}
	}

	return
}

// TODO: unused? useless?
// subTreeLeafRange gives the range of leaves under a node
func subTreeLeafRange(
	subroot uint64, forestRows uint8) (uint64, uint64) {

	h := detectRow(subroot, forestRows)
	left := childMany(subroot, h, forestRows)
	run := uint64(1 << h)

	return left, run
}

// to leaves takes a arrow and returns a slice of arrows that are all the
// leaf arrows below it
func (a *arrow) toLeaves(h, forestRows uint8) []arrow {
	if h == 0 {
		return []arrow{*a}
	}

	run := uint64(1 << h)
	fromStart := childMany(a.from, h, forestRows)
	toStart := childMany(a.to, h, forestRows)

	leaves := make([]arrow, run)
	for i := uint64(0); i < run; i++ {
		leaves[i] = arrow{from: fromStart + i, to: toStart + i}
	}

	return leaves
}

// it'd be cool if you just had .sort() methods on slices of builtin types...
func sortUint64s(s []uint64) {
	sort.Slice(s, func(a, b int) bool { return s[a] < s[b] })
}

func sortNodeSlice(s []node) {
	sort.Slice(s, func(a, b int) bool { return s[a].Pos < s[b].Pos })
}

// checkSortedNoDupes returns true for strictly increasing slices
func checkSortedNoDupes(s []uint64) bool {
	for i, _ := range s {
		if i == 0 {
			continue
		}
		if s[i-1] >= s[i] {
			return false
		}
	}
	return true
}

// TODO is there really no way to just... reverse any slice?  Like with
// interface or something?  it's just pointers and never touches the actual
// type...

// reverseArrowSlice does what it says.  Maybe can get rid of if we return
// the slice top-down instead of bottom-up
func reverseArrowSlice(a []arrow) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

// exact same code twice, couldn't you have a reverse *any* slice func...?
// but maybe that's generics or something
func reverseUint64Slice(a []uint64) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func reversePolNodeSlice(a []polNode) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

// mergeSortedSlices takes two slices (of uint64s; though this seems
// genericizable in that it's just < and > operators) and merges them into
// a single sorted slice, discarding duplicates.  (but not detecting or discarding
// duplicates within a single slice)
// (eg [1, 5, 8, 9], [2, 3, 4, 5, 6] -> [1, 2, 3, 4, 5, 6, 8, 9]
func mergeSortedSlices(a []uint64, b []uint64) (c []uint64) {
	maxa := len(a)
	maxb := len(b)

	// shortcuts:
	if maxa == 0 {
		return b
	}
	if maxb == 0 {
		return a
	}

	// make it (potentially) too long and truncate later
	c = make([]uint64, maxa+maxb)

	idxa, idxb := 0, 0
	for j := 0; j < len(c); j++ {
		// if we're out of a or b, just use the remainder of the other one
		if idxa >= maxa {
			// a is done, copy remainder of b
			j += copy(c[j:], b[idxb:])
			c = c[:j] // truncate empty section of c
			break
		}
		if idxb >= maxb {
			// b is done, copy remainder of a
			j += copy(c[j:], a[idxa:])
			c = c[:j] // truncate empty section of c
			break
		}

		vala, valb := a[idxa], b[idxb]
		if vala < valb { // a is less so append that
			c[j] = vala
			idxa++
		} else if vala > valb { // b is less so append that
			c[j] = valb
			idxb++
		} else { // they're equal
			c[j] = vala
			idxa++
			idxb++
		}
	}
	return
}

// kindof like mergeSortedSlices.  Takes 2 sorted slices a, b and removes
// all elements of b from a and returns a.
// in this case b is arrow.to
func dedupeSwapDirt(a []uint64, b []arrow) []uint64 {
	maxa := len(a)
	maxb := len(b)
	var c []uint64
	// shortcuts:
	if maxa == 0 || maxb == 0 {
		return a
	}
	idxb := 0
	for j := 0; j < maxa; j++ {
		// skip over swap destinations less than current dirt
		for idxb < maxb && a[j] < b[idxb].to {
			idxb++
		}
		if idxb == maxb { // done
			c = append(c, a[j:]...)
			break
		}
		if a[j] != b[idxb].to {
			c = append(c, a[j])
		}
	}

	return c
}

// BinString prints out the whole thing.  Only viable for small forests
func BinString(leaves uint64) string {
	fh := treeRows(leaves)

	// tree rows should be 6 or less
	if fh > 6 {
		return "forest too big to print "
	}

	output := make([]string, (fh*2)+1)
	var pos uint8
	for h := uint8(0); h <= fh; h++ {
		rowlen := uint8(1 << (fh - h))

		for j := uint8(0); j < rowlen; j++ {
			//			if pos < uint8(leaves) {
			output[h*2] += fmt.Sprintf("%05b ", pos)
			//			} else {
			//				output[h*2] += fmt.Sprintf("       ")
			//			}

			if h > 0 {
				//				if x%2 == 0 {
				output[(h*2)-1] += "|-----"
				for q := uint8(0); q < ((1<<h)-1)/2; q++ {
					output[(h*2)-1] += "------"
				}
				output[(h*2)-1] += "\\     "
				for q := uint8(0); q < ((1<<h)-1)/2; q++ {
					output[(h*2)-1] += "      "
				}

				//				}

				for q := uint8(0); q < (1<<h)-1; q++ {
					output[h*2] += "      "
				}

			}
			pos++
		}

	}
	var s string
	for z := len(output) - 1; z >= 0; z-- {
		s += output[z] + "\n"
	}
	return s
}

// BtU32 : 4 byte slice to uint32.  Returns ffffffff if something doesn't work.
func BtU32(b []byte) uint32 {
	if len(b) != 4 {
		log.Error(fmt.Sprintf("Got %x to BtU32 (%d bytes)\n", b, len(b)))
		return 0xffffffff
	}
	var i uint32
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.BigEndian, &i)
	return i
}

// U32tB : uint32 to 4 bytes.  Always works.
func U32tB(i uint32) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, i)
	return buf.Bytes()
}

// BtU64 : 8 bytes to uint64.  returns ffff. if it doesn't work.
func BtU64(b []byte) uint64 {
	if len(b) != 8 {
		log.Error(fmt.Sprintf("Got %x to BtU64 (%d bytes)\n", b, len(b)))
		return 0xffffffffffffffff
	}
	var i uint64
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.BigEndian, &i)
	return i
}

// U64tB : uint64 to 8 bytes.  Always works.
func U64tB(i uint64) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, i)
	return buf.Bytes()
}

// BtU8 : 1 byte to uint8.  returns ffff. if it doesn't work.
func BtU8(b []byte) uint8 {
	if len(b) != 1 {
		log.Error(fmt.Sprintf("Got %x to BtU8 (%d bytes)\n", b, len(b)))
		return 0xff
	}
	var i uint8
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.BigEndian, &i)
	return i
}

// U8tB : uint8 to a byte.  Always works.
func U8tB(i uint8) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, i)
	return buf.Bytes()
}
