
(
SynthDef(\pad, {
	var snd;
	snd = Saw.ar(\freq.kr(440) * ((-8..8).linlin(-8, 8, -1, 1) * 0.1).midiratio);
	snd = Splay.ar(snd);
	snd = snd * Env.asr(3, 1, 3).ar(Done.freeSelf, \gate.kr(1));
	snd = snd * \amp.kr(-30.dbamp);
	snd = snd ! 2;
	Out.ar(\out.kr(0), snd);
}).add;

SynthDef(\fx, {
	var snd;
	snd = In.ar(\out.kr(0), 2);
	snd = snd + NHHall.ar(snd, 3, modDepth: 1);
	snd = BHiShelf.ar(snd, 3000, 2, -10);
	snd = BLowShelf.ar(snd, 100, 1, 10);
	snd = Limiter.ar(snd);
	Out.ar(\out.kr(0), snd);
}).add;
)

(
Routine {
	var frequencies;
	frequencies = 50 * (1..8).nthPrime;
	Synth.tail(nil, \fx);
	loop {
		Routine {
			var synth;
			synth = Synth(\pad, [freq: frequencies.choose] );
			rrand(0.0, 16.0).wait;
			synth.set(\gate, 0);
		}.play;
		rrand(0.3, 1.0).wait;
	}
}.play;
)