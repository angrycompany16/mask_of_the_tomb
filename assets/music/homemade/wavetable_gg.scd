(
s.reboot;
s = Server.local;
)

(
s.plotTree;
s.meter;
FreqScope.new;
)
s.scope;

(
~wavetables = (1..10).collect({|i|
	var numSegments = rrand(4, 6 + i), levels;
	levels = ({((rrand(0.0, 1.0)) * [1, (-1)]).choose}!(numSegments-1));
	Env(
		[0] ++ levels ++ [0],
		{exprand(2.0, 8.0)}!numSegments,
		{exprand(-0.1, (-8.0) / (i+1))}!numSegments,
	).asSignal(1024).asWavetable;
});

~bufs = Buffer.allocConsecutive(10, s, 2048, 1, {|buf, i|
	buf.setnMsg(0, ~wavetables[i]);
});
)

Buffer.freeAll;

(
SynthDef(\osc, {
	arg buf=0, amp=(-10), out=0, freq=440, detune=0.1, atk=0.01, sus=0.0, rel=1.0, curve=(-4);
	var snd, masterEnv, freqs;

	freqs = freq+LFNoise1.kr(0.1!8).bipolar(detune.midiratio);
	snd = Osc.ar(buf, freqs, {rrand(0.0, 2pi)}!8);
	snd = Splay.ar(snd);
	snd = Balance2.ar(snd[0], snd[1]);

	snd = LeakDC.ar(snd);
	masterEnv = EnvGen.ar(Env.linen(atk, sus, rel, curve:curve), doneAction: Done.freeSelf);
	snd = snd * masterEnv;

	snd = snd * amp.dbamp;
	Out.ar(out, snd);
}).add;
)

(
var bpm;
bpm = 60;
~clock = TempoClock.default();
~clock.tempo = bpm / 60;
)

// Let's make some cool patterns
(
Pdef(\perc_basic,
	Pbind(
		\instrument, \osc,
		\dur, Pseq(1/2 ! 4, inf),

		\buf, Prand(~bufs, inf),
		\amp, Pseq([6] ++ (-4!3), inf),
		\freq, Pexprand(36, 60).midicps,
		\detune, Pexprand(0.01, 0.4),
		\atk, 0.0,
		\sus, 0.0,
		\rel, Pexprand(0.02, 0.04),
		\curve, -8,
	)
);

Pdef(\melody,
	Pbind(
		\instrument, \osc,
		\scale, Scale.minor,
		\degree, Pseq([
			[-3,   -1, 1  ],
			[-4+7, -2, 0  ],
			[0,     2, 4-7],
			[-2,    0, 2  ],
		], inf),
		\dur, 4,

		\buf, Pxrand(~bufs[0..4], inf),
		\amp, -17,
		\detune, 0.4,
		\atk, Pexprand(0.3, 0.5),
		\sus, Pexprand(1.0, 1.5),
		\rel, Pexprand(1.5, 2.0),
		\curve, 0,
	);
);

Pdef(\bass,
	Pbind(
		\instrument, \osc,
		\dur, Pseq([1, 7], inf),
		\scale, Scale.minor,
		\octave, 2,
		\degree, Pseq([Pdup(2, Prand([0, 6, 12]))], inf),

		\buf, Pseq([Pdup(2, Prand(~bufs[0..6]))], inf),
		\amp, -12,
		\detune, Pexprand(0.1, 1.0),
		\atk, 0.05,
		\sus, Pseq([1.0, 5.0], inf),
		\rel, Pseq([0.5, 2.5], inf),
		\curve, -4,
	)
);

Pdef(\pad,
	Pbind(
		\instrument, \osc,
		\dur, 7,
		\scale, Scale.minor,
		\octave, 5,
		\degree, Pwhite(0.0, 12.0).round(2),

		\buf, Prand(~bufs[0..2], inf),
		\amp, -20,
		\detune, Pexprand(0.4, 1.0),
		\atk, Pexprand(3.0, 5.0),
		\sus, Pexprand(2.0, 4.0),
		\rel, Pexprand(3.0, 5.0),
		\curve, 0,
	)
);

Pdef(\perc_secondary,
	Pbind(
		\instrument, \osc,
		\dur, Pseq(((1/8 ! 10) ++ [Pwhite(1/4, 4, 1).round(1/4)]), inf),
		\degree, Pexprand(5, 12).round(),
		\scale, Scale.minor,
		\octave, 5,

		\buf, Pshuf(~bufs, inf),
		\amp, -8,
		\detune, 0.4,
		\atk, 0.0,
		\sus, 0.0,
		\rel, Pexprand(0.03, 0.09),
		\curve, -8,
	)
);
)

Pbindef(\perc_basic).play.quant(1);

(
Pbindef(\perc_secondary).play.quant(4);
Pbindef(\pad).play.quant(4);
)

Pbindef(\bass).play.quant(4);

Pbindef(\melody).play.quant(4);

Pbindef(\perc_basic).stop;
Pbindef(\melody).stop;
Pbindef(\bass).stop;
Pbindef(\perc_secondary).stop;
Pbindef(\pad).stop;










