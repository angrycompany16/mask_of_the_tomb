(
s = Server.local;
s.options.memSize = 32768;
s.reboot();
)

(
SynthDef(\sine, {
	var freq = \freq.kr(440), atk = \atk.kr(1.0), sus = \sus.kr(1.0), rel = \rel.kr(1.0), amp = \amp.kr(1.0), out = \out.kr(0);
	var snd;

	snd = (1..10).collect({ |n|
		var n_snd;
		n_snd = Saw.ar(freq + n);
		n_snd = n_snd * Env.linen(atk * (n ** 2), sus, rel).ar(Done.freeSelf);

		n_snd = DelayN.ar(n_snd , 1 * (n ** -1.5), 1 * (n ** -1.5));
		n_snd = n_snd + (FreqShift.ar(n_snd, n * 0.2));
		n_snd;
	});

	snd.sum();

	snd = snd.tanh;
	snd = BLowShelf.ar(snd, 1000);
	snd = snd * -15.dbamp;
	snd = LPF.ar(snd, 2000);
	snd = snd ! 2;
	snd = snd * amp.dbamp;
	Out.ar(out, snd);
}).add;
)

(
var playNote;
var bpm, beat, wait;
playNote = { |note|
	Routine({
		s.bind({
			Synth(\sine, [freq: note.midicps, atk: 0.01, sus: 1.0, rel: 2.0, amp: -20]);
		});
	}).play;
};

bpm = 90;
beat = 60 / bpm;

wait = { |duration|
	(duration * beat).wait;
};

Routine({
	var chords;
	chords = [
		[68,          68 + 4,     68 + 4 + 3],
		[68 + 2,      68 + 2 + 4, 68 + 2 + 4 + 3 - 12],
		[68 - 5 + 12, 68 - 5 + 4, 68 - 5 + 4 + 3],
		[68 - 5 + 12, 68 - 5 + 3, 68 - 5 + 4 + 3],
	];
	chords.do { |chord|
		chord.do { |note|
			playNote.(note);
		};
		wait.(4.0);
	};
}).play;
)

