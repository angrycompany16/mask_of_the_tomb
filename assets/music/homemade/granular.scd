(
s.reboot;
s = Server.local;
)

(
s.plotTree;
s.meter;
)

~b = Buffer.readChannel(s, "/home/angrycompany16/gamedev/go/mask_of_the_tomb/assets/music/homemade/samples/vietnamese_talking.wav", channels: [0]);
~b.free;

(
SynthDef(\gs, {
	arg buf = 0, sync = 0, freq = 1, amp = 1, duration = 1, durRand = 0, rate = 1, pos = 0, posSpeed = 0, posRand = 0.01, pan = 0, panRand = 0, panStrength = 0, env = (-1), atk = 0.01, sus = 1.0, rel=1.0, numChannels = 1;

	var snd, trig, dur, posCtrl, panCtrl;

	trig = Select.ar(sync, [Impulse.ar(freq), Dust.ar(freq)]);

	dur = duration + (LFNoise1.ar(freq).range(0.0002*duration, duration)*durRand);

	posCtrl = Phasor.ar(0, posSpeed * BufRateScale.ir(buf), 0, BufSampleRate.ir(buf));
	posCtrl = posCtrl + LFNoise1.kr(100).bipolar(posRand * SampleRate.ir(s));
	posCtrl = posCtrl / BufSamples.ir(buf);
	posCtrl = posCtrl + pos;

	panCtrl = (pan + (panStrength * LFNoise1.kr(panRand))).range(-1, 1);

	snd = GrainBuf.ar(
		2,
		trig,
		dur,
		buf,
		rate,
		posCtrl,
		2,
		panCtrl,
		env,
		512
	);

	snd = snd * EnvGen.ar(Env.linen(atk, sus, rel), doneAction: Done.freeSelf);

	snd = snd * amp.dbamp;
	Out.ar(\out.kr(0), snd);
}).add;
)

(
var bpm;
bpm = 120;
~clock = TempoClock.default();
~clock.tempo = bpm / 60;
)

(
~drums = Routine({
	var drumPatterns;
	drumPatterns = [
		[0.01, 0.01, 0.01, 0.01, 0.03, \, 0.01, 0.01],
		[0.01, \, 0.01, \, 0.01, \, 0.01, \],
		// [0.01, \, 0.01, \, 0.03, 0.03, 0.03, 0.03],
		// [0.01, \, \, 0.01, 0.01, \, 0.01, 0.01],
	];
	loop {
		drumPatterns.choose.do { |duration|
			if(duration != \) {
				Synth(\gs, [
					buf: ~b,
					sync: 0,
					freq: 1,
					amp: 30,
					duration: duration,
					durRand: 0.01,
					posRand: 0.04,
					pos: 0.01,
					atk: 0.001,
					sus: 0,
					rel: 0.3,
					env: ~grainEnv,
				]);
			};
			0.25.wait;
		}
	}
});


~chatter = Routine({
	loop {
		var dur = rrand(0.1, 0.6);

		Synth(\gs, [
			buf: ~b,
			sync: 0,
			freq: 1 / dur,
			amp: rrand(-4, 0),
			duration: dur,
			rate: rrand(0.8, 1.0),
			durRand: 0,
			posRand: 0,
			pos: rrand(0.0, 1.0),
			atk: 0.001,
			sus: dur,
			rel: 0.001,
			pan: 0,
			panRand: 1,
			panStrength: 3,
		]);
		[2.0, 4.0, 8.0].choose.wait;
	}
});

~chords = Routine({
	var chords = [
		[2, 5, 9], // D minor
		[9, 12, 16], // A minor
		[5, 9, 12],  // F major
		[7, 10, 14, 17], // D minor 7th
	];
	loop {
		chords.choose.do { |note|
			var freq = note.midicps + (12 * 7);
			Synth(\gs, [
				buf: ~b,
				sync: 0,
				freq: freq,
				rate: 1,
				amp: -10,
				duration: 0.02,
				durRand: 0,
				posRand: 0.0,
				pos: 0.6,
				atk: exprand(1.0, 2.0),
				sus: exprand(2.0, 4.0),
				rel: exprand(1.0, 3.0),
				pan: 0,
				panRand: 0.5,
				panStrength: 1,
			]);
		};
		[8.0, 12.0, 16.0].choose.wait;
	}
});
)

~drums.play(~clock, quant: 4);
~chatter.play(~clock, quant: 4);
~chords.play(~clock, quant: 4);

~drums.stop;
~chatter.stop;
~chords.stop;


(
e = Env.perc(0.001, 0.3);
~grainEnv = Buffer.sendCollection(s, e.discretize(1024));
)