package day16

import (
	"container/heap"
	"fmt"

	"github.com/mikeramage/aoc2024/position"
	"github.com/mikeramage/aoc2024/priorityqueue"
	"github.com/mikeramage/aoc2024/utils"
)

const (
	N int = iota //Order indexes nicely into position.Directions
	E
	S
	W
)

const (
	rotCost  = 1000
	moveCost = 1
)

type GridContent byte

var (
	wall GridContent = '#'
	// empty GridContent = '.'
	start GridContent = 'S'
	end   GridContent = 'E'
)

type NodeState struct {
	Row, Col, Orientation int
}

type Node struct {
	State  NodeState
	Parent *Node
}

type DFSNode struct {
	State    NodeState
	Parent   *DFSNode
	PathCost int
}

func Day16() (int, int) {
	var part1, part2 int

	var grid [][]GridContent
	lines := utils.Lines("./input/day16.txt")
	var startR, startC, endR, endC int
	// rows := len(lines)
	// cols := len(lines[0])
	for r, line := range lines {
		var row []GridContent
		for c, content := range line {
			row = append(row, GridContent(content))
			if GridContent(content) == start {
				startR = r
				startC = c
			} else if GridContent(content) == end {
				endR = r
				endC = c
			}
		}
		grid = append(grid, row)
	}

	//Create initial node, explored set and frontier queue
	initialNodeState := NodeState{Row: startR, Col: startC, Orientation: E}
	initialNode := Node{State: initialNodeState, Parent: nil}
	explored := make(map[NodeState]int)
	frontier := priorityqueue.NewPriorityQueue[*Node]()
	initialQueueItem := priorityqueue.Item[*Node]{Value: &initialNode, Priority: 0}
	frontier.Push(&initialQueueItem)
	heap.Init(frontier)

	var solutions []*Node
	for frontier.Len() > 0 {
		queueItem := frontier.PopItem()
		node := queueItem.Value
		// fmt.Println(node.QueueItem.Priority, nodeState.Row, nodeState.Col, endR, endC)
		if node.State.Row == endR && node.State.Col == endC {
			part1 = queueItem.Priority
			solutions = append(solutions, node)
			break
		}
		explored[node.State] = queueItem.Priority
		pushNextNodes(frontier, queueItem, grid, explored, -1)
	}

	//Part2 - reinitialize
	initialNodeState = NodeState{Row: startR, Col: startC, Orientation: E}
	initialNode = Node{State: initialNodeState, Parent: nil}
	explored = make(map[NodeState]int)
	frontier = priorityqueue.NewPriorityQueue[*Node]()
	initialQueueItem = priorityqueue.Item[*Node]{Value: &initialNode, Priority: 0}
	frontier.Push(&initialQueueItem)
	heap.Init(frontier)

	var solutions2 []*Node
	for frontier.Len() > 0 {
		queueItem := frontier.PopItem()
		node := queueItem.Value
		// fmt.Println(node.QueueItem.Priority, nodeState.Row, nodeState.Col, endR, endC)
		if node.State.Row == endR && node.State.Col == endC {
			part1 = queueItem.Priority
			solutions2 = append(solutions, node)
		}
		explored[node.State] = queueItem.Priority
		pushNextNodes(frontier, queueItem, grid, explored, part1)
	}

	//Use a map so we're not counting parent nodes that are rotations of each other
	tilesOnPath := make(map[position.Position]bool)
	for _, solution := range solutions2 {
		pos := position.Position{Row: solution.State.Row, Col: solution.State.Col}
		tilesOnPath[pos] = true
		parent := solution.Parent
		for parent != nil {
			pos = position.Position{Row: parent.State.Row, Col: parent.State.Col}
			tilesOnPath[pos] = true
			parent = parent.Parent
		}
	}

	part2 = len(tilesOnPath)

	visualiseMaze(grid, solutions)
	visualiseMaze(grid, solutions2)
	return part1, part2
}

func pushNextNodes(frontier *priorityqueue.PriorityQueue[*Node], queueItem *priorityqueue.Item[*Node], grid [][]GridContent, explored map[NodeState]int, maxCost int) {
	node := queueItem.Value
	nodeState := node.State
	cost := queueItem.Priority

	orientations := []int{(nodeState.Orientation + 1) % 4, (nodeState.Orientation + 3) % 4} //Use +3 rather than -1 because mod negative numbers is bad :)
	var newQueueItems []*priorityqueue.Item[*Node]

	if maxCost < 0 || cost+rotCost <= maxCost { //maxCost < 0 ==> no constraint.
		for _, orientation := range orientations {
			newQueueItems = append(newQueueItems,
				createNewQueueItem(
					nodeState.Row,
					nodeState.Col,
					orientation,
					cost+rotCost,
					node))
		}
	}

	if maxCost < 0 || cost+moveCost <= maxCost {
		dr, dc := position.Directions[int(nodeState.Orientation)][0], position.Directions[int(nodeState.Orientation)][1]
		if grid[nodeState.Row+dr][nodeState.Col+dc] != wall { //No need to check boundary conditions because all states end up with wall
			//Only valid direction given current orientation
			newQueueItems = append(newQueueItems,
				createNewQueueItem(
					nodeState.Row+dr,
					nodeState.Col+dc,
					nodeState.Orientation,
					cost+moveCost,
					node))
		}
	}

	for _, newQueueItem := range newQueueItems {
		pathCost, seen := explored[newQueueItem.Value.State]
		if (maxCost < 0 && !seen) || (maxCost > 0 && (!seen || pathCost >= newQueueItem.Priority)) {
			frontier.PushItem(newQueueItem)
			heap.Fix(frontier, newQueueItem.Index)
		}
	}
}

func createNewQueueItem(row, col, orientation, priority int, parent *Node) *priorityqueue.Item[*Node] {
	return &priorityqueue.Item[*Node]{
		Value: &Node{
			State:  NodeState{Row: row, Col: col, Orientation: orientation},
			Parent: parent},
		Priority: priority}
}

func visualiseMaze(grid [][]GridContent, nodes []*Node) {
	var solutionPaths []map[position.Position]int

	for _, node := range nodes {
		pathMap := make(map[position.Position]int)
		parent := node.Parent
		for parent != nil {
			pathMap[position.Position{Row: parent.State.Row, Col: parent.State.Col}] = parent.State.Orientation
			parent = parent.Parent
		}
		solutionPaths = append(solutionPaths, pathMap)
	}

	for r, row := range grid {
		for c, content := range row {
			pos := position.Position{Row: r, Col: c}
			exists := false
			orientation := -1
			for _, pathMap := range solutionPaths {
				orientation, exists = pathMap[pos]
				if exists {
					break
				}
			}

			if exists {
				switch orientation {
				case N:
					fmt.Printf("^")
				case E:
					fmt.Printf(">")
				case S:
					fmt.Printf("v")
				case W:
					fmt.Printf("<")
				}
			} else {
				fmt.Printf("%v", string(content))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
