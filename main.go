package main

/*
	This program tests hash and matrix graph classes.
	Author: Ivan Liljeqvist
	Date: 23-04-2015
*/

import (
	graph "./graph"
	"fmt"
	"math/rand"
	"time"
)

/*
	This interface describes the hash and matrix classes.
	It declares all the methods.

	These methods are matched in matrix.go and hash.go
*/

type Grapher interface {
	NumVertices() int
	NumEdges() int
	Degree(int) int
	DoNeighbors(int, func(int, interface{}))
	HasEdge(int, int) bool
	Label(int, int) interface{}
	Add(int, int)
	AddLabel(int, int, interface{})
	AddBi(int, int)
	AddBiLabel(int, int, interface{})
	Remove(int, int)
	RemoveBi(int, int)
}

/*
	This function need a Grapher interface instance.

	This function returns the size of the largest component in the graph and the number of components
	in that graph.
*/

func getLargestSizeAndNumOfComponents(g Grapher) (largest_component_size, number_of_components int) {

	//we'll be searching for largest component - set it to lowest possible number first
	largest_component_size = 0
	number_of_components = 0

	//we'll need to pass a state array of booleans to graph.DFS
	//make an array of booleans, one element for each vertex
	state := make([]bool, g.NumVertices())

	//we'll loop through all components - size of each will be saved here
	current_component_size := 0

	//loop through the state boolean array. We get index of each vertex as 'vertex' and
	//'visited' is the boolean if it's visited or not
	for vertex, visited := range state {
		//if not visited - a new component we have't traversed
		if visited == false {

			number_of_components++
			current_component_size = 0

			graph.DFS(g, vertex, state, func(w int) {
				//increase the size each 'traverse-loop'
				current_component_size++
				if current_component_size > largest_component_size {
					largest_component_size = current_component_size
				}

			})

		}

	}

	//now we have calculated largest_component size and number_of_components and are ready to return

	return

}

/*
	Takes in a boolean - either creates a hash or matrix.
    Takes in n - number of edges and verticies.
    Returns two Grapher object - one matrix and one hash.
    The vertices in the graphs returned are conencted exactly the same.
*/

func setupGraphs(n int) (hash_graph_to_return Grapher, matrix_graph_to_return Grapher) {

	//create empty graphs
	hash_graph_to_return = graph.NewHash(n)
	matrix_graph_to_return = graph.NewMatrix(n)

	//init randomness
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	//randomly insert edges
	for count := 0; count < n; {

		from := random.Intn(n)
		to := random.Intn(n)

		//only insert if edge doesn't exist
		if !hash_graph_to_return.HasEdge(from, to) && !matrix_graph_to_return.HasEdge(from, to) {

			//connect both graphs in the same way
			hash_graph_to_return.Add(from, to)
			matrix_graph_to_return.Add(from, to)

			count++
		}

	}

	return

}

/*
	This function runs deep-first-search for each component.
	Takes in Grapher object.
*/

func runDFS(g Grapher) {

	//boolean array indicating if a vertex has been visited by DFS or not
	//used to detect new components
	state := make([]bool, g.NumVertices())

	//take out vertex-index 'v' and boolean 'visited' for each vertex
	for v, visited := range state {
		//if not visited - then we've found a new component
		if visited == false {
			//run DFS starting from this vertex.
			graph.DFS(g, v, state, func(w int) {
				//do nothing each 'traverse-loop'
				//we'll just measure time in the main funciton
			})
		}

	}
}

/*
	This function will initialize a hash and a matrix that have n vertecies
	and n random edges. It will then print the number of components and size of largest component
	for matrix and hash graphs.
*/

func analyzeGraphInfo(n int) {

	//generate the graphs.
	//the vertices are connected in the same way in both
	hashGraph, matrixGraph := setupGraphs(n)

	//get graph info
	largest_component_size_hash, number_of_components_hash := getLargestSizeAndNumOfComponents(hashGraph)
	largest_component_size_matrix, number_of_components_matrix := getLargestSizeAndNumOfComponents(matrixGraph)

	//print graph info
	fmt.Println("Largest component size in Hash: ", largest_component_size_hash)
	fmt.Println("Number of components in Hash ", number_of_components_hash)
	fmt.Println("------------------")
	fmt.Println("Largest component size in Matrix: ", largest_component_size_matrix)
	fmt.Println("Number of components in Matrix ", number_of_components_matrix)
}

/*
	This function takes in a slice of ints that holds the graph sizes
	we will analyze the performance of.

	The function will print the time it took to run DFS 100 times for a hash and matrix graph.
*/

func analyzeGraphPerformance(graph_sizes []int) {

	const TEST_ITERATIONS = 100

	//go through all sizes
	for _, size := range graph_sizes {

		fmt.Println("TESTING GRAPH SIZE: ", size)

		hashGraph, matrixGraph := setupGraphs(size)

		//ANALYZE HASH
		//remember time
		before_hash := time.Now()

		for i := 0; i < TEST_ITERATIONS; i++ {
			runDFS(hashGraph)
		}

		time_taken_hash := time.Since(before_hash)

		fmt.Println("TIME FOR HASH: ", time_taken_hash)

		//ANALYZE MATRIX

		before_matrix := time.Now()

		for i := 0; i < TEST_ITERATIONS; i++ {
			runDFS(matrixGraph)
		}

		time_taken_matrix := time.Since(before_matrix)

		fmt.Println("TIME FOR MATRIX: ", time_taken_matrix)

		fmt.Println("------------------")
	}

}

func main() {

	//analyze number of components and size of largest component
	const GRAPH_SIZE_TO_ANALYZE = 1000
	analyzeGraphInfo(GRAPH_SIZE_TO_ANALYZE)

	fmt.Println("------------------")

	//analyze performace with different graph sizes
	graph_sizes_to_analyze := []int{3, 100, 500, 1000, 1500, 2000, 2500, 3000, 3500, 4000, 4500, 5000}
	analyzeGraphPerformance(graph_sizes_to_analyze)
}
