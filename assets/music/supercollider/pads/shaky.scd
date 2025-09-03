(
s = Server.local;
s.reboot();
)

(
SynthDef(\pad, {
	var snd;

	snd = Saw.ar(\freq.kr(440) * ({ LFNoise2.kr(8) } ! 16 * 0.1).midiratio * [1, 2, 1, 2, 1, 2, 1, 2, 2, 1, 2, 1, 2, 1, 2, 1]);
	snd = DelayC.ar(snd, 0.01, { Rand(0, 0.01) } ! snd.size);
	snd = Splay.ar(snd);
	snd = LPF.ar(snd, \freq.kr(440) * XLine.kr(16, 4, 8));
	snd = BPeakEQ.ar(snd, \freq.kr(440) * 8, 0.1, 8);
	snd = snd * Env.asr(1, 1, 8).ar(Done.freeSelf, \gate.kr(1));
	snd = snd * \amp.kr(-30.dbamp);
	snd = snd ! 2;
	Out.ar(\out.kr(0), snd);
}).add;

SynthDef(\fx, {
	var snd;
	snd = In.ar(\out.kr(0), 2);
	snd = snd + NHHall.ar(snd, 3, modDepth: 1);
	snd = BHiShelf.ar(snd, 3000, 2, -10 * LFNoise2.kr(16));
	snd = BPeakEQ.ar(snd, 240, 0.1, -8 * LFNoise2.kr(16));
	snd = BPeakEQ.ar(snd, 460, 0.1, -8 * LFNoise2.kr(16));
	snd = BPeakEQ.ar(snd, 820, 0.1, -8 * LFNoise2.kr(16));
	snd = BLowShelf.ar(snd, 100, 1, 10 * LFNoise2.kr(16));
	snd = Limiter.ar(snd);
	Out.ar(\out.kr(0), snd);
}).add;
)

(
Routine {
	var frequencies;
	frequencies = (40 + [3, 7, 12, 14, 15, 17, 24, 26]).midicps;
	Synth.tail(nil, \fx);
	frequencies.scramble.do { |frequency|
		Routine {
			var synth;
			synth = Synth(\pad, [freq: frequency] );
			rrand(0.0, 4.0).wait;
			synth.set(\gate, 0);
		}.play;
		rrand(0.3, 1.0).wait;
	}
}.play;
)