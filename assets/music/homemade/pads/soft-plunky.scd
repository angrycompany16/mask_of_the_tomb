(
s = Server.local;
s.reboot();
)

(
SynthDef(\pad, {
	var snd;

	snd = Saw.ar(\freq.kr(440) * ({ LFNoise2.kr(8) } ! 16 * Line.kr(0, 0.5, 3)).midiratio * [1, 2, 1, 2, 1, 2, 1, 2, 2, 1, 2, 1, 2, 1, 2, 1]);
	snd = DelayC.ar(snd, 0.01, { Rand(0, 0.01) } ! snd.size);
	snd = Splay.ar(snd);
	snd = LPF.ar(snd, XLine.kr(1000, 100, 8));
	snd = BPeakEQ.ar(snd, \freq.kr(440) * 8, 0.1, 8);
	snd = snd * Env.asr(1, 1, 8).ar(Done.freeSelf, \gate.kr(1));
	snd = snd * \amp.kr(-30.dbamp);
	snd = snd ! 2;
	Out.ar(\out.kr(0), snd);
}).add;

SynthDef(\ping, {
	var snd;
	snd = SinOsc.ar(\freq.kr(440) * (LFNoise2.kr(0.3) * 0.1).midiratio);
	snd = snd * Env.perc(0.5, 3).ar(Done.freeSelf);
	snd = snd * \amp.kr(-20.dbamp);
	snd = snd ! 2;
	Out.ar(\out.kr(0), snd);
}).add;

SynthDef(\fx, {
	var snd;
	snd = In.ar(\out.kr(0), 2);
	snd = snd + PitchShift.ar(snd, 0.2, 0.3, 0, 0.1);
	snd = snd + NHHall.ar(snd, 3, modDepth: 1);
	snd = BHiShelf.ar(snd, 3000, 2, -10 + (5 * LFNoise2.kr(16 ! 2)));
	snd = BPeakEQ.ar(snd, 240, 0.1, -8 * LFNoise2.kr(16 ! 2));
	snd = BPeakEQ.ar(snd, 460, 0.1, -8 * LFNoise2.kr(16 ! 2));
	snd = BPeakEQ.ar(snd, 820, 0.1, -8 * LFNoise2.kr(16 ! 2));
	snd = BLowShelf.ar(snd, 100, 1, 10);
	snd = snd * -10.dbamp;
	snd = snd.tanh;
	snd = Limiter.ar(snd);
	ReplaceOut.ar(\out.kr(0), snd);
}).add;
)

(
Routine {
	var frequencies, chords;
	chords = [
		[3, 7, 12, 14, 15, 24, 26],
		[3, 7, 12, 14, 15, 24, 26],
		[3, 7, 12, 14, 15, 24, 26, 27, 29],
		[3, 7, 12, 14, 15, 24, 26, 27, 29],
		[0, 15, 24, 26, 27, 31],
	];

	Synth.tail(nil, \fx);

	Routine {
		loop {
			rrand(5.0, 8.0).wait;
			6.do {
				if (frequencies.notNil) {
					Synth(\ping,  [freq: frequencies.choose * 2, amp: -28.dbamp]);
				};


				exprand(0.1, 0.5).wait;
			};
		};
	}.play;

	loop {
		chords.do { |chord|
			frequencies = (39 + chord +.x [0, 12]).midicps;
			frequencies.scramble.do { |frequency|
				Routine {
					var synth;
					synth = Synth(\pad, [freq: frequency, amp: -20.dbamp] );
					rrand(1.0, 4.0).wait;
					synth.set(\gate, 0);
				}.play;
				rrand(0.3, 1.0).wait;
			};
			rrand(5.0, 10.0).wait;
		};
	};
}.play;
)