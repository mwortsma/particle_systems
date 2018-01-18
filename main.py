import argparse
import os
import subprocess

parser=argparse.ArgumentParser()

parser.add_argument('-keep', action='store_true')
parser.add_argument('-plot', action='store_true')
parser.add_argument('-binary_path', action='store')

args=parser.parse_args()

print(args)


commands = [
	"dtcp -full_ring -steps=10000",
	"dtcp -local_ring -steps=10000",
]

files = [
	"tmp_files/out1.txt",
	"tmp_files/out2.txt",
]

prefix = "" if args.binary_path is None else args.binary_path + "/"

for i in range(len(commands)):
	revised_command = prefix + commands[i] + " -file=" + files[i]
	print "running: " + revised_command
	print subprocess.check_output(revised_command.split())

# plot if args.plot
if args.plot:
	pass

# Delete the files if keep=True
if not args.keep:
	for f in files:
		try:
			os.remove(f)
		except OSError:
			pass