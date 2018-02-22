import matplotlib.pyplot as plt
from matplotlib import pylab
import numpy as np

def get_keys(distributions):
	for d1 in distributions:
		for d2 in distributions:
			for k in d1:
				if k not in d2: d2[k] = 0
	return sorted(distributions[0].keys())

def plot_discrete(distributions, labels, show, save):
	keys = get_keys(distributions)
	for i in range(len(distributions)):
		d = distributions[i]
		plt.plot([d[k] for k in keys], label=labels[i])
	plt.legend(loc=2)
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)


def plot_continuous(distributions, labels, show, save):
	for i in range(len(distributions)):
		d = distributions[i]
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k-1):
			print len(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']))
			print len(arr[:,j])
			plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), arr[:,j], label=(labels[i]+" P(X="+str(j))+")")
	plt.legend(loc=2)
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)

def plot_continuous2(distributions, labels, show, save):
	for i in range(len(distributions)):
		d = distributions[i]
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k):
			print len(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']))
			print len(arr[:,j])
			plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), arr[:,j], label=(labels[i]+" P(X="+str(j))+")")
	plt.legend(loc=2)
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)

def plot_guess(distributions, labels, show, save):
	for i in range(len(distributions)):
		d = distributions[i]
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k-1):
			print len(np.arange(0, len(arr[50:,j])*d['Dt'],d['Dt']))
			print len(arr[:,j])
			plt.plot(np.arange(50, 50+len(arr[50:,j])*d['Dt'],d['Dt']), arr[50:,j], label=(labels[i]+" P(X="+str(j))+")")
	
	f0 = np.array(distributions[1]['Distr'])
	f0 = f0[50:,0]
	f1 = np.array(distributions[2]['Distr'])
	f1 = f1[50:,0]
	f2 = np.array(distributions[3]['Distr'])
	f2 = f2[50:,0]

	d1 = abs(f1-f0)
	d2 = abs(f2-f1)
	k = (d2/d1)

	guess1 = f2 + d2*(1/(1-k))

	real = np.array(distributions[0]['Distr'])
	real = np.mean(real[100:,0])
	print real - guess1[-1]





	d = distributions[0]
	plt.plot(np.arange(50, 50+len(arr[50:,j])*d['Dt'],d['Dt']), guess1, label="guess" )
	plt.plot(np.arange(50, 50+len(arr[50:,j])*d['Dt'],d['Dt']), real*np.ones((150,1)), label="mean" )


	plt.legend(loc=2)
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)

'''
def plot_error(distributions, labels, show, save):
	truth = np.array(distributions[0]['Distr'])
	for i in range(1, len(distributions)):
		d = distributions[i]
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k-1):
			plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), abs(truth[:,j]-arr[:,j]), label=(" dist "+labels[i] + " " +labels[0]))
	plt.legend(loc=2)
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)
'''

def plot_error(distributions, labels, show, save):
	truth = np.array(distributions[0]['Distr'])
	for i in range(1, len(distributions)):
		d = distributions[i]
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k-1):
			plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), abs(truth[:,j]-arr[:,j]), label=(" dist "+labels[i] + " " +labels[0]))
	
	plt.legend(loc=0)
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)


def plot_relative_error(distributions, labels, show, save):
	truth = np.array(distributions[0]['Distr'])
	for i in range(1, len(distributions)):
		d = distributions[i]
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k-1):
			plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), abs(truth[:,j]-arr[:,j])/truth[:,j], label=(" relative error "+labels[i] + " " +labels[0]))
	
	for i in range(0, len(distributions)-1):
		d = distributions[i]
		truth = np.array(distributions[i+1]['Distr'])
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k-1):
			plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), abs(truth[:,j]-arr[:,j])/truth[:,j], label=(" relative error "+labels[i] + " " +labels[i+1]))
	
	plt.legend(loc=0)
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)



def plot_error_fraction(distributions, labels, show, save):
	'''
	truth = np.array(distributions[0]['Distr'])
	for i in range(1, len(distributions)):
		d = distributions[i]
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k-1):
			plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), truth[:,j]/arr[:,j], label=(" ratio "+labels[0] + ", " +labels[i]))
	'''
	truth = np.array(distributions[0]['Distr'])
	data1 = []
	for i in range(1, len(distributions)):
		d = distributions[i]
		#truth = np.array(distributions[i+1]['Distr'])
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k-1):
			#plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), abs(truth[:,j]/arr[:,j]), label=(" ratio "+labels[i+1] + " " +labels[i]))

			data1.append(abs(arr[-1,j]-truth[-1,j])/truth[-1,j])

	truth = np.array(distributions[-1]['Distr'])
	data0 = []
	for i in range(1, len(distributions)-1):
		d = distributions[i]
		#truth = np.array(distributions[i+1]['Distr'])
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k-1):
			#plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), abs(truth[:,j]/arr[:,j]), label=(" ratio "+labels[i+1] + " " +labels[i]))

			data0.append(abs(arr[-1,j]-truth[-1,j])/truth[-1,j])

	truth = np.array(distributions[-2]['Distr'])
	data3 = []
	for i in range(1, len(distributions)-2):
		d = distributions[i]
		#truth = np.array(distributions[i+1]['Distr'])
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k-1):
			#plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), abs(truth[:,j]/arr[:,j]), label=(" ratio "+labels[i+1] + " " +labels[i]))

			data3.append(abs(arr[-1,j]-truth[-1,j])/truth[-1,j])


	data2 = []
	for i in range(1, len(distributions)-1):
		d = distributions[i]
		truth = np.array(distributions[i+1]['Distr'])
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k-1):
			#plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), abs(truth[:,j]/arr[:,j]), label=(" ratio "+labels[i+1] + " " +labels[i]))

			data2.append(abs(arr[-1,j]-truth[-1,j]))
	
	''''
	plt.legend(loc=0)
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)
	'''
	print 'here'
	print (np.log(data2[-1]) - np.log(data2[0]))/ (len(data2)-1)
	data1 = (np.array(data1))
	data0 = (np.array(data0))
	data3 = (np.array(data3))
	data2 = (np.array(data2))
	#print len(data1)
	plt.subplot(211)
	plt.plot([0,1,2,3], data2, label="L1 error f_i and f_i+1")
	plt.xlabel("i")
	plt.legend()
	plt.subplot(212)
	plt.plot([0,1,2,3], np.log(data2), label="L1 error f_i and f_i+1")
	plt.xlabel("i")
	plt.legend()
	print np.log(data2)
	#plt.plot([0,1,2,3], data0)
	#plt.plot([0,1,2], data3)
	#plt.plot([0,1,2,3], data2)
	#plt.xlabel('i')
	#plt.ylabel('log of relative error between tau = i + 1 and tau = i')
	plt.show()


def plot_error_consecutive_pairs(distributions, labels, show, save):
	for i in range(0, len(distributions)-1):
		d = distributions[i]
		truth = np.array(distributions[i+1]['Distr'])
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k-1):
			plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), abs(truth[:,j]-arr[:,j]), label=(" dist "+labels[i] + " " +labels[i+1]))
	plt.legend(loc=2)
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)


def plot_error_all_pairs(distributions, labels, show, save):
	for l in range(0, len(distributions)):
		truth = np.array(distributions[l]['Distr'])
		for i in range(l+1, len(distributions)):
			d = distributions[i]
			k = d['K']
			arr = np.array(d['Distr'])
			for j in range(k-1):
				plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), abs(truth[:,j]-arr[:,j]), label=(" dist "+labels[i] + " " +labels[l]))
	plt.legend(loc=2)
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)
