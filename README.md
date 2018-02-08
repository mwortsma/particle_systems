# Particle Systems
A study of particle systems.

TODO:

General

- Fix rand.Rand bug (threadsafe?)

CTCP

- Fix bug in mean field fixed point (dipping ~ T = 3)

- Fix bug in python main.py -commands="dtcp -local_tree_realization -tau=2, dtcp -local_tree_realization -tau=0, dtcp -local_tree_realization -tau=1, dtcp -full_tree, dtcp -local_tree_fp -tau=1 -iters=3" -shared="-d=3 -T=4 -steps=100000" -show_plot


E.g.

python main.py -commands=" dtcp -full_tree, dtcp -local_tree_fp -tau=1 -iters=3, dtcp -local_tree_fp -tau=2 -iters=3, dtcp -local_tree_fp -tau=3" -shared="-d=3 -T=7 -steps=100000" -show_plot

E.g.

python main.py -commands=" dtcp -full_tree, dtcp -local_tree_realization -tau=1 -iters=3, dtcp -local_tree_realization -tau=2 -iters=3, dtcp -local_tree_realization -tau=3" -shared="-d=3 -T=7 -steps=100000" -show_plot

E.g.

python main.py -commands="dtlb -local_ring_fp -steps=100000 -tau=1, dtlb -local_ring_fp -steps=100000, dtlb -full_ring -n=60 -steps=100000" -shared="-T=3 -k=5" -show_plot


E.g.

python main.py -commands="dtlb -local_ring_realization -steps=10000 -tau=1, dtlb -local_ring_realization -steps=10000, dtlb -full_ring -n=60 -steps=100000" -shared="-T=3 -k=5" -show_plot

Most Recent
1. python main.py -commands="dtlb -local_ring_realization -tau=2, dtlb -local_ring_realization -tau=0, dtlb -local_ring_realization -tau=1, dtlb -full_tree, dtlb -local_ring_fp -tau=1" -shared="-d=3 -T=3 -steps=10000" -show_plot

2. python main.py -commands="dtcp -local_ring_realization -tau=2, dtcp -local_ring_realization -tau=0, dtcp -local_ring_realization -tau=1, dtcp -full_tree, dtcp -local_ring_fp -tau=1" -shared="-d=3 -T=3 -steps=100000" -show_plot

3. dtcp_fp_T7
python main.py -commands=" dtcp -full_tree -T=7, dtcp -local_tree_fp -tau=1 -iters=7 -T=7, dtcp -local_tree_fp -tau=2 -iters=7 -T=7, dtcp -local_tree_fp -tau=3 -iters=7 -T=7" -shared="-d=3 -steps=100000" -show_plot

4. dtcp_real_T7
python main.py -commands=" dtcp -full_tree -T=7, dtcp -local_tree_realization -tau=1  -T=7, dtcp -local_tree_realization -tau=2 -T=7, dtcp -local_tree_realization -tau=3 -T=7" -shared="-d=3 -steps=100000" -show_plot

3. dtcp_fp_T4
python main.py -commands=" dtcp -full_tree -T=7, dtcp -local_tree_fp -tau=1 -iters=4 -T=4, dtcp -local_tree_fp -tau=2 -iters=4 -T=4, dtcp -local_tree_fp -tau=3 -iters=4 -T=4" -shared="-d=3 -steps=100000" -show_plot

4. dtcp_real_T4
python main.py -commands=" dtcp -full_tree -T=4, dtcp -local_tree_realization -tau=1  -T=4, dtcp -local_tree_realization -tau=2 -T=4, dtcp -local_tree_realization -tau=3 -T=7" -shared="-d=3 -steps=100000" -show_plot