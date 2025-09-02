(
s = Server.local;
s.reboot();
)

(
SynthDef(\pad, {
	var snd;
	snd = Saw.ar(\freq.kr(440));
	snd = snd * Env.asr(3, 1, 3).ar(Done.freeSelf, \gate.kr(1));
	snd = snd * \amp.kr(-10.dbamp);
	snd = snd ! 2;
	Out.ar(\out.kr(0), snd);
}).add;
)

(
Routine {
	var frequencies;
	frequencies = 50 * (1..16).nthPrime;
	loop {
		Routine {
			var synth;
			synth = Synth(\pad, [freq: frequencies.choose] );
			rrand(0.0, 16.0).wait;
			synth.set(\gate, 0);
		}.play;
		rrand(1.0, 4.0).wait;
	}
}.play;
)