package services

import (
	"workflow/src/modules/definitions/domain/entities"
)

/*
 * CYCLE DETECTOR
 * Port of the frontend workflowCycleDetector.ts.
 * Detects tight cycles (no async pause point) in workflow graphs.
 * Called on CREATE/UPDATE only — never at runtime.
 *
 * Algorithm: iterative DFS, 3-color (white=0, gray=1, black=2).
 * Cycles through async nodes (wait_signal, wait_for, delay) are allowed.
 * Loop body back-edges (core/loop → sourceHandle "body") are excluded.
 * Annotation nodes/edges are excluded.
 */

// asyncNodeTypes are node types that represent async pause points.
// Cycles passing through these are allowed (they break execution).
var asyncNodeTypes = map[string]bool{
	"core/wait_signal": true,
	"core/wait_for":    true,
	"core/delay":       true,
}

// annotationTypes are visual-only nodes excluded from graph analysis.
var annotationTypes = map[string]bool{
	"core/text_note":   true,
	"core/group_frame": true,
}

// DetectTightCycles returns node IDs participating in cycles that have no async pause point.
// Cycles through wait_signal, wait_for, or delay are allowed (they break execution).
// Loop body back-edges (core/loop → body handle) are excluded.
// Annotation nodes/edges are excluded.
func DetectTightCycles(nodes []entities.WorkflowNode, edges []entities.WorkflowEdge) []string {
	// Build nodeTypeMap: nodeId → nodeType
	nodeTypeMap := make(map[string]string, len(nodes))
	for _, node := range nodes {
		nodeTypeMap[node.ID] = node.Type
	}

	// Build adjacency list (excluding annotations and loop body back-edges)
	adjacency := make(map[string][]string, len(nodes))
	for _, node := range nodes {
		if !annotationTypes[node.Type] {
			adjacency[node.ID] = nil
		}
	}

	for _, edge := range edges {
		// Skip annotation edges
		if isAnnotationEdge(edge) {
			continue
		}

		// Skip edges involving annotation nodes
		sourceType := nodeTypeMap[edge.Source]
		targetType := nodeTypeMap[edge.Target]
		if annotationTypes[sourceType] || annotationTypes[targetType] {
			continue
		}

		// Skip loop body back-edges
		if isLoopBodyEdge(edge, nodeTypeMap) {
			continue
		}

		adjacency[edge.Source] = append(adjacency[edge.Source], edge.Target)
	}

	// Find all cycles via iterative DFS
	cycles := findAllCycles(adjacency)

	// Filter tight cycles (no async node in cycle)
	tightNodes := make(map[string]struct{})
	for _, cycle := range cycles {
		hasAsync := false
		for _, nodeID := range cycle {
			if asyncNodeTypes[nodeTypeMap[nodeID]] {
				hasAsync = true
				break
			}
		}
		if !hasAsync {
			for _, nodeID := range cycle {
				tightNodes[nodeID] = struct{}{}
			}
		}
	}

	if len(tightNodes) == 0 {
		return nil
	}

	result := make([]string, 0, len(tightNodes))
	for nodeID := range tightNodes {
		result = append(result, nodeID)
	}
	return result
}

// isAnnotationEdge checks whether an edge connects annotation handles.
func isAnnotationEdge(edge entities.WorkflowEdge) bool {
	return edge.SourceHandle == "__note_out" || edge.TargetHandle == "__note"
}

// isLoopBodyEdge checks whether an edge is a loop-body back-edge.
// Loop nodes (core/loop) have a 'body' sourceHandle that feeds back
// into the loop body — this is an intentional structural cycle.
func isLoopBodyEdge(edge entities.WorkflowEdge, nodeTypeMap map[string]string) bool {
	return nodeTypeMap[edge.Source] == "core/loop" && edge.SourceHandle == "body"
}

// findAllCycles finds all elementary cycles using iterative DFS with 3-color marking.
// Returns cycles as slices of node IDs.
func findAllCycles(adjacency map[string][]string) [][]string {
	var cycles [][]string

	// 0 = white (unvisited), 1 = gray (in stack), 2 = black (done)
	color := make(map[string]int, len(adjacency))
	parent := make(map[string]string, len(adjacency))

	for nodeID := range adjacency {
		color[nodeID] = 0
	}

	// stackEntry holds the current node and the index of the next neighbor to visit
	type stackEntry struct {
		nodeID  string
		nbIndex int
	}

	for startNode := range adjacency {
		if color[startNode] != 0 {
			continue
		}

		stack := []stackEntry{{nodeID: startNode, nbIndex: 0}}
		color[startNode] = 1
		parent[startNode] = ""

		for len(stack) > 0 {
			top := &stack[len(stack)-1]
			neighbors := adjacency[top.nodeID]

			if top.nbIndex < len(neighbors) {
				neighbor := neighbors[top.nbIndex]
				top.nbIndex++

				neighborColor := color[neighbor]
				if neighborColor == 0 {
					color[neighbor] = 1
					parent[neighbor] = top.nodeID
					stack = append(stack, stackEntry{nodeID: neighbor, nbIndex: 0})
				} else if neighborColor == 1 {
					// Back-edge found → reconstruct cycle
					cycle := reconstructCycle(neighbor, top.nodeID, parent)
					cycles = append(cycles, cycle)
				}
			} else {
				color[top.nodeID] = 2
				stack = stack[:len(stack)-1]
			}
		}
	}

	return cycles
}

// reconstructCycle traces the parent chain from cycleEnd back to cycleStart.
func reconstructCycle(cycleStart, cycleEnd string, parent map[string]string) []string {
	path := []string{cycleEnd}
	current := cycleEnd

	for current != cycleStart {
		p := parent[current]
		if p == "" {
			break
		}
		path = append(path, p)
		current = p
	}

	// Reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}
