package correlation

// GraphNode represents an entity in the security graph.
type GraphNode struct {
	ID   string `json:"id"`
	Type string `json:"type"` // e.g., "URL", "IP", "ASN"
}

// GraphEdge represents a relationship between entities.
type GraphEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Rel  string `json:"rel"` // e.g., "RESOLVES_TO", "REDIRECTS_TO"
}

// RelationshipGraph provides a structural view of the threat infrastructure.
type RelationshipGraph struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

type GraphBuilder struct{}

func NewGraphBuilder() *GraphBuilder {
	return &GraphBuilder{}
}

// BuildFromAnalysis constructs a graph from analysis data.
func (gb *GraphBuilder) BuildFromAnalysis(url string, resolvedIP string, asn string) *RelationshipGraph {
	graph := &RelationshipGraph{
		Nodes: []GraphNode{
			{ID: url, Type: "URL"},
			{ID: resolvedIP, Type: "IP"},
		},
		Edges: []GraphEdge{
			{From: url, To: resolvedIP, Rel: "RESOLVES_TO"},
		},
	}

	if asn != "" {
		graph.Nodes = append(graph.Nodes, GraphNode{ID: asn, Type: "ASN"})
		graph.Edges = append(graph.Edges, GraphEdge{From: resolvedIP, To: asn, Rel: "HOSTED_ON"})
	}

	return graph
}
