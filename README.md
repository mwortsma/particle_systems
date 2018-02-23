# Particle Systems
A study of particle systems.

TODO:

General

- Fix rand.Rand bug (threadsafe?)



Errors:
python main.py -commands="cmd -dtsir_rec -tau=0, cmd -dtsir_rec -tau=1, cmd -dtsir_full_continuous -steps=100000, cmd -dtsir_full_continuous -steps=1000000" -shared="-T=100  -d=2 -q=0.1 -p=0.8" -show_plot -type=continuous2 -labels="Tau=0 Approx, Tau=1 Approx, Full, Full2"

Error Tau=0 Approx Tau=1 Approx:
j : 3.45279360658e-14
j : 4.41092068469e-15
j : 3.2852193188e-14
Error Tau=0 Approx Full:
j : 0.0358612958319
j : 0.0364196063854
j : 0.0446630339432
Error Tau=0 Approx Full2:
j : 0.0157920521057
j : 0.0123702108794
j : 0.0234832426466
Error Tau=1 Approx Full:
j : 0.0358612958319
j : 0.0364196063854
j : 0.0446630339432
Error Tau=1 Approx Full2:
j : 0.0157920521057
j : 0.0123702108794
j : 0.0234832426467
Error Full Full2:
j : 0.029409000011
j : 0.0352520000097
j : 0.0463429997751
