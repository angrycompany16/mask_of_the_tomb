// Boot
(
s = Server.local;
s.options.memSize = 32768;
s.reboot();
)

// Synths
(
SynthDef(\lead, {
	var snd, pop, noise, freq, atk, sus, rel;
	// Generate soft (main) sound
	freq = \freq.kr(440);
	atk = \atk.kr(0.01);
	sus = \sus.kr(0.2);
	rel = \rel.kr(1.0);

	snd = [-10, -3, -1, 0, 1, 3, 10].collect({ |n|
		var n_snd, n_freq;
		n_freq = freq + ((0.7 * n.midiratio) ** n);
		n_snd = SinOsc.ar(n_freq);
		n_snd = n_snd * Env.linen.kr(atk / (abs(n) + 1), sus, rel * abs(n));
		n_snd = (n_snd ** 1.3);
		n_snd = n_snd * (-1 * abs(n)).dbamp;
		n_snd;
	});

	snd = snd.sum;
	snd = BPF.ar(snd, freq, 3.0);

	// Add a small pop
	pop = BrownNoise.ar();
	pop = pop * Env.perc(0.001, 0.06).ar();
	pop = BPF.ar(pop, freq, 0.3);
	pop = pop * 10.dbamp;
	snd = snd + pop;

	snd = (snd).tanh();
	snd = LPF.ar(snd, 2000);

	// Master + Output
	snd = snd * EnvGen.ar(Env.linen(atk, sus, rel), doneAction: Done.freeSelf);
	snd = snd * \amp.kr(-20).dbamp;
	snd = snd ! 2;
	Out.ar(\out.kr(0), snd)
}).add;

SynthDef(\bass, {
	var snd, freq, noise, atk, rel;
	freq = \freq.kr(440);
	atk = \atk.kr(0.01);
	rel = \rel.kr(1.0);
	snd = Pulse.ar(freq, 0.5);
	snd = snd + (Pulse.ar(freq, 0.2) * Env.perc.ar(atk, rel / 2));
	snd = snd * Env.perc.ar(atk, rel);

	noise = BrownNoise.ar();
	noise = BPF.ar(noise, freq);
	noise = noise * Env.perc.ar(atk, rel * 3);
	snd = snd + (7.dbamp * noise);

	snd = LPF.ar(snd, 500);
	// snd = snd * Env.perc.ar(atk, rel);

	snd = snd * \amp.kr(-2).dbamp;
	snd = snd ! 2;
	Out.ar(\out.kr(0), snd);
}).add;
)

// Synth testing ground
Synth(\bass, [freq: 12.midicps, atk: 1.0, rel: 10.0]);
Synth(\lead);

// Composition
(
~bpm = 90;
~clock = TempoClock.default;
~clock.tempo = ~bpm / 60;

~noteToFreq = { |note|
	var scale;
	// C major
	scale = [0, 2, 3, 5, 7, 8, 10];
	(60 - 4 + scale.wrapAt(note) + (12 * (note.div(scale.size)))).midicps;
};

~play = { | synthDef, args = #[] |
	s.bind( { Synth(synthDef, args) }; )
};
)

(
~chords = Routine({
	var chords;

	chords = [
		[0, 2, 4], // C major chord
		[5, 7, 9], // A minor chord
		[4, 6, 8], // G major chord
		[3, 5, 7], // E minor chord
	];

	chords.do { |chord|
		Routine({
			chord.do { |note|
				0.1.wait;
				~play.(\lead, [freq: ~noteToFreq.(note), atk: 0.05, sus: 3.0 * TempoClock.default.beatDur, rel: 2.0 * TempoClock.default.beatDur]);
			};
		}).play;
		4.wait;
	};
});


~bass = Routine({
	var notes;

	notes = [-28, -28+ 5, -28+ 4, -28+ 3];

	notes.do { |note|
		~play.(\bass, [freq: ~noteToFreq.(note)]);
		4.wait;
	};
});
)

~chords.play(~clock, quant: 4);
~chords.stop;

~bass.play(~clock, quant: 4);
~bass.stop;







