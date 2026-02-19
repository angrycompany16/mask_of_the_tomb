(
s = Server.local;
s.reboot();
)
(
SynthDef(\piano, {
	var snd, freq, n, knock;
	freq = \freq.kr(440);
	n = (1..15);
	snd = [-0.05, 0, 0.05].collect { |detune|
		var snd, freqs;

		freqs = freq * n * (1 + (n * n * 2e-3)).sqrt * detune.midiratio;
		snd = SinOsc.ar(freqs, n.collect({ Rand(0, 2pi)}) );
		snd = snd * Env.perc(snd.collect { Rand(0, 0.001) } * n, 20 * ((260 / freqs) ** 1.4)).ar;
		snd = snd * (1 + (0.5 * (n / Rand(3, 5)).sin));
		snd = snd * (n ** -1.5);
		snd = snd.sum;
		snd;
	};
	snd = snd.sum * -7.dbamp;
	snd = snd * -9.dbamp;
	snd = snd * (1 + (3 * Env.perc(0, 0.1).ar));
	snd = snd * Env.perc(0.01, 10.0).ar;

	knock = PinkNoise.ar;
	knock = BPF.ar(knock, freq, 0.3) * 15.dbamp;
	knock = knock * Env.perc(0.001, 0.01).ar;

	snd = snd + knock;
	snd = snd.tanh;

	snd = snd * -5.dbamp;
	snd = snd * Env.linen(0, \duration.kr(10.0), 0.1, curve: -4).ar(Done.freeSelf);
	snd = snd ! 2;
	snd = snd * -6.dbamp;
	Out.ar(\out.kr(0), snd);
}).add;

SynthDef(\pianoFx, {
	var snd;
	snd = In.ar(\in.kr(0), 2);
	snd = snd + CombC.ar(LPF.ar(snd * -10.dbamp, 3000), 0.2, 1 / [120, 132], 0.3);
	snd = snd + CombC.ar(LPF.ar(snd * -10.dbamp, 3000), 0.2, 1 / [95, 73], 0.3);
	snd = snd + CombC.ar(LPF.ar(snd * -10.dbamp, 3000), 0.2, 1 / [35, 46], 0.3);
	Out.ar(\out.kr(0), snd);
}).add;
)

Synth(\piano);

(
var s, bpm, beat;
s = Server.local;

bpm = 90;
beat = 60 / bpm;

Routine({
	var pianoBus;
	var play, playParallel, wait;
	var playPiano;

	s.bind {
		pianoBus = Bus.audio(nil, 2);
		Synth.tail(nil, \pianoFx, [in: pianoBus, out: 0])
	};

	wait = { |duration|
		(duration * beat).wait;
	};

	playParallel = { |synthDef, duration args = #[] |
		s.bind( {Synth(synthDef, [duration: duration * beat] ++ args)}; );
	};

	play = { |synthDef, duration args = #[] |
		playParallel.(synthDef, duration, args);
		wait.(duration);
	};

	playPiano = { |note|
		playParallel.(\piano, 10.0, [freq: (68 + note - 12).midicps, out: pianoBus]);
	};

	[-3, 0, 2, 4, 7, \, \, \, 5, 4, 0, \, \, -3, 0, 2, 4, 9, 7, \, \, 5, 4, 0].do { |note|
		if(note != \) {
			playPiano.(note);
		};

		wait.(1);
	}
}).play;
)