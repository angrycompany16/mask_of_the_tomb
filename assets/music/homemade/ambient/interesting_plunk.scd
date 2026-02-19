(
s = Server.local;
s.reboot();
)
(
SynthDef(\piano, {
	var snd, freq, n;
	freq = \freq.kr(440);
	n = (1..10);
	snd = SinOsc.ar(freq * n * (1 + (n * n * 0.1e-3)).sqrt, n.collect({ Rand(0, 2pi)}) );
	snd = snd * Env.perc(0, 10 / (n ** 5.0)).ar;
	snd = snd * (n * 1.3).cos;
	snd = snd.sum;
	snd = snd * -9.dbamp;
	snd = snd * Env.perc(0.01, 3.0).ar;
	snd = snd * Env.linen(0, \duration.kr(1.0), 0.1, curve: -4).ar(Done.freeSelf);
	snd = snd ! 2;
	snd = snd * -6.dbamp;
	Out.ar(\out.kr(0), snd);
}).add;
)

Synth(\piano);

(
var s;
s = Server.local;

Routine({
	s.bind( {Synth(\piano);} );
}).play;
)