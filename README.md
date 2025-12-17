# min-cost-max-flow

This repository is only public for technical reasons.
Feel free to use it, however don't take it as maintained.

Despite the name, the project now contains the following algorithms:
- Min-cost-max-flow in directed, acyclic network graphs via successive shortest paths
- Minimax for simple symmetric deterministic games of perfect information
- Alpha-beta-pruning for the same type of games (WIP)
- Minimum Cost Bipartite Matching via the Min-Cost-Max-Flow implementation from above
- Graph traversal with BFS and DFS
- Shortest-Path using a variant of Belman-Ford-Moore which also minimizes the hop-distance using BFS
- Gauss-Elimination and Back Substitution using a maximum absolute value pivot rule
- Matrix inversion via Gauss-Elimination implementation
- Simplex algorithm using Bland's pivot rule on natural systems of the natural form max c@x s.t. A@x <= b, where a start basis is given
- Phase One Simplex to find a start basis of systems of the form named above (WIP)
- LP maximization and minimization of systems in natural form using the above implementation
- Simple algorithm to generate random permutations
- A generic binary search that uses Phase One from above to be used in the elipsoid method (WIP)
- Number theoretic algorithm that computes for any rational a fractional repesentation