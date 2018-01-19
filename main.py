import argparse
import os
import json
import plot
import subprocess
import string

parser=argparse.ArgumentParser()

parser.add_argument('-keep', action='store_true')
parser.add_argument('-plot', action='store_true')
parser.add_argument('-binary_path', action='store')
parser.add_argument('-commands', action='store')
parser.add_argument('-labels', action='store')
parser.add_argument('-files', action='store')

args=parser.parse_args()

# Commands
commands = [
	"dtcp -full_ring -steps=10000 -T=5",
	"dtcp -local_ring -steps=10000 -T=5",
] if args.commands is None else args.commands.split(",")
commands = [c.strip() for c in commands]

# Lables
labels = commands if args.labels is None else args.labels.split(",")
labels = [s.strip() for s in labels]

# Files
if args.files is None:
	files = [string.replace(s.strip(), ' ', '_') + ".txt" for s in commands]
else:
	files = args.files.split(",")
files = [f.strip() for f in files]

tmp_folder = 'tmp_files'
if os.path.isdir(tmp_folder):
	files = [os.path.join(tmp_folder, f.strip()) for f in files]
	print files

# Check length equality
if len(commands) != len(labels) or len(labels) != len(files):
	print("Error: Length of labels, files, commands need to be equal.")
	exit(0)

# Binary Path
if args.binary_path and args.binary_path[-1] != '/': args.binary_path += '/'
prefix = "" if args.binary_path is None else args.binary_path

# Run commands
for i in range(len(commands)):
	revised_command = prefix + commands[i] + " -file=" + files[i]
	print "running: " + revised_command
	print subprocess.check_output(revised_command.split())

# Get distributions
distributions = []
for f in files:
	with open(f) as json_data:
		d = json.load(json_data)
		distributions.append(d)

# Plot
if args.plot:
	plot.plot(distributions, labels)

# Delete the files if keep=True
if not args.keep:
	for f in files:
		try:
			os.remove(f)
		except OSError:
			pass