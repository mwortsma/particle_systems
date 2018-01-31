import argparse
import os
import json
import plot
import subprocess
import string

parser=argparse.ArgumentParser()

parser.add_argument('-keep', action='store_true')
parser.add_argument('-show_plot', action='store_true')
parser.add_argument('-save_plot', action='store')
parser.add_argument('-binary_path', action='store')
parser.add_argument('-commands', action='store')
parser.add_argument('-shared', action='store')
parser.add_argument('-labels', action='store')
parser.add_argument('-files', action='store')
parser.add_argument('-type', action='store')

args=parser.parse_args()

# Commands
commands = [] if args.commands is None else args.commands.split(",")
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

# Binary Path
if args.binary_path and args.binary_path[-1] != '/': args.binary_path += '/'
prefix = "" if args.binary_path is None else args.binary_path

# Run commands
for i in range(len(commands)):
	revised_command = "time " + prefix + commands[i] + " -file=" + files[i]
	if args.shared is not None:
		revised_command = revised_command + " " + args.shared
	print "running: " + revised_command
	print subprocess.check_output(revised_command.split())

# Get distributions
distributions = []
for f in files:
	with open(f) as json_data:
		d = json.load(json_data)
		distributions.append(d)

# Plot
if args.save_plot or args.show_plot:
	print args.type
	if args.type is None or args.type == 'discrete':
		plot.plot_discrete(distributions, labels, args.show_plot, args.save_plot)
	elif args.type == 'continuous':
		plot.plot_continuous(distributions, labels, args.show_plot, args.save_plot)


# Delete the files if keep=True
if not args.keep:
	for f in files:
		try:
			os.remove(f)
		except OSError:
			pass