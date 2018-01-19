import matplotlib.pyplot as plt

def shared_sorted_keys(distributions):
	for d1 in distributions:
		for d2 in distributions:
			for k in d1:
				if k not in d2: d2[k] = 0
	return sorted(distributions[0].keys())

def plot(distributions, labels, show, save):
	keys = shared_sorted_keys(distributions)
	for i in range(len(distributions)):
		d = distributions[i]
		plt.plot([d[k] for k in keys], label=labels[i])
	plt.legend(loc=2)
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)