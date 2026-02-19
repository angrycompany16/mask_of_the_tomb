(
s = Server.local;
s.options.memSize = 32768;
s.reboot();
)
(
SynthDef(\piano, {
	var snd, freq, n, knock;
	freq = \freq.kr(440);
	n = (1..15);
	snd = [-0.05, 0, 0.05].collect { |detune|
		var snd, freqs;

		freqs = freq * n * (1 + (n * n * 1e-3)).sqrt * detune.midiratio;
		snd = SinOsc.ar(freqs, n.collect({ Rand(0, 2pi)}) );
		snd = snd * Env.perc(snd.collect { Rand(0, 0.001) } * n, 20 * ((260 / freqs) ** 1.4)).ar;
		snd = snd * (1 + (0.5 * (n / Rand(3, 8)).sin));
		snd = snd * (n ** -1.5);
		snd = snd.sum;
		snd;
	};
	snd = snd.sum * -7.dbamp;
	snd = snd * -9.dbamp;
	snd = snd * (1 + Env.perc(0, 0.1).ar);
	snd = snd * Env.perc(0.01, 10.0).ar;

	knock = PinkNoise.ar;
	knock = BPF.ar(knock, freq * [1, 2.3], 0.7).sum * 10.dbamp;
	knock = knock * Env.perc(0.001, 0.06).ar;

	snd = snd + knock;
	snd = BHiShelf.ar(snd, 1230, 0.3, -3);
	snd = (snd * 2.dbamp).tanh;
	snd = BHiShelf.ar(snd, 3230, 0.3, 5);

	snd = snd * -7.dbamp;
	snd = snd * Env.linen(0, \duration.kr(10.0), 0.1, curve: -4).ar(Done.freeSelf);
	snd = snd ! 2;
	snd = snd * -6.dbamp;
	Out.ar(\out.kr(0), snd);
}).add;

SynthDef(\pianoFx, {
	var snd;
	snd = In.ar(\in.kr(0), 2);
	snd = snd + CombC.ar(LPF.ar(snd * -10.dbamp, 3000), 0.2, 1 / 120, 0.3);
	snd = snd + CombC.ar(LPF.ar(snd * -10.dbamp, 3000), 0.2, 1 / 95, 0.3);
	snd = snd + CombC.ar(LPF.ar(snd * -10.dbamp, 3000), 0.2, 1 / 35, 0.3);
	Out.ar(\out.kr(0), snd);
}).add;

SynthDef(\pad, {
	var snd, freq, env;
	env = Env([0, 1, 0], [1, 1].normalizeSum * \duration.kr(1), \sin).ar(Done.freeSelf);
	freq = \freq.kr(440);
	snd = Saw.ar((freq * (LFNoise2.kr(3 ! 3) * 0.1).midiratio));
	snd = Splay.ar(snd, 0.5);
	snd = RLPF.ar([snd], env.linexp(0, 1, 100, 700) * [1, 1.3], 0.5).sum * -5.dbamp;
	snd = snd.tanh;
	snd = snd + FreqShift.ar(snd, 0.5);
	snd = BPF.ar([snd], (1..30).normalize.linexp(0, 1, 100, 7000), 0.3);
	snd = snd + CombC.ar(snd * -10.dbamp, 0.2, 1 / 70, 1.0);
	snd = snd * env;
	snd = snd * -28.dbamp;
	Out.ar(\out.kr(0), snd);
}).add;

SynthDef(\padHigh, {
	var snd, freq, env;
	env = Env([0, 1, 0], [1, 1].normalizeSum * \duration.kr(1), \sin).ar(Done.freeSelf);
	freq = \freq.kr(440);
	snd = Pulse.ar((freq * (LFNoise2.kr(3 ! 3) * 0.1).midiratio) + 10);
	snd = Splay.ar(snd, 0.5);
	snd = FreqShift.ar(snd, -10);
	snd = (BPF.ar([snd], [73, 1440, 3200], 0.3) * [0, -5, -10].dbamp).sum;
	snd = LPF.ar(snd, 1200);
	snd = snd + DelayC.ar(snd, 0.2, SinOsc.ar(0.3, [0, pi]).linlin(-1, 1, 0, 1) * 1e-3);
	snd = snd * env;
	snd = snd * -23.dbamp;
	Out.ar(\out.kr(0), snd);
}).add;

SynthDef(\padLow, {
	var snd, freq, env;
	env = Env([0, 1, 0], [1, 1].normalizeSum * \duration.kr(1), \sin).ar(Done.freeSelf);
	freq = \freq.kr(440);
	snd = LFTri.ar(freq * (LFNoise2.kr(3) * 0.1).midiratio);
	snd = snd ! 2;
	snd = snd * env;
	snd = snd * -12.dbamp;
	Out.ar(\out.kr(0), snd);
}).add;

SynthDef(\ping, {
	var snd;
	snd = SinOsc.ar(\freq.kr(440));
	snd = snd * Env.perc(0.001, 1.0).ar;
	snd = snd * (1 + Env.perc(0.0, 0.1).ar);
	snd = LPF.ar(snd, 3000);
	snd = Latch.ar(snd, Impulse.ar(3000));
	snd = snd + GVerb.ar(DelayC.ar(snd, 0.2, LFNoise2.ar(3).linlin(-1, 1, 0, 1) * 5e-3) * -5.dbamp, 30, 10);
	snd = snd * Env.perc(0.0, 15.0).ar(Done.freeSelf);
	snd = snd * -35.dbamp;
	Out.ar(\out.kr(0), snd);
}).add;
)

Synth(\ping, [freq: (68 + 12).midicps]);


(
var s, bpm, beat, noteToFreq;
s = Server.local;

bpm = 72;
beat = 60 / bpm;

noteToFreq = { |note|
	var scale;
	// C major
	scale = [0, 2, 4, 5, 7, 9, 11];
	(60 - 4 + scale.wrapAt(note) + (12 * (note.div(scale.size)))).midicps;
};

Routine({
	var pianoBus;
	var play, playParallel, wait;
	var playPiano;
	var playChord;

	s.bind {
		pianoBus = Bus.audio(nil, 2);
		Synth.tail(nil, \pianoFx, [in: pianoBus, out: 0])
	};

	wait = { |duration|
		(duration * beat).wait;
	};

	playParallel = { |synthDef, duration, args = #[] |
		Routine({
			(1.0.rand.squared * 0.1).wait;
			s.bind( {Synth(synthDef, [duration: duration * beat] ++ args)}; );
		}).play;
	};

	play = { |synthDef, duration args = #[] |
		playParallel.(synthDef, duration, args);
		wait.(duration);
	};

	playPiano = { |note|
		playParallel.(\piano, 10.0, [freq: noteToFreq.(note), out: pianoBus]);
	};

	// Play some chords
	Routine({
		var oct, notes, root;
		oct = 5;
		root = 4;
		// major
		// notes = [12 * oct + root, 12 * oct + root + 4, 12 * oct + root + 4 + 3];
		// minor
		// notes = [12 * oct + root, 12 * oct + root + 3, 12 * oct + root + 4 + 3];

		notes.do { |note|
 			Routine({
				s.bind( {Synth(\piano, [duration: 10.0 * beat, freq: note.midicps])}; );
			}).play;
		};
	}).play;

	// Routine({
	// 	loop {
	// 		[0, 5, 4, 3, 2, -3, -2].do { |note|
	// 			playParallel.(\padHigh, 20.0, [freq: noteToFreq.(7 + note)]);
	// 			wait.(10.0);
	// 		};
	// 	};
	// }).play;
	//
	// Routine({
	// 	loop {
	// 		[2, -4, 0, 4].do{ |note|
	// 			playParallel.(\pad, 20.0, [freq: noteToFreq.(note)]);
	// 			wait.(17);
	// 		};
	// 		wait.(17);
	// 	};
	// }).play;
	//
	// Routine({
	// 	loop {
	// 		[3, 0, 5].do{ |note|
	// 			play.(\padLow, 50, [freq: noteToFreq.((-7  * 2) + note)]);
	// 			wait.(10);
	// 		};
	// 	};
	// }).play;
	//
	// Routine({
	// 	loop {
	// 		wait.(23);
	// 		play.(\ping, 67, [freq: noteToFreq.(7 * 2)]);
	// 	};
	// }).play;
	//
	// loop {
	// 	var fragments;
	//
	// 	fragments = [
	// 		[-2, 0, 1, 2, 4],
	// 		[3, 2, 0],
	// 		[-2, 0, 1, 2, 5, 4],
	// 		[3, 2, 0],
	// 	];
	//
	// 	[1, 2, 3, 4, 5, 6, 5, 4, 5, 6].do { |cutoff|
	// 		fragments.do { |fragment|
	// 			fragment[..cutoff].do { |note|
	// 				if(note != \) {
	// 					playPiano.(note);
	// 				};
	// 				wait.(1);
	// 			};
	//
	// 			wait.(5);
	// 		};
	// 		wait.(4);
	// 	};
	// };
}).play;
)