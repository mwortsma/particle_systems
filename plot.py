import matplotlib.pyplot as plt

def plot(distributions, labels):
	for i in range(len(distributions)):
		d = distributions[i]
		plt.plot([d[k] for k in sorted(d.iterkeys())], label=labels[i])
	plt.legend(loc=2)
	plt.show()