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
	QueueItem *priorityqueue.Item[NodeState]
	Parent    *Node
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
	initialNode := &Node{QueueItem: &priorityqueue.Item[NodeState]{Value: initialNodeState, Priority: 0}, Parent: nil}
	nodeMap := make(map[NodeState]*Node)
	nodeMap[initialNodeState] = initialNode
	explored := make(map[NodeState]bool)
	frontier := priorityqueue.NewPriorityQueue[NodeState]()
	frontier.Push(initialNode.QueueItem)
	heap.Init(frontier)

	var solutions []*Node
	for frontier.Len() > 0 {
		node, nodeState := popNode(frontier, nodeMap)
		// fmt.Println(node.QueueItem.Priority, nodeState.Row, nodeState.Col, endR, endC)
		if nodeState.Row == endR && nodeState.Col == endC {
			part1 = node.QueueItem.Priority
			solutions = append(solutions, node)
			break
		}
		explored[*nodeState] = true
		pushNextNodes(frontier, node, nodeState, grid, nodeMap, explored)
	}

	//Part 2 - just do a depth first search with limited explored set, but discard nodes if the cost is longer than the answer to part 1.
	initialDFSNode := DFSNode{State: NodeState{Row: startR, Col: startC, Orientation: E}, PathCost: 0, Parent: nil}
	dfsFrontier := []*DFSNode{&initialDFSNode}
	exploredDFS := make(map[NodeState]int)

	var solutions2 []*DFSNode
	for len(dfsFrontier) > 0 {
		dfsNode := dfsFrontier[len(dfsFrontier)-1]
		dfsFrontier = dfsFrontier[:len(dfsFrontier)-1] //do the pop
		exploredDFS[dfsNode.State] = dfsNode.PathCost
		if dfsNode.State.Row == endR && dfsNode.State.Col == endC {
			solutions2 = append(solutions2, dfsNode)
		}
		dfsFrontier = pushNextDfsNodes(dfsFrontier, dfsNode, grid, exploredDFS, part1)
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
	visualiseMaze2(grid, solutions2)
	return part1, part2
}

func pushNextDfsNodes(dfsFrontier []*DFSNode, dfsNode *DFSNode, grid [][]GridContent, exploredDFS map[NodeState]int, maxCost int) []*DFSNode {
	orientations := []int{(dfsNode.State.Orientation + 1) % 4, (dfsNode.State.Orientation + 3) % 4}
	var newNodes []*DFSNode

	if dfsNode.PathCost+rotCost <= maxCost {
		for _, orientation := range orientations {
			newNodes = append(newNodes,
				&DFSNode{
					State: NodeState{
						Row:         dfsNode.State.Row,
						Col:         dfsNode.State.Col,
						Orientation: orientation},
					Parent:   dfsNode,
					PathCost: dfsNode.PathCost + rotCost})
		}
	}

	if dfsNode.PathCost+moveCost <= maxCost {
		dr, dc := position.Directions[int(dfsNode.State.Orientation)][0], position.Directions[int(dfsNode.State.Orientation)][1]
		if grid[dfsNode.State.Row+dr][dfsNode.State.Col+dc] != wall { //No need to check boundary conditions because all states end up with wall
			//Only valid direction given current orientation
			newNodes = append(newNodes,
				&DFSNode{
					State: NodeState{
						Row:         dfsNode.State.Row + dr,
						Col:         dfsNode.State.Col + dc,
						Orientation: dfsNode.State.Orientation},
					Parent:   dfsNode,
					PathCost: dfsNode.PathCost + moveCost})
		}
	}

	for _, newNode := range newNodes {
		cost, exists := exploredDFS[newNode.State]
		if !exists || newNode.PathCost <= cost {
			dfsFrontier = append(dfsFrontier, newNode)

		}
	}
	return dfsFrontier
}

func pushNextNodes(frontier *priorityqueue.PriorityQueue[NodeState], node *Node, nodeState *NodeState, grid [][]GridContent, nodeMap map[NodeState]*Node, explored map[NodeState]bool) {
	orientations := []int{(nodeState.Orientation + 1) % 4, (nodeState.Orientation + 3) % 4} //Use +3 rather than -1 because mod negative numbers is bad :)
	var newStates []NodeState
	var costs []int
	pathCost := node.QueueItem.Priority
	for _, orientation := range orientations {
		newStates = append(newStates, NodeState{Row: nodeState.Row, Col: nodeState.Col, Orientation: orientation})
		costs = append(costs, rotCost)
	}

	dr, dc := position.Directions[int(nodeState.Orientation)][0], position.Directions[int(nodeState.Orientation)][1]
	if grid[nodeState.Row+dr][nodeState.Col+dc] != wall { //No need to check boundary conditions because all states end up with wall
		//Only valid direction given current orientation
		newStates = append(newStates, NodeState{Row: nodeState.Row + dr, Col: nodeState.Col + dc, Orientation: nodeState.Orientation})
		costs = append(costs, moveCost)
	}

	for i, state := range newStates {
		if !explored[state] {
			existingItem := frontier.GetItem(state)
			cost := pathCost + costs[i]
			if existingItem == nil {
				pushNode(frontier, nodeMap, state, cost, node)
			} else if existingItem.Priority > cost {
				frontier.Update(existingItem, state, cost)
			}
		}
	}
}

func pushNode(frontier *priorityqueue.PriorityQueue[NodeState], nodeMap map[NodeState]*Node, nodeState NodeState, priority int, parent *Node) {
	node := &Node{
		QueueItem: &priorityqueue.Item[NodeState]{Value: nodeState, Priority: priority},
		Parent:    parent,
	}
	nodeMap[nodeState] = node
	frontier.Push(node.QueueItem)
	heap.Fix(frontier, node.QueueItem.Index)
}

func popNode(frontier *priorityqueue.PriorityQueue[NodeState], nodeMap map[NodeState]*Node) (*Node, *NodeState) {
	item := frontier.PopItem()
	nodeState := item.Value
	node := nodeMap[nodeState]
	delete(nodeMap, nodeState)
	return node, &nodeState
}

func visualiseMaze(grid [][]GridContent, nodes []*Node) {
	var solutionPaths []map[position.Position]int

	for _, node := range nodes {
		pathMap := make(map[position.Position]int)
		parent := node.Parent
		for parent != nil {
			pathMap[position.Position{Row: parent.QueueItem.Value.Row, Col: parent.QueueItem.Value.Col}] = parent.QueueItem.Value.Orientation
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

// Duplication, but can't be bothered to fix
func visualiseMaze2(grid [][]GridContent, nodes []*DFSNode) {
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
