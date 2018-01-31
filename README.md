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