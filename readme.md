# BenthosWorkflowV2 

### Problem 

The Benthos provided processor [Workflow](https://www.benthos.dev/docs/components/processors/workflow/) executes a [DAG](https://en.wikipedia.org/wiki/Directed_acyclic_graph) of Nodes, "performing them in parallel where possible".

However the current implementation uses this [dependency solver](https://github.com/quipo/dependencysolver) which is noted on the project readme as being incorrectly implemented. 
It takes the approach: resolve the DAG into series of steps where the steps are performed sequentially but the nodes in the step are performed in parallel.

This means that there can be situation where a step could be waiting for all the nodes in the previous step: even though all dependencies for the step are ready.

Consider the following DAG:

```
      /--> B -------------|--> D
     /                   /
A --|          /--> E --|
     \--> C --|          \
               \----------|--> F
```

The dependency solver would resolve the DAG into: ```[ [ A ], [ B, C ], [ E ], [ D, F ] ]```.
When we consider the node E, we can see the that full dependency of this node would be : ```A -> C -> E```, however in the stage before ```[ E ]```, there is the node B so in the current Benthos Workflow implementation *E would not execute until B even though there is no dependency of B for E*.

### Solution 

An alternative process is proposed here: 

Store the DAG in an [Adjacency Matrix](https://en.wikipedia.org/wiki/Adjacency_matrix):

For the DAG above this would be : 

|-|A|B|C|D|E|F|
|-|-|-|-|-|-|-|
|A|0|1|1|0|0|0|
|B|0|0|0|1|0|0|
|C|0|0|0|0|1|1|
|D|0|0|0|0|0|0|
|E|0|0|0|1|0|1|
|F|0|0|0|0|0|0|


From this Adjacency Matrix we can ascertain the following: 

  - The Node A has no dependencies because all of the values in Column A are 0.
  - The Nodes D and F are terminal Nodes because all of the values in Row D and F are 0.

Then execute the following psuedocode: 

```python
while !all_elements_are_zero(adjacency_matrix):
     for node in adjacency_matrix:
         if all_elements_are_zero(node_col):
               execute(node)
               update_row_to_all_zero(node_row)
```

To illustrate: 

|-|A|B|C|D|E|F|
|-|-|-|-|-|-|-|
|A|0|1|1|0|0|0|
|B|0|0|0|1|0|0|
|C|0|0|0|0|1|1|
|D|0|0|0|0|0|0|
|E|0|0|0|1|0|1|
|F|0|0|0|0|0|0|

A column is all 0 so then A is ready to be executed, all other steps are not ready.

A then finishes and the matrix is updated so that all values in row A are 0: 

|-|A|B|C|D|E|F|
|-|-|-|-|-|-|-|
|A|0|0|0|0|0|0|
|B|0|0|0|1|0|0|
|C|0|0|0|0|1|1|
|D|0|0|0|0|0|0|
|E|0|0|0|1|0|1|
|F|0|0|0|0|0|0|

So then the columns B and C are all 0 so they will be executed next. 
The process would repeat until the end condition of all elements were equal to 0.

There is some other details that have been ignored for the sake of simplicity in the illustration: 
     - Once a stage is finished do not execute it again 
     - The nodes are executed on a seperate thread

There is an implementation of the above in [./go_test/main.go](./go_test/main.go).
And will output: 

```
Node: id=A, started
Node: id=A, finished
Node: id=C, started
Node: id=B, started
Node: id=C, finished 
* here E starts before B is finished *
Node: id=E, started
Node: id=E, finished
Node: id=F, started
Node: id=F, finished
Node: id=B, finished
Node: id=D, started
Node: id=D, finished
```

